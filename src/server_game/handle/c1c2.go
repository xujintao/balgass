package handle

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server_game/game"
	"github.com/xujintao/balgass/src/server_game/game/model"
)

func init() {
	C1C2Handle.init(apiIns[:], apiOuts[:])
}

// C1C2Handle a c1c2 handle
var C1C2Handle c1c2Handle

type c1c2Handle struct {
	apiIns  map[int]*apiIn
	apiOuts map[any]*apiOut
}

func (h *c1c2Handle) init(apiIns []*apiIn, apiOuts []*apiOut) {
	// ingress
	h.apiIns = make(map[int]*apiIn)
	for _, v := range apiIns {
		if vv, ok := h.apiIns[v.code]; ok {
			log.Printf("duplicated api code[%d] handle[%s] with code[%d] handle[%s]",
				v.code, v.action, vv.code, vv.action)
		}
		h.apiIns[v.code] = v
	}
	// egress
	h.apiOuts = make(map[any]*apiOut)
	for _, v := range apiOuts {
		t := reflect.TypeOf(v.msg)
		if t.Kind() != reflect.Ptr {
			log.Printf("api code[%d] name[%s] msg field must be a pointer",
				v.code, v.name)
		}
		h.apiOuts[t] = v
	}
}

// *c1c2Handle.Handle implements c1c2.Handler.Handle
func (h *c1c2Handle) Handle(ctx context.Context, req *c1c2.Request) {
	v := ctx.Value(c1c2.UserContextKey)
	if v == nil {
		return
	}
	id := v.(int)

	var api *apiIn
	var ok bool
	code := int(req.Body[0])
	if api, ok = h.apiIns[code]; !ok {
		if len(req.Body) < 2 {
			log.Printf("invalid api [body]%s\n", hex.EncodeToString(req.Body))
			return
		}
		codes := []byte{req.Body[0], req.Body[1]}
		code = int(binary.BigEndian.Uint16(codes))
		if api, ok = h.apiIns[code]; !ok {
			log.Printf("invalid api [body]%s\n", hex.EncodeToString(req.Body))
			return
		}
		req.Body = req.Body[1:]
	}
	req.Body = req.Body[1:]
	if os.Getenv("DEBUG") == "1" {
		switch api.action {
		case "Live", "DefineKey":
			return
		case "Move", "Action", "Attack":
		default:
			log.Printf("[player]%d [action]%s\n", id, api.action)
		}
	}

	// validate encrypt
	if api.enc && !req.Encrypt {
		log.Printf("[%d][%s] not encrypt", api.code, api.action)
		// return
	}

	// authenticate/authorize
	// if game.GetAuthLevel(ctx) < int(api.level) {
	// 	log.Printf("[%d][%s] not authrized", api.code, api.name)
	// 	return
	// }

	// cbi.Unmarshal(req.Body, api.msg)
	t := reflect.TypeOf(api.msg)
	if _, ok := t.MethodByName("Unmarshal"); !ok {
		log.Printf("can't find Unmarshal method [msg]%s\n", t.String())
		return
	}
	msg := reflect.New(t.Elem())
	in := []reflect.Value{reflect.ValueOf(req.Body)}
	out := msg.MethodByName("Unmarshal").Call(in)
	err := out[0].Interface()
	if err != nil {
		log.Printf("Unmarshal failed [msg]%s [err]%v\n", msg.String(), err)
		return
	}
	game.Game.PlayerAction(id, api.action, msg.Interface())
}

func (h *c1c2Handle) marshal(msg any) (*c1c2.Response, error) {
	v := reflect.ValueOf(msg)
	t := v.Type()
	api, ok := h.apiOuts[t]
	if !ok {
		err := fmt.Errorf("%s has not yet be registered to api table", t.String())
		return nil, err
	}
	if _, ok := t.MethodByName("Marshal"); !ok {
		err := fmt.Errorf("%s has no Marshal Method", t.String())
		return nil, err
	}
	rets := v.MethodByName("Marshal").Call(nil)
	if len(rets) != 2 {
		err := fmt.Errorf("%s Marshal Method signature is invalid", t.String())
		return nil, err
	}
	data := rets[0].Bytes()
	err := rets[1].Interface()
	if err != nil {
		return nil, err.(error)
	}
	var buf bytes.Buffer
	if api.code>>8 != 0 {
		var codes [2]uint8
		binary.BigEndian.PutUint16(codes[:], uint16(api.code))
		buf.Write(codes[:])
	} else {
		buf.WriteByte(uint8(api.code))
	}
	buf.Write(data)

	var resp c1c2.Response
	resp.WriteHead(uint8(api.flag))
	resp.Write(buf.Bytes())
	return &resp, nil
}

type conn struct {
	*c1c2.Conn
}

func (c *conn) Addr() string {
	return c.RemoteAddr
}

func (c *conn) Write(msg any) error {
	resp, err := C1C2Handle.marshal(msg)
	if err != nil {
		return err
	}
	return c.Conn.Write(resp)
}

func (c *conn) Close() error {
	return c.Conn.Close()
}

// OnConn implements c1c2.Handler.OnConn
func (h *c1c2Handle) OnConn(c *c1c2.Conn) (any, error) {
	conn := conn{c}
	return game.Game.PlayerConn(&conn)
}

// OnClose implements c1c2.Handler.OnConn
func (h *c1c2Handle) OnClose(ctx context.Context) {
	v := ctx.Value(c1c2.UserContextKey)
	if v == nil {
		return
	}
	id := v.(int)
	game.Game.PlayerCloseConn(id)
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
	code   int
	action string
	msg    any
}

type apiOut struct {
	id   int
	enc  bool
	flag int // 0xC1 or 0xC2
	code int
	name string
	msg  any
}

var apiIns = [...]*apiIn{
	{0, false, Player, 0x00, "Chat", (*model.MsgChat)(nil)},
	{0, false, Player, 0x02, "Whisper", (*model.MsgWhisper)(nil)},
	{0, false, Player, 0x0E, "Live", (*model.MsgLive)(nil)},
	{0, false, Player, 0x11, "Attack", (*model.MsgAttack)(nil)}, // s9
	{0, false, Player, 0x18, "Action", (*model.MsgAction)(nil)},
	{0, false, Player, 0x26, "UseItem", (*model.MsgUseItem)(nil)},
	{0, false, Player, 0x4E11, "MuunSystem", (*model.MsgMuunSystem)(nil)},
	{0, false, Player, 0xD4, "Move", (*model.MsgMove)(nil)},     // s9
	{0, false, Player, 0xD7, "Move", (*model.MsgMove)(nil)},     // 1.04R
	{0, false, Player, 0xD9, "Attack", (*model.MsgAttack)(nil)}, // 1.04R
	{0, false, Guest, 0xF101, "Login", (*model.MsgLogin)(nil)},
	{0, false, Player, 0xF102, "Logout", (*model.MsgLogout)(nil)},
	{0, false, Player, 0xF103, "Hack", (*model.MsgHack)(nil)},
	{0, false, Player, 0xF300, "GetCharacterList", (*model.MsgGetCharacterList)(nil)},
	{0, false, Player, 0xF301, "CreateCharacter", (*model.MsgCreateCharacter)(nil)},
	{0, false, Player, 0xF302, "DeleteCharacter", (*model.MsgDeleteCharacter)(nil)},
	{0, false, Player, 0xF303, "LoadCharacter", (*model.MsgLoadCharacter)(nil)},
	{0, false, Player, 0xF315, "CheckCharacter", (*model.MsgCheckCharacter)(nil)},
	{0, false, Player, 0xF330, "DefineKey", (*model.MsgDefineKey)(nil)},
	{0, false, Player, 0xF352, "LearnMasterSkill", (*model.MsgLearnMasterSkill)(nil)},
	{0, false, Player, 0xFFFF, "Test", (*model.MsgTest)(nil)},
}

var apiOuts = [...]*apiOut{
	{0, false, 0xC1, 0x02, "out_chat_whisper", (*model.MsgWhisper)(nil)},
	{0, false, 0xC1, 0x11, "AttackDamageReply", (*model.MsgAttackDamageReply)(nil)},
	{0, false, 0xC2, 0x12, "CreateViewportPlayerReply", (*model.MsgCreateViewportPlayerReply)(nil)},
	{0, false, 0xC2, 0x13, "CreateViewportMonsterReply", (*model.MsgCreateViewportMonsterReply)(nil)},
	{0, false, 0xC1, 0x14, "DestroyViewportObjectReply", (*model.MsgDestroyViewportObjectReply)(nil)},
	{0, false, 0xC1, 0x17, "AttackDieReply", (*model.MsgAttackDieReply)(nil)},
	{0, false, 0xC1, 0x18, "ActionReply", (*model.MsgActionReply)(nil)},
	{0, false, 0xC1, 0x26, "HPReply", (*model.MsgHPReply)(nil)},
	{0, false, 0xC1, 0x27, "MPReply", (*model.MsgMPReply)(nil)},
	{0, false, 0xC1, 0xD4, "MoveReply", (*model.MsgMoveReply)(nil)},
	{0, false, 0xC1, 0xDE00, "EnableCharacterClassReply", (*model.MsgEnableCharacterClassReply)(nil)},
	{0, false, 0xC1, 0xEC10, "AttackHPReply", (*model.MsgAttackHPReply)(nil)},
	{0, false, 0xC1, 0xF100, "ConnectReply", (*model.MsgConnectReply)(nil)},
	{0, false, 0xC1, 0xF101, "LoginReply", (*model.MsgLoginReply)(nil)},
	{0, false, 0xC3, 0xF102, "LogoutReply", (*model.MsgLogoutReply)(nil)},
	{0, false, 0xC1, 0xF300, "GetCharacterListReply", (*model.MsgGetCharacterListReply)(nil)},
	{0, false, 0xC1, 0xF301, "CreateCharacterReply", (*model.MsgCreateCharacterReply)(nil)},
	{0, false, 0xC1, 0xF302, "DeleteCharacterReply", (*model.MsgDeleteCharacterReply)(nil)},
	{0, false, 0xC3, 0xF303, "LoadCharacterReply", (*model.MsgLoadCharacterReply)(nil)},
	{0, false, 0xC3, 0xF304, "ReloadCharacterReply", (*model.MsgReloadCharacterReply)(nil)},
	{0, false, 0xC1, 0xF311, "out_skill_list", (*model.MsgSkillList)(nil)},
	{0, false, 0xC1, 0xF315, "CheckCharacterReply", (*model.MsgCheckCharacterReply)(nil)},
	{0, false, 0xC1, 0xFA05, "AttackEffectReply", (*model.MsgAttackEffectReply)(nil)},
	{0, false, 0xC1, 0xFA0A, "ResetCharacterReply", (*model.MsgResetCharacterReply)(nil)},
	{0, false, 0xC1, 0xFFFF, "out_test", (*model.MsgTest)(nil)},
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
