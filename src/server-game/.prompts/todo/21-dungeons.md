# 21. 副本系统

副本系统拥有 BloodCastle、DevilSquare、ChaosCastle、IllusionTemple、ImperialGuardian、DoppelGanger 等玩法状态；跨服入场认证、DataServer 积分/排名保存、MapServer 迁移通道归 `28-external-comm.md`。

本模块覆盖 BloodCastle、DevilSquare、ChaosCastle、IllusionTemple、ImperialGuardian、DoppelGanger 这类限时、入场、实例化或半实例化副本玩法。副本系统不拥有对象、地图、道具、经验、掉落、组队等基础能力，而是通过统一副本上下文接口被这些系统调用：对象系统通知玩家/怪物行为，地图系统提供地图和阻挡能力，组队系统提供入场授权，经验/掉落系统消费副本修正和奖励上下文。DoppelGanger 的副本状态、怪群和奖励业务归本模块，Lua VM、Lua 回调和绑定基础设施归 `26-script.md`。CastleSiege、Crywolf、Kanturu、Raklion、Arca、Acheron 暂归世界事件候选，不放入本模块核心范围。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 模块边界与总入口 | DungeonManager 总管理器 | `g_BloodCastle`、`g_DevilSquare`、`g_ChaosCastle`、`g_IT_Event`、`g_ImperialGuardian`、`g_DoppelGanger` | 暂无 `game/dungeon` | 未覆盖 | 建立统一副本服务，负责注册各副本实现和对外接口。 |
| 2 | 模块边界与总入口 | 副本类型枚举 | BC/DS/CC/IT/IG/DG 常量 | 地图枚举存在，副本类型无 | 未覆盖 | 定义副本类型，避免业务层直接用地图号判断玩法。 |
| 3 | 模块边界与总入口 | 副本实例标识 | bridge/floor/temple/zone 等索引 | 对象注释字段有部分索引 | 未覆盖 | 统一表达副本类型、层级、实例、队伍和场次。 |
| 4 | 模块边界与总入口 | 副本启停配置 | `IsEventEnable`、event XML | `conf.Events` 已读取 server 节点 | 部分覆盖 | 将 `events.xml` 中各副本 server 配置接入启停判断。 |
| 5 | 模块边界与总入口 | 副本时间表 | `*_START_TIME`、脚本加载 | `conf.Events` 结构较粗 | 未覆盖 | 加载开放时间、准备时间、进行时间、结算时间。 |
| 6 | 模块边界与总入口 | 副本 Tick 调度 | `Run`、`Proc*` | 运行时无副本 tick | 未覆盖 | 接入 server-game 主循环或定时器，按状态机推进副本。 |
| 7 | 模块边界与总入口 | GM/调试入口 | CheatOpenTime、WarpZoneGM、SetDayOfWeekGM | 暂无 | 未覆盖 | 预留 GM 开启、跳阶段、清场、传送和诊断入口。 |
| 8 | 通用副本状态机 | 通用状态模型 | CLOSED/OPEN/PLAYING/PLAYEND/NONE | 暂无 | 未覆盖 | 抽象准备、入场、进行、结束、清理等通用状态。 |
| 9 | 通用副本状态机 | 状态切换通知 | `SendNoticeState`、state packet | 多数 handler 仅占位 | 未覆盖 | 状态变化时下发倒计时、开始、结束、关闭消息。 |
| 10 | 通用副本状态机 | 剩余时间查询 | `GetRemainTime`、`GetCurrentRemainSec` | `0x91` 等入口占位 | 未覆盖 | 给客户端和其他模块查询副本剩余时间。 |
| 11 | 通用副本状态机 | 开放前广播 | before enter/play/end/quit flags | 暂无 | 未覆盖 | 实现进入前、开始前、结束前、踢出前的阶段通知。 |
| 12 | 通用副本状态机 | 无效用户清理 | `CheckInvalidUser`、`KickInvalidUser` | 暂无 | 未覆盖 | 周期检查离线、地图错误、状态错误用户并清理。 |
| 13 | 通用副本状态机 | 副本结束清理 | `ClearMonster`、`ItemClear`、kick flow | 暂无 | 未覆盖 | 清理怪物、地图物品、玩家状态、临时变量和动态阻挡。 |
| 14 | 通用副本状态机 | 副本失败判定 | every user die、timeout、mission fail | 暂无 | 未覆盖 | 超时、全员死亡、任务失败时统一进入失败结算。 |
| 15 | 通用副本状态机 | 副本成功判定 | mission success、winner exist | 暂无 | 未覆盖 | 按副本实现判定通关、获胜或最终存活者。 |
| 16 | 入场与资格 | 入场协议分发 | `0x90`、`0x9A`、`0xAF01`、`0xBF0E`、`0xDBxx` | `handle/c1c2.go` 有占位 | 部分覆盖 | 将副本入场协议委托到副本服务。 |
| 17 | 入场与资格 | 地图/Gate 入场 | `GetDevilSquareIndex`、Gate 常量 | `maps` 和 `MapMove` 有基础能力 | 部分覆盖 | 使用地图系统 Gate 能力传送到副本地图。 |
| 18 | 入场与资格 | 等级段校验 | `CheckEnterLevel`、level tables | 暂无副本等级规则 | 未覆盖 | 按职业、普通/特殊职业、大师等级判断层级。 |
| 19 | 入场与资格 | 门票校验 | `CheckEnterItem`、`CheckEnterFreeTicket` | 道具类型存在 | 未覆盖 | 检查门票类型、等级、位置、数量和免费入场券。 |
| 20 | 入场与资格 | 门票消耗 | entry item consume | 背包操作存在但未接入 | 未覆盖 | 入场成功后原子删除或扣除门票。 |
| 21 | 入场与资格 | 入场次数限制 | `reqEventEnterCount` | `0x9F` handler 占位 | 未覆盖 | 查询和限制每日/账号/角色副本入场次数。 |
| 22 | 入场与资格 | PK 入场限制 | `CanEnterEventWithPK` | 配置已读取 | 部分覆盖 | 按配置禁止 PK 状态进入副本。 |
| 23 | 入场与资格 | 人数上限 | `MAX_*_USER`、max user config | 暂无 | 未覆盖 | 控制每层、每实例、每队伍最大人数。 |
| 24 | 入场与资格 | 最小人数 | `MIN_CC_USER_NEED_PLAY`、IT party min | 暂无 | 未覆盖 | 达不到最小人数时不开始或直接失败。 |
| 25 | 入场与资格 | 互斥状态检查 | `m_IfState`、trade/shop/chaos checks | 各模块待接入 | 未覆盖 | 交易、个人商店、合成、仓库等状态下禁止入场。 |
| 26 | 入场与资格 | 组队入场授权 | `EnterITR_PartyAuth`、DSF party enter | `15-party.md` 已记录边界 | 未覆盖 | 组队副本需要全员确认、队伍资格和授权清理。 |
| 27 | 玩家生命周期 | 玩家加入副本 | `EnterUserBridge`、`AddUser`、`AddDoppelgangerUser` | 暂无 | 未覆盖 | 记录玩家到副本实例，并设置对象副本上下文。 |
| 28 | 玩家生命周期 | 玩家离开副本 | `LeaveUserBridge`、`LeaveDevilSquare`、`LeaveDoppelganger` | 暂无 | 未覆盖 | 主动离开、回城、换图时清理副本成员。 |
| 29 | 玩家生命周期 | 玩家断线清理 | `CheckUsersOnConnect`、invalid user flow | `03-characters.md` 有下线边界 | 未覆盖 | 断线时按副本规则保留、踢出或失败结算。 |
| 30 | 玩家生命周期 | 玩家死亡处理 | `IllusionTempleUserDie`、CC/BC die checks | 对象死亡已有入口 | 未覆盖 | 对象死亡后通知副本系统判定复活、淘汰或失败。 |
| 31 | 玩家生命周期 | 玩家重生位置 | `UserDieRegen`、rollback pos | `SpawnPosition` 有活动地图例外 | 部分覆盖 | 副本内死亡、复活、失败回城应有独立规则。 |
| 32 | 玩家生命周期 | 玩家状态字段 | `BloodCastleSubIndex`、`ChaosCastleIndex` 等 | `object.Object` 有注释字段 | 未覆盖 | 恢复或重设对象上的副本运行时状态。 |
| 33 | 玩家生命周期 | 玩家积分/经验字段 | BC exp/score、CC score、DS score | 对象无正式字段 | 未覆盖 | 副本期间记录临时积分、经验、击杀、贡献。 |
| 34 | 玩家生命周期 | 玩家踢出 | all kick、quit msg | 暂无 | 未覆盖 | 结束、失败、非法状态、GM 操作时踢出副本。 |
| 35 | 玩家生命周期 | 重连状态恢复 | TemporaryUser/Event reconnect | `03-characters.md` 记录临时用户缺口 | 未覆盖 | 允许后续实现断线重连回副本或安全回城。 |
| 36 | 对象系统接口 | OnPlayerEnterMap | map event check | 对象进入地图无副本 hook | 未覆盖 | 玩家进入地图后，副本系统识别地图并绑定上下文。 |
| 37 | 对象系统接口 | OnPlayerLeaveMap | leave flow | 暂无 | 未覆盖 | 玩家离开副本地图时清理副本成员和状态。 |
| 38 | 对象系统接口 | OnPlayerDie | user die event funcs | `Object.Die` 路径存在 | 未覆盖 | 玩家死亡时由对象系统通知副本系统。 |
| 39 | 对象系统接口 | OnMonsterDie | `DieProcDevilSquare`、BC/CC kill count | 怪物死亡掉落入口存在 | 未覆盖 | 怪物死亡先推进副本积分/目标，再进入经验和掉落。 |
| 40 | 对象系统接口 | CanMove | `CheckWalk`、CC hollow zone | `Object.Move` 有基础移动 | 未覆盖 | 副本动态阻挡、塌陷区域、阶段门禁要裁决移动。 |
| 41 | 对象系统接口 | CanAttack | IT team/protection、IG attackable | 攻击入口存在 | 未覆盖 | 副本保护阶段、阵营、不可攻击怪物需要裁决攻击。 |
| 42 | 对象系统接口 | OnObjectViewport | BC/PS state notify 类似模式 | 视野系统基础存在 | 未覆盖 | 进入视野时同步副本内特殊对象、门、雕像、状态。 |
| 43 | 地图系统接口 | 副本地图识别 | `BC_MAP_RANGE`、`DS_MAP_RANGE`、`CC_MAP_RANGE` | `game/maps/const.go` 已有枚举 | 部分覆盖 | 地图系统提供 `IsDungeonMap` 和具体副本类型查询。 |
| 44 | 地图系统接口 | 动态阻挡设置 | `PMSG_SETMAPATTR`、`BlockCastleDoor` | 地图属性读取存在 | 未覆盖 | 支持副本中按阶段修改地图阻挡。 |
| 45 | 地图系统接口 | 副本地图掉落限制 | `IsCanNotItemDropInDevilSquare` 等 | `14-drops.md` 记录缺口 | 未覆盖 | 地图系统和掉落系统需要查询副本禁掉规则。 |
| 46 | 地图系统接口 | 副本离场回城 | rollback/move home | `MapMove` 有基础能力 | 未覆盖 | 失败、结束、断线恢复时回到安全地图或入口。 |
| 47 | 地图系统接口 | 副本移动命令限制 | `CheckMainToMove`、event map checks | `06-maps.md` 已记录 | 未覆盖 | 副本进行中禁止普通移动、传送或换线。 |
| 48 | 怪物与机关 | 怪物脚本加载 | `LoadMonster`、`LoadMonsterScript` | 怪物配置基础存在 | 未覆盖 | 加载副本怪物出生点、层级、波次、区域和门/机关对象。 |
| 49 | 怪物与机关 | 怪物生成 | `SetMonster`、`RegenMonster`、`AddMonsterHerd` | 怪物对象基础存在 | 未覆盖 | 按副本状态生成怪物、Boss、机关、宝箱。 |
| 50 | 怪物与机关 | 怪物清理 | `ClearMonster`、`DeleteMonster` | 对象管理器可删除对象 | 未覆盖 | 副本结束或阶段切换时清理副本怪。 |
| 51 | 怪物与机关 | 怪物死亡计数 | `SetMonsterKillCount`、`GetLiveMonsterCount` | 暂无 | 未覆盖 | 怪物死亡推进任务目标、波次和通关条件。 |
| 52 | 怪物与机关 | 副本怪 AI 边界 | DG herd、IG monster base act | `25-monster-ai.md` | 未覆盖 | 副本系统只触发特殊 AI，上层 AI 行为归怪物 AI 系统。 |
| 53 | 怪物与机关 | 门/雕像/机关对象 | BC door/statue、IG gate、IT relic NPC | NPC/怪物对象基础存在 | 未覆盖 | 建模可被攻击、可交互、可阻挡的副本机关。 |
| 54 | 怪物与机关 | 宝箱对象 | DG treasure、IG item bag statue | 暂无 | 未覆盖 | 结算或阶段后生成宝箱，并转给掉落/奖励接口。 |
| 55 | 经验与奖励接口 | 副本经验入口 | DS `gObjMonsterExpSingle`、BC `AddExperience` | `09-exp.md` 记录缺口 | 未覆盖 | 经验系统调用副本修正或副本直接结算经验。 |
| 56 | 经验与奖励接口 | 副本经验倍率 | BC 50%、CC bonus、DG rate | `ExpManager` 无副本接入 | 未覆盖 | 对普通/大师经验分别提供副本倍率。 |
| 57 | 经验与奖励接口 | 积分转经验 | DS score、BC reward exp、IG reward exp | 暂无 | 未覆盖 | 副本积分、通关、击杀转换为经验奖励。 |
| 58 | 经验与奖励接口 | 副本奖励 Zen | `CalcSendRewardZEN` | 金币字段存在 | 未覆盖 | 通关或结算按副本规则发 Zen。 |
| 59 | 经验与奖励接口 | 副本道具奖励 | `DropReward`、EventDungeonItemBag | `14-drops.md` 有副本掉落缺口 | 未覆盖 | 使用掉落系统/Bag 系统生成副本奖励。 |
| 60 | 经验与奖励接口 | 奖励结果下发 | `SendRewardScore`、mission result packet | 暂无 | 未覆盖 | 向客户端下发积分、经验、奖励、胜负结果。 |
| 61 | BloodCastle | BC 管理器 | `CBloodCastle` | 暂无 | 未覆盖 | 实现血色城堡独立副本实现。 |
| 62 | BloodCastle | BC 桥实例 | `BLOODCASTLE_BRIDGE` | 对象注释字段有 BC 索引 | 未覆盖 | 管理 8 个等级桥、子槽位、状态和成员。 |
| 63 | BloodCastle | BC 状态机 | `BC_STATE_CLOSED/PLAYING/PLAYEND` | 暂无 | 未覆盖 | 实现关闭、进行、结束和倒计时通知。 |
| 64 | BloodCastle | BC 入场等级 | `CheckEnterLevel` | 暂无 | 未覆盖 | 按斗篷等级/角色等级判断可进桥。 |
| 65 | BloodCastle | BC 门票检查 | `CheckEnterItem`、`CheckEnterFreeTicket` | 门票分类存在 | 未覆盖 | 检查并消耗透明斗篷或免费入场券。 |
| 66 | BloodCastle | BC 队伍入场 | `CheckCanParty`、`CheckPartyExist` | `15-party.md` 待接入 | 未覆盖 | 队伍进入和通关归属按 GameServer 规则处理。 |
| 67 | BloodCastle | BC 城门阻挡 | `BlockCastleDoor`、`ReleaseCastleDoor` | 地图动态阻挡未实现 | 未覆盖 | 阶段性阻挡城门、桥、入口。 |
| 68 | BloodCastle | BC 城门生命 | `m_iCastleDoorHealth` | 暂无 | 未覆盖 | 城门作为可攻击目标，死亡推进任务。 |
| 69 | BloodCastle | BC 雕像生命 | `m_iCastleStatueHealth` | 暂无 | 未覆盖 | 雕像作为可攻击目标，死亡后掉任务物。 |
| 70 | BloodCastle | BC 怪物击杀计数 | `SetMonsterKillCount` | 暂无 | 未覆盖 | 杀够怪物后开放下一阶段或发送提示。 |
| 71 | BloodCastle | BC Boss 击杀计数 | `CheckBossKillSuccess` | 暂无 | 未覆盖 | 处理最后阶段 Boss 击杀目标。 |
| 72 | BloodCastle | BC 任务物品 | quest item serial/user index | 暂无 | 未覆盖 | 掉落、拾取、持有、交付大天使武器。 |
| 73 | BloodCastle | BC 玩家行走限制 | `CheckWalk` | 移动无副本裁决 | 未覆盖 | 按城门、桥、阶段阻挡玩家移动。 |
| 74 | BloodCastle | BC 胜利判定 | `CheckWinnerExist/Valid` | 暂无 | 未覆盖 | 判断持任务物且交付成功的玩家或队伍。 |
| 75 | BloodCastle | BC 奖励 | `GiveReward_Win/Fail`、`DropReward` | 暂无 | 未覆盖 | 胜利和失败分别结算经验、Zen、道具、积分。 |
| 76 | DevilSquare | DS 管理器 | `CDevilSquare` | 暂无 | 未覆盖 | 实现恶魔广场独立副本实现。 |
| 77 | DevilSquare | DS 状态机 | CLOSE/OPEN/PLAYING/NONE | 暂无 | 未覆盖 | 实现关闭、开放入场、进行中、空状态。 |
| 78 | DevilSquare | DS 时间同步 | `CheckSync`、start times | 暂无 | 未覆盖 | 按时间表自动开启和关闭。 |
| 79 | DevilSquare | DS 入场层级 | `GetUserLevelToEnter` | 暂无 | 未覆盖 | 按等级、职业决定进入 DS1-7 和 Gate。 |
| 80 | DevilSquare | DS 移动入口 | `moveDevilSquare` | `handle` 0x90 占位 | 部分覆盖 | 完成移动请求、门票校验和传送。 |
| 81 | DevilSquare | DS 剩余时间 | `devilSquareRemainTime` | `handle` 0x91 占位 | 部分覆盖 | 返回当前开放/进行剩余时间。 |
| 82 | DevilSquare | DS Ground | `CDevilSquareGround` | 暂无 | 未覆盖 | 管理每个广场玩家、怪物数、积分和排名。 |
| 83 | DevilSquare | DS 怪物生成 | `SetMonster`、`gDevilSquareMonsterRegen` | 暂无 | 未覆盖 | 按广场和波次生成怪物。 |
| 84 | DevilSquare | DS 经验分流 | `gObjMonsterExpSingle/gObjExpParty` | `09-exp.md` 记录缺口 | 未覆盖 | DS 中使用专用经验和积分计算。 |
| 85 | DevilSquare | DS 积分结算 | `CalcScore` | 暂无 | 未覆盖 | 统计击杀、伤害、职业奖励和最终积分。 |
| 86 | DevilSquare | DS 掉落/材料 | eye/key drop rate | `14-drops.md` 有活动掉落缺口 | 未覆盖 | 门票材料掉落和 DS 内掉落规则接入掉落系统。 |
| 87 | ChaosCastle | CC 管理器 | `CChaosCastle` | 暂无 | 未覆盖 | 实现赤色要塞独立副本实现。 |
| 88 | ChaosCastle | CC 层级实例 | `CHAOSCASTLE_DATA` | 对象注释字段有 CC 索引 | 未覆盖 | 管理 7 个等级层、成员、怪物和状态。 |
| 89 | ChaosCastle | CC 入场等级 | `g_sttCHAOSCASTLE_LEVEL` | 暂无 | 未覆盖 | 按职业和等级判断可进层级。 |
| 90 | ChaosCastle | CC 入场费用 | `g_iChaosCastle_EnterCost` | 暂无 | 未覆盖 | 入场扣 Zen，并处理失败回滚。 |
| 91 | ChaosCastle | CC 入场协议 | `CGReqEnterChaosCastle` | `handle` 0xAF01 占位 | 部分覆盖 | 完成赤色要塞入场请求业务。 |
| 92 | ChaosCastle | CC 伪装/重定位 | reposition user | `handle` 0xAF02/0xAF06 占位 | 部分覆盖 | 入场后重定位、伪装和同步客户端。 |
| 93 | ChaosCastle | CC 地图塌陷 | hollow zone/trap step | 地图动态阻挡未实现 | 未覆盖 | 分阶段缩小地图并处理掉落/击退/死亡。 |
| 94 | ChaosCastle | CC 吹飞伤害 | blow out distance/damage | 暂无 | 未覆盖 | 角色被击退出区域时应用伤害和淘汰。 |
| 95 | ChaosCastle | CC 存活判定 | winner / live user count | 暂无 | 未覆盖 | 玩家和怪物混战后判定最后存活者。 |
| 96 | ChaosCastle | CC 怪物掉落 | monster item table | 暂无 | 未覆盖 | 副本怪物特殊掉落和获胜奖励接入掉落系统。 |
| 97 | ChaosCastle | CCF 生存版本 | `ChaosCastleFinal`、0xAF03/05/07/08 | handler 有占位 | 部分覆盖 | 记录 Survival/Final 版本边界，后续可独立子实现。 |
| 98 | IllusionTemple | IT 管理器 | `CIllusionTempleEvent_Renewal` | 暂无 | 未覆盖 | 实现幻影寺院 Renewal 副本。 |
| 99 | IllusionTemple | IT 配置加载 | `Load_ITR_EventInfo/Script/NPC_Position` | `conf.Events.IllusionTemple` 粗配置 | 部分覆盖 | 加载开放时间、NPC 位置、规则和奖励。 |
| 100 | IllusionTemple | IT 队伍授权 | `Find_EmptySlot`、`Update_PartyInfo` | `15-party.md` 记录授权缺口 | 未覆盖 | 队伍入场、槽位、同意状态和清理。 |
| 101 | IllusionTemple | IT 入场等级 | `CheckEnterLevel` | 暂无 | 未覆盖 | 按角色等级和寺院层级校验入场。 |
| 102 | IllusionTemple | IT 阵营/队伍 | `GetUserTeam` | 暂无 | 未覆盖 | 管理两个阵营、队伍归属和分数。 |
| 103 | IllusionTemple | IT 圣物交互 | `ActRelicsGetOrRegister` | 暂无 | 未覆盖 | 拾取、登记、丢弃圣物并同步玩家状态。 |
| 104 | IllusionTemple | IT 特殊技能 | `ITR_USeSkill`、`EventSkillProc` | `10-skills.md` 记录事件技能边界 | 未覆盖 | 事件技能使用限制和效果由副本裁决。 |
| 105 | IllusionTemple | IT 死亡重生 | `IllusionTempleUserDie/Regen` | 暂无 | 未覆盖 | 玩家死亡后按副本规则重生并清理圣物。 |
| 106 | IllusionTemple | IT 奖励请求 | `ReqEventReward` | 暂无 | 未覆盖 | 客户端领取或结算幻影寺院奖励。 |
| 107 | ImperialGuardian | IG 管理器 | `CImperialGuardian` | 暂无 | 未覆盖 | 实现帝国要塞副本。 |
| 108 | ImperialGuardian | IG 区域状态 | `MAX_ZONE`、`_stZoneInfo` | 暂无 | 未覆盖 | 管理 4 个 Zone 的状态、成员、怪物和门。 |
| 109 | ImperialGuardian | IG 状态机 | READY/TIMEATTACK/LOOTTIME 等 | 暂无 | 未覆盖 | 实现准备、限时击杀、等待、拾取、传送、失败、通关。 |
| 110 | ImperialGuardian | IG 门票/材料 | `CheckGaionOrderPaper`、`CheckFullSecromicon` | 道具系统未接入 | 未覆盖 | 检查并消耗帝国要塞入场道具。 |
| 111 | ImperialGuardian | IG 传送门 | `CGEnterPortal`、`SetGateBlockState` | 暂无 | 未覆盖 | 区域通关后开启传送门并传送队伍。 |
| 112 | ImperialGuardian | IG 怪物可攻击 | `IsAttackAbleMonster`、`SetAtackAbleState` | 攻击无副本裁决 | 未覆盖 | 阶段性控制怪物可攻击状态。 |
| 113 | ImperialGuardian | IG 门/机关破坏 | `DestroyGate` | 暂无 | 未覆盖 | 攻击并摧毁区域门或机关。 |
| 114 | ImperialGuardian | IG 奖励经验 | `CImperialGuardianRewardExp` | `09-exp.md` 未接入 | 未覆盖 | 根据区域、玩家等级和通关结果发经验。 |
| 115 | ImperialGuardian | IG ItemBag | `EventDungeonItemBag` | `14-drops.md` 有副本 Bag 缺口 | 未覆盖 | 通关、宝箱、雕像奖励接入 Bag。 |
| 116 | DoppelGanger | DG 管理器 | `CDoppelGanger` | 暂无 | 未覆盖 | 实现生魂广场/DoppelGanger 副本；副本系统消费 Lua 配置和回调结果，Lua 运行时、注册和调用归 `26-script.md`。 |
| 117 | DoppelGanger | DG 状态机 | NONE/READY/PLAY/END | 暂无 | 未覆盖 | 实现准备、进行、结束状态和通知。 |
| 118 | DoppelGanger | DG 入场协议 | `EnterDoppelgangerEvent` | `handle` 0xBF0E 占位 | 部分覆盖 | 解析入场道具位置并进入对应地图。 |
| 119 | DoppelGanger | DG 队伍人数 | `MAX_DOPPELGANGER_USER_INFO` | `15-party.md` 待接入 | 未覆盖 | 队伍/成员限制、等级统计和入场授权。 |
| 120 | DoppelGanger | DG 等级统计 | `CalUserLevel`、average/min/max | 暂无 | 未覆盖 | 根据队伍等级生成怪物强度和奖励。 |
| 121 | DoppelGanger | DG 怪群推进 | `CDoppelGangerMonsterHerd` | 怪物 AI 未实现 | 未覆盖 | 按路线移动怪群，处理终点、目标数和进度。 |
| 122 | DoppelGanger | DG 冰工人 | `SetIceWorkerRegen`、`CheckIceWorker` | 暂无 | 未覆盖 | 生成和检查冰工人目标。 |
| 123 | DoppelGanger | DG 幼虫目标 | Golden Larva | 暂无 | 未覆盖 | 管理幼虫生成、击杀和目标计数。 |
| 124 | DoppelGanger | DG 宝箱 | `OpenTreasureBox/OpenLastTreasureBox` | 暂无 | 未覆盖 | 中间和最终宝箱交互与奖励。 |
| 125 | DoppelGanger | DG 奖励经验 | mission result exp | `09-exp.md` 记录 DG 倍率缺口 | 未覆盖 | 通关/失败后发送任务结果和经验奖励。 |
| 126 | DSF/扩展副本 | DevilSquareFinal 协议 | `DevilSquareFinal`、`0xDBxx` | handler 有 0xDB00-0xDB09 占位 | 部分覆盖 | DSF 作为 DS 扩展版本记录，后续可拆子实现。 |
| 127 | DSF/扩展副本 | DSF 日程查询 | `reqDSFSchedule` | `handle` 0xDB00 占位 | 部分覆盖 | 返回 DSF 开放日程。 |
| 128 | DSF/扩展副本 | DSF 队伍入场检查 | `reqDSFCanPartyEnter` | `handle` 0xDB01 占位 | 部分覆盖 | 与组队系统协作检查全队资格。 |
| 129 | DSF/扩展副本 | DSF 接受入场 | `reqAcceptEnterDSF` | `handle` 0xDB02 占位 | 部分覆盖 | 处理队员接受/拒绝进入 DSF。 |
| 130 | DSF/扩展副本 | DSF 奖励领取 | `reqDSFGetReward` | `handle` 0xDB09 占位 | 部分覆盖 | 完成 DSF 奖励发放和记录。 |
| 131 | 跨系统接口 | 与对象系统联动 | object hooks | `04-objects.md` 需接入 | 未覆盖 | 对象系统只调用副本接口，不实现具体副本规则。 |
| 132 | 跨系统接口 | 与地图系统联动 | map range/block | `06-maps.md` 需接入 | 未覆盖 | 地图系统提供地图、Gate、阻挡和回城能力。 |
| 133 | 跨系统接口 | 与道具系统联动 | ticket/item consume | `07-items.md` 需接入 | 未覆盖 | 门票、任务物、奖励物品都经道具系统处理。 |
| 134 | 跨系统接口 | 与经验系统联动 | event exp | `09-exp.md` 需接入 | 未覆盖 | 经验系统查询副本倍率或消费副本结算。 |
| 135 | 跨系统接口 | 与技能系统联动 | IT event skills | `10-skills.md` 需接入 | 未覆盖 | 副本裁决事件技能可用性和目标合法性。 |
| 136 | 跨系统接口 | 与任务系统联动 | `FinishEventProc`、BC/CC/DS quest | `11-quests.md` 已记录缺口 | 未覆盖 | 副本完成、击杀、积分触发任务进度。 |
| 137 | 跨系统接口 | 与合成系统联动 | BC/DS ticket mix | `12-mix.md` 已记录 | 未覆盖 | 门票合成归合成系统，副本只消费入场物。 |
| 138 | 跨系统接口 | 与掉落系统联动 | EventDungeonItemBag、DoppelGangerItemBag | `14-drops.md` 需接入 | 未覆盖 | 副本掉落和奖励 Bag 由掉落系统落地。 |
| 139 | 跨系统接口 | 与组队系统联动 | party auth | `15-party.md` 已记录 | 未覆盖 | 组队副本入场、资格、队伍失败和成员踢出需接入。 |
| 140 | 协议与测试 | 副本 handler 路由测试 | protocol dispatch | `handle/c1c2.go` 有占位 | 部分覆盖 | 确认各副本 opcode 都委托到副本服务。 |
| 141 | 协议与测试 | 入场资格测试 | CheckEnter* | 待实现 | 未覆盖 | 覆盖等级、门票、PK、人数、互斥状态、队伍授权。 |
| 142 | 协议与测试 | 状态机测试 | `Run/Proc*` | 待实现 | 未覆盖 | 覆盖开放、开始、结束、失败、清理和倒计时。 |
| 143 | 协议与测试 | 玩家生命周期测试 | enter/leave/die/disconnect | 待实现 | 未覆盖 | 覆盖进出副本、死亡、断线、重连和踢出。 |
| 144 | 协议与测试 | 怪物死亡测试 | kill count/score | 待实现 | 未覆盖 | 覆盖怪物死亡推进积分、目标、奖励和掉落上下文。 |
| 145 | 协议与测试 | 地图阻挡测试 | BC block、CC hollow、IG gates | 待实现 | 未覆盖 | 覆盖副本动态阻挡、区域开放和非法移动。 |
| 146 | 协议与测试 | 奖励结算测试 | reward exp/zen/item | 待实现 | 未覆盖 | 覆盖成功、失败、超时、背包满和奖励日志。 |
| 147 | 协议与测试 | 跨系统互斥测试 | trade/shop/chaos/warehouse/move | 待实现 | 未覆盖 | 覆盖副本入场和进行中与其他界面状态互斥。 |
| 148 | 协议与测试 | 世界事件边界测试 | CastleSiege/Crywolf/Kanturu 等 | 候选系统 | 未覆盖 | 确认世界事件不误路由到副本系统。 |
