# 26. 脚本系统

本模块覆盖脚本基础设施：Lua VM 生命周期、Lua 文件加载、Lua 函数调用、Go/Lua 绑定、WZ/SMD 文本脚本解析、脚本错误处理、脚本热重载、自检和测试工具。脚本系统不拥有公式、掉落、任务、Buff、副本、活动、怪物 AI 等业务语义；这些模块只调用脚本系统执行脚本或解析配置。GameServer 中 `MuLua`、`LuaFun`、`LuaExport`、`LuaBag`、`QuestExpLuaBind`、`DoppelgangerLua`、`BuffScriptLoader`、`ReadScript/WzMemScript` 均存在真实调用链，server-game 当前只有公式包私有 Lua wrapper，尚无通用脚本系统。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 模块边界与总入口 | ScriptManager 总管理器 | `MULua` 分散在各业务对象中 | 暂无 `game/script` | 未覆盖 | 建立统一脚本基础设施，管理 Lua VM、脚本路径、加载、自检和调用。 |
| 2 | 模块边界与总入口 | 脚本系统初始化 | `GameMain` 多处 `DoFile/Load` | `formula.init` 私有加载 | 部分覆盖 | 启动时统一初始化脚本运行时，并保留业务模块自有加载顺序。 |
| 3 | 模块边界与总入口 | 脚本路径解析 | `g_ConfigRead.GetPath` | `conf.PathCommon` | 部分覆盖 | 统一处理 `Scripts/`、`ItemBags/`、事件配置等脚本路径。 |
| 4 | 模块边界与总入口 | 脚本模块注册表 | 分散全局对象 | 暂无 | 未覆盖 | 记录公式、Bag、Quest、DoppelGanger 等脚本模块及其 VM。 |
| 5 | 模块边界与总入口 | 脚本健康检查 | GameServer 启动日志 | 暂无 | 未覆盖 | 启动期校验脚本文件存在、关键函数存在、绑定函数完整。 |
| 6 | 模块边界与总入口 | 脚本错误统一日志 | `g_Log.Add` | `slog.Error` 分散 | 部分覆盖 | 标准化脚本文件、函数名、参数签名、调用者和错误信息。 |
| 7 | 模块边界与总入口 | 脚本重载入口 | BagManager 重载式 Init | 暂无 | 未覆盖 | 支持开发期脚本重载，线上默认关闭或需 GM 权限。 |
| 8 | 模块边界与总入口 | 脚本关闭清理 | `Release`、析构 | 暂无 | 未覆盖 | 进程退出或重载时关闭 VM、释放注册对象和缓存。 |
| 9 | Lua 运行时 | Lua VM 创建 | `MULua::Create` | `lua.NewState` | 部分覆盖 | 封装 gopher-lua VM 创建，支持按模块独立 VM。 |
| 10 | Lua 运行时 | Lua VM 释放 | `MULua::Release` | 暂无 Close 调用 | 未覆盖 | 在重载或关闭时调用 `L.Close()`。 |
| 11 | Lua 运行时 | DoFile | `MULua::DoFile(const char*)` | `l.DoFile(file)` | 部分覆盖 | 封装文件加载，返回 error 而不是直接退出进程。 |
| 12 | Lua 运行时 | DoFile buffer | `MULua::DoFile(char* buff, size)`、`ScriptLuaMem` | 暂无 | 未覆盖 | 支持从内存加载 Lua 脚本，便于加密脚本或测试。 |
| 13 | Lua 运行时 | DoString | `MULua::DoString` | 暂无 | 未覆盖 | 支持测试和 GM/调试场景执行字符串脚本。 |
| 14 | Lua 运行时 | GetLua | `MULua::GetLua` | 直接持有 `*lua.LState` | 部分覆盖 | 提供受控访问 Lua state，避免业务随意操作 VM。 |
| 15 | Lua 运行时 | 同步锁 | `MULua(bool UseSync)`、`CRITICAL_SECTION` | 暂无 | 未覆盖 | 明确 gopher-lua 非并发安全，按 VM 加 mutex 或串行调用。 |
| 16 | Lua 运行时 | 多 VM 策略 | 每业务对象一个 `MULua` | formula 多个 LState | 部分覆盖 | 保持公式、Bag、Quest、副本可隔离，避免全局污染。 |
| 17 | Lua 运行时 | Lua 标准库策略 | `luaL_openlibs` | gopher-lua 默认库 | 部分覆盖 | 决定开放哪些库，线上禁用危险 IO/debug/package 能力。 |
| 18 | Lua 运行时 | Lua panic 保护 | `Generic_Call` 错误日志 | `Protect: true` | 部分覆盖 | 统一 panic/error 转换，不让脚本错误中断主循环。 |
| 19 | Lua 运行时 | Lua 控制台 | `CreateWinConsole` | 不适用 | 不适用 | Windows 调试控制台不迁移，Go 侧用日志和调试接口替代。 |
| 20 | Lua 调用封装 | Generic_Call 入口 | `MULua::Generic_Call`、`g_Generic_Call` | `formula.call` | 部分覆盖 | 抽出通用 `Call`，公式包不再私有实现签名解析。 |
| 21 | Lua 调用封装 | 签名字符串解析 | `"iiii>iiii"` | `strings.Split(sig, ">")` | 部分覆盖 | 支持输入输出分隔、参数数量检查和未知类型报错。 |
| 22 | Lua 调用封装 | int 参数 | `i` | `lua.LNumber` | 已覆盖 | 保持 int 到 Lua number 的映射。 |
| 23 | Lua 调用封装 | float/double 参数 | `d` | `lua.LNumber` | 已覆盖 | 保持 float64 到 Lua number 的映射。 |
| 24 | Lua 调用封装 | int64 参数 | `l` | 暂无 | 未覆盖 | DoppelGanger callback 使用 tick/time，需要 int64 支持。 |
| 25 | Lua 调用封装 | bool 参数 | C++ 可扩展 | 暂无 | 未覆盖 | 支持脚本配置开关和条件调用。 |
| 26 | Lua 调用封装 | string 参数 | C++ 可扩展 | 暂无 | 未覆盖 | 支持脚本文件名、角色名、事件名等。 |
| 27 | Lua 调用封装 | table 参数 | C++ 分散绑定 | 暂无 | 未覆盖 | 复杂参数优先封装为 table 或 userdata。 |
| 28 | Lua 调用封装 | 返回值数量校验 | `Generic_Call` 签名 | `NRet` 存在 | 部分覆盖 | Lua 返回数量不足时应报错。 |
| 29 | Lua 调用封装 | 返回值类型校验 | C++ 日志 | formula 只记录部分错误 | 部分覆盖 | 类型不匹配必须返回 error 给业务模块。 |
| 30 | Lua 调用封装 | 调用失败策略 | `Generic_Call` 返回 bool | formula 只打日志 | 需修正 | 调用失败时返回 error，业务决定降级、跳过或启动失败。 |
| 31 | Lua 调用封装 | 保护调用 | `lua_pcall` | `CallByParam Protect` | 已覆盖 | 通用封装保持 protected call。 |
| 32 | Lua 调用封装 | 栈清理 | `lua_pop` | `defer ls.Pop(nRet)` | 已覆盖 | 确保成功和失败路径都清理栈。 |
| 33 | Lua 调用封装 | 函数存在检查 | `lua_getglobal` 后调用 | 暂无启动自检 | 未覆盖 | 启动期检查关键函数，避免运行时才失败。 |
| 34 | Lua 调用封装 | 调用指标 | 无统一 | 暂无 | 未覆盖 | 记录调用次数、失败次数、耗时和慢调用。 |
| 35 | Lua 调用封装 | 调用上下文 | 分散日志 | 暂无 | 未覆盖 | 携带模块名、玩家 index、怪物 index、事件名等上下文。 |
| 36 | Lua 绑定与导出 | Go 函数注册 | `MULua::Register` | 暂无通用接口 | 未覆盖 | 将 Go 函数注册给 Lua，用于 Bag/Quest/副本脚本。 |
| 37 | Lua 绑定与导出 | 数据表注册 | `RegisterData` | 暂无 | 未覆盖 | 暴露 item/player/map 等只读常量或数据表。 |
| 38 | Lua 绑定与导出 | LuaExport 管理器 | `LuaExport` | 暂无 | 未覆盖 | 建立一组基础导出函数，如随机数、Bag 注册、物品变量。 |
| 39 | Lua 绑定与导出 | GetRandomValue | `LuaGetRandomValue` | 暂无 | 未覆盖 | Lua 侧可请求统一随机数。 |
| 40 | Lua 绑定与导出 | GetLargeRandomValue | `LuaGetLargeRandomValue` | 暂无 | 未覆盖 | 支持大范围随机数。 |
| 41 | Lua 绑定与导出 | AddItemBag | `LuaAddItemBag` | 暂无 | 未覆盖 | Lua 的 `LoadItemBag` 调用 Go 注册 Bag，业务归掉落系统。 |
| 42 | Lua 绑定与导出 | GetBagItemLevel | `LuaGetBagItemLevel` | 暂无 | 未覆盖 | Bag Lua 查询随机 item level。 |
| 43 | Lua 绑定与导出 | GetAncientOpt | `LuaGetSetItemOption` | 暂无 | 未覆盖 | Bag Lua 生成套装/远古选项。 |
| 44 | Lua 绑定与导出 | GetExcellentOptByKind | `LuaGetExcOptionByConfig` | 暂无 | 未覆盖 | Bag Lua 生成卓越选项。 |
| 45 | Lua 绑定与导出 | Lua userdata 策略 | Luna 模板绑定 | 暂无 | 未覆盖 | Go 侧优先用 table/API，不直接模拟 C++ 类绑定。 |
| 46 | Lua 绑定与导出 | 绑定函数安全校验 | 分散 | 暂无 | 未覆盖 | 校验参数数量、类型、对象存在和权限。 |
| 47 | 公式脚本接入 | CalcCharacter.lua 加载 | `ObjCalCharacter::Init` | `formula.load("Character/CalcCharacter.lua")` | 已覆盖 | 公式业务归 `05-formula.md`，通用加载和调用归脚本系统。 |
| 48 | 公式脚本接入 | RegularSkillCalc.lua 加载 | `ObjUseSkill::m_Lua.DoFile` | 已加载 | 部分覆盖 | 技能公式调用应复用通用脚本调用器。 |
| 49 | 公式脚本接入 | ItemCalc.lua 加载 | `configread.m_ItemCalcLua` | 已加载 | 部分覆盖 | 道具/翅膀/宠物公式统一通过脚本系统加载。 |
| 50 | 公式脚本接入 | ExpCalc.lua 加载 | `MasterSkillSystem`、`GameMain` | 已加载 | 部分覆盖 | 经验表和经验公式调用复用脚本系统。 |
| 51 | 公式脚本接入 | StatSpec.lua 加载 | `StatSpecialize` | 已加载 | 部分覆盖 | StatSpec 业务归公式系统，脚本 VM 归本模块。 |
| 52 | 公式脚本接入 | SkillSpec.lua 加载 | `SkillSpecialize` | 暂无 | 未覆盖 | 技能专精脚本加载基础设施。 |
| 53 | 公式脚本接入 | MasterSkillPoint.lua | `MasterLevelSkillTreeSystem` | 暂无 | 未覆盖 | 大师技能点脚本加载基础设施。 |
| 54 | 公式脚本接入 | MasterSkillCalc.lua | `MasterLevelSkillTreeSystem` | 暂无 | 未覆盖 | 大师技能数值脚本加载基础设施。 |
| 55 | Bag/LuaBag 接入 | LuaBag 全局对象 | `gLuaBag` | 暂无 | 未覆盖 | 提供 Bag 脚本执行适配器，业务对象仍在掉落系统。 |
| 56 | Bag/LuaBag 接入 | ItemBagScript.lua 加载 | `CLuaBag::Init` | 暂无 | 未覆盖 | 加载 `Scripts/ItemBags/ItemBagScript.lua`。 |
| 57 | Bag/LuaBag 接入 | LoadItemBag 调用 | `CLuaBag::LoadItemBag` | 暂无 | 未覆盖 | 调用 Lua `LoadItemBag` 注册全部 Bag。 |
| 58 | Bag/LuaBag 接入 | Bag item 变量映射 | `m_ItemInfo`、`GetVariableItem` | 暂无 | 未覆盖 | Lua 可读写当前 Bag item 的类型、等级、选项等变量。 |
| 59 | Bag/LuaBag 接入 | CommonBagItemDrop | `DropCommonBag` | 暂无 | 未覆盖 | 脚本系统执行 Lua，掉落系统处理最终地图落物。 |
| 60 | Bag/LuaBag 接入 | MonsterBagItemDrop | `DropMonsterBag` | 暂无 | 未覆盖 | 支持怪物上下文参数。 |
| 61 | Bag/LuaBag 接入 | EventBagItemDrop | `DropEventBag` | 暂无 | 未覆盖 | 支持事件上下文参数。 |
| 62 | Bag/LuaBag 接入 | EventBagMakeItem | `MakeItemFromBag` | 暂无 | 未覆盖 | Lua 生成物品字段，掉落/道具系统构造最终物品。 |
| 63 | Bag/LuaBag 接入 | GremoryCase item 生成 | `MakeItemFromBagForGremoryCase` | 暂无 | 未覆盖 | 只提供脚本执行和字段输出，Gremory 业务后续独立。 |
| 64 | Bag/LuaBag 接入 | BagManager 初始化链 | `CBagManager::InitBagManager` | 暂无 | 未覆盖 | 脚本系统提供 LuaBag 初始化，掉落系统拥有 BagManager 业务。 |
| 65 | QuestExp 脚本接入 | Quest Lua 全局 VM | `g_MuLuaQuestExp` | 暂无 | 未覆盖 | 建立 QuestExp 脚本 VM 或由任务系统持有 VM。 |
| 66 | QuestExp 脚本接入 | QuestExpLuaBind 注册 | `Luna<QuestExpLuaBind>::Register` | 暂无 | 未覆盖 | 注册任务脚本可调用的 Go 函数。 |
| 67 | QuestExp 脚本接入 | Quest_Info.lua 加载 | `GameMain` | 暂无 | 未覆盖 | 加载任务定义脚本。 |
| 68 | QuestExp 脚本接入 | Quest_Main.lua 加载 | `GameMain` | 暂无 | 未覆盖 | 加载任务主逻辑脚本。 |
| 69 | QuestExp 脚本接入 | CGReqQuestSwitch 调用 | `QuestExpProtocol` | 暂无 | 未覆盖 | 协议入口调用 Lua 任务选择处理。 |
| 70 | QuestExp 脚本接入 | CGReqQuestProgress 调用 | `QuestExpProtocol` | 暂无 | 未覆盖 | 协议入口调用 Lua 任务进度分支。 |
| 71 | QuestExp 脚本接入 | CGReqQuestComplete 调用 | `QuestExpProtocol` | 暂无 | 未覆盖 | 协议入口调用 Lua 完成任务。 |
| 72 | QuestExp 脚本接入 | ItemAndEvent 调用 | `QuestExpProtocol` | 暂无 | 未覆盖 | 请求活动/任务物品列表时调用 Lua。 |
| 73 | QuestExp 脚本接入 | NpcTalkClick 调用 | `QuestExpProtocol` | 暂无 | 未覆盖 | NPC 对话触发 Lua 任务逻辑。 |
| 74 | QuestExp 脚本接入 | GetNpcIndex 绑定 | `QuestExpLuaBind::GetNpcIndex` | 暂无 | 未覆盖 | Lua 查询当前 NPC。 |
| 75 | QuestExp 脚本接入 | GetCharClass 绑定 | `GetCharClass` | 暂无 | 未覆盖 | Lua 查询职业。 |
| 76 | QuestExp 脚本接入 | IsMasterLevel 绑定 | `IsMasterLevel` | 暂无 | 未覆盖 | Lua 判断大师等级状态。 |
| 77 | QuestExp 脚本接入 | GetGensInfluence 绑定 | `GetGensInfluence` | 暂无 | 未覆盖 | Lua 查询 Gens 阵营。 |
| 78 | QuestExp 脚本接入 | GetUserLv 绑定 | `GetUserLv` | 暂无 | 未覆盖 | Lua 查询等级。 |
| 79 | QuestExp 脚本接入 | GetInvenItemFind 绑定 | `GetInvenItemFind` | 暂无 | 未覆盖 | Lua 查询背包任务物品。 |
| 80 | QuestExp 脚本接入 | SetQuestSwitch 绑定 | `SetQuestSwitch` | 暂无 | 未覆盖 | Lua 设置任务开关。 |
| 81 | QuestExp 脚本接入 | SetQuestMonsterKill 绑定 | `SetQuestMonsterKill` | 暂无 | 未覆盖 | Lua 注册击杀目标。 |
| 82 | QuestExp 脚本接入 | DeleteInvenItem 绑定 | `DeleteInvenItem` | 暂无 | 未覆盖 | Lua 消耗背包物品。 |
| 83 | QuestExp 脚本接入 | SetQuestGetItem 绑定 | `SetQuestGetItem` | 暂无 | 未覆盖 | Lua 注册收集物品目标。 |
| 84 | QuestExp 脚本接入 | SetQuestSkillLearn 绑定 | `SetQuestSkillLearn` | 暂无 | 未覆盖 | Lua 注册学习技能目标。 |
| 85 | QuestExp 脚本接入 | SetQuestLevelUp 绑定 | `SetQuestLevelUp` | 暂无 | 未覆盖 | Lua 注册升级目标。 |
| 86 | QuestExp 脚本接入 | SetQuestBuff 绑定 | `SetQuestBuff` | 暂无 | 未覆盖 | Lua 注册 Buff 条件。 |
| 87 | QuestExp 脚本接入 | SetQuestRewardExp 绑定 | `SetQuestRewardExp` | 暂无 | 未覆盖 | Lua 设置经验奖励。 |
| 88 | QuestExp 脚本接入 | SetQuestRewardZen 绑定 | `SetQuestRewardZen` | 暂无 | 未覆盖 | Lua 设置金币奖励。 |
| 89 | QuestExp 脚本接入 | SetQuestRewardItem 绑定 | `SetQuestRewardItem` | 暂无 | 未覆盖 | Lua 设置道具奖励。 |
| 90 | QuestExp 脚本接入 | SendQuestProgress 绑定 | `SendQuestProgress` | 暂无 | 未覆盖 | Lua 触发进度协议下发。 |
| 91 | QuestExp 脚本接入 | SendQuestReward 绑定 | `SendQuestReward` | 暂无 | 未覆盖 | Lua 触发奖励发放。 |
| 92 | QuestExp 脚本接入 | SetQuestNPCTalk 绑定 | `SetQuestNPCTalk` | 暂无 | 未覆盖 | Lua 注册 NPC 对话条件。 |
| 93 | DoppelGanger 脚本接入 | DG Lua VM | `CDoppelGanger::m_MULua` | 暂无 | 未覆盖 | 副本系统持有 VM，脚本系统提供运行时封装。 |
| 94 | DoppelGanger 脚本接入 | CDoppelgangerLua 注册 | `Luna<CDoppelgangerLua>::Register` | 暂无 | 未覆盖 | 注册 DG Lua 可调用的副本函数。 |
| 95 | DoppelGanger 脚本接入 | DoppelGanger.lua 加载 | `m_MULua.DoFile` | 暂无 | 未覆盖 | 加载 `Scripts/Events/DoppelGanger.lua`。 |
| 96 | DoppelGanger 脚本接入 | DG 初始化回调 | `FN_LuaDopplegangerInit` | 暂无 | 未覆盖 | 副本初始化时调用 Lua 初始化。 |
| 97 | DoppelGanger 脚本接入 | DG 每秒回调 | `FN_LuaDoppelgangerCallback` | 暂无 | 未覆盖 | PLAY 状态每秒调用 Lua 驱动部分行为。 |
| 98 | DoppelGanger 脚本接入 | LuaGetMaxUserLevel | `CDoppelgangerLua` | 暂无 | 未覆盖 | Lua 查询参与者最高等级。 |
| 99 | DoppelGanger 脚本接入 | LuaGetMinUserLevel | `CDoppelgangerLua` | 暂无 | 未覆盖 | Lua 查询参与者最低等级。 |
| 100 | DoppelGanger 脚本接入 | LuaGetUserCount | `CDoppelgangerLua` | 暂无 | 未覆盖 | Lua 查询参与人数。 |
| 101 | DoppelGanger 脚本接入 | LuaSetReadyTime | `CDoppelgangerLua` | 暂无 | 未覆盖 | Lua 设置准备时间。 |
| 102 | DoppelGanger 脚本接入 | LuaSetPlayTime | `CDoppelgangerLua` | 暂无 | 未覆盖 | Lua 设置进行时间。 |
| 103 | DoppelGanger 脚本接入 | LuaSetEndTime | `CDoppelgangerLua` | 暂无 | 未覆盖 | Lua 设置结束时间。 |
| 104 | DoppelGanger 脚本接入 | LuaSetHerdStartPosInfo | `CDoppelgangerLua` | 暂无 | 未覆盖 | Lua 设置怪群起点。 |
| 105 | DoppelGanger 脚本接入 | LuaSetHerdEndPosInfo | `CDoppelgangerLua` | 暂无 | 未覆盖 | Lua 设置怪群终点。 |
| 106 | DoppelGanger 脚本接入 | LuaAddMonsterHerd | `CDoppelgangerLua` | 暂无 | 未覆盖 | Lua 添加怪群。 |
| 107 | DoppelGanger 脚本接入 | LuaAddMonsterNormal | `CDoppelgangerLua` | 暂无 | 未覆盖 | Lua 添加普通怪。 |
| 108 | DoppelGanger 脚本接入 | LuaMonsterHerdStart | `CDoppelgangerLua` | 暂无 | 未覆盖 | Lua 启动怪群。 |
| 109 | DoppelGanger 脚本接入 | Lua boss state | `LuaGet/SetBossMonsterState` | 暂无 | 未覆盖 | Lua 读写 Boss 状态。 |
| 110 | WZ/SMD 文本脚本 | ReadScript token parser | `ReadScript.h::GetToken` | 暂无 | 未覆盖 | 实现兼容 GameServer 的老式 token parser。 |
| 111 | WZ/SMD 文本脚本 | WzMemScript parser | `CWzMemScript` | 暂无 | 未覆盖 | 支持内存 buffer 脚本解析。 |
| 112 | WZ/SMD 文本脚本 | WZScriptEncode parser | `CWZScriptEncode` | 暂无 | 未覆盖 | 如需兼容编码脚本，再实现编码/解码解析器。 |
| 113 | WZ/SMD 文本脚本 | NAME token | `NAME` | 暂无 | 未覆盖 | 解析标识符和双引号字符串。 |
| 114 | WZ/SMD 文本脚本 | NUMBER token | `NUMBER` | 暂无 | 未覆盖 | 解析整数、浮点数和负数。 |
| 115 | WZ/SMD 文本脚本 | COMMAND token | `#` | 暂无 | 未覆盖 | 解析命令符。 |
| 116 | WZ/SMD 文本脚本 | LP/RP token | `{`、`}` | 暂无 | 未覆盖 | 解析块结构。 |
| 117 | WZ/SMD 文本脚本 | COMMA/SEMICOLON token | `,`、`;` | 暂无 | 未覆盖 | 解析分隔符。 |
| 118 | WZ/SMD 文本脚本 | 注释处理 | `//` | 暂无 | 未覆盖 | 跳过行注释。 |
| 119 | WZ/SMD 文本脚本 | 空白处理 | `isspace` | 暂无 | 未覆盖 | 跳过空白、换行和制表。 |
| 120 | WZ/SMD 文本脚本 | 错误 token | `SMD_ERROR` | 暂无 | 未覆盖 | 非法字符返回错误并携带位置。 |
| 121 | WZ/SMD 文本脚本 | GetNumber/GetString | `GetNumber/GetString` | 暂无 | 未覆盖 | 提供读取当前 token 值的接口。 |
| 122 | WZ/SMD 文本脚本 | 文件 parser | `ReadScript` 全局 FILE | 暂无 | 未覆盖 | Go 实现不使用全局变量，封装为 parser 实例。 |
| 123 | WZ/SMD 文本脚本 | 内存 parser | `SetBuffer` | 暂无 | 未覆盖 | 测试和脚本解码后可直接从内存解析。 |
| 124 | WZ/SMD 文本脚本 | 兼容性测试样本 | 多个 `GetToken` 使用者 | 暂无 | 未覆盖 | 用 GameServer 配置样本验证 token 行为。 |
| 125 | Buff/配置加载边界 | BuffScriptLoader 加载 | `g_BuffScript.Load` | 暂无 | 未覆盖 | 脚本系统提供通用加载框架，Buff 定义业务归 `13-buffs.md`。 |
| 126 | Buff/配置加载边界 | Buff data 查询 | `GetBuffData` | 暂无 | 未覆盖 | Buff 系统拥有查询接口，脚本系统只记录加载能力边界。 |
| 127 | Buff/配置加载边界 | PeriodBuff 加载 | `AddPeriodBuffEffectInfo` | 暂无 | 未覆盖 | 期限 Buff 业务归 Buff 系统。 |
| 128 | 业务配置加载边界 | LoadScript 命名统一 | 多模块 `LoadScript` | 分散 XML/INI 读取 | 部分覆盖 | 不是所有 LoadScript 都是 Lua，脚本系统只抽通用解析/加载能力。 |
| 129 | 业务配置加载边界 | XML 配置不强行纳入脚本 | pugi/xml、conf.XML | `conf.XML` | 已覆盖 | XML 业务配置可留在业务模块，不必都迁入脚本系统。 |
| 130 | 业务配置加载边界 | INI 配置不强行纳入脚本 | `CIniReader` | conf ini | 已覆盖 | INI 读取归配置系统，脚本系统只处理脚本/DSL 解析。 |
| 131 | 安全与沙箱 | 禁止危险库 | 无统一 | 暂无 | 未覆盖 | 线上 Lua 默认禁用 OS/IO/debug/package 等危险能力。 |
| 132 | 安全与沙箱 | 调用超时 | 无统一 | 暂无 | 未覆盖 | 防止 Lua 死循环阻塞游戏主循环。 |
| 133 | 安全与沙箱 | 指令预算 | 无统一 | 暂无 | 未覆盖 | gopher-lua 可通过 hook/上下文策略限制执行。 |
| 134 | 安全与沙箱 | 内存限制 | 无统一 | 暂无 | 未覆盖 | 控制脚本表膨胀和内存泄漏风险。 |
| 135 | 安全与沙箱 | 脚本权限模型 | 无统一 | 暂无 | 未覆盖 | 区分公式脚本、Bag 脚本、Quest 脚本、GM 脚本可访问 API。 |
| 136 | 安全与沙箱 | 只读数据导出 | 无统一 | 暂无 | 未覆盖 | 常量表、物品表、地图表默认只读。 |
| 137 | 热重载与运维 | 开发期热重载 | BagManager 可重 Init | 暂无 | 未覆盖 | 支持指定模块重载脚本。 |
| 138 | 热重载与运维 | 热重载事务 | 无统一 | 暂无 | 未覆盖 | 新脚本加载和自检成功后再替换旧 VM。 |
| 139 | 热重载与运维 | 热重载回滚 | 无统一 | 暂无 | 未覆盖 | 加载失败时保留旧脚本。 |
| 140 | 热重载与运维 | 脚本版本信息 | 无统一 | 暂无 | 未覆盖 | 记录脚本路径、mtime、hash、加载时间。 |
| 141 | 热重载与运维 | 脚本诊断命令 | 无统一 | 暂无 | 未覆盖 | 输出已加载脚本、函数检查、调用统计。 |
| 142 | 跨系统接口 | 公式系统接口 | `ObjCalCharacter` 等 | `05-formula.md` | 部分覆盖 | 公式系统调用脚本系统，不把公式业务迁入脚本系统。 |
| 143 | 跨系统接口 | 掉落系统接口 | `LuaBag/BagManager` | `14-drops.md` | 未覆盖 | 掉落系统拥有 Bag 业务，脚本系统提供 LuaBag 执行。 |
| 144 | 跨系统接口 | 任务系统接口 | `QuestExpLuaBind` | `11-quests.md` | 未覆盖 | 任务系统拥有 QuestExp 业务，脚本系统提供绑定和调用。 |
| 145 | 跨系统接口 | 副本系统接口 | `DoppelgangerLua` | `21-dungeons.md` | 未覆盖 | 副本系统拥有 DG 业务，脚本系统提供 Lua 回调执行。 |
| 146 | 跨系统接口 | Buff 系统接口 | `BuffScriptLoader` | `13-buffs.md` | 未覆盖 | Buff 系统拥有 Buff 定义，脚本系统提供通用加载边界。 |
| 147 | 跨系统接口 | 活动系统接口 | Event `LoadScript`、EventBag | `22-events.md` | 未覆盖 | 活动系统可通过脚本系统加载活动脚本或调用奖励脚本。 |
| 148 | 跨系统接口 | 怪物 AI 接口 | AI `LoadData/LoadScript` | `25-monster-ai.md` | 未覆盖 | AI 配置解析可复用 WZ/SMD parser，行为归怪物 AI。 |
| 149 | 跨系统接口 | 宠物召唤接口 | Pet/Muun script configs | `24-pets-summons.md` | 未覆盖 | 宠物业务归宠物系统，脚本系统只提供加载/调用能力。 |
| 150 | 测试与验收 | Lua VM 生命周期测试 | `MULua` | 待实现 | 未覆盖 | 覆盖创建、加载、关闭、重复关闭。 |
| 151 | 测试与验收 | Lua Call 签名测试 | `Generic_Call` | 待实现 | 未覆盖 | 覆盖 int、float、int64、bool、string、返回值类型错误。 |
| 152 | 测试与验收 | Lua 错误测试 | `lua_pcall` | 待实现 | 未覆盖 | 覆盖函数不存在、脚本 panic、参数错误。 |
| 153 | 测试与验收 | 并发访问测试 | `UseSync` | 待实现 | 未覆盖 | 覆盖同一 VM 并发调用必须串行。 |
| 154 | 测试与验收 | 公式脚本回归测试 | CalcCharacter/Skill/Exp | 待实现 | 未覆盖 | 确认抽通用脚本系统后公式结果不变。 |
| 155 | 测试与验收 | LuaBag 加载测试 | `LoadItemBag` | 待实现 | 未覆盖 | 覆盖 Lua 注册 Bag 和基础 item 变量。 |
| 156 | 测试与验收 | Bag 掉落脚本测试 | Common/Monster/Event drop | 待实现 | 未覆盖 | 覆盖脚本生成物品字段并交给掉落系统。 |
| 157 | 测试与验收 | Quest Lua 绑定测试 | `QuestExpLuaBind` | 待实现 | 未覆盖 | 覆盖任务条件、奖励、NPC talk 关键绑定。 |
| 158 | 测试与验收 | DoppelGanger 回调测试 | init/callback | 待实现 | 未覆盖 | 覆盖初始化回调和每秒回调。 |
| 159 | 测试与验收 | SMD parser token 测试 | `ReadScript/GetToken` | 待实现 | 未覆盖 | 覆盖 NAME、NUMBER、注释、括号、错误 token。 |
| 160 | 测试与验收 | 热重载测试 | reload/rollback | 待实现 | 未覆盖 | 覆盖成功替换、失败回滚和并发调用保护。 |
