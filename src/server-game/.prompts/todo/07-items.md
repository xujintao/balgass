# 7. 道具系统

本模块覆盖物品基础表、物品实例、协议编码、背包/仓库容器、装备穿戴、拾取/丢弃/买卖/修理/使用、掉落生成、卓越/幸运/套装/380/Harmony/Socket/Pentagram/期限道具。地图地面物品容器、商店、合成分别归地图系统、商店系统、合成系统；本模块只记录道具侧的实例、属性、容器和调用边界。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 物品基础表 | 物品编码规则 | `ITEMGET`、`ItemGetNumberMake` | `item.Code(section,index)` | 已覆盖 | Go 侧采用 `section*512+index`，需保持与协议和配置表一致。 |
| 2 | 物品基础表 | ItemList 配置加载 | `OpenItemScript`、`ItemAttribute` | `ItemTable.init` 读取 `Items/IGC_ItemList.xml` | 已覆盖 | 当前已按 Section 加载所有 `ItemBase`。 |
| 3 | 物品基础表 | Section 映射表 | `MAX_TYPE_ITEMS`、`ItemAttribute` | `itemTable []map[int]*ItemBase` | 已覆盖 | Go 侧 map 查询比数组灵活，需防止缺 section panic。 |
| 4 | 物品基础表 | ItemBase 字段 | `ITEM_ATTRIBUTE` | `ItemBase` | 已覆盖 | 基础攻击、防御、需求、价格、交易、修理等字段已有结构。 |
| 5 | 物品基础表 | 职业需求映射 | `CItem::IsClass` | `ReqClass`、`ItemTable.init` | 部分覆盖 | GrowLancer 仍注释，职业进阶等级语义需继续核对。 |
| 6 | 物品基础表 | 物品类型枚举 | `ItemAttribute.Type` 等配置 | `itemType` | 已覆盖 | 已覆盖 common/regular/socket/380/lucky/event/archangel/chaos。 |
| 7 | 物品基础表 | 物品大类枚举 | `GetItemKindA` | `itemKindA` | 已覆盖 | 武器、防具、翅膀、宝石、技能、Pentagram 等已建模。 |
| 8 | 物品基础表 | 物品小类枚举 | `GetItemKindB` | `itemKindB` | 已覆盖 | 细分武器、翅膀、药水、宝石、Muun 等，后续需补客户端全部类别。 |
| 9 | 物品基础表 | 基础表查询 | `IsItem`、`ItemAttribute[item]` | `GetItemBase`、`GetItemBaseMust` | 已覆盖 | `GetItemBaseMust` 适合内部强依赖，外部输入应使用带 error 版本。 |
| 10 | 物品基础表 | 掉落等级计算 | `zzzItemLevel`、`GetLevelItem` | `ItemTable.GetItemLevel` | 部分覆盖 | 已处理普通掉落、首饰、宝石、召唤石和变身戒指；需与 GameServer 掉落表对齐。 |
| 11 | 物品实例与计算 | Item 实例结构 | `CItem` | `item.Item` | 已覆盖 | Go 结构已包含基础字段、优秀、套装、380、Harmony、Socket、Pentagram 等。 |
| 12 | 物品实例与计算 | 创建物品实例 | `CItem::Convert` | `NewItem` | 部分覆盖 | 当前只创建基础实例，未完整应用 option/serial/socket/period 参数。 |
| 13 | 物品实例与计算 | 复制物品实例 | `CItem` 拷贝赋值 | `Item.Copy` | 已覆盖 | 浅拷贝足够复制实例字段，但共享 `ItemBase` 是预期行为。 |
| 14 | 物品实例与计算 | 卓越判断 | `CItem::IsExtItem`、`m_NewOption` | `Item.IsExcellent` | 已覆盖 | 已覆盖普通优秀、翅膀优秀、2.5/3代翅膀优秀字段。 |
| 15 | 物品实例与计算 | 卓越数量统计 | `CalcExcOption`、优秀 bit 数 | `Item.ExcellentCount` | 已覆盖 | 用于价格、掉落和后续限制规则。 |
| 16 | 物品实例与计算 | 套装判断 | `CItem::IsSetItem` | `Item.IsSet` | 已覆盖 | 依赖 `Item.Set > 0`，和 `SetManager` 配合。 |
| 17 | 物品实例与计算 | 技能索引获取 | `m_Option1`、技能字段 | `Item.GetSkillIndex` | 已覆盖 | 装备技能和技能书学习都依赖该字段。 |
| 18 | 物品实例与计算 | 翅膀追加随机 | `CItem::Convert` 翅膀 option | `Item.RandWingAdditionKind` | 部分覆盖 | Go 已按翅膀类型随机追加方向，需补掉落/合成调用点。 |
| 19 | 物品实例与计算 | 卓越 bit 解码 | `m_NewOption` 解析 | `Item.DecodeExcellent` | 已覆盖 | Go 已将 bit 映射到不同种类优秀属性。 |
| 20 | 物品实例与计算 | 实例重算入口 | `CItem::Value` | `Item.Calc` | 已覆盖 | 统一重算耐久、攻击、魔攻、防御、防御率、追加和价格。 |
| 21 | 物品协议与持久化 | 12 字节协议编码 | `ItemByteConvert`、`ItemByteConvert32` | `Item.Marshal` | 部分覆盖 | 基础字段可编码，但需要用测试锁定 bit 优先级和客户端兼容。 |
| 22 | 物品协议与持久化 | Zen 物品编码 | `MoneyItemDrop` 特殊编码 | `Item.Marshal` 中 `Code(14,15)` | 已覆盖 | Zen 使用 durability 保存金额并特殊编码。 |
| 23 | 物品协议与持久化 | 幸运/技能/等级编码 | `ItemByteConvert` option byte | `Item.Marshal` data[1] | 需核对 | 表达式 `item.Addition & 0x0C >> 2` 需要括号测试。 |
| 24 | 物品协议与持久化 | 卓越 bit 编码 | `ItemByteConvert` NewOption | `Item.Marshal` data[3] | 部分覆盖 | 已覆盖多类优秀字段，需补完整 round-trip 测试。 |
| 25 | 物品协议与持久化 | 套装 tier 编码 | `SetOption` | `SetManager.GetTierIndex` | 已覆盖 | `data[4]` 写入套装 tier。 |
| 26 | 物品协议与持久化 | 380/期限字段编码 | `m_ItemOptionEx`、`m_PeriodItemOption` | `Item.Marshal` data[5] | 部分覆盖 | 字段可编码，但期限道具系统未完整落地。 |
| 27 | 物品协议与持久化 | Harmony 字段编码 | `m_JewelOfHarmonyOption` | `Item.Marshal` data[6] | 部分覆盖 | 当前只写 `HarmonyOption`，需确认与 effect/level 合成规则。 |
| 28 | 物品协议与持久化 | Socket/Pentagram 编码 | `m_SocketOption`、`m_BonusSocketOption` | `SocketBonus`、`SocketSlot1-5` | 部分覆盖 | 字段可下发，效果应用和镶嵌流程未完整实现。 |
| 29 | 物品协议与持久化 | 二进制反序列化 | `BufferItemtoConvert3`、DB item decode | `Item.Unmarshal` | 未覆盖 | 当前为空，是协议反解和二进制 DB 导入的核心缺口。 |
| 30 | 物品协议与持久化 | JSON/DB 持久化 | GameServer DB byte array | `PositionedItems.MarshalJSON`、`Value/Scan` | 部分覆盖 | Go 侧采用 JSON 保存，需保持字段迁移兼容和坏数据修复。 |
| 31 | 背包容器 | 背包尺寸 | `INVENTORY_SIZE`、`INVENTORY_MAP_SIZE` | `Inventory.UnmarshalJSON` Size=237 | 已覆盖 | 需明确 237 的区域划分和客户端扩展背包对应关系。 |
| 32 | 背包容器 | 装备栏切片 | `INVETORY_WEAR_SIZE` | `INVENTORY_WEAR_SIZE`、`WearingItems` | 已覆盖 | 前 12 格作为穿戴装备栏。 |
| 33 | 背包容器 | 背包占格 flags | `pInventoryMap` | `Inventory.Flags` | 已覆盖 | Go 使用 bool flags 表示占格。 |
| 34 | 背包容器 | 装备栏占位规则 | `gObjInventoryItemSet` | `CheckFlagsForItem(position<12)` | 部分覆盖 | 只检查槽位空闲，缺槽位类型和穿戴条件校验。 |
| 35 | 背包容器 | 背包矩形检查 | `InventoryExtentCheck`、`gObjOnlyInventoryRectCheck` | `Inventory.CheckFlagsForItem` | 已覆盖 | 8 列、多段高度边界已实现。 |
| 36 | 背包容器 | 扩展背包边界 | `CheckOutOfInventory`、扩展字段 | hardcoded 8/12/16 高度 | 部分覆盖 | Go 有固定扩展边界，但未接入 `inventoryExpansion` 动态开关。 |
| 37 | 背包容器 | 空位查找 | `gObjInventoryInsertItem` | `FindFreePositionForItem` | 已覆盖 | 从背包区开始查找，不占装备栏。 |
| 38 | 背包容器 | 堆叠查找 | `IsOverlapItem`、堆叠逻辑 | `FindFreePositionForItem` overlap 分支 | 部分覆盖 | 已支持同 code/level 耐久堆叠，需补不同宝石包等级规则。 |
| 39 | 背包容器 | 放入物品 | `gObjInventoryInsertItemPos` | `Inventory.AddItem` | 已覆盖 | 写入 item 并设置 flags。 |
| 40 | 背包容器 | 移除物品 | `gObjInventoryDeleteItem` | `Inventory.RemoveItem` | 已覆盖 | 清空 item 并释放 flags。 |
| 41 | 仓库容器 | 仓库尺寸 | `WAREHOUSE_SIZE` | `Warehouse.UnmarshalJSON` Size=240 | 已覆盖 | Go 仓库大小和 flags 初始化已实现。 |
| 42 | 仓库容器 | 仓库占格 flags | `pWarehouseMap` | `Warehouse.Flags` | 已覆盖 | 与背包结构类似。 |
| 43 | 仓库容器 | 仓库矩形检查 | `WarehouseExtentCheck` | `Warehouse.CheckFlagsForItem` | 已覆盖 | 8 列、两页高度边界已实现。 |
| 44 | 仓库容器 | 仓库扩展边界 | `CheckOutOfWarehouse` | `maxHeight1/maxHeight2` 固定 | 部分覆盖 | 未接入角色仓库扩展证书状态。 |
| 45 | 仓库容器 | 仓库空位查找 | `gObjWarehouseInsertItem` | `Warehouse.FindFreePositionForItem` | 已覆盖 | 查找可容纳物品的空位。 |
| 46 | 仓库容器 | 仓库放入 | `gObjWarehouseInsertItemPos` | `Warehouse.AddItem` | 已覆盖 | 物品写入和 flags 更新已实现。 |
| 47 | 仓库容器 | 仓库移除 | `gObjWarehouseDeleteItem` | `Warehouse.RemoveItem` | 已覆盖 | 移除后释放占格。 |
| 48 | 仓库容器 | 仓库 JSON 加载 | GameServer DB byte decode | `Warehouse.UnmarshalJSON` | 部分覆盖 | 可加载 JSON，但不是 GameServer 原始二进制格式。 |
| 49 | 仓库容器 | 仓库 DB Scan | DB warehouse load | `Warehouse.Scan` | 已覆盖 | 当前接收 `[]byte` 并走 JSON 解析。 |
| 50 | 仓库容器 | 仓库金币/锁 | `WarehouseMoney`、`WarehouseLock/PW` | `MsgWarehouseMoneyReply`、player 关闭窗口 | 未覆盖 | 仓库物品容器存在，但金币、密码、锁和保存状态未完整实现。 |
| 51 | 装备穿戴与角色重算 | 装备槽定义 | `pInventory[0..11]` | `InventoryWearSize`、`WearingItems` | 已覆盖 | Go 装备槽数量已定义为 12。 |
| 52 | 装备穿戴与角色重算 | 装备穿戴移动 | `CGInventoryItemMove` | `Object.MoveItem` 背包内移动 | 部分覆盖 | 移动到装备栏会触发重算，但穿戴条件不完整。 |
| 53 | 装备穿戴与角色重算 | 装备卸下移动 | `CGInventoryItemMove` | `Object.MoveItem` 装备栏到背包/仓库 | 部分覆盖 | 卸下会重算，需补特殊状态限制。 |
| 54 | 装备穿戴与角色重算 | 穿戴需求校验 | `CItem::IsClass`、需求属性检查 | `Player.limitUseItem` | 部分覆盖 | 函数存在，但 `MoveItem` 穿戴路径未明确调用。 |
| 55 | 装备穿戴与角色重算 | 职业限制校验 | `IsClass(Class,ChangeUp)` | `ReqClass[p.Class]` | 部分覆盖 | 已按 class/changeUp 判定，需补召唤/格斗/未来职业细节。 |
| 56 | 装备穿戴与角色重算 | 双手武器限制 | `TwoHand`、左右手规则 | 当前无完整校验 | 未覆盖 | 穿戴双手武器时应限制副手，盾/箭筒/弩箭需特殊处理。 |
| 57 | 装备穿戴与角色重算 | 装备技能增删 | `gObjInventoryEquipment` | `Player.EquipmentChanged` | 已覆盖 | 武器技能会自动学习/遗忘并推送技能变动。 |
| 58 | 装备穿戴与角色重算 | 装备外观广播 | `GCEquipmentChange` | `pushChangedEquipment` | 部分覆盖 | `UsePet` 已广播，装备移动路径也应完整广播外观。 |
| 59 | 装备穿戴与角色重算 | 装备后角色重算 | `ObjCalCharacter::CalcCharacter` | `EquipmentChanged -> calc` | 已覆盖 | 所有装备变更后应走该入口。 |
| 60 | 装备穿戴与角色重算 | 宠物启用状态 | `m_btInvenPetPos`、宠物槽 | `UsePet`、`findAndUsePet` | 需修正 | 当前借用 `HarmonyOption` 表示启用状态，语义冲突，应拆独立字段。 |
| 61 | 拾取/丢弃/买卖 | 地面拾取入口 | `CGItemGetRequest`、`MapClass::ItemGive` | `Object.PickItem` | 部分覆盖 | 基础拾取可用，缺归属、距离、权限和防抢规则。 |
| 62 | 拾取/丢弃/买卖 | Zen 拾取 | `MoneyItemDrop` 拾取 | `PickItem` 处理 `Code(14,15)` | 已覆盖 | 会增加角色 Money 并返回 `position=-2`。 |
| 63 | 拾取/丢弃/买卖 | 普通物品拾取 | `ItemGive` | 背包 `FindFreePositionForItem/AddItem` | 已覆盖 | 背包有空位时加入物品。 |
| 64 | 拾取/丢弃/买卖 | 堆叠拾取 | `IsOverlapItem` | `PickItem` 合并 Durability | 部分覆盖 | 支持基础堆叠，需补堆叠上限和协议位置语义测试。 |
| 65 | 拾取/丢弃/买卖 | 背包满失败 | `ItemGive` 失败回复 | `PickItem` position=-1 | 已覆盖 | 失败不会删除地图物品。 |
| 66 | 拾取/丢弃/买卖 | 地图物品删除 | `MapClass::StateSetDestroy` | `maps.MapManager.RemoveItem` | 已覆盖 | 拾取成功后删除地图物品。 |
| 67 | 拾取/丢弃/买卖 | 丢弃入口 | `CGItemDropRequest` | `Object.DropItem` | 部分覆盖 | 可丢弃到地图，但缺交易锁、绑定、贵重物品确认等规则。 |
| 68 | 拾取/丢弃/买卖 | 丢弃坐标校验 | `MapClass::ItemDrop`、地图属性 | `MapManager.AddItem` | 部分覆盖 | 地图系统只校验 valid，仍需过滤阻挡/安全区/事件限制。 |
| 69 | 拾取/丢弃/买卖 | NPC 买入 | `CShop::Buy`、`gObjInventoryInsertItem` | `Object.BuyItem` | 部分覆盖 | 基础距离、NPC 类型、金钱、背包空位已校验；缺商店限制和税率。 |
| 70 | 拾取/丢弃/买卖 | NPC 卖出 | `CShop::Sell` | `Object.SellItem` | 部分覆盖 | 基础卖出和加钱已实现，缺不可卖、全优秀、绑定等配置限制。 |
| 71 | 移动物品与协议入口 | MoveItem 协议解析 | `CGInventoryItemMove` | `MsgMoveItem.Unmarshal` | 已覆盖 | 已解析源/目标 flag 和 position。 |
| 72 | 移动物品与协议入口 | 背包到背包移动 | `gObjInventoryMoveItem` | `MoveItem` SrcFlag=0 DstFlag=0 | 已覆盖 | 支持空位移动和堆叠。 |
| 73 | 移动物品与协议入口 | 背包到仓库移动 | `gObjInventoryMoveItem` 仓库目标 | `MoveItem` SrcFlag=0 DstFlag=2 | 已覆盖 | 已检查仓库 flags 并移动。 |
| 74 | 移动物品与协议入口 | 仓库到背包移动 | `gObjWarehouse...` | `MoveItem` SrcFlag=2 DstFlag=0 | 已覆盖 | 已检查背包 flags 并移动。 |
| 75 | 移动物品与协议入口 | 仓库到仓库移动 | `pWarehouseMap` 移动 | `MoveItem` SrcFlag=2 DstFlag=2 | 已覆盖 | 已支持仓库内部移动。 |
| 76 | 移动物品与协议入口 | 堆叠移动 | `IsOverlapItem` | `MoveItem` overlap 分支 | 部分覆盖 | 只覆盖背包到背包堆叠，仓库堆叠未覆盖。 |
| 77 | 移动物品与协议入口 | 移动失败回滚 | `gObjInventoryRollback` | 当前靠先检查后变更 | 部分覆盖 | 简单路径安全，但复杂移动/交换/堆叠应增加事务式回滚。 |
| 78 | 移动物品与协议入口 | 移动后耐久同步 | `GCItemDurSend` | `MsgItemDurabilityReply` | 已覆盖 | 堆叠导致 durability 变化时会下发。 |
| 79 | 移动物品与协议入口 | 移动后装备重算 | `gObjInventoryEquipment` | `EquipmentChanged` 调用 | 部分覆盖 | 使用特殊位置 `236` 可疑，应确认实际装备/宠物槽语义。 |
| 80 | 移动物品与协议入口 | MoveItem 结果码 | `PMSG_INVENTORYITEMMOVE_RESULT` | `MsgMoveItemReply` | 部分覆盖 | 结果值目前较粗，需对齐客户端失败原因。 |
| 81 | 使用/消耗/耐久修理 | UseItem 入口 | `CGUseItemRecv` | `Object.UseItem` | 部分覆盖 | 已按 code 分支，很多消耗品仍是空 case。 |
| 82 | 使用/消耗/耐久修理 | HP 药水 | `gObjUseDrink` 等 | `Code(14,0..3)` | 已覆盖 | 延迟恢复 HP 并扣耐久。 |
| 83 | 使用/消耗/耐久修理 | MP 药水 | `gObjUseDrink` 等 | `Code(14,4..6)` | 已覆盖 | 立即恢复 MP 并推送 MP/AG。 |
| 84 | 使用/消耗/耐久修理 | SD 药水 | SD potion 使用逻辑 | `Code(14,35..37)` | 已覆盖 | 延迟恢复 SD 并扣耐久。 |
| 85 | 使用/消耗/耐久修理 | 复合药水 | compound potion | `Code(14,38..40)` | 已覆盖 | 同时延迟恢复 HP 和 SD。 |
| 86 | 使用/消耗/耐久修理 | 技能书学习 | `SkillKeyRecv`/item skill learn | `KindASkill` 分支 | 部分覆盖 | 已学习技能并删除道具，缺职业/等级/重复学习完整提示。 |
| 87 | 使用/消耗/耐久修理 | 宠物启用 UsePet | `CGInventoryEquipment`、宠物装备 | `Player.UsePet` | 部分覆盖 | 可启用/取消宠物，但状态字段需拆分。 |
| 88 | 使用/消耗/耐久修理 | 耐久扣减 | `CItem::DurabilityDown*` | `decreaseItemDurability` | 部分覆盖 | 仅消耗品简单 -1，缺攻击/受击/幸运/宠物/武器分类扣耐久。 |
| 89 | 使用/消耗/耐久修理 | 单件修理 | `ItemRepair` | `RepairItem` 指定 position | 已覆盖 | 使用 `CalculateRepairMoney` 扣钱并恢复耐久。 |
| 90 | 使用/消耗/耐久修理 | 全部修理 | `AllRepair` | `RepairItem` position `0xFF` | 已覆盖 | 遍历可修理物品，失败时停止。 |
| 91 | 掉落与地面物品生成 | 怪物掉落表加载 | `CMonsterItemMng::LoadMonsterItemDropRate` | `DropManager.init` | 部分覆盖 | 已读取 `IGC_MonsterItemDropRate.xml`，需补更多 GameServer 掉落规则。 |
| 92 | 掉落与地面物品生成 | 普通物品池 | `NormalGiveItemSearch` | `dropManager.makeItem` | 已覆盖 | 按 monster level 生成普通候选池。 |
| 93 | 掉落与地面物品生成 | 魔法书池 | `MagicBookGiveItemSearch` | `dropManager.makeMagicBook` | 已覆盖 | 技能书按怪物等级进入候选池。 |
| 94 | 掉落与地面物品生成 | 宝石池 | `MakeJewelItem` | `dropManager.makeJewel` | 已覆盖 | 已生成基础宝石候选。 |
| 95 | 掉落与地面物品生成 | 普通掉落选择 | `CMonsterItemMng::GetItem` | `DropManager.DropItem` | 部分覆盖 | 可随机物品并 `Calc`，缺完整掉落概率配置融合。 |
| 96 | 掉落与地面物品生成 | 卓越掉落选择 | `GetItemExcel` | `DropItemExcellent` | 部分覆盖 | 已接入优秀随机，但需要对齐优秀数量和类别限制。 |
| 97 | 掉落与地面物品生成 | 地图物品创建 | `MapClass::ItemDrop`、`CMapItem::CreateItem` | `MapManager.AddItem` | 部分覆盖 | 地图容器在地图系统；本模块负责实例生成与属性。 |
| 98 | 掉落与地面物品生成 | Zen 掉落生成 | `CItemDrop::GetZenAmount` | `maps.GetZen`、怪物掉落调用 | 部分覆盖 | Zen 计算在地图系统，怪物掉落处应统一调用并配置化。 |
| 99 | 掉落与地面物品生成 | Bag 掉落系统 | `CBag`、`CBagManager`、`CMonsterBag` | 当前无等价系统 | 未覆盖 | 各类宝箱/事件包/怪物包需要独立实现。 |
| 100 | 掉落与地面物品生成 | AppointItemDrop | `CAppointItemDrop::AppointItemDrop` | 当前无 | 未覆盖 | 指定怪物/地图/事件掉落需要后续补齐。 |
| 101 | 卓越/幸运/套装/380 | 卓越配置加载 | `ItemAddOption`、优秀掉落配置 | `ExcellentDropManager.init` | 已覆盖 | Go 已读取优秀掉落配置。 |
| 102 | 卓越/幸运/套装/380 | 卓越数量随机 | `CalcExcOption` | `dropExcellentCount` | 已覆盖 | 需要用概率测试验证分布。 |
| 103 | 卓越/幸运/套装/380 | 卓越属性随机 | `CalcExcOption` | `DropExcellent` | 已覆盖 | 按物品 kindA/kindB 选择可用优秀 bit。 |
| 104 | 卓越/幸运/套装/380 | 幸运字段 | `m_Option2` | `Item.Lucky` | 部分覆盖 | 字段存在并参与价格，掉落/宝石强化规则待补。 |
| 105 | 卓越/幸运/套装/380 | 追加字段 | `m_Option3` | `Item.Addition`、`Item.Calc` | 部分覆盖 | 攻击、防御、翅膀追加已算，宝石追加流程待补。 |
| 106 | 卓越/幸运/套装/380 | 套装类型加载 | `CSetItemOption::LoadTypeInfo` | `SetManager.init` 读取 set type | 已覆盖 | 已加载套装物品到 set/tier 映射。 |
| 107 | 卓越/幸运/套装/380 | 套装效果加载 | `LoadOptionInfo` | `SetManager.init` 读取 set option | 已覆盖 | 已加载普通套装效果和 full set 效果。 |
| 108 | 卓越/幸运/套装/380 | 套装效果应用 | `GetSetOption`、`GetGetFullSetOption` | `Player.calcSetItem` | 部分覆盖 | 已接入多个效果类型，仍需覆盖所有 GameServer 套装效果。 |
| 109 | 卓越/幸运/套装/380 | 380 物品识别 | `CItemSystemFor380::Is380Item` | `Item380Manager.Is380Item` | 已覆盖 | 基础 380 识别已实现。 |
| 110 | 卓越/幸运/套装/380 | 380 效果应用 | `ApplyFor380Option` | `Player.calc380Item` | 部分覆盖 | 目前主要应用 HP/SD，需要补攻击/防御等完整 380 效果。 |
| 111 | Harmony/Socket/Pentagram/期限道具 | Harmony 配置加载 | `CJewelOfHarmonySystem::LoadScript` | `HarmonyManager` init | 部分覆盖 | Go 已建配置和随机属性框架，需确认完整 XML 字段。 |
| 112 | Harmony/Socket/Pentagram/期限道具 | Harmony 随机属性 | `_GetSelectRandomOption` | `addRandEffect` | 部分覆盖 | 已按物品类型随机，但概率和可选属性需对齐。 |
| 113 | Harmony/Socket/Pentagram/期限道具 | Harmony 强化入口 | `StrengthenItemByJewelOfHarmony` | `StrengthenItem` | 部分覆盖 | Go 只提供道具层函数，协议和宝石消耗流程未接入。 |
| 114 | Harmony/Socket/Pentagram/期限道具 | Harmony 效果应用 | `SetApplyStrengthenItem` | `Player.calc` 中部分字段 | 未覆盖 | 需要把 Harmony effect/level 转成角色属性加成。 |
| 115 | Harmony/Socket/Pentagram/期限道具 | Socket 槽位字段 | `MakeSocketSlot` | `SocketSlot1-5`、`SocketSlots` | 部分覆盖 | 字段存在，但 slot 数生成和镶嵌流程缺失。 |
| 116 | Harmony/Socket/Pentagram/期限道具 | Seed/Sphere 识别 | `IsSeedItem`、`IsSphereItem`、`IsSeedSphereItem` | `item_socket.go` 基础结构 | 未覆盖 | Go 侧几乎只有结构占位，需要完整配置和识别函数。 |
| 117 | Harmony/Socket/Pentagram/期限道具 | Socket 效果应用 | `ApplySeedSphereEffect`、`SetApplySocketEffect` | 当前无完整接入 | 未覆盖 | 需要在 `Player.calc` 中应用 socket 攻防属性。 |
| 118 | Harmony/Socket/Pentagram/期限道具 | Pentagram/Errtel 字段 | Pentagram/Errtel 系统 | `PentagramBonus`、Socket 字段复用 | 部分覆盖 | 字段和协议编码有基础，元素效果系统未完成。 |
| 119 | Harmony/Socket/Pentagram/期限道具 | LuckyItem/PeriodItem | `LuckyItemManager`、`CPeriodItemEx` | `TypeLucky`、`Period` 字段 | 未覆盖 | 期限、过期、幸运装备修理/数据库同步都需要独立实现。 |
| 120 | Harmony/Socket/Pentagram/期限道具 | 道具系统回归测试总集 | GameServer 运行验证 | 当前缺系统性 item tests | 未覆盖 | 建议覆盖编码、重算、背包占格、移动、拾取、修理、套装、380、Harmony、Socket。 |
