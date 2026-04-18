# 25. 怪物 AI 系统

本模块覆盖怪物行为决策层：怪物找敌、仇恨、游走、追击、攻击、AI 状态机、AI 规则、AI 行为元素、怪物技能、怪群、特殊怪和怪物再生策略。对象系统负责创建和持有怪物对象、视野、移动执行、攻击入口、死亡状态和调度入口；地图系统负责地形、阻挡、寻路和坐标能力；技能、公式、Buff、掉落、副本、普通活动、世界事件、宠物与召唤系统通过事件或服务接口与怪物 AI 协作，不在各自模块重复实现怪物行为决策。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 模块边界与总入口 | MonsterAIManager 总管理器 | `TMonsterAI` | 暂无独立 `game/monsterai` | 未覆盖 | 建立怪物 AI 服务，承接怪物行为 tick、规则、仇恨、技能、怪群和特殊行为。 |
| 2 | 模块边界与总入口 | AI 初始化入口 | `LoadData` 系列、`MonsterAIProc` | `SpawnMonster` 只初始化 MonsterTable/DropManager | 未覆盖 | 启服时加载 AIUnit、Automata、AIElement、AIRule、MovePath、MonsterSkill 配置。 |
| 3 | 模块边界与总入口 | 对象层委托入口 | `TMonsterAI::RunAI`、`gObjMonsterProcess` | `Monster.ProcessAction` 直接执行简化行为 | 部分覆盖 | 将 `ProcessAction` 改为对象层入口，复杂决策委托给怪物 AI 服务。 |
| 4 | 模块边界与总入口 | AI 100ms 调度 | `MonsterMoveProc`、`MonsterAIProc` | `ObjectManager.Process100ms` | 部分覆盖 | 保持对象主循环，增加按怪物 AI delay 控制的调度。 |
| 5 | 模块边界与总入口 | AI 1000ms 调度 | `MonsterSkillProc`、周期检查 | `Monster.Process1000ms` 空实现 | 未覆盖 | 处理仇恨衰减、延迟技能、怪群状态、特殊怪周期逻辑。 |
| 6 | 模块边界与总入口 | AI 延迟时间 | `TMonsterAIUnit::m_iDelayTime` | `nextActionInterval`、`delayActionInterval` | 部分覆盖 | 对齐 AIUnit 和行为元素 delay，而不是只用移动/攻击间隔。 |
| 7 | 模块边界与总入口 | AI 当前单元字段 | `s_iMonsterCurrentAIUnitTable` | 暂无 | 未覆盖 | 为怪物 class 或实例记录当前 AIUnit，支持事件切换。 |
| 8 | 模块边界与总入口 | AI 数据热清理 | `DelAllAIUnit/DelAllAutomata/DelAllAIElement` | 暂无 | 未覆盖 | 支持重载前清空 AI 配置，避免旧指针或旧索引残留。 |
| 9 | 模块边界与总入口 | AI 配置索引 | `FindAIUnit`、`FindAutomata`、`FindAIElement` | 暂无 | 未覆盖 | 建立按编号查找的只读配置索引。 |
| 10 | 模块边界与总入口 | AI 错误降级 | GameServer 默认基础行为 | 简化行为已存在 | 部分覆盖 | 配置缺失时回退到基础游走/追击/攻击，避免怪物卡死。 |
| 11 | 基础怪物行为 | Monster 属性加载 | `MonsterAttr` | `MonsterTable.init` | 部分覆盖 | 保留怪物基础属性读取，补齐 AI 编号、脚本 HP、抗性和 MonsterSkill 语义。 |
| 12 | 基础怪物行为 | Monster Spawn 数据 | `MonsterSetBase` | `SpawnMonster` 读取 `IGC_MonsterSpawn.xml` | 部分覆盖 | 对齐 NPC/单点/范围/元素怪刷出语义。 |
| 13 | 基础怪物行为 | 出生位置随机 | `CMonsterSetBase::GetPosition/GetBoxPosition` | `SpawnPosition`、`RandPosition` | 部分覆盖 | 完善出生范围、阻挡、同点再生、特殊地图例外。 |
| 14 | 基础怪物行为 | 怪物待机 | `gObjMonsterBaseAct` | `baseAction` emotion 0/2/3 | 部分覆盖 | 明确 idle/rest 状态、下一次动作时间和目标清理。 |
| 15 | 基础怪物行为 | 视野找敌 | `gObjMonsterSearchEnemy` | `searchEnemy` | 部分覆盖 | 当前只取最近目标，后续接入仇恨、对象类型和特殊规则。 |
| 16 | 基础怪物行为 | 守卫找敌 | `gObjGuardSearchEnemy` | `searchEnemy` 对 247/249 有 PK 过滤 | 部分覆盖 | 抽出 Guard AI，按 PK 等级、安全区和地图属性搜敌。 |
| 17 | 基础怪物行为 | 随机游走 | `gObjMonsterMoveAction` | `roamMove` | 部分覆盖 | 保留 moveRange 随机游走，补齐路径、站位占用和失败重试策略。 |
| 18 | 基础怪物行为 | 追击移动 | `gObjMonsterGetTargetPos`、`gObjGetTargetPos` | `chaseMove` | 部分覆盖 | 从单步追击升级为可受 AIElement 和寻路策略控制的追击。 |
| 19 | 基础怪物行为 | 攻击触发 | `gObjMonsterAttack` | `Monster.attack` | 部分覆盖 | 普攻和技能选择应由 AI/MonsterSkill 决定，攻击执行仍调用对象/技能系统。 |
| 20 | 基础怪物行为 | 目标过远处理 | `gObjMonsterBaseAct` | `baseAction` 中 `ViewRange<<1` 判断 | 部分覆盖 | 目标过远后丢失、返回出生范围或切换状态。 |
| 21 | 基础怪物行为 | 返回出生范围 | `gObjMonsterMoveRegen` | `overDis` 和重置 StartX/StartY | 部分覆盖 | 怪物超出活动半径后回归出生点或重新选点。 |
| 22 | 基础怪物行为 | 不行动怪物 | monster attribute 0 | `Attribute == 0` 直接 return | 已覆盖 | 保持 NPC/静态怪不执行主动 AI。 |
| 23 | 基础怪物行为 | 移动可行性检查 | `gObjMonsterMoveCheck` | `CheckMapNoWall`、`GetMapAttr` | 部分覆盖 | 统一检查地图、阻挡、安全区、站位和特殊怪限制。 |
| 24 | 基础怪物行为 | 路径移动消息 | `PathFindMoveMsgSend` | `Monster.move` 构造 `MsgMove` | 部分覆盖 | 对齐路径长度、方向序列、失败处理和视野广播。 |
| 25 | 基础怪物行为 | 死亡行为清理 | `MonsterDieAction`、`gObjMonsterDieGiveItem` | `Object.processRegen`、掉落入口分散 | 部分覆盖 | 死亡时清理 AI 状态、仇恨、怪群成员状态和延迟技能。 |
| 26 | 目标与合法性 | 目标对象存在检查 | `ObjectMaxRange` | `ObjectManager.GetObject` | 部分覆盖 | 所有 AI 行为统一检查目标存在。 |
| 27 | 目标与合法性 | 目标存活检查 | `lpTargetObj->Live` | `!tobj.Live` | 部分覆盖 | 目标死亡时切换状态并清理仇恨。 |
| 28 | 目标与合法性 | 目标地图检查 | map number check | `tobj.MapNumber != m.MapNumber` | 部分覆盖 | 目标跨图或离线时重置目标。 |
| 29 | 目标与合法性 | 目标安全区检查 | map attr safe zone | `attr&1 == 0` | 部分覆盖 | 安全区、不可攻击区域、事件保护区统一走地图接口。 |
| 30 | 目标与合法性 | 攻击距离检查 | attack range/view range | `CalcDistance`、`AttackRange` | 部分覆盖 | 普攻、远程、技能、范围技能应分别校验距离。 |
| 31 | 目标与合法性 | 墙体检查 | `CheckWall`/path check | `CheckMapNoWall` | 部分覆盖 | 技能攻击和普通攻击都需要墙体/阻挡校验。 |
| 32 | 目标与合法性 | 玩家类型过滤 | character/object type | `Viewports` 遍历 | 部分覆盖 | 允许按玩家、召唤物、怪物、NPC 类型过滤目标。 |
| 33 | 目标与合法性 | PK 条件过滤 | guard PK logic | Guard class 特判 | 部分覆盖 | 独立 GuardTargetPolicy，避免 class 分支散落。 |
| 34 | 目标与合法性 | 事件目标过滤 | event state target check | 暂无 | 未覆盖 | 副本/活动/世界事件可限制怪物攻击对象。 |
| 35 | 目标与合法性 | 召唤物目标过滤 | summoned monster target | 暂无 | 未覆盖 | 召唤物只能攻击主人可攻击目标或指定阵营目标。 |
| 36 | 目标与合法性 | 目标切换条件 | AI state transition | 简单最近目标 | 未覆盖 | 按仇恨、距离、状态机和特殊规则切换目标。 |
| 37 | 目标与合法性 | 目标丢失冷却 | actionState emotion 3 | emotionCount 简化 | 部分覆盖 | 目标丢失后延迟重新搜敌，避免每 tick 抖动。 |
| 38 | 仇恨系统 | Agro 容器 | `TMonsterAIAgro` | `Object` 注释 `argo` | 未覆盖 | 给怪物实例增加仇恨表。 |
| 39 | 仇恨系统 | Agro entry | `TMonsterAIAgroInfo` | 暂无 | 未覆盖 | 记录目标 index 和仇恨值。 |
| 40 | 仇恨系统 | ResetAll | `ResetAll` | 暂无 | 未覆盖 | 怪物出生、死亡、回收时清空仇恨。 |
| 41 | 仇恨系统 | SetAgro | `SetAgro` | 暂无 | 未覆盖 | 直接设置目标仇恨，用于强制目标或事件脚本。 |
| 42 | 仇恨系统 | DelAgro | `DelAgro` | 暂无 | 未覆盖 | 目标死亡、离线、跨图时删除仇恨。 |
| 43 | 仇恨系统 | GetAgro | `GetAgro` | 暂无 | 未覆盖 | 查询指定目标仇恨。 |
| 44 | 仇恨系统 | IncAgro | `IncAgro` | 暂无 | 未覆盖 | 受击、技能、治疗、嘲讽等增加仇恨。 |
| 45 | 仇恨系统 | DecAgro | `DecAgro` | 暂无 | 未覆盖 | 距离过远、时间衰减或技能效果减少仇恨。 |
| 46 | 仇恨系统 | DecAllAgro | `DecAllAgro` | 暂无 | 未覆盖 | 周期性衰减所有仇恨。 |
| 47 | 仇恨系统 | MaxAgro target | `GetMaxAgroUserIndex` | 暂无 | 未覆盖 | 选择最高仇恨且合法的目标。 |
| 48 | 仇恨系统 | 受击加仇恨 | `BeenAttacked`、hit damage | `attack.go` 有伤害流程 | 未覆盖 | 怪物被攻击后将攻击者加入仇恨表。 |
| 49 | 仇恨系统 | 伤害归属接口 | `gObjMonsterSetHitDamage` | 掉落/经验归属未完整 | 未覆盖 | 仇恨和伤害表分别服务 AI 目标选择与奖励归属。 |
| 50 | 仇恨系统 | 最高伤害玩家 | `gObjMonsterTopHitDamageUser` | 暂无 | 未覆盖 | 给掉落、经验、任务、活动排名提供最高伤害者。 |
| 51 | 仇恨系统 | 最后一击玩家 | `gObjMonsterLastHitDamageUser` | 暂无 | 未覆盖 | 给任务、成就、活动击杀判断提供最后一击者。 |
| 52 | 仇恨系统 | 伤害表清理 | `gObjMonsterHitDamageInit/UserDel` | 暂无 | 未覆盖 | 怪物出生、死亡、玩家离线时清理伤害记录。 |
| 53 | AI 状态机 | Automata 配置 | `TMonsterAIAutomata` | 暂无 | 未覆盖 | 按 automata number 加载状态机。 |
| 54 | AI 状态机 | AIState 数据 | `TMonsterAIState` | `actionState` 简化 | 部分覆盖 | 用配置化状态替代固定 emotion 分支。 |
| 55 | AI 状态机 | 状态优先级 | `m_iPriority` | 暂无 | 未覆盖 | 同状态多转移按优先级决策。 |
| 56 | AI 状态机 | 当前状态 | `m_iCurrentState` | `actionState.emotion` | 部分覆盖 | 保留实例当前 AI 状态。 |
| 57 | AI 状态机 | 下一状态 | `m_iNextState` | 暂无 | 未覆盖 | 状态转移后写入下一状态。 |
| 58 | AI 状态机 | 无敌人转移 | `MAI_STATE_TRANS_NO_ENEMY` | emotion 0/3 | 部分覆盖 | 无目标时转入游走、待机或返回。 |
| 59 | AI 状态机 | 有敌人转移 | `MAI_STATE_TRANS_IN_ENEMY` | `TargetNumber >= 0` | 部分覆盖 | 搜到合法目标后转入追击或攻击。 |
| 60 | AI 状态机 | 敌人离开转移 | `MAI_STATE_TRANS_OUT_ENEMY` | `ViewRange<<1` | 部分覆盖 | 目标超出范围后转入搜索或返回。 |
| 61 | AI 状态机 | 敌人死亡转移 | `MAI_STATE_TRANS_DIE_ENEMY` | `!tobj.Live` | 部分覆盖 | 目标死亡时转入下一个目标或待机。 |
| 62 | AI 状态机 | HP 下降转移 | `MAI_STATE_TRANS_DEC_HP` | 暂无 | 未覆盖 | 受伤达到数值阈值触发逃跑、治疗、特殊技能。 |
| 63 | AI 状态机 | HP 百分比转移 | `MAI_STATE_TRANS_DEC_HP_PER` | 暂无 | 未覆盖 | Boss 阶段切换和濒死技能依赖 HP 百分比。 |
| 64 | AI 状态机 | 立即转移 | `MAI_STATE_TRANS_IMMEDIATELY` | 暂无 | 未覆盖 | 无条件转移用于脚本化状态链。 |
| 65 | AI 状态机 | 仇恨上升转移 | `MAI_STATE_TRANS_AGRO_UP` | 暂无 | 未覆盖 | 仇恨变化触发集火或特殊行为。 |
| 66 | AI 状态机 | 仇恨下降转移 | `MAI_STATE_TRANS_AGRO_DOWN` | 暂无 | 未覆盖 | 仇恨降低触发回归或换目标。 |
| 67 | AI 状态机 | 群体召唤转移 | `MAI_STATE_TRANS_GROUP_SOMMON` | 暂无 | 未覆盖 | 怪群/召唤类 Boss 转入召唤状态。 |
| 68 | AI 状态机 | 群体治疗转移 | `MAI_STATE_TRANS_GROUP_HEAL` | 暂无 | 未覆盖 | 支援怪触发治疗队友。 |
| 69 | AI 状态机 | 转移概率 | `m_iTransitionRate` | 暂无 | 未覆盖 | 状态切换支持概率。 |
| 70 | AI 状态机 | 转移延迟 | `m_iDelayTime` | `nextActionInterval` | 部分覆盖 | 状态转移后应用延迟。 |
| 71 | AI 规则与 AIUnit | AIRule 加载 | `TMonsterAIRule::LoadData` | 暂无 | 未覆盖 | 加载怪物 class 到 AIUnit 的映射规则。 |
| 72 | AI 规则与 AIUnit | 当前 AIUnit 获取 | `GetCurrentAIUnit` | 暂无 | 未覆盖 | 按怪物 class 或事件状态取当前 AIUnit。 |
| 73 | AI 规则与 AIUnit | AIRule 周期处理 | `MonsterAIRuleProc` | 暂无 | 未覆盖 | 支持周期性更新 AIUnit。 |
| 74 | AI 规则与 AIUnit | AIUnit 数据 | `TMonsterAIUnit` | 暂无 | 未覆盖 | AIUnit 保存名称、编号、delay、automata 和各状态行为槽。 |
| 75 | AI 规则与 AIUnit | AIUnit 执行 | `RunAIUnit` | `baseAction` | 未覆盖 | 按当前状态选择 AIElement 并执行。 |
| 76 | AI 规则与 AIUnit | Normal 行为槽 | `m_lpAIClassNormal` | emotion 0 简化 | 部分覆盖 | 待机、普通搜索、普通巡逻。 |
| 77 | AI 规则与 AIUnit | Move 行为槽 | `m_lpAIClassMove` | `roamMove/chaseMove` | 部分覆盖 | 移动行为由配置元素控制。 |
| 78 | AI 规则与 AIUnit | Attack 行为槽 | `m_lpAIClassAttack` | `attack` 随机技能 | 部分覆盖 | 攻击行为由怪物技能和 AIElement 控制。 |
| 79 | AI 规则与 AIUnit | Heal 行为槽 | `m_lpAIClassHeal` | 暂无 | 未覆盖 | 自疗或群疗怪物行为。 |
| 80 | AI 规则与 AIUnit | Avoid 行为槽 | `m_lpAIClassAvoid` | 暂无 | 未覆盖 | 逃跑、保持距离、规避行为。 |
| 81 | AI 规则与 AIUnit | Help 行为槽 | `m_lpAIClassHelp` | 暂无 | 未覆盖 | 援护、Buff、转移目标。 |
| 82 | AI 规则与 AIUnit | Special 行为槽 | `m_lpAIClassSpecial` | class 特判 LearnSkill | 部分覆盖 | Boss 特殊技能、免疫、召唤、传送。 |
| 83 | AI 规则与 AIUnit | Event 行为槽 | `m_lpAIClassEvent` | 副本/活动未接入 | 未覆盖 | 活动阶段驱动的特殊行为。 |
| 84 | AI 行为元素 | AIElement 配置 | `TMonsterAIElement` | 暂无 | 未覆盖 | 加载行为元素编号、类型、成功率、delay、目标类型和坐标。 |
| 85 | AI 行为元素 | Common normal | `MAE_TYPE_COMMON_NORMAL` | `baseAction` | 部分覆盖 | 普通待机/普通状态处理。 |
| 86 | AI 行为元素 | Move normal | `MAE_TYPE_MOVE_NORMAL` | `roamMove` | 部分覆盖 | 随机或指定路径移动。 |
| 87 | AI 行为元素 | Move target | `MAE_TYPE_MOVE_TARGET` | `chaseMove` | 部分覆盖 | 朝目标移动。 |
| 88 | AI 行为元素 | Group move | `MAE_TYPE_GROUP_MOVE` | 暂无 | 未覆盖 | 怪群整体移动。 |
| 89 | AI 行为元素 | Group move target | `MAE_TYPE_GROUP_MOVE_TARGET` | 暂无 | 未覆盖 | 怪群朝目标推进。 |
| 90 | AI 行为元素 | Attack normal | `MAE_TYPE_ATTACK_NORMAL` | `Monster.attack` | 部分覆盖 | 普通攻击元素。 |
| 91 | AI 行为元素 | Attack area | `MAE_TYPE_ATTACK_AREA` | 暂无 | 未覆盖 | 范围攻击元素。 |
| 92 | AI 行为元素 | Attack penetration | `MAE_TYPE_ATTACK_PENETRATION` | 暂无 | 未覆盖 | 直线/穿透攻击元素。 |
| 93 | AI 行为元素 | Heal self | `MAE_TYPE_HEAL_SELF` | 暂无 | 未覆盖 | 怪物自疗。 |
| 94 | AI 行为元素 | Heal group | `MAE_TYPE_HEAL_GROUP` | 暂无 | 未覆盖 | 怪群或附近友方治疗。 |
| 95 | AI 行为元素 | Avoid normal | `MAE_TYPE_AVOID_NORMAL` | 暂无 | 未覆盖 | 远离目标、逃跑、拉开距离。 |
| 96 | AI 行为元素 | Help HP | `MAE_TYPE_HELP_HP` | 暂无 | 未覆盖 | 按友方 HP 触发支援。 |
| 97 | AI 行为元素 | Help Buff | `MAE_TYPE_HELP_BUFF` | 暂无 | 未覆盖 | 给友方施加 Buff。 |
| 98 | AI 行为元素 | Help target | `MAE_TYPE_HELP_TARGET` | 暂无 | 未覆盖 | 协助友方攻击目标。 |
| 99 | AI 行为元素 | Special summon | `MAE_TYPE_SPECIAL_SOMMON` | 昆顿 LearnSkill 200 | 部分覆盖 | 召唤小怪或事件怪。 |
| 100 | AI 行为元素 | Special immune | `MAE_TYPE_SPECIAL_IMMUNE` | 昆顿 LearnSkill 201/202 | 部分覆盖 | 物理/魔法/全技能免疫。 |
| 101 | AI 行为元素 | Nightmare summon | `MAE_TYPE_SPECIAL_NIGHTMARE_SUMMON` | 暂无 | 未覆盖 | Nightmare 特殊召唤。 |
| 102 | AI 行为元素 | Special warp | `MAE_TYPE_SPECIAL_WARP` | 暂无 | 未覆盖 | 怪物传送或拉人。 |
| 103 | AI 行为元素 | Skill attack | `MAE_TYPE_SPECIAL_SKILLATTACK` | `UseSkill` 随机释放 | 部分覆盖 | 按配置释放指定怪物技能。 |
| 104 | AI 行为元素 | Change AI | `MAE_TYPE_SPECIAL_CHANGEAI` | 暂无 | 未覆盖 | 切换 AIUnit 或状态机。 |
| 105 | AI 行为元素 | Event element | `MAE_TYPE_EVENT` | 暂无 | 未覆盖 | 事件专用行为元素。 |
| 106 | AI 路径 | MovePath 管理 | `TMonsterAIMovePath` | `maps.FindMapPath` 只做寻路 | 未覆盖 | 区分预设 AI 路径和地图寻路算法。 |
| 107 | AI 路径 | 路径配置加载 | `LoadData` | 暂无 | 未覆盖 | 按地图和 section 加载路径点。 |
| 108 | AI 路径 | 路径点数据 | `TMonsterAIMovePathInfo` | `maps.Pot` | 部分覆盖 | 保存路径类型、地图、X、Y。 |
| 109 | AI 路径 | 路径点计数 | `m_iMovePathSpotCount` | 暂无 | 未覆盖 | 限制和统计路径点数量。 |
| 110 | AI 路径 | 巡逻路径选择 | move path spot | 暂无 | 未覆盖 | 怪物按路径点巡逻或推进。 |
| 111 | AI 路径 | 事件路径选择 | event move path | 暂无 | 未覆盖 | DoppelGanger、入侵、世界事件使用路径推进。 |
| 112 | 怪物技能系统 | MonsterSkillManager | `TMonsterSkillManager` | `Monster.attack` 随机 `Skills` | 部分覆盖 | 建立怪物技能管理器，替代随机从已学技能里选。 |
| 113 | 怪物技能系统 | 技能配置加载 | `TMonsterSkillManager::LoadData` | `MonsterConfig.MonsterSkill` 未使用 | 未覆盖 | 加载怪物 class 到技能单元映射。 |
| 114 | 怪物技能系统 | CheckMonsterSkill | `CheckMonsterSkill` | 暂无 | 未覆盖 | 判断怪物是否拥有配置化怪物技能。 |
| 115 | 怪物技能系统 | FindMonsterSkillUnit | `FindMonsterSkillUnit` | 暂无 | 未覆盖 | 按类型和编号查找技能单元。 |
| 116 | 怪物技能系统 | UseMonsterSkill | `UseMonsterSkill` | `UseSkill` | 部分覆盖 | 执行配置化怪物技能，并调用通用技能/战斗系统。 |
| 117 | 怪物技能系统 | 延迟技能队列 | `s_MonsterSkillDelayInfoArray` | 对象 delay msg 与 AI 无关 | 未覆盖 | 支持怪物技能延迟释放。 |
| 118 | 怪物技能系统 | 技能 delay info | `_ST_MONSTER_SKILL_DELAYTIME_INFO` | 暂无 | 未覆盖 | 记录怪物、目标、释放时间和技能单元。 |
| 119 | 怪物技能系统 | MonsterSkillProc | `MonsterSkillProc` | `Process1000ms` 空 | 未覆盖 | 周期扫描并释放到期怪物技能。 |
| 120 | 怪物技能系统 | MagicInfo 映射 | `FindMagicInf` | `skill.Skill` | 部分覆盖 | 将怪物技能单元映射到 Go 技能定义。 |
| 121 | 怪物技能系统 | SkillUnit 数据 | `TMonsterSkillUnit` | 暂无 | 未覆盖 | 保存目标类型、范围类型、范围值、delay 和元素列表。 |
| 122 | 怪物技能系统 | SkillUnit RunSkill | `RunSkill` | 暂无 | 未覆盖 | 按目标范围执行多个技能元素。 |
| 123 | 怪物技能系统 | SkillElement 数据 | `TMonsterSkillElement` | `Object.monsterSkillElementInfo` 注释 | 未覆盖 | 保存元素类型、成功率、持续时间、增减类型和值。 |
| 124 | 怪物技能系统 | Stun 元素 | `ApplyElementStun` | Buff/状态未接 | 未覆盖 | 怪物技能造成眩晕。 |
| 125 | 怪物技能系统 | Move 元素 | `ApplyElementMove` | 暂无 | 未覆盖 | 怪物技能移动目标或自身。 |
| 126 | 怪物技能系统 | HP/MP/AG 元素 | `ApplyElementHP/MP/AG` | 暂无 | 未覆盖 | 改变目标资源。 |
| 127 | 怪物技能系统 | Defense/Attack 元素 | `ApplyElementDefense/Attack` | 暂无 | 未覆盖 | 临时改变攻防，通常接 Buff 系统。 |
| 128 | 怪物技能系统 | Summon 元素 | `ApplyElementSummon` | 暂无 | 未覆盖 | 怪物技能召唤怪物。 |
| 129 | 怪物技能系统 | Push 元素 | `ApplyElementPush` | 击退能力不完整 | 未覆盖 | 怪物技能击退目标。 |
| 130 | 怪物技能系统 | Stat 元素 | `ApplyElementStat*` | 暂无 | 未覆盖 | 临时改变力量、敏捷、体力、智力。 |
| 131 | 怪物技能系统 | Remove/Resist/Immune 元素 | `ApplyElementRemoveSkill/ResistSkill/ImmuneSkill` | 暂无 | 未覆盖 | 移除技能效果、抗性、免疫处理。 |
| 132 | 怪物技能系统 | Teleport 元素 | `ApplyElementTeleportSkill` | 暂无 | 未覆盖 | 传送目标或怪物。 |
| 133 | 怪物技能系统 | Double HP 元素 | `ApplyElementDoubleHP` | 暂无 | 未覆盖 | Boss 临时翻倍 HP。 |
| 134 | 怪物技能系统 | Poison 元素 | `ApplyElementPoison` | Buff/毒效果未完整 | 未覆盖 | 怪物技能施加毒。 |
| 135 | 怪物技能系统 | Percent damage 元素 | `ApplyElementPercentDamageNormalAttack` | 暂无 | 未覆盖 | 百分比伤害攻击。 |
| 136 | 怪群系统 | MonsterHerd 基类 | `MonsterHerd` | `Object` 注释 `monsterHerd` | 未覆盖 | 建立怪群对象，管理地图、半径、成员和活动状态。 |
| 137 | 怪群系统 | SetTotalInfo | `SetTotalInfo` | 暂无 | 未覆盖 | 设置怪群地图、半径、起始坐标。 |
| 138 | 怪群系统 | AddMonster | `AddMonster` | 暂无 | 未覆盖 | 向怪群添加怪物类型、是否重生、是否先攻。 |
| 139 | 怪群系统 | Start/Stop | `Start/Stop` | 暂无 | 未覆盖 | 启动或停止怪群。 |
| 140 | 怪群系统 | CheckInRadius | `CheckInRadius` | `overDis` 仅单怪 | 部分覆盖 | 判断成员是否仍在怪群半径内。 |
| 141 | 怪群系统 | CurrentLocation | `GetCurrentLocation` | 暂无 | 未覆盖 | 获取怪群当前位置。 |
| 142 | 怪群系统 | RandomLocation | `GetRandomLocation/CheckLocation` | `RandPosition` | 部分覆盖 | 在怪群半径内选合法位置。 |
| 143 | 怪群系统 | MoveHerd | `MoveHerd` | 暂无 | 未覆盖 | 移动整个怪群到目标点。 |
| 144 | 怪群系统 | Herd item drop | `MonsterHerdItemDrop` | 掉落系统未接怪群上下文 | 未覆盖 | 只传递怪群掉落上下文，最终掉落归 `14-drops.md`。 |
| 145 | 怪群系统 | BeenAttacked | `BeenAttacked` | 暂无 | 未覆盖 | 成员受击后通知怪群。 |
| 146 | 怪群系统 | OrderAttack | `OrderAttack` | 暂无 | 未覆盖 | 怪群按概率命令成员集火目标。 |
| 147 | 怪群系统 | Herd base act | `MonsterBaseAct` | `baseAction` 单怪 | 未覆盖 | 怪群级基础行为。 |
| 148 | 怪群系统 | Herd move action | `MonsterMoveAction` | 暂无 | 未覆盖 | 怪群级移动行为。 |
| 149 | 怪群系统 | Herd attack action | `MonsterAttackAction` | 暂无 | 未覆盖 | 怪群级攻击行为。 |
| 150 | 怪群系统 | Herd die action | `MonsterDieAction` | 暂无 | 未覆盖 | 成员死亡后的怪群处理。 |
| 151 | 怪群系统 | Herd regen action | `MonsterRegenAction` | `processRegen` 单怪 | 未覆盖 | 怪群成员重生策略。 |
| 152 | 怪物再生系统 | 普通死亡重生 | `gObjMonsterRegen` | `Object.processRegen`、`Monster.Regen` | 部分覆盖 | 补齐怪物死亡状态、AI 清理、站位、特殊地图重生。 |
| 153 | 怪物再生系统 | 移动重生 | `gObjMonsterMoveRegen` | `SpawnPosition` | 部分覆盖 | 指定位置重生或移动后重生。 |
| 154 | 怪物再生系统 | MonsterRegenSystem | `CMonsterRegenSystem` | 暂无 | 未覆盖 | 实现按组、时间表、区域刷 Boss/事件怪。 |
| 155 | 怪物再生系统 | Regen script 加载 | `LoadScript` | 暂无 | 未覆盖 | 加载再生组、区域、怪物列表、时间表、公告。 |
| 156 | 怪物再生系统 | Regen Run | `Run` | 暂无 | 未覆盖 | 周期检查是否到达再生时间。 |
| 157 | 怪物再生系统 | MonsterKillCheck | `MonsterKillCheck` | 怪物死亡事件未统一 | 未覆盖 | 怪物死亡后更新再生组存活数和 Boss 状态。 |
| 158 | 怪物再生系统 | RegenMonster | `RegenMonster` | `ObjectManager.AddMonster` 可用 | 未覆盖 | 按组批量创建怪物。 |
| 159 | 怪物再生系统 | SetPosMonster | `SetPosMonster` | `RandPosition` | 部分覆盖 | 在再生区域放置怪物。 |
| 160 | 怪物再生系统 | IsLiveBossState | `IsLiveBossState` | 暂无 | 未覆盖 | 避免 Boss 存活时重复刷。 |
| 161 | 怪物再生系统 | DeleteMonster | `DeleteMonster` | `ObjectManager` 无完整删除怪物接口 | 未覆盖 | 事件结束或持续时间到期删除怪物组。 |
| 162 | 怪物再生系统 | Spawn notice | `SendAllUserAnyMsg` | 公告系统未独立 | 未覆盖 | Boss 刷新或消失时广播。 |
| 163 | 怪物再生系统 | Regen time table | `_stRegenTimeTable` | 暂无 | 未覆盖 | 按小时、分钟、概率控制刷出。 |
| 164 | 特殊怪与陷阱 | 昆顿 AI | class 161/181/189/197/267/275 | `LearnSkill` 若干技能 | 部分覆盖 | 免疫、召唤、技能释放应迁入配置化 AI。 |
| 165 | 特殊怪与陷阱 | 暗黑巫师 AI | class 149/179/187/195/265/273/335 | `LearnSkill(1)` | 部分覆盖 | 补齐技能选择和特殊行为。 |
| 166 | 特殊怪与陷阱 | 美杜莎 AI | `gObjMonsterSelectSkillForMedusa` | `LearnSkill` 9/38/237/238 | 部分覆盖 | 美杜莎技能选择和阶段逻辑独立实现。 |
| 167 | 特殊怪与陷阱 | Lord Silvester 召唤 | `gObjMonsterSummonSkillForLordSilvester` | `LearnSkill` 占位 | 未覆盖 | 实现 Boss 召唤和召唤物清理。 |
| 168 | 特殊怪与陷阱 | Silvester 召唤物清理 | `KillLordSilvesterRecallMon` | 暂无 | 未覆盖 | Boss 死亡或阶段结束时清理召唤物。 |
| 169 | 特殊怪与陷阱 | 怪物特殊能力 | `gObjUseMonsterSpecialAbillity` | 暂无 | 未覆盖 | 统一处理特殊技能、免疫、召唤、传送等能力。 |
| 170 | 特殊怪与陷阱 | 陷阱怪行为 | `gObjMonsterTrapAct` | 暂无 | 未覆盖 | 固定陷阱怪不移动，只按方向或范围攻击。 |
| 171 | 特殊怪与陷阱 | 陷阱 X 搜敌 | `gObjTrapAttackEnemySearchX` | 暂无 | 未覆盖 | 横向陷阱攻击。 |
| 172 | 特殊怪与陷阱 | 陷阱 Y 搜敌 | `gObjTrapAttackEnemySearchY` | 暂无 | 未覆盖 | 纵向陷阱攻击。 |
| 173 | 特殊怪与陷阱 | 陷阱范围搜敌 | `gObjTrapAttackEnemySearchRange` | 暂无 | 未覆盖 | 范围陷阱攻击。 |
| 174 | 特殊怪与陷阱 | Quest NPC teleport | `CQeustNpcTeleport` | NPC 对话/移动未覆盖 | 未覆盖 | 如果保留该逻辑，应划入任务/NPC 与怪物行为边界。 |
| 175 | 跨系统接口 | 对象系统接口 | `OBJECTSTRUCT` monster fields | `object.Object` | 部分覆盖 | 对象系统提供怪物实例、坐标、视野、死亡、移动和攻击入口。 |
| 176 | 跨系统接口 | 地图系统接口 | `MapClass`、path/wall | `maps.MapManager` | 部分覆盖 | AI 请求地图合法点、阻挡、路径、安全区，不直接读地图细节。 |
| 177 | 跨系统接口 | 技能系统接口 | `CMagicInf` | `skill.Skill`、`UseSkill` | 部分覆盖 | 怪物技能最终调用技能系统执行通用技能效果。 |
| 178 | 跨系统接口 | 公式系统接口 | monster stat calc | `05-formula.md` | 未覆盖 | 怪物数值修正、Boss 属性倍率、元素属性等归公式。 |
| 179 | 跨系统接口 | Buff 系统接口 | skill element buffs | `13-buffs.md` | 未覆盖 | 怪物技能施加/移除状态，生命周期归 Buff。 |
| 180 | 跨系统接口 | 掉落系统接口 | `gObjMonsterDieGiveItem` | `14-drops.md` | 部分覆盖 | AI 只触发死亡事件和上下文，掉落生成归掉落系统。 |
| 181 | 跨系统接口 | 经验系统接口 | monster kill exp | `09-exp.md` | 未覆盖 | 提供击杀者、伤害归属和怪物 exp type。 |
| 182 | 跨系统接口 | 任务系统接口 | kill quest | `11-quests.md` | 未覆盖 | 怪物死亡后通知任务系统。 |
| 183 | 跨系统接口 | 副本系统接口 | DG/IG herd、BC/DS/CC 怪 | `21-dungeons.md` | 未覆盖 | 副本系统驱动事件状态，怪物行为由 AI 系统执行。 |
| 184 | 跨系统接口 | 普通活动接口 | 入侵怪、节日怪、怪群 | `22-events.md` | 未覆盖 | 活动创建怪物和触发特殊 AI。 |
| 185 | 跨系统接口 | 世界事件接口 | Crywolf/Raklion/Kanturu AI | `23-world-events.md` | 未覆盖 | 世界事件按阶段切换怪物状态或 AIUnit。 |
| 186 | 跨系统接口 | 宠物召唤接口 | summoned monster AI | `24-pets-summons.md` | 未覆盖 | 宠物系统创建召唤物，怪物 AI 负责召唤物行为。 |
| 187 | 协议与测试 | 基础 AI 测试 | `gObjMonsterBaseAct` 行为 | 待补测试 | 未覆盖 | 覆盖无目标、发现目标、追击、攻击、丢失目标、返回。 |
| 188 | 协议与测试 | 地图阻挡测试 | move/check wall | `maps/path_test.go` 只测路径 | 部分覆盖 | 覆盖怪物追击时墙体、安全区、不可站点处理。 |
| 189 | 协议与测试 | 仇恨测试 | `TMonsterAIAgro` | 待实现 | 未覆盖 | 覆盖增减仇恨、衰减、最高仇恨、目标失效清理。 |
| 190 | 协议与测试 | 状态机测试 | `TMonsterAIAutomata` | 待实现 | 未覆盖 | 覆盖 HP、敌人、仇恨和立即转移。 |
| 191 | 协议与测试 | AIUnit 测试 | `TMonsterAIUnit` | 待实现 | 未覆盖 | 覆盖八类行为槽选择和 delay。 |
| 192 | 协议与测试 | AIElement 测试 | `TMonsterAIElement` | 待实现 | 未覆盖 | 覆盖移动、攻击、治疗、召唤、免疫、AI 切换。 |
| 193 | 协议与测试 | 怪物技能测试 | `TMonsterSkillManager` | 待实现 | 未覆盖 | 覆盖技能单元、元素、延迟释放、范围目标和成功率。 |
| 194 | 协议与测试 | 怪群测试 | `MonsterHerd` | 待实现 | 未覆盖 | 覆盖创建、移动、受击联动、集火、死亡、重生。 |
| 195 | 协议与测试 | 再生系统测试 | `CMonsterRegenSystem` | 待实现 | 未覆盖 | 覆盖时间表、概率、Boss 存活、删除过期怪。 |
| 196 | 协议与测试 | 特殊怪测试 | Kuntun/Medusa/Silvester | 待实现 | 未覆盖 | 覆盖特殊技能、召唤、免疫和阶段逻辑。 |
| 197 | 协议与测试 | 陷阱怪测试 | trap search/action | 待实现 | 未覆盖 | 覆盖 X/Y/范围陷阱攻击。 |
| 198 | 协议与测试 | 死亡归属测试 | hit damage/top/last | 待实现 | 未覆盖 | 覆盖最高伤害、最后一击、离线清理、奖励归属输入。 |
| 199 | 协议与测试 | 跨系统边界测试 | object/map/skill/drop/event | 待实现 | 未覆盖 | 确认怪物 AI 不重复实现地图、技能、掉落、经验和活动规则。 |
| 200 | 协议与测试 | 回归兼容测试 | existing `Monster.ProcessAction` | 待实现 | 未覆盖 | 新 AI 未启用或配置缺失时，基础怪物仍可游走、追击和攻击。 |
