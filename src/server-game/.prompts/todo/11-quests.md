# 11. 任务系统

本模块覆盖 `GameServer` 中老任务 `QuestInfo`、扩展任务 `Quests`、新任务 `QuestExp` 三套任务逻辑，并映射到 `server-game` 当前协议入口、角色状态、道具、NPC、怪物击杀、奖励与持久化的实现缺口。当前 `server-game` 主要只有任务相关 opcode 入口，尚未形成独立任务业务系统。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 老任务 QuestInfo | QuestInfo 初始化 | `CQuestInfo::Init` | 暂无独立 quest 初始化模块 | 未覆盖 | 建立老任务管理器初始化流程，明确任务表、任务物品表、玩家任务状态的装载时机。 |
| 2 | 老任务 QuestInfo | LoadQuestInfo 配置加载 | `CQuestInfo::LoadQuestInfo` | 暂无老任务配置加载 | 未覆盖 | 解析老任务配置并转换为 Go 结构，保留任务编号、条件、奖励、NPC、怪物和物品字段。 |
| 3 | 老任务 QuestInfo | InitQuestItem | `CQuestInfo::InitQuestItem` | `game/item` 仅有 `KindBQuest` 类型 | 未覆盖 | 建立任务物品定义索引，用于掉落、拾取、出售限制和任务完成判断。 |
| 4 | 老任务 QuestInfo | IsQuest | `CQuestInfo::IsQuest` | 暂无任务 ID 判定 | 未覆盖 | 提供任务定义存在性检查，避免协议请求引用不存在的任务。 |
| 5 | 老任务 QuestInfo | GetQuestState | `CQuestInfo::GetQuestState` | 玩家模型未保存老任务状态数组 | 未覆盖 | 为玩家提供老任务状态读取能力，兼容原 C++ 的任务状态语义。 |
| 6 | 老任务 QuestInfo | GetQuestStateBYTE | `CQuestInfo::GetQuestStateBYTE` | 暂无任务状态字节编码 | 未覆盖 | 明确任务状态在协议和 DB 中的字节表达，避免 Go 侧状态与客户端不一致。 |
| 7 | 老任务 QuestInfo | SetQuestState | `CQuestInfo::SetQuestState` | 暂无任务状态写入 | 未覆盖 | 实现任务状态变更入口，并处理接受、完成、重置等状态转移约束。 |
| 8 | 老任务 QuestInfo | ReSetQuestState | `CQuestInfo::ReSetQuestState` | 暂无老任务重置 | 未覆盖 | 支持任务状态清空或回退，用于任务失败、重置、调试和异常修复。 |
| 9 | 老任务 QuestInfo | QuestAccept | `CQuestInfo::QuestAccept` | `handle/c1c2.go` 仅有 `0xA0` 入口 | 未覆盖 | 实现老任务接受流程，包含条件检查、状态写入和客户端响应。 |
| 10 | 老任务 QuestInfo | QuestClear | `CQuestInfo::QuestClear` | 暂无老任务完成流程 | 未覆盖 | 实现老任务完成流程，包含条件复核、奖励发放、状态保存和响应。 |
| 11 | 老任务条件/NPC/掉落 | QuestInfoSave | `CQuestInfo::QuestInfoSave` | `game/model/db.go` 无任务持久化字段 | 未覆盖 | 设计老任务状态持久化字段或表，保证登录后恢复任务进度。 |
| 12 | 老任务条件/NPC/掉落 | QuestRunConditionCheck | `CQuestInfo::QuestRunConditionCheck` | 暂无老任务运行条件检查 | 未覆盖 | 实现接任务条件检查，包括等级、职业、状态、物品、地图等限制。 |
| 13 | 老任务条件/NPC/掉落 | QuestClearConditionCheck | `CQuestInfo::QuestClearConditionCheck` | 暂无老任务完成条件检查 | 未覆盖 | 实现完成条件检查，避免客户端伪造完成请求。 |
| 14 | 老任务条件/NPC/掉落 | CompareCondition | `CQuestInfo::CompareCondition` | 暂无通用任务条件比较器 | 未覆盖 | 抽象条件比较逻辑，供老任务与后续 QuestExp 复用。 |
| 15 | 老任务条件/NPC/掉落 | NpcTalk | `CQuestInfo::NpcTalk` | NPC 对话未接入任务流程 | 未覆盖 | 将 NPC 点击/对话事件接入任务接受、完成、提示和状态查询。 |
| 16 | 老任务条件/NPC/掉落 | MonsterItemDrop | `CQuestInfo::MonsterItemDrop` | 怪物掉落未接入任务物品规则 | 未覆盖 | 根据玩家任务状态决定怪物是否掉落任务物品。 |
| 17 | 老任务条件/NPC/掉落 | MonsterItemDropParty | `CQuestInfo::MonsterItemDropParty` | 组队任务掉落未实现 | 未覆盖 | 处理队伍成员任务状态、地图距离、归属和掉落资格。 |
| 18 | 老任务条件/NPC/掉落 | AddMonsterKillCount | `CQuestInfo::AddMonsterKillCount` | 怪物死亡未接入任务击杀进度 | 未覆盖 | 在怪物死亡结算处增加任务击杀计数钩子。 |
| 19 | 老任务条件/NPC/掉落 | GetQuestKillCount | `CQuestInfo::GetQuestKillCount` | 暂无击杀进度读取 | 未覆盖 | 提供任务击杀进度查询，用于完成检查和进度下发。 |
| 20 | 老任务条件/NPC/掉落 | SendQuestMonsterKill | `CQuestInfo::SendQuestMonsterKill` | `ThirdQuestMonsterCountMsg` 配置未消费 | 未覆盖 | 按配置控制任务击杀进度通知，避免刷屏或客户端进度不同步。 |
| 21 | 扩展任务 Quests 定义 | CQuests 配置加载 | `CQuests::LoadData` | 暂无扩展任务配置加载 | 未覆盖 | 加载扩展任务定义，区分任务数组、阶段、目标、奖励和 NPC 入口。 |
| 22 | 扩展任务 Quests 定义 | QUEST_INTERNAL_DATA 模型 | `QUEST_INTERNAL_DATA` | 暂无对应 Go 模型 | 未覆盖 | 定义扩展任务内部结构，表达任务索引、状态、目标、奖励和条件。 |
| 23 | 扩展任务 Quests 定义 | CGReqQuestList | `CQuests::CGReqQuestList` | `handle/c1c2.go` 有任务 opcode 入口但缺业务 | 未覆盖 | 实现任务列表请求，按角色状态返回可接、进行中、已完成任务。 |
| 24 | 扩展任务 Quests 定义 | CGReqQuestListLogin | `CQuests::CGReqQuestListLogin` | 登录流程未加载任务列表 | 未覆盖 | 登录后初始化客户端任务状态，避免客户端任务 UI 空白或旧状态残留。 |
| 25 | 扩展任务 Quests 定义 | GCInitQuest | `CQuests::GCInitQuest` | 暂无任务初始化包体 | 未覆盖 | 定义并发送任务初始化响应，兼容客户端预期字段。 |
| 26 | 扩展任务 Quests 定义 | CGReqQuestSpecificInfo | `CQuests::CGReqQuestSpecificInfo` | 暂无指定任务详情请求 | 未覆盖 | 支持客户端请求单个任务详情，返回目标、奖励、NPC 和状态。 |
| 27 | 扩展任务 Quests 定义 | GCAnsQuestSpecificInfo | `CQuests::GCAnsQuestSpecificInfo` | 暂无任务详情响应 | 未覆盖 | 实现任务详情响应包体，保证 UI 能显示目标与奖励信息。 |
| 28 | 扩展任务 Quests 定义 | FindQuestArrayNum(Index, State) | `CQuests::FindQuestArrayNum` | 暂无任务索引查询 | 未覆盖 | 按任务 ID 与状态查找配置项，减少重复扫描和状态歧义。 |
| 29 | 扩展任务 Quests 定义 | ConvertByteToData | `CQuests::ConvertByteToData` | 暂无任务状态反序列化 | 未覆盖 | 将 DB/协议字节状态转换为 Go 内部任务数据。 |
| 30 | 扩展任务 Quests 定义 | ConvertDataToByte | `CQuests::ConvertDataToByte` | 暂无任务状态序列化 | 未覆盖 | 将 Go 内部任务状态转换为 DB/协议字节格式。 |
| 31 | 扩展任务状态与流程 | AddQuest | `CQuests::AddQuest` | 暂无扩展任务接受 | 未覆盖 | 增加任务到玩家任务列表，处理重复接取、上限和前置条件。 |
| 32 | 扩展任务状态与流程 | DelQuest | `CQuests::DelQuest` | 暂无扩展任务删除 | 未覆盖 | 从玩家任务列表删除任务，处理放弃、完成清理和异常状态。 |
| 33 | 扩展任务状态与流程 | CGReqAnswer | `CQuests::CGReqAnswer` | 暂无任务确认答案处理 | 未覆盖 | 处理客户端任务接受/拒绝/确认类请求，校验 NPC 和任务状态。 |
| 34 | 扩展任务状态与流程 | CGReqDialog | `CQuests::CGReqDialog` | NPC 对话未接入扩展任务 | 未覆盖 | 根据 NPC、角色状态和任务配置返回对应任务对话。 |
| 35 | 扩展任务状态与流程 | CGReqQuestFinish | `CQuests::CGReqQuestFinish` | 暂无扩展任务完成请求 | 未覆盖 | 实现扩展任务完成请求，包含条件校验、奖励、扣物品和状态更新。 |
| 36 | 扩展任务状态与流程 | CGReqQuestDelete | `CQuests::CGReqQuestDelete` | 暂无扩展任务放弃请求 | 未覆盖 | 实现任务放弃协议，清理临时进度并返回客户端结果。 |
| 37 | 扩展任务状态与流程 | CheckRequirements | `CQuests::CheckRequirements` | 暂无扩展任务条件聚合检查 | 未覆盖 | 聚合等级、职业、物品、任务状态、NPC、事件等接取/完成条件。 |
| 38 | 扩展任务状态与流程 | CheckNPC | `CQuests::CheckNPC` | 暂无任务 NPC 校验 | 未覆盖 | 防止玩家在错误 NPC 或距离外触发任务操作。 |
| 39 | 扩展任务状态与流程 | CheckStage | `CQuests::CheckStage` | 暂无任务阶段检查 | 未覆盖 | 校验任务阶段是否满足指定操作，避免跳阶段完成。 |
| 40 | 扩展任务状态与流程 | SetQuestStage | `CQuests::SetQuestStage` | 暂无任务阶段写入 | 未覆盖 | 写入任务阶段并触发必要的进度通知与持久化。 |
| 41 | 扩展任务目标/事件 | KillMonsterProc | `CQuests::KillMonsterProc` | 怪物击杀未接入扩展任务 | 未覆盖 | 在怪物死亡事件中更新扩展任务击杀目标。 |
| 42 | 扩展任务目标/事件 | KillGensProc | `CQuests::KillGensProc` | Gens/PVP 相关任务未实现 | 未覆盖 | 为 Gens 击杀类任务预留事件入口和进度更新。 |
| 43 | 扩展任务目标/事件 | FinishEventProc | `CQuests::FinishEventProc` | 事件完成未接入任务 | 未覆盖 | 将事件地图完成结果转换为任务目标进度。 |
| 44 | 扩展任务目标/事件 | DestroyBloodCastleGate | `CQuests::DestroyBloodCastleGate` | 血色城堡事件任务未实现 | 未覆盖 | 接入血色城堡城门破坏类任务目标。 |
| 45 | 扩展任务目标/事件 | ChaosCastleKillProc | `CQuests::ChaosCastleKillProc` | 赤色要塞击杀任务未实现 | 未覆盖 | 接入赤色要塞击杀积分或击杀计数任务。 |
| 46 | 扩展任务目标/事件 | AddDevilSquarePoint | `CQuests::AddDevilSquarePoint` | 恶魔广场积分任务未实现 | 未覆盖 | 接入恶魔广场积分变化并更新任务目标。 |
| 47 | 扩展任务目标/事件 | DropQuestItem | `CQuests::DropQuestItem` | 任务物品掉落未统一 | 未覆盖 | 为扩展任务生成任务物品掉落，处理地图、怪物和玩家资格。 |
| 48 | 扩展任务目标/事件 | DropQuestItemParty | `CQuests::DropQuestItemParty` | 组队扩展任务掉落未实现 | 未覆盖 | 实现组队场景下任务物品掉落资格和归属。 |
| 49 | 扩展任务目标/事件 | GetItemCount | `CQuests::GetItemCount` | 背包任务物品计数未封装 | 未覆盖 | 统计玩家背包中指定任务物品数量，用于完成条件。 |
| 50 | 扩展任务目标/事件 | GetRequiredCount | `CQuests::GetRequiredCount` | 暂无任务目标需求数量查询 | 未覆盖 | 查询目标所需数量，供 UI 进度和完成检查使用。 |
| 51 | 扩展任务奖励/日常 | GiveReward | `CQuests::GiveReward` | 奖励系统未接入任务 | 未覆盖 | 聚合发放经验、金币、道具、属性点和特殊奖励。 |
| 52 | 扩展任务奖励/日常 | CreateItem | `CQuests::CreateItem` | 任务奖励道具生成未实现 | 未覆盖 | 生成固定奖励道具，并处理背包空间、属性和失败回滚。 |
| 53 | 扩展任务奖励/日常 | CreateRandomItem | `CQuests::CreateRandomItem` | 随机任务奖励未实现 | 未覆盖 | 实现随机奖励池选择、生成和日志记录。 |
| 54 | 扩展任务奖励/日常 | DeleteItem | `CQuests::DeleteItem` | 任务完成扣物品未实现 | 未覆盖 | 完成任务时按条件扣除任务物品，避免重复提交。 |
| 55 | 扩展任务奖励/日常 | LevelUp | `CQuests::LevelUp` | 升级事件未接入任务 | 未覆盖 | 玩家升级后刷新任务可接状态或触发等级类目标。 |
| 56 | 扩展任务奖励/日常 | CalcMoney | `CQuests::CalcMoney` | 任务金币计算未实现 | 未覆盖 | 计算任务金币奖励或需求，处理上限和溢出。 |
| 57 | 扩展任务奖励/日常 | GCNotifyQuest | `CQuests::GCNotifyQuest` | 暂无任务通知包 | 未覆盖 | 向客户端推送任务状态变化、可完成、可接取等通知。 |
| 58 | 扩展任务奖励/日常 | StartTutorial | `CQuests::StartTutorial` | `tutorialKeyCOmplete` 入口存在但缺业务 | 未覆盖 | 实现新手教程任务开始和完成状态同步。 |
| 59 | 扩展任务奖励/日常 | IsDailyQuest | `CQuests::IsDailyQuest` | 暂无每日任务判断 | 未覆盖 | 区分普通任务和每日任务，供刷新、重置和限制使用。 |
| 60 | 扩展任务奖励/日常 | GetDailyQuest | `CQuests::GetDailyQuest` | 暂无每日任务获取 | 未覆盖 | 根据日期、角色状态和配置获取每日任务实例。 |
| 61 | QuestExp 定义与 Lua 配置 | SetQuestExpAsk | `QuestExpInfo::SetQuestExpAsk` | 暂无 QuestExp ask 定义 | 未覆盖 | 建立 QuestExp 目标条件定义模型，表达击杀、物品、技能、等级、事件等 ask；Lua VM 和绑定基础设施归 `26-script.md`。 |
| 62 | QuestExp 定义与 Lua 配置 | SetQuestReward | `QuestExpInfo::SetQuestReward` | 暂无 QuestExp reward 定义 | 未覆盖 | 建立 QuestExp 奖励定义模型，表达经验、金币、道具、随机奖励和贡献等奖励。 |
| 63 | QuestExp 定义与 Lua 配置 | QuestExpItemInit | `QuestExpManager::QuestExpItemInit` | QuestExp 物品初始化未实现 | 未覆盖 | 初始化 QuestExp 任务物品索引，支持任务掉落和物品属性判断。 |
| 64 | QuestExp 定义与 Lua 配置 | QuestInfoDelete | `QuestExpManager::QuestInfoDelete` | 暂无 QuestExp 配置清理 | 未覆盖 | 支持重载或关闭时清理 QuestExp 配置缓存。 |
| 65 | QuestExp 定义与 Lua 配置 | SetQuestExpInfo Ask | `QuestExpManager::SetQuestExpInfo(QuestExpAsk)` | 暂无 ask 注册入口 | 未覆盖 | 将配置或脚本定义的 ask 注册到 QuestExp 管理器。 |
| 66 | QuestExp 定义与 Lua 配置 | SetQuestExpInfo Reward | `QuestExpManager::SetQuestExpInfo(QuestExpReward)` | 暂无 reward 注册入口 | 未覆盖 | 将配置或脚本定义的 reward 注册到 QuestExp 管理器。 |
| 67 | QuestExp 定义与 Lua 配置 | AddQuestExpInfoList | `QuestExpManager::AddQuestExpInfoList` | 暂无 QuestExp 列表索引 | 未覆盖 | 按 Episode、QuestSwitch、QuestIndex 建立 QuestExp 定义索引。 |
| 68 | QuestExp 定义与 Lua 配置 | AddQuestExpRewardKind | `QuestExpManager::AddQuestExpRewardKind` | 暂无奖励类型索引 | 未覆盖 | 建立奖励类型到处理器的映射，便于扩展奖励种类。 |
| 69 | QuestExp 定义与 Lua 配置 | AddQuestDropItemInfo | `QuestExpManager::AddQuestDropItemInfo` | 暂无任务掉落定义 | 未覆盖 | 注册 QuestExp 掉落规则，供怪物死亡时查询。 |
| 70 | QuestExp 定义与 Lua 配置 | AddQuestItemInfo/IsQuestItemInfo/IsQuestItemAtt/GetQuestItemEp | `QuestExpManager` 任务物品函数组 | `KindBQuest` 存在但无 QuestExp 属性判断 | 未覆盖 | 完成任务物品登记、属性识别、Episode 查询和客户端限制。 |
| 71 | QuestExp 用户状态/持久化 | QuestExpUserInfo Init/Clear | `QuestExpUserInfo::Init/Clear` | 暂无 QuestExp 玩家状态 | 未覆盖 | 建立玩家 QuestExp 运行态，支持初始化、清空和登出释放。 |
| 72 | QuestExp 用户状态/持久化 | SetEpisode/GetEpisode | `QuestExpUserInfo::SetEpisode/GetEpisode` | 暂无 Episode 字段 | 未覆盖 | 保存玩家当前 QuestExp Episode，用于任务线进度判断。 |
| 73 | QuestExp 用户状态/持久化 | SetQuestSwitch/GetQuestSwitch | `QuestExpUserInfo::SetQuestSwitch/GetQuestSwitch` | `questSwitch` opcode 存在但缺状态 | 未覆盖 | 实现 QuestSwitch 状态读写，支撑章节或任务线切换。 |
| 74 | QuestExp 用户状态/持久化 | SetAskCnt/GetAskCnt | `QuestExpUserInfo::SetAskCnt/GetAskCnt` | 暂无 ask 计数 | 未覆盖 | 保存每个任务条件的当前计数，用于进度显示和完成判断。 |
| 75 | QuestExp 用户状态/持久化 | SetStartDate/SetEndDate | `QuestExpUserInfo::SetStartDate/SetEndDate` | 暂无 QuestExp 时间字段 | 未覆盖 | 支持限时任务、每日任务和过期检查。 |
| 76 | QuestExp 用户状态/持久化 | SetQuestProgState/GetQuestProgState | `QuestExpUserInfo::SetQuestProgState/GetQuestProgState` | 暂无 QuestExp 进度状态 | 未覆盖 | 记录任务未接、进行中、可完成、已完成等状态。 |
| 77 | QuestExp 用户状态/持久化 | AddUserQuestAskInfo | `QuestExpUserMng::AddUserQuestAskInfo` | 暂无玩家 ask 初始化 | 未覆盖 | 接取任务时初始化玩家 ask 列表和初始进度。 |
| 78 | QuestExp 用户状态/持久化 | UserQuestInfoSave | `QuestExpUserMng::UserQuestInfoSave` | DB 模型无任务明细 | 未覆盖 | 保存单个任务状态，降低频繁全量保存成本。 |
| 79 | QuestExp 用户状态/持久化 | UserAllQuestInfoSave | `QuestExpUserMng::UserAllQuestInfoSave` | 登出保存未接入任务 | 未覆盖 | 登出或切线时保存玩家全部 QuestExp 状态。 |
| 80 | QuestExp 用户状态/持久化 | DB_ReqUserQuestInfo/UserQuestInfoLoad | `QuestExpUserMng` DB 函数组 | 登录加载未接入 QuestExp | 未覆盖 | 登录时从 DB 读取 QuestExp 状态并恢复到玩家对象。 |
| 81 | QuestExp 进度条件 | IsQuestAccept/IsQuestComplete/IsProgQuestInfo | `QuestExpUserMng` 状态查询函数组 | 暂无 QuestExp 状态查询 | 未覆盖 | 提供是否已接、是否完成、是否进行中的统一查询接口。 |
| 82 | QuestExp 进度条件 | ReqQuestAskComplete | `QuestExpProgMng::ReqQuestAskComplete` | `questComplete` opcode 存在但缺业务 | 未覆盖 | 处理客户端完成任务条件请求，重新校验进度和奖励资格。 |
| 83 | QuestExp 进度条件 | QuestExpGiveUpBtnClick | `QuestExpProgMng::QuestExpGiveUpBtnClick` | `questGiveUp` opcode 存在但缺业务 | 未覆盖 | 实现 QuestExp 放弃任务，清理进度、任务物品和状态。 |
| 84 | QuestExp 进度条件 | SendQuestProgress | `QuestExpProgMng::SendQuestProgress` | `questProgress` opcode 存在但缺响应 | 未覆盖 | 下发单个任务进度，保持客户端 UI 同步。 |
| 85 | QuestExp 进度条件 | SendQuestProgressInfo | `QuestExpProgMng::SendQuestProgressInfo` | `questProgressInfo` opcode 存在但缺响应 | 未覆盖 | 下发任务详细进度，包括 ask 列表、计数和状态。 |
| 86 | QuestExp 进度条件 | ChkQuestAsk | `QuestExpProgMng::ChkQuestAsk` | 暂无 QuestExp ask 总检查 | 未覆盖 | 聚合检查所有 ask 是否满足完成条件。 |
| 87 | QuestExp 进度条件 | ChkUserQuestTypeMonsterKill/Party | `QuestExpProgMng` 击杀检查函数 | 怪物死亡未接入 QuestExp | 未覆盖 | 支持单人和组队击杀类 QuestExp 目标。 |
| 88 | QuestExp 进度条件 | ChkUserQuestTypeItem/DeleteInventoryItem | `QuestExpProgMng` 物品检查函数 | 背包未接入 QuestExp 目标 | 未覆盖 | 支持收集物品、扣除物品和背包物品计数。 |
| 89 | QuestExp 进度条件 | ChkUserQuestTypeSkillLearn/LevelUp/Buff/NeedZen | `QuestExpProgMng` 状态检查函数 | 技能、升级、Buff、Zen 未接入任务 | 未覆盖 | 将角色技能、等级、Buff、金币变化接入 QuestExp 条件。 |
| 90 | QuestExp 进度条件 | ChkUserQuestTypeEventMap | `QuestExpProgMng::ChkUserQuestTypeEventMap` | 事件地图未接入 QuestExp | 未覆盖 | 支持事件地图击杀、清场、积分等任务目标。 |
| 91 | QuestExp 奖励/背包/随机 | SetQuestTimeLimit/CheckExpireDate | `QuestExpProgMng::SetQuestTimeLimit/CheckExpireDate` | 暂无限时任务检查 | 未覆盖 | 设置任务期限并在请求、登录、进度更新时检查过期状态。 |
| 92 | QuestExp 奖励/背包/随机 | InvenChk_EnableReward | `QuestExpUserMng::InvenChk_EnableReward` | 奖励背包检查未实现 | 未覆盖 | 发奖励前检查背包空间，避免奖励丢失或状态已完成但物品未发。 |
| 93 | QuestExp 奖励/背包/随机 | QuestExpCheckInventoryEmptySpace | `QuestExpUserMng::QuestExpCheckInventoryEmptySpace` | 暂无 QuestExp 背包空间检查 | 未覆盖 | 按奖励物品大小和数量判断背包可用空间。 |
| 94 | QuestExp 奖励/背包/随机 | QuestExpOnlyInventoryRectCheck | `QuestExpUserMng::QuestExpOnlyInventoryRectCheck` | 暂无矩形空间校验 | 未覆盖 | 对指定背包区域做矩形占用检查。 |
| 95 | QuestExp 奖励/背包/随机 | QuestExpTempInventoryItemSet | `QuestExpUserMng::QuestExpTempInventoryItemSet` | 暂无奖励临时落位 | 未覆盖 | 在真正发放奖励前计算临时背包落点。 |
| 96 | QuestExp 奖励/背包/随机 | GetRandomRewardItemIndex | `QuestExpUserMng::GetRandomRewardItemIndex` | 暂无 QuestExp 随机奖励选择 | 未覆盖 | 根据随机奖励组选择奖励项，并保证概率与配置一致。 |
| 97 | QuestExp 奖励/背包/随机 | IsRandRewardIndex/IsRandResultReward | `QuestExpUserMng` 随机奖励校验函数 | 暂无随机奖励校验 | 未覆盖 | 校验随机奖励索引和结果是否合法，防止异常配置或伪造。 |
| 98 | QuestExp 奖励/背包/随机 | SendQuestReward | `QuestExpUserMng::SendQuestReward` | 暂无 QuestExp 奖励发放 | 未覆盖 | 发放 QuestExp 奖励并回写状态、通知客户端和记录日志。 |
| 99 | QuestExp 奖励/背包/随机 | SendQuestAskInfoUpdate | `QuestExpProgMng::SendQuestAskInfoUpdate` | 暂无 ask 增量通知 | 未覆盖 | ask 计数变化时下发增量更新，减少全量同步。 |
| 100 | QuestExp 奖励/背包/随机 | GCANSQuestCompleteBtnClick | `QuestExpProgMng::GCANSQuestCompleteBtnClick` | 暂无完成按钮响应 | 未覆盖 | 对客户端完成按钮请求返回成功、失败原因和奖励结果。 |
| 101 | QuestExp 协议/NPC | SetQuestProg | `QuestExpProgMng::SetQuestProg` | 暂无 QuestExp 进度写入入口 | 未覆盖 | 提供统一进度写入函数，供击杀、物品、NPC、事件等钩子调用。 |
| 102 | QuestExp 协议/NPC | CGReqQuestExp | `CGReqQuestExp` | `handle/c1c2.go` 有 `0xF630: questExp` | 部分覆盖 | 实现 QuestExp 主请求解析、参数校验、业务分发和响应。 |
| 103 | QuestExp 协议/NPC | CGReqQuestSwitch | `CGReqQuestSwitch` | `handle/c1c2.go` 有 `0xF60A: questSwitch` | 部分覆盖 | 实现 QuestSwitch 请求，返回指定 Episode 或任务线状态。 |
| 104 | QuestExp 协议/NPC | CGReqQuestProgress | `CGReqQuestProgress` | `handle/c1c2.go` 有 `0xF60B: questProgress` | 部分覆盖 | 实现客户端查询任务进度的协议处理。 |
| 105 | QuestExp 协议/NPC | CGReqQuestComplete | `CGReqQuestComplete` | `handle/c1c2.go` 有 `0xF60D: questComplete` | 部分覆盖 | 实现 QuestExp 完成请求，接入条件检查和奖励发放。 |
| 106 | QuestExp 协议/NPC | CGReqQuestGiveUp | `CGReqQuestGiveUp` | `handle/c1c2.go` 有 `0xF60F: questGiveUp` | 部分覆盖 | 实现 QuestExp 放弃请求，清理状态并通知客户端。 |
| 107 | QuestExp 协议/NPC | CGReqTutorialKeyComplete | `CGReqTutorialKeyComplete` | `handle/c1c2.go` 有 `0xF610: tutorialKeyCOmplete` | 部分覆盖 | 实现教程按键完成请求并推进教程任务状态。 |
| 108 | QuestExp 协议/NPC | CGReqProgressQuestList | `CGReqProgressQuestList` | `handle/c1c2.go` 有 `0xF61A: questProgressList` | 部分覆盖 | 返回玩家进行中的 QuestExp 列表。 |
| 109 | QuestExp 协议/NPC | CGReqProgressQuestInfo | `CGReqProgressQuestInfo` | `handle/c1c2.go` 有 `0xF61B: questProgressInfo` | 部分覆盖 | 返回指定 QuestExp 的详细进度和 ask 状态。 |
| 110 | QuestExp 协议/NPC | CGReqEventItemQuestList | `CGReqEventItemQuestList` | `handle/c1c2.go` 有 `0xF621: eventItemQuestList` | 部分覆盖 | 返回事件道具相关任务列表，支持客户端任务物品 UI。 |
| 111 | server-game 落地点/测试边界 | CGReqAttDefPowerInc | `CGReqAttDefPowerInc` | `handle/c1c2.go` 有 `0xF631: AttDefPowerInc` | 部分覆盖 | 实现攻击/防御能力提升类任务或奖励请求，并校验属性变化来源。 |
| 112 | server-game 落地点/测试边界 | handle/c1c2.go 任务 opcode 入口 | `protocol.cpp` 分发 `0xA0` 与 `0xF6xx` | 已声明多个任务处理函数 | 部分覆盖 | 将现有空入口连接到 Go 任务服务，保持协议层薄、业务层集中。 |
| 113 | server-game 落地点/测试边界 | 任务数据模型缺失 | `QuestInfo`、`Quests`、`QuestExpUserInfo` | `game/model/db.go` 仅有 `ChangeUp` | 未覆盖 | 设计角色任务状态表或字段，区分老任务、扩展任务和 QuestExp。 |
| 114 | server-game 落地点/测试边界 | 任务协议结构缺失 | `QuestExpDefine.h`、`protocol.h` | 缺少任务请求/响应结构体 | 未覆盖 | 补齐 C1/C2 任务协议结构体，明确字段大小、序列化顺序和错误码。 |
| 115 | server-game 落地点/测试边界 | 任务道具拾取/出售/使用集成 | `CGItemGetRequest`、出售限制、`ItemUseQuest` | `KindBQuest` 存在但未接业务 | 未覆盖 | 任务物品应接入掉落、拾取、出售限制、使用和完成扣除链路。 |
| 116 | server-game 落地点/测试边界 | ChangeUp/职业任务集成 | 老任务与转职状态相关逻辑 | `model.Character.ChangeUp` 已存在 | 部分覆盖 | 将职业任务完成结果写入 `ChangeUp`，并影响职业、技能、属性和协议表现。 |
| 117 | server-game 落地点/测试边界 | 经验/金币/道具/技能奖励集成 | `GiveReward`、`SendQuestReward` | 奖励系统未由任务统一调用 | 未覆盖 | 建立任务奖励调用边界，保证奖励发放原子性和失败回滚。 |
| 118 | server-game 落地点/测试边界 | NPC Talk 任务入口集成 | `NpcTalk`、`AddQuestExpNpcTalk` | NPC 交互未接入任务 | 未覆盖 | NPC 点击后应按任务状态返回可接、进行中、可完成和普通对话。 |
| 119 | server-game 落地点/测试边界 | 怪物击杀/事件钩子 | `AddMonsterKillCount`、`KillMonsterProc`、事件任务函数 | 怪物与事件流程未通知任务系统 | 未覆盖 | 在怪物死亡、事件积分、事件完成处调用任务进度更新。 |
| 120 | server-game 落地点/测试边界 | 持久化与端到端回归 | `QuestInfoSave`、`UserQuestInfoSave`、`UserQuestInfoLoad` | 暂无任务持久化测试 | 未覆盖 | 覆盖登录加载、接取、进度、完成、放弃、奖励、登出保存和重登恢复。 |
