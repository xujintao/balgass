# 14. 掉落系统

本模块覆盖怪物死亡掉落、普通物品、卓越物品、金币、Bag 掉落、事件掉落、指定掉落、套装掉落、掉落倍率、掉落归属、地图落物和掉落物品属性生成。`server-game` 已有基础怪物掉落雏形，本模块不是从零开始；但 `GameServer` 中 Bag、特殊掉落、活动掉落、归属与倍率等完整掉落链路仍大部分未覆盖。本模块只记录“掉落”边界，不把任务奖励、商城奖励、活动结算奖励、GremoryCase、积分或现金币奖励作为核心范围。Bag 业务和掉落结果归本模块，LuaBag 的脚本运行时、Lua 函数调用和 Go/Lua 绑定归 `26-script.md`。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | server-game 已有基础掉落 | DropManager 全局实例 | `g_MonsterItemMng` | `game/object/monster/drop.go::DropManager` | 部分覆盖 | 已有基础掉落管理器，后续需明确是否保留在 monster 包或迁移为独立 drop 模块。 |
| 2 | server-game 已有基础掉落 | DropManager.init | `CMonsterItemMng::Init` | `dropManager.init` | 部分覆盖 | 已加载基础掉落率并构建候选池，缺 Bag、特殊掉落和归属上下文。 |
| 3 | server-game 已有基础掉落 | IGC_MonsterItemDropRate 加载 | `LoadMonsterItemDropRate` | `conf.XML(...IGC_MonsterItemDropRate.xml...)` | 已覆盖 | 基础怪物等级掉落率已读取。 |
| 4 | server-game 已有基础掉落 | makeMagicBook | `MagicBookGiveItemSearch` | `dropManager.makeMagicBook` | 部分覆盖 | 已按怪物等级生成魔法书候选，需对齐 GameServer 过滤规则。 |
| 5 | server-game 已有基础掉落 | makeItem | `NormalGiveItemSearch` | `dropManager.makeItem` | 部分覆盖 | 已按物品表和怪物等级生成普通候选，缺地图/事件/职业等额外过滤。 |
| 6 | server-game 已有基础掉落 | makeJewel | `CMonsterItemMng::MakeJewelItem` | `dropManager.makeJewel` | 已覆盖 | 已生成祝福、灵魂、生命、创造、玛雅基础宝石。 |
| 7 | server-game 已有基础掉落 | DropItem | `CMonsterItemMng::GetItemEx` | `dropManager.DropItem` | 部分覆盖 | 已能按等级掉普通掉落，缺特殊属性和高级过滤。 |
| 8 | server-game 已有基础掉落 | DropItemExcellent | `CMonsterItemMng::GetItemExcel` | `dropManager.DropItemExcellent` | 部分覆盖 | 已能生成卓越候选，卓越属性由 `ExcellentDropManager` 补充。 |
| 9 | server-game 已有基础掉落 | Monster.DieDropItem | `gObjMonsterDieGiveItem` | `game/object/monster/attack.go::DieDropItem` | 部分覆盖 | 已在怪物死亡后掉物品或金币，但编排能力远少于 GameServer。 |
| 10 | server-game 已有基础掉落 | 地图落物 | `MapC[].MonsterItemDrop/MoneyItemDrop` | `maps.MapManager.AddItem` | 部分覆盖 | 已能把物品放到地图，缺所有权、拾取限制和完整生命周期。 |
| 11 | GameServer 怪物死亡掉落入口 | gObjMonsterDieGiveItem | `gObjMonster.cpp::gObjMonsterDieGiveItem` | `Monster.DieDropItem` | 部分覆盖 | 对齐 GameServer 的死亡掉落主编排，但 Go 侧应避免函数继续膨胀。 |
| 12 | GameServer 怪物死亡掉落入口 | 延迟掉落入口 | `gObjAddMsgSendDelay(... die drop item ...)` | `Object.processDelayMsg` | 部分覆盖 | Go 侧已有延迟消息入口，需确认掉落延迟、目标有效性和死亡状态。 |
| 13 | GameServer 怪物死亡掉落入口 | 最大伤害玩家归属 | `MaxHitUser/iMaxHitUser` | 当前主要使用击杀者 `tobj` | 未覆盖 | 需要记录伤害归属，掉落不一定属于最后一击。 |
| 14 | GameServer 怪物死亡掉落入口 | 怪物 Herd 掉落 | `MonsterHerdItemDrop` | 暂无 MonsterHerd | 未覆盖 | 怪物组或 Herd 专属掉落应优先于普通掉落。 |
| 15 | GameServer 怪物死亡掉落入口 | 任务掉落短路 | `g_QuestInfo.MonsterItemDrop` | 任务掉落未接入 | 未覆盖 | 任务物品掉落应在普通掉落前处理并可能短路。 |
| 16 | GameServer 怪物死亡掉落入口 | 事件掉落短路 | `gEventMonsterItemDrop` | 暂无事件掉落入口 | 未覆盖 | 事件物品掉落应独立判断并可阻止普通掉落。 |
| 17 | GameServer 怪物死亡掉落入口 | 特殊掉落顺序 | `Bag/NewPVP/Arca/SetItem/CItemDrop` | 暂无统一顺序 | 未覆盖 | 明确各类掉落模块的优先级和短路策略。 |
| 18 | GameServer 怪物死亡掉落入口 | 普通掉落兜底 | `g_MonsterItemMng.GetItemEx` | `DropManager.DropItem` | 部分覆盖 | 当特殊掉落都未命中时走基础物品/金币掉落。 |
| 19 | GameServer 怪物死亡掉落入口 | 掉落坐标选择 | `gObjGetRandomItemDropLocation` | 当前使用怪物坐标 | 未覆盖 | 应在怪物附近寻找可落物坐标，避免阻挡点或安全区异常。 |
| 20 | GameServer 怪物死亡掉落入口 | 禁止掉落地图规则 | `IsCanNotItemDtopInDevilSquare` 等 | 暂无地图掉落过滤 | 未覆盖 | 某些地图或事件应禁止部分物品落地。 |
| 21 | MonsterItemMng 基础掉落表 | CMonsterItemMng::Init | `CMonsterItemMng::Init` | `DropManager.init` | 部分覆盖 | 对齐初始化顺序、清理旧数据和 Jewel/普通/卓越表构建。 |
| 22 | MonsterItemMng 基础掉落表 | Clear | `CMonsterItemMng::Clear` | 暂无 Clear | 未覆盖 | 支持重载配置或测试隔离时清理掉落候选池。 |
| 23 | MonsterItemMng 基础掉落表 | InsertItem | `CMonsterItemMng::InsertItem` | 暂无手动插入候选 | 未覆盖 | 支持从脚本或配置插入指定怪物等级候选物品。 |
| 24 | MonsterItemMng 基础掉落表 | GetItem | `CMonsterItemMng::GetItem` | `DropItem` 近似 | 部分覆盖 | 普通获取接口需区分旧版与扩展版掉落逻辑。 |
| 25 | MonsterItemMng 基础掉落表 | gObjGiveItemSearch | `gObjGiveItemSearch` | `makeItem` 近似 | 部分覆盖 | 按怪物等级和最大物品等级生成候选池。 |
| 26 | MonsterItemMng 基础掉落表 | LoadMonsterItemDropRate | `LoadMonsterItemDropRate` | 已读取 XML | 已覆盖 | 掉落率表基本对应，后续校验字段和比例精度。 |
| 27 | MonsterItemMng 基础掉落表 | gObjGiveItemSearchEx | `gObjGiveItemSearchEx` | 暂无 Ex 候选生成 | 未覆盖 | GameServer 扩展候选生成逻辑需单独对齐。 |
| 28 | MonsterItemMng 基础掉落表 | NormalGiveItemSearchEx | `NormalGiveItemSearchEx` | 暂无 Ex 普通候选生成 | 未覆盖 | 扩展普通物品候选过滤规则未实现。 |
| 29 | MonsterItemMng 基础掉落表 | CheckMonsterDropItem | `CheckMonsterDropItem` | `ItemTable.GetItemLevel` 部分过滤 | 部分覆盖 | 应补全怪物可掉物品过滤规则。 |
| 30 | MonsterItemMng 基础掉落表 | MonsterItemDropRate 总和 | `m_TotalRate` | 当前手工累加分段 | 部分覆盖 | 统一总概率和分段边界，避免概率漏区。 |
| 31 | 普通/卓越/金币掉落 | 基础掉落概率 | `gItemDropPer` | `ItemDropPercent` + `ItemDropRate` | 部分覆盖 | Go 侧已有基础概率，需接入倍率叠加和地图加成。 |
| 32 | 普通/卓越/金币掉落 | 卓越掉落概率 | `ExcelItemDropPercent` | `conf.CommonServer.GameServerInfo.ExcelItemDropPercent` | 部分覆盖 | 已有卓越掉落概率判断，需对齐 GameServer 范围和分母。 |
| 33 | 普通/卓越/金币掉落 | 卓越选项数量 | Excellent option drop | `ExcellentDropManager.dropExcellentCount` | 部分覆盖 | 已随机卓越数量，需对齐各类型权重和上限。 |
| 34 | 普通/卓越/金币掉落 | 卓越选项类型 | `DropExcellent` | `ExcellentDropManager.DropExcellent` | 部分覆盖 | 已生成卓越 bit，需补所有物品类型合法性。 |
| 35 | 普通/卓越/金币掉落 | 普通物品 skill/luck | `ItemSkillDropPercent/ItemLuckyDropPercent` | `DieDropItem` 已判断 | 部分覆盖 | 普通掉落已随机 skill/luck，需对齐 option 与其他属性。 |
| 36 | 普通/卓越/金币掉落 | 卓越物品 skill/luck | `ExcelItemSkillDropPercent/ExcelItemLuckyDropPercent` | `DieDropItem` 已判断 | 部分覆盖 | 卓越物品 skill/luck 概率已接入，需补与卓越属性互斥规则。 |
| 37 | 普通/卓越/金币掉落 | 物品等级 | `DropItem->m_Level` | `it.Level = dit.level` | 部分覆盖 | 掉落等级已设置，需补随机等级、最大等级和特殊物品等级规则。 |
| 38 | 普通/卓越/金币掉落 | 金币掉落概率 | `MonsterMoneyDrop/MoneyDropRate` | `MoneyDropRate` + `MoneyDrop` | 部分覆盖 | Go 侧已有金币掉落，需对齐地图和角色倍率。 |
| 39 | 普通/卓越/金币掉落 | 金币数量计算 | `MoneyItemDrop` 分支 | `MoneyDrop * ZenDropMultiplier` | 部分覆盖 | 已有基础金币数量，缺地图、Buff、事件、玩家属性倍率。 |
| 40 | 普通/卓越/金币掉落 | 掉落失败兜底 | `if (!DropItem) return` | `nil` 检查 | 部分覆盖 | 需明确无候选、无坐标、无地图容量时日志和回退策略。 |
| 41 | BagManager 与 Bag 基础 | InitBagManager | `CBagManager::InitBagManager` | 暂无 BagManager | 未覆盖 | 初始化全部 Bag 配置和索引；底层 LuaBag 加载能力归 `26-script.md`。 |
| 42 | BagManager 与 Bag 基础 | AddItemBag | `CBagManager::AddItemBag` | 暂无 Bag 注册 | 未覆盖 | 支持 CommonBag、MonsterBag、EventBag 按类型注册；Lua `AddItemBag` 绑定归 `26-script.md`。 |
| 43 | BagManager 与 Bag 基础 | DeleteItemBags | `CBagManager::DeleteItemBags` | 暂无 Bag 清理 | 未覆盖 | 支持配置重载或关闭时释放 Bag。 |
| 44 | BagManager 与 Bag 基础 | IsBag | `CBagManager::IsBag` | 暂无 Bag 命中判断 | 未覆盖 | 判断当前条件是否存在可用 Bag。 |
| 45 | BagManager 与 Bag 基础 | SearchAndUseBag | `CBagManager::SearchAndUseBag` | 暂无 Bag 使用 | 未覆盖 | 按类型和参数查找 Bag 并生成掉落。 |
| 46 | BagManager 与 Bag 基础 | GetItemFromBag | `CBagManager::GetItemFromBag` | 暂无 Bag 取物品 | 未覆盖 | 从 Bag 中抽取一个待掉落物品和期限。 |
| 47 | BagManager 与 Bag 基础 | CBag::LoadBag | `CBag::LoadBag` | 暂无 Bag 文件解析 | 未覆盖 | 解析 Bag 的 DropSection、ItemsSection、金币、召唤物等配置。 |
| 48 | BagManager 与 Bag 基础 | CBag::GetDropSection | `CBag::GetDropSection` | 暂无 DropSection | 未覆盖 | 按概率选择 Bag 掉落段。 |
| 49 | BagManager 与 Bag 基础 | CBag::GetItemsSection | `CBag::GetItemsSection` | 暂无 ItemsSection | 未覆盖 | 从掉落段中选择物品段。 |
| 50 | BagManager 与 Bag 基础 | CBag::GetItem | `CBag::GetItem` | 暂无 Bag 物品生成 | 未覆盖 | 根据 Bag 配置生成具体掉落物品。 |
| 51 | CommonBag/MonsterBag/EventBag | CommonBag::SetBagInfo | `CCommonBag::SetBagInfo` | 暂无 CommonBag | 未覆盖 | CommonBag 绑定物品 ID 和等级条件。 |
| 52 | CommonBag/MonsterBag/EventBag | CommonBag::CheckCondition | `CCommonBag::CheckCondition` | 暂无 CommonBag 条件 | 未覆盖 | 使用物品或触发条件时判断 CommonBag 是否匹配。 |
| 53 | CommonBag/MonsterBag/EventBag | CommonBag::UseBag | `CCommonBag::UseBag` | 暂无 CommonBag 使用 | 未覆盖 | CommonBag 命中后生成掉落结果。 |
| 54 | CommonBag/MonsterBag/EventBag | MonsterBag::SetBagInfo | `CMonsterBag::SetBagInfo` | 暂无 MonsterBag | 未覆盖 | MonsterBag 绑定怪物 Class 和条件。 |
| 55 | CommonBag/MonsterBag/EventBag | MonsterBag::CheckCondition | `CMonsterBag::CheckCondition` | 暂无 MonsterBag 条件 | 未覆盖 | 怪物死亡时判断是否命中特定 MonsterBag。 |
| 56 | CommonBag/MonsterBag/EventBag | MonsterBag::UseBag | `CMonsterBag::UseBag` | 暂无 MonsterBag 使用 | 未覆盖 | 怪物专属 Bag 命中后生成掉落。 |
| 57 | CommonBag/MonsterBag/EventBag | EventBag::SetBagInfo | `CEventBag::SetBagInfo` | 暂无 EventBag | 未覆盖 | EventBag 绑定事件 ID 和条件。 |
| 58 | CommonBag/MonsterBag/EventBag | EventBag::CheckCondition | `CEventBag::CheckCondition` | 暂无 EventBag 条件 | 未覆盖 | 事件上下文判断是否可触发 EventBag。 |
| 59 | CommonBag/MonsterBag/EventBag | EventBag::UseBag | `CEventBag::UseBag` | 暂无 EventBag 使用 | 未覆盖 | 事件掉落包命中后生成地图掉落。 |
| 60 | CommonBag/MonsterBag/EventBag | Bag 外部出口 | `UseBag_GremoryCase/AddCashCoin` | 暂无外部出口 | 未覆盖 | 只记录 Bag 可有外部出口，具体 Gremory/现金币不归掉落核心。 |
| 61 | 特殊掉落系统 | CItemDrop::LoadFile | `CItemDrop::LoadFile` | 暂无 CItemDrop 配置 | 未覆盖 | 加载通用特殊掉落配置。 |
| 62 | 特殊掉落系统 | CItemDrop::DropItem | `CItemDrop::DropItem` | 暂无通用特殊掉落 | 未覆盖 | 按地图、怪物、玩家条件尝试特殊掉落。 |
| 63 | 特殊掉落系统 | CItemDrop::LoadZenDropFile | `CItemDrop::LoadZenDropFile` | `maps.map.go` 已加载 ZenDrop | 部分覆盖 | Go 侧有 ZenDrop 配置，需接入怪物掉落流程。 |
| 64 | 特殊掉落系统 | CItemDrop::IsZenDropActive | `CItemDrop::IsZenDropActive` | `mapManager.zen` 部分存在 | 部分覆盖 | 判断地图是否启用特殊 ZenDrop。 |
| 65 | 特殊掉落系统 | AppointItemDrop 脚本加载 | `CAppointItemDrop::LoadAppointItemDropScript` | 暂无指定掉落配置 | 未覆盖 | 加载指定怪物、地图、时间或条件掉落。 |
| 66 | 特殊掉落系统 | AppointItemDrop 执行 | `CAppointItemDrop::AppointItemDrop` | 暂无指定掉落执行 | 未覆盖 | 怪物死亡时优先执行定向掉落。 |
| 67 | 特殊掉落系统 | SetItemDrop 加载 | `CSetItemDrop::LoadFile` | `item_set.go` 只加载套装定义 | 未覆盖 | 加载套装掉落规则，而不是套装属性规则。 |
| 68 | 特殊掉落系统 | SetItemDrop 执行 | `CSetItemDrop::DropItem` | 暂无套装掉落执行 | 未覆盖 | 按概率掉落 ancient/set item。 |
| 69 | 特殊掉落系统 | Event1 当日限额掉落 | `gEvent1ItemDropToday*` | 暂无每日限额掉落 | 未覆盖 | 支持按日计数和上限控制的事件掉落。 |
| 70 | 特殊掉落系统 | MuRummy 掉落 | `g_CMuRummyMng.GetMuRummyEventItemDropRate` | 暂无 MuRummy | 未覆盖 | 事件卡牌掉落归活动触发，掉落落地规则归本模块记录。 |
| 71 | 活动/任务/Pentagram/套装掉落 | QuestInfo MonsterItemDrop | `g_QuestInfo.MonsterItemDrop` | 任务掉落未接入 | 未覆盖 | 老任务物品掉落在普通掉落前处理。 |
| 72 | 活动/任务/Pentagram/套装掉落 | QuestExp MonsterItemDrop | `g_QuestExpProgMng.QuestMonsterItemDrop` | QuestExp 掉落未接入 | 未覆盖 | 新任务物品掉落需按玩家任务状态判断。 |
| 73 | 活动/任务/Pentagram/套装掉落 | NewPVP DropItem | `g_NewPVP.DropItem` | NewPVP 未实现 | 未覆盖 | NewPVP 特殊掉落由活动/PVP 系统触发，落地归掉落系统。 |
| 74 | 活动/任务/Pentagram/套装掉落 | ArcaBattle DropItem | `g_ArcaBattle.DropItem` | Arca 未实现 | 未覆盖 | Arca 活动掉落需接入死亡掉落编排。 |
| 75 | 活动/任务/Pentagram/套装掉落 | Pentagram 属性怪掉落 | `g_PentagramSystem.AttributeMonsterItemDrop` | Pentagram 掉落未实现 | 未覆盖 | 属性怪掉落 Pentagram 相关物品。 |
| 76 | 活动/任务/Pentagram/套装掉落 | Socket Sphere 掉落 | `GetSphereDropInfo` 相关分支 | Socket 掉落未实现 | 未覆盖 | Socket Sphere/Tetra 等特殊掉落需要接入。 |
| 77 | 活动/任务/Pentagram/套装掉落 | Fenrir 材料掉落 | `gEventMonsterItemDrop` Fenrir 分支 | Fenrir 材料掉落未实现 | 未覆盖 | 炎狼兽材料按地图/怪物条件掉落。 |
| 78 | 活动/任务/Pentagram/套装掉落 | 变身戒指掉落 | `gIsItemDropRingOfTransform` | 戒指事件掉落未实现 | 未覆盖 | 变身戒指掉落按事件开关和概率控制。 |
| 79 | 活动/任务/Pentagram/套装掉落 | 活动道具掉落 | `gEventMonsterItemDrop` 多分支 | 活动道具掉落未实现 | 未覆盖 | 万圣节、樱花、活动入场物等掉落后续接入。 |
| 80 | 活动/任务/Pentagram/套装掉落 | 副本专属掉落 | EventDungeon/DoppelGanger bags | 副本掉落未实现 | 未覆盖 | 副本内怪物或结算掉落需要独立配置和落地规则。 |
| 81 | 掉落倍率与归属 | VIP 掉落加成 | `g_VipSystem.GetDropBonus` | VIP 未接入掉落 | 未覆盖 | 掉落概率应叠加 VIP 掉落加成。 |
| 82 | 掉落倍率与归属 | BonusEvent 掉落加成 | `g_BonusEvent.GetAddDrop` | BonusEvent 未实现 | 未覆盖 | 全局活动掉落加成应进入掉落概率。 |
| 83 | 掉落倍率与归属 | Gens 掉落加成 | `g_GensSystem.GetBattleZoneDropBonus` | Gens 未实现 | 未覆盖 | BattleZone 内 Gens 掉落加成需接入。 |
| 84 | 掉落倍率与归属 | Buff 掉落加成 | `gObjGetTotalValueOfEffect(EFFECTTYPE_ITEMDROPRATE)` | Buff 系统未接掉落 | 未覆盖 | 掉落率应读取 Buff/期限道具效果。 |
| 85 | 掉落倍率与归属 | 地图掉落加成 | `g_MapAttr.GetItemDropBonus` | `MapAttr.ItemDropRateBonus` 已加载 | 部分覆盖 | 地图加成配置已读，需接入掉落计算。 |
| 86 | 掉落倍率与归属 | 套装掉落加成 | `SetOpImproveItemDropRate` | 套装效果未接掉落 | 未覆盖 | 套装属性提高掉落率应接入倍率。 |
| 87 | 掉落倍率与归属 | 掉落所有者 | `iLootIndex/MaxHitUser` | 地图物品无所有者 | 未覆盖 | 物品落地应记录初始可拾取玩家。 |
| 88 | 掉落倍率与归属 | 队伍归属 | Party 最大贡献者/成员 | 组队未实现 | 未覆盖 | 组队掉落应支持队伍共享、轮流或最大伤害归属。 |
| 89 | 掉落倍率与归属 | 所有权时间 | `ItemSerialCreateSend` loot index | 暂无所有权超时 | 未覆盖 | 掉落一定时间内只允许归属者拾取，超时开放。 |
| 90 | 掉落倍率与归属 | 掉落位置归属 | `joinmuDropItemUnderCharacter` | 暂无落到玩家脚下配置 | 未覆盖 | 支持掉落在怪物位置或归属玩家脚下。 |
| 91 | 地图落物与拾取边界 | MapC.MonsterItemDrop | `MapC[].MonsterItemDrop` | `MapManager.AddItem` | 部分覆盖 | Go 侧地图落物基础存在，需补所有权和协议细节。 |
| 92 | 地图落物与拾取边界 | MapC.MoneyItemDrop | `MapC[].MoneyItemDrop` | 金币作为 item.NewItem(14,15) | 部分覆盖 | Go 侧金币表示需要确认是否符合客户端金币包体语义。 |
| 93 | 地图落物与拾取边界 | ItemSerialCreateSend | `ItemSerialCreateSend` | 暂无序列号创建服务 | 未覆盖 | 高价值物品应创建序列号并落地。 |
| 94 | 地图落物与拾取边界 | 随机落点 | `gObjGetRandomItemDropLocation` | 暂无随机落点 | 未覆盖 | 在怪物附近寻找可站立、可见、可拾取坐标。 |
| 95 | 地图落物与拾取边界 | 地图容量 | `MapItem` 容器容量 | `map_item.go` 有固定容器 | 部分覆盖 | 地图物品满时的失败策略和日志需补齐。 |
| 96 | 地图落物与拾取边界 | 地图物品过期 | `MapItem` 生命周期 | `MapManager.ExpireItem` 已有 | 部分覆盖 | 已有过期入口，需对齐不同物品类型的过期时间。 |
| 97 | 地图落物与拾取边界 | 拾取权限 | `iLootIndex` 检查 | `PickItem` 未查归属 | 未覆盖 | 拾取时校验所有者、队伍、超时和地图状态。 |
| 98 | 地图落物与拾取边界 | 拾取后移除 | `MapC.RemoveItem` | `MapManager.RemoveItem` 已有 | 部分覆盖 | 需保证并发拾取只成功一次。 |
| 99 | 地图落物与拾取边界 | 视野创建掉落 | 地图物品视野协议 | `MsgCreateViewportItemReply` 已有 | 部分覆盖 | 需确保新进入视野能看到已有掉落。 |
| 100 | 地图落物与拾取边界 | 掉落日志 | drop log | 暂无系统化掉落日志 | 未覆盖 | 记录掉落来源、归属、物品属性、地图坐标和拾取结果。 |
| 101 | 掉落结果生成与物品属性 | 掉落 type/index | `DropItem->m_Type` | `item.NewItem(section,index)` | 部分覆盖 | 基础物品生成已覆盖，需补特殊掉落的类型转换。 |
| 102 | 掉落结果生成与物品属性 | 掉落 level | `DropItem->m_Level` | `it.Level` | 部分覆盖 | 对齐普通、卓越、事件、Bag 物品等级规则。 |
| 103 | 掉落结果生成与物品属性 | durability | `ItemGetDurability`、`m_Durability` | `Item.Calc` 计算耐久 | 部分覆盖 | 掉落物品应设置正确初始耐久。 |
| 104 | 掉落结果生成与物品属性 | skill option | `Option1` | `it.Skill` | 部分覆盖 | 掉落时按概率设置技能选项。 |
| 105 | 掉落结果生成与物品属性 | luck option | `Option2` | `it.Luck` | 部分覆盖 | 掉落时按概率设置幸运。 |
| 106 | 掉落结果生成与物品属性 | add option | `Option3` | `it.Option`/追加字段 | 部分覆盖 | 追加属性随机规则需对齐 GameServer。 |
| 107 | 掉落结果生成与物品属性 | excellent option | `NewOption` | `Excellent` 字段组 | 部分覆盖 | 卓越 bit 生成后需正确序列化和数值生效。 |
| 108 | 掉落结果生成与物品属性 | ancient/set option | `SetOption` | `it.Set` | 部分覆盖 | 套装掉落应写入 set/ancient 标记。 |
| 109 | 掉落结果生成与物品属性 | socket/380 option | `MakeSocketSlot`、`Is380Item` | Socket 未完整、380 字段存在 | 未覆盖 | 掉落可带 Socket 槽或 380 属性时需生成对应字段。 |
| 110 | 掉落结果生成与物品属性 | period duration | `DurationItem`、期限物品 | 暂无掉落期限物品 | 未覆盖 | Bag 或特殊掉落可生成期限物品。 |
| 111 | server-game 落地点/测试边界 | 独立 drop 包边界 | GameServer 分散模块 | 当前在 `object/monster` | 设计调整 | 后续可抽成 `game/drop`，对象系统只调用掉落服务。 |
| 112 | server-game 落地点/测试边界 | Monster.DieDropItem 精简 | `gObjMonsterDieGiveItem` 臃肿 | `Monster.DieDropItem` 已开始膨胀 | 部分覆盖 | 保留死亡入口，具体规则委托掉落系统。 |
| 113 | server-game 落地点/测试边界 | DropContext | `lpObj/lpTargetObj/MaxHitUser` | 暂无上下文结构 | 未覆盖 | 定义怪物、归属玩家、地图、事件、队伍、倍率等上下文。 |
| 114 | server-game 落地点/测试边界 | DropResult | `CItem/ItemSerialCreateSend` | 暂无统一结果结构 | 未覆盖 | 定义地图物品、金币、短路、失败原因等结果。 |
| 115 | server-game 落地点/测试边界 | 掉落顺序测试 | `gObjMonsterDieGiveItem` 顺序 | 暂无测试 | 未覆盖 | 覆盖任务、Bag、特殊、普通、金币的优先级。 |
| 116 | server-game 落地点/测试边界 | 基础掉落回归 | `MonsterItemMng` | `DropManager` | 部分覆盖 | 覆盖不同怪物等级、无候选、宝石、魔法书、普通物品。 |
| 117 | server-game 落地点/测试边界 | 卓越掉落回归 | `GetItemExcel` | `DropItemExcellent` | 部分覆盖 | 覆盖卓越候选、卓越数量、卓越 bit 和 skill/luck。 |
| 118 | server-game 落地点/测试边界 | 金币掉落回归 | `MoneyItemDrop` | `DieDropItem` 金币分支 | 部分覆盖 | 覆盖金币概率、倍率、地图落物和拾取。 |
| 119 | server-game 落地点/测试边界 | 归属与拾取回归 | `iLootIndex` | 暂无归属 | 未覆盖 | 覆盖归属者、非归属者、队伍成员、超时开放。 |
| 120 | server-game 落地点/测试边界 | 地图落物生命周期回归 | `MapItem` | `MapManager.AddItem/ExpireItem/RemoveItem` | 部分覆盖 | 覆盖掉落创建、视野显示、拾取、过期、并发拾取。 |
