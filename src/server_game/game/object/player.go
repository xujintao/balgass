package object

import (
	"context"
	"log"
	"math"
	"sort"
	"time"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/class"
	"github.com/xujintao/balgass/src/server_game/game/formula"
	"github.com/xujintao/balgass/src/server_game/game/item"
	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/model"
	"github.com/xujintao/balgass/src/server_game/game/shop"
	"github.com/xujintao/balgass/src/server_game/game/skill"
	"gorm.io/gorm"
)

func init() {
	var characterList []*model.Character
	conf.JSON(conf.PathCommon, "players/character.json", &characterList)
	CharacterTable = make(characterTable)
	for _, c := range characterList {
		CharacterTable[c.Class] = *c
	}
}

var CharacterTable characterTable

type characterTable map[int]model.Character

type Conn interface {
	Addr() string
	Write(any) error
	Close() error
}

type actioner interface {
	PlayerAction(int, string, any)
}

func NewPlayer(conn Conn, actioner actioner) *Player {
	// create a new player
	player := &Player{}
	player.init()
	// player.LoginMsgSend = false
	// player.LoginMsgCount = 0
	player.conn = conn
	player.msgChan = make(chan any, 100)
	ctx, cancel := context.WithCancel(context.Background())
	player.cancel = cancel
	// player.ConnectCheckTime = time.Now()
	// player.AutoSaveTime = player.ConnectCheckTime
	player.ConnectState = ConnectStateConnected
	// player.CheckSpeedHack = false
	// player.EnableCharacterCreate = false
	player.Type = ObjectTypePlayer
	player.actioner = actioner

	// new a new goroutine to reply message
	go func() {
		for {
			select {
			case msg := <-player.msgChan:
				err := player.conn.Write(msg)
				if err != nil {
					log.Printf("conn.Write failed [err]%v [player]%d [name]%s [msg]%v\n",
						err, player.index, player.Name, msg)
				}
			case <-ctx.Done():
				close(player.msgChan)
				player.conn.Close()
				return // return ctx.Err()
			}
		}
	}()
	return player
}

type Player struct {
	object
	offline              bool
	conn                 Conn
	msgChan              chan any
	cancel               context.CancelFunc
	actioner             actioner
	AccountID            int
	AccountName          string
	AccountPassword      string
	AuthLevel            int
	Experience           int
	NextExperience       int
	LevelPoint           int
	MasterExperience     int
	MasterNextExperience int
	MasterLevel          int
	MasterPoint          int
	MasterPointUsed      int
	FruitPoint           int
	Money                int
	Strength             int
	Dexterity            int
	Vitality             int
	Energy               int
	Leadership           int
	AddStrength          int
	AddDexterity         int
	AddVitality          int
	AddEnergy            int
	AddLeadership        int
	autoRecoverHPTick    int
	autoRecoverMPTick    int
	autoRecoverSDTime    time.Time
	delayRecoverHP       int
	delayRecoverHPMax    int
	delayRecoverSD       int
	delayRecoverSDMax    int
	magic                int
	magicAttackMin       int // 魔攻min
	magicAttackMax       int // 魔攻max
	curse                int
	curseAttackMin       int // 诅咒min
	curseAttackMax       int // 诅咒max
	magicSpeed           int // 魔攻速度
	// curseSpell           int
	DamageMinus        int // 伤害减少
	DamageReflect      int // 伤害反射
	MonsterDieGetMoney int // 杀怪加钱
	MonsterDieGetLife  int // 杀怪回生
	MonsterDieGetMana  int // 杀怪回蓝
	item380Effect      item.Item380Effect
	// criticalDamage       int
	excellentDamage    int // 卓越一击概率
	Inventory          item.Inventory
	InventoryExpansion int
	Warehouse          item.Warehouse
	WarehouseExpansion int
	WarehouseMoney     int
	// dbClass              uint8
	ChangeUp int // 1=1转 2=2转 3=3转
	// PKCount                    int
	PKLevel int
	// PKTime                     int
	// PKTotalCount               int
	// guild                *guild.GuildInfo
	// guildName                     string
	// guildStatus                   int
	// guildUnionTimeStamp           int
	// guildNumber                   int
	// lastMoveTime                  time.Time
	// resets                        int
	// vipType                       uint8
	// vipEffect                     uint8
	// santaCount                    uint8
	// goblinTime                    time.Time
	// securityCheck                 bool
	// securityCode                  int
	// RegisterLMS                   uint8
	// registerLMSRoom               uint8
	// jewelHarmonyEffect            item.JewelHarmonyItemEffect
	// kanturuEntranceByNPC          bool
	// gensInfoLoad                  bool
	// questInfoLoad                 bool
	// wCoinP                        int
	// wCoinC                        int
	// goblinPoint                   int
	// periodItemEffectIndex         int
	// seedOptionList                [35]item.SocketOptionList
	// bonusOptionList               [7]item.SocketOptionList
	// setOptionList                 [2]item.SocketOptionList
	// refillHPSocketOption          uint16
	// refillMPSocketOption          uint16
	// socketOptionMonsterDieGetHP   uint16
	// socketOptionMonsterDieGetMana uint16
	// AGReduceRate                  uint8
	// muBotEnable                   bool
	// muBotTotalTime                time.Duration
	// muBotPayTime                  time.Duration
	// muBotTick                     time.Time
	// WarehouseExpansion            int
	// LastAuthTime                  time.Time
	// LastXorKey1                   [4]int
	// LaskXorKey2                   [4]int
	// bot                           bool
	// botIndex                      int
	// skillHellFire2State           int
	// skillHellFire2Count           int
	// skillHellFire2Time            time.Time
	// skillStrengthenHellFire2State int
	// skillStrengthenHellFire2Count int
	// skillStrengthenHellFire2Time  time.Time
	// reqWarehouseOpen              int
	// set
	setEffectIncSkillAttack     int
	setEffectIncExcelDamage     int
	setEffectIncExcelDamageRate int
	setEffectIncCritiDamage     int
	setEffectIncCritiDamageRate int
	setEffectIncAG              int
	setEffectIncDamage          int
	setEffectIncAttackMin       int
	setEffectIncAttackMax       int
	// setEffectIncAttack             int
	setEffectIncDefense int
	// setEffectIncDefenseRate        int
	setEffectIncMagicAttack        int
	setEffectIgnoreDefense         int
	setEffectDoubleDamage          int
	setEffectTwoHandSwordIncDamage int
	setEffectIncAttackRate         int
	// setEffectReflectDamage         int
	setEffectIncShieldDefense int
	// setEffectDecAG            int
	// setEffectIncItemDropRate  int
	setFull bool
	// excel wing
	excelWingEffectIgnoreDefense int
	excelWingEffectReboundDamage int
	excelWingEffectRecoveryHP    int
	excelWingEffectRecoveryMP    int
	excelWingEffectDoubleDamage  int
}

func (p *Player) addr() string {
	return p.conn.Addr()
}

func (p *Player) Offline() {
	if p.offline {
		return
	}
	p.offline = true
	// todo
	p.SaveCharacter()
	p.cancel()
}

func (p *Player) push(msg any) {
	if p.offline {
		log.Printf("Still pushing [msg]%#v to [player]%d that already offline\n",
			msg, p.index)
		return
	}
	if len(p.msgChan) > 80 {
		p.Offline()
		return
	}
	p.msgChan <- msg
}

func (p *Player) pushMaxHP(hp, sd int) {
	p.push(&model.MsgHPReply{Position: -2, HP: hp, SD: sd})
}

func (p *Player) pushHP(hp, sd int) {
	p.push(&model.MsgHPReply{Position: -1, HP: hp, SD: sd})
}

func (p *Player) pushMaxMP(mp, ag int) {
	p.push(&model.MsgMPReply{Position: -2, MP: mp, AG: ag})
}

func (p *Player) pushMP(mp, ag int) {
	p.push(&model.MsgMPReply{Position: -1, MP: mp, AG: ag})
}

func (p *Player) pushItemDurability(position, dur int) {
	p.push(&model.MsgItemDurabilityReply{Position: position, Durability: dur, Flag: 1})
}

func (p *Player) pushDeleteItem(position int) {
	p.push(&model.MsgDeleteInventoryItemReply{Position: position, Flag: 1})
}

func (p *Player) decreaseItemDurability(position int) {
	it := p.Inventory.Items[position]
	it.Durability--
	if it.Durability > 0 {
		p.pushItemDurability(position, it.Durability)
	} else {
		p.Inventory.DropItem(position, it)
		p.pushDeleteItem(position)
	}
}

func (p *Player) spawnPosition() {
	gate := 0
	switch p.MapNumber {
	case maps.Arena, // 古战场
		maps.DuelArena, // 竞技场
		maps.Exile,     // 流放地
		maps.SantaTown: // 圣诞之地

	case maps.LorenMarket, // 罗兰市场
		maps.ImperialGuardian1, // 帝国要塞1
		maps.ImperialGuardian2, // 帝国要塞2
		maps.ImperialGuardian3, // 帝国要塞3
		maps.ImperialGuardian4, // 帝国要塞4
		maps.IllusionTemple1,   // 幻影寺院1
		maps.IllusionTemple2,   // 幻影寺院2
		maps.IllusionTemple3,   // 幻影寺院3
		maps.IllusionTemple4,   // 幻影寺院4
		maps.IllusionTemple5,   // 幻影寺院5
		maps.IllusionTemple6,   // 幻影寺院6
		maps.IllusionTemple7,   // 幻影寺院7
		maps.IllusionTemple8,   // 幻影寺院8
		maps.Doppelganger1,     //生魂广场1
		maps.Doppelganger2,     //生魂广场2
		maps.Doppelganger3,     //生魂广场3
		maps.Doppelganger4:     //生魂广场4
		gate = 333
	case maps.Lorencia, // 勇者大陆
		maps.Dungeon, // 地下城1~3
		maps.Kalima1, // 卡利玛1
		maps.Kalima2, // 卡利玛2
		maps.Kalima3, // 卡利玛3
		maps.Kalima4, // 卡利玛4
		maps.Kalima5, // 卡利玛5
		maps.Kalima6, // 卡利玛6
		maps.Kalima7: // 卡利玛7
		gate = 17 // Lorencia
	case maps.Devias, // 冰风谷1~4
		maps.Icarus,           // 天空之城
		maps.DevilSquare,      // 恶魔广场1~4
		maps.DevilSquare2,     // 恶魔广场5~7
		maps.DevilSquareFinal, // 恶魔广场
		maps.BloodCastle1,     // 血色城堡1
		maps.BloodCastle2,     // 血色城堡2
		maps.BloodCastle3,     // 血色城堡3
		maps.BloodCastle4,     // 血色城堡4
		maps.BloodCastle5,     // 血色城堡5
		maps.BloodCastle6,     // 血色城堡6
		maps.BloodCastle7,     // 血色城堡7
		maps.BloodCastle8,     // 血色城堡8
		maps.ChaosCastle1,     // 赤色要塞1
		maps.ChaosCastle2,     // 赤色要塞2
		maps.ChaosCastle3,     // 赤色要塞3
		maps.ChaosCastle4,     // 赤色要塞4
		maps.ChaosCastle5,     // 赤色要塞5
		maps.ChaosCastle6,     // 赤色要塞6
		maps.ChaosCastle7:     // 赤色要塞7
		gate = 22 // Devias
	case maps.Noria: // 仙踪林
		gate = 27 // Noria
	case maps.Elbeland: // 幻术园
		gate = 267 // Elbeland
	case maps.LostTower: // 失落之塔1~7
		gate = 42 // LostTower
	case maps.Atlans: // 亚特兰蒂斯1~3
		gate = 49 // Atlans
	case maps.Tarkan: // 死亡沙漠1~2
		gate = 57 // Tarkan
	case maps.Aida: // 幽暗森林1~2
		gate = 119 // Aida
	case maps.Kanturu: // 坎特鲁废墟1~3
		gate = 138 // Kanturu
	case maps.KanturuRemain: // 坎特鲁遗址
		gate = 139 // KanturuRemain
	case maps.Karutan1, // 卡伦特1
		maps.Karutan2: // 卡伦特2
		gate = 335 // Karutan1
	case maps.Raklion, // 冰霜之城
		maps.RaklionBoss: // 冰霜之城Boss
		gate = 287
	case maps.SwampOfCalmness: // 安宁池
		gate = 273 // SwampOfCalmness
	case maps.Acheron, // 阿卡伦
		maps.AcheronArcaWar, // 阿卡伦战役
		maps.Debenter,
		maps.DebenterArcaWar,
		maps.UrkMontain,
		maps.UrkMontainArcaWar:
		gate = 417
	case maps.Vulcanus: // 囚禁之岛
		gate = 294
	case maps.ValleyOfLoren: // 罗兰峡谷

	case maps.Crywolf:
		gate = 114 // Crywolf
	case maps.BalgassBarracks, // 巴卡斯兵营
		maps.BalgassRefuge: // 巴卡斯休息室
		gate = 256
	}
	maps.GateMoveManager.Move(gate, func(mapNumber, x, y, dir int) {
		p.MapNumber = mapNumber
		p.X, p.Y = x, y
		p.TX, p.TY = x, y
		p.Dir = dir
	})
	maps.MapManager.SetMapAttrStand(p.MapNumber, p.X, p.Y)
	p.createFrustrum()
}

func (p *Player) MuunSystem(msg *model.MsgMuunSystem) {
	log.Println("MuunSystem placeholder")
}

func (p *Player) Chat(msg *model.MsgChat) {
	l := len(msg.Msg)
	if l == 0 {
		return
	}
	switch {
	case msg.Msg[0] == '!' && l > 2: // global announcement
		return
	case msg.Msg[0] == '/' && l > 1: // command
		return
	case msg.Msg[0] == '~' || msg.Msg[0] == ']': // party
		return
	case msg.Msg[0] == '$': // gens
		return
	case msg.Msg[0] == '@': //guild
		return
	default:
		reply := model.MsgChatReply{MsgChat: *msg}
		p.pushViewport(&reply)
	}
}

func (p *Player) Whisper(msg *model.MsgWhisper) {
	if len(msg.Name) == 0 {
		return
	}
	if p.Name == msg.Name {
		return
	}
	tobj := ObjectManager.GetPlayerByName(msg.Name)
	if tobj == nil {
		reply := model.MsgWhisperReplyFailed{
			Flag: 0,
		}
		p.push(&reply)
		return
	}
	reply := model.MsgWhisperReply{}
	reply.Name = p.Name
	reply.Msg = msg.Msg
	tobj.push(&reply)
}

// func (p *Player) Live(msg *model.MsgLive) {

// }

func (p *Player) Login(msg *model.MsgLogin) {
	// validate msg
	resp := model.MsgLoginReply{Result: 1}
	defer p.push(&resp)
	account, err := model.DB.GetAccountByName(msg.Account)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.Result = 2
			return
		}
		log.Printf("model.DB.GetAccountByAccount failed [err]%v\n", err)
		resp.Result = 7
		return
	}
	if account.Password != msg.Password {
		resp.Result = 0
		return
	}
	p.AccountID = account.ID
	p.AccountName = account.Name
	p.AccountPassword = account.Password
	p.WarehouseExpansion = account.WarehouseExpansion
	p.ConnectState = ConnectStateLogged

	// async
	// go func() {
	// 	account, err := model.DB.GetAccountByName(msg.Account)
	// 	p.actioner.PlayerAction(p.index, "SetAccount", &model.MsgSetAccount{
	// 		MsgLogin: msg,
	// 		Account:  account,
	// 		Err:      err,
	// 	})
	// }()
}

// func (p *Player) SetAccount(msg *model.MsgSetAccount) {
// 	resp := model.MsgLoginReply{Result: 1}
// 	defer p.push(&resp)
// 	login := msg.MsgLogin
// 	account := msg.Account
// 	err := msg.Err
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			resp.Result = 2
// 			return
// 		}
// 		log.Printf("model.DB.GetAccountByAccount failed [err]%v\n", err)
// 		resp.Result = 7
// 		return
// 	}
// 	if account.Password != login.Password {
// 		resp.Result = 0
// 		return
// 	}
// }

func (p *Player) Logout(msg *model.MsgLogout) {
	defer p.push(&model.MsgLogoutReply{Flag: msg.Flag})
	switch msg.Flag {
	case 0: // close game
	case 1: // back to pick character
		// offline to login state.\
		p.ConnectState = ConnectStateLogged
		p.SaveCharacter()
		p.reset()
	case 2: // back to pick server
		// offline to init state.
		// Do not close connection
	default:
		log.Printf("Logout failed [flag]%d\n", msg.Flag)
	}
}

func (p *Player) Hack(msg *model.MsgHack) {
	log.Printf("Hack [flag1]%02x [flag2]%02x\n", msg.Flag1, msg.Flag2)
}

func (p *Player) GetCharacterList(msg *model.MsgGetCharacterList) {
	reply := model.MsgGetCharacterListReply{}

	// get account
	reply.EnableCharacterClass = 0xFF
	reply.WarehouseExpansion = p.WarehouseExpansion
	p.push(&model.MsgEnableCharacterClassReply{
		Class: reply.EnableCharacterClass,
	})
	p.push(&model.MsgResetCharacterReply{
		Reset: "012345678901234567",
	})

	// get character list
	chars, err := model.DB.GetCharacterList(p.AccountID)
	if err != nil {
		log.Printf("model.DB.GetCharacterList failed [err]%v\n", err)
		return
	}

	reply.CharacterList = make([]*model.MsgCharacter, len(chars))
	for i, c := range chars {
		reply.CharacterList[i] = &model.MsgCharacter{
			Index: c.Position,
			Name:  c.Name,
			Level: c.Level,
			// Level:       c.Level + c.MasterLevel,
			Class:       c.Class,
			ChangeUp:    c.ChangeUp,
			Inventory:   [9]*item.Item(c.Inventory.Items[:9]),
			GuildStatus: 0xFF,
			PKLevel:     0,
		}
	}
	p.push(&reply)
}

func (p *Player) CreateCharacter(msg *model.MsgCreateCharacter) {
	reply := model.MsgCreateCharacterReply{Result: 0}
	defer p.push(&reply)

	// validate msg
	if msg.Name == "" || msg.Class > 6 {
		log.Printf("CreateCharacter validate msg failed [msg]%v\n", msg)
		return
	}

	// try to get an empty postion
	chars, err := model.DB.GetCharacterList(p.AccountID)
	if err != nil {
		log.Printf("model.DB.GetCharacterList failed [err]%v\n", err)
		return
	}
	position := 0
	for i, char := range chars {
		if i != char.Position {
			position = i
			break
		}
		position++
	}
	if position > 4 {
		log.Printf("over max character count [account]%s\n", p.AccountName)
		return
	}
	// create character
	c := CharacterTable[msg.Class]
	c.AccountID = p.AccountID
	c.Position = position
	c.Name = msg.Name
	if err := model.DB.CreateCharacter(&c); err != nil {
		log.Printf("model.DB.CreateCharacter failed [err]%v\n", err)
		return
	}
	// reply
	reply.Result = 1
	reply.Name = c.Name
	reply.Index = c.Position
	reply.Level = c.Level
	reply.Class = c.Class
}

func (p *Player) DeleteCharacter(msg *model.MsgDeleteCharacter) {
	reply := model.MsgDeleteCharacterReply{Result: 0}
	defer p.push(&reply)

	if p.ConnectState == ConnectStatePlaying {
		return
	}

	// validate msg
	if msg.Name == "" || msg.Password == "" {
		log.Printf("DeleteCharacter validate msg failed [msg]%v\n", msg)
		return
	}

	// check password
	if msg.Password != "1234567" {
		reply.Result = 2
		return
	}

	// delete character
	if err := model.DB.DeleteCharacterByName(p.AccountID, msg.Name); err != nil {
		log.Printf("model.DB.DeleteCharacterByName failed [err]%v\n", err)
		return
	}

	// reply
	reply.Result = 1
}

func (p *Player) CheckCharacter(msg *model.MsgCheckCharacter) {
	p.push(&model.MsgCheckCharacterReply{
		Result: 0,
	})
}

func (p *Player) LoadCharacter(msg *model.MsgLoadCharacter) {
	// validate msg
	if msg.Name == "" || msg.Position < 0 || msg.Position > 4 {
		log.Printf("LoadCharacter validate msg failed [msg]%v [account]%s \n",
			msg, p.AccountName)
		return
	}

	// load character data from db
	c, err := model.DB.GetCharacterByName(p.AccountID, msg.Name)
	if err != nil {
		log.Printf("model.DB.GetCharacterByName failed [err]%v\n", err)
		return
	}

	// set player with character data
	p.Name = c.Name
	p.Annotation = c.Name
	p.Class = c.Class
	p.ChangeUp = c.ChangeUp
	p.Level = c.Level
	p.LevelPoint = c.LevelPoint
	p.Experience = c.Experience
	p.Strength = c.Strength
	p.Dexterity = c.Dexterity
	p.Vitality = c.Vitality
	p.Energy = c.Energy
	p.Leadership = c.Leadership
	p.MasterLevel = c.MasterLevel
	p.MasterPoint = c.MasterPoint
	p.MasterExperience = c.MasterExperience
	p.MasterNextExperience = 100000
	p.HP = c.HP
	p.MP = c.MP
	p.skills = c.Skills
	p.skills.FillSkillData(p.Class)
	p.Inventory = c.Inventory
	p.InventoryExpansion = c.InventoryExpansion
	p.Money = c.Money
	p.MapNumber = c.MapNumber
	p.X, p.TX = c.X, c.X
	p.Y, p.TY = c.Y, c.Y
	p.Dir = c.Dir
	if p.Level <= 10 {
		p.spawnPosition()
	}
	p.createFrustrum()

	// mc := MonsterTable[249]
	mc := MonsterTable[10]
	p.moveSpeed = mc.MoveSpeed
	p.attackRange = mc.AttackRange
	p.attackType = mc.AttackType
	p.viewRange = mc.ViewRange
	p.maxRegenTime = 4
	p.ConnectState = ConnectStatePlaying
	p.Live = true
	p.State = 1

	// p.push(&model.MsgResetGameReply{})
	// reply
	p.push(&model.MsgLoadCharacterReply{
		X:                  p.X,
		Y:                  p.Y,
		MapNumber:          p.MapNumber,
		Dir:                p.Dir,
		Experience:         p.Experience,
		NextExperience:     p.Experience + 100,
		LevelPoint:         p.LevelPoint,
		Strength:           p.Strength,
		Dexterity:          p.Dexterity,
		Vitality:           p.Vitality,
		Energy:             p.Energy,
		Leadership:         p.Leadership,
		Money:              p.Money,
		PKLevel:            p.PKLevel,
		CtlCode:            0,
		AddPoint:           0,
		MaxAddPoint:        122,
		MinusPoint:         0,
		MaxMinusPoint:      122,
		InventoryExpansion: p.InventoryExpansion,
	})
	p.loadMiniMap()

	// reply inventory
	p.push(&model.MsgItemListReply{
		Items: p.Inventory.Items,
	})
	// reply master data
	p.push(&model.MsgMasterDataReply{
		MasterLevel:          p.MasterLevel,
		MasterExperience:     p.MasterExperience,
		MasterNextExperience: p.MasterNextExperience,
		MasterPoint:          p.MasterPoint,
	})
	p.pushSkillList()
	p.pushMasterSkillList()

	// client will calculate character after receiving inventory msg and master msg
	// calculate
	p.calc()

	// go func() {
	// 	time.Sleep(100 * time.Millisecond) // get character info
	// 	p.actioner.PlayerAction(p.index, "SetCharacter", &model.MsgSetCharacter{Name: msg.Name})
	// }()
}

func (p *Player) calc() {
	leftHand := p.Inventory.Items[0]
	rightHand := p.Inventory.Items[1]
	glove := p.Inventory.Items[5]
	boot := p.Inventory.Items[6]
	wing := p.Inventory.Items[7]

	p.AddHP = 0
	p.AddMP = 0
	p.AddSD = 0
	p.AddAG = 0
	p.AddStrength = 0
	p.AddDexterity = 0
	p.AddVitality = 0
	p.AddEnergy = 0
	p.AddLeadership = 0
	strength := p.Strength + p.AddStrength
	dexterity := p.Dexterity + p.AddDexterity
	vitality := p.Vitality + p.AddVitality
	energy := p.Energy + p.AddEnergy
	leadership := p.Leadership + p.AddLeadership
	level := p.Level + p.MasterLevel

	// base attack
	leftAttackMin, leftAttackMax := 0, 0
	rightAttackMin, rightAttackMax := 0, 0
	attackMin, attackMax := 0, 0
	magic := 0
	magicAttackMin, magicAttackMax := 0, 0
	curse := 0
	curseAttackMin, curseAttackMax := 0, 0
	switch class.Class(p.Class) {
	case class.Wizard:
		formula.WizardDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.WizardMagicDamageCalc(energy, &magicAttackMin, &magicAttackMax)
	case class.Knight:
		formula.KnightDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.KnightMagicDamageCalc(energy, &magicAttackMin, &magicAttackMax)
	case class.Elf:
		if leftHand != nil && leftHand.KindB == item.KindBCrossbow ||
			rightHand != nil && rightHand.KindB == item.KindBBow {
			formula.ElfWithBowDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		} else {
			formula.ElfWithoutBowDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		}
		formula.ElfMagicDamageCalc(energy, &magicAttackMin, &magicAttackMax)
	case class.Magumsa:
		formula.GladiatorDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.GladiatorMagicDamageCalc(energy, &magicAttackMin, &magicAttackMax)
	case class.DarkLord:
		formula.LordDamageCalc(strength, dexterity, vitality, energy, leadership, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.LordMagicDamageCalc(energy, &magicAttackMin, &magicAttackMax)
	case class.Summoner:
		formula.SummonerDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.SummonerMagicDamageCalc(energy, &magicAttackMin, &magicAttackMax, &curseAttackMin, &curseAttackMax)
	case class.RageFighter:
		formula.RageFighterDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.RageFighterMagicDamageCalc(energy, &magicAttackMin, &magicAttackMax)
	}

	// Stat Specialization
	// calc Stat Specialization: increase attack power
	// http://muonline.webzen.com/guides/219/1976/season-9/season-9-character-renewal
	// - (Bonus stat generated by equipping items like weapon, armor, wing, or master skill,
	// etc is not applied to specialization calculation.)
	var options []*model.MsgStatSpec
	var percent float32

	// STAT_OPTION_INC_ATTACK_POWER
	formula.StatSpec_GetPercent(p.Class, 1, strength, dexterity, vitality, energy, leadership, &percent)
	min := leftAttackMin * int(percent) / 100
	max := leftAttackMax * int(percent) / 100
	leftAttackMin += min
	leftAttackMax += max
	rightAttackMin += min
	rightAttackMax += max
	options = append(options, &model.MsgStatSpec{ID: 1, Min: min, Max: max})

	// STAT_OPTION_INC_MAGIC_DAMAGE
	formula.StatSpec_GetPercent(p.Class, 9, strength, dexterity, vitality, energy, leadership, &percent)
	min = magicAttackMin * int(percent) / 100
	max = magicAttackMax * int(percent) / 100
	magicAttackMin += min
	magicAttackMax += max
	options = append(options, &model.MsgStatSpec{ID: 9, Min: min, Max: max})

	// STAT_OPTION_INC_CURSE_DAMAGE
	formula.StatSpec_GetPercent(p.Class, 10, strength, dexterity, vitality, energy, leadership, &percent)
	min = curseAttackMin * int(percent) / 100
	max = curseAttackMax * int(percent) / 100
	curseAttackMin += min
	curseAttackMax += max
	options = append(options, &model.MsgStatSpec{ID: 10, Min: min, Max: max})

	// base defense
	defense := 0
	formula.CalcDefense(p.Class, dexterity, &defense)

	// Stat Specialization
	// STAT_OPTION_INC_DEFENSE
	formula.StatSpec_GetPercent(p.Class, 4, strength, dexterity, vitality, energy, leadership, &percent)
	min = defense * int(percent) / 100
	defense += min
	options = append(options, &model.MsgStatSpec{ID: 4, Min: min})

	// base attack/defense success rate
	attackRate := 0
	attackRatePVP := 0
	defenseRate := 0
	defenseRatePVP := 0
	formula.CalcAttackSuccessRate_PvM(p.Class, strength, dexterity, leadership, level, &attackRate)
	formula.CalcAttackSuccessRate_PvP(p.Class, dexterity, level, &attackRatePVP)
	formula.CalcDefenseSuccessRate_PvM(p.Class, dexterity, &defenseRate)
	formula.CalcDefenseSuccessRate_PvP(p.Class, dexterity, level, &defenseRatePVP)

	// Stat Specialization
	// STAT_OPTION_INC_ATTACK_RATE
	formula.StatSpec_GetPercent(p.Class, 2, strength, dexterity, vitality, energy, leadership, &percent)
	min = attackRate * int(percent) / 100
	attackRate += min
	options = append(options, &model.MsgStatSpec{ID: 2, Min: min})

	// STAT_OPTION_INC_ATTACK_RATE_PVP
	formula.StatSpec_GetPercent(p.Class, 3, strength, dexterity, vitality, energy, leadership, &percent)
	min = attackRatePVP * int(percent) / 100
	attackRatePVP += min
	options = append(options, &model.MsgStatSpec{ID: 3, Min: min})

	// STAT_OPTION_INC_DEFENSE_RATE
	formula.StatSpec_GetPercent(p.Class, 6, strength, dexterity, vitality, energy, leadership, &percent)
	min = defenseRate * int(percent) / 100
	defenseRate += min
	options = append(options, &model.MsgStatSpec{ID: 6, Min: min})

	// STAT_OPTION_INC_DEFENSE_RATE_PVP
	formula.StatSpec_GetPercent(p.Class, 7, strength, dexterity, vitality, energy, leadership, &percent)
	min = defenseRatePVP * int(percent) / 100
	defenseRatePVP += min
	options = append(options, &model.MsgStatSpec{ID: 7, Min: min})

	// base speed
	attackSpeed := 0
	magicSpeed := 0
	formula.CalcAttackSpeed(p.Class, dexterity, &attackSpeed, &magicSpeed)

	// weapon and weapon addition attack
	if leftHand != nil {
		leftAttackMin += leftHand.DamageMin + leftHand.AdditionAttack
		leftAttackMax += leftHand.DamageMax + leftHand.AdditionAttack
		magicAttackMin += leftHand.AdditionMagicAttack
		magicAttackMax += leftHand.AdditionMagicAttack
		curseAttackMin += leftHand.AdditionCurseAttack
		curseAttackMax += leftHand.AdditionCurseAttack
		curse = leftHand.MagicPower
	}
	if rightHand != nil {
		rightAttackMin += rightHand.DamageMin + rightHand.AdditionAttack
		rightAttackMax += rightHand.DamageMax + rightHand.AdditionAttack
		magicAttackMin += rightHand.AdditionMagicAttack
		magicAttackMax += rightHand.AdditionMagicAttack
		curseAttackMin += leftHand.AdditionCurseAttack
		curseAttackMax += leftHand.AdditionCurseAttack
		magic = rightHand.MagicPower
	}

	// wing addition attack
	if wing != nil {
		leftAttackMin += wing.AdditionAttack
		leftAttackMax += wing.AdditionAttack
		rightAttackMin += wing.AdditionAttack
		rightAttackMax += wing.AdditionAttack
		magicAttackMin += wing.AdditionMagicAttack
		magicAttackMax += wing.AdditionMagicAttack
		curseAttackMin += wing.AdditionCurseAttack
		curseAttackMax += wing.AdditionCurseAttack
	}

	// convert left/right attack to attack
	left, right := false, false
	if leftHand != nil &&
		leftHand.KindA == item.KindAWeapon &&
		leftHand.Code != item.Code(4, 7) &&
		leftHand.Code != item.Code(4, 15) {
		left = true
	}
	if rightHand != nil &&
		rightHand.KindA == item.KindAWeapon &&
		rightHand.Code != item.Code(4, 7) &&
		rightHand.Code != item.Code(4, 15) {
		right = true
	}
	switch {
	case left, right:
		switch class.Class(p.Class) {
		case class.Knight, class.Magumsa:
			if leftHand.Code == rightHand.Code {
				formula.CalcTwoSameWeaponBonus(
					leftAttackMin, leftAttackMax, rightAttackMin, rightAttackMax,
					&leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
			} else {
				formula.CalcTwoDifferentWeaponBonus(
					leftAttackMin, leftAttackMax, rightAttackMin, rightAttackMax,
					&leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
			}
		case class.RageFighter:
			formula.CalcRageFighterTwoWeaponBonus(
				leftAttackMin, leftAttackMax, rightAttackMin, rightAttackMax,
				&leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		}
		attackMin = leftAttackMin + rightAttackMin
		attackMax = leftAttackMax + rightAttackMax
		attackSpeed += (leftHand.AttackSpeed + rightHand.AttackSpeed) / 2
		magicSpeed += (leftHand.AttackSpeed + rightHand.AttackSpeed) / 2
	case left:
		attackMin = leftAttackMin
		attackMax = leftAttackMax
		attackSpeed += leftHand.AttackSpeed
		magicSpeed += leftHand.AttackSpeed
	case right:
		attackMin = rightAttackMin
		attackMax = rightAttackMax
		attackSpeed += rightHand.AttackSpeed
		magicSpeed += rightHand.AttackSpeed
	default:
		attackMin = (leftAttackMin + rightAttackMin) / 2
		attackMax = (leftAttackMax + rightAttackMax) / 2
	}

	// armor(shield|armor|wing) addition defense
	for i := 1; i <= 7; i++ {
		it := p.Inventory.Items[i]
		if it != nil {
			defense += it.Defense
			defense += it.AdditionDefense
		}
	}

	// pet defense

	// shield defense rate
	if rightHand != nil {
		defenseRate += rightHand.SuccessfulBlocking
		defenseRate += rightHand.AdditionDefenseRate
	}

	// armor bonus defense and defense rate
	// defense level>=10 bonus of item the same type contributed
	// defense success rate bonus of item the same type contributed
	sameCount := 0
	level10Count := 0
	level11Count := 0
	level12Count := 0
	level13Count := 0
	level14Count := 0
	level15Count := 0
	if boot != nil {
		refer := boot.Code % item.MaxItemIndex
		for i := 2; i <= 6; i++ {
			it := p.Inventory.Items[i]
			if it != nil && it.Code%item.MaxItemIndex == refer {
				sameCount++
				if it.Level > 9 {
					level10Count++
				}
				if it.Level > 10 {
					level11Count++
				}
				if it.Level > 11 {
					level12Count++
				}
				if it.Level > 12 {
					level13Count++
				}
				if it.Level > 13 {
					level14Count++
				}
				if it.Level > 14 {
					level15Count++
				}
			}
		}
		if p.Class == int(class.Magumsa) || p.Class == int(class.RageFighter) {
			sameCount++
			level10Count++
			level11Count++
			level12Count++
			level13Count++
			level14Count++
			level15Count++
		}
		if sameCount == 5 {
			defenseRate += defenseRate / 10
			switch {
			case level15Count == 5:
				defense += defense * 30 / 100
			case level14Count == 5:
				defense += defense * 25 / 100
			case level13Count == 5:
				defense += defense * 20 / 100
			case level12Count == 5:
				defense += defense * 15 / 100
			case level11Count == 5:
				defense += defense * 10 / 100
			case level10Count == 5:
				defense += defense * 5 / 100
			}
		}
	}

	// glove speed
	if glove != nil {
		attackSpeed += glove.AttackSpeed
		magicSpeed += glove.AttackSpeed
	}

	// ...

	// hp,mp
	c := CharacterTable[p.Class]
	hp := c.HP + int(float32(level-1)*c.LevelHP) + int(float32(vitality-c.Vitality)*c.VitalityHP)
	mp := c.MP + int(float32(level-1)*c.LevelMP) + int(float32(energy-c.Energy)*c.EnergyMP)

	// sd
	sdGageConstA := conf.CommonServer.GameServerInfo.SDGageConstA
	sdGageConstB := conf.CommonServer.GameServerInfo.SDGageConstB
	expressionA := strength + dexterity + vitality + energy
	if p.Class == int(class.DarkLord) {
		expressionA += leadership
	}
	expressionB := level * level / sdGageConstB
	sd := expressionA*sdGageConstA/10 + expressionB + defense

	// ag
	ag := 0
	f := func(s, d, v, e, l float32) int {
		return int(float32(strength)*s + float32(dexterity)*d + float32(vitality)*v + float32(energy)*e + float32(leadership)*l)
	}
	switch class.Class(p.Class) {
	case class.Wizard:
		ag = f(0.2, 0.4, 0.3, 0.2, 0)
	case class.Knight:
		ag = f(0.15, 0.2, 0.3, 1.0, 0)
	case class.Elf:
		ag = f(0.3, 0.2, 0.3, 0.2, 0)
	case class.Magumsa:
		ag = f(0.2, 0.25, 0.3, 0.15, 0)
	case class.DarkLord:
		ag = f(0.3, 0.2, 0.1, 0.15, 0.3)
	case class.Summoner:
		ag = f(0.2, 0.25, 0.3, 0.15, 0)
	case class.RageFighter:
		ag = f(0.15, 0.2, 0.3, 1.0, 0)
	}

	// sumary
	p.attackMin, p.attackMax = attackMin, attackMax
	p.magic = magic
	p.magicAttackMin, p.magicAttackMax = magicAttackMin, magicAttackMax
	p.curse = curse
	p.curseAttackMin, p.curseAttackMax = curseAttackMin, curseAttackMax
	p.defense = defense
	p.attackRate = attackRate
	p.defenseRate = defenseRate
	p.attackSpeed, p.magicSpeed = attackSpeed, magicSpeed
	p.MaxHP = hp
	if p.HP > p.MaxHP+p.AddHP {
		p.HP = p.MaxHP + p.AddHP
	}
	p.MaxMP = mp
	if p.MP > p.MaxMP+p.AddMP {
		p.MP = p.MaxMP + p.AddMP
	}
	p.MaxSD = sd
	if p.SD > p.MaxSD+p.AddSD {
		p.SD = p.MaxSD + p.AddSD
	}
	p.MaxAG = ag
	if p.AG > p.MaxAG+p.AddAG {
		p.AG = p.MaxAG + p.AddAG
	}

	// push
	p.push(&model.MsgStatSpecReply{
		Options: options,
	})
	p.push(&model.MsgAttackSpeedReply{
		AttackSpeed: attackSpeed,
		MagicSpeed:  magicSpeed,
	})
	p.pushMaxHP(p.MaxHP+p.AddHP, p.MaxSD+p.AddSD)
	p.pushMaxMP(p.MaxMP+p.AddMP, p.MaxAG+p.AddAG)
	p.pushHP(p.HP, p.SD)
	p.pushMP(p.MP, p.AG)
}

func (p *Player) loadMiniMap() {
	maps.MiniManager.ForEachMapNpc(p.MapNumber, func(id, display, x, y int, name string) {
		reply := model.MsgMiniMapReply{
			ID:          id,
			IsNpc:       1,
			DisplayType: display,
			Type:        0,
			X:           x,
			Y:           y,
			Name:        name,
		}
		p.push(&reply)
	})
	maps.MiniManager.ForEachMapEntrance(p.MapNumber, func(id, display, x, y int, name string) {
		reply := model.MsgMiniMapReply{
			ID:          id,
			IsNpc:       0,
			DisplayType: 1,
			Type:        0,
			X:           x,
			Y:           y,
			Name:        name,
		}
		p.push(&reply)
	})
}

func (p *Player) MapDataLoadingOK(msg *model.MsgMapDataLoadingOK) {}

func (p *Player) SaveCharacter() {
	if p.Name == "" {
		return
	}
	c := model.Character{
		ChangeUp:           p.ChangeUp,
		Level:              p.Level,
		LevelPoint:         p.LevelPoint,
		Experience:         p.Experience,
		Strength:           p.Strength,
		Dexterity:          p.Dexterity,
		Vitality:           p.Vitality,
		Energy:             p.Energy,
		Leadership:         p.Leadership,
		MasterLevel:        p.MasterLevel,
		MasterPoint:        p.MasterPoint,
		MasterExperience:   p.MasterExperience,
		HP:                 p.HP,
		MP:                 p.MP,
		Skills:             p.skills,
		Inventory:          p.Inventory,
		InventoryExpansion: p.InventoryExpansion,
		Money:              p.Money,
		MapNumber:          p.MapNumber,
		X:                  p.X,
		Y:                  p.Y,
		Dir:                p.Dir,
	}
	err := model.DB.UpdateCharacter(p.Name, &c)
	if err != nil {
		log.Printf("model.DB.SaveCharacter failed [err]%v\n", err)
		return
	}
}

func (p *Player) getPKLevel() int {
	return p.PKLevel
}

func (player *Player) GetMasterLevel() bool {
	return player.ChangeUp == 2 && player.Level >= conf.Common.General.MaxLevelNormal
}

func (player *Player) addExcelCommonEffect(opt *item.ExcelCommon, wItem *item.Item, position int) {
	id := opt.ID
	value := opt.Value
	switch id {
	case item.ExcelCommonIncMPMonsterDie: // 杀怪回蓝
		player.MonsterDieGetMana++
	case item.ExcelCommonIncHPMonsterDie: // 杀怪回红
		player.MonsterDieGetLife++
	case item.ExcelCommonIncAttackSpeed: // 攻速
		player.attackSpeed += value
		player.magicSpeed += value
	case item.ExcelCommonIncAttackPercent: // 2%
		if wItem.Section == 5 || // 法杖
			wItem.Code == item.Code(13, 12) || // 雷链子
			wItem.Code == item.Code(13, 25) || // 冰链子
			wItem.Code == item.Code(13, 27) { // 水链子
			player.magicAttackMin += player.magicAttackMin * value / 100
			player.magicAttackMax += player.magicAttackMax * value / 100
		} else {
			if position == 0 || position == 9 {
				// player.attackDamageLeftMin += player.attackDamageLeftMin * value / 100
				// player.attackDamageLeftMax += player.attackDamageLeftMax * value / 100
			}
			if position == 1 || position == 9 {
				// player.attackDamageRightMin += player.attackDamageRightMin * value / 100
				// player.attackDamageRightMax += player.attackDamageRightMax * value / 100
			}
		}
	case item.ExcelCommonIncAttackLevel: // =20
		if wItem.Section == 5 || // 法杖
			wItem.Code == item.Code(13, 12) || // 雷链子
			wItem.Code == item.Code(13, 25) || // 冰链子
			wItem.Code == item.Code(13, 27) { // 水链子
			player.magicAttackMin += (player.Level + player.MasterLevel) / value
			player.magicAttackMax += (player.Level + player.MasterLevel) / value
		} else {
			if position == 0 || position == 9 {
				// player.attackDamageLeftMin += (player.Level + player.MasterLevel) / value
				// player.attackDamageLeftMax += (player.Level + player.MasterLevel) / value
			}
			if position == 1 || position == 9 {
				// player.attackDamageRightMin += (player.Level + player.MasterLevel) / value
				// player.attackDamageRightMax += (player.Level + player.MasterLevel) / value
			}
		}
	case item.ExcelCommonIncExcelDamage: // 一击
		player.excellentDamage += value
	case item.ExcelCommonIncZen: // 加钱
		player.MonsterDieGetMoney += value
	case item.ExcelCommonIncDefenseRate: // f10
		player.successfulBlocking += player.successfulBlocking * value / 100
	case item.ExcelCommonReflectDamage: // 反伤
		player.DamageReflect += value
	case item.ExcelCommonDecDamage: // 减伤
		player.DamageMinus += value
	case item.ExcelCommonIncMaxMP: // 加魔
		player.AddMP += player.MaxMP * value / 100
	case item.ExcelCommonIncMaxHP: // 加生
		player.AddHP += player.MaxHP * value / 100
	}
}

func (player *Player) addExcelWingEffect(opt *item.ExcelWing, wItem *item.Item) {
	id := opt.ID
	value := opt.Value
	switch id {
	case 0, 9, 13: // incHP
		player.AddHP += value + wItem.Level*5
	case 1, 10, 14: // incMP
		player.AddMP += value + wItem.Level*5
	case 2, 5, 11, 15, 16, 18, 23: // ingore
		player.excelWingEffectIgnoreDefense = value
	// case 3: // AG ?
	// 	player.AddAG += value
	// case 4: // speed ?
	// 	player.attackSpeed += value
	// 	player.magicSpeed += value
	case 6, 19:
		player.excelWingEffectReboundDamage = value
	case 7, 17, 20, 24:
		player.excelWingEffectRecoveryHP = value
	case 8, 21:
		player.excelWingEffectRecoveryMP = value
	case 12:
		player.AddLeadership += value + wItem.Level*5
	case 22:
		player.excelWingEffectDoubleDamage = value
	}
}

func (player *Player) addSetEffect(index item.SetEffectType, value int) {
	switch index {
	case item.SetEffectIncStrength:
		player.AddStrength += value
	case item.SetEffectIncAgility:
		player.AddDexterity += value
	case item.SetEffectIncEnergy:
		player.AddEnergy += value
	case item.SetEffectIncVitality:
		player.AddVitality += value
	case item.SetEffectIncLeadership:
		player.AddLeadership += value
	case item.SetEffectIncMaxHP:
		player.AddHP += value
	case item.SetEffectIncMaxMP:
		player.AddMP += value
	case item.SetEffectIncMaxAG:
		player.AddAG += value
	case item.SetEffectDoubleDamage:
		player.setEffectDoubleDamage += value
	case item.SetEffectIncShieldDefense:
		player.setEffectIncShieldDefense += value
	case item.SetEffectIncTwoHandSwordDamage:
		player.setEffectTwoHandSwordIncDamage += value
	case item.SetEffectIncAttackMin:
		player.setEffectIncAttackMin += value
	case item.SetEffectIncAttackMax:
		player.setEffectIncAttackMax += value
	case item.SetEffectIncMagicAttack:
		player.setEffectIncMagicAttack += value
	case item.SetEffectIncDamage:
		player.setEffectIncDamage += value
	case item.SetEffectIncAttackRate:
		player.setEffectIncAttackRate += value
	case item.SetEffectIncDefense:
		player.setEffectIncDefense += value
	case item.SetEffectIgnoreDefense:
		player.setEffectIgnoreDefense += value
	case item.SetEffectIncAG:
		player.setEffectIncAG += value
	case item.SetEffectIncCritiDamage:
		player.setEffectIncCritiDamage += value
	case item.SetEffectIncCritiDamageRate:
		player.setEffectIncCritiDamageRate += value
	case item.SetEffectIncExcelDamage:
		player.setEffectIncExcelDamage += value
	case item.SetEffectIncExcelDamageRate:
		player.setEffectIncExcelDamageRate += value
	case item.SetEffectIncSkillAttack:
		player.setEffectIncSkillAttack += value
	}
}

func (player *Player) CalcExcelItem() {
	// for i, wItem := range player.Inventory[0:InventoryWearSize] {
	// 	if wItem.Durability == 0 {
	// 		continue
	// 	}
	// 	if wItem.Excel == 0 {
	// 		continue
	// 	}
	// 	if i == 7 {
	// 		for _, opt := range item.ExcelManager.Wings.Options {
	// 			if wItem.KindA == opt.ItemKindA && wItem.KindB == opt.ItemKindB {
	// 				if wItem.Excel&opt.Number == opt.Number {
	// 					player.addExcelWingEffect(opt, wItem)
	// 				}
	// 			}
	// 		}
	// 	} else {
	// 		for _, opt := range item.ExcelManager.Common.Options {
	// 			switch wItem.KindA {
	// 			case opt.ItemKindA1, opt.ItemKindA2, opt.ItemKindA3:
	// 				if wItem.Excel&opt.Number == opt.Number {
	// 					player.addExcelCommonEffect(opt, wItem, i)
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}

func (player *Player) CalcSetItem() {
	type set struct {
		index int
		count int
	}
	var sets []set

	sameWeapon := 0
	sameRing := 0
	for i, wItem := range player.Inventory.Items[0:InventoryWearSize] {
		if wItem.Durability == 0 {
			continue
		}
		tierIndex := wItem.GetSetTierIndex()
		if tierIndex == 0 {
			continue
		}
		index := item.SetManager.GetSetIndex(wItem.Section, wItem.Index, tierIndex)
		if index <= 0 {
			continue
		}
		if i == 0 {
			sameWeapon = index
		}
		if i == 1 && sameWeapon > 0 {
			continue
		}
		if i == 10 {
			sameRing = index
		}
		if i == 11 && sameRing > 0 {
			continue
		}
		ok := false
		for i := range sets {
			if sets[i].index == index {
				sets[i].count++
				ok = true
				break
			}
		}
		if !ok {
			sets = append(sets, set{index, 1})
		}
	}

	for _, v := range sets {
		index := v.index
		count := v.count
		if count >= 2 {
			for i := 0; i < count-1; i++ {
				setEffect := item.SetManager.GetSet(index, i)
				player.addSetEffect(setEffect.Index, setEffect.Value)
			}

			if count > item.SetManager.GetSetEffectCount(index) {
				player.setFull = true
				setEffects := item.SetManager.GetSetFull(index)
				for _, setEffect := range setEffects {
					player.addSetEffect(setEffect.Index, setEffect.Value)
				}
			}
		}
	}
}

func (player *Player) Calc380Item() {
	for _, wItem := range player.Inventory.Items[0:InventoryWearSize] {
		if wItem.Durability == 0 {
			continue
		}
		if !wItem.Option380 || !item.Item380Manager.Is380Item(wItem.Section, wItem.Index) {
			continue
		}
		item.Item380Manager.Apply380ItemEffect(wItem.Section, wItem.Index, &player.item380Effect)
	}
	player.AddHP += player.item380Effect.Item380EffectIncMaxHP
	player.AddSD += player.item380Effect.Item380EffectIncMaxSD
}

func (p *Player) pushSkillList() {
	var skills []*skill.Skill
	p.skills.ForEachActiveSkill(func(s *skill.Skill) {
		skills = append(skills, s)
	})
	sort.Sort(skill.SortedSkillSlice(skills))
	p.push(&model.MsgSkillListReply{Skills: skills})
}

func (p *Player) pushMasterSkillList() {
	var skills []*model.MsgMasterSkill
	p.skills.ForEachMasterSkill(func(index int, level int, curValue float32, nextValue float32) {
		skills = append(skills, &model.MsgMasterSkill{
			MasterSkillUIIndex:   index,
			MasterSkillLevel:     level,
			MasterSkillCurValue:  curValue,
			MasterSkillNextValue: nextValue,
		})
	})
	p.push(&model.MsgMasterSkillListReply{Skills: skills})
}

func (p *Player) LearnMasterSkill(msg *model.MsgLearnMasterSkill) {
	reply := model.MsgLearnMasterSkillReply{
		Result:           0,
		MasterPoint:      p.MasterPoint,
		MasterSkillIndex: -1,
	}
	defer p.push(&reply)

	if p.MasterLevel <= 0 ||
		p.MasterPoint <= 0 {
		return
	}

	p.skills.GetMaster(p.Class, msg.SkillIndex, p.MasterPoint, func(point, uiIndex, index, level int, curValue, NextValue float32) {
		p.MasterPoint -= point
		reply.Result = 1
		reply.MasterPoint -= point
		reply.MasterSkillUIIndex = uiIndex
		reply.MasterSkillIndex = index
		reply.MasterSkillLevel = level
		reply.MasterSkillCurValue = curValue
		reply.MasterSkillNextValue = NextValue
	})
}

func (p *Player) ProcessAction() {}

func (p *Player) Action(msg *model.MsgAction) {
	reply := model.MsgActionReply{
		Index:  p.index,
		Action: msg.Action,
		Dir:    msg.Dir,
	}
	p.pushViewport(&reply)
}

func (p *Player) Process1000ms() {
	if p.ConnectState == ConnectStatePlaying {
		p.recoverHPSD()
		p.recoverMPAG()
	}
}

func (p *Player) recoverHPSD() {
	change := false
	totalHP := p.MaxHP + p.AddHP
	totalSD := p.MaxSD + p.AddSD
	if p.HP < totalHP {
		// auto recover HP
		p.autoRecoverHPTick++
		if p.autoRecoverHPTick >= 7 {
			p.autoRecoverHPTick = 0
			percent := 0.0
			// base item recover HP
			positions := []int{7, 9, 10, 11} // wing/pendant/ring
			for _, n := range positions {
				it := p.Inventory.Items[n]
				if it != nil && it.Durability != 0 {
					percent += float64(it.AdditionRecoverHP)
				}
			}
			// master skill recover HP
			percent += 0.0
			if percent != 0.0 {
				hp := p.HP
				hp += int(float64(totalHP) * percent / 100)
				switch {
				case hp < 1:
					hp = 1
				case hp > totalHP:
					hp = totalHP
				}
				p.HP = hp
				change = true
			}
		}
	} else {
		p.autoRecoverHPTick = 0
	}

	// auto recover SD
	if p.SD < totalSD {
		if conf.CommonServer.GameServerInfo.SDAutoRefillSafeZoneEnable {
			attr := maps.MapManager.GetMapAttr(p.MapNumber, p.X, p.Y)
			if attr&1 == 1 {
				now := time.Now()
				if now.After(p.autoRecoverSDTime.Add(time.Second * 1)) {
					p.autoRecoverSDTime = now
					expressionA := totalSD / 30
					expressionB := 100 // 380 option
					sd := p.SD
					sd += expressionA * expressionB / 100 / 25
					switch {
					case sd < 1:
						sd = 1
					case sd > totalSD:
						sd = totalSD
					}
					p.SD = sd
					change = true
				}
			}
		}
	}

	// posion delay recover hp
	if p.delayRecoverHP > 0 {
		hp := p.delayRecoverHPMax / 2
		if p.delayRecoverHP > hp {
			p.delayRecoverHP -= hp
		} else {
			hp = p.delayRecoverHP
			p.delayRecoverHP = 0
		}
		if p.HP < totalHP {
			hp += p.HP
			switch {
			case hp < 1:
				hp = 1
			case hp > totalHP:
				hp = totalHP
			}
			p.HP = hp
			change = true
		}
	}

	// posion delay recover SD
	if p.delayRecoverSD > 0 {
		sd := p.delayRecoverSDMax / 2
		if p.delayRecoverSD > sd {
			p.delayRecoverSD -= sd
		} else {
			sd = p.delayRecoverSD
			p.delayRecoverSD = 0
		}
		if p.SD < totalSD {
			sd += p.SD
			if sd > totalSD {
				sd = totalSD
			}
			p.SD = sd
			change = true
		}
	}

	if change {
		p.pushHP(p.HP, p.SD)
	}
}

func (p *Player) recoverMPAG() {
	p.autoRecoverMPTick++
	if p.autoRecoverMPTick < 3 {
		return
	}
	p.autoRecoverMPTick = 0
	change := false
	totalMP := p.MaxMP + p.AddMP
	if p.MP < totalMP {
		// base recover MP
		percent := 3.7
		// master skill recover MP
		percent += 0
		mp := p.MP
		mp += int(float64(totalMP) * percent / 100)
		switch {
		case mp < 1:
			mp = 1
		case mp > totalMP:
			mp = totalMP
		}
		p.MP = mp
		change = true
	}
	totalAG := p.MaxAG + p.AddAG
	if p.AG < totalAG {
		// base recover AG
		percent := 3.0
		// master skill recover AG
		percent += 0
		if p.Class == int(class.Knight) {
			percent = 5
		}
		ag := p.AG
		ag += 5 + int(float64(totalAG)*percent/100)
		switch {
		case ag < 1:
			ag = 1
		case ag > totalAG:
			ag = totalAG
		}
		p.AG = ag
		change = true
	}
	if change {
		p.pushMP(p.MP, p.AG)
	}
}

func (p *Player) Die(obj *object) {

}

func (p *Player) Regen() {
	p.HP = p.MaxHP + p.AddHP
	p.SD = p.MaxSD + p.AddSD
	p.MP = p.MaxMP + p.AddMP
	p.AG = p.MaxAG + p.AddAG
	reply := model.MsgReloadCharacterReply{
		X:          p.X,
		Y:          p.Y,
		MapNumber:  p.MapNumber,
		Dir:        p.Dir,
		HP:         p.HP,
		MP:         p.MP,
		SD:         p.SD,
		AG:         p.AG,
		Experience: int(p.Experience),
		Money:      p.Money,
	}
	p.push(&reply)
}

func (p *Player) GetChangeUp() int {
	return p.ChangeUp
}

func (p *Player) GetInventory() [9]*item.Item {
	return [9]*item.Item(p.Inventory.Items[:9])
}

func (p *Player) gateMove(gateNumber int) bool {
	success := false
	maps.GateMoveManager.Move(gateNumber, func(mapNumber, x, y, dir int) {
		reply := model.MsgTeleportReply{
			GateNumber: gateNumber,
			MapNumber:  mapNumber,
			X:          x,
			Y:          y,
			Dir:        dir,
		}
		p.push(&reply)
		if p.MapNumber != mapNumber {
			p.MapNumber = mapNumber
			p.loadMiniMap()
		}
		p.X, p.Y = x, y
		p.TX, p.TY = x, y
		p.Dir = dir
		p.createFrustrum()
		success = true
	})
	return success
}

func (p *Player) Teleport(msg *model.MsgTeleport) {
	p.gateMove(msg.GateNumber)
}

func (p *Player) inventoryChanged() {
	// 1, change skill
	newItemSkills := make(map[int]struct{})
	primaryHandWeapon := p.Inventory.Items[0]
	if primaryHandWeapon != nil && primaryHandWeapon.SkillIndex != 0 {
		newItemSkills[primaryHandWeapon.SkillIndex] = struct{}{}
	}
	secondaryHandWeapon := p.Inventory.Items[1]
	if secondaryHandWeapon != nil && secondaryHandWeapon.SkillIndex != 0 {
		newItemSkills[secondaryHandWeapon.SkillIndex] = struct{}{}
	}
	oldItemSkills := make(map[int]struct{})
	for _, s := range p.skills {
		if s.Index < 300 && s.SkillBase.ItemSkill {
			oldItemSkills[s.Index] = struct{}{}
		}
	}
	var needLearnSkills []int
	for newSkill := range newItemSkills {
		if _, ok := oldItemSkills[newSkill]; !ok {
			needLearnSkills = append(needLearnSkills, newSkill)
		}
	}
	var needForgetSkills []int
	for oldSkill := range oldItemSkills {
		if _, ok := newItemSkills[oldSkill]; !ok {
			needForgetSkills = append(needForgetSkills, oldSkill)
		}
	}
	for _, index := range needLearnSkills {
		if s, ok := p.learnSkill(index); ok {
			p.push(&model.MsgSkillOneReply{
				Flag:  -2,
				Skill: s,
			})
		}
	}
	for _, index := range needForgetSkills {
		if s, ok := p.forgetSkill(index); ok {
			p.push(&model.MsgSkillOneReply{
				Flag:  -1,
				Skill: s,
			})
		}
	}
	// 2, calculate player
	p.calc()
}

func (p *Player) GetItem(msg *model.MsgGetItem) {
	reply := model.MsgGetItemReply{
		Result: -1,
	}
	itemDurChanged := false
	position := -1
	var it *item.Item
	defer func() {
		p.push(&reply)
		if itemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   position,
				Durability: it.Durability,
				Flag:       0,
			}
			p.push(&reply)
		}
	}()
	item := maps.MapManager.PickItem(p.MapNumber, msg.Index)
	if item == nil {
		return
	}
	position = p.Inventory.FindFreePositionForItem(item)
	if position == -1 {
		return
	}
	it = p.Inventory.Items[position]
	if it == nil {
		p.Inventory.GetItem(position, item)
	} else {
		it.Durability += item.Durability
		itemDurChanged = true
	}
	maps.MapManager.PopItem(p.MapNumber, msg.Index)
	reply.Result = position
	reply.Item = item
}

func (p *Player) DropInventoryItem(msg *model.MsgDropInventoryItem) {
	reply := model.MsgDropInventoryItemReply{
		Result:   0,
		Position: msg.Position,
	}
	defer p.push(&reply)
	// validate
	if msg.Position >= len(p.Inventory.Items) {
		return
	}
	item := p.Inventory.Items[msg.Position]
	if item == nil {
		return
	}
	ok := maps.MapManager.PushItem(p.MapNumber, msg.X, msg.Y, item)
	if !ok {
		return
	}
	p.Inventory.DropItem(msg.Position, item)
	reply.Result = 1
}

func (p *Player) MoveItem(msg *model.MsgMoveItem) {
	reply := model.MsgMoveItemReply{
		Result: -1,
	}
	sitemDurChanged := false
	titemDurChanged := false
	var sitem *item.Item
	var titem *item.Item
	defer func() {
		p.push(&reply)
		if sitemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   msg.SrcPosition,
				Durability: sitem.Durability,
				Flag:       0,
			}
			p.push(&reply)
		}
		if titemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   msg.DstPosition,
				Durability: titem.Durability,
				Flag:       0,
			}
			p.push(&reply)
		}
	}()

	// get source item
	switch msg.SrcFlag {
	case 0:
		if msg.SrcPosition >= p.Inventory.Size {
			return
		}
		sitem = p.Inventory.Items[msg.SrcPosition]
	case 2:
		if msg.SrcPosition >= p.Warehouse.Size {
			return
		}
		sitem = p.Warehouse.Items[msg.SrcPosition]

	}
	if sitem == nil {
		return
	}
	// get destination item
	switch msg.DstFlag {
	case 0:
		if msg.DstPosition >= p.Inventory.Size {
			return
		}
		titem = p.Inventory.Items[msg.DstPosition]
	case 2:
		if msg.DstPosition >= p.Warehouse.Size {
			return
		}
		titem = p.Warehouse.Items[msg.DstPosition]
	}

	switch msg.SrcFlag {
	case 0:
		switch msg.DstFlag {
		case 0:
			switch {
			case titem == nil: // move
				ok := p.Inventory.CheckFlagsForItem(msg.DstPosition, sitem)
				if !ok {
					return
				}
				p.Inventory.DropItem(msg.SrcPosition, sitem)
				p.Inventory.GetItem(msg.DstPosition, sitem)
				if msg.SrcPosition < 12 || msg.SrcPosition == 236 ||
					msg.DstPosition < 12 || msg.DstPosition == 236 {
					p.inventoryChanged()
				}
				reply.Result = msg.DstFlag
			case titem.Overlap != 0 && // overlap
				titem.Code == sitem.Code &&
				titem.Level == sitem.Level &&
				titem.Durability < titem.Overlap:
				delta := titem.Overlap - titem.Durability
				if delta > sitem.Durability {
					delta = sitem.Durability
				}
				sitem.Durability -= delta
				sitemDurChanged = true
				if sitem.Durability <= 0 {
					reply.Result = msg.DstFlag
					sitem.Durability = 0
					sitemDurChanged = false
					p.Inventory.DropItem(msg.SrcPosition, sitem)
					reply := model.MsgDeleteInventoryItemReply{
						Position: msg.SrcPosition,
						Flag:     1,
					}
					p.push(&reply)
				} else {
					reply.Result = -1
				}
				titem.Durability += delta
				titemDurChanged = true
			default:
				return
			}
		case 2:
			if titem == nil {
				ok := p.Warehouse.CheckFlagsForItem(msg.DstPosition, sitem)
				if !ok {
					return
				}
				p.Inventory.DropItem(msg.SrcPosition, sitem)
				p.Warehouse.GetItem(msg.DstPosition, sitem)
				if msg.SrcPosition < 12 || msg.SrcPosition == 236 {
					p.inventoryChanged()
				}
				reply.Result = msg.DstFlag
			}
		default:
			return
		}
	case 2:
		switch msg.DstFlag {
		case 0:
			if titem == nil {
				ok := p.Inventory.CheckFlagsForItem(msg.DstPosition, sitem)
				if !ok {
					return
				}
				p.Warehouse.DropItem(msg.SrcPosition, sitem)
				p.Inventory.GetItem(msg.DstPosition, sitem)
				if msg.DstPosition < 12 || msg.DstPosition == 236 {
					p.inventoryChanged()
				}
				reply.Result = msg.DstFlag
			}
		case 2:
			if titem == nil {
				ok := p.Warehouse.CheckFlagsForItem(msg.DstPosition, sitem)
				if !ok {
					return
				}
				p.Warehouse.DropItem(msg.SrcPosition, sitem)
				p.Warehouse.GetItem(msg.DstPosition, sitem)
				reply.Result = msg.DstFlag
			}
		default:
			return
		}
	default:
		return
	}
	reply.Position = msg.DstPosition
	reply.Item = sitem
}

func (p *Player) LimitUseItem(it *item.Item) bool {
	if p.Level < it.ReqLevel ||
		p.Strength+p.AddStrength < it.ReqStrength ||
		p.Dexterity+p.AddDexterity < it.ReqDexterity ||
		p.Vitality+p.AddVitality < it.ReqVitality ||
		p.Energy+p.AddEnergy < it.ReqEnergy ||
		p.Leadership+p.AddLeadership < it.ReqCommand {
		return true
	}
	reqClass := it.ReqClass[p.Class]
	if reqClass == 0 || (reqClass > p.ChangeUp+1) {
		return true
	}
	return false
}

func (p *Player) UseItem(msg *model.MsgUseItem) {
	// validate the position
	if msg.SrcPosition < 0 || msg.SrcPosition >= p.Inventory.Size ||
		msg.DstPosition < 0 || msg.DstPosition >= p.Inventory.Size ||
		msg.SrcPosition == msg.DstPosition {
		return
	}
	it := p.Inventory.Items[msg.SrcPosition]
	if it == nil {
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
		addRate := 0
		switch it.Code {
		case item.Code(14, 0): // Apple 苹果
			addRate = 10
		case item.Code(14, 1): // Small Healing Potion 小瓶治疗药水
			addRate = 20
		case item.Code(14, 2): // Healing Potion 中瓶治疗药水
			addRate = 30
		case item.Code(14, 3): // Large Healing Potion 大瓶治疗药水
			addRate = 40
		}
		if it.Level >= 1 {
			addRate += 5
		}
		hp := 0
		hp += (p.MaxHP + p.AddHP) * addRate / 100
		// defer recover hp
		p.delayRecoverHP = hp
		p.delayRecoverHPMax = hp
		// decrease durability
		p.decreaseItemDurability(msg.SrcPosition)
	case it.Code >= item.Code(14, 4) && it.Code <= item.Code(14, 6): // MP Potion
		addRate := 0
		switch it.Code {
		case item.Code(14, 4):
			addRate = 20
		case item.Code(14, 5):
			addRate = 30
		case item.Code(14, 6):
			addRate = 40
		}
		totalMP := p.MaxMP + p.AddMP
		mp := totalMP * addRate / 100
		// recover mp immediately
		if p.MP < totalMP {
			p.MP += mp
			if p.MP > totalMP {
				p.MP = totalMP
			}
			p.pushMP(p.MP, p.AG)
		}
		// decrease durability
		p.decreaseItemDurability(msg.SrcPosition)
	case it.Code == item.Code(14, 7): // Siege Potion 攻城药水
	case it.Code == item.Code(14, 8): // Antidote 解毒剂
	case it.Code == item.Code(14, 9) || it.Code == item.Code(14, 20): // Ale 酒 / Remedy of Love 爱情的魔力
	case it.Code == item.Code(14, 10): // Town Portal Scroll 回城卷轴
	case it.Code == item.Code(14, 13): // Jewel of Bless
	case it.Code == item.Code(14, 14): // Jewel of Soul
	case it.Code == item.Code(14, 16): // Jewel of Life
	case it.Code >= item.Code(14, 35) && it.Code <= item.Code(14, 37): // SD Potion
		addRate := 0
		switch it.Code {
		case item.Code(14, 35):
			addRate = 25
		case item.Code(14, 36):
			addRate = 35
		case item.Code(14, 37):
			addRate = 45
		}
		sd := (p.MaxSD + p.AddSD) * addRate / 100
		p.delayRecoverSD = sd
		p.delayRecoverSDMax = sd
		// decrease durability
		p.decreaseItemDurability(msg.SrcPosition)
	case it.Code >= item.Code(14, 38) && it.Code <= item.Code(14, 40): // comples/compound Potion
		addHPRate, addSDRate := 0, 0
		switch it.Code {
		case item.Code(14, 38):
			addHPRate = 10
			addSDRate = 5
		case item.Code(14, 39):
			addHPRate = 25
			addSDRate = 10
		case item.Code(14, 40):
			addHPRate = 45
			addSDRate = 20
		}
		hp := (p.MaxHP + p.AddHP) * addHPRate / 100
		sd := (p.MaxSD + p.AddSD) * addSDRate / 100
		// defer recover hp sd
		p.delayRecoverHP = hp
		p.delayRecoverHPMax = hp
		p.delayRecoverSD = sd
		p.delayRecoverSDMax = sd
		// decrease durability
		p.decreaseItemDurability(msg.SrcPosition)
	// case it.Code == item.Code(14, 42) && it2.Type != item.TypeSocket: // 再生强化
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
		if s, ok := p.learnSkill(skillIndex); ok {
			p.push(&model.MsgSkillOneReply{
				Flag:  -2,
				Skill: s,
			})
			p.Inventory.DropItem(msg.SrcPosition, it)
			p.push(&model.MsgDeleteInventoryItemReply{
				Position: msg.SrcPosition,
				Flag:     1,
			})
		}
	}
}

func (p *Player) Talk(msg *model.MsgTalk) {
	reply := model.MsgTalkReply{
		Result: 0,
	}
	// validate
	if msg.Target >= len(ObjectManager.objects) {
		return
	}
	tobj := ObjectManager.objects[msg.Target]
	if tobj == nil {
		return
	}
	if math.Abs(float64(p.X-tobj.X)) > 5 ||
		math.Abs(float64(p.Y-tobj.Y)) > 5 {
		return
	}
	p.targetNumber = tobj.index
	switch tobj.NpcType {
	case NpcTypeShop:
		inventory := shop.ShopManager.GetShopInventory(tobj.Class, tobj.MapNumber)
		if tobj.Class == 492 {
			reply.Result = 34
		}
		p.push(&reply)
		shopItemListReply := model.MsgTypeItemListReply{Type: 0}
		shopItemListReply.Items = inventory
		p.push(&shopItemListReply)
	case NpcTypeWarehouse:
		account, err := model.DB.GetAccountByID(p.AccountID)
		if err != nil {
			log.Printf("Talk model.DB.GetAccountByID failed [err]%v\n", err)
			return
		}
		reply.Result = 2
		p.push(&reply)
		p.Warehouse = account.Warehouse
		p.WarehouseMoney = account.WarehouseMoney
		warehouseItemListReply := model.MsgTypeItemListReply{Type: 0}
		warehouseItemListReply.Items = account.Warehouse.Items
		p.push(&warehouseItemListReply)
		warehouseMoneyReply := model.MsgWarehouseMoneyReply{
			Result:         1,
			WarehouseMoney: p.WarehouseMoney,
			InventoryMoney: p.Money,
		}
		p.push(&warehouseMoneyReply)
	case NpcTypeChaosMix:
	case NpcTypeGoldarcher:
	case NpcTypePentagramMix:
	}
}

func (p *Player) CloseTalkWindow(msg *model.MsgCloseTalkWindow) {
	p.targetNumber = -1
}

func (p *Player) BuyItem(msg *model.MsgBuyItem) {
	reply := model.MsgBuyItemReply{
		Result: -1,
	}
	itemDurChanged := false
	position := -1
	var it *item.Item
	defer func() {
		p.push(&reply)
		if itemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   position,
				Durability: it.Durability,
				Flag:       0,
			}
			p.push(&reply)
		}
	}()
	// validate
	if msg.Position < 0 ||
		msg.Position >= shop.MaxShopItemCount ||
		p.targetNumber < 0 ||
		p.targetNumber >= len(ObjectManager.objects) {
		return
	}
	tobj := ObjectManager.objects[p.targetNumber]
	if tobj == nil {
		return
	}
	if math.Abs(float64(p.X-tobj.X)) > 5 ||
		math.Abs(float64(p.Y-tobj.Y)) > 5 {
		return
	}
	if tobj.NpcType != NpcTypeShop {
		return
	}
	item := shop.ShopManager.GetShopItem(tobj.Class, tobj.MapNumber, msg.Position)
	if item == nil {
		return
	}
	position = p.Inventory.FindFreePositionForItem(item)
	if position == -1 {
		return
	}
	it = p.Inventory.Items[position]
	if it == nil {
		p.Inventory.GetItem(position, item)
	} else {
		it.Durability += item.Durability
		itemDurChanged = true
	}
	reply.Result = position
	reply.Item = item
}

func (p *Player) SellItem(msg *model.MsgSellItem) {
	reply := model.MsgSellItemReply{
		Result: 0,
	}
	defer p.push(&reply)
	// validate
	if msg.Position < 0 ||
		msg.Position >= p.Inventory.Size ||
		p.targetNumber < 0 ||
		p.targetNumber >= len(ObjectManager.objects) {
		return
	}
	tobj := ObjectManager.objects[p.targetNumber]
	if tobj == nil {
		return
	}
	if math.Abs(float64(p.X-tobj.X)) > 5 ||
		math.Abs(float64(p.Y-tobj.Y)) > 5 {
		return
	}
	if tobj.NpcType != NpcTypeShop {
		return
	}
	item := p.Inventory.Items[msg.Position]
	if item == nil {
		return
	}
	p.Inventory.DropItem(msg.Position, item)
	reply.Result = 1
	reply.Money = p.Money
}

func (p *Player) CloseWarehouseWindow(msg *model.MsgCloseWarehouseWindow) {
	account := model.Account{
		ID:             p.AccountID,
		Warehouse:      p.Warehouse,
		WarehouseMoney: p.WarehouseMoney,
	}
	err := model.DB.UpdateAccountWarehouse(&account)
	if err != nil {
		log.Printf("CloseWarehouseWindow UpdateAccountWarehouse failed [err]%v\n", err)
		return
	}
	p.SaveCharacter()
	reply := model.MsgCloseWarehouseWindowReply{}
	p.push(&reply)
}

func (p *Player) MapMove(msg *model.MsgMapMove) {
	reply := model.MsgMapMoveReply{
		Result: 0,
	}
	defer p.push(&reply)
	maps.MapMoveManager.Move(msg.MoveIndex, func(gateNumber, level, money int) {
		if p.Level < level || p.Money < money {
			return
		}
		ok := p.gateMove(gateNumber)
		if ok {
			p.Money -= money
			reply := model.MsgMoneyReply{
				Result: -2,
				Money:  p.Money,
			}
			p.push(&reply)
		}
	})
}
