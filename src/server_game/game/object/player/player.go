package player

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
	"github.com/xujintao/balgass/src/server_game/game/object"
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

func NewPlayer(conn object.Conn, actioner object.Actioner) (int, error) {
	return object.ObjectManager.AddPlayer(conn, actioner, newPlayer)
}

func newPlayer(conn object.Conn, actioner object.Actioner) *object.Object {
	// create a new player
	p := Player{}
	p.Init()
	// p.LoginMsgSend = false
	// p.LoginMsgCount = 0
	p.conn = conn
	p.msgChan = make(chan any, 100)
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	// player.ConnectCheckTime = time.Now()
	// player.AutoSaveTime = player.ConnectCheckTime
	p.ConnectState = object.ConnectStateConnected
	// p.CheckSpeedHack = false
	// p.EnableCharacterCreate = false
	p.Type = object.ObjectTypePlayer
	p.actioner = actioner

	// new a new goroutine to reply message
	go func() {
		for {
			select {
			case msg := <-p.msgChan:
				err := p.conn.Write(msg)
				if err != nil {
					log.Printf("conn.Write failed [err]%v [player]%d [name]%s [msg]%v\n",
						err, p.Index, p.Name, msg)
				}
			case <-ctx.Done():
				close(p.msgChan)
				p.conn.Close()
				return // return ctx.Err()
			}
		}
	}()
	p.Objecter = &p
	return &p.Object
}

type Player struct {
	object.Object
	offline              bool
	conn                 object.Conn
	msgChan              chan any
	cancel               context.CancelFunc
	actioner             object.Actioner
	AccountID            int
	AccountName          string
	AccountPassword      string
	CharacterID          int
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
	RecoverHP            int // 恢复生命(翅膀+项链+大师技能)
	magic                int
	magicAttackMin       int // 魔攻min
	magicAttackMax       int // 魔攻max
	curse                int
	curseAttackMin       int // 诅咒min
	curseAttackMax       int // 诅咒max
	attackRatePVP        int
	defenseRatePVP       int
	magicSpeed           int // 魔攻速度
	// curseSpell           int
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
	IncreaseAttackMin             int // 增加最小攻击(套装+洞装+大师技能)
	IncreaseAttackMax             int // 增加最大攻击(套装+洞装+大师技能)
	IncreaseMagicAttack           int // 增加魔攻(套装+洞装+大师技能)
	IncreaseSkillAttack           int // 增加技能攻击(套装+卓越+大师技能)
	SetAddDamage                  int // 增加伤害
	CriticalAttackRate            int // 幸运一击概率
	CriticalAttackDamage          int // 幸运一击伤害
	ExcellentAttackRate           int // 卓越一击概率
	ExcellentAttackDamage         int // 卓越一击伤害
	MonsterDieGetHP               int // 杀怪回生(套装+卓越+大师技能)
	MonsterDieGetMP               int // 杀怪回蓝(套装+卓越+大师技能)
	MonsterDieGetMoney            int // 杀怪加钱(卓越)
	ArmorReduceDamage             int // 防具减少伤害(卓越+洞装)
	WingIncreaseDamage            int // 翅膀增加伤害
	WingReduceDamage              int // 翅膀减少伤害
	HelperReduceDamage            int // 天使减少伤害
	ArmorReflectDamage            int // 防具伤害反射(卓越+洞装)
	DoubleDamageRate              int // 双倍伤害(套装+大师技能)
	IgnoreDefenseRate             int // 无视防御(套装+翅膀+大师技能)
	ReturnDamage                  int // 反弹伤害(翅膀+大师技能)
	RecoverMaxHP                  int // 恢复最大生命(翅膀+大师技能)
	RecoverMaxMP                  int // 恢复最大魔法(翅膀+大师技能)
	IncreaseTwoHandWeaponDamage   int // 增加双手武器伤害
	item380Effect                 item.Item380Effect
	setFull                       bool
	KnightGladiatorCalcSkillBonus float64
}

func (p *Player) Addr() string {
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

func (p *Player) Push(msg any) {
	if p.offline {
		log.Printf("Still pushing [msg]%#v to [player]%d that already offline\n",
			msg, p.Index)
		return
	}
	if len(p.msgChan) > 80 {
		p.Offline()
		return
	}
	p.msgChan <- msg
}

func (p *Player) pushMaxHP(hp, sd int) {
	p.Push(&model.MsgHPReply{Position: -2, HP: hp, SD: sd})
}

func (p *Player) pushHP(hp, sd int) {
	p.Push(&model.MsgHPReply{Position: -1, HP: hp, SD: sd})
}

func (p *Player) pushMaxMP(mp, ag int) {
	p.Push(&model.MsgMPReply{Position: -2, MP: mp, AG: ag})
}

func (p *Player) PushMPAG(mp, ag int) {
	p.Push(&model.MsgMPReply{Position: -1, MP: mp, AG: ag})
}

func (p *Player) pushItemDurability(position, dur int) {
	p.Push(&model.MsgItemDurabilityReply{Position: position, Durability: dur, Flag: 1})
}

func (p *Player) pushDeleteItem(position int) {
	p.Push(&model.MsgDeleteInventoryItemReply{Position: position, Flag: 1})
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

func (p *Player) SpawnPosition() {
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
	p.CreateFrustrum()
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
		p.PushViewport(&reply)
	}
}

func (p *Player) Whisper(msg *model.MsgWhisper) {
	if len(msg.Name) == 0 {
		return
	}
	if p.Name == msg.Name {
		return
	}
	tobj := object.ObjectManager.GetPlayerByName(msg.Name)
	if tobj == nil {
		reply := model.MsgWhisperReplyFailed{
			Flag: 0,
		}
		p.Push(&reply)
		return
	}
	reply := model.MsgWhisperReply{}
	reply.Name = p.Name
	reply.Msg = msg.Msg
	tobj.Push(&reply)
}

// func (p *Player) Live(msg *model.MsgLive) {

// }

func (p *Player) Login(msg *model.MsgLogin) {
	// validate msg
	resp := model.MsgLoginReply{Result: 1}
	defer p.Push(&resp)
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
	p.ConnectState = object.ConnectStateLogged

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
// 	defer p.Push(&resp)
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
	defer p.Push(&model.MsgLogoutReply{Flag: msg.Flag})
	switch msg.Flag {
	case 0: // close game
	case 1: // back to pick character
		// offline to login state.\
		p.ConnectState = object.ConnectStateLogged
		p.SaveCharacter()
		p.Reset()
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
	p.Push(&model.MsgEnableCharacterClassReply{
		Class: reply.EnableCharacterClass,
	})
	p.Push(&model.MsgResetCharacterReply{
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
	p.Push(&reply)
}

func (p *Player) CreateCharacter(msg *model.MsgCreateCharacter) {
	reply := model.MsgCreateCharacterReply{Result: 0}
	defer p.Push(&reply)

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
	defer p.Push(&reply)

	if p.ConnectState == object.ConnectStatePlaying {
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
	p.Push(&model.MsgCheckCharacterReply{
		Result: 0,
	})
}

func (p *Player) DefineMuKey(msg *model.MsgDefineMuKey) {
	if p.ConnectState != object.ConnectStatePlaying {
		return
	}
	err := model.DB.UpdateCharacterMuKey(p.CharacterID, msg)
	if err != nil {
		log.Printf("model.DB.UpdateCharacterMuKey failed [err]%v\n", err)
		return
	}
}

func (p *Player) DefineMuBot(msg *model.MsgDefineMuBot) {
	err := model.DB.UpdateCharacterMuBot(p.CharacterID, msg)
	if err != nil {
		log.Printf("model.DB.UpdateCharacterMuBot failed [err]%v\n", err)
		return
	}
}

func (p *Player) EnableMuBot(msg *model.MsgEnableMuBot) {
	reply := model.MsgEnableMuBotReply{Flag: msg.Flag}
	defer p.Push(&reply)
	switch msg.Flag {
	case 0:
	case 1:
	}
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
	p.CharacterID = c.ID
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
	p.Skills = c.Skills
	p.Skills.FillSkillData(p.Class)
	p.Inventory = c.Inventory
	p.InventoryExpansion = c.InventoryExpansion
	p.Money = c.Money
	p.MapNumber = c.MapNumber
	p.X, p.TX = c.X, c.X
	p.Y, p.TY = c.Y, c.Y
	p.Dir = c.Dir
	if p.Level <= 10 {
		p.SpawnPosition()
	}
	p.CreateFrustrum()
	p.MoveSpeed = 1000
	p.MaxRegenTime = 4 * time.Second
	p.ConnectState = object.ConnectStatePlaying
	p.Live = true
	p.State = 1

	// p.Push(&model.MsgResetGameReply{})
	// reply
	p.Push(&model.MsgLoadCharacterReply{
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
	p.Push(&model.MsgItemListReply{
		Items: p.Inventory.Items,
	})
	// reply master data
	p.Push(&model.MsgMasterDataReply{
		MasterLevel:          p.MasterLevel,
		MasterExperience:     p.MasterExperience,
		MasterNextExperience: p.MasterNextExperience,
		MasterPoint:          p.MasterPoint,
	})
	p.pushSkillList()
	p.pushMasterSkillList()

	p.Push(&model.MsgMuKeyReply{
		MsgMuKey: c.MuKey,
	})
	p.Push(&model.MsgMuBotReply{
		MsgMuBot: c.MuBot,
	})

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
	helper := p.Inventory.Items[8]

	p.AddStrength = 0
	p.AddDexterity = 0
	p.AddVitality = 0
	p.AddEnergy = 0
	p.AddLeadership = 0
	p.CriticalAttackRate = 0
	p.ExcellentAttackRate = 0
	p.MonsterDieGetHP = 0
	p.MonsterDieGetMP = 0
	p.MonsterDieGetMoney = 0
	p.ArmorReduceDamage = 0
	p.ArmorReflectDamage = 0
	p.IgnoreDefenseRate = 0
	p.ReturnDamage = 0
	p.RecoverMaxHP = 0
	p.RecoverMaxMP = 0

	// wing item contribution
	if wing != nil && wing.ExcellentWing2Leadership {
		p.AddLeadership += 10 + 5*wing.Level
	}
	// set item contribution
	p.CalcSetItem(true)
	// master skill contribution

	strength := p.Strength + p.AddStrength
	dexterity := p.Dexterity + p.AddDexterity
	vitality := p.Vitality + p.AddVitality
	energy := p.Energy + p.AddEnergy
	leadership := p.Leadership + p.AddLeadership
	level := p.Level + p.MasterLevel

	// base attack
	leftAttackMin, leftAttackMax := 0, 0
	rightAttackMin, rightAttackMax := 0, 0
	switch class.Class(p.Class) {
	case class.Wizard:
		formula.WizardDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.WizardMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.Knight:
		formula.KnightDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.KnightMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.Elf:
		if leftHand != nil && leftHand.KindB == item.KindBCrossbow ||
			rightHand != nil && rightHand.KindB == item.KindBBow {
			formula.ElfWithBowDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		} else {
			formula.ElfWithoutBowDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		}
		formula.ElfMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.Magumsa:
		formula.GladiatorDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.GladiatorMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.DarkLord:
		formula.LordDamageCalc(strength, dexterity, vitality, energy, leadership, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.LordMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.Summoner:
		formula.SummonerDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.SummonerMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax, &p.curseAttackMin, &p.curseAttackMax)
	case class.RageFighter:
		formula.RageFighterDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.RageFighterMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	}

	// base defense
	formula.CalcDefense(p.Class, dexterity, &p.Defense)

	// base attack/defense success rate
	formula.CalcAttackSuccessRate_PvM(p.Class, strength, dexterity, leadership, p.Level, &p.AttackRate)
	formula.CalcAttackSuccessRate_PvP(p.Class, dexterity, p.Level, &p.attackRatePVP)
	formula.CalcDefenseSuccessRate_PvM(p.Class, dexterity, &p.DefenseRate)
	formula.CalcDefenseSuccessRate_PvP(p.Class, dexterity, p.Level, &p.defenseRatePVP)

	// base speed
	formula.CalcAttackSpeed(p.Class, dexterity, &p.AttackSpeed, &p.magicSpeed)

	// Stat Specialization
	// calc Stat Specialization: increase attack power
	// http://muonline.webzen.com/guides/219/1976/season-9/season-9-character-renewal
	// - (Bonus stat generated by equipping items like weapon, armor, wing, or master skill,
	// etc is not applied to specialization calculation.)
	var options []*model.MsgStatSpec
	var percent float64

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
	min = p.magicAttackMin * int(percent) / 100
	max = p.magicAttackMax * int(percent) / 100
	p.magicAttackMin += min
	p.magicAttackMax += max
	options = append(options, &model.MsgStatSpec{ID: 9, Min: min, Max: max})

	// STAT_OPTION_INC_CURSE_DAMAGE
	formula.StatSpec_GetPercent(p.Class, 10, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.curseAttackMin * int(percent) / 100
	max = p.curseAttackMax * int(percent) / 100
	p.curseAttackMin += min
	p.curseAttackMax += max
	options = append(options, &model.MsgStatSpec{ID: 10, Min: min, Max: max})

	// STAT_OPTION_INC_DEFENSE
	formula.StatSpec_GetPercent(p.Class, 4, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.Defense * int(percent) / 100
	p.Defense += min
	options = append(options, &model.MsgStatSpec{ID: 4, Min: min})

	// STAT_OPTION_INC_ATTACK_RATE
	formula.StatSpec_GetPercent(p.Class, 2, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.AttackRate * int(percent) / 100
	p.AttackRate += min
	options = append(options, &model.MsgStatSpec{ID: 2, Min: min})

	// STAT_OPTION_INC_ATTACK_RATE_PVP
	formula.StatSpec_GetPercent(p.Class, 3, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.attackRatePVP * int(percent) / 100
	p.attackRatePVP += min
	options = append(options, &model.MsgStatSpec{ID: 3, Min: min})

	// STAT_OPTION_INC_DEFENSE_RATE
	formula.StatSpec_GetPercent(p.Class, 6, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.DefenseRate * int(percent) / 100
	p.DefenseRate += min
	options = append(options, &model.MsgStatSpec{ID: 6, Min: min})

	// STAT_OPTION_INC_DEFENSE_RATE_PVP
	formula.StatSpec_GetPercent(p.Class, 7, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.defenseRatePVP * int(percent) / 100
	p.defenseRatePVP += min
	options = append(options, &model.MsgStatSpec{ID: 7, Min: min})

	// weapon and weapon addition attack
	if leftHand != nil {
		leftAttackMin += leftHand.AttackMin + leftHand.AdditionAttack
		leftAttackMax += leftHand.AttackMax + leftHand.AdditionAttack
		p.magicAttackMin += leftHand.AdditionMagicAttack
		p.magicAttackMax += leftHand.AdditionMagicAttack
		p.curseAttackMin += leftHand.AdditionCurseAttack
		p.curseAttackMax += leftHand.AdditionCurseAttack
		p.curse = leftHand.Magic
	}
	if rightHand != nil {
		rightAttackMin += rightHand.AttackMin + rightHand.AdditionAttack
		rightAttackMax += rightHand.AttackMax + rightHand.AdditionAttack
		p.magicAttackMin += rightHand.AdditionMagicAttack
		p.magicAttackMax += rightHand.AdditionMagicAttack
		p.curseAttackMin += rightHand.AdditionCurseAttack
		p.curseAttackMax += rightHand.AdditionCurseAttack
		p.magic = rightHand.Magic
	}

	// wing addition attack
	if wing != nil {
		leftAttackMin += wing.AdditionAttack
		leftAttackMax += wing.AdditionAttack
		rightAttackMin += wing.AdditionAttack
		rightAttackMax += wing.AdditionAttack
		p.magicAttackMin += wing.AdditionMagicAttack
		p.magicAttackMax += wing.AdditionMagicAttack
		p.curseAttackMin += wing.AdditionCurseAttack
		p.curseAttackMax += wing.AdditionCurseAttack
	}

	// glove speed
	if glove != nil {
		p.AttackSpeed += glove.AttackSpeed
		p.magicSpeed += glove.AttackSpeed
	}

	// armor(shield|armor|wing) addition defense(rate)
	for i := 1; i <= 7; i++ {
		it := p.Inventory.Items[i]
		if it != nil {
			p.Defense += it.Defense
			p.Defense += it.AdditionDefense
			p.DefenseRate += it.SuccessfulBlocking
			p.DefenseRate += it.AdditionDefenseRate
			p.RecoverHP = it.AdditionRecoverHP
		}
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
			n := 0
			switch {
			case level15Count == 5:
				n = 30
			case level14Count == 5:
				n = 25
			case level13Count == 5:
				n = 20
			case level12Count == 5:
				n = 15
			case level11Count == 5:
				n = 10
			case level10Count == 5:
				n = 5
			}
			p.Defense += p.Defense * n / 100
			p.DefenseRate += p.DefenseRate / 10
		}
	}

	// hp/mp
	c := CharacterTable[p.Class]
	p.MaxHP = c.HP + int(float32(level-1)*c.LevelHP) + int(float32(vitality-c.Vitality)*c.VitalityHP)
	p.MaxMP = c.MP + int(float32(level-1)*c.LevelMP) + int(float32(energy-c.Energy)*c.EnergyMP)

	// sd
	sdGageConstA := conf.CommonServer.GameServerInfo.SDGageConstA
	sdGageConstB := conf.CommonServer.GameServerInfo.SDGageConstB
	expressionA := strength + dexterity + vitality + energy
	if p.Class == int(class.DarkLord) {
		expressionA += leadership
	}
	expressionB := level * level / sdGageConstB
	p.MaxSD = expressionA*sdGageConstA/10 + expressionB + p.Defense

	// ag
	f := func(s, d, v, e, l float32) int {
		return int(float32(strength)*s + float32(dexterity)*d + float32(vitality)*v + float32(energy)*e + float32(leadership)*l)
	}
	switch class.Class(p.Class) {
	case class.Wizard:
		p.MaxAG = f(0.2, 0.4, 0.3, 0.2, 0)
	case class.Knight:
		p.MaxAG = f(0.15, 0.2, 0.3, 1.0, 0)
	case class.Elf:
		p.MaxAG = f(0.3, 0.2, 0.3, 0.2, 0)
	case class.Magumsa:
		p.MaxAG = f(0.2, 0.25, 0.3, 0.15, 0)
	case class.DarkLord:
		p.MaxAG = f(0.3, 0.2, 0.1, 0.15, 0.3)
	case class.Summoner:
		p.MaxAG = f(0.2, 0.25, 0.3, 0.15, 0)
	case class.RageFighter:
		p.MaxAG = f(0.15, 0.2, 0.3, 1.0, 0)
	}

	p.CalcSetItem(false)
	// set item contribution
	leftAttackMin += p.IncreaseAttackMin
	leftAttackMax += p.IncreaseAttackMax
	rightAttackMin += p.IncreaseAttackMin
	rightAttackMax += p.IncreaseAttackMax

	// excellent item contribution
	for i := 0; i < 12; i++ {
		it := p.Inventory.Items[i]
		if it != nil {
			if it.Lucky {
				p.CriticalAttackRate += 5
			}
			if it.ExcellentAttackRate {
				p.ExcellentAttackRate += 10
			}
			if it.ExcellentAttackLevel {
				value := level / 20
				switch i {
				case 0:
					leftAttackMin += value
					leftAttackMax += value
				case 1:
					rightAttackMin += value
					rightAttackMax += value
				}
				p.magicAttackMin += value
				p.magicAttackMax += value
			}
			if it.ExcellentAttackPercent {
				switch i {
				case 0:
					leftAttackMin += leftAttackMin * 2 / 100
					leftAttackMax += leftAttackMax * 2 / 100
				case 1:
					rightAttackMin += rightAttackMin * 2 / 100
					rightAttackMax += rightAttackMax * 2 / 100
				}
				p.magicAttackMin += p.magicAttackMin * 2 / 100
				p.magicAttackMax += p.magicAttackMax * 2 / 100
			}
			if it.ExcellentAttackSpeed {
				p.AttackSpeed += 7
				p.magicSpeed += 7
			}
			if it.ExcellentAttackHP {
				p.MonsterDieGetHP++
			}
			if it.ExcellentAttackMP {
				p.MonsterDieGetMP++
			}
			if it.ExcellentDefenseHP {
				p.MaxHP += p.MaxHP * 4 / 100
			}
			if it.ExcellentDefenseMP {
				p.MaxMP += p.MaxMP * 4 / 100
			}
			if it.ExcellentDefenseReduce {
				p.ArmorReduceDamage += 4
			}
			if it.ExcellentDefenseReflect {
				p.ArmorReflectDamage += 5
			}
			if it.ExcellentDefenseRate {
				p.DefenseRate += p.DefenseRate * 10 / 100
			}
			if it.ExcellentDefenseMoney {
				p.MonsterDieGetMoney += 30
			}
			if it.ExcellentWing2Speed {
				p.AttackSpeed += 5
				p.magicSpeed += 5
			}
			if it.ExcellentWing2AG {
				p.MaxAG += 50
			}
			if it.ExcellentWing2Ignore {
				p.IgnoreDefenseRate += 3
			}
			if it.ExcellentWing2MP {
				p.MaxMP += 50 + 5*it.Level
			}
			if it.ExcellentWing2HP {
				p.MaxHP += 50 + 5*it.Level
			}
			if it.ExcellentWing3Ignore {
				p.IgnoreDefenseRate += 5
			}
			if it.ExcellentWing3Return {
				p.ReturnDamage += 5
			}
			if it.ExcellentWing3HP {
				p.RecoverMaxHP += 5
			}
			if it.ExcellentWing3MP {
				p.RecoverMaxMP += 5
			}
		}
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
	case left && right:
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
		p.AttackMin = leftAttackMin + rightAttackMin
		p.AttackMax = leftAttackMax + rightAttackMax
		p.AttackSpeed += (leftHand.AttackSpeed + rightHand.AttackSpeed) / 2
		p.magicSpeed += (leftHand.AttackSpeed + rightHand.AttackSpeed) / 2
	case left:
		p.AttackMin = leftAttackMin
		p.AttackMax = leftAttackMax
		p.AttackSpeed += leftHand.AttackSpeed
		p.magicSpeed += leftHand.AttackSpeed
	case right:
		p.AttackMin = rightAttackMin
		p.AttackMax = rightAttackMax
		p.AttackSpeed += rightHand.AttackSpeed
		p.magicSpeed += rightHand.AttackSpeed
	default:
		p.AttackMin = (leftAttackMin + rightAttackMin) / 2
		p.AttackMax = (leftAttackMax + rightAttackMax) / 2
	}

	// convert magic/curse to magic/curse attack
	p.magicAttackMin += p.magicAttackMax * p.magic / 100
	p.magicAttackMax += p.magicAttackMax * p.magic / 100
	p.curseAttackMin += p.curseAttackMin * p.curseAttackMax / 100
	p.curseAttackMax += p.curseAttackMin * p.curseAttackMax / 100

	// wing
	if wing != nil {
		wingIncreaseDamage := 0
		formula.Wings_CalcIncAttack(wing.Code, wing.Level, &wingIncreaseDamage)
		p.WingIncreaseDamage = wingIncreaseDamage - 100
		wingReduceDamage := 0
		formula.Wings_CalcAbsorb(wing.Code, wing.Level, &wingReduceDamage)
		p.WingReduceDamage = 100 - wingReduceDamage
	}

	// helper item contribution
	switch {
	case p.equippedItem(helper, item.Code(13, 1)): // 小恶魔
		percent := conf.PetRing.Pets.Imp.AddAttackPercent
		p.AttackMin += p.AttackMin * percent / 100
		p.AttackMax += p.AttackMax * percent / 100
		p.magicAttackMin += p.magicAttackMin * percent / 100
		p.magicAttackMax += p.magicAttackMax * percent / 100
		p.curseAttackMin += p.curseAttackMin * percent / 100
		p.curseAttackMax += p.curseAttackMax * percent / 100
	case p.equippedItem(helper, item.Code(13, 64)): // 强化恶魔
		percent := conf.PetRing.Pets.Demon.AddAttackPercent
		p.AttackMin += p.AttackMin * percent / 100
		p.AttackMax += p.AttackMax * percent / 100
		p.magicAttackMin += p.magicAttackMin * percent / 100
		p.magicAttackMax += p.magicAttackMax * percent / 100
		p.curseAttackMin += p.curseAttackMin * percent / 100
		p.curseAttackMax += p.curseAttackMax * percent / 100
		p.AttackSpeed += conf.PetRing.Pets.Demon.AddAttackSpeed
		p.magicSpeed += conf.PetRing.Pets.Demon.AddAttackSpeed
	case p.equippedItem(helper, item.Code(13, 0)): // 守护天使
		p.MaxHP += conf.PetRing.Pets.Angel.AddHP
		p.HelperReduceDamage = conf.PetRing.Pets.Angel.ReduceDamagePercent
	case p.equippedItem(helper, item.Code(13, 65)): // 强化天使
		p.MaxHP += conf.PetRing.Pets.SpiritAngel.AddHP
		p.HelperReduceDamage = conf.PetRing.Pets.SpiritAngel.ReduceDamagePercent
	case p.equippedItem(helper, item.Code(13, 80)): // 熊猫
		p.Defense += conf.PetRing.Pets.Panda.AddDefenseValue
	}

	// ...

	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}
	if p.MP > p.MaxMP {
		p.MP = p.MaxMP
	}
	if p.SD > p.MaxSD {
		p.SD = p.MaxSD
	}
	if p.AG > p.MaxAG {
		p.AG = p.MaxAG
	}

	// skill bonus
	KnightGladiatorCalcSkillBonus := 0.0
	formula.Knight_Gladiator_CalcSkillBonus(p.Class, energy, &KnightGladiatorCalcSkillBonus)
	p.KnightGladiatorCalcSkillBonus = KnightGladiatorCalcSkillBonus

	// Push
	p.Push(&model.MsgStatSpecReply{
		Options: options,
	})
	p.Push(&model.MsgAttackSpeedReply{
		AttackSpeed: p.AttackSpeed,
		MagicSpeed:  p.magicSpeed,
	})
	p.pushMaxHP(p.MaxHP+p.AddHP, p.MaxSD+p.AddSD)
	p.pushMaxMP(p.MaxMP+p.AddMP, p.MaxAG+p.AddAG)
	p.pushHP(p.HP, p.SD)
	p.PushMPAG(p.MP, p.AG)
}

func (p *Player) equippedItem(it *item.Item, code int) bool {
	if it == nil || it.Durability == 0 {
		return false
	}
	return it.Code == code
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
		p.Push(&reply)
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
		p.Push(&reply)
	})
}

func (p *Player) MapDataLoadingOK(msg *model.MsgMapDataLoadingOK) {}

func (p *Player) SaveCharacter() {
	if p.Name == "" {
		return
	}
	c := model.Character{
		ID:                 p.CharacterID,
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
		Skills:             p.Skills,
		Inventory:          p.Inventory,
		InventoryExpansion: p.InventoryExpansion,
		Money:              p.Money,
		MapNumber:          p.MapNumber,
		X:                  p.X,
		Y:                  p.Y,
		Dir:                p.Dir,
	}
	err := model.DB.UpdateCharacter(&c)
	if err != nil {
		log.Printf("model.DB.SaveCharacter failed [err]%v\n", err)
		return
	}
}

func (p *Player) GetPKLevel() int {
	return p.PKLevel
}

func (p *Player) GetMasterLevel() int {
	return p.MasterLevel
}

func (p *Player) GetSkillMPAG(s *skill.Skill) (int, int) {
	return s.ManaUsage, s.BPUsage
}

func (p *Player) addSetEffect(index item.SetEffectType, value int, base bool) {
	if base {
		switch index {
		case item.SetEffectIncStrength:
			p.AddStrength += value
		case item.SetEffectIncAgility:
			p.AddDexterity += value
		case item.SetEffectIncEnergy:
			p.AddEnergy += value
		case item.SetEffectIncVitality:
			p.AddVitality += value
		case item.SetEffectIncLeadership:
			p.AddLeadership += value
		}
	} else {
		switch index {
		case item.SetEffectIncAttackMin:
			p.IncreaseAttackMin += value
		case item.SetEffectIncAttackMax:
			p.IncreaseAttackMax += value
		case item.SetEffectIncMagicAttack:
			p.magicAttackMin += value
			p.magicAttackMax += value
		case item.SetEffectIncDamage:
			p.SetAddDamage += value
		case item.SetEffectIncAttackRate:
			p.AttackRate += value
		case item.SetEffectIncDefense:
			p.Defense += value
		case item.SetEffectIncMaxHP:
			p.MaxHP += value
		case item.SetEffectIncMaxMP:
			p.MaxMP += value
		case item.SetEffectIncMaxAG:
			p.MaxAG += value
		case item.SetEffectIncAG:
			p.MaxAG += value
		case item.SetEffectIncCritiDamageRate:
			p.CriticalAttackRate += value
		case item.SetEffectIncCritiDamage:
			p.CriticalAttackDamage += value
		case item.SetEffectIncExcelDamageRate:
			p.ExcellentAttackRate += value
		case item.SetEffectIncExcelDamage:
			p.ExcellentAttackDamage += value
		case item.SetEffectIncSkillAttack:
			p.IncreaseSkillAttack += value
		case item.SetEffectDoubleDamage:
			p.DoubleDamageRate += value
		case item.SetEffectIgnoreDefense:
			p.IgnoreDefenseRate += value
		case item.SetEffectIncShieldDefense:
			p.MaxSD += value
		case item.SetEffectIncTwoHandSwordDamage:
			p.IncreaseTwoHandWeaponDamage += value
		}
	}
}

func (p *Player) CalcSetItem(base bool) {
	type set struct {
		index int
		count int
	}
	var sets []set

	sameWeapon := 0
	sameRing := 0
	for i, it := range p.Inventory.Items[0:object.InventoryWearSize] {
		if it == nil {
			continue
		}
		if it.Durability == 0 {
			continue
		}
		tierIndex := it.GetSetTierIndex()
		if tierIndex == 0 {
			continue
		}
		index := item.SetManager.GetSetIndex(it.Section, it.Index, tierIndex)
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
				p.addSetEffect(setEffect.Index, setEffect.Value, base)
			}

			if count > item.SetManager.GetSetEffectCount(index) {
				p.setFull = true
				setEffects := item.SetManager.GetSetFull(index)
				for _, setEffect := range setEffects {
					p.addSetEffect(setEffect.Index, setEffect.Value, base)
				}
			}
		}
	}
}

func (player *Player) Calc380Item() {
	for _, wItem := range player.Inventory.Items[0:object.InventoryWearSize] {
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
	p.Skills.ForEachActiveSkill(func(s *skill.Skill) {
		skills = append(skills, s)
	})
	sort.Sort(skill.SortedSkillSlice(skills))
	p.Push(&model.MsgSkillListReply{Skills: skills})
}

func (p *Player) pushMasterSkillList() {
	var skills []*model.MsgMasterSkill
	p.Skills.ForEachMasterSkill(func(index int, level int, curValue float32, nextValue float32) {
		skills = append(skills, &model.MsgMasterSkill{
			MasterSkillUIIndex:   index,
			MasterSkillLevel:     level,
			MasterSkillCurValue:  curValue,
			MasterSkillNextValue: nextValue,
		})
	})
	p.Push(&model.MsgMasterSkillListReply{Skills: skills})
}

func (p *Player) LearnMasterSkill(msg *model.MsgLearnMasterSkill) {
	reply := model.MsgLearnMasterSkillReply{
		Result:           0,
		MasterPoint:      p.MasterPoint,
		MasterSkillIndex: -1,
	}
	defer p.Push(&reply)

	if p.MasterLevel <= 0 ||
		p.MasterPoint <= 0 {
		return
	}

	p.Skills.GetMaster(p.Class, msg.SkillIndex, p.MasterPoint, func(point, uiIndex, index, level int, curValue, NextValue float32) {
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
		Index:  p.Index,
		Action: msg.Action,
		Dir:    msg.Dir,
	}
	p.PushViewport(&reply)
}

func (p *Player) Process1000ms() {
	if p.ConnectState == object.ConnectStatePlaying {
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
		p.PushMPAG(p.MP, p.AG)
	}
}

func (p *Player) Die(obj *object.Object) {

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
	p.Push(&reply)
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
		p.Push(&reply)
		if p.MapNumber != mapNumber {
			p.MapNumber = mapNumber
			p.loadMiniMap()
		}
		p.X, p.Y = x, y
		p.TX, p.TY = x, y
		p.Dir = dir
		p.CreateFrustrum()
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
	for _, s := range p.Skills {
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
		if s, ok := p.LearnSkill(index); ok {
			p.Push(&model.MsgSkillOneReply{
				Flag:  -2,
				Skill: s,
			})
		}
	}
	for _, index := range needForgetSkills {
		if s, ok := p.ForgetSkill(index); ok {
			p.Push(&model.MsgSkillOneReply{
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
		p.Push(&reply)
		if itemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   position,
				Durability: it.Durability,
				Flag:       0,
			}
			p.Push(&reply)
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
	defer p.Push(&reply)
	// validate
	if msg.Position >= len(p.Inventory.Items) {
		return
	}
	it := p.Inventory.Items[msg.Position]
	if it == nil {
		return
	}
	switch it.Code {
	case item.Code(14, 63):
		cmdReply := model.MsgServerCMDReply{
			Type: 0,
			X:    p.X,
			Y:    p.Y,
		}
		p.PushViewport(&cmdReply)
	default:
		ok := maps.MapManager.PushItem(p.MapNumber, msg.X, msg.Y, it)
		if !ok {
			return
		}
	}
	p.Inventory.DropItem(msg.Position, it)
	if msg.Position < 12 || msg.Position == 126 {
		p.inventoryChanged()
	}
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
		p.Push(&reply)
		if sitemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   msg.SrcPosition,
				Durability: sitem.Durability,
				Flag:       0,
			}
			p.Push(&reply)
		}
		if titemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   msg.DstPosition,
				Durability: titem.Durability,
				Flag:       0,
			}
			p.Push(&reply)
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
					p.Push(&reply)
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
			p.PushMPAG(p.MP, p.AG)
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
		if s, ok := p.LearnSkill(skillIndex); ok {
			p.Push(&model.MsgSkillOneReply{
				Flag:  -2,
				Skill: s,
			})
			p.Inventory.DropItem(msg.SrcPosition, it)
			p.Push(&model.MsgDeleteInventoryItemReply{
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
	// if msg.Target >= len(object.ObjectManager.objects) {
	// 	return
	// }
	// tobj := object.ObjectManager.objects[msg.Target]
	tobj := object.ObjectManager.GetObject(msg.Target)
	if tobj == nil {
		return
	}
	if math.Abs(float64(p.X-tobj.X)) > 5 ||
		math.Abs(float64(p.Y-tobj.Y)) > 5 {
		return
	}
	p.TargetNumber = tobj.Index
	switch tobj.NpcType {
	case object.NpcTypeShop:
		inventory := shop.ShopManager.GetShopInventory(tobj.Class, tobj.MapNumber)
		if tobj.Class == 492 {
			reply.Result = 34
		}
		p.Push(&reply)
		shopItemListReply := model.MsgTypeItemListReply{Type: 0}
		shopItemListReply.Items = inventory
		p.Push(&shopItemListReply)
	case object.NpcTypeWarehouse:
		account, err := model.DB.GetAccountByID(p.AccountID)
		if err != nil {
			log.Printf("Talk model.DB.GetAccountByID failed [err]%v\n", err)
			return
		}
		reply.Result = 2
		p.Push(&reply)
		p.Warehouse = account.Warehouse
		p.WarehouseMoney = account.WarehouseMoney
		warehouseItemListReply := model.MsgTypeItemListReply{Type: 0}
		warehouseItemListReply.Items = account.Warehouse.Items
		p.Push(&warehouseItemListReply)
		warehouseMoneyReply := model.MsgWarehouseMoneyReply{
			Result:         1,
			WarehouseMoney: p.WarehouseMoney,
			InventoryMoney: p.Money,
		}
		p.Push(&warehouseMoneyReply)
	case object.NpcTypeChaosMix:
	case object.NpcTypeGoldarcher:
	case object.NpcTypePentagramMix:
	}
}

func (p *Player) CloseTalkWindow(msg *model.MsgCloseTalkWindow) {
	p.TargetNumber = -1
}

func (p *Player) BuyItem(msg *model.MsgBuyItem) {
	reply := model.MsgBuyItemReply{
		Result: -1,
	}
	itemDurChanged := false
	position := -1
	var it *item.Item
	defer func() {
		p.Push(&reply)
		if itemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   position,
				Durability: it.Durability,
				Flag:       0,
			}
			p.Push(&reply)
		}
	}()
	// validate
	if msg.Position < 0 ||
		msg.Position >= shop.MaxShopItemCount ||
		p.TargetNumber < 0 {
		return
	}
	tobj := object.ObjectManager.GetObject(p.TargetNumber)
	if tobj == nil {
		return
	}
	if math.Abs(float64(p.X-tobj.X)) > 5 ||
		math.Abs(float64(p.Y-tobj.Y)) > 5 {
		return
	}
	if tobj.NpcType != object.NpcTypeShop {
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
	defer p.Push(&reply)
	// validate
	if msg.Position < 0 ||
		msg.Position >= p.Inventory.Size ||
		p.TargetNumber < 0 {
		return
	}
	tobj := object.ObjectManager.GetObject(p.TargetNumber)
	if tobj == nil {
		return
	}
	if math.Abs(float64(p.X-tobj.X)) > 5 ||
		math.Abs(float64(p.Y-tobj.Y)) > 5 {
		return
	}
	if tobj.NpcType != object.NpcTypeShop {
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
	p.Push(&reply)
}

func (p *Player) MapMove(msg *model.MsgMapMove) {
	reply := model.MsgMapMoveReply{
		Result: 0,
	}
	defer p.Push(&reply)
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
			p.Push(&reply)
		}
	})
}
