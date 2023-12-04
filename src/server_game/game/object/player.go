package object

import (
	"context"
	"log"
	"math"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/item"
	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/model"
	"github.com/xujintao/balgass/src/server_game/game/shop"
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
	FillHP               int
	FillHPCount          int
	Experience           int
	ExperienceNext       int
	ExperienceMaster     int
	ExperienceMasterNext int
	masterLevel          int
	LevelUpPoint         int
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
	// attackDamageLeft     int // 物攻左
	// attackDamageRight    int // 物攻右
	attackDamageLeftMin  int // 物攻左min
	attackDamageLeftMax  int // 物攻左min
	attackDamageRightMin int // 物攻右min
	attackDamageRightMax int // 物攻右max
	magicDamageMin       int // 魔攻min
	magicDamageMax       int // 魔攻max
	magicSpeed           int // 魔攻速度
	// curseDamageMin       int // 诅咒min
	// curseDamageMax       int // 诅咒max
	// curseSpell           int
	DamageMinus        int // 伤害减少
	DamageReflect      int // 伤害反射
	MonsterDieGetMoney int // 杀怪加钱
	MonsterDieGetLife  int // 杀怪回生
	MonsterDieGetMana  int // 杀怪回蓝
	item380Effect      item.Item380Effect
	// criticalDamage       int
	excellentDamage    int // 卓越一击概率
	Inventory          model.Inventory
	InventoryExpansion int
	Warehouse          model.Warehouse
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
			Index:       c.Position,
			Name:        c.Name,
			Level:       c.Level,
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
	p.LevelUpPoint = c.LevelUpPoint
	p.MapNumber = c.MapNumber
	p.X, p.TX = c.X, c.X
	p.Y, p.TY = c.Y, c.Y
	p.Dir = c.Dir
	if p.Level <= 10 {
		p.spawnPosition()
	}
	p.createFrustrum()
	p.Strength = c.Strength
	p.Dexterity = c.Dexterity
	p.Vitality = c.Vitality
	p.Energy = c.Energy
	p.Leadership = c.Leadership
	p.skills = c.Skills
	p.Inventory = c.Inventory
	p.InventoryExpansion = c.InventoryExpansion
	p.Money = c.Money
	p.Experience = c.Experience

	// calculate
	// mc := MonsterTable[249]
	mc := MonsterTable[10]
	p.attackPanelMin = mc.DamageMin
	p.attackPanelMax = mc.DamageMax
	p.attackRate = mc.AttackRate
	p.attackSpeed = mc.AttackSpeed
	p.defense = mc.Defense
	p.magicDefense = mc.MagicDefense
	p.defenseRate = mc.BlockRate
	p.HP = mc.HP
	p.MaxHP = mc.HP
	p.MP = 1000
	p.MaxMP = 1000
	p.SD = 100
	p.MaxSD = 100
	p.AG = 200
	p.MaxAG = 200
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
		LevelUpPoint:       p.LevelUpPoint,
		Strength:           p.Strength,
		Dexterity:          p.Dexterity,
		Vitality:           p.Vitality,
		Energy:             p.Energy,
		Leadership:         p.Leadership,
		HP:                 p.HP,
		MaxHP:              p.MaxHP,
		MP:                 p.MP,
		MaxMP:              p.MaxMP,
		SD:                 p.SD,
		MaxSD:              p.MaxSD,
		AG:                 p.AG,
		MaxAG:              p.MaxAG,
		Money:              p.Money,
		PKLevel:            p.PKLevel,
		CtlCode:            0,
		AddPoint:           0,
		MaxAddPoint:        122,
		MinusPoint:         0,
		MaxMinusPoint:      122,
		InventoryExpansion: p.InventoryExpansion,
	})

	// reply inventory
	// client will calculate character after receiving inventory msg
	p.push(&model.MsgItemListReply{
		Items: p.Inventory.Items,
	})
	p.push(&model.MsgSkillListReply{
		Skills: p.skills,
	})
	p.loadMiniMap()
	p.pushMaxHP(p.MaxHP+p.AddHP, p.MaxSD+p.AddSD)
	p.pushHP(p.HP, p.SD)
	p.pushMaxMP(p.MaxMP+p.AddMP, p.MaxAG+p.AddAG)
	p.pushMP(p.MP, p.AG)
	// go func() {
	// 	time.Sleep(100 * time.Millisecond) // get character info
	// 	p.actioner.PlayerAction(p.index, "SetCharacter", &model.MsgSetCharacter{Name: msg.Name})
	// }()
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
		LevelUpPoint:       p.LevelUpPoint,
		MapNumber:          p.MapNumber,
		X:                  p.X,
		Y:                  p.Y,
		Dir:                p.Dir,
		Strength:           p.Strength,
		Dexterity:          p.Dexterity,
		Vitality:           p.Vitality,
		Energy:             p.Energy,
		Leadership:         p.Leadership,
		Skills:             p.skills,
		Inventory:          p.Inventory,
		InventoryExpansion: p.InventoryExpansion,
		Money:              p.Money,
		Experience:         p.Experience,
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

func (player *Player) MasterLevel() bool {
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
			player.magicDamageMin += player.magicDamageMin * value / 100
			player.magicDamageMax += player.magicDamageMax * value / 100
		} else {
			if position == 0 || position == 9 {
				player.attackDamageLeftMin += player.attackDamageLeftMin * value / 100
				player.attackDamageLeftMax += player.attackDamageLeftMax * value / 100
			}
			if position == 1 || position == 9 {
				player.attackDamageRightMin += player.attackDamageRightMin * value / 100
				player.attackDamageRightMax += player.attackDamageRightMax * value / 100
			}
		}
	case item.ExcelCommonIncAttackLevel: // =20
		if wItem.Section == 5 || // 法杖
			wItem.Code == item.Code(13, 12) || // 雷链子
			wItem.Code == item.Code(13, 25) || // 冰链子
			wItem.Code == item.Code(13, 27) { // 水链子
			player.magicDamageMin += (player.Level + player.masterLevel) / value
			player.magicDamageMax += (player.Level + player.masterLevel) / value
		} else {
			if position == 0 || position == 9 {
				player.attackDamageLeftMin += (player.Level + player.masterLevel) / value
				player.attackDamageLeftMax += (player.Level + player.masterLevel) / value
			}
			if position == 1 || position == 9 {
				player.attackDamageRightMin += (player.Level + player.masterLevel) / value
				player.attackDamageRightMax += (player.Level + player.masterLevel) / value
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

func (player *Player) LearnMasterSkill(msg *model.MsgLearnMasterSkill) {
	// // validate skillIndex
	// skillBase, ok := skill.SkillTable[skillIndex]
	// if !ok {
	// 	log.Printf("player[%s] learn invalid skill index[%d]", player.Name, skillIndex)
	// 	return false
	// }
	// if skillBase.STID == 0 && skillBase.UseType == 0 {
	// 	return player.learnSkill(skillIndex, 0)
	// }

	// // validate player level
	// if !player.MasterLevel() {
	// 	return false // 2
	// }

	// // validate skill level
	// level := 0
	// if skill, ok := player.skills[skillIndex]; ok {
	// 	level = skill.Level
	// }
	// level += skillBase.ReqMLPoint
	// skillMaster, _ := skill.SkillMasterTable[skillIndex]
	// if level > skillMaster.MaxPoint {
	// 	return false // 4
	// }

	// // validate master point
	// if player.MasterPoint < skillBase.ReqMLPoint {
	// 	return false // 4
	// }

	// // validate new skill
	// if level == 1 {

	// }
}

func (p *Player) processAction() {}

func (p *Player) Action(msg *model.MsgAction) {
	reply := model.MsgActionReply{
		Index:  p.index,
		Action: msg.Action,
		Dir:    msg.Dir,
	}
	p.pushViewport(&reply)
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
		if s, ok := p.learnSkill(index, 0); ok {
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

		// fill object hp
		p.HP += hp
		if p.HP > p.MaxHP+p.AddHP {
			p.HP = p.MaxHP + p.AddHP
		}
		reply := model.MsgHPReply{
			Position: -1,
			HP:       p.HP,
			Flag:     0,
			SD:       p.SD,
		}
		p.push(&reply)
		it.Durability--
		if it.Durability > 0 {
			reply := model.MsgItemDurabilityReply{
				Position:   msg.SrcPosition,
				Durability: it.Durability,
				Flag:       1,
			}
			p.push(&reply)
		} else {
			p.Inventory.DropItem(msg.SrcPosition, it)
			reply := model.MsgDeleteInventoryItemReply{
				Position: msg.SrcPosition,
				Flag:     1,
			}
			p.push(&reply)
		}

		// defer fill object hp
		p.FillHP = hp
		switch it.Level {
		case 0:
			p.FillHPCount = 0
		case 1:
			p.FillHPCount = 2
		default:
			p.FillHPCount = 3
		}
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
		if s, ok := p.learnSkill(skillIndex, 0); ok {
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
	if msg.Target >= len(p.objectManager.objects) {
		return
	}
	tobj := p.objectManager.objects[msg.Target]
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
		p.targetNumber >= len(p.objectManager.objects) {
		return
	}
	tobj := p.objectManager.objects[p.targetNumber]
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
		p.targetNumber >= len(p.objectManager.objects) {
		return
	}
	tobj := p.objectManager.objects[p.targetNumber]
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
