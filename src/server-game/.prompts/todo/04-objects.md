# 4. 对象系统

本模块覆盖 server-game 的对象运行时承载层，粒度按现有 Go 函数和 GameServer 对象系统语义拆分。对象系统只负责对象生命周期、状态承载、移动/视野/攻击/技能/道具/交互入口；完整地图规则、道具规则、公式、经验、技能定义、任务、合成、掉落表和社交系统分别归入其他正式模块。移动异常、攻击速度、技能距离、安全区违规和处罚状态归 `27-security.md`，对象系统在入口处调用安全检查。GM 移动、召唤、踢人、禁言等运营命令由 `29-ops.md` 发起，对象系统只执行对象级动作。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 对象管理器 | ObjectManager.init 索引段初始化 | user.cpp:gObjInit 与对象数组初始化语义 | game/object/object.go:objectManager.init | 已覆盖 | 校验怪物、召唤怪、玩家、user 的索引段边界和配置上限，补充越界测试。 |
| 2 | 对象管理器 | ObjectManager.AddMonster 怪物槽分配 | user.cpp:gObjAddMonster 与 gObjMonsterSetBase | game/object/object.go:objectManager.AddMonster | 部分覆盖 | 保留当前环形查找逻辑，补充满员错误、panic 边界和怪物计数一致性测试。 |
| 3 | 对象管理器 | ObjectManager.AddPlayer 玩家槽分配 | user.cpp:gObjAddSearch 与连接对象注册 | game/object/object.go:objectManager.AddPlayer | 部分覆盖 | 补齐满员响应、连接关闭策略、对象初始化失败回滚和 ConnectReply 测试。 |
| 4 | 对象管理器 | ObjectManager.DeletePlayer 玩家槽释放 | user.cpp:gObjDel 与 gObjCloseSet | game/object/object.go:objectManager.DeletePlayer | 部分覆盖 | 明确 Offline、Reset、视野清理、地图移除和 playerCount 递减的顺序。 |
| 5 | 对象管理器 | ObjectManager.GetObject 对象查找 | user.cpp:gObjIsConnected 与对象索引访问 | game/object/object.go:objectManager.GetObject | 已覆盖 | 补充负数索引保护、空槽返回约定和调用方 nil 处理测试。 |
| 6 | 对象管理器 | ObjectManager.GetPlayerByName 名称查找 | user.cpp:gObjFind | game/object/object.go:objectManager.GetPlayerByName | 部分覆盖 | 补充大小写规则、离线对象过滤、重名异常和性能边界说明。 |
| 7 | 对象管理器 | ObjectManager.GetPlayerPercent 在线百分比 | TServerInfoDisplayer.cpp 在线统计语义 | game/object/object.go:objectManager.GetPlayerPercent | 已覆盖 | 增加除零保护和配置 MaxPlayerCount 异常时的行为测试。 |
| 8 | 对象管理器 | ObjectManager.OfflineAllObjects 全员下线 | GameMain.cpp 关闭流程与 gObjAllLogOut 语义 | game/object/object.go:objectManager.OfflineAllObjects | 部分覆盖 | 补齐保存顺序、重复调用保护、下线超时和关闭期间拒绝新连接策略。 |
| 9 | 对象管理器 | ObjectManager.Process100ms 对象快速 Tick | GameMain.cpp 主循环与 gObjSecondProc 相关语义 | game/object/object.go:objectManager.Process100ms | 部分覆盖 | 明确每 100ms 处理移动、视野、延迟消息、怪物行为的顺序和耗时监控。 |
| 10 | 对象管理器 | ObjectManager.Process1000ms 对象慢 Tick | user.cpp:gObjSecondProc 与周期恢复语义 | game/object/object.go:objectManager.Process1000ms | 部分覆盖 | 明确玩家恢复、怪物慢逻辑、订阅发布和保存触发边界。 |
| 11 | Object 基础结构 | Object 字段集合与对象语义 | user.h:OBJECTSTRUCT | game/object/object.go:Object | 部分覆盖 | 对照 OBJECTSTRUCT 梳理 Go 版保留字段、迁移字段和不照搬字段。 |
| 12 | Object 基础结构 | Objecter 接口承载多态行为 | user.h:LPOBJ 类型和函数指针式语义 | game/object/object.go:Objecter | 部分覆盖 | 明确 Player、Monster、user 对 Push、Offline、Regen 等方法的职责差异。 |
| 13 | Object 基础结构 | Object.Init 初始化对象组件 | user.cpp:gObjClear 与初始化语义 | game/object/object.go:Object.Init | 已覆盖 | 确认技能、视野、延迟消息初始化顺序，补充重复 Init 测试。 |
| 14 | Object 基础结构 | Object.Reset 清理可复用对象 | user.cpp:gObjClear 与 gObjDel | game/object/object.go:Object.Reset | 部分覆盖 | 补齐移动路径、延迟消息、状态字段、地图容器和死亡状态清理。 |
| 15 | Object 基础结构 | Object.RandPosition 随机可站立坐标 | user.cpp:gObjGetRandomFreeLocation | game/object/object.go:Object.RandPosition | 部分覆盖 | 把地图属性位含义移到地图系统，并在对象层补充失败回退和日志测试。 |
| 16 | Object 基础结构 | Object.PushHPSD 推送 HP/SD | protocol.cpp:GCLifeSend | game/object/object.go:Object.PushHPSD | 已覆盖 | 确认 Position 语义、死亡状态下推送策略和最大值同步时机。 |
| 17 | Object 基础结构 | Object.PushMPAG 推送 MP/AG | protocol.cpp:GCManaSend | game/object/object.go:Object.PushMPAG | 已覆盖 | 补充 AG 为零、MP 溢出和客户端显示兼容测试。 |
| 18 | Object 基础结构 | Object.PushMoney 推送金币 | protocol.cpp:GCMoneySend | game/object/object.go:Object.PushMoney | 已覆盖 | 明确金币变化来源、溢出保护和对象道具入口调用约定。 |
| 19 | Object 基础结构 | Object.PushSystemMsg 系统消息 | TNotice.cpp 与 MsgOutput | game/object/object.go:Object.PushSystemMsg | 部分覆盖 | 统一系统消息来源、语言包接入和聊天通道复用策略。 |
| 20 | Object 基础结构 | Object.PushWeather 天气推送 | protocol.cpp:GCWeatherSend | game/object/object.go:Object.PushWeather | 已覆盖 | 明确天气由地图系统驱动，对象层只负责单对象发送。 |
| 21 | 延迟消息与重生状态 | Object.AddDelayMsg 延迟事件注册 | user.cpp:gObjAddMsgSendDelay | game/object/object.go:Object.AddDelayMsg | 部分覆盖 | 替换固定 20 槽或补充满槽行为，定义延迟事件枚举和日志。 |
| 22 | 延迟消息与重生状态 | Object.processDelayMsg 处理死亡掉落 | gObjMonster.cpp:gObjMonsterDieGiveItem | game/object/object.go:Object.processDelayMsg | 部分覆盖 | 明确延迟掉落只触发入口，掉落表和 Bag 规则归 14-drops。 |
| 23 | 延迟消息与重生状态 | Object.processDelayMsg 处理击退 | ObjUseSkill.cpp:SkillKnightRush 与击退技能语义 | game/object/object.go:Object.processDelayMsg | 部分覆盖 | 统一延迟击退和即时 Knockback 的调用边界。 |
| 24 | 延迟消息与重生状态 | Object.processDelayMsg 处理吸血回蓝 | ObjAttack.cpp 击杀恢复语义 | game/object/object.go:Object.processDelayMsg | 部分覆盖 | 确认击杀者仍在线、目标仍有效和恢复值公式归属。 |
| 25 | 延迟消息与重生状态 | Object.processRegen 死亡后重生 | user.cpp:gObjLifeCheck 与 gObjViewportClose | game/object/object.go:Object.processRegen | 部分覆盖 | 明确 State 4 到 State 8，再到 Regen 的阶段含义。 |
| 26 | 延迟消息与重生状态 | Object.dieTime 死亡时间记录 | gObjMonster.cpp:gObjMonsterRegen | game/object/object.go:Object.dieTime | 部分覆盖 | 统一玩家和怪物死亡时间用途，补充断线和地图切换边界。 |
| 27 | 延迟消息与重生状态 | Object.dieRegen 重生开关 | gObjMonster.cpp:gObjMonsterRegen | game/object/object.go:Object.dieRegen | 部分覆盖 | 定义哪些对象允许自动重生，NPC 和特殊怪物应明确例外。 |
| 28 | 延迟消息与重生状态 | Object.State 对象阶段状态 | user.h:OBJECTSTRUCT State 字段 | game/object/object.go:Object.State | 部分覆盖 | 把 1 初始、2 视野、4 死亡、8 清理定义成常量并补测试。 |
| 29 | 延迟消息与重生状态 | Object.Live 存活状态 | user.h:OBJECTSTRUCT Live 字段 | game/object/object.go:Object.Live | 部分覆盖 | 统一攻击、移动、视野、道具入口对 Live 的校验。 |
| 30 | 延迟消息与重生状态 | Object.MaxRegenTime 重生时间 | gObjMonster.cpp:MaxRegenTime 语义 | game/object/object.go:Object.MaxRegenTime | 部分覆盖 | 区分怪物、玩家、召唤物、活动对象的重生时长来源。 |
| 31 | Player 运行时 | player.NewPlayer 注册玩家对象 | user.cpp:gObjAddSearch 与玩家创建 | game/object/player/player.go:NewPlayer | 部分覆盖 | 明确 actioner、conn、channel、context 的生命周期和失败回滚。 |
| 32 | Player 运行时 | player.newPlayer 初始化 Player | user.cpp:gObjCharZeroSet | game/object/player/player.go:newPlayer | 部分覆盖 | 补齐默认连接状态、对象类型、恢复 Tick、背包仓库初始化说明。 |
| 33 | Player 运行时 | Player.Addr 返回连接地址 | wsGameServer.cpp 连接地址日志语义 | game/object/player/player.go:Player.Addr | 已覆盖 | 补充 conn nil 时的行为，避免离线流程 panic。 |
| 34 | Player 运行时 | Player.Offline 离线流程 | user.cpp:gObjCloseSet | game/object/player/player.go:Player.Offline | 部分覆盖 | 明确 cancel、conn.Close、保存角色、视野销毁和重复 Offline 保护。 |
| 35 | Player 运行时 | Player.Push 发送消息 | protocol.cpp:DataSend | game/object/player/player.go:Player.Push | 部分覆盖 | 补充写失败处理、离线丢弃策略和发送队列背压约定。 |
| 36 | Player 运行时 | Player.Process1000ms 慢 Tick | user.cpp:gObjSecondProc | game/object/player/player.go:Player.Process1000ms | 部分覆盖 | 明确恢复、自动保存、Buff、MuBot 等慢 Tick 的扩展点。 |
| 37 | Player 运行时 | Player.recoverHPSD 生命护盾恢复 | user.cpp:gProcessAutoRecuperation 与 gObjShieldAutoRefill | game/object/player/player.go:Player.recoverHPSD | 部分覆盖 | 将恢复公式归 05-formula，对象层补状态校验和推送规则。 |
| 38 | Player 运行时 | Player.recoverMPAG 魔法 AG 恢复 | user.cpp:gObjManaCheck 与自动恢复语义 | game/object/player/player.go:Player.recoverMPAG | 部分覆盖 | 明确恢复频率、最大值裁剪、死亡状态跳过和协议推送。 |
| 39 | Player 运行时 | Player.SpawnPosition 玩家出生位置 | user.cpp:gObjSetMapMove 与 Regen 位置 | game/object/player/player.go:Player.SpawnPosition | 部分覆盖 | 对接地图系统的复活点、死亡点、登录点和活动地图例外。 |
| 40 | Player 运行时 | Player.Regen 玩家重生属性恢复 | user.cpp:gObjRegenLive | game/object/player/player.go:Player.Regen | 部分覆盖 | 明确 HP/MP/SD/AG 恢复比例、Buff 清理和视野重建时机。 |
| 41 | Player 角色承载入口 | Player.Login 登录入口 | wsJoinServerCli.cpp 与 login protocol | game/object/player/player.go:Player.Login | 部分覆盖 | 对象层只负责绑定账号状态，认证细节归 02-accounts。 |
| 42 | Player 角色承载入口 | Player.Logout 登出入口 | user.cpp:gObjCloseSet | game/object/player/player.go:Player.Logout | 部分覆盖 | 明确登出包、断线、主动退出和服务器关闭的统一路径。 |
| 43 | Player 角色承载入口 | Player.GetCharacterList 角色列表入口 | DSProtocol.cpp:角色列表请求语义 | game/object/player/player.go:Player.GetCharacterList | 部分覆盖 | 对象层保留协议入口，角色查询规则归 03-characters。 |
| 44 | Player 角色承载入口 | Player.CreateCharacter 创建角色入口 | DSProtocol.cpp:角色创建请求语义 | game/object/player/player.go:Player.CreateCharacter | 部分覆盖 | 对象层保留调用入口，职业开放和初始数据归角色系统。 |
| 45 | Player 角色承载入口 | Player.DeleteCharacter 删除角色入口 | DSProtocol.cpp:角色删除请求语义 | game/object/player/player.go:Player.DeleteCharacter | 部分覆盖 | 补齐在线角色限制、密码校验调用和错误码映射。 |
| 46 | Player 角色承载入口 | Player.CheckCharacter 检查角色入口 | DSProtocol.cpp:角色检查语义 | game/object/player/player.go:Player.CheckCharacter | 部分覆盖 | 明确对象层仅转发检查结果，不承载 DB 规则。 |
| 47 | Player 角色承载入口 | Player.LoadCharacter 加载角色入口 | DSProtocol.cpp:角色加载请求语义 | game/object/player/player.go:Player.LoadCharacter | 部分覆盖 | 明确加载后对象状态、地图注册、视野创建和数据下发顺序。 |
| 48 | Player 角色承载入口 | Player.saveCharacter 保存角色 | DbSave.cpp 与 DSProtocol.cpp 保存语义 | game/object/player/player.go:Player.saveCharacter | 部分覆盖 | 补齐保存触发点、失败日志、关闭等待和字段完整性测试。 |
| 49 | Player 角色承载入口 | Player.EquipmentChanged 装备变化入口 | user.cpp:gObjCalCharacter 与装备变更 | game/object/player/player.go:Player.EquipmentChanged | 部分覆盖 | 对象层触发重算和广播，装备规则归 07-items，公式归 05-formula。 |
| 50 | Player 角色承载入口 | Player.Talk NPC 对话分发 | NpcTalk.cpp:NpcTalk | game/object/player/player.go:Player.Talk | 部分覆盖 | 对象层保留 NPC 距离和类型入口，具体商店/任务/合成归对应模块。 |
| 51 | Monster 运行时 | monster.SpawnMonster 加载刷怪 | MonsterSetBase.cpp 与 gObjSetMonster | game/object/monster/monster.go:SpawnMonster | 部分覆盖 | 明确配置加载失败、重复刷怪、地图无效和对象注册错误处理。 |
| 52 | Monster 运行时 | monsterTable.init 怪物属性表 | MonsterAttr.cpp | game/object/monster/monster.go:monsterTable.init | 部分覆盖 | 怪物属性读取可保留在对象层，复杂数值修正归公式或 `25-monster-ai.md`。 |
| 53 | Monster 运行时 | newMonster 构造怪物对象 | gObjMonster.cpp:gObjSetMonster | game/object/monster/monster.go:newMonster | 部分覆盖 | 补齐 Object.Init 调用、技能学习、掉落字段、NPC 类型和 Live 状态测试。 |
| 54 | Monster 运行时 | Monster.SpawnPosition 怪物出生坐标 | gObjMonster.cpp:gObjMonsterRegen | game/object/monster/monster.go:Monster.SpawnPosition | 部分覆盖 | 对接地图可站立检查和特殊地图例外，补充随机失败日志。 |
| 55 | Monster 运行时 | Monster.Regen 怪物重生恢复 | gObjMonster.cpp:gObjMonsterRegen | game/object/monster/monster.go:Monster.Regen | 部分覆盖 | 补齐 MP、目标、状态机、视野、延迟消息和掉落归属清理。 |
| 56 | Monster 运行时 | Monster.ProcessAction 怪物行为 Tick | gObjMonster.cpp:gObjMonsterProcess | game/object/monster/monster.go:Monster.ProcessAction | 部分覆盖 | 明确 Tick 间隔、Live 检查、ConnectState 检查和耗时监控。 |
| 57 | Monster 运行时 | Monster.Process1000ms 怪物慢 Tick | gObjMonster.cpp:gObjMonsterProcess 周期逻辑 | game/object/monster/monster.go:Monster.Process1000ms | 未覆盖 | 预留回血、特殊怪物状态、活动怪清理和召唤物寿命处理。 |
| 58 | Monster 运行时 | Monster.GetSkillMPAG 怪物技能消耗 | ObjUseSkill.cpp:GetUseMana 与怪物技能语义 | game/object/monster/monster.go:Monster.GetSkillMPAG | 部分覆盖 | 明确怪物技能默认不消耗资源，特殊 Boss 例外归 `25-monster-ai.md`。 |
| 59 | Monster 运行时 | Monster.Offline 怪物离线空实现 | gObjMonster.cpp:怪物对象无连接语义 | game/object/monster/monster.go:Monster.Offline | Go 化替代 | 保持空实现，但说明怪物删除应走 ObjectManager 槽释放而非连接离线。 |
| 60 | Monster 运行时 | Monster.Push 怪物发送空实现 | GameServer 怪物无客户端连接语义 | game/object/monster/monster.go:Monster.Push | Go 化替代 | 保持空实现，怪物广播必须通过视野或地图广播完成。 |
| 61 | Monster 行为函数 | Monster.overDis 出生范围判断 | gObjMonster.cpp:gObjMonsterMoveCheck | game/object/monster/monster.go:Monster.overDis | 部分覆盖 | 补齐 spawnDis 为 0、地图边界和追击回出生点策略。 |
| 62 | Monster 行为函数 | Monster.searchEnemy 视野选敌 | gObjMonster.cpp:gObjMonsterSearchEnemy | game/object/monster/monster.go:Monster.searchEnemy | 部分覆盖 | 补齐仇恨、最后攻击者、PK 条件、隐身目标和安全区策略。 |
| 63 | Monster 行为函数 | Monster.roamMove 随机游走 | gObjMonster.cpp:gObjMonsterMoveAction | game/object/monster/monster.go:Monster.roamMove | 部分覆盖 | 明确守卫、安全区、阻挡属性和移动范围失败回退。 |
| 64 | Monster 行为函数 | Monster.chaseMove 追击移动 | gObjMonster.cpp:gObjMonsterGetTargetPos | game/object/monster/monster.go:Monster.chaseMove | 部分覆盖 | 补齐寻路失败、目标阻挡、距离过远和返回出生点策略。 |
| 65 | Monster 行为函数 | Monster.baseAction 行为状态机 | gObjMonster.cpp:gObjMonsterBaseAct | game/object/monster/monster.go:Monster.baseAction | 部分覆盖 | 将 emotion、move、attack 状态定义成常量并补状态转移测试。 |
| 66 | Monster 行为函数 | Monster.move 寻路移动执行 | gObjMonster.cpp:PathFindMoveMsgSend | game/object/monster/monster.go:Monster.move | 部分覆盖 | 明确 FindMapPath 失败、路径方向生成和 Move 包广播。 |
| 67 | Monster 行为函数 | Monster.attack 普攻或技能选择 | gObjMonster.cpp:gObjMonsterAttack 与 gObjMonsterMagicAttack | game/object/monster/monster.go:Monster.attack | 部分覆盖 | 补齐技能权重、攻击冷却、目标失效和普通攻击概率配置。 |
| 68 | Monster 行为函数 | Monster.TargetNumber 目标字段维护 | gObjMonster.cpp:TargetNumber 相关逻辑 | game/object/object.go:Object.TargetNumber | 部分覆盖 | 明确目标重置、死亡清理、跨地图清理和并发访问边界。 |
| 69 | Monster 行为函数 | Monster.actionState 行为状态字段 | gObjMonster.cpp:m_ActState 语义 | game/object/monster/monster.go:actionState | 部分覆盖 | 对齐 GameServer actionState 中 rest、attack、move、emotion 的语义。 |
| 70 | Monster 行为函数 | Monster 特殊行为扩展点 | TMonsterAI.cpp 与 TMonsterSkillManager.cpp | game/object/monster/monster.go:ProcessAction | 未覆盖 | 对象层保留扩展入口，完整 AI、仇恨和特殊技能归 `25-monster-ai.md`。 |
| 71 | 对象移动 | Object.CalcDistance 距离计算 | user.cpp:gObjCalDistanceTX | game/object/move.go:Object.CalcDistance | 已覆盖 | 统一与地图系统距离算法，明确欧式或格子距离使用场景。 |
| 72 | 对象移动 | Object.processMove 路径推进 | user.cpp:gObjSetPosition 与移动 Tick | game/object/move.go:Object.processMove | 部分覆盖 | 补齐路径越界、速度间隔、移动结束广播和异常路径清理。 |
| 73 | 对象移动 | Object.Move 客户端移动入口 | protocol.cpp:CGMoveRecv | game/object/move.go:Object.Move | 部分覆盖 | 明确路径包解析、`27-security.md` 移动校验调用、坐标更新和 MoveReply 广播。 |
| 74 | 对象移动 | Object.LoadMiniMap 小地图加载入口 | MinimapData.cpp 与 SendNPCInfo.cpp | game/object/move.go:Object.LoadMiniMap | 部分覆盖 | 对象层只负责触发下发，地图点位数据归 06-maps。 |
| 75 | 对象移动 | Object.gateMove Gate 移动 | user.cpp:gObjMoveGate | game/object/move.go:Object.gateMove | 部分覆盖 | 明确 Gate 查找、目标坐标、视野清理、地图容器迁移和失败回复。 |
| 76 | 对象移动 | Object.Teleport 传送入口 | user.cpp:gObjTeleport | game/object/move.go:Object.Teleport | 部分覆盖 | 补齐安全区、地图限制、传送冷却和客户端同步。 |
| 77 | 对象移动 | Object.MapMove 地图移动命令 | MoveCommand.cpp 与 gObjMoveGate | game/object/move.go:Object.MapMove | 部分覆盖 | 地图命令规则归 06-maps，对象层负责状态切换和协议回复。 |
| 78 | 对象移动 | Object.SetPosition 强制设坐标 | user.cpp:gObjSetPosition | game/object/move.go:Object.SetPosition | 部分覆盖 | GM 强制移动入口归 `29-ops.md`，对象层负责坐标设置、权限结果执行和广播。 |
| 79 | 对象移动 | Object.Knockback 击退 | ObjUseSkill.cpp:SkillKnightRush | game/object/move.go:Object.Knockback | 部分覆盖 | 补齐阻挡检查、地图边界、击退失败和技能延迟调用。 |
| 80 | 对象移动 | Object.Action 动作广播 | protocol.cpp:CGActionRecv | game/object/move.go:Object.Action | 部分覆盖 | 明确动作状态是否影响攻击、移动、交互和客户端表现。 |
| 81 | 对象视野 | Object.CreateFrustum 视野几何初始化 | user.cpp:SkillFrustrum 与视野方向语义 | game/object/viewport.go:Object.CreateFrustum | 部分覆盖 | 区分普通视野几何和技能扇形几何，补充方向边界测试。 |
| 82 | 对象视野 | Object.initViewport 初始化视野容器 | user.cpp:gObjViewportListCreate | game/object/viewport.go:Object.initViewport | 已覆盖 | 补充重复初始化、数组清空和 ViewportsNum 一致性测试。 |
| 83 | 对象视野 | Object.checkViewport 坐标可见性 | user.cpp:gObjViewportListCreate 距离检查 | game/object/viewport.go:Object.checkViewport | 部分覆盖 | 明确 ViewRange、地图一致、死亡状态和隐身状态的过滤规则。 |
| 84 | 对象视野 | Object.addViewport 主动视野加入 | user.cpp:gObjViewportListProtocolCreate | game/object/viewport.go:Object.addViewport | 部分覆盖 | 补齐重复加入、容量满、对象类型和距离字段更新策略。 |
| 85 | 对象视野 | Object.addViewportPassive 被动推送列表 | user.cpp:gObjViewportListProtocolCreate | game/object/viewport.go:Object.addViewportPassive | 部分覆盖 | 明确攻击用视野和推送用视野的区别，避免双向不一致。 |
| 86 | 对象视野 | Object.removeViewport 移除单个视野 | user.cpp:gObjViewportClose | game/object/viewport.go:Object.removeViewport | 部分覆盖 | 补齐顺序压缩、计数一致性和销毁消息触发边界。 |
| 87 | 对象视野 | Object.clearViewport 清空视野 | user.cpp:gObjViewportClose | game/object/viewport.go:Object.clearViewport | 部分覆盖 | 明确下线、地图切换、死亡和重生时的清理顺序。 |
| 88 | 对象视野 | Object.createViewport 创建视野消息 | protocol.cpp:GCViewportListProtocolCreate | game/object/viewport.go:Object.createViewport | 部分覆盖 | 完善玩家、怪物、NPC、地图物品不同 CreateViewport 包。 |
| 89 | 对象视野 | Object.destroyViewport 销毁视野消息 | protocol.cpp:GCViewportClose | game/object/viewport.go:Object.destroyViewport | 部分覆盖 | 补齐离开视野、死亡、拾取物品、地图切换的销毁原因。 |
| 90 | 对象视野 | Object.processViewport 差量处理 | user.cpp:gObjViewportListProtocol | game/object/viewport.go:Object.processViewport | 部分覆盖 | 明确 create/destroy 调用频率、排序稳定性和性能边界。 |
| 91 | 对象攻击死亡 | Object.CheckMiss 命中判定入口 | ObjBaseAttack.cpp:MissCheck 与 MissCheckPvP | game/object/attack.go:Object.CheckMiss | 部分覆盖 | 公式归 05-formula，对象层补目标类型、PVP/PVM 分支和随机测试。 |
| 92 | 对象攻击死亡 | Object.getDefense 防御读取 | ObjBaseAttack.cpp:GetTargetDefense | game/object/attack.go:Object.getDefense | 部分覆盖 | 明确技能类型、忽防、SD、防御公式调用和最小伤害保护。 |
| 93 | 对象攻击死亡 | Object.getDamage 伤害读取 | ObjAttack.cpp:gObjAttack 与 SkillAttack | game/object/attack.go:Object.getDamage | 部分覆盖 | 将技能倍率、卓越幸运、装备修正交给公式，保留对象层调用顺序。 |
| 94 | 对象攻击死亡 | Object.attack 伤害落地 | ObjAttack.cpp:gObjAttack | game/object/attack.go:Object.attack | 部分覆盖 | 明确扣血、死亡、广播、延迟消息和经验掉落触发顺序。 |
| 95 | 对象攻击死亡 | Object.Attack 普攻协议入口 | protocol.cpp:CGAttackRecv | game/object/attack.go:Object.Attack | 部分覆盖 | 补齐目标查找、Live 校验，并调用 `27-security.md` 做距离、攻速、安全区限制和异常处理。 |
| 96 | 对象攻击死亡 | Object.DieGiveExperience 死亡给经验入口 | gObjMonster.cpp:gObjMonsterDieGiveItem 与经验语义 | game/object/attack.go:Object.DieGiveExperience | 部分覆盖 | 对象层只触发经验系统接口，分配和倍率归 09-exp。 |
| 97 | 对象攻击死亡 | Object.DieDropExperience 玩家死亡掉经验 | user.cpp:玩家死亡惩罚语义 | game/object/attack.go:Object.DieDropExperience | 未覆盖 | 明确 PK、地图、等级、活动等死亡惩罚触发点。 |
| 98 | 对象攻击死亡 | Object.DieRecoverHPMP 击杀恢复 | ObjAttack.cpp 击杀恢复语义 | game/object/attack.go:Object.DieRecoverHPMP | 部分覆盖 | 恢复比例归公式，对象层负责目标校验和推送。 |
| 99 | 对象攻击死亡 | Player.Die 玩家死亡实现 | user.cpp:gObjUserDie | game/object/player/attack.go:Player.Die | 部分覆盖 | 明确玩家死亡状态、复活点、掉落惩罚和 PVP 结果。 |
| 100 | 对象攻击死亡 | Monster.Die 怪物死亡实现 | gObjMonster.cpp:gObjMonsterDieGiveItem | game/object/monster/attack.go:Monster.Die | 部分覆盖 | 明确怪物死亡状态、重生定时、经验入口和掉落入口。 |
| 101 | 对象技能入口 | Object.initSkill 初始化技能容器 | user.cpp:gObjMagicAdd 与 Magic 数组初始化 | game/object/skill.go:Object.initSkill | 已覆盖 | 补充重复初始化和登录加载技能前后的容器状态测试。 |
| 102 | 对象技能入口 | Object.clearSkill 清空技能容器 | user.cpp:gObjMagicDel 与对象清理 | game/object/skill.go:Object.clearSkill | 已覆盖 | 明确下线、重置、怪物删除时是否清空技能。 |
| 103 | 对象技能入口 | Object.LearnSkill 学习技能入口 | user.cpp:gObjMagicAdd | game/object/skill.go:Object.LearnSkill | 部分覆盖 | 学习规则归 10-skills，对象层只负责容器变更和错误日志。 |
| 104 | 对象技能入口 | Object.ForgetSkill 遗忘技能入口 | user.cpp:gObjMagicDel | game/object/skill.go:Object.ForgetSkill | 部分覆盖 | 补齐不存在技能处理、协议通知和持久化触发边界。 |
| 105 | 对象技能入口 | Object.UseSkill 技能使用入口 | ObjUseSkill.cpp:CObjUseSkill::UseSkill | game/object/skill.go:Object.UseSkill | 部分覆盖 | 对象层负责目标查找和通用入口，技能规则归 10-skills。 |
| 106 | 对象技能入口 | Object.canUseSkill 技能分发 | ObjUseSkill.cpp:CObjUseSkill::RunningSkill | game/object/skill.go:Object.canUseSkill | 部分覆盖 | 当前 switch 应逐步迁移到技能系统，保留对象层 dispatch 边界。 |
| 107 | 对象技能入口 | Object.UseSkillReply 技能广播 | protocol.cpp:GCMagicAttackNumberSend | game/object/skill.go:Object.UseSkillReply | 已覆盖 | 补齐失败原因、目标无效、范围技能多目标回复策略。 |
| 108 | 对象技能入口 | Object.getAngle 技能角度 | ObjUseSkill.cpp:GetAngle | game/object/skill.go:Object.getAngle | 已覆盖 | 和地图数学工具统一，补充八方向边界测试。 |
| 109 | 对象技能入口 | Object.CreateSkillFrustum 技能扇形 | ObjUseSkill.cpp:SkillFrustrum | game/object/skill.go:Object.CreateSkillFrustum | 部分覆盖 | 技能范围规则归 10-skills，对象层保留几何计算入口。 |
| 110 | 对象技能入口 | Object.UseSkillDeathStab 具体技能入口 | ObjUseSkill.cpp:SkillKnightBlow 与相关技能语义 | game/object/skill.go:Object.UseSkillDeathStab | 部分覆盖 | 作为迁移样例，明确从对象层移动到技能系统的目标边界。 |
| 111 | 对象道具与交互入口 | Object.PickItem 拾取地图物品 | user.cpp:gObjInventoryInsertItem 与 MapItem.cpp | game/object/item.go:Object.PickItem | 部分覆盖 | 对象层负责拾取入口，背包规则归 07-items，归属和掉落归 14-drops。 |
| 112 | 对象道具与交互入口 | Object.DropItem 丢弃物品 | user.cpp:gObjInventoryItemDrop | game/object/item.go:Object.DropItem | 部分覆盖 | 补齐地图限制、物品限制、视野创建和失败回滚。 |
| 113 | 对象道具与交互入口 | Object.BuyItem 购买入口 | Shop.cpp 与 NPC 商店购买语义 | game/object/item.go:Object.BuyItem | 部分覆盖 | 商店规则归 08-shops，对象层负责背包和金币变更入口。 |
| 114 | 对象道具与交互入口 | Object.SellItem 出售入口 | Shop.cpp 与出售语义 | game/object/item.go:Object.SellItem | 部分覆盖 | 明确出售限制调用、金币上限和删除物品顺序。 |
| 115 | 对象道具与交互入口 | Object.MoveItem 背包移动入口 | user.cpp:gObjInventoryMove | game/object/item.go:Object.MoveItem | 部分覆盖 | 物品占格归 07-items，对象层负责协议入口和装备变化触发。 |
| 116 | 对象道具与交互入口 | Object.UseItem 使用物品入口 | user.cpp:gObjUseItem 与各种 Use 函数 | game/object/item.go:Object.UseItem | 部分覆盖 | 具体药水、宝石、卷轴效果归道具或合成模块，对象层保留分发。 |
| 117 | 对象道具与交互入口 | Object.RepairItem 修理入口 | user.cpp:gObjRepairItem | game/object/item.go:Object.RepairItem | 部分覆盖 | 修理价格归 05-formula，商店/NPC 限制归 08-shops。 |
| 118 | 对象道具与交互入口 | Object.ValidateItemPosition 位置校验 | user.cpp:背包位置校验语义 | game/object/item.go:Object.ValidateItemPosition | 部分覆盖 | 明确背包、装备、仓库、扩展格的合法范围和错误码。 |
| 119 | 对象道具与交互入口 | Object.Chat 与 Object.Whisper | protocol.cpp:聊天协议与 WhisperCash.cpp | game/object/chat.go:Object.Chat,Object.Whisper | 部分覆盖 | 对象层保留普通聊天和私聊入口；禁言运营入口归 `29-ops.md`，过滤审计归 `27-security.md`。 |
| 120 | 对象道具与交互入口 | Player.Talk 与窗口关闭入口 | NpcTalk.cpp:NpcTalk 与 Warehouse 入口 | game/object/player/player.go:Player.Talk,CloseTalkWindow,CloseWarehouseWindow | 部分覆盖 | 对象层保留 NPC 交互状态，商店、仓库、任务、合成分别归对应模块。 |
