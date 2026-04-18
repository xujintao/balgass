# 23. 世界事件系统

本模块覆盖 CastleSiege、Crywolf、Kanturu、Raklion、ArcaBattle、AcheronGuardian 这类世界地图级大型事件。世界事件系统是状态机和规则编排层，不拥有战盟、对象、地图、道具、掉落、经验、Buff、商店等底层能力；它负责事件阶段、地图状态、NPC/机关、Boss/怪物、占领/排名、税率/惩罚/奖励等规则，再调用对应基础系统完成实际操作。GM 手动开关、跳阶段、设置城主/状态和运营广播入口归 `29-ops.md`。BloodCastle、DevilSquare、ChaosCastle、IllusionTemple、ImperialGuardian、DoppelGanger 归 `21-dungeons.md`；EventChip、Rena、Lotto、LuckyCoin、BonusEvent、节日掉落和普通地图入侵归 `22-events.md`。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 模块边界与总入口 | WorldEventManager 总管理器 | `g_CastleSiege`、`g_Crywolf`、`g_Kanturu`、`g_Raklion`、`g_ArcaBattle`、`g_AcheronGuardianEvent` | 暂无 `game/worldevent` | 未覆盖 | 建立世界事件服务，统一注册和调度各大型地图事件。 |
| 2 | 模块边界与总入口 | 世界事件类型枚举 | CastleSiege/Crywolf/Kanturu/Raklion/Arca/Acheron | 暂无 | 未覆盖 | 定义世界事件类型，避免业务层直接按地图或 opcode 分支。 |
| 3 | 模块边界与总入口 | 世界事件实例上下文 | event state/map/group/guild | 暂无 | 未覆盖 | 表达事件类型、地图、状态、参与战盟/玩家和当前阶段。 |
| 4 | 模块边界与总入口 | 世界事件 Tick 调度 | 各事件 `Run` | 运行时无世界事件 tick | 未覆盖 | 接入主循环或定时器，驱动事件状态机。 |
| 5 | 模块边界与总入口 | GM/调试入口 | `OperateGmCommand`、manual start | 暂无 | 未覆盖 | GM 命令和审计归 `29-ops.md`，世界事件系统执行状态变更。 |
| 6 | 配置与时间表 | events.xml server 节点 | `EventManagement`/各事件 Load | `conf.Events` 已读取相关节点 | 部分覆盖 | 接入 CastleSiege、LorenDeep、Crywolf、Kanturu、Raklion、ArcaBattle、AcheronGuardian。 |
| 7 | 配置与时间表 | 世界事件专用配置 | `LoadData/LoadScript` | 配置结构不完整 | 未覆盖 | 加载地图属性、状态时间、怪物、NPC、奖励、税率和报名规则。 |
| 8 | 配置与时间表 | 事件时间同步 | `CheckSync`、state time info | 暂无 | 未覆盖 | 根据日程自动切换事件阶段。 |
| 9 | 配置与时间表 | 事件状态广播 | notify state/left time | 多数协议占位 | 未覆盖 | 向地图玩家或全服广播阶段、剩余时间、开始/结束。 |
| 10 | 配置与时间表 | 事件数据重载 | `LoadData`、DB load | 暂无 | 未覆盖 | 支持重载事件配置并安全更新运行时状态。 |
| 11 | 地图与进入限制 | 世界事件地图识别 | ValleyOfLoren/Crywolf/Kanturu/Raklion/Acheron | `game/maps/const.go` 已有枚举 | 部分覆盖 | 地图系统提供世界事件地图归属查询。 |
| 12 | 地图与进入限制 | 事件地图属性 | Crywolf/Kanturu map attr | `game/maps/map.go` 有 Crywolf/Kanturu EventAttr | 部分覆盖 | 事件状态切换时替换地图属性、阻挡和安全区。 |
| 13 | 地图与进入限制 | 进入资格检查 | `CheckEnter*` | `MapMove` 有基础能力 | 未覆盖 | 按事件状态、道具、等级、战盟、队伍或地图服限制进入。 |
| 14 | 地图与进入限制 | 事件中移动限制 | move checks | `06-maps.md` 记录事件限制 | 未覆盖 | 世界事件进行中限制传送、换线、进入 Boss 区或战区。 |
| 15 | 地图与进入限制 | 事件离场回城 | kick/rollback | `SpawnPosition` 有特殊地图分支 | 部分覆盖 | 事件结束、失败、断线或非法状态时回到安全位置。 |
| 16 | 对象/NPC/机关接口 | OnPlayerEnterMap | event map enter | 暂无 | 未覆盖 | 玩家进入世界事件地图时绑定事件上下文。 |
| 17 | 对象/NPC/机关接口 | OnPlayerLeaveMap | event leave | 暂无 | 未覆盖 | 玩家离开地图时清理事件临时状态。 |
| 18 | 对象/NPC/机关接口 | OnMonsterDie | boss/obelisk/monster death | 怪物死亡入口存在 | 未覆盖 | 怪物死亡后推进世界事件状态、积分、掉落和奖励。 |
| 19 | 对象/NPC/机关接口 | OnNPCTalk | Castle NPC/Crywolf altar/Kanturu NPC | `Player.Talk` 基础存在 | 未覆盖 | NPC 对话委托世界事件系统裁决。 |
| 20 | 对象/NPC/机关接口 | CanMove | gate/door/block/obelisk zone | 移动无事件裁决 | 未覆盖 | 世界事件动态阻挡、Boss 区、攻城区域限制移动。 |
| 21 | 对象/NPC/机关接口 | CanAttack | siege side/altar/boss egg/obelisk | 攻击入口存在 | 未覆盖 | 世界事件决定阵营、机关、Boss、保护状态是否可攻击。 |
| 22 | 对象/NPC/机关接口 | 事件机关对象 | crown/switch/altar/statue/gate/obelisk/egg | NPC/怪物对象基础存在 | 未覆盖 | 建模世界事件中的可交互或可攻击机关。 |
| 23 | CastleSiege | 攻城总管理器 | `CCastleSiege` | 暂无 | 未覆盖 | 实现罗兰攻城状态机和规则入口。 |
| 24 | CastleSiege | 攻城状态机 | `CASTLESIEGE_STATE_*` | 暂无 | 未覆盖 | 实现报名、注册标识、通知、准备、开始、结束、周期结束。 |
| 25 | CastleSiege | 攻城日程 | `_CS_SCHEDULE_DATA` | `conf.Events.CastleSiege` 粗配置 | 部分覆盖 | 加载攻城周期、阶段偏移和持续时间。 |
| 26 | CastleSiege | 攻城状态查询协议 | `castleSiegeState` | `handle` 0xB200 占位 | 部分覆盖 | 返回当前攻城状态、剩余时间和城主信息。 |
| 27 | CastleSiege | 战盟报名 | `castleSiegeReg` | `handle` 0xB201 占位 | 部分覆盖 | 战盟盟主报名攻城，调用 `16-guild.md` 校验。 |
| 28 | CastleSiege | 放弃报名 | `castleSiegeGiveUp` | `handle` 0xB202 占位 | 部分覆盖 | 取消战盟攻城报名并同步 DB。 |
| 29 | CastleSiege | 战盟报名信息 | `guildRegInfo` | `handle` 0xB203 占位 | 部分覆盖 | 查询战盟报名状态、标识数量和排名。 |
| 30 | CastleSiege | 战盟标识注册 | `guildRegMark` | `handle` 0xB204 占位 | 部分覆盖 | 上交城主标识并更新排名分数。 |
| 31 | CastleSiege | 攻城战盟列表 | `csGuildRegList/csGuildAttackList` | `handle` 0xB4/0xB5 占位 | 部分覆盖 | 查询报名和进攻战盟列表。 |
| 32 | CastleSiege | 城主战盟标识 | `guildMarkOfCastleOwner` | `handle` 0xB902 占位 | 部分覆盖 | 返回当前城主战盟标识、名称和归属。 |
| 33 | CastleSiege | 攻城参与方 | `_CS_TOTAL_GUILD_DATA` | 战盟基础存在 | 未覆盖 | 管理守方、攻方、联盟和参与标记。 |
| 34 | CastleSiege | 攻城 NPC 数据 | `_CS_NPC_DATA` | 对象/NPC 无攻城扩展 | 未覆盖 | 管理城门、守护石像、开关等 NPC 数据。 |
| 35 | CastleSiege | NPC DB 列表 | `npcDBList` | `handle` 0xB3 占位 | 部分覆盖 | 加载攻城 NPC 状态、等级、HP 和位置。 |
| 36 | CastleSiege | NPC 购买 | `npcBuy` | `handle` 0xB205 占位 | 部分覆盖 | 购买或放置攻城 NPC。 |
| 37 | CastleSiege | NPC 修理 | `npcRepair` | `handle` 0xB206 占位 | 部分覆盖 | 修理城门、守护石像等攻城 NPC。 |
| 38 | CastleSiege | NPC 升级 | `npcUpgrade` | `handle` 0xB207 占位 | 部分覆盖 | 升级防御、恢复、最大 HP 等 NPC 属性。 |
| 39 | CastleSiege | 城门操作 | `csGateOperate` | `handle` 0xB212 占位 | 部分覆盖 | 操作城门开关、权限和状态同步。 |
| 40 | CastleSiege | 皇冠交互 | `CastleCrownAct` | 暂无 | 未覆盖 | 处理皇冠登记、占领、打断和胜利判定。 |
| 41 | CastleSiege | 皇冠开关 | `CastleCrownSwitch` | 暂无 | 未覆盖 | 两侧开关占领和皇冠登记条件。 |
| 42 | CastleSiege | 迷你地图数据 | `csMiniMapData/csminiMapDataStop` | `handle` 0xB21B/0xB21C 占位 | 部分覆盖 | 攻城期间同步战盟成员和 NPC 位置。 |
| 43 | CastleSiege | 攻城指挥命令 | `csSendCommand` | `handle` 0xB21D 占位 | 部分覆盖 | 盟主/指挥发送进攻、防守、集结等命令。 |
| 44 | CastleSiege | 狩猎区进入开关 | `csSetEnterHuntZone` | `handle` 0xB21F 占位 | 部分覆盖 | 城主设置是否允许进入狩猎区。 |
| 45 | CastleSiege | 税金信息 | `taxMoneyInfo` | `handle` 0xB208 占位 | 部分覆盖 | 查询城堡税金、商店税、Chaos 税等。 |
| 46 | CastleSiege | 税率修改 | `taxRateChange` | `handle` 0xB209 占位 | 部分覆盖 | 城主修改税率并影响商店/合成。 |
| 47 | CastleSiege | 税金提取 | `moneyDrawOut` | `handle` 0xB210 占位 | 部分覆盖 | 城主提取税金，处理权限和金额上限。 |
| 48 | CastleSiege | 商店税率联动 | Castle Siege tax shop | `08-shops.md` 已记录 | 未覆盖 | NPC 商店买入、出售、修理查询攻城税率。 |
| 49 | CastleSiege | 合成税率联动 | Chaos tax | `12-mix.md` 可接入 | 未覆盖 | ChaosBox 费用按城堡税率修正。 |
| 50 | CastleSiege | 攻城药水 | Siege Potion | `game/item/item.go` 有道具分支 | 部分覆盖 | 攻城药水使用效果、限制和持续时间。 |
| 51 | CastleSiege | 攻城数据保存 | CSP NPC/Guild save | 暂无 | 未覆盖 | 保存 NPC、税率、报名、城主、周期等 DB 数据。 |
| 52 | LorenDeep/CastleDeep | LorenDeep 配置 | `conf.Events.LorenDeep` | 已读取 server 节点 | 部分覆盖 | 作为罗兰峡谷深处事件或攻城附属事件接入。 |
| 53 | LorenDeep/CastleDeep | CastleDeep 管理器 | `CCastleDeepEvent` | 暂无 | 未覆盖 | 实现 CastleDeep 刷怪/突袭状态机。 |
| 54 | LorenDeep/CastleDeep | CastleDeep 状态机 | `CD_STATE_*` | 暂无 | 未覆盖 | 实现关闭、进行中和清理。 |
| 55 | LorenDeep/CastleDeep | CastleDeep 怪物配置 | `CASTLEDEEP_MONSTERINFO` | 暂无 | 未覆盖 | 加载怪物组、类型、数量和区域。 |
| 56 | LorenDeep/CastleDeep | CastleDeep 波次时间 | `CASTLEDEEP_SPAWNTIME` | 暂无 | 未覆盖 | 按突袭类型和时间刷新怪物。 |
| 57 | LorenDeep/CastleDeep | CastleDeep 清怪 | `ClearMonster` | 对象删除能力存在 | 未覆盖 | 事件结束或重载时清理怪物。 |
| 58 | Crywolf | Crywolf 管理器 | `CCrywolf` | 暂无 | 未覆盖 | 实现狼魂要塞世界事件。 |
| 59 | Crywolf | Crywolf 状态机 | `CRYWOLF_STATE_*` | 暂无 | 未覆盖 | 实现通知、准备、开始、结束和周期结束。 |
| 60 | Crywolf | Crywolf 地图属性 | `LoadCrywolfMapAttr/SetCrywolfMapAttr` | `game/maps/map.go` 有 CrywolfEventAttr | 部分覆盖 | 按占领状态切换地图属性。 |
| 61 | Crywolf | Crywolf 占领状态 | `m_iOccupationState` | 暂无 | 未覆盖 | 管理成功/失败后的占领状态和全局影响。 |
| 62 | Crywolf | Crywolf 信息协议 | `reqCrywolfInfo` | `handle` 0xBD00 占位 | 部分覆盖 | 返回占领状态和事件状态。 |
| 63 | Crywolf | Crywolf 剩余时间 | `PMSG_ANS_CRYWOLF_LEFTTIME` | 暂无 | 未覆盖 | 返回事件剩余小时和分钟。 |
| 64 | Crywolf | Crywolf 祭坛 | `CrywolfAltar` | 暂无 | 未覆盖 | 祭坛绑定精灵、状态同步、失败条件和 Buff。 |
| 65 | Crywolf | Crywolf 雕像 | `CrywolfStatue` | 暂无 | 未覆盖 | 雕像 HP、状态、被攻击和防守目标。 |
| 66 | Crywolf | Crywolf NPC 信息 | `CCrywolfObjInfo` | 暂无 | 未覆盖 | 管理普通 NPC、特殊 NPC、普通怪、特殊怪集合。 |
| 67 | Crywolf | Crywolf 怪物 AI 切换 | `ChangeAI`、`SetCrywolfAllCommonMonsterState2` | `25-monster-ai.md` | 未覆盖 | 事件阶段驱动怪物 AI 状态。 |
| 68 | Crywolf | Crywolf Boss 出现 | `TurnUpBoss` | 暂无 | 未覆盖 | Boss 刷新、通知和击杀状态。 |
| 69 | Crywolf | Crywolf 怪物死亡 | `CrywolfMonsterDieProc` | 怪物死亡入口存在 | 未覆盖 | 计算分数、排名、奖励和事件成败。 |
| 70 | Crywolf | Crywolf 英雄榜 | `NotifyCrywolfHeroList` | 暂无 | 未覆盖 | 计算并下发英雄榜和 Top5 奖励。 |
| 71 | Crywolf | Crywolf 个人排名 | `NotifyCrywolfPersonalRank` | 暂无 | 未覆盖 | 下发个人 MVP 积分和排名。 |
| 72 | Crywolf | Crywolf 奖励经验 | `CalcGettingRewardExp/GiveRewardExp` | `09-exp.md` 待接入 | 未覆盖 | 根据 MVP 排名发经验奖励。 |
| 73 | Crywolf | Crywolf 经验惩罚 | `CrywolfSync.GetGettingExpPenaltyRate` | `09-exp.md` 已记录 | 未覆盖 | 事件失败后全局或地图经验惩罚。 |
| 74 | Crywolf | Crywolf 合成加成 | `reqPlusChaosRate`/Crywolf bonus | `12-mix.md` 已记录 | 未覆盖 | ChaosBox 查询 Crywolf 成功率加成或惩罚。 |
| 75 | Crywolf | Crywolf DB 保存 | `CrywolfInfoDBSave/Load` | 暂无 | 未覆盖 | 保存事件状态、占领状态和历史结果。 |
| 76 | Kanturu | Kanturu 管理器 | `CKanturu` | 暂无 | 未覆盖 | 实现坎特鲁世界事件。 |
| 77 | Kanturu | Kanturu 状态机 | `KANTURU_STATE_*` | 暂无 | 未覆盖 | 实现 Standby、Maya、Nightmare、Refinement、End。 |
| 78 | Kanturu | Kanturu 地图属性 | `LoadKanturuMapAttr/SetKanturuMapAttr` | `game/maps/map.go` 有 KanturuEventAttr | 部分覆盖 | 按阶段切换坎特鲁地图属性。 |
| 79 | Kanturu | Kanturu 状态查询协议 | `reqKanturuStateInfo` | `handle` 0xD100 占位 | 部分覆盖 | 返回当前状态、详细状态和剩余时间。 |
| 80 | Kanturu | Kanturu BossMap 入场 | `reqEnterKanturuBossMap` | `handle` 0xD101 占位 | 部分覆盖 | 校验月之石、状态、人数和传送。 |
| 81 | Kanturu | 月之石检查 | `CheckEqipmentMoonStone` | 道具系统未接入 | 未覆盖 | 检查并控制是否需要月之石。 |
| 82 | Kanturu | Battle Standby | `CKanturuBattleStanby` | 暂无 | 未覆盖 | 准备阶段、玩家等待和开始条件。 |
| 83 | Kanturu | Battle of Maya | `CKanturuBattleOfMaya` | 暂无 | 未覆盖 | Maya 阶段怪物、目标、成功/失败。 |
| 84 | Kanturu | Battle of Nightmare | `CKanturuBattleOfNightmare` | 暂无 | 未覆盖 | Nightmare Boss 阶段和结果。 |
| 85 | Kanturu | Refinement Tower | `CKanturuTowerOfRefinement` | 暂无 | 未覆盖 | 提炼之塔开放、入口 NPC 和状态。 |
| 86 | Kanturu | Kanturu 怪物死亡 | `KanturuMonsterDieProc` | 怪物死亡入口存在 | 未覆盖 | 推进 Maya/Nightmare 阶段和掉落。 |
| 87 | Kanturu | Kanturu 掉落率 | MoonStone/JewelOfHarmony drop rate | 配置散落 | 未覆盖 | 掉落系统查询月之石、再生宝石相关掉率。 |
| 88 | Kanturu | Kanturu 计数与日期 | battle counter/date | 暂无 | 未覆盖 | 记录每日战斗次数和状态。 |
| 89 | Raklion | Raklion 管理器 | `CRaklion` | 暂无 | 未覆盖 | 实现冰霜之城 Boss 世界事件。 |
| 90 | Raklion | Raklion 状态机 | `RAKLION_STATE_*` | 暂无 | 未覆盖 | 实现通知、准备、Boss 战、关门、全灭、结束。 |
| 91 | Raklion | Raklion Boss 启用 | `SetRaklionBossEnable` | `conf.Events.Raklion` 已读取 | 部分覆盖 | 根据配置和状态控制 Boss 事件开启。 |
| 92 | Raklion | Raklion 入场检查 | `CheckEnterRaklion` | 地图移动基础存在 | 未覆盖 | 检查 Boss 区入场条件和状态。 |
| 93 | Raklion | Raklion BossMap 玩家检查 | `CheckUserOnRaklionBossMap` | 暂无 | 未覆盖 | 周期清理非法或死亡玩家。 |
| 94 | Raklion | Boss 蛋生成 | `RegenBossEgg` | 暂无 | 未覆盖 | 生成 Selupan 前置 Boss 蛋。 |
| 95 | Raklion | Boss 蛋死亡计数 | `BossEggDieIncrease/Decrease` | 暂无 | 未覆盖 | 达到数量后推进 Boss 出现。 |
| 96 | Raklion | Selupan 战斗 | `CRaklionBattleOfSelupan` | 暂无 | 未覆盖 | Selupan Boss 状态、技能、击杀和奖励。 |
| 97 | Raklion | Raklion 怪物管理 | `RaklionMonsterMng` | 暂无 | 未覆盖 | 事件怪物生成、删除和状态。 |
| 98 | Raklion | 全员死亡处理 | `RAKLION_STATE_ALL_USER_DIE` | 暂无 | 未覆盖 | Boss 区全员死亡后失败和清理。 |
| 99 | ArcaBattle | ArcaBattle 管理器 | `CArcaBattle` | 暂无 | 未覆盖 | 实现阿卡伦战役世界事件。 |
| 100 | ArcaBattle | ArcaBattle 状态机 | `ARCA_STATE_*` | 暂无 | 未覆盖 | 实现同步、报名、准备、组队、进行、结果、关闭。 |
| 101 | ArcaBattle | 盟主报名 | `arcaBattleGuildMasterJoin` | `handle` 0xF830 占位 | 部分覆盖 | 战盟盟主报名 ArcaBattle。 |
| 102 | ArcaBattle | 成员报名 | `arcaBattleGuildMemberJoin` | `handle` 0xF832 占位 | 部分覆盖 | 战盟成员报名和资格校验。 |
| 103 | ArcaBattle | 战场进入 | `arcaBattleEnter` | `handle` 0xF834 占位 | 部分覆盖 | 校验报名状态并传送入战场。 |
| 104 | ArcaBattle | 报名人数查询 | `reqRegisteredMemberCnt` | `handle` 0xF841 占位 | 部分覆盖 | 查询已报名成员数量。 |
| 105 | ArcaBattle | 标识注册 | `arcaBattleMarkReg` | `handle` 0xF843 占位 | 部分覆盖 | 注册战盟标识并更新排名。 |
| 106 | ArcaBattle | 标识排名 | `arcaBattleMarkRank` | `handle` 0xF845 占位 | 部分覆盖 | 查询报名标识排名。 |
| 107 | ArcaBattle | 方尖碑配置 | `_tagOBELISK_INFO` | 暂无 | 未覆盖 | 加载方尖碑位置、属性、HP 和守护怪。 |
| 108 | ArcaBattle | 方尖碑状态 | `ARCA_BATTLE_OBELISK_STATE` | 暂无 | 未覆盖 | 管理未占领、占领中、已占领状态。 |
| 109 | ArcaBattle | Aura 状态 | `_tagAURA_STATE` | 暂无 | 未覆盖 | 管理方尖碑周围光环状态和释放战盟。 |
| 110 | ArcaBattle | 贡献/击杀积分 | `_stABAcquiredPoints` | 暂无 | 未覆盖 | 记录战利品、贡献和击杀积分。 |
| 111 | ArcaBattle | 战利品兑换 | `arcaBattleBootyExchange` | `handle` 0xF836 占位 | 部分覆盖 | 按 Booty 配置兑换奖励。 |
| 112 | ArcaBattle | Arca 掉落 | `g_ArcaBattle.DropItem` | `14-drops.md` 已记录 | 未覆盖 | 怪物死亡时触发 Arca 特殊掉落。 |
| 113 | ArcaBattle | 战盟联动 | guild master/member checks | `16-guild.md` 已记录边界 | 未覆盖 | 使用战盟系统校验盟主、成员、联盟和战盟归属。 |
| 114 | AcheronGuardian | AcheronGuardian 管理器 | `CAcheronGuardianEvent` | 暂无 | 未覆盖 | 实现阿卡伦守护者世界事件。 |
| 115 | AcheronGuardian | Acheron 状态机 | Closed/Ready/Playing/PlayEnd/ChannelClose | 暂无 | 未覆盖 | 实现关闭、准备、进行、结束和频道关闭。 |
| 116 | AcheronGuardian | Acheron 进入 | `acheronEnter`、`CGReqAcheronEventEnter` | `handle` 0xF820/0xF84B 占位 | 部分覆盖 | 普通 Acheron 和事件 Acheron 入场。 |
| 117 | AcheronGuardian | SpiritMap 兑换 | `spritMapExchange` | `handle` 0xF83C 占位 | 部分覆盖 | 兑换 SpiritMap，校验材料和地图条件。 |
| 118 | AcheronGuardian | Obelisk 生成 | `GenObelisk/DelObelisk` | 暂无 | 未覆盖 | 生成和删除 Acheron 方尖碑。 |
| 119 | AcheronGuardian | Obelisk 属性随机 | `SetRandomObeliskAttr` | 暂无 | 未覆盖 | 随机方尖碑属性并影响怪物/奖励。 |
| 120 | AcheronGuardian | 怪物组配置 | `_stAEMonPosition/_stAEMonGroupInfo` | 暂无 | 未覆盖 | 加载怪物组、区域、数量、刷新间隔。 |
| 121 | AcheronGuardian | 怪物刷新 | `RegenMonsterRun/RegenMonster` | 暂无 | 未覆盖 | 按组和时间刷新 Acheron 事件怪物。 |
| 122 | AcheronGuardian | 怪物删除 | `DeleteMonster/DeleteAcheronEventAllMonster` | 对象删除能力存在 | 未覆盖 | 阶段结束或重载时清理事件怪物。 |
| 123 | AcheronGuardian | Obelisk 破坏 | `DestroyObelisk` | 暂无 | 未覆盖 | 根据最大伤害玩家、地图和坐标处理破坏与奖励。 |
| 124 | AcheronGuardian | 地图服组播 | `SendMapServerGroupMsg` | 暂无 | 未覆盖 | 多地图服/跨服状态同步通道归 `28-external-comm.md`。 |
| 125 | 跨系统接口 | 与战盟系统联动 | CastleSiege/Arca guild checks | `16-guild.md` 已记录 | 未覆盖 | 世界事件调用战盟系统，不实现战盟基础数据。 |
| 126 | 跨系统接口 | 与对象系统联动 | NPC/Monster/User hooks | `04-objects.md` 需接入 | 未覆盖 | 对象系统通知事件，世界事件裁决具体规则。 |
| 127 | 跨系统接口 | 与地图系统联动 | map attr/gate/move | `06-maps.md` 需接入 | 未覆盖 | 地图系统提供地图属性、Gate、阻挡和回城。 |
| 128 | 跨系统接口 | 与商店系统联动 | Castle tax shop | `08-shops.md` 已记录 | 未覆盖 | 商店买卖和修理查询城堡税率。 |
| 129 | 跨系统接口 | 与经验系统联动 | Crywolf penalty | `09-exp.md` 已记录 | 未覆盖 | 经验系统查询 Crywolf 惩罚和世界事件加成/惩罚。 |
| 130 | 跨系统接口 | 与技能系统联动 | Castle/event skills | `10-skills.md` 已记录 | 未覆盖 | 攻城/世界事件技能由世界事件裁决可用性。 |
| 131 | 跨系统接口 | 与合成系统联动 | Crywolf plus chaos、Castle tax | `12-mix.md` 已记录 | 未覆盖 | 合成系统查询世界事件的成功率和税率修正。 |
| 132 | 跨系统接口 | 与 Buff 系统联动 | Crywolf altar buff | `13-buffs.md` 已记录 | 未覆盖 | 世界事件触发 Buff，Buff 生命周期归 Buff 系统。 |
| 133 | 跨系统接口 | 与掉落系统联动 | Arca/Kanturu/Raklion drops | `14-drops.md` 需接入 | 未覆盖 | 世界事件选择特殊掉落上下文，实际生成归掉落系统。 |
| 134 | 跨系统接口 | 与副本系统边界 | BC/DS/CC/IT/IG/DG | `21-dungeons.md` 已转正 | 未覆盖 | 世界事件不得实现副本状态机。 |
| 135 | 跨系统接口 | 与普通活动边界 | EventChip/Bonus/节日/入侵 | `22-events.md` 已转正 | 未覆盖 | 世界事件不得实现普通活动编排。 |
| 136 | 协议与测试 | CastleSiege 协议测试 | 0xB2/0xB3/0xB4/0xB5/0xB902 | 待实现 | 未覆盖 | 覆盖状态、报名、标识、NPC、税率、城门、皇冠、迷你地图。 |
| 137 | 协议与测试 | Crywolf 协议测试 | 0xBD00、state packets | 待实现 | 未覆盖 | 覆盖状态查询、祭坛/雕像、Boss、排名、经验惩罚。 |
| 138 | 协议与测试 | Kanturu 协议测试 | 0xD100/0xD101 | 待实现 | 未覆盖 | 覆盖状态查询、BossMap 入场、月之石、阶段切换。 |
| 139 | 协议与测试 | Raklion 测试 | Selupan/Boss egg flow | 待实现 | 未覆盖 | 覆盖 Boss 蛋、Boss 区、Selupan、全员死亡和奖励。 |
| 140 | 协议与测试 | ArcaBattle 协议测试 | 0xF830/32/34/36/41/43/45 | 待实现 | 未覆盖 | 覆盖报名、进入、标识、方尖碑、战利品和排名。 |
| 141 | 协议与测试 | AcheronGuardian 协议测试 | 0xF820/0xF83C/0xF84B | 待实现 | 未覆盖 | 覆盖入场、SpiritMap、方尖碑、怪物刷新和组播。 |
| 142 | 协议与测试 | 地图属性测试 | Crywolf/Kanturu/Raklion/Arca | 待实现 | 未覆盖 | 覆盖事件阶段切换后的地图属性、阻挡和进入限制。 |
| 143 | 协议与测试 | 跨系统边界测试 | guild/shop/exp/mix/buff/drop | 待实现 | 未覆盖 | 覆盖世界事件对基础系统的查询和触发，不反向耦合规则。 |
