package handle

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/xujintao/balgass/cmd/server_game/game/user"
	"github.com/xujintao/balgass/cmd/server_game/network"
)

// CMDHandle tcp cmd handle
type CMDHandle struct{}

// Handle *CMDHandle implements network.Handler
func (CMDHandle) Handle(req *network.Request, res *network.Response) bool {
	code := req.Code
	if h, ok := cmds[int(code)]; ok {
		return h(req, res)
	}
	subcode := req.Body[0]
	codes := []byte{code, subcode}
	code16 := binary.BigEndian.Uint16(codes)
	if h, ok := cmds[int(code16)]; ok {
		req.Body = req.Body[1:]
		return h(req, res)
	}
	log.Printf("invalid cmd, code:%02dx, body: %v", code, req.Body)
	return false
}

// HandleUDP implements network.UDPHandler
func (CMDHandle) HandleUDP(req *network.Request, res *network.Response) bool {
	code := req.Code
	if h, ok := udpcmds[int(code)]; ok {
		return h(req, res)
	}
	subcode := req.Body[0]
	codes := []byte{code, subcode}
	code16 := binary.BigEndian.Uint16(codes)
	if h, ok := udpcmds[int(code16)]; ok {
		req.Body = req.Body[1:]
		return h(req, res)
	}
	log.Printf("invalid cmd, code:%02dx, body: %v", code, req.Body)
	return false
}

// OnConn implements network.Handler.OnConn
func (CMDHandle) OnConn(conn *network.Conn) error {

	index, err := user.ObjectAdd(conn)
	if err != nil {
		return err
	}
	conn.Index = index

	res := &network.Response{}
	res.WriteHead2(0xC1, 0x00, 0x01)
	conn.Write(res)
	return nil
}

// TrackConnState track the connect state
func (CMDHandle) TrackConnState(c net.Conn, state network.ConnState) {
	switch state {
	case network.StateNew:
		log.Printf("%s connected", c.RemoteAddr().String())
	case network.StateClosed:
		log.Printf("%s disconnected", c.RemoteAddr().String())
	}
}

var udpcmds = map[int]func(req *network.Request, res *network.Response) bool{
	// 0x0100: registerServer,
}

var cmds = map[int]func(req *network.Request, res *network.Response) bool{
	/*
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
	*/
}
