# 31. 助手挂机系统

本模块覆盖玩家自己的自动化挂机能力，包括 MuHelper、MuBot、OfflineLevelling、挂机配置保存、启停状态、自动战斗、自动拾取、自动修理、挂机计费、挂机限制和状态同步。助手挂机系统不拥有对象、地图、技能、经验、掉落、道具、经济、安全或运营规则；它只编排玩家自动化行为，并调用对应业务系统完成实际动作。

明确排除：GameServer `BotSystem` 的 Buffer Bot、Trade Bot、BotShop、HideSeek Bot 归 `32-service-bots.md`；server-game `game/bot` 的模拟连接、压测、开发调试归候选“测试Bot/压测系统”。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 模块边界与总入口 | HelperManager 总入口 | `COfflineLevelling`、MuBot 逻辑分散在 `user.cpp`/`protocol.cpp` | 暂无 `game/helper` | 未覆盖 | 建立统一助手挂机入口，管理 MuBot 与离线挂机生命周期。 |
| 2 | 模块边界与总入口 | 助手系统初始化 | `COfflineLevelling::Initiate`、MuBot 配置读取 | `conf.Common.GameServerInfo.MUHelper/MuBot` | 部分覆盖 | 启动时加载开关、计费、限制、技能、拾取和地图配置。 |
| 3 | 模块边界与总入口 | 助手系统关闭 | `DeleteUser`/线程退出分散 | 暂无 | 未覆盖 | 服务关闭时停用所有挂机玩家并保存必要状态。 |
| 4 | 模块边界与总入口 | 在线 MuBot 与离线挂机边界 | MuBot 在线，OfflineLevelling 可离线执行 | 只实现 MuBot 协议雏形 | 部分覆盖 | 明确在线官方挂机和服务端离线挂机的状态模型。 |
| 5 | 模块边界与总入口 | 业务编排边界 | 直接调用 `gObj`/技能/拾取 | 暂无统一接口 | 未覆盖 | 助手只发起动作，实际战斗、拾取、修理由对应系统执行。 |
| 6 | 模块边界与总入口 | 并发安全策略 | `CRITICAL_SECTION m_OfflevelCriti` | Go channel/对象协程模型 | 未覆盖 | 所有挂机请求通过对象或游戏主循环串行化，避免跨协程改对象。 |
| 7 | 模块边界与总入口 | 助手状态模型 | `m_MuBotEnable`、`m_OffPlayerData` | `MsgEnableMuBot` 无真实状态 | 部分覆盖 | 定义关闭、在线挂机、离线挂机、暂停、计费失败等状态。 |
| 8 | 模块边界与总入口 | 助手错误模型 | 协议 status 分散 | `MsgEnableMuBotReply.Flag` | 部分覆盖 | 统一成功、拒绝、余额不足、等级不足、地图禁止、状态冲突。 |
| 9 | 模块边界与总入口 | 助手配置快照 | `g_ConfigRead.data.mubot`、离线挂机配置 | conf 字段已读 | 部分覆盖 | 玩家启用时使用配置快照，避免运行中配置变化造成状态不一致。 |
| 10 | 模块边界与总入口 | 助手测试夹具 | GameServer 运行验证 | 暂无 | 未覆盖 | 提供 fake player/map/skill/item 测试自动化行为。 |
| 11 | MuBot 协议 | MuBot 配置结构 | `PMSG_MUBOT_DATASAVE`、257 字节 | `MsgMuBot.Data [257]byte` | 已覆盖 | 保持客户端 257 字节配置兼容。 |
| 12 | MuBot 协议 | MuBot 配置解析 | 原样保存并下发 | `MsgMuBot.Unmarshal` | 已覆盖 | 当前按原始字节保存，后续可增加结构化解析。 |
| 13 | MuBot 协议 | MuBot 配置序列化 | `GCAnsMuBotData` | `MsgMuBot.Marshal` | 已覆盖 | 保持协议字节顺序。 |
| 14 | MuBot 协议 | 保存 MuBot 配置请求 | `CGReqMuBotSaveData` | `DefineMuBot` | 已覆盖 | 玩家请求保存配置时写入角色持久化。 |
| 15 | MuBot 协议 | 保存 MuBot 到 DataServer | `GDReqMuBotSave` | `DB.UpdateCharacterMuBot` | 部分覆盖 | server-game 直接 DB 保存，外部 DataServer 可后续接入。 |
| 16 | MuBot 协议 | 加载 MuBot 配置应答 | `DGAnsMuBotData`、`GCAnsMuBotData` | `MsgMuBotReply` | 部分覆盖 | 角色进入游戏时下发已保存配置。 |
| 17 | MuBot 协议 | 启用 MuBot 请求 | `CGReqMuBotUse` | `EnableMuBot` | 部分覆盖 | 增加等级、地图、费用、状态校验后再启用。 |
| 18 | MuBot 协议 | 启用 MuBot 应答 | `GCAnsMuBotUse` | `MsgEnableMuBotReply` | 部分覆盖 | 返回状态、时间、倍率/时间单位、费用。 |
| 19 | MuBot 协议 | 关闭 MuBot 请求 | `CGReqMuBotUse` flag off | `EnableMuBot` case 1 空实现 | 未覆盖 | 关闭时清理计时和状态并下发应答。 |
| 20 | MuBot 协议 | 重连保护 | protocol reconnect bypass MuBot | 暂无 | 未覆盖 | 玩家重连时恢复或关闭 MuBot 状态，避免重复计费。 |
| 21 | MuBot 状态 | MuBot 启用字段 | `m_MuBotEnable` | 暂无玩家字段 | 未覆盖 | Player 保存当前 MuBot 是否启用。 |
| 22 | MuBot 状态 | MuBot 总时长 | `m_MuBotTotalTime` | Reply 有 `Time` | 部分覆盖 | 记录本次启用累计分钟或服务器定义时间单位。 |
| 23 | MuBot 状态 | MuBot 计费时间 | `m_MuBotPayTime` | 暂无 | 未覆盖 | 每个计费周期累积并触发扣费。 |
| 24 | MuBot 状态 | MuBot tick | `m_MuBotTick` | 暂无 | 未覆盖 | 用服务器时钟驱动每分钟统计。 |
| 25 | MuBot 状态 | MuBot 自动关闭时间 | `autodisabletime` | `AutoDisableTime` 配置 | 部分覆盖 | 达到配置时长自动关闭并通知客户端。 |
| 26 | MuBot 状态 | MuBot 状态同步 | `GCAnsMuBotUse` | `EnableMuBotReply` | 部分覆盖 | 状态改变必须同步客户端。 |
| 27 | MuBot 状态 | MuBot 下线处理 | user close 清状态 | 暂无 | 未覆盖 | 玩家下线时按在线挂机语义关闭或转离线挂机。 |
| 28 | MuBot 状态 | MuBot PVP 中断 | `NewPVP` 关闭 MuBot | 暂无 | 未覆盖 | 进入 PVP/决斗等状态时强制关闭或暂停。 |
| 29 | MuBot 状态 | MuBot GM 操作中断 | `GMMng` 关闭 MuBot | 暂无 | 未覆盖 | GM 移动、踢线、状态修改时处理挂机状态。 |
| 30 | MuBot 状态 | MuBot 状态持久化边界 | 配置持久，启用状态运行期 | DB 只保存配置 | 部分覆盖 | 不默认持久化启用状态，除非离线挂机明确需要。 |
| 31 | MuBot 限制 | MuBot 总开关 | `mubot.enable` | `MuBot.Enable` | 部分覆盖 | 关闭时拒绝所有启用请求。 |
| 32 | MuBot 限制 | MuBot 最低等级 | `mubot.minlevel` | `MinLevel` | 部分覆盖 | 玩家等级不足拒绝启用。 |
| 33 | MuBot 限制 | MuBot VIP 等级 | `NeedVIPLevel` | `NeedVIPLevel` | 部分覆盖 | VIP 不足拒绝或走不同费用。 |
| 34 | MuBot 限制 | MuBot 地图限制 | GameServer 分散限制 | 暂无 | 未覆盖 | 禁止安全区、特殊副本、PVP/交易关键状态启用。 |
| 35 | MuBot 限制 | MuBot 状态冲突 | 交易/商店/PVP 中断 | 暂无 | 未覆盖 | 交易、个人商店、合成、仓库等操作中不能启用。 |
| 36 | MuBot 限制 | MuBot 死亡限制 | 死亡状态不可继续 | 暂无 | 未覆盖 | 死亡后暂停或关闭挂机。 |
| 37 | MuBot 限制 | MuBot 移动限制 | 特殊移动关闭 | 暂无 | 未覆盖 | 使用 Gate/MoveCommand 后重新校验地图规则。 |
| 38 | MuBot 限制 | MuBot 账号限制 | 账号状态/封禁 | 暂无 | 未覆盖 | 账号被封禁、冻结、踢线时关闭挂机。 |
| 39 | MuBot 限制 | MuBot 防滥用限制 | AntiHack/速度检测协作 | 暂无 | 未覆盖 | 自动行为必须标记来源，避免误判或绕过风控。 |
| 40 | MuBot 限制 | MuBot 限制测试 | GameServer 无统一 | 暂无 | 未覆盖 | 覆盖等级、VIP、地图、状态冲突、关闭开关。 |
| 41 | MuBot 计费 | MuBot 费用配置 | `mubot.cost` | `Cost` | 部分覆盖 | 读取每周期费用。 |
| 42 | MuBot 计费 | MuBot 货币类型 | GameServer 结合 coin/money | `MuBot` 配置未完整使用 | 未覆盖 | 明确 Zen/WCoin/GoblinPoint/Ruud 等货币来源。 |
| 43 | MuBot 计费 | MuBot 扣费周期 | `m_MuBotPayTime == 5` | 暂无 | 未覆盖 | 达到周期后扣费，失败则关闭。 |
| 44 | MuBot 计费 | MuBot 扣费函数 | `gObjMuBotPayForUse` | 暂无 | 未覆盖 | 调用经济/账号货币系统执行扣费。 |
| 45 | MuBot 计费 | MuBot 余额不足 | 扣费失败关闭 | 暂无 | 未覆盖 | 余额不足时关闭并下发失败应答。 |
| 46 | MuBot 计费 | MuBot 免费策略 | VIP/配置可能免费 | 暂无 | 未覆盖 | VIP 或活动减免费用由配置和运营系统控制。 |
| 47 | MuBot 计费 | MuBot 扣费审计 | GameServer 日志分散 | 暂无 | 未覆盖 | 记录角色、账号、费用、余额、周期。 |
| 48 | MuBot 计费 | MuBot 金额应答 | `GCAnsMuBotUse` money | Reply 有 `Money` | 部分覆盖 | 应答客户端本周期扣费金额。 |
| 49 | MuBot 计费 | MuBot 计费事务 | 扣费与状态分散 | 暂无 | 未覆盖 | 扣费成功才继续挂机，失败原子关闭。 |
| 50 | MuBot 计费 | MuBot 计费测试 | GameServer 无统一 | 暂无 | 未覆盖 | 覆盖扣费成功、余额不足、免费、周期边界。 |
| 51 | OfflineLevelling 配置 | 离线挂机总开关 | `b_Enabled` | 暂无 | 未覆盖 | 独立于 MuBot 的服务端离线挂机开关。 |
| 52 | OfflineLevelling 配置 | 通用配置加载 | `LoadFile`、`OFF_LEVEL_GENERAL_SETTINGS` | 暂无 | 未覆盖 | 加载 VIP、货币、费用、攻击间隔、最大时长。 |
| 53 | OfflineLevelling 配置 | 技能定义加载 | `LoadSkillDefinitions` | 暂无 | 未覆盖 | 加载职业可用离线挂机技能。 |
| 54 | OfflineLevelling 配置 | 技能类型配置 | `OFF_SKILL_CATEGORIES` | 暂无 | 未覆盖 | 区分普通攻击、范围攻击、持续魔法。 |
| 55 | OfflineLevelling 配置 | 地图属性配置 | `OFF_LEVEL_PER_MAP_ATTRIBUTE` | 暂无 | 未覆盖 | 每地图禁用、费用、货币覆盖。 |
| 56 | OfflineLevelling 配置 | 拾取配置 | `OFF_LEVEL_ITEM_PICK_SETTINGS` | 暂无 | 未覆盖 | 配置 Zen、卓越、Socket、Ancient 拾取。 |
| 57 | OfflineLevelling 配置 | 拾取白名单 | `OFF_LEVEL_ITEM_PICK_LIST` | 暂无 | 未覆盖 | 指定物品 ID 自动拾取。 |
| 58 | OfflineLevelling 配置 | 自动修理配置 | `AutoRepairItems` | 暂无 | 未覆盖 | 配置是否挂机自动修理装备。 |
| 59 | OfflineLevelling 配置 | 最大时长配置 | `MaxDuration`、`MaxVipDuration` | 暂无 | 未覆盖 | VIP 与普通玩家使用不同最大时长。 |
| 60 | OfflineLevelling 配置 | 配置校验 | LoadFile 返回 BOOL | 暂无 | 未覆盖 | 校验地图范围、技能 ID、费用、间隔、时长。 |
| 61 | OfflineLevelling 玩家 | 添加离线挂机玩家 | `AddUser` | 暂无 | 未覆盖 | 玩家断线或命令启用时注册离线挂机。 |
| 62 | OfflineLevelling 玩家 | 删除离线挂机玩家 | `DeleteUser` | 暂无 | 未覆盖 | 重新登录、死亡、余额不足、超时后移除。 |
| 63 | OfflineLevelling 玩家 | 查找离线挂机玩家 | `FindUser` | 暂无 | 未覆盖 | 按对象 index/角色 ID 查找挂机记录。 |
| 64 | OfflineLevelling 玩家 | 挂机玩家表 | `m_OffPlayerData` | 暂无 | 未覆盖 | 保存 aIndex、技能、离线开始时间。 |
| 65 | OfflineLevelling 玩家 | 离线对象保留 | gObj 保持对象 | 暂无 | 未覆盖 | 明确 server-game 是否保留 Player/Object 参与地图。 |
| 66 | OfflineLevelling 玩家 | 重新登录恢复 | reconnect 处理 | 暂无 | 未覆盖 | 玩家登录时停止离线挂机并恢复正常连接。 |
| 67 | OfflineLevelling 玩家 | 离线挂机数量统计 | `GetOffLevelerCount` | 暂无 | 未覆盖 | 提供运营查看和容量限制。 |
| 68 | OfflineLevelling 玩家 | 离线挂机容量限制 | GameServer map/对象限制 | 暂无 | 未覆盖 | 防止离线对象耗尽对象池。 |
| 69 | OfflineLevelling 玩家 | 离线挂机状态审计 | 日志分散 | 暂无 | 未覆盖 | 记录加入、移除、原因和持续时间。 |
| 70 | OfflineLevelling 玩家 | 离线挂机测试 | GameServer 无统一 | 暂无 | 未覆盖 | 覆盖 Add/Delete/Find/relogin/timeout。 |
| 71 | 自动战斗 | 离线挂机主循环 | `Run`、`OffLevelThread` | 暂无 | 未覆盖 | 定时扫描挂机玩家并执行自动战斗。 |
| 72 | 自动战斗 | 找怪逻辑 | `FindAndAttack` | 暂无 | 未覆盖 | 在地图视野/范围内选择可攻击怪物。 |
| 73 | 自动战斗 | 目标过滤 | `FindAndAttack` 内过滤 | 暂无 | 未覆盖 | 排除死亡、不可见、安全区、非怪物目标。 |
| 74 | 自动战斗 | 技能查找 | `Magic` 查找失败日志 | 暂无 | 未覆盖 | 根据配置技能 ID 找玩家已学技能。 |
| 75 | 自动战斗 | 技能攻击类型 | `GetSkillAttackType` | 暂无 | 未覆盖 | 根据技能类型选择单体、范围、持续技能逻辑。 |
| 76 | 自动战斗 | 技能距离检查 | `SkillDistanceCheck` | 暂无 | 未覆盖 | 调用地图/技能范围规则判断是否可释放。 |
| 77 | 自动战斗 | 攻击间隔 | `HitInterval` | 暂无 | 未覆盖 | 使用配置控制自动攻击频率。 |
| 78 | 自动战斗 | 普通攻击 fallback | GameServer 可 fallback | 暂无 | 未覆盖 | 技能不可用时是否普通攻击由配置决定。 |
| 79 | 自动战斗 | MP/AG 消耗检查 | 技能系统判断 | 暂无 | 未覆盖 | 自动释放仍必须通过技能系统消耗校验。 |
| 80 | 自动战斗 | 自动战斗来源标记 | GameServer 直接调用 | 暂无 | 未覆盖 | 标记为 helper action，供安全与日志区分。 |
| 81 | 自动拾取 | 拾取入口 | `CheckAndPickUpItem` | 暂无 | 未覆盖 | 地图掉落出现或循环扫描时检查拾取。 |
| 82 | 自动拾取 | Zen 拾取 | `PickUpZen` | 暂无 | 未覆盖 | 根据配置自动拾取金币。 |
| 83 | 自动拾取 | 卓越物品拾取 | `PickUpExc` | 暂无 | 未覆盖 | 根据道具属性判断卓越物品。 |
| 84 | 自动拾取 | Socket 物品拾取 | `PickUpSocket` | 暂无 | 未覆盖 | 根据道具 socket 属性判断。 |
| 85 | 自动拾取 | Ancient 物品拾取 | `PickUpAncient` | 暂无 | 未覆盖 | 根据套装/Ancient 属性判断。 |
| 86 | 自动拾取 | 白名单物品拾取 | `m_PickItems` | 暂无 | 未覆盖 | 指定 ItemId 优先拾取。 |
| 87 | 自动拾取 | 拾取距离检查 | 地图 item 距离 | 暂无 | 未覆盖 | 只有范围内掉落可拾取。 |
| 88 | 自动拾取 | 背包空间检查 | 道具系统 | 暂无 | 未覆盖 | 背包满则跳过或关闭拾取。 |
| 89 | 自动拾取 | 掉落归属检查 | 掉落系统 | 暂无 | 未覆盖 | 不得拾取无权归属的掉落。 |
| 90 | 自动拾取 | 拾取事务 | `CMapItem`/Inventory | 暂无 | 未覆盖 | 调用道具/掉落系统完成原子拾取。 |
| 91 | 自动修理 | 自动修理入口 | `CheckRepairItems` | 暂无 | 未覆盖 | 周期性检查装备耐久。 |
| 92 | 自动修理 | 修理开关 | `AutoRepairItems` | 暂无 | 未覆盖 | 配置关闭时不自动修理。 |
| 93 | 自动修理 | 耐久阈值 | GameServer 内部判断 | 暂无 | 未覆盖 | 装备耐久低于阈值时尝试修理。 |
| 94 | 自动修理 | 修理费用 | 道具/经济系统 | 暂无 | 未覆盖 | 调用商店/道具修理公式计算费用。 |
| 95 | 自动修理 | 修理事务 | 道具系统 | 暂无 | 未覆盖 | 扣费和恢复耐久必须原子执行。 |
| 96 | 自动修理 | 修理失败处理 | 日志/跳过 | 暂无 | 未覆盖 | 余额不足或物品不可修理时记录并继续/关闭。 |
| 97 | 使用时长 | 使用时间检查 | `CheckUseTime` | 暂无 | 未覆盖 | 检查普通/VIP 最大离线挂机时长。 |
| 98 | 使用时长 | 离线开始时间 | `dwOffTime` | 暂无 | 未覆盖 | 注册挂机时记录服务器时间。 |
| 99 | 使用时长 | 普通最大时长 | `MaxDuration` | 暂无 | 未覆盖 | 普通玩家超时移除。 |
| 100 | 使用时长 | VIP 最大时长 | `MaxVipDuration` | 暂无 | 未覆盖 | VIP 玩家使用 VIP 时长。 |
| 101 | 使用时长 | 超时移除 | `CheckUseTime` false | 暂无 | 未覆盖 | 超时关闭挂机并记录原因。 |
| 102 | 使用时长 | 时间统计应答 | 日志/运营统计 | 暂无 | 未覆盖 | 提供累计时长给运营系统查看。 |
| 103 | 计费与货币 | 离线挂机计费入口 | `ChargePlayer` | 暂无 | 未覆盖 | 按配置周期扣除货币。 |
| 104 | 计费与货币 | 离线挂机货币类型 | `CoinType` | 暂无 | 未覆盖 | 支持地图覆盖和通用配置货币类型。 |
| 105 | 计费与货币 | 离线挂机费用 | `CoinValue` | 暂无 | 未覆盖 | 普通费用和地图费用覆盖。 |
| 106 | 计费与货币 | VIP 费用策略 | `VipType` | 暂无 | 未覆盖 | VIP 降低/免除/改变计费策略。 |
| 107 | 计费与货币 | 计费间隔 | `ChargeInterval` | 暂无 | 未覆盖 | 达到间隔才扣费。 |
| 108 | 计费与货币 | 扣费失败关闭 | `ChargePlayer` false | 暂无 | 未覆盖 | 余额不足移除离线挂机。 |
| 109 | 计费与货币 | 计费审计 | 离线挂机日志 | 暂无 | 未覆盖 | 记录扣费金额、货币类型和剩余余额。 |
| 110 | 计费与货币 | 计费事务测试 | GameServer 无统一 | 暂无 | 未覆盖 | 覆盖扣费成功、失败、地图覆盖、VIP。 |
| 111 | 跨系统协作 | 对象系统协作 | `gObj` | `object.Player` | 部分覆盖 | 读取/修改玩家状态必须通过对象系统接口。 |
| 112 | 跨系统协作 | 地图系统协作 | 地图属性/怪物查找 | 暂无 | 未覆盖 | 地图提供可挂机、怪物查询、掉落查询。 |
| 113 | 跨系统协作 | 技能系统协作 | `ObjUseSkill` | 技能系统未完整 | 未覆盖 | 自动释放必须走技能系统。 |
| 114 | 跨系统协作 | 经验系统协作 | 杀怪经验 | 经验系统 | 未覆盖 | 自动杀怪收益由经验系统结算。 |
| 115 | 跨系统协作 | 掉落系统协作 | `CMapItem` | 掉落系统 | 未覆盖 | 掉落创建和归属不由助手系统决定。 |
| 116 | 跨系统协作 | 道具系统协作 | 背包/修理 | 道具系统 | 未覆盖 | 拾取、修理、空间检查走道具系统。 |
| 117 | 跨系统协作 | 安全系统协作 | PVP/AntiHack 限制 | 安全系统 | 未覆盖 | 防止挂机绕过风控和 PVP 状态限制。 |
| 118 | 跨系统协作 | 运营系统协作 | GM 可干预 | 运营系统 | 未覆盖 | 提供查询、关闭、统计、审计入口。 |
| 119 | 测试与验收 | MuBot 协议测试 | GameServer 运行验证 | 暂无 | 未覆盖 | 覆盖保存、加载、启用、关闭、错误应答。 |
| 120 | 测试与验收 | 离线挂机集成测试 | GameServer 运行验证 | 暂无 | 未覆盖 | 覆盖自动战斗、拾取、修理、计费、超时、重登。 |
