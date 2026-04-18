# 32. 服务Bot系统

本模块覆盖服务器放置在地图中的服务型 Bot，包括 Buffer Bot、Trade/Alchemist Bot、BotShop、HideSeek Bot、Bot 对象创建、外观装备、地图坐标、交互事务和 Bot 专属日志。服务Bot系统不等同于玩家自己的 MuBot/OfflineLevelling，也不等同于 server-game 当前 `game/bot` 的模拟连接/压测工具。

明确排除：玩家自动挂机归 `31-helper.md`；模拟连接、模拟登录、压测玩家和开发调试命令归候选“测试Bot/压测系统”；普通 NPC 商店归 `08-shops.md`；玩家个人商店归 `20-personal-shops.md`，BotShop 只在本模块编排并调用个人商店/道具/经济能力。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 模块边界与总入口 | ServiceBotManager 总入口 | `CBotSystem g_BotSystem` | 暂无业务服务 Bot | 未覆盖 | 建立服务型 Bot 管理器，区别于测试 Bot。 |
| 2 | 模块边界与总入口 | BotSystem 初始化 | `CBotSystem::Init` | 暂无 | 未覆盖 | 启动时加载 Bot 配置、专精配置并创建 Bot 对象。 |
| 3 | 模块边界与总入口 | BotSystem 配置加载 | `LoadData` | 暂无 | 未覆盖 | 加载 Bot 基础属性、类型、地图、坐标、装备、文件路径。 |
| 4 | 模块边界与总入口 | Bot 专精加载 | `LoadBotSpecializationData` | 暂无 | 未覆盖 | 按 Bot 类型加载 Buff、炼金、商店等专精数据。 |
| 5 | 模块边界与总入口 | BotSystem 关闭 | 析构/对象释放 | 暂无 | 未覆盖 | 服务关闭时移除所有 Bot 对象并刷出日志。 |
| 6 | 模块边界与总入口 | Bot 类型枚举 | `BOT_BUFFER`、`BOT_HIDENSEEK`、`BOT_TRADE`、`BOT_SHOP` | 暂无 | 未覆盖 | 定义服务 Bot 类型并避免和玩家/NPC 类型混淆。 |
| 7 | 模块边界与总入口 | Bot 数据表 | `m_BotData` | 暂无 | 未覆盖 | 按 Bot ID 保存配置、对象 index 和专精数据。 |
| 8 | 模块边界与总入口 | Bot 数量统计 | `iCount` | 暂无 | 未覆盖 | 提供 Bot 数量和类型统计。 |
| 9 | 模块边界与总入口 | Bot 错误模型 | C++ bool/int 分散 | 暂无 | 未覆盖 | 统一配置错误、创建失败、交互失败、事务失败。 |
| 10 | 模块边界与总入口 | Bot 测试夹具 | GameServer 运行验证 | 暂无 | 未覆盖 | 提供 fake object/player/item 测试 Bot 交互。 |
| 11 | Bot 基础配置 | Bot ID | `_sBOT_SETTINGS::wID` | 暂无 | 未覆盖 | 每个服务 Bot 有稳定配置 ID。 |
| 12 | Bot 基础配置 | Bot 启用开关 | `bEnabled` | 暂无 | 未覆盖 | 未启用 Bot 不创建对象。 |
| 13 | Bot 基础配置 | Bot 名称 | `sName` | 暂无 | 未覆盖 | Bot 对象显示名称和日志名称。 |
| 14 | Bot 基础配置 | Bot 类型 | `btType` | 暂无 | 未覆盖 | 决定交互入口和专精配置。 |
| 15 | Bot 基础配置 | Bot VIP 限制 | `btVipType` | 暂无 | 未覆盖 | 交互时检查玩家 VIP 等级。 |
| 16 | Bot 基础配置 | Bot 货币类型 | `btCoinType` | 暂无 | 未覆盖 | 部分 Bot 交互需要扣费。 |
| 17 | Bot 基础配置 | Bot 货币值 | `iCoinValue` | 暂无 | 未覆盖 | 配置每次交互费用或服务费用。 |
| 18 | Bot 基础配置 | Bot 职业 | `btClass` | 暂无 | 未覆盖 | 用于外观、动作和技能表现。 |
| 19 | Bot 基础配置 | Bot 属性点 | `wStrength/wDexterity/wVitality/wEnergy` | 暂无 | 未覆盖 | 创建对象时填充基础属性。 |
| 20 | Bot 基础配置 | Bot 配置路径 | `sPathActionFile` | 暂无 | 未覆盖 | 指向专精配置文件。 |
| 21 | Bot 位置与对象 | Bot 地图 | `btMap` | 暂无 | 未覆盖 | Bot 创建在指定地图。 |
| 22 | Bot 位置与对象 | Bot 坐标 X/Y | `btX`、`btY` | 暂无 | 未覆盖 | Bot 创建在固定坐标。 |
| 23 | Bot 位置与对象 | Bot 朝向 | `btDir` | 暂无 | 未覆盖 | 创建对象时设置朝向。 |
| 24 | Bot 位置与对象 | Bot 对象创建 | `AddBot` | 暂无 | 未覆盖 | 分配对象池 index 并注册到地图。 |
| 25 | Bot 位置与对象 | 批量创建 Bot | `SetAllBots` | 暂无 | 未覆盖 | 启动时创建所有启用 Bot。 |
| 26 | Bot 位置与对象 | Bot 对象标记 | `ISBOT`、`wBotIndex` | 暂无 | 未覆盖 | 对象上标记服务 Bot 和配置索引。 |
| 27 | Bot 位置与对象 | Bot 类型查询 | `GetBotType` | 暂无 | 未覆盖 | 通过对象 index 查 Bot 类型。 |
| 28 | Bot 位置与对象 | Bot 不可移动 | GameServer 固定位置 | 暂无 | 未覆盖 | 默认服务 Bot 不参与普通移动 AI。 |
| 29 | Bot 位置与对象 | Bot 视野广播 | gObj 创建广播 | 暂无 | 未覆盖 | Bot 创建/删除时刷新地图可见对象。 |
| 30 | Bot 位置与对象 | Bot 重载重建 | Init/SetAllBots | 暂无 | 未覆盖 | 配置重载时安全删除旧 Bot 并创建新 Bot。 |
| 31 | Bot 外观装备 | 穿戴装备配置 | `_sBOT_INVENTORY_WEAR_ITEMS` | 暂无 | 未覆盖 | 配置 Bot 展示装备。 |
| 32 | Bot 外观装备 | 装备 item id | `wItemID` | 暂无 | 未覆盖 | Bot 外观装备道具 ID。 |
| 33 | Bot 外观装备 | 装备等级 | `btItemLv` | 暂无 | 未覆盖 | Bot 外观装备等级。 |
| 34 | Bot 外观装备 | 装备卓越 | `btItemExc` | 暂无 | 未覆盖 | Bot 外观装备卓越标记。 |
| 35 | Bot 外观装备 | PreviewCharSet 生成 | `MakePreviewCharSet` | 暂无 | 未覆盖 | 生成客户端显示角色外观。 |
| 36 | Bot 外观装备 | 职业外观 | `btClass` + preview | 暂无 | 未覆盖 | 根据职业和装备合成外观。 |
| 37 | Bot 外观装备 | 外观协议同步 | gObj viewport | 暂无 | 未覆盖 | Bot 进入视野时下发正确外观。 |
| 38 | Bot 外观装备 | 装备合法性校验 | 配置读取分散 | 暂无 | 未覆盖 | 校验装备槽位、ItemId、等级、卓越值。 |
| 39 | Bot 外观装备 | 外观缓存 | GameServer 即时生成 | 暂无 | 未覆盖 | 可缓存 preview，减少重复计算。 |
| 40 | Bot 外观装备 | 外观测试 | GameServer 运行验证 | 暂无 | 未覆盖 | 覆盖职业、装备、空槽、非法配置。 |
| 41 | Buffer Bot | Buffer Bot 类型 | `BOT_BUFFER` | 暂无 | 未覆盖 | 支持点击 Bot 后给玩家加 Buff。 |
| 42 | Buffer Bot | Buff 配置加载 | `LoadBotSpecializationData` | 暂无 | 未覆盖 | 加载 Bot 可提供的 Buff 列表。 |
| 43 | Buffer Bot | 最大 Buff 数 | `MAX_BUFFS_PER_BOT` | 暂无 | 未覆盖 | 限制每个 Buffer Bot 的 Buff 数量。 |
| 44 | Buffer Bot | BuffPlayer 入口 | `BuffPlayer` | 暂无 | 未覆盖 | 玩家与 Buffer Bot 交互时触发 Buff。 |
| 45 | Buffer Bot | Buff 权限检查 | VIP/费用/状态 | 暂无 | 未覆盖 | 检查等级、VIP、费用、已有 Buff。 |
| 46 | Buffer Bot | Buff 持续时间 | `GetSkillTime` | 暂无 | 未覆盖 | 获取服务 Bot 提供技能的持续时间。 |
| 47 | Buffer Bot | Buff 效果下发 | `gObjAddBuffEffect` | Buff 系统未完整 | 未覆盖 | 调用 Buff 系统添加效果并同步客户端。 |
| 48 | Buffer Bot | Buff 费用扣除 | Bot coin 配置 | 暂无 | 未覆盖 | 需要收费时调用经济系统扣费。 |
| 49 | Buffer Bot | Buff 冷却限制 | GameServer 分散 | 暂无 | 未覆盖 | 防止频繁点击刷 Buff。 |
| 50 | Buffer Bot | Buff Bot 测试 | GameServer 运行验证 | 暂无 | 未覆盖 | 覆盖成功、权限不足、扣费失败、重复 Buff。 |
| 51 | Trade Bot | Trade Bot 类型 | `BOT_TRADE` | 暂无 | 未覆盖 | 支持炼金/合成型服务 Bot。 |
| 52 | Trade Bot | 炼金配置结构 | `sBOT_CRAFTING` | 暂无 | 未覆盖 | 定义需求物品、奖励、成功率、需求数量。 |
| 53 | Trade Bot | 需求物品结构 | `s_BOT_CRAFTING_ITEM_STRUCT` | 暂无 | 未覆盖 | ItemId、等级、幸运、技能、追加、卓越、套装、数量。 |
| 54 | Trade Bot | 奖励结构 | `sBOT_REWARD_STRUCT` | 暂无 | 未覆盖 | 生成奖励 CItem、数量和成功率。 |
| 55 | Trade Bot | 打开炼金交易 | `AlchemistTradeOpen` | 暂无 | 未覆盖 | 玩家点击 Trade Bot 时打开交易/合成界面。 |
| 56 | Trade Bot | 检查炼金条件 | `CheckAlchemist` | 暂无 | 未覆盖 | 检查目标是否 Trade Bot、玩家状态、配置存在。 |
| 57 | Trade Bot | 统计交易栏物品 | `AlchemistTradeItemCount` | 暂无 | 未覆盖 | 统计玩家提交的物品数量。 |
| 58 | Trade Bot | 校验需求物品 | `AlchemistVerifyItem` | 暂无 | 未覆盖 | 按配置校验物品属性。 |
| 59 | Trade Bot | 确认合成成功 | `ConfirmMixSuccess` | 暂无 | 未覆盖 | 根据成功率和需求生成奖励。 |
| 60 | Trade Bot | 炼金交易确认 | `AlchemistTradeOk` | 暂无 | 未覆盖 | 执行扣除材料、发放奖励、失败回滚。 |
| 61 | Trade Bot | 材料数量校验 | `iReqCount/iTotalReqCount` | 暂无 | 未覆盖 | 确认每类和总需求数量满足。 |
| 62 | Trade Bot | 成功率随机 | `iSuccessRate` | 暂无 | 未覆盖 | 调用基础随机工具判定成功。 |
| 63 | Trade Bot | 失败处理 | 合成失败分散 | 暂无 | 未覆盖 | 失败时按配置消耗材料或返还。 |
| 64 | Trade Bot | 炼金事务边界 | 直接改 trade/inventory | 暂无 | 未覆盖 | 调用合成/道具系统执行事务，不在 Bot 内直接改背包。 |
| 65 | BotShop | BotShop 类型 | `BOT_SHOP` | 暂无 | 未覆盖 | 支持服务 Bot 作为商店卖物品。 |
| 66 | BotShop | BotShop 配置结构 | `sBOT_PSHOP` | 暂无 | 未覆盖 | 保存商品、货币类型、店名、商品数量。 |
| 67 | BotShop | BotShop 商品结构 | `s_BOT_SHOP_ITEM` | 暂无 | 未覆盖 | 商品价格、ItemId、等级、属性、socket。 |
| 68 | BotShop | BotShop 店名 | `szBotShopName` | 暂无 | 未覆盖 | 展示个人商店标题。 |
| 69 | BotShop | BotShop 货币类型 | `iCoinType` | 暂无 | 未覆盖 | 支持 Zen/WCoin/GoblinPoint 等货币购买。 |
| 70 | BotShop | BotShop 商品数量 | `iItemCount` | 暂无 | 未覆盖 | 限制最多 32 个商品。 |
| 71 | BotShop | BotShop 添加商品 | `StoreAddItems` | 暂无 | 未覆盖 | 将配置商品放入 Bot 商店容器。 |
| 72 | BotShop | BotShop 空间检查 | `PShopCheckSpace` | 暂无 | 未覆盖 | 检查商品尺寸是否能摆入商店格子。 |
| 73 | BotShop | 临时格子检查 | `gObjTempPShopRectCheck` | 暂无 | 未覆盖 | 复用个人商店格子占用算法。 |
| 74 | BotShop | BotShop 打开 | `PersonalStore.cpp` Bot 分支 | 暂无 | 未覆盖 | 玩家点击 BotShop 时打开商品列表。 |
| 75 | BotShop | BotShop 购买 | `PersonalStore.cpp` Bot 分支 | 暂无 | 未覆盖 | 调用个人商店/道具/经济系统完成购买。 |
| 76 | BotShop | BotShop 货币扣除 | `iCoinType` switch | 暂无 | 未覆盖 | 按商品货币类型扣费。 |
| 77 | BotShop | BotShop 日志 | `BotShopLog` | 暂无 | 未覆盖 | 记录购买者、商品、价格、Bot 名称。 |
| 78 | BotShop | BotShop 库存策略 | 配置静态商品 | 暂无 | 未覆盖 | 明确无限库存或每次重载刷新。 |
| 79 | BotShop | BotShop 搜索边界 | PShop 系统 | 暂无 | 未覆盖 | 是否出现在个人商店搜索由个人商店系统决定。 |
| 80 | BotShop | BotShop 测试 | GameServer 运行验证 | 暂无 | 未覆盖 | 覆盖开店、购买、货币不足、背包满、日志。 |
| 81 | HideSeek Bot | HideSeek Bot 类型 | `BOT_HIDENSEEK` | 暂无 | 未覆盖 | 支持隐藏寻找类服务 Bot。 |
| 82 | HideSeek Bot | HideSeek 位置 | Bot map/x/y | 暂无 | 未覆盖 | Bot 放置在活动指定位置。 |
| 83 | HideSeek Bot | HideSeek 交互 | BotSystem 类型分支 | 暂无 | 未覆盖 | 玩家找到后触发奖励或活动进度。 |
| 84 | HideSeek Bot | HideSeek 活动边界 | 活动系统 | 暂无 | 未覆盖 | 活动规则归普通活动或世界事件系统。 |
| 85 | HideSeek Bot | HideSeek 奖励边界 | 道具/奖励系统 | 暂无 | 未覆盖 | 奖励发放调用掉落/道具/运营奖励能力。 |
| 86 | HideSeek Bot | HideSeek 重置 | 活动重置 | 暂无 | 未覆盖 | 活动结束后删除或隐藏 Bot。 |
| 87 | 交互入口 | NPC 点击协议接入 | `protocol.cpp` 点击对象分支 | handler/object 未完整 | 未覆盖 | 点击服务 Bot 时按 Bot 类型路由。 |
| 88 | 交互入口 | 目标对象校验 | `TargetNumber`/ISBOT | 暂无 | 未覆盖 | 校验目标存在、可见、距离合法、类型正确。 |
| 89 | 交互入口 | 距离校验 | 对象距离 | 暂无 | 未覆盖 | 玩家必须在交互距离内。 |
| 90 | 交互入口 | 玩家状态校验 | 交易/死亡/移动限制 | 暂无 | 未覆盖 | 死亡、交易、个人商店等状态不能交互。 |
| 91 | 交互入口 | Bot 类型路由 | `GetBotType` 分支 | 暂无 | 未覆盖 | Buffer/Trade/BotShop/HideSeek 分别路由。 |
| 92 | 交互入口 | 失败应答 | C++ 分散返回 | 暂无 | 未覆盖 | 交互失败给客户端明确提示。 |
| 93 | 交互入口 | 交互冷却 | 暂无统一 | 暂无 | 未覆盖 | 防止短时间重复点击触发事务。 |
| 94 | 交互入口 | 交互审计 | BotShopLog/日志分散 | 暂无 | 未覆盖 | 关键交易、合成、购买写审计日志。 |
| 95 | 跨系统协作 | 对象系统协作 | `gObj`、ISBOT | `object` | 未覆盖 | 服务 Bot 是对象系统中的特殊对象。 |
| 96 | 跨系统协作 | 地图系统协作 | map/x/y | `map` | 未覆盖 | Bot 注册到地图对象容器。 |
| 97 | 跨系统协作 | Buff 系统协作 | `gObjAddBuffEffect` | `13-buffs.md` | 未覆盖 | Buffer Bot 调用 Buff 系统。 |
| 98 | 跨系统协作 | 技能系统协作 | `GetSkillTime` | `10-skills.md` | 未覆盖 | Buff 技能持续时间和效果来自技能/Buff 规则。 |
| 99 | 跨系统协作 | 道具系统协作 | `CItem`/Inventory | `07-items.md` | 未覆盖 | BotShop 和炼金都必须通过道具系统改物品。 |
| 100 | 跨系统协作 | 合成系统协作 | Alchemist 事务 | `12-mix.md` | 未覆盖 | 炼金 Bot 可调用合成系统完成材料和奖励事务。 |
| 101 | 跨系统协作 | 个人商店协作 | `PersonalStore.cpp` Bot 分支 | `20-personal-shops.md` | 未覆盖 | BotShop 复用个人商店展示和购买能力。 |
| 102 | 跨系统协作 | 商店系统协作 | 普通 NPC 商店 | `08-shops.md` | 未覆盖 | 如果 BotShop 走 NPC 商店模型，需要清晰分界。 |
| 103 | 跨系统协作 | 经济系统协作 | coin/zen 扣除 | 暂无独立经济模块 | 未覆盖 | 所有扣费必须通过统一经济/货币接口。 |
| 104 | 跨系统协作 | 运营系统协作 | BotShopLog/GM reload | `29-ops.md` | 未覆盖 | 运营查看、重载、关闭 Bot。 |
| 105 | 与测试Bot区分 | 测试 Bot 管理器 | 无等价业务 Bot | `game/bot.BotManager` | 已区分 | 当前 Go Bot 是模拟连接，不属于服务 Bot。 |
| 106 | 与测试Bot区分 | AddBot 命令边界 | GameServer `AddBots` 业务创建 | `cmd.AddBot` 模拟玩家 | 已区分 | 不把压测 AddBot 当作服务 Bot 创建。 |
| 107 | 与测试Bot区分 | fake connection 边界 | 无 | `botConn` | 已区分 | fake conn 只用于测试/压测，不参与 BotSystem。 |
| 108 | 与测试Bot区分 | 自动选角边界 | 无 | `pickCharacter` | 已区分 | 模拟登录行为不属于服务 Bot。 |
| 109 | 配置与重载 | Bot 配置文件校验 | XML load | 暂无 | 未覆盖 | 校验 Bot ID、类型、地图、坐标、装备、专精文件。 |
| 110 | 配置与重载 | 专精配置校验 | `LoadBotSpecializationData` | 暂无 | 未覆盖 | 校验 Buff、Crafting、Shop 配置完整性。 |
| 111 | 配置与重载 | 重复 Bot ID 检测 | map 覆盖风险 | 暂无 | 未覆盖 | 重复 ID 必须报错。 |
| 112 | 配置与重载 | 坐标合法性校验 | 地图校验分散 | 暂无 | 未覆盖 | Bot 坐标必须可站立且非阻挡。 |
| 113 | 配置与重载 | 商品合法性校验 | item attr 校验分散 | 暂无 | 未覆盖 | BotShop 商品必须是有效道具和合法价格。 |
| 114 | 配置与重载 | 热重载策略 | GameServer 启动加载 | 暂无 | 未覆盖 | 热重载时处理玩家正在交互的 Bot。 |
| 115 | 日志与审计 | BotSystem 日志 | `g_Log` | 暂无 | 未覆盖 | 创建、删除、配置错误写结构化日志。 |
| 116 | 日志与审计 | BotShop 购买日志 | `BotShopLog` | 暂无 | 未覆盖 | 单独记录 BotShop 购买流水。 |
| 117 | 日志与审计 | 炼金交易日志 | 日志分散 | 暂无 | 未覆盖 | 记录材料、奖励、成功率、结果。 |
| 118 | 日志与审计 | Buff 服务日志 | 暂无 | 暂无 | 未覆盖 | 记录 Buffer Bot 服务玩家与扣费。 |
| 119 | 测试与验收 | Bot 创建测试 | GameServer 运行验证 | 暂无 | 未覆盖 | 覆盖配置加载、对象创建、视野展示、删除。 |
| 120 | 测试与验收 | Bot 交互集成测试 | GameServer 运行验证 | 暂无 | 未覆盖 | 覆盖 Buffer、Trade、BotShop、HideSeek 的成功和失败路径。 |
