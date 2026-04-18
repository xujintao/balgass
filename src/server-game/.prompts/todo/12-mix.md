# 12. 合成系统

本模块覆盖 `GameServer` 中 ChaosBox 主合成、宝石组合/拆分、翅膀/披风/高级翅膀、事件物品、宠物/坐骑、Harmony、380、Socket、Pentagram、概率费用与合成事务。当前 `server-game` 已有部分协议入口、NPC 类型、配置结构和道具字段，但核心合成流程、材料校验、事务处理与结果生成仍基本缺失。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 合成入口/ChaosBox 状态 | GCUserChaosBoxSend | `GameProtocol::GCUserChaosBoxSend` | 暂无 ChaosBox 物品列表响应 | 未覆盖 | 实现 ChaosBox 容器下发，包含窗口类型、格子物品序列化和客户端打开状态。 |
| 2 | 合成入口/ChaosBox 状态 | CGChaosBoxItemMixButtonClick | `GameProtocol::CGChaosBoxItemMixButtonClick` | `handle/c1c2.go` 有 `0x86: chaosBoxItemMixButtonOK` | 部分覆盖 | 解析合成按钮请求，校验玩家状态、ChaosBox 状态、混合类型和材料。 |
| 3 | 合成入口/ChaosBox 状态 | CGChaosBoxUseEnd | `GameProtocol::CGChaosBoxUseEnd` | `handle/c1c2.go` 有 `0x87: chaosBoxUseEnd` | 部分覆盖 | 关闭 ChaosBox 时清理合成状态、解锁容器并同步物品。 |
| 4 | 合成入口/ChaosBox 状态 | CheckEmptySpace_MultiMix | `CMixSystem::CheckEmptySpace_MultiMix` | `handle/c1c2.go` 有 `0x88: checkMultiMix` | 部分覆盖 | 多合成前检查背包可用空间，避免成功后无处放置结果。 |
| 5 | 合成入口/ChaosBox 状态 | ChaosBoxCheck | `CMixSystem::ChaosBoxCheck` | 暂无 ChaosBox 状态检查 | 未覆盖 | 判断玩家是否可使用 ChaosBox，排除交易、商店、仓库、个人商店等冲突状态。 |
| 6 | 合成入口/ChaosBox 状态 | ChaosBoxInit | `CMixSystem::ChaosBoxInit` | `object.Object` 中 ChaosBox 字段仍是注释 | 未覆盖 | 初始化玩家 ChaosBox 容器、费用、成功率、锁状态和多合成状态。 |
| 7 | 合成入口/ChaosBox 状态 | ChaosBoxItemDown | `CMixSystem::ChaosBoxItemDown` | 暂无 ChaosBox 材料落回背包 | 未覆盖 | 关闭或失败回滚时将 ChaosBox 中的材料安全移回背包。 |
| 8 | 合成入口/ChaosBox 状态 | ChaosBoxMix | `CMixSystem::ChaosBoxMix` | 暂无合成主判定 | 未覆盖 | 对 ChaosBox 材料进行总判定，识别合成类型、基础成功率、费用和目标结果。 |
| 9 | 合成入口/ChaosBox 状态 | ChaosMixCharmItemUsed | `CMixSystem::ChaosMixCharmItemUsed` | 暂无合成符咒消耗 | 未覆盖 | 处理合成符咒或概率增益道具的识别、消耗和成功率加成。 |
| 10 | 合成入口/ChaosBox 状态 | bIsChaosMixCompleted/ChaosLock 状态 | `OBJECTSTRUCT::bIsChaosMixCompleted`、`ChaosLock` | `object.Object` 中对应字段为注释 | 未覆盖 | 建立合成进行中状态，防止重复点击、移动材料和并发修改背包。 |
| 11 | ChaosBox 主分发 | DefaultChaosMix | `CMixSystem::DefaultChaosMix` | 暂无默认 Chaos 合成 | 未覆盖 | 实现普通 Chaos 合成主流程，覆盖玛雅武器等基础合成。 |
| 12 | ChaosBox 主分发 | IsMixPossibleItem | `CMixSystem::IsMixPossibleItem` | 暂无材料可合成判断 | 未覆盖 | 判定材料是否允许放入合成逻辑，过滤锁定、不可交易、任务、过期等物品。 |
| 13 | ChaosBox 主分发 | IsDeleteItem | `CMixSystem::IsDeleteItem` | 暂无失败删除规则 | 未覆盖 | 定义失败时哪些材料删除、保留、降级或耐久变化。 |
| 14 | ChaosBox 主分发 | LogChaosItem | `CMixSystem::LogChaosItem` | 暂无合成日志 | 未覆盖 | 记录合成材料、结果、概率、随机数、费用和玩家信息。 |
| 15 | ChaosBox 主分发 | LogPlusItemLevelChaosItem | `CMixSystem::LogPlusItemLevelChaosItem` | 暂无强化合成日志 | 未覆盖 | 为 +10~+15 强化记录专用日志，便于追踪高价值物品变化。 |
| 16 | ChaosBox 主分发 | LogDQChaosItem | `CMixSystem::LogDQChaosItem` | 暂无事件材料合成日志 | 未覆盖 | 记录恶魔广场等事件入场材料合成结果。 |
| 17 | ChaosBox 主分发 | PlusChaosRate 查询 | `GameProtocol::CGReqPlusChaosRate` | `handle/c1c2.go` 有 `0xBD09: reqPlusChaosRate` | 部分覆盖 | 返回 Crywolf 等系统提供的 Chaos 成功率加成。 |
| 18 | ChaosBox 主分发 | 合成类型识别 | `GameProtocol::CGChaosBoxItemMixButtonClick` 分发逻辑 | 暂无统一 mix type 枚举 | 未覆盖 | 将客户端 mix type 映射到具体 Go 合成处理器。 |
| 19 | ChaosBox 主分发 | 合成结果响应 | `PMSG_CHAOSMIX` 相关响应 | 暂无 ChaosBox 结果包体 | 未覆盖 | 统一成功、失败、材料不足、钱不足、空间不足等结果响应。 |
| 20 | ChaosBox 主分发 | 合成失败材料处理 | `ChaosBox.cpp` 各失败分支 | 暂无失败事务策略 | 未覆盖 | 按不同合成类型处理失败后的材料删除、降级、保留和容器同步。 |
| 21 | 装备强化 +10~+15 | PlusItemLevelChaosMix +10 | `CMixSystem::PlusItemLevelChaosMix(CHAOS_TYPE_UPGRADE_10)` | ChaosBox 强化未实现 | 未覆盖 | 实现 +10 强化材料校验、费用、成功率和结果更新。 |
| 22 | 装备强化 +10~+15 | PlusItemLevelChaosMix +11 | `CMixSystem::PlusItemLevelChaosMix(CHAOS_TYPE_UPGRADE_11)` | ChaosBox 强化未实现 | 未覆盖 | 实现 +11 强化，并继承 +10 后的材料与失败规则。 |
| 23 | 装备强化 +10~+15 | PlusItemLevelChaosMix +12 | `CMixSystem::PlusItemLevelChaosMix(CHAOS_TYPE_UPGRADE_12)` | ChaosBox 强化未实现 | 未覆盖 | 实现 +12 强化，确保目标装备等级与材料数量匹配。 |
| 24 | 装备强化 +10~+15 | PlusItemLevelChaosMix +13 | `CMixSystem::PlusItemLevelChaosMix(CHAOS_TYPE_UPGRADE_13)` | ChaosBox 强化未实现 | 未覆盖 | 实现 +13 强化，包含失败降级或销毁策略。 |
| 25 | 装备强化 +10~+15 | PlusItemLevelChaosMix +14 | `CMixSystem::PlusItemLevelChaosMix(CHAOS_TYPE_UPGRADE_14)` | ChaosBox 强化未实现 | 未覆盖 | 实现 +14 强化并读取对应成功率配置。 |
| 26 | 装备强化 +10~+15 | PlusItemLevelChaosMix +15 | `CMixSystem::PlusItemLevelChaosMix(CHAOS_TYPE_UPGRADE_15)` | ChaosBox 强化未实现 | 未覆盖 | 实现 +15 强化、公告开关和高等级失败处理。 |
| 27 | 装备强化 +10~+15 | Level10-Level15 成功率配置 | `IGC_ChaosBox.xml`、`RateChaosBoxMixs` | `conf.configChaosBox` 已读取部分字段 | 部分覆盖 | 使用 `Level10` 到 `Level15` 配置驱动强化成功率。 |
| 28 | 装备强化 +10~+15 | AddLuck 幸运加成 | `rateChaosBoxMix::AddLuck` | `conf.rateChaosBoxMix.AddLuck` 已存在 | 部分覆盖 | 装备幸运属性应按配置增加合成成功率，并受上限限制。 |
| 29 | 装备强化 +10~+15 | Level15Notice | `EnableLevel15Notice` | `conf.rateChaosBoxMix.EnableLevel15Notice` 已存在 | 部分覆盖 | +15 成功时按配置向玩家或全服广播。 |
| 30 | 装备强化 +10~+15 | 强化失败降级/销毁规则 | `PlusItemLevelChaosMix` 失败分支 | 暂无强化失败规则 | 未覆盖 | 明确失败时目标装备等级、Option、Luck、Excellent、Socket 等字段如何变化。 |
| 31 | 翅膀/披风/高级翅膀 | CheckWingItem | `CMixSystem::CheckWingItem` | `item.KindBWing1st` 等分类存在 | 部分覆盖 | 判断一代翅膀材料是否合法，用于翅膀合成和高级翅膀前置。 |
| 32 | 翅膀/披风/高级翅膀 | Check2ndWingItem | `CMixSystem::Check2ndWingItem` | `item.KindBWing2nd` 分类存在 | 部分覆盖 | 判断二代翅膀材料是否合法。 |
| 33 | 翅膀/披风/高级翅膀 | Check3rdWingItem | `CMixSystem::Check3rdWingItem` | `item.KindBWing3rd` 分类存在 | 部分覆盖 | 判断三代翅膀材料是否合法。 |
| 34 | 翅膀/披风/高级翅膀 | WingChaosMix | `CMixSystem::WingChaosMix` | 暂无一代/二代翅膀合成 | 未覆盖 | 实现基础翅膀合成，包含材料识别、职业结果、随机属性和失败处理。 |
| 35 | 翅膀/披风/高级翅膀 | ThirdWingLevel1ChaosMix | `CMixSystem::ThirdWingLevel1ChaosMix` | 暂无三代翅膀阶段一 | 未覆盖 | 合成神鹰之羽等三代翅膀中间材料。 |
| 36 | 翅膀/披风/高级翅膀 | ThirdWingLevel2ChaosMix | `CMixSystem::ThirdWingLevel2ChaosMix` | 暂无三代翅膀阶段二 | 未覆盖 | 使用中间材料合成三代翅膀或披风。 |
| 37 | 翅膀/披风/高级翅膀 | AdvancedWingMix | `CMixSystem::AdvancedWingMix` | 暂无高级翅膀合成 | 未覆盖 | 实现高级翅膀或特殊翅膀合成规则。 |
| 38 | 翅膀/披风/高级翅膀 | ThirdWingMixFail | `CMixSystem::ThirdWingMixFail` | 暂无三代翅膀失败处理 | 未覆盖 | 三代翅膀失败时按 GameServer 规则处理材料和目标物品。 |
| 39 | 翅膀/披风/高级翅膀 | ThirdWingMixFailItemPanalty | `CMixSystem::ThirdWingMixFailItemPanalty` | 暂无失败惩罚函数 | 未覆盖 | 对失败后保留的材料执行等级、Option、Durability 等惩罚。 |
| 40 | 翅膀/披风/高级翅膀 | MonsterWingMix | `CMixSystem::MonsterWingMix` | `item.KindBWingMonster` 分类存在 | 部分覆盖 | 实现 2.5 代或怪物翅膀合成与结果随机。 |
| 41 | 事件/宠物/坐骑合成 | DevilSquareItemChaosMix | `CMixSystem::DevilSquareItemChaosMix` | 暂无恶魔广场材料合成 | 未覆盖 | 合成恶魔广场入场券，校验眼睛、钥匙、玛雅等材料等级。 |
| 42 | 事件/宠物/坐骑合成 | DevilSquareItemChaosMix_Multi | `CMixSystem::DevilSquareItemChaosMix_Multi` | 多合成入口存在但无业务 | 未覆盖 | 支持批量合成恶魔广场入场券并检查背包空间。 |
| 43 | 事件/宠物/坐骑合成 | BloodCastleItemChaosMix | `CMixSystem::BloodCastleItemChaosMix` | 暂无血色城堡材料合成 | 未覆盖 | 合成血色城堡入场券，校验材料等级与数量。 |
| 44 | 事件/宠物/坐骑合成 | BloodCastleItemChaosMix_Multi | `CMixSystem::BloodCastleItemChaosMix_Multi` | 多合成入口存在但无业务 | 未覆盖 | 支持批量合成血色城堡入场券。 |
| 45 | 事件/宠物/坐骑合成 | DarkHorseChaosMix | `CMixSystem::DarkHorseChaosMix` | 暂无黑王马合成 | 未覆盖 | 实现黑王马材料检查、宠物生成、等级经验字段初始化。 |
| 46 | 事件/宠物/坐骑合成 | DarkSpiritChaosMix | `CMixSystem::DarkSpiritChaosMix` | 暂无黑王鸟合成 | 未覆盖 | 实现黑王鸟材料检查和宠物道具生成。 |
| 47 | 事件/宠物/坐骑合成 | PegasiaChaosMix | `CMixSystem::PegasiaChaosMix` | 暂无独角兽/特殊坐骑合成 | 未覆盖 | 实现 Pegasia 类坐骑或事件道具合成。 |
| 48 | 事件/宠物/坐骑合成 | PegasiaChaosMix_Multi | `CMixSystem::PegasiaChaosMix_Multi` | 多合成入口存在但无业务 | 未覆盖 | 支持 Pegasia 类道具批量合成。 |
| 49 | 事件/宠物/坐骑合成 | Fenrir_01Level_Mix/Fenrir_02Level_Mix | `CMixSystem::Fenrir_01Level_Mix`、`Fenrir_02Level_Mix` | `conf.Common` 有 Fenrir 成功率字段 | 部分覆盖 | 实现炎狼兽低阶材料到中间材料的分阶段合成。 |
| 50 | 事件/宠物/坐骑合成 | Fenrir_03Level_Mix/Fenrir_04Upgrade_Mix | `CMixSystem::Fenrir_03Level_Mix`、`Fenrir_04Upgrade_Mix` | `conf.Common` 有 Fenrir 成功率字段 | 部分覆盖 | 实现炎狼兽最终合成和升级强化。 |
| 51 | 药水/果实/特殊箱子 | BlessPotionChaosMix | `CMixSystem::BlessPotionChaosMix` | 暂无祝福药水合成 | 未覆盖 | 合成攻城或事件用祝福药水。 |
| 52 | 药水/果实/特殊箱子 | BlessPotionChaosMix_Multi | `CMixSystem::BlessPotionChaosMix_Multi` | 多合成入口存在但无业务 | 未覆盖 | 支持批量合成祝福药水。 |
| 53 | 药水/果实/特殊箱子 | SoulPotionChaosMix | `CMixSystem::SoulPotionChaosMix` | 暂无灵魂药水合成 | 未覆盖 | 合成攻城或事件用灵魂药水。 |
| 54 | 药水/果实/特殊箱子 | SoulPotionChaosMix_Multi | `CMixSystem::SoulPotionChaosMix_Multi` | 多合成入口存在但无业务 | 未覆盖 | 支持批量合成灵魂药水。 |
| 55 | 药水/果实/特殊箱子 | ShieldPotionLv1_Mix | `CMixSystem::ShieldPotionLv1_Mix` | `conf.Common` 有 SDPotion1 配置 | 部分覆盖 | 实现一级防护药水合成。 |
| 56 | 药水/果实/特殊箱子 | ShieldPotionLv2_Mix | `CMixSystem::ShieldPotionLv2_Mix` | `conf.Common` 有 SDPotion2 配置 | 部分覆盖 | 实现二级防护药水合成。 |
| 57 | 药水/果实/特殊箱子 | ShieldPotionLv3_Mix | `CMixSystem::ShieldPotionLv3_Mix` | `conf.Common` 有 SDPotion3 配置 | 部分覆盖 | 实现三级防护药水合成。 |
| 58 | 药水/果实/特殊箱子 | CircleChaosMix | `CMixSystem::CircleChaosMix` | 暂无果实合成 | 未覆盖 | 实现果实类道具合成，并支持多合成版本。 |
| 59 | 药水/果实/特殊箱子 | LotteryItemMix/PremiumBoxMix | `CMixSystem::LotteryItemMix`、`PremiumBoxMix` | 暂无抽奖/高级箱子合成 | 未覆盖 | 实现特殊箱子或抽奖道具的材料合成与结果生成。 |
| 60 | 药水/果实/特殊箱子 | CherryBlossomMix/HiddenTreasureBoxItemMix | `CMixSystem::CherryBlossomMix`、`HiddenTreasureBoxItemMix` | `conf.ChaosBox.CherryBlossom` 已有配置 | 部分覆盖 | 实现樱花、隐藏宝箱等活动材料合成。 |
| 61 | 宝石组合/拆分 | CGReqJewelMix | `GameProtocol::CGReqJewelMix` | `handle/c1c2.go` 有 `0xBC00: jewelMix` | 部分覆盖 | 实现宝石组合请求解析和参数校验。 |
| 62 | 宝石组合/拆分 | GCAnsJewelMix | `GameProtocol::GCAnsJewelMix` | 暂无宝石组合响应 | 未覆盖 | 返回宝石组合成功、失败、钱不足、材料不足、空间不足等结果。 |
| 63 | 宝石组合/拆分 | CGReqJewelUnMix | `GameProtocol::CGReqJewelUnMix` | `handle/c1c2.go` 有 `0xBC01: jewelUnmix` | 部分覆盖 | 实现宝石拆分请求，校验组合宝石类型、等级和位置。 |
| 64 | 宝石组合/拆分 | GCAnsJewelUnMix | `GameProtocol::GCAnsJewelUnMix` | 暂无宝石拆分响应 | 未覆盖 | 返回宝石拆分结果并同步背包。 |
| 65 | 宝石组合/拆分 | GetJewelCount | `CJewelMixSystem::GetJewelCount` | 背包有道具结构但无宝石计数封装 | 未覆盖 | 统计背包中指定类型散宝石数量。 |
| 66 | 宝石组合/拆分 | LoadMixJewelPrice | `CJewelMixSystem::LoadMixJewelPrice` | `conf.Price` 有宝石价格字段 | 部分覆盖 | 加载或复用宝石组合/拆分费用配置。 |
| 67 | 宝石组合/拆分 | GetJewelCountPerLevel | `CJewelMixSystem::GetJewelCountPerLevel` | 暂无组合等级到数量映射 | 未覆盖 | 定义 10/20/30 等组合宝石对应散宝石数量。 |
| 68 | 宝石组合/拆分 | MixJewel | `CJewelMixSystem::MixJewel` | 暂无宝石组合业务 | 未覆盖 | 扣除散宝石和费用，生成指定等级的组合宝石。 |
| 69 | 宝石组合/拆分 | UnMixJewel | `CJewelMixSystem::UnMixJewel` | 暂无宝石拆分业务 | 未覆盖 | 删除组合宝石，生成对应数量散宝石并检查背包空间。 |
| 70 | 宝石组合/拆分 | 宝石组合背包空间/费用校验 | `CJewelMixSystem` 组合/拆分分支 | 背包空间工具存在但未被合成使用 | 未覆盖 | 合成前统一检查材料、钱、目标空格，确保事务可完成。 |
| 71 | Harmony 再生系统 | LoadScript | `CJewelOfHarmonySystem::LoadScript` | `game/item/item_harmoney.go` 已加载 Harmony XML | 部分覆盖 | 对齐 Harmony 属性配置解析，确认武器、法杖、防具分类和权重。 |
| 72 | Harmony 再生系统 | LoadScriptOfSmelt | `CJewelOfHarmonySystem::LoadScriptOfSmelt` | `item_harmoney.go` 已加载 Smelt XML | 部分覆盖 | 对齐进化石材料配置解析和可熔炼物品列表。 |
| 73 | Harmony 再生系统 | PurityJewelOfHarmony | `CJewelOfHarmonySystem::PurityJewelOfHarmony` | 暂无再生原石提炼业务 | 未覆盖 | 将再生原石提炼为再生宝石，处理成功率、费用和失败。 |
| 74 | Harmony 再生系统 | PurityJewelOfHarmony_MultiMix | `CJewelOfHarmonySystem::PurityJewelOfHarmony_MultiMix` | 暂无 Harmony 多合成 | 未覆盖 | 支持批量提炼再生宝石。 |
| 75 | Harmony 再生系统 | MakeSmeltingStoneItem | `CJewelOfHarmonySystem::MakeSmeltingStoneItem` | 暂无进化宝石合成 | 未覆盖 | 使用装备或材料合成低级/高级进化宝石。 |
| 76 | Harmony 再生系统 | MakeSmeltingStoneItem_MultiMix | `CJewelOfHarmonySystem::MakeSmeltingStoneItem_MultiMix` | 暂无进化宝石多合成 | 未覆盖 | 支持批量合成进化宝石。 |
| 77 | Harmony 再生系统 | StrengthenItemByJewelOfHarmony | `CJewelOfHarmonySystem::StrengthenItemByJewelOfHarmony` | `HarmonyManager.StrengthenItem` 有雏形 | 部分覆盖 | 使用再生宝石给装备添加 Harmony 属性，并同步物品与角色属性。 |
| 78 | Harmony 再生系统 | StrengthenItemByJewelOfRise | `CJewelOfHarmonySystem::StrengthenItemByJewelOfRise` | 暂无提高宝石强化 | 未覆盖 | 实现提高宝石或强化等级提升逻辑。 |
| 79 | Harmony 再生系统 | SmeltItemBySmeltingStone | `CJewelOfHarmonySystem::SmeltItemBySmeltingStone` | 暂无进化宝石熔炼装备 | 未覆盖 | 使用进化宝石提升 Harmony 属性等级或效果。 |
| 80 | Harmony 再生系统 | RestoreStrengthenItem | `CJewelOfHarmonySystem::RestoreStrengthenItem` | 暂无 Harmony 还原 | 未覆盖 | 移除装备 Harmony 属性并按配置扣除费用。 |
| 81 | Harmony 效果/限制 | GetItemStrengthenOption | `CJewelOfHarmonySystem::GetItemStrengthenOption` | `item.Item` 有 `HarmonyEffect` 字段 | 部分覆盖 | 提供装备 Harmony 属性类型读取接口。 |
| 82 | Harmony 效果/限制 | GetItemOptionLevel | `CJewelOfHarmonySystem::GetItemOptionLevel` | `item.Item` 有 `HarmonyLevel` 字段 | 部分覆盖 | 提供装备 Harmony 属性等级读取接口。 |
| 83 | Harmony 效果/限制 | IsStrengthenByJewelOfHarmony | `CJewelOfHarmonySystem::IsStrengthenByJewelOfHarmony` | 可通过 `HarmonyEffect` 判断 | 部分覆盖 | 判断装备是否已有 Harmony 强化，防止重复强化。 |
| 84 | Harmony 效果/限制 | IsActive | `CJewelOfHarmonySystem::IsActive` | 暂无 Harmony 生效检查 | 未覆盖 | 判断装备 Harmony 属性是否满足等级、职业、装备状态等生效条件。 |
| 85 | Harmony 效果/限制 | _GetSelectRandomOption | `CJewelOfHarmonySystem::_GetSelectRandomOption` | `addRandEffect` 已有权重随机雏形 | 部分覆盖 | 对齐 GameServer 的随机属性选择和等级要求。 |
| 86 | Harmony 效果/限制 | _MakeOption | `CJewelOfHarmonySystem::_MakeOption` | 暂无统一 Harmony 写入函数 | 未覆盖 | 封装 Harmony 属性类型和等级写入，统一校验与持久化。 |
| 87 | Harmony 效果/限制 | SetApplyStrengthenItem | `CJewelOfHarmonySystem::SetApplyStrengthenItem` | `Player.calc` 中 Harmony 效果仍未完整启用 | 未覆盖 | 装备变化时应用 Harmony 属性到玩家数值。 |
| 88 | Harmony 效果/限制 | GetItemEffectValue | `CJewelOfHarmonySystem::GetItemEffectValue` | 暂无 Harmony 效果数值查询 | 未覆盖 | 按属性类型和等级查询具体加成数值。 |
| 89 | Harmony 效果/限制 | _CalcItemEffectValue | `CJewelOfHarmonySystem::_CalcItemEffectValue` | `HarmonyEffect` 结构存在但未完整计算 | 部分覆盖 | 将 Harmony 属性转换为攻击、防御、HP、AG、SD 等角色加成。 |
| 90 | Harmony 效果/限制 | NpcJewelOfHarmony/IsEnableToTrade | `NpcJewelOfHarmony`、`IsEnableToTrade` | NPC 分支与交易限制未接入 | 未覆盖 | 接入 Harmony NPC 入口，并限制 Harmony 物品交易规则。 |
| 91 | 380 强化系统 | Load380ItemOptionInfo | `CItemSystemFor380::Load380ItemOptionInfo` | `item_380.go` 已加载 `IGC_Item380Option.xml` | 部分覆盖 | 对齐 380 道具配置加载和效果字段。 |
| 92 | 380 强化系统 | Is380Item | `CItemSystemFor380::Is380Item` | `Item380Manager.Is380Item` 已实现 | 部分覆盖 | 复核 380 可强化物品判断与 GameServer 一致。 |
| 93 | 380 强化系统 | Is380OptionItem | `CItemSystemFor380::Is380OptionItem` | `item.Item.Option380` 字段存在 | 部分覆盖 | 判断装备是否已拥有 380 属性。 |
| 94 | 380 强化系统 | InitEffectValue | `CItemSystemFor380::InitEffectValue` | `Item380Effect` 结构存在 | 部分覆盖 | 初始化 380 效果容器，避免重复累加旧值。 |
| 95 | 380 强化系统 | ApplyFor380Option | `CItemSystemFor380::ApplyFor380Option` | `Player.calc380Item` 已应用部分效果 | 部分覆盖 | 装备变化时应用 380 属性到角色 PVP、HP、SD 等数值。 |
| 96 | 380 强化系统 | _CalcItemEffectValue | `CItemSystemFor380::_CalcItemEffectValue` | `Apply380ItemEffect` 已有映射 | 部分覆盖 | 对齐 380 属性类型到 Go 效果字段的映射。 |
| 97 | 380 强化系统 | _SetOption | `CItemSystemFor380::_SetOption` | 暂无 380 选项写入业务 | 未覆盖 | 合成成功后写入或移除 `Option380` 标记。 |
| 98 | 380 强化系统 | SetOptionItemByMacro | `CItemSystemFor380::SetOptionItemByMacro` | 暂无 GM/宏设置入口 | 未覆盖 | 为调试或脚本提供 380 属性直接设置能力。 |
| 99 | 380 强化系统 | ChaosMix380ItemOption | `CItemSystemFor380::ChaosMix380ItemOption` | 暂无 380 ChaosBox 合成 | 未覆盖 | 实现 380 强化材料检查、概率、费用、成功写入和失败消耗。 |
| 100 | 380 强化系统 | 380 材料/概率/费用配置 | `Item380Option.xml` Mix 节点 | `Item380Manager.mix` 已读取 | 部分覆盖 | 使用 `JewelOfHarmonyCount`、`JewelOfGuardianCount`、`NeedZen`、`Rate` 驱动合成。 |
| 101 | Socket 系统 | LoadScript | `CItemSocketOptionSystem::LoadScript` | 暂无 Socket 系统加载器 | 未覆盖 | 加载 Socket 选项、Seed、Sphere、Slot rate 等配置。 |
| 102 | Socket 系统 | LoadOptionScript | `CItemSocketOptionSystem::LoadOptionScript` | 暂无 Socket option 配置 | 未覆盖 | 解析 Seed/Sphere 效果类型、数值类型和属性值。 |
| 103 | Socket 系统 | LoadSocketSlotRateFile | `CItemSocketOptionSystem::LoadSocketSlotRateFile` | 暂无 Socket 槽概率配置 | 未覆盖 | 加载装备生成 Socket 槽数量的概率表。 |
| 104 | Socket 系统 | IsEnableSocketItem | `CItemSocketOptionSystem::IsEnableSocketItem` | `item.TypeSocket` 判断存在但不完整 | 部分覆盖 | 判断装备是否允许拥有 Socket 槽。 |
| 105 | Socket 系统 | GetEmptySlotCount | `CItemSocketOptionSystem::GetEmptySlotCount` | 暂无 Socket 槽数据模型 | 未覆盖 | 统计装备空 Socket 槽数量。 |
| 106 | Socket 系统 | SeedExtractMix | `CMixSystem::SeedExtractMix` | 暂无 Seed 提取合成 | 未覆盖 | 从装备或材料中提取 Seed，处理材料消耗和随机结果。 |
| 107 | Socket 系统 | SeedSphereCompositeMix | `CMixSystem::SeedSphereCompositeMix` | 暂无 SeedSphere 合成 | 未覆盖 | 将 Seed 与 Sphere 合成为 SeedSphere。 |
| 108 | Socket 系统 | SetSeedSphereMix | `CMixSystem::SetSeedSphereMix` | 暂无 Socket 镶嵌 | 未覆盖 | 将 SeedSphere 镶嵌到目标装备指定槽位。 |
| 109 | Socket 系统 | SeedSphereRemoveMix | `CMixSystem::SeedSphereRemoveMix` | 暂无 Socket 拆除 | 未覆盖 | 从装备指定槽移除 Socket 属性，并处理费用和结果。 |
| 110 | Socket 系统 | ApplySeedSphereEffect/SetApplySocketEffect | `ApplySeedSphereEffect`、`SetApplySocketEffect` | 玩家数值未接入 Socket 效果 | 未覆盖 | 装备变化时应用 Socket 效果、Bonus Option 和 Set Option。 |
| 111 | Pentagram 与 Go 落地点 | LoadMixNeedSourceScript | `CPentagramMixSystem::LoadMixNeedSourceScript` | 暂无 Pentagram 合成材料配置 | 未覆盖 | 加载 Pentagram 精炼和升级所需材料。 |
| 112 | Pentagram 与 Go 落地点 | LoadJewelOptionScript | `CPentagramMixSystem::LoadJewelOptionScript` | 暂无 Errtel 属性配置加载 | 未覆盖 | 加载 Pentagram Jewel/Errtel 选项配置。 |
| 113 | Pentagram 与 Go 落地点 | PentagramMixBoxInit | `CPentagramMixSystem::PentagramMixBoxInit` | `NpcTypePentagramMix` 分支为空 | 未覆盖 | 打开 Pentagram 合成窗口并初始化相关容器状态。 |
| 114 | Pentagram 与 Go 落地点 | PentagramJewelRefine | `CPentagramMixSystem::PentagramJewelRefine` | `handle/c1c2.go` 有 `0xEC02: reqRefinePentagramJewel` | 部分覆盖 | 实现 Errtel 精炼请求、材料检查、成功率和结果更新。 |
| 115 | Pentagram 与 Go 落地点 | PentagramJewel_Upgrade | `CPentagramMixSystem::PentagramJewel_Upgrade` | `handle/c1c2.go` 有 `0xEC03: reqUpgradePentagramJewel` | 部分覆盖 | 实现 Errtel 升级请求，处理等级、目标值、材料和失败规则。 |
| 116 | Pentagram 与 Go 落地点 | PentagramJewel_IN | `g_PentagramSystem.PentagramJewel_IN` | `handle/c1c2.go` 有 `0xEC00: inJewelPentagramItem` | 部分覆盖 | 实现 Errtel 镶嵌到 Pentagram 指定槽位。 |
| 117 | Pentagram 与 Go 落地点 | PentagramJewel_OUT | `g_PentagramSystem.PentagramJewel_OUT` | `handle/c1c2.go` 有 `0xEC01: outJewelPentagramItem` | 部分覆盖 | 实现 Errtel 从 Pentagram 槽位拆除并返回背包。 |
| 118 | Pentagram 与 Go 落地点 | handle/c1c2.go 合成 opcode 入口 | `protocol.cpp` 中 `0x86/0x87/0x88/0xBC/0xBD/0xEC` 分发 | 已有 handler 映射，业务多为空 | 部分覆盖 | 补齐合成相关请求/响应结构体，并将 handler 委托到合成服务。 |
| 119 | Pentagram 与 Go 落地点 | Player.Talk ChaosMix/PentagramMix NPC 分支 | ChaosBox NPC 与 Pentagram NPC 打开逻辑 | `Player.Talk` 中两个分支为空 | 未覆盖 | NPC 对话时打开对应合成窗口，设置玩家接口状态和目标 NPC。 |
| 120 | Pentagram 与 Go 落地点 | 合成事务、并发保护、回归测试 | `ChaosLock`、ChaosBox 容器、各 mix 函数 | Go 玩家单协程可承载事务，但未实现合成事务 | 未覆盖 | 合成必须在玩家动作协程内原子执行，覆盖材料锁定、扣费、结果生成、失败回滚、断线恢复和端到端测试。 |
