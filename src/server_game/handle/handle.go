package handle

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"log"
	"reflect"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server_game/item"
	"github.com/xujintao/balgass/src/server_game/model"
	"github.com/xujintao/balgass/src/server_game/object"
)

func init() {
	APIHandleDefault.init(apiIns[:], apiOuts[:])
	APIHandleDefault.start()
}

// APIHandleDefault default api handle
var APIHandleDefault apiHandle

type apiHandle struct {
	apiIns      map[int]*apiIn
	apiOuts     map[interface{}]*apiOut
	onConnChan  chan context.Context
	onCloseChan chan context.Context
	apiChan     chan context.Context
}

func (h *apiHandle) init(apiIns []*apiIn, apiOuts []*apiOut) {
	// ingress
	h.apiIns = make(map[int]*apiIn)
	for _, v := range apiIns {
		if vv, ok := h.apiIns[v.code]; ok {
			log.Printf("duplicated api code[%d] name[%s] with code[%d] name[%s]", v.code, v.name, vv.code, vv.name)
		}
		h.apiIns[v.code] = v
	}
	// egress
	h.apiOuts = make(map[interface{}]*apiOut)
	for _, v := range apiOuts {
		t := reflect.TypeOf(v.msg)
		if t.Kind() != reflect.Ptr {
			log.Printf("api code[%d] name[%s] msg field must be a pointer", v.code, v.name)
		}
		h.apiOuts[t] = v
	}

	h.onConnChan = make(chan context.Context, 100)
	h.onCloseChan = make(chan context.Context, 100)
	h.apiChan = make(chan context.Context, 1000)
}

func (h *apiHandle) start() {
	go func() {
		for {
			select {
			case ctx := <-h.onConnChan:
				object.ObjectManager.AddPlayer(ctx)
			case ctx := <-h.onCloseChan:
				object.ObjectManager.DeletePlayer(ctx)
			case ctx := <-h.apiChan:
				api := ctx.Value(nil).(*apiIn)
				obj := ctx.Value(nil)
				handle := reflect.ValueOf(api.handle)
				handle.Call([]reflect.Value{reflect.ValueOf(obj), reflect.ValueOf(api.msg)})
			}
		}
	}()
}

// Handle *apiHandle implements c1c2.Handler
func (h *apiHandle) Handle(ctx context.Context, req *c1c2.Request) {
	var api *apiIn
	var ok bool
	code := int(req.Body[0])
	if api, ok = h.apiIns[code]; !ok {
		codes := []byte{req.Body[0], req.Body[1]}
		code = int(binary.BigEndian.Uint16(codes))
		if api, ok = h.apiIns[code]; !ok {
			log.Printf("invalid api, body: %v", req.Body)
			return
		}
		req.Body = req.Body[1:]
	}
	req.Body = req.Body[1:]

	// validate encrypt
	if api.enc && !req.Encrypt {
		log.Printf("[%d][%s] not encrypt", api.code, api.name)
		// return
	}

	// authenticate/authorize
	// if game.GetAuthLevel(ctx) < int(api.level) {
	// 	log.Printf("[%d][%s] not authrized", api.code, api.name)
	// 	return
	// }

	// cbi.Unmarshal(req.Body, api.msg)
	// api.handle(obj, api.msg) // reflect call
	ctx = context.WithValue(ctx, "msg", api.msg)
	h.apiChan <- ctx
}

func (h *apiHandle) Push(w c1c2.ConnWriter, msg interface{}) {
	t := reflect.TypeOf(msg)
	api, ok := h.apiOuts[t]
	if !ok {
		log.Printf("%s has not yet be registered to api table", t.String())
		return
	}
	data, _ := json.Marshal(msg)
	var buf bytes.Buffer
	if api.subcode {
		var codes [2]uint8
		binary.BigEndian.PutUint16(codes[:], uint16(api.code))
		buf.Write(codes[:])
	} else {
		buf.WriteByte(uint8(api.code))
	}
	buf.Write(data)
	res := &c1c2.Response{}
	res.WriteHead(uint8(api.flag)).Write(buf.Bytes())
	w.Write(res)
}

// OnConn implements c1c2.Handler.OnConn
func (h *apiHandle) OnConn(ctx context.Context) {
	// msg := model.MsgConnectResult{}
	// ctx, err = game.OnConn(addr, conn, h)
	// if err != nil {
	// 	msg.Result = 0
	// } else {
	// 	msg.Result = ctx.(int)
	// }
	// h.Push(conn, &msg)
	// return
	h.onConnChan <- ctx
}

// OnClose implements c1c2.Handler.OnConn
func (h *apiHandle) OnClose(ctx context.Context) {
	h.onCloseChan <- ctx
}

type AuthLevel int

const (
	Guest AuthLevel = iota
	Player
	GM
	Admin
)

type apiIn struct {
	id     int
	enc    bool
	level  AuthLevel
	name   string
	code   int
	handle interface{}
	msg    interface{}
}

type apiOut struct {
	id      int
	enc     bool
	name    string
	flag    int // 0xC1 or 0xC2
	code    int
	subcode bool
	msg     interface{}
}

var apiIns = [...]*apiIn{
	{0, false, Player, "in_use_item", 0x26, useItem, (*model.MsgUseItem)(nil)},
	{0, false, Player, "in_learn_master_skill", 0xF352, learnMasterSkill, (*model.MsgLearnMasterSkill)(nil)},
	// {0, false, Guest, "in_login", 0xF101, game.Login, game.Login, game.SetAuthLevel},
}

var apiOuts = [...]*apiOut{
	{0, false, "out_connect_result", 0xC1, 0xF100, true, (*model.MsgConnectResult)(nil)},
	{0, false, "out_skill_list", 0xC1, 0xF311, true, (*model.MsgSkillList)(nil)},
}

func useItem(player *object.Player, msg *model.MsgUseItem) {
	// validate the position

	it := &player.Inventory[msg.InventoryPos]
	it2 := &player.Inventory[msg.InventoryPosTarget]
	// validate item serial/id
	if player.LimitUseItem(it) {
		return
	}
	switch {
	case it.Code == item.Code(12, 30): // Bundle Jewel of Bless
	case it.Code == item.Code(13, 15): // Fruits 果实
	case it.Code >= item.Code(13, 43) && it.Code <= item.Code(13, 45): // Seal 印章
	case it.Code == item.Code(13, 48): // Kalima Ticket 卡利玛自由入场券
	case it.Code >= item.Code(13, 54) && it.Code <= item.Code(13, 58): // Reset Fruit 洗点果实
	case it.Code == item.Code(13, 60): // Indulgence 免罪符
	case it.Code == item.Code(13, 66): // Invitation to Santa Village 圣诞之地入场券
	case it.Code == item.Code(13, 69): // Talisman of Resurrection 复活符咒
	case it.Code == item.Code(13, 70): // Talisman of Mobility 移动符咒
	case it.Code == item.Code(13, 82): // Talisman of Item Protection 装备保护符咒
	case it.Code >= item.Code(13, 152) && it.Code <= item.Code(13, 159): // Scroll of Oblivion 忘却卷轴
	case it.Code >= item.Code(14, 0) && it.Code <= item.Code(14, 3): // HP Potion
	case it.Code >= item.Code(14, 4) && it.Code <= item.Code(14, 6): // MP Potion
	case it.Code == item.Code(14, 7): // Siege Potion 攻城药水
	case it.Code == item.Code(14, 8): // Antidote 解毒剂
	case it.Code == item.Code(14, 9) || it.Code == item.Code(14, 20): // Ale 酒 / Remedy of Love 爱情的魔力
	case it.Code == item.Code(14, 10): // Town Portal Scroll 回城卷轴
	case it.Code == item.Code(14, 13): // Jewel of Bless
	case it.Code == item.Code(14, 14): // Jewel of Soul
	case it.Code == item.Code(14, 16): // Jewel of Life
	case it.Code >= item.Code(14, 38) && it.Code <= item.Code(14, 40): // comples/compound Potion
	case it.Code >= item.Code(14, 35) && it.Code <= item.Code(14, 37): // SD Potion
	case it.Code == item.Code(14, 42) && it2.Type != item.TypeSocket: // 再生强化
	case it.Code >= item.Code(14, 43) && it.Code <= item.Code(14, 44): // 进化道具
	case it.Code >= item.Code(14, 46) && it.Code <= item.Code(14, 50): // Jack O'Lantern 南瓜灯饮料
	case it.Code == item.Code(14, 70): // Elite HP Potion 精华HP药水
	case it.Code == item.Code(14, 71): // Elite MP Potion 精华MP药水
	case it.Code >= item.Code(14, 78) && it.Code <= item.Code(14, 82): // kindBPremiumElixir 会员圣水
	case it.Code >= item.Code(14, 85) && it.Code <= item.Code(14, 87): // Cherry Blossom 樱花
	case it.Code == item.Code(14, 133): // Elite SD Potion 精华防护值药水
	case it.Code == item.Code(14, 160): // Jewel of Extension 延长宝石
	case it.Code == item.Code(14, 161): // Jewel of Elevation 提高宝石
	case it.Code == item.Code(14, 162): // Magic Backpack 魔法背书
	case it.Code == item.Code(14, 163): // Vault Expansion Certificate 仓库拓展证书
	case it.Code == item.Code(14, 209): // Tradeable Seal 交易印章
	case it.Code == item.Code(14, 224): // Bless of Light (Greater) 光的祝福
	case it.Code >= item.Code(14, 263) && it.Code <= item.Code(14, 264): // Bless of Light 光之祝福
	case it.KindA == item.KindASkill:
		// (15, 18) // Scroll of Nova 星辰一怒术
		skillIndex := it.SkillIndex
		if it.Code == item.Code(12, 11) { // Orb of Summoning 召唤之石
			skillIndex += it.Level
		}
		if player.SkillLearn(skillIndex) {
			player.PushSkillOne()
		}

	}
}

func learnMasterSkill(player *object.Player, msg *model.MsgLearnMasterSkill) {
	if player.SkillLearn(msg.SkillIndex) {

	}
}

/*
var apis = map[int]func(req *c1c2.Request, res *c1c2.Response) bool{
	0x00:   chatProc,
	0x01:   chatGet,
	0x02:   chatWhisperGet,
	0x03:   mainCheck,
	0x0E:   connAlive,
	0x11:   attack, // s9
	0x15:   positionSet,
	0x18:   actionGet,
	0x19:   magicAttack,
	0x1B:   magicCancel,
	0x1C:   teleportGet,
	0x1D:   attackedGet, // 1.04R
	0x1E:   durationMagicGet,
	0x22:   itemGet,
	0x23:   itemDrop,
	0x24:   inventoryItemMove,
	0x26:   itemUse,
	0x30:   talk,
	0x31:   windowClose,
	0x32:   buy,
	0x33:   sell,
	0x34:   itemModify,
	0x36:   tradeReq,
	0x37:   tradeReqResult,
	0x3A:   tradeMoney,
	0x3C:   tradeConfirm,
	0x3D:   tradeCancel,
	0x3F01: shopItemSetPrice,
	0x3F02: shopOpen,
	0x3F03: shopClose,
	0x3F05: shopItemList,
	0x3F06: shopItemBuy,
	0x3F07: shopDealClose,
	0x40:   partyReq,
	0x41:   partyReqResult,
	0x42:   partyList,
	0x43:   partDelMember,
	0x4A:   rageAttack,
	0x4B:   rageAttackRange,
	0x4C00: mineTwinkle,
	0x4C01: mineTwinkleReward,
	0x4C03: mineTwinkleFail,
	0x4D00: eventItemGet,
	0x4D01: eventItemDrop,
	0x4D0F: eventInventoryOPen,
	0x4D10: rummyStart,
	0x4D11: cardReveal,
	0x4D12: cardMove,
	0x4D13: cardReMove,
	0x4D14: cardMatch,
	0x4D15: rummyEnd,
	0x4E00: muunItemGet,
	0x4E08: muunInventoryItemUse,
	0x4E09: muunItemSell,
	0x4E11: rideSelect,
	0x4E13: muunItemExchange,
	0x4F02: gremoryCaseItemGet,
	0x4F05: gremoryCaseOPen,
	0x50:   guildReq,
	0x51:   guildReqResult,
	0x52:   guildList,
	0x53:   guildDelMember,
	0x54:   guildMasterAnswer,
	0x55:   guildMasterInfoSave,
	0x57:   guildMasterCreateCancel,
	0x61:   guildWarReq,
	0x66:   guildViewportInfo,
	0x6F00: itemSoldList,
	0x6F01: itemSoldCancelSale,
	0x6F02: itemsoldReBuy,
	0x71:   pingSend,
	0x72:   packetCheckSum,
	0x81:   warehouseMoneyInOut,
	0x82:   warehouseUseEnd,
	0x83:   warehousePassword,
	0x86:   chaosBoxItemMixButtonOK,
	0x87:   chaosBoxUseEnd,
	0x88:   checkMultiMix,
	0x8E02: mapMove,
	0x90:   moveDevilSquare,
	0x91:   devilSquareRemainTime,
	0x95:   registerEventChip,
	0x96:   getMutoNum,
	0x97:   endEventChip,
	0x98:   useRenaChangeZen,
	0x99:   moveToOtherServer,
	0x9A:   reqEnterBloodCastle,
	0x9D:   reqLottoRegister,
	0x9F:   reqEventEnterCount,
	0xA0:   getQuestInfo,
	0xA2:   setQusetState,
	0xA7:   getPetItemCommand,
	0xA9:   getPetItemInfo,
	0xAA01: pkInvite,
	0xAA02: pkInviteResult,
	0xAA03: pkLeave,
	0xAA07: pkJoinChannel,
	0xAA09: pkLeaveChannel,
	0xAE:   reqMuBotSaveData,
	0xAF01: reqEnterChaosCastle,
	0xAF02: reqRepositionUserInChaoCastle,
	0xAF03: reqCCFDayTime,
	0xAF05: reqCCFEnterCheck,
	0xAF06: reqRepositionUserInCCF,
	0xAF07: reqCCFRanking,
	0xAF08: reqCCFUIOnOff,
	0xB0:   targetTeleportGet,
	0xB101: mapServerAuth,
	0xB200: castleSiegeState,
	0xB201: castleSiegeReg,
	0xB202: castleSiegeGiveUp,
	0xB203: guildRegInfo,
	0xB204: guildRegMark,
	0xB205: npcBuy,
	0xB206: npcRepair,
	0xB207: npcUpgrade,
	0xB208: taxMoneyInfo,
	0xB209: taxRateChange,
	0xB210: moneyDrawOut,
	0xB212: csGateOperate,
	0xB21B: csMiniMapData,
	0xB21C: csminiMapDataStop,
	0xB21D: csSendCommand,
	0xB21F: csSetEnterHuntZone,
	0xB3:   npcDBList,
	0xB4:   csGuildRegList,
	0xB5:   csGuildAttackList,
	0xB701: weaponUse,
	0xB704: weaponDamage,
	0xB902: guildMarkOfCastleOwner,
	0xB905: castleHuntZoneEntrance,
	0xBC00: jewelMix,
	0xBC01: jewelUnmix,
	0xBD00: reqCrywolfInfo,
	0xBD03: reqAltarContract,
	0xBD09: reqPlusChaosRate,
	0xBF02: illusionTempleUseSkill,
	0xBF05: illusionTempleReward,
	0xBF0B: reqLuckyCoinInfo,
	0xBF0C: reqLuckyCoinRegister,
	0xBF0D: reqLuckyCoinTrade,
	0xBF0E: reqEnterDoppelGanger,
	0xBF17: moveGate,
	0xBF20: inventoryEquipment,
	0xBF51: reqMuBotUse,
	0xBF6A: reqITLRelics,
	0xBF70: reqEnterITR,
	0xBF71: reqAcceptEnterITR,
	0xC0:   friendList,
	0xC1:   friendAdd,
	0xC2:   friendAddWait,
	0xC3:   friendDel,
	0xC4:   friendState,
	0xC5:   friendMemo,
	0xC7:   friendMemoRead,
	0xC8:   friendMemoDel,
	0xC9:   friendMemoList,
	0xCA:   friendChatRoomCreate,
	0xCB:   friendRoomInvitation,
	0xCD01: reqUBFMyCharInfo,
	0xCD02: reqUBFJoin,
	0xCD06: reqUBFGetReward,
	0xCD07: reqUBFCancel,
	0xCD08: reqUBFGetRealName,
	0xD007: wereWolfQuarrelCheck,
	0xD008: gateKeeperCheck,
	0xD00A: checkGateLevel,
	0xD010: reqSantaGift,
	0xD100: reqKanturuStateInfo,
	0xD101: reqEnterKanturuBossMap,
	0xD201: cashPoint,
	0xD202: cashShopOpen,
	0xD203: cashItemBuf,
	0xD204: cashItemGift,
	0xD205: cashInventoryItemCount,
	0xD20B: cashInventoryItemUse,
	0xD4:   moveProc, // s9
	0xD7:   moveProc, // 1.04R
	0xD9:   attack,   //1.04R
	0xDB00: reqDSFSchedule,
	0xDB01: reqDSFCanPartyEnter,
	0xDB02: reqAcceptEnterDSF,
	0xDB03: reqDSFGoFinalParty,
	0xDB09: reqDSFGetReward,
	0xDF:   attackedGet, // s9
	0xE1:   guildAssignStatus,
	0xE2:   guildAssignType,
	0xE5:   relationshipReqJoinBreakOff,
	0xE6:   relationshipAnsJoinBrreakOff,
	0xE701: reqSendMemberPosInfoStart,
	0xE702: reqSendMemberPosInfoEnd,
	0xE703: reqSendNpcPosInfo,
	0xE9:   unionList,
	0xEB01: relationshipReqKickOutUnionMember,
	0xEC00: inJewelPentagramItem,
	0xEC01: outJewelPentagramItem,
	0xEC02: reqRefinePentagramJewel,
	0xEC03: reqUpgradePentagramJewel,
	0xEC31: reqSearchItemInPShop,
	0xEC33: reqPShopLog,
	0xED00: guildMatchList,
	0xED01: guildMatchListSearchWord,
	0xED02: guildRegMatchList,
	0xED03: guildMatchListCancel,
	0xED04: guidMatchJoin,
	0xED05: guildMatchJoinCancel,
	0xED06: guildMatchJoinAllow,
	0xED07: guildMatchJoinAllowList,
	0xED08: guildMatchWaitStateList,
	0xEF00: partyRegWantedMember,
	0xEF01: partyMatchList,
	0xEF02: partyMatchJoinMember,
	0xEF03: partyMemberJoinStateList,
	0xEF04: partyMemberJoinStateListLeader,
	0xEF05: partyMatchJoinMemberAccept,
	0xEF06: partyMatchCancel,
	0xF101: login,
	0xF102: closeNotify,
	0xF103: clientMsg,
	0xF300: charListGet,
	0xF301: charCreate,
	0xF302: charDel,
	0xF303: charJoin,
	0xF306: levelUpPointAdd,
	0xF312: moveDataLoadOK,
	0xF315: charCheck,
	0xF321: transformationRingUse,
	0xF326: popupTypeUse,
	0xF330: skillKeyGet,
	0xF331: trapMsg,
	0xF352: masterLevelSkillGet,
	0xF60A: questSwitch,
	0xF60B: questProgress,
	0xF60D: questComplete,
	0xF60F: questGiveUp,
	0xF610: tutorialKeyCOmplete,
	0xF61A: questProgressList,
	0xF61B: questProgressInfo,
	0xF621: eventItemQuestList,
	0xF630: questExp,
	0xF631: AttDefPowerInc,
	0xF701: enterZone,
	0xF801: gensMemberJoin,
	0xF803: gensMemberSecede,
	0xF809: gensReward,
	0xF80B: gensMemberInfo,
	0xF820: acheronEnter,
	0xF830: arcaBattleGuildMasterJoin,
	0xF832: arcaBattleGuildMemberJoin,
	0xF834: arcaBattleEnter,
	0xF836: arcaBattleBootyExchange,
	0xF83C: spritMapExchange,
	0xF841: reqRegisteredMemberCnt,
	0xF843: arcaBattleMarkReg,
	0xF845: arcaBattleMarkRank,
	0xF84B: acheronEventEnter,
	0xFA08: antiCheat,
	0xFA09: dgGuildMember,
	0xFA0A: antiHackBreach,
	0xFA0D: fileCRC,
	0xFA11: antiHackCheck,
	0xFA15: hitHack,
}
*/
