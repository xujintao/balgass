package handle

import (
	"log"
	"math/rand"

	"github.com/golang/protobuf/proto"

	"github.com/xujintao/balgass/cmd/server_data/model"
	"github.com/xujintao/balgass/cmd/server_data/service"
	"github.com/xujintao/balgass/network"
)

const (
	chanNum  = 3
	chanSize = 1000
	gNum     = 10
)

type data struct {
	v   interface{}
	req *network.Request
}

type handleData struct {
	handleBase
	queue []chan data
	exit  chan struct{}
}

// Handle *CMDHandle implements network.Handler
func (h *handleData) Handle(v interface{}, req *network.Request) {
	if h.queue == nil {
		h.exit = make(chan struct{}, 1)
		h.queue = make([]chan data, chanNum)
		for i := range h.queue {
			h.queue[i] = make(chan data, chanSize)
		}
		for i := 0; i < gNum; i++ {
			go func() {
				for {
					select {
					case d := <-h.queue[i%chanNum]:
						h.handleBase.Handle(d.v, d.req)
					case <-h.exit:
						log.Println("goroutine exit")
					}
				}
			}()
		}
	}
	h.queue[rand.Int()%chanNum] <- data{v, req}
}

func (h *handleData) Exit() {
	h.exit <- struct{}{}
}

var cmdsData = map[int]func(index interface{}, req *network.Request){
	// join
	0x00: serverLogin,
	0x01: userLogin,
	0x02: userLoginFailed,
	0x04: userBlock,
	0x05: userExit,
	0x06: vipCheck,
	// 0x30: loveHeartEventReq,
	// 0x31: loveHeartCreate,
	0x7A: userServerMove,
	0x7B: userServerMoveAuth,
	0x7C: serverUserCountSet,
	0x80: userOffTradeSet,
	0xD3: vipAdd,
	0xFE: userKill,
	// 0xFF: serverDisconnect, // service.ServerManager.ServerExit(index) instead
	/*
		// data
		0x00:   loginDataServer,
		0x01:   charListGet,
		0x02:   classDefDataReq,
		0x03:   warehouseSwitchReq,
		0x04:   charCreateReq,
		0x05:   charDelReq,
		0x06:   charInfoGet,
		0x07:   charInfoSet,
		0x08:   warehouseListGet,
		0x09:   warehouseListSet,
		0x11:   userItemSave,
		0x12:   warehouseMoneySet,
		0x4C00: mineModifyUPTUserInfoReq,
		0x4C01: mineCheckReqIsUPTWhenUserConnect,
		0x4F00: gremoryCaseItemListReq,
		0x4F01: gremoryCaseItemAdd,
		0x4F02: gremoryCaseItemUseCheck,
		0x4F03: gremoryCaseItemDel,
		0x52:   ItemSerialCreate,
		0x55:   ItemPetSerialCreate,
		0x56:   ItemPetInfoGet,
		0x57:   ItemPetInfoSave,
		0x60:   optionDataGet,
		0x6F00: ItemReBuyListGet,
		0x6F01: ItemReBuyAdd,
		0x6F02: ItemReBuyGet,
		0x6F03: ItemSoldDel,
		0x8001: ansGuildMasterOwner,
		0x8003: ansCastleNPCBuy,
		0x8004: ansCastleNPCRepair,
		0x8005: ansCastleNPCUpgrade,
		0x8006: ansTaxInfo,
		0x8007: ansTaxRateChange,
		0x8008: ansCastleMoneyChange,
		0x8009: ansSiegeDateChange,
		0x800A: ansGuildMarkRegInfo,
		0x800B: ansSiegeEndedChange,
		0x800C: ansCastleOwnerChange,
		0x800D: ansAttackGuildReg,
		0x800E: ansCastleStateRestart,
		0x800F: ansMapSvrMsgMulticast,
		0x8010: ansGuildMarkReg,
		0x8011: ansGuildMarkReset,
		0x8012: ansGuildSetGiveUp,
		0x8016: ansNPCRemove,
		0x8017: ansCastleStateSync,
		0x8018: ansCastleTributeMoney,
		0x8019: ansCastleTaxInfoReset,
		0x801A: ansSiegeGuildInfoReset,
		0x801B: ansSiegeInfoReset,
		0x801F: ansGlobalPOstMulticast,
		0x81:   ansCastleDataInit,
		0x83:   ansAllGuildMarkRegInfo,
		0x84:   ansFirstNpcCreate,
		0x85:   ansCalcGuildRegList,
		0x86:   ansCsGuildUnionInfo,
		0x87:   ansCsGuildInfoSaveTotal,
		0x88:   ansCsGuildInfoLoadTotal,
		0x89:   ansCastleNpcUpdate,
		0xAE:   muBotDataSave,
		0xB0:   ansCrywolfSync,
		0xB1:   ansCrywolfInfoLoad,
		0xB2:   ansCrywolfInfoSave,
		0xBD01: eventDevilSqureScore,
		0xBD02: eventBloodCastleEnterCount,
		0xBD03: eventBloodCastleScore5th,
		0xBD04: eventIllusionTempleScore,
		0xBD05: eventChaosCastleScore,
		0xBE01: eventChipInfo,
		0xBE02: eventChipInfoReg,
		0xBE04: eventChipInfoReset,
		0xBE05: eventChipInfoStone,
		0xBE06: eventInfoStoneReg,
		0xBE07: stoneDel,
		0xBE08: anivRegSerial,
		0xBE09: StoneInfoReset,
		0xBE15: offlineGiftRegCC,
		0xBE16: offlineGiftRegDL,
		0xBE17: offlineGiftRegHT,
		0xBE18: luckyCoinReg,
		0xBE19: luckyCoinInfo,
		0xBE20: santaCheck,
		0xBE21: santaGift,
		0xC2:   whisperOtherChannel,
		0xC300: mapSrvGroupServerCount,
		0xC302: whisperResponse,
		0xC304: disconnectOtherChannel,
		0xCE:   tempUserInfoDel,
		0xCF:   moveOtherServer,
		0xD001: ItemPeriodExInsert,
		0xD002: ItemPeriodExDelete,
		0xD003: ItemPeriodExSelect,
		0xD020: ItemLuckySelect,
		0xD022: ItemLuckyItemInsert,
		0xD023: ItemLuckyDelete,
		0xD024: ItemLuckyInsert2nd,
		0xD1:   gameShopPointsGet,
		0xD2:   gameShopItemListGet,
		0xD4:   gameShopCharCardBuy,
		0xD5:   gameShopItemBuf,
		0xD6:   gameShopItemGift,
		0xD7:   gameShopPointsAdd,
		0xD801: gameShopPackageBuy,
		0xD802: gameShopPacketGift,
		0xD9:   gameShopItemUse,
		0xDA:   warehouseExpand,
		0xDB:   gameShopItemDel,
		0xDC:   gameShopItemRollbackUse,
		0xE0:   jewelPentagramGet, // 可能是指艾尔特
		0xE1:   jewelPentagramSet,
		0xE2:   jewelPentagramDel,
		0xE3:   jewelPentagramInsert,
		0xE401: buffPeriodInsert,
		0xE402: buffPeriodDelete,
		0xE403: buffPeriodSelect,
		0xE5:   arcaBattleUserJoin,
		0xE6:   itemEventInventoryLoad,
		0xE7:   itemEventInventorySave,
		0xE800: cardInfoGet,
		0xE801: cardInfoInsert,
		0xE802: scoreUpdate,
		0xE803: cardInfoUpdate,
		0xE804: scoreDelete,
		0xE805: slotInfoUpdate,
		0xE806: muRummyInfoUpdate,
		0xE807: muRummyDBLogGet,
		0xE9:   pShopItemValueGet,
		0xEB:   pShopItemValueSaveAll,
		0xEC:   pShopItemValueDel,
		0xED:   pShopItemMove,
		0xF1:   muunInventoryItemLoad,
		0xF2:   muunInventoryItemSave,
		0xF301: ubfCheckIsUserJoin,
		0xF302: ubfUserJoin,
		0xF303: ubfcharCopy,
		0xF305: ubfRewardGainSet,
		0xF306: ubfRewardGet,
		0xF307: ubfUserCancel,
		0xF308: ubfRealNameGet,
		0xF4:   ubfItemPetCopy,
		0xF6:   questExpInfoSave,
		0xF7:   questExpInfoLoad,
		0xF830: arcaBattleGuildJoin,
		0xF832: arcaBattleGuildMemberJoin,
		0xF834: arcaBattleEnter,
		0xF838: arcaBattleGuildWinInfoInsert,
		0xF839: arcaBattleGuildWInInfoGet,
		0xF83D: arcaBattleInfoDel,
		0xF83E: arcaBattleProcMulticast,
		0xF83F: arcaBattleProcState,
		0xF841: arcaBattleMemberJoinUnder,
		0xF843: arcaBattleMemberJoinUnderReq,
		0xF845: arcaBattleMemberRegCnt,
		0xF847: arcaBattleAllGuildBuffRemoveMulticast,
		0xF84D: arcaBattleMarkCntGet,
		0xF84F: arcaBattleMarkRegGet,
		0xF851: arcaBattleMarkRankGet,
		0xF853: arcaBattleMarkRegDel,
		0xF854: arcaBattleIsTopRank,
		0xF855: arcaBattleMarkRegDelAll,
		0xF8F0: acheronGuardianProcMulticast,
		0xF8FB: arcaBattleAllGuildMarkCnt,
		0xF8FD: arcaBattleMarkRegSet,
		0xF8FE: arcaBattleGuildRegInit,
		0xF902: chaosCastleUBFSetReward,
		0xF9A1: chaosCastleFinalSave,
		0xF9A2: chaosCastleFinalPermission,
		0xF9A3: chaosCastleFinalLoad,
		0xF9A4: chaosCastleFinalRenew,
		0xF9A5: chaosCastleFinalSendMsgAllSvr,
		0xF9B0: chaosCastleFinalRewardSet,
		0xFA:   banUserGet,
		0xFB:   securityLockGet,
		0xFC:   monsterCountSave,
		0xFD00: dsfPartyCanEnter,
		0xFD02: dsfPartyPointSave,
		0xFD03: dsfPartyRankRenew,
		0xFD05: dsfPartyGoFinal,
		0xFD06: dsfRewardUserInsert,
		0xFD07: dsfRewardGet,
		0xFD08: dsfRewardUBFSet,
		0xFE:   playerKillerSave,

		// exData
		0x00:   loginExDataServer,
		0x02:   charClose,
		0x30:   guildCreate,
		0x31:   guildDestory,
		0x32:   guildMemberAdd,
		0x33:   guildMemberDel,
		0x35:   guildMemberInfoGet,
		0x37:   guildScoreUpdate,
		0x38:   guildNoticeSave,
		0x39:   guildMemberAddWithoutUserIndex,
		0x50:   guildServerGroupChatting,
		0x51:   unionServerGroupChatting,
		0x52:   gensServerGroupChatting,
		0x5301: guildPeriodBuffInsert,
		0x5302: guildPeriodBuffDelete,
		0x60:   friendListGet,
		0x62:   friendStateClientGet,
		0x63:   friendAdd,
		0x64:   friendAddWait,
		0x65:   friendDel,
		0x66:   friendChatRoomCreate,
		0x70:   friendMemo,
		0x71:   friendMemoList,
		0x72:   friendMemoRead,
		0x73:   friendMemoDel,
		0x74:   friendChatRoomInvite,
		0xA0:   friendChatRoomCreateResult,
		0xA300: guildMatchListGet,
		0xA301: guildMatchListSearchWordGet,
		0xA302: guildMatchListReg,
		0xA304: guildMatchListJoin,
		0xA305: guildMatchListJoinCancel,
		0xA306: guildMatchJoinAllow,
		0xA307: guildMatchListWait,
		0xA308: guildMatchListWaitStateGet,
		0xA400: partyMemberWantedReg,
		0xA401: partyMatchListGet,
		0xA402: partyMatchMemberJoin,
		0xA403: partyMemberStateListGet,
		0xA404: partyMemberStateListLeader,
		0xA405: partyMemberJoinAccept,
		0xA406: parthMatchCancel,
		0xA407: partyUserDel,
		0xA411: partyMemberListGet,
		0xA412: partyMatchChatMsg,
		0xE1:   guildAssignStatus,
		0xE2:   guildAssignType,
		0xE5:   relationshipJoin,
		0xE6:   relationshipBreak,
		0xE9:   unionListGet,
		0xEB01: relationshipKickOutMember,
		0xF801: gensInfoGet,
		0xF803: gensMemberReg,
		0xF805: gensMemberExit,
		0xF807: gensContributePointSave,
		0xF808: gensAbuseKillUserNameSave,
		0xF809: gensAbuseInfoGet,
		0xF80C: gensRewardCheck,
		0xF80E: gensRewardComplete,
		0xF80F: gensMemberCountGet,
		0xF811: gensMemberCountGet,
		0xF812: gensRankRefreshManual,
		0xF813: gensRewardDayGet,
	*/
}

func serverLogin(index interface{}, req *network.Request) {
	msgReq := &model.ServerLoginReq{}
	if err := proto.Unmarshal(req.Body, msgReq); err != nil {
		log.Println(err)
		return
	}

	msgRes, err := service.ServerManager.ServerLogin(index, msgReq)
	if err != nil {
		log.Printf("ServerLogin failed, %v", err)
	}

	buf, err := proto.Marshal(msgRes)
	if err != nil {
		log.Printf("Marshal failed, %v", err)
		return
	}

	// write
	res := &network.Response{}
	res.WriteHead(0xC1, 0x00).Write(buf)
	if err := service.ServerManager.Send(index, res); err != nil {
		log.Printf("Send failed, %v", err)
	}
}

func userLogin(index interface{}, req *network.Request) {
	msgReq := &model.UserLoginReq{}
	if err := proto.Unmarshal(req.Body, msgReq); err != nil {
		log.Println(err)
		return
	}
	// validate username and passwd

	msgRes, err := service.UserManager.UserLogin(index, msgReq)
	if err != nil {
		log.Println(err)
		return
	}

	buf, err := proto.Marshal(msgRes)
	if err != nil {
		log.Printf("Marshal failed, %v", err)
		return
	}

	res := &network.Response{}
	res.WriteHead(0xC1, 0x01).Write(buf)
	if err := service.ServerManager.Send(index, res); err != nil {
		log.Printf("Send failed, %v", err)
	}
}

func userLoginFailed(index interface{}, req *network.Request) {
	msgReq := &model.UserLoginFailedReq{}
	if err := proto.Unmarshal(req.Body, msgReq); err != nil {
		log.Println(err)
		return
	}
	// validate username

	err := service.UserManager.UserLoginFailed(index, msgReq)
	if err != nil {
		log.Println(err)
		return
	}
}

func userBlock(index interface{}, req *network.Request) {
	msgReq := &model.UserBlockReq{}
	if err := proto.Unmarshal(req.Body, msgReq); err != nil {
		log.Println(err)
		return
	}
	// validate username

	err := service.UserManager.UserBlock(index, msgReq)
	if err != nil {
		log.Println(err)
		return
	}
}

func userExit(index interface{}, req *network.Request) {
	msgReq := &model.UserExitReq{}
	if err := proto.Unmarshal(req.Body, msgReq); err != nil {
		log.Println(err)
		return
	}
	// validate username

	err := service.UserManager.UserExit(index, msgReq)
	if err != nil {
		log.Println(err)
		return
	}
}

func vipCheck(index interface{}, req *network.Request) {
	msgReq := &model.VipCheckReq{}
	if err := proto.Unmarshal(req.Body, msgReq); err != nil {
		log.Println(err)
		return
	}
	// validate username
	msgRes, err := service.VIPManager.VIPCheck(index, msgReq)
	if err != nil {
		log.Println(err)
		return
	}

	buf, err := proto.Marshal(msgRes)
	if err != nil {
		log.Printf("Marshal failed, %v", err)
		return
	}

	res := &network.Response{}
	res.WriteHead(0xC1, 0x06).Write(buf)
	if err := service.ServerManager.Send(index, res); err != nil {
		log.Printf("Send failed, %v", err)
	}
}

func userServerMove(index interface{}, req *network.Request) {
	msgReq := &model.UserServerMoveReq{}
	if err := proto.Unmarshal(req.Body, msgReq); err != nil {
		log.Println(err)
		return
	}

	msgRes, err := service.UserManager.UserServerMove(index, msgReq)
	if err != nil {
		log.Println(err)
		return
	}

	buf, err := proto.Marshal(msgRes)
	if err != nil {
		log.Printf("Marshal failed, %v", err)
		return
	}

	res := &network.Response{}
	res.WriteHead(0xC1, 0x7A).Write(buf)
	if err := service.ServerManager.Send(index, res); err != nil {
		log.Printf("Send failed, %v", err)
	}
}

func userServerMoveAuth(index interface{}, req *network.Request) {
	msgReq := &model.UserServerMoveAuthReq{}
	if err := proto.Unmarshal(req.Body, msgReq); err != nil {
		log.Println(err)
		return
	}

	msgRes, err := service.UserManager.UserServerMoveAuth(index, msgReq)
	if err != nil {
		log.Println(err)
		return
	}

	buf, err := proto.Marshal(msgRes)
	if err != nil {
		log.Printf("Marshal failed, %v", err)
		return
	}

	res := &network.Response{}
	res.WriteHead(0xC1, 0x7B).Write(buf)
	if err := service.ServerManager.Send(index, res); err != nil {
		log.Printf("Send failed, %v", err)
	}
}

func serverUserCountSet(index interface{}, req *network.Request) {
	msgReq := &model.ServerUserCountSetReq{}
	if err := proto.Unmarshal(req.Body, msgReq); err != nil {
		log.Println(err)
		return
	}

	// service
	err := service.ServerManager.ServerUserCountSet(index, msgReq)
	if err != nil {
		log.Println(err)
		return
	}
}

func userOffTradeSet(index interface{}, req *network.Request) {
	msgReq := &model.UserOffTradeSetReq{}
	if err := proto.Unmarshal(req.Body, msgReq); err != nil {
		log.Println(err)
		return
	}
	// validate username
	err := service.UserManager.UserOffTradeSet(index, msgReq)
	if err != nil {
		log.Println(err)
		return
	}
}

func vipAdd(index interface{}, req *network.Request) {
	msgReq := &model.VipAddReq{}
	if err := proto.Unmarshal(req.Body, msgReq); err != nil {
		log.Println(err)
		return
	}
	// validate username
	err := service.VIPManager.VIPAdd(index, msgReq)
	if err != nil {
		log.Println(err)
		return
	}
}

func userKill(index interface{}, req *network.Request) {
	msgReq := &model.UserKillReq{}
	if err := proto.Unmarshal(req.Body, msgReq); err != nil {
		log.Println(err)
		return
	}
	// validate username

	err := service.UserManager.UserKill(index, msgReq)
	if err != nil {
		log.Println(err)
		return
	}
}
