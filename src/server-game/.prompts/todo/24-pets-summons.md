# 24. 宠物与召唤系统

本模块覆盖玩家附属实体和附属效果：Helper 宠物、坐骑、DarkHorse、DarkSpirit、Fenrir、Muun、变身戒指、召唤技能和召唤物。宠物与召唤系统不拥有基础物品、背包、公式、Buff 生命周期、合成和普通怪物 AI；它负责宠物启用状态、宠物经验/等级/耐久、宠物命令、Muun 背包/效果/进化/兑换、召唤物创建和清理，再调用对象、道具、公式、经验、技能、Buff、合成等系统完成底层操作。宠物/坐骑合成归 `12-mix.md`，普通 Buff 生命周期归 `13-buffs.md`，普通怪物 AI 归 `25-monster-ai.md`。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 模块边界与总入口 | PetSummonManager 总管理器 | `MuunSystem`、`DarkSpirit`、`Guardian` 等分散模块 | 暂无 `game/pets` 或 `game/summon` | 未覆盖 | 建立统一宠物与召唤服务，承接对象、道具、技能、经验事件。 |
| 2 | 模块边界与总入口 | 宠物类型枚举 | Helper/DarkHorse/DarkSpirit/Fenrir/Muun/Summon | 道具分类和 Code 判断分散 | 部分覆盖 | 定义内部宠物、坐骑、战斗宠物、Muun、召唤物类型。 |
| 3 | 模块边界与总入口 | 玩家宠物状态容器 | `OBJECTSTRUCT` 宠物字段、Muun period data | `Player.pet`、注释字段、外观字段 | 部分覆盖 | 统一玩家当前 helper、ride pet、DarkSpirit、Muun、召唤物状态。 |
| 4 | 模块边界与总入口 | 宠物启用状态字段 | pet equipped/used state | `UsePet` 借用 `HarmonyOption` | 需修正 | 拆出独立宠物启用字段或状态结构，禁止继续复用 Harmony 语义。 |
| 5 | 模块边界与总入口 | 加载角色后初始化 | Muun/Period/Pet item load | `LoadCharacter` 后 `findAndUsePet` | 部分覆盖 | 角色进入游戏后恢复宠物、坐骑、Muun 和召唤状态。 |
| 6 | 模块边界与总入口 | 装备变化刷新 | equipment changed hooks | `EquipmentChanged` 已调用 | 部分覆盖 | 装备/背包变化后重新计算宠物效果和外观。 |
| 7 | 配置与物品识别 | PetRing 配置加载 | `IGC_PetSettings.xml` | `conf.PetRing` 已读取 | 部分覆盖 | 对齐 Helper 宠物和 Ring 配置字段。 |
| 8 | 配置与物品识别 | Helper 宠物识别 | Angel/Imp/Demon/SpiritAngel/Panda | `equippedItem(helper, code)` | 部分覆盖 | 统一识别小恶魔、守护天使、强化恶魔、强化天使、熊猫等。 |
| 9 | 配置与物品识别 | 坐骑识别 | Uniria/Dinorant/DarkHorse/Fenrir | `findAndUsePet`、`UsePet` | 部分覆盖 | 统一识别兽角、彩云兽、黑王马、炎狼兽。 |
| 10 | 配置与物品识别 | 变身戒指识别 | transformation rings | `KindBRingTransform`、`transformationRingUse` | 部分覆盖 | 统一识别变身戒指和活动戒指。 |
| 11 | 配置与物品识别 | Muun 物品识别 | `CMuunSystem::IsMuunItem` | `KindAMuun/KindBMuun` | 部分覆盖 | 通过物品分类和配置判断 Muun 道具。 |
| 12 | 配置与物品识别 | 进化石识别 | `IsStoneofEvolution` | `MuunRank`/Code 字段存在 | 未覆盖 | 判断 Muun 进化石和能量转换道具。 |
| 13 | 配置与物品识别 | 宠物可交易检查 | `IsEnableToTrade` 宠物分支 | `19-trade.md` 已记录 | 未覆盖 | 宠物、Muun、期限宠物需参与交易限制。 |
| 14 | 配置与物品识别 | 宠物可修理检查 | pet durability repair | `08-shops.md` 已记录边界 | 未覆盖 | 炎狼兽、DarkHorse 等宠物修理规则。 |
| 15 | UsePet/坐骑启用 | UsePet 协议 | `CGUsePet`/client helper use | `handle` 0xBF20、`Player.UsePet` | 部分覆盖 | 完成启用/取消坐骑宠物请求和结果码。 |
| 16 | UsePet/坐骑启用 | UsePet 消息结构 | pet position/value | `MsgUsePet` 已解析 | 已覆盖 | 已解析位置和值。 |
| 17 | UsePet/坐骑启用 | UsePet 回复 | `MsgUsePetReply` | 已实现 Marshal | 部分覆盖 | 对齐 GameServer/客户端结果码语义。 |
| 18 | UsePet/坐骑启用 | UsePet 位置校验 | inventory position check | `12..204` 检查存在 | 部分覆盖 | 确认客户端宠物启用位置范围和扩展背包。 |
| 19 | UsePet/坐骑启用 | UsePet 物品校验 | pet item code check | 支持 13,2/3/4/37 | 部分覆盖 | 补齐所有可骑乘/启用宠物类型。 |
| 20 | UsePet/坐骑启用 | UsePet 耐久校验 | durability > 0 | 已检查 | 部分覆盖 | 耐久为 0 时应返回明确失败码。 |
| 21 | UsePet/坐骑启用 | UsePet 单宠互斥 | one ride pet active | `p.pet != nil` 粗略实现 | 部分覆盖 | 同一时间只能启用一个坐骑宠物，切换时清理旧状态。 |
| 22 | UsePet/坐骑启用 | UsePet 外观广播 | `GCEquipmentChange` | `pushChangedEquipment` | 部分覆盖 | 启用/取消坐骑后同步外观和视野。 |
| 23 | UsePet/坐骑启用 | UsePet 角色重算 | `ObjCalCharacter` | `EquipmentChanged` | 部分覆盖 | 启用/取消宠物后重新计算属性。 |
| 24 | UsePet/坐骑启用 | 登录恢复启用宠物 | pet used flag restore | `findAndUsePet` 扫背包 | 需修正 | 当前依赖 `HarmonyOption == 1`，需改为独立状态。 |
| 25 | Helper 宠物效果 | 小恶魔增伤 | Imp AddAttackPercent | `Player.calc` 已接入 | 部分覆盖 | 复核物理、魔法、诅咒攻击加成一致性。 |
| 26 | Helper 宠物效果 | 强化恶魔增伤/攻速 | Demon AddAttack/AddAttackSpeed | `Player.calc` 已接入 | 部分覆盖 | 复核攻速、魔速、诅咒攻击、Buff 表示。 |
| 27 | Helper 宠物效果 | 守护天使 HP/减伤 | Angel AddHP/ReduceDamage | `Player.calc` 已接入 | 部分覆盖 | 接入统一 Buff/Effect 来源。 |
| 28 | Helper 宠物效果 | 强化天使 HP/减伤 | SpiritAngel | `Player.calc` 已接入 | 部分覆盖 | 接入统一 Buff/Effect 来源。 |
| 29 | Helper 宠物效果 | 熊猫防御 | Panda AddDefense | `Player.calc` 已接入 | 部分覆盖 | 补齐经验/掉落等熊猫效果。 |
| 30 | Helper 宠物效果 | 骷髅宠物效果 | Skeleton pet | 配置已读，未计算 | 未覆盖 | 幼龙骨架/骷髅宠物的属性和经验加成。 |
| 31 | Helper 宠物效果 | 独角兽效果 | Unicorn pet | 配置已读，未计算 | 未覆盖 | 兽角或独角兽类宠物效果接入。 |
| 32 | Helper 宠物效果 | Helper 耐久扣减 | helper durability down | 耐久系统未分类 | 未覆盖 | 攻击/受击/时间触发 Helper 耐久消耗。 |
| 33 | 变身戒指与活动戒指 | 变身戒指使用协议 | transformation ring use | `handle` 0xF321 占位 | 部分覆盖 | 使用变身戒指改变外观和 Buff。 |
| 34 | 变身戒指与活动戒指 | WizardRing 效果 | WizardRing | 配置已读，未完整计算 | 未覆盖 | 魔法戒指攻击/防御/外观效果。 |
| 35 | 变身戒指与活动戒指 | SkeletonRing 效果 | SkeletonRing | 配置已读，未完整计算 | 未覆盖 | 骷髅变身戒指效果。 |
| 36 | 变身戒指与活动戒指 | ChristmasRing 效果 | ChristmasRing | 配置已读，未完整计算 | 未覆盖 | 圣诞变身戒指效果。 |
| 37 | 变身戒指与活动戒指 | PandaRing 效果 | PandaRing/PandaBrown/PandaPink | 配置已读，未完整计算 | 未覆盖 | 熊猫类变身戒指效果。 |
| 38 | 变身戒指与活动戒指 | Robot/Mage/Decoration Ring | RobotKnight/MiniRobot/Mage/Decoration | 配置已读，未完整计算 | 未覆盖 | 后续版本活动戒指效果。 |
| 39 | 变身戒指与活动戒指 | 变身外观同步 | Viewport equipment/appearance | 外观序列化有基础 | 未覆盖 | 变身戒指启用后广播外观变化。 |
| 40 | 变身戒指与活动戒指 | 变身戒指 Buff 边界 | Ring Buff | `13-buffs.md` 已记录 | 未覆盖 | 宠物系统触发，Buff 生命周期归 Buff 系统。 |
| 41 | DarkHorse 黑王马 | DarkHorse 物品识别 | Code(13,4) | `UsePet` 支持 | 部分覆盖 | 统一黑王马物品、等级、经验、耐久字段。 |
| 42 | DarkHorse 黑王马 | DarkHorse 减伤 | DarkHorse damage reduce | `petReduceDamage=(30+Level)/2` | 部分覆盖 | 对齐 GameServer 黑王马减伤公式。 |
| 43 | DarkHorse 黑王马 | DarkHorse 防御加成 | Dark Lord pet defense | 公式未完整 | 未覆盖 | 黑王马防御、吸收或属性加成接入公式。 |
| 44 | DarkHorse 黑王马 | DarkHorse 技能 | `SkillDarkHorseAttack` | `10-skills.md` 记录缺口 | 未覆盖 | 技能系统调用宠物系统执行黑王马攻击。 |
| 45 | DarkHorse 黑王马 | DarkHorse 经验 | pet exp | `09-exp.md` 记录缺口 | 未覆盖 | 玩家获得经验时同步黑王马经验。 |
| 46 | DarkHorse 黑王马 | DarkHorse 升级 | pet level up | 暂无 | 未覆盖 | 宠物经验达到阈值后升级并下发消息。 |
| 47 | DarkHorse 黑王马 | DarkHorse 耐久 | pet durability | 暂无 | 未覆盖 | 攻击/受击/死亡/修理影响耐久。 |
| 48 | DarkHorse 黑王马 | DarkHorse 合成边界 | `DarkHorseChaosMix` | `12-mix.md` 已记录 | 未覆盖 | 本模块只消费已生成宠物，道具生成归合成系统。 |
| 49 | DarkSpirit 黑王鸟 | DarkSpirit 管理器 | `CDarkSpirit` | 暂无 | 未覆盖 | 实现黑王鸟战斗宠物管理器。 |
| 50 | DarkSpirit 黑王鸟 | DarkSpirit 物品信息协议 | `PMSG_SEND_PET_ITEMINFO` | `getPetItemInfo` 0xA9 占位 | 部分覆盖 | 下发黑王鸟等级、经验、耐久等。 |
| 51 | DarkSpirit 黑王鸟 | DarkSpirit 命令协议 | `ChangeCommand` | `getPetItemCommand` 0xA7 占位 | 部分覆盖 | 设置攻击模式和目标。 |
| 52 | DarkSpirit 黑王鸟 | DarkSpirit 模式 Normal | `ModeNormal` | 暂无 | 未覆盖 | 普通待机模式。 |
| 53 | DarkSpirit 黑王鸟 | DarkSpirit 随机攻击 | `ModeAttackRandom` | 暂无 | 未覆盖 | 自动选择目标攻击。 |
| 54 | DarkSpirit 黑王鸟 | DarkSpirit 随主人攻击 | `ModeAttackWithMaster` | 暂无 | 未覆盖 | 跟随主人目标攻击。 |
| 55 | DarkSpirit 黑王鸟 | DarkSpirit 指定目标 | `ModeAttackTarget` | 暂无 | 未覆盖 | 按客户端目标攻击。 |
| 56 | DarkSpirit 黑王鸟 | DarkSpirit 目标设置 | `SetTarget/ReSetTarget` | 暂无 | 未覆盖 | 记录和重置黑王鸟目标对象。 |
| 57 | DarkSpirit 黑王鸟 | DarkSpirit 攻击 | `Attack` | 暂无 | 未覆盖 | 黑王鸟独立攻击流程，调用战斗系统。 |
| 58 | DarkSpirit 黑王鸟 | DarkSpirit 伤害计算 | `GetAttackDamage` | `conf.CalcCharacter.DarkSpiritDamageRate` | 部分覆盖 | 使用配置计算 PVE/PVP 伤害倍率。 |
| 59 | DarkSpirit 黑王鸟 | DarkSpirit 命中 | `MissCheck/MissCheckPvP` | 暂无 | 未覆盖 | 实现 PVE/PVP 命中率。 |
| 60 | DarkSpirit 黑王鸟 | DarkSpirit 盾伤害 | `GetShieldDamage` | 暂无 | 未覆盖 | 计算 SD/Shield 伤害分摊。 |
| 61 | DarkSpirit 黑王鸟 | DarkSpirit 攻击消息 | `SendAttackMsg` | 暂无 | 未覆盖 | 向视野广播黑王鸟攻击。 |
| 62 | DarkSpirit 黑王鸟 | DarkSpirit 经验倍率 | `DarkSpiritAddExperience` | 配置已读 | 部分覆盖 | 玩家经验同步到黑王鸟时应用倍率。 |
| 63 | DarkSpirit 黑王鸟 | DarkSpirit 合成边界 | `DarkSpiritChaosMix` | `12-mix.md` 已记录 | 未覆盖 | 黑王鸟生成归合成系统，本模块负责使用和成长。 |
| 64 | Fenrir 炎狼兽 | Fenrir 物品识别 | Code(13,37) | `UsePet` 支持 | 部分覆盖 | 统一红/蓝/黑/金狼状态和效果。 |
| 65 | Fenrir 炎狼兽 | Fenrir 增伤 | Fenrir damage increase | `petIncreaseDamage=33` | 部分覆盖 | 对齐红狼/金狼等不同类型增伤。 |
| 66 | Fenrir 炎狼兽 | Fenrir 减伤 | Fenrir damage reduce | `petReduceDamage=10` | 部分覆盖 | 对齐蓝狼/黑狼/金狼减伤差异。 |
| 67 | Fenrir 炎狼兽 | Fenrir 技能 | `SkillFenrirAttack` | `10-skills.md` 记录缺口 | 未覆盖 | 技能系统调用宠物系统执行炎狼兽技能。 |
| 68 | Fenrir 炎狼兽 | Fenrir 耐久上限 | Fenrir max durability | `ItemFenrirDefaultMaxDurSmall` 等配置 | 部分覆盖 | 按职业和配置计算耐久。 |
| 69 | Fenrir 炎狼兽 | Fenrir 修理 | `FenrirRepairRate` | 配置已读，商店未接 | 未覆盖 | 修理费用、成功率和耐久恢复。 |
| 70 | Fenrir 炎狼兽 | Fenrir 材料掉落边界 | Fenrir stuff drop | `14-drops.md` 已记录 | 未覆盖 | 掉落系统负责材料掉落，本模块只记录消费关系。 |
| 71 | Fenrir 炎狼兽 | Fenrir 合成边界 | `Fenrir_*_Mix` | `12-mix.md` 已记录 | 未覆盖 | 合成和升级生成归合成系统。 |
| 72 | Muun 配置与信息 | MuunSystem 管理器 | `CMuunSystem` | `MuunSystem` 空函数 | 未覆盖 | 实现 Muun 总管理器和协议入口。 |
| 73 | Muun 配置与信息 | MuunInfo | `CMuunInfo` | `item.MuunRank` 等字段 | 部分覆盖 | 保存 Muun item num、type、rank、option、期限和进化目标。 |
| 74 | Muun 配置与信息 | MuunOpt | `CMuunOpt` | 暂无 | 未覆盖 | 保存 Muun 选项类型、等级值和条件。 |
| 75 | Muun 配置与信息 | MuunInfoMng | `CMuunInfoMng` | 暂无 | 未覆盖 | 加载并索引 Muun 配置。 |
| 76 | Muun 配置与信息 | Muun 系统配置加载 | `LoadScriptMuunSystemInfo` | 暂无 | 未覆盖 | 加载 Muun 基础配置。 |
| 77 | Muun 配置与信息 | Muun 选项配置加载 | `LoadScriptMuunSystemOption` | 暂无 | 未覆盖 | 加载 Muun 效果配置。 |
| 78 | Muun 配置与信息 | Muun 兑换配置加载 | `LoadScriptMuunExchange` | 暂无 | 未覆盖 | 加载 Muun 兑换需求和奖励 Bag。 |
| 79 | Muun 配置与信息 | Muun 选项枚举 | `MUUN_OPTIONS` | 暂无 | 未覆盖 | 支持攻击、防御、卓越、暴击、技能、元素、攻击技能等。 |
| 80 | Muun 背包与 DB | Muun 背包模型 | `MUUN_INVENTORY_SIZE` | 角色外观字段存在，背包无 | 未覆盖 | 建立 Muun 专用背包和装备槽。 |
| 81 | Muun 背包与 DB | Muun DB 加载 | `GDReqLoadMuunInvenItem` | 暂无 | 未覆盖 | 登录后加载 Muun 背包。 |
| 82 | Muun 背包与 DB | Muun DB 响应 | `DGLoadMuunInvenItem` | 暂无 | 未覆盖 | 应用 DB 返回的 Muun 背包物品。 |
| 83 | Muun 背包与 DB | Muun DB 保存 | `GDReqSaveMuunInvenItem` | 暂无 | 未覆盖 | 登出、移动、使用、升级后保存 Muun 背包。 |
| 84 | Muun 背包与 DB | Muun 期限数据 | `MUUN_ITEM_PERIOD_DATA` | 暂无 | 未覆盖 | 管理期限 Muun 的账号、角色、serial、过期时间。 |
| 85 | Muun 背包与 DB | Muun 期限加载 | `AddMuunItemPeriodInfo` | 暂无 | 未覆盖 | 登录后恢复期限 Muun 数据。 |
| 86 | Muun 背包与 DB | Muun 期限清理 | `ClearPeriodMuunItemData` | 暂无 | 未覆盖 | 过期、删除、下线时清理期限数据。 |
| 87 | Muun 装备与效果 | Muun 装备 | `MuunItemEquipment` | 暂无 | 未覆盖 | 将 Muun 从背包装备到主/副/骑乘槽。 |
| 88 | Muun 装备与效果 | Muun 外观字段 | MuunItem/SubItem/RideItem | `model.Player` 有字段 | 部分覆盖 | 装备后同步角色外观中的 Muun 字段。 |
| 89 | Muun 装备与效果 | 设置 Muun 效果 | `SetUserMuunEffect` | 暂无 | 未覆盖 | 装备 Muun 后写入玩家 Muun 效果列表。 |
| 90 | Muun 装备与效果 | 移除 Muun 效果 | `RemoveUserMuunEffect` | 暂无 | 未覆盖 | 卸下、过期、条件失效时移除效果。 |
| 91 | Muun 装备与效果 | Muun 效果数值查询 | `GetMuunItemValueOfOptType` | 暂无 | 未覆盖 | 公式、战斗、掉落系统查询 Muun 效果。 |
| 92 | Muun 装备与效果 | Muun 角色属性计算 | `CalCharacterStat` | `05-formula.md` 记录缺口 | 未覆盖 | 将 Muun 效果接入角色属性重算。 |
| 93 | Muun 装备与效果 | Muun 条件状态下发 | `GCSendConditionStatus` | 暂无 | 未覆盖 | 条件满足/失效时通知客户端。 |
| 94 | Muun 条件检查 | Muun 总条件检查 | `CheckMuunItemCondition` | 暂无 | 未覆盖 | 对装备 Muun 判断是否生效。 |
| 95 | Muun 条件检查 | 时间条件 | `ChkMuunOptConditionTime` | 暂无 | 未覆盖 | 按开始/结束时间判断。 |
| 96 | Muun 条件检查 | 日期/星期条件 | `ChkMuunOptConditionDay` | 暂无 | 未覆盖 | 按星期或日期判断。 |
| 97 | Muun 条件检查 | 等级条件 | `ChkMuunOptConditionLevel` | `09-exp.md` 记录升级刷新 | 未覆盖 | 角色升级后刷新 Muun 条件。 |
| 98 | Muun 条件检查 | 地图条件 | `ChkMuunOptConditionMap` | 暂无 | 未覆盖 | 角色换图后刷新 Muun 条件。 |
| 99 | Muun 条件检查 | 装备条件检查流程 | `CheckEquipMuunItemConditionProc` | 暂无 | 未覆盖 | 周期或事件驱动刷新装备 Muun 生效状态。 |
| 100 | Muun 升级/进化/兑换 | Muun 使用物品协议 | `CGMuunInventoryUseItemRecv` | `handle` 0x4E11、`MsgMuunSystem` 空 | 部分覆盖 | 解析 Muun 背包使用、升级、进化、兑换请求。 |
| 101 | Muun 升级/进化/兑换 | Muun 使用结果 | `GCMuunInventoryUseItemResult` | 暂无 | 未覆盖 | 下发成功、材料不足、等级不足、位置错误等结果。 |
| 102 | Muun 升级/进化/兑换 | Muun 进化 | `MuunItemEvolution` | 暂无 | 未覆盖 | 使用进化石将 Muun 进化到目标物品。 |
| 103 | Muun 升级/进化/兑换 | Muun 升级 | `MuunItemLevelUp` | 暂无 | 未覆盖 | 使用材料提升 Muun 等级。 |
| 104 | Muun 升级/进化/兑换 | Muun 生命宝石 | `MuunItemLifeGem` | 暂无 | 未覆盖 | 使用生命宝石处理 Muun 生命/耐久类能力。 |
| 105 | Muun 升级/进化/兑换 | Muun 能量转换 | `MuunItemEnergyGenerator` | 暂无 | 未覆盖 | 根据 rank/level 点数转换能量。 |
| 106 | Muun 升级/进化/兑换 | Muun 兑换请求 | `CGMuunExchangeItem` | 暂无 | 未覆盖 | 根据兑换配置消耗物品并生成奖励。 |
| 107 | Muun 升级/进化/兑换 | Muun 兑换材料检查 | `ChkMuunExchangeInvenNeedItem` | 暂无 | 未覆盖 | 检查背包所需材料和数量。 |
| 108 | Muun 升级/进化/兑换 | Muun 兑换背包空间 | `ChkMuunExchangeInvenEmpty` | 背包空间工具存在 | 未覆盖 | 生成奖励前确认空间。 |
| 109 | Muun 升级/进化/兑换 | Muun 兑换奖励生成 | `GDMuunExchangeInsertInven` | 暂无 | 未覆盖 | 生成兑换奖励并写入背包。 |
| 110 | Muun 攻击与骑乘 | Muun 攻击对象 | `CMuunAttack` | 暂无 | 未覆盖 | 实现 Muun 攻击技能执行器。 |
| 111 | Muun 攻击与骑乘 | Muun 攻击消息 | `SendAttackMsg` | 暂无 | 未覆盖 | 向视野广播 Muun 攻击。 |
| 112 | Muun 攻击与骑乘 | Muun 技能处理 | `SkillProc` | 暂无 | 未覆盖 | 根据 Muun 攻击选项触发技能。 |
| 113 | Muun 攻击与骑乘 | Muun 伤害吸收 | `DamageAbsorb` | 暂无 | 未覆盖 | 实现 Muun 减伤/吸收效果。 |
| 114 | Muun 攻击与骑乘 | Muun Stun | `Stun` | 暂无 | 未覆盖 | 实现 Muun 眩晕效果。 |
| 115 | Muun 攻击与骑乘 | Muun 攻击伤害 | `GetAttackDamage` | 暂无 | 未覆盖 | 计算 Muun 攻击伤害。 |
| 116 | Muun 攻击与骑乘 | Muun 骑乘选择 | `CGReqRideSelect` | 角色外观有 RideItem 字段 | 未覆盖 | 选择 Muun 骑乘并同步外观。 |
| 117 | 召唤技能与召唤物 | 召唤技能入口 | `SkillMonsterCall`、`SkillSummon` | `10-skills.md` 已记录 | 未覆盖 | 技能系统调用宠物与召唤系统创建召唤物。 |
| 118 | 召唤技能与召唤物 | 基础召唤技能 | SummonGoblin/StoneGolem/... | 技能索引已定义 | 部分覆盖 | 根据技能类型召唤对应怪物。 |
| 119 | 召唤技能与召唤物 | 星云召唤 | `SkillIndexSummon` | 技能索引已定义 | 部分覆盖 | 召唤队友或目标到指定位置的规则边界。 |
| 120 | 召唤技能与召唤物 | 召唤数量限制 | `MaxSummonMonsterCount` | `ObjectManager.maxCallMonsterCount` | 部分覆盖 | 控制每玩家或全服召唤怪数量。 |
| 121 | 召唤技能与召唤物 | 召唤物对象创建 | gObjAddMonster/summon object | 对象管理器可创建怪物 | 未覆盖 | 创建召唤怪并绑定主人。 |
| 122 | 召唤技能与召唤物 | 召唤物主人绑定 | owner/master index | 暂无正式字段 | 未覆盖 | 记录召唤物主人、持续时间和阵营。 |
| 123 | 召唤技能与召唤物 | 召唤物生命周期 | summon life/recall/delete | 暂无 | 未覆盖 | 超时、主人下线、换图、死亡时清理召唤物。 |
| 124 | 召唤技能与召唤物 | 召唤物 AI 边界 | summoned monster AI | `25-monster-ai.md` | 未覆盖 | 宠物系统创建和绑定，具体 AI 由怪物 AI 系统承载。 |
| 125 | 宠物经验/等级/耐久 | 宠物经验同步 | `gObjSetExpPetItem` | `09-exp.md` 已记录 | 未覆盖 | 玩家获得经验时同步 DarkHorse/DarkSpirit 等宠物经验。 |
| 126 | 宠物经验/等级/耐久 | 宠物升级消息 | `CDarkSpirit::SendLevelmsg` | 暂无 | 未覆盖 | 宠物升级后下发等级消息。 |
| 127 | 宠物经验/等级/耐久 | 宠物死亡/耐久清零 | pet dead/durability zero | 暂无 | 未覆盖 | 宠物耐久为 0 时禁用效果和外观。 |
| 128 | 宠物经验/等级/耐久 | 宠物耐久扣减 | attack/hit/time durability | `07-items.md` 已记录耐久缺口 | 未覆盖 | 按攻击、受击、技能或时间扣耐久。 |
| 129 | 宠物经验/等级/耐久 | 宠物修理费用 | repair pet item | `08-shops.md` 已记录边界 | 未覆盖 | 计算宠物修理价格并恢复耐久。 |
| 130 | 跨系统接口 | 与对象系统联动 | player/object/summon hooks | `04-objects.md` 需接入 | 未覆盖 | 对象系统通知加载、装备、死亡、下线、召唤物创建清理。 |
| 131 | 跨系统接口 | 与公式系统联动 | pet/ring/Muun stat formulas | `05-formula.md` 已记录 | 未覆盖 | 公式系统查询宠物与 Muun 效果，不重复公式。 |
| 132 | 跨系统接口 | 与地图系统联动 | mount map restrictions | `06-maps.md` 已记录 | 未覆盖 | 地图移动检查坐骑/翅膀/宠物限制。 |
| 133 | 跨系统接口 | 与道具系统联动 | item move/use/durability | `07-items.md` 需接入 | 未覆盖 | 宠物启用、Muun 背包、耐久、交易、修理都依赖道具系统。 |
| 134 | 跨系统接口 | 与经验系统联动 | pet exp | `09-exp.md` 已记录 | 未覆盖 | 经验系统通知宠物系统分配宠物经验。 |
| 135 | 跨系统接口 | 与技能系统联动 | summon/Fenrir/DarkHorse/DarkSpirit skills | `10-skills.md` 已记录 | 未覆盖 | 技能系统分发到宠物系统执行特殊宠物技能。 |
| 136 | 跨系统接口 | 与合成系统联动 | DarkHorse/DarkSpirit/Fenrir mix | `12-mix.md` 已记录 | 未覆盖 | 合成系统生成宠物道具，本模块负责使用和成长。 |
| 137 | 跨系统接口 | 与 Buff 系统联动 | pet/ring/period buffs | `13-buffs.md` 已记录 | 未覆盖 | 宠物系统触发效果来源，Buff 生命周期由 Buff 系统维护。 |
| 138 | 跨系统接口 | 与掉落系统联动 | Fenrir materials/ring drops | `14-drops.md` 已记录 | 未覆盖 | 掉落系统生成宠物材料和戒指，本模块消费效果。 |
| 139 | 协议与测试 | UsePet 测试 | 0xBF20 | 待实现/补齐 | 部分覆盖 | 覆盖位置错误、耐久 0、启用、取消、切换、外观同步。 |
| 140 | 协议与测试 | PetItemInfo/Command 测试 | 0xA7/0xA9 | 待实现 | 未覆盖 | 覆盖 DarkSpirit 信息、命令、目标、模式切换。 |
| 141 | 协议与测试 | MuunSystem 测试 | 0x4E11 | 待实现 | 未覆盖 | 覆盖 Muun 加载、装备、卸下、升级、进化、兑换、期限。 |
| 142 | 协议与测试 | 变身戒指测试 | 0xF321 | 待实现 | 未覆盖 | 覆盖使用、外观、Buff、过期和取消。 |
| 143 | 协议与测试 | 宠物经验测试 | pet exp flow | 待实现 | 未覆盖 | 覆盖玩家经验获得、宠物倍率、升级和消息下发。 |
| 144 | 协议与测试 | 召唤物测试 | summon skills | 待实现 | 未覆盖 | 覆盖召唤数量、创建、主人绑定、死亡、下线清理。 |
| 145 | 协议与测试 | 跨系统边界测试 | items/formula/skill/buff/mix/drop | 待实现 | 未覆盖 | 确认宠物系统只编排规则，不重复基础系统实现。 |
