# 8. 商店系统

本模块覆盖 NPC 商店配置、库存、NPC 绑定、Talk 打开、购买、出售、修理、价格限制、售出回购、特殊商店、协议与测试。个人商店已独立为 `20-personal-shops.md`；现金商城只作为协议边界记录，不在本模块展开完整业务。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | NPC 商店配置与加载 | 商店列表配置加载 | `CShopMng::LoadShopList` | `shopManager.init` 读取 `IGC_ShopList.xml` | 部分覆盖 | 已加载 NPC、地图、坐标和文件名；需补加载失败策略、重复项诊断和配置热重载边界。 |
| 2 | NPC 商店配置与加载 | 商店数据结构 | `SHOP_DATA`、`CShop` | `Shop` | 部分覆盖 | Go 侧已有 NPCIndex、MapNumber、PosX、PosY、Dir、VIP/GM/Moss/BattleCore 字段；特殊字段多数未参与业务判定。 |
| 3 | NPC 商店配置与加载 | NPC+地图索引表 | `CShopMng::GetShop` | `shopTable map[int]map[int]*Shop` | 部分覆盖 | 需将未知 NPC 或地图的查询改为 nil-safe，避免直接索引导致 panic。 |
| 4 | NPC 商店配置与加载 | 共享商店库存文件 | `CShopMng::LoadShopList` 复用商店数据 | `shopByName` 复用同 FileName 的 `PositionedItems` | 已覆盖 | 需用测试固定同一 shop 文件被多个 NPC 共享库存的语义。 |
| 5 | NPC 商店配置与加载 | 商店 XML 商品节点解析 | `CShop::LoadShopItem` | `Shop.Scan` 解析 `<Item>` 属性 | 部分覆盖 | 已解析 Cat、Index、Level、Durability、Skill、Luck、Option、Exc 等；需补完整字段语义。 |
| 6 | NPC 商店配置与加载 | 商品实例创建 | `CItem::Convert`、`CShop::InsertItem` | `item.NewItem` | 部分覆盖 | 当前创建基础物品后再设置等级、耐久、技能、幸运、追加、优秀；需对齐 GameServer 初始化顺序。 |
| 7 | NPC 商店配置与加载 | 商品等级与耐久应用 | `CShop::LoadShopItem` | `SetLevel`、`Durability`、`DefaultDurability` | 部分覆盖 | 需明确配置耐久为 0 时使用默认耐久，非 0 时是否完全信任配置。 |
| 8 | NPC 商店配置与加载 | 商品优秀属性解码 | `m_NewOption` | `DecodeExcellent` | 部分覆盖 | 当前支持 Exc bit 解码；需确认商店优秀字段和物品协议编码完全一致。 |
| 9 | NPC 商店配置与加载 | 套装/Socket/元素/序列字段 | `CShop::LoadShopItem` 商品扩展字段 | `Shop.Scan` 已读 `SetItem`、`SocketCount`、`Elemental`、`Serial` | 未覆盖 | Go 侧读取但未应用这些字段，需要明确是否生成套装、Socket、元素属性和固定 serial。 |
| 10 | NPC 商店配置与加载 | 商店加载诊断 | `LogAdd`、配置错误日志 | 当前主要 `panic` 或静默跳过 | 未覆盖 | 需要把坏配置、无空位、未知物品、重复商店等错误转成可定位日志和测试。 |
| 11 | 商店库存与格子 | 商店最大格子数 | `MAX_ITEM_IN_SHOP=120` | `MaxShopItemCount=120` | 已覆盖 | 8x15 商店容量已对齐。 |
| 12 | 商店库存与格子 | 商店 8 列布局 | `ShopInventoryMap[120]` | `PositionedItems` Size=120 | 已覆盖 | 需用测试锁定位置到行列的换算。 |
| 13 | 商店库存与格子 | 商店占位图初始化 | `CShop::Init` | `PositionedItems.Flags` | 部分覆盖 | Go 侧依赖 `NewPositionedItems` 初始化；需确认每个商店加载前 flags 被清空。 |
| 14 | 商店库存与格子 | 商店空位查找 | `CShop::InentoryMapCheck` | `FindFreePositionForItem` | 已覆盖 | 当前扫描 8x15 并校验物品宽高。 |
| 15 | 商店库存与格子 | 商店占格写入 | `CShop::InsertItem` | `SetFlagsForItem` | 已覆盖 | 需补非法 position 或尺寸越界的保护。 |
| 16 | 商店库存与格子 | 商品重叠校验 | `ShopInventoryMap` | `FindFreePositionForItem` flags 检查 | 已覆盖 | 当前拒绝已占用格；需加入回归测试。 |
| 17 | 商店库存与格子 | 商品尺寸越界校验 | `InventoryMapCheck` | `FindFreePositionForItem` 宽高检查 | 部分覆盖 | 需补物品宽高为 0、超过 8 列、超过 15 行时的错误处理。 |
| 18 | 商店库存与格子 | 固定库存语义 | `CShop` 固定售卖列表 | `Shop.Items` | 已覆盖 | NPC 商店是静态库存，购买不会从商店移除商品。 |
| 19 | 商店库存与格子 | 动态库存边界 | GameServer 固定商店为主 | 当前无动态库存 | 未覆盖 | 如果后续支持限量或动态刷新，应独立设计，不混入当前固定商店路径。 |
| 20 | 商店库存与格子 | 商店库存热重载 | GameServer 启动加载为主 | 当前无热重载 | 未覆盖 | 需决定是否支持运行时 reload；若不支持，应在文档和配置层明确。 |
| 21 | 商店 NPC 生成与绑定 | 商店 NPC 枚举 | `CShopMng::SetShopNpcs` | `ShopManager.ForEachShop` | 部分覆盖 | Go 侧可枚举 class/map/x/y/dir，用于生成 NPC；需确认调用链覆盖所有商店。 |
| 22 | 商店 NPC 生成与绑定 | NPC 类型设置 | `IsShopNpc`、NPC talk 路由 | `NpcTypeShop` | 部分覆盖 | 需确保由商店配置生成的 NPC 都设置为 `NpcTypeShop`。 |
| 23 | 商店 NPC 生成与绑定 | NPC class 绑定 | `SHOP_DATA.iNpcIndex` | `Shop.NPCIndex`、`Object.Class` | 已覆盖 | 商店通过 NPC class 查表。 |
| 24 | 商店 NPC 生成与绑定 | NPC 地图绑定 | `SHOP_DATA.btMapNumber` | `Shop.MapNumber`、`Object.MapNumber` | 已覆盖 | 同 class 可在不同地图绑定不同商店。 |
| 25 | 商店 NPC 生成与绑定 | NPC 坐标与方向 | `btPosX/btPosY/btDir` | `PosX`、`PosY`、`Dir` | 已覆盖 | 需确认生成怪物/NPC 时使用这些坐标和方向。 |
| 26 | 商店 NPC 生成与绑定 | GM 商店绑定 | `btOnlyForGameMaster` | `Shop.GMShop` | 未覆盖 | Go 侧字段存在但购买/打开时未限制 GM 权限。 |
| 27 | 商店 NPC 生成与绑定 | VIP 商店绑定 | `btReqVipLevel` | `VipType`、`VIPType` | 未覆盖 | Go 侧字段存在但未校验玩家 VIP 等级，且字段命名重复需统一。 |
| 28 | 商店 NPC 生成与绑定 | Moss 商人绑定 | `btIsMossMerchant` | `Shop.MossMerchant` | 未覆盖 | 需明确 Moss 商人是普通商店、抽奖商店还是特殊协议。 |
| 29 | 商店 NPC 生成与绑定 | BattleCore 商店绑定 | `btBattleCoreShop` | `Shop.BattleCore` | 未覆盖 | 需按 GameServer 规则限制 BattleCore/跨服场景。 |
| 30 | 商店 NPC 生成与绑定 | 商店存在校验 | `CShopMng::GetShop` 返回空判断 | `GetShopInventory`、`GetShopItem` | 未覆盖 | Go 侧需在 Talk/Buy 前安全处理“NPC 是商店类型但无商店配置”。 |
| 31 | Talk 打开商店 | Talk 请求解析 | `CGNPCChatRecv`、`NpcTalk` | `MsgTalk.Unmarshal` | 已覆盖 | 已解析目标编号。 |
| 32 | Talk 打开商店 | Talk 目标对象获取 | `gObj` 目标索引 | `MapManager.GetObject` | 部分覆盖 | Go 侧已获取目标对象；需处理目标离线、跨图、已销毁对象。 |
| 33 | Talk 打开商店 | Talk 距离校验 | `NpcTalk` 距离和状态检查 | `Distance > 5` 返回 | 已覆盖 | 需用测试覆盖超过距离不能打开商店。 |
| 34 | Talk 打开商店 | 移动中位置同步 | GameServer talk 前位置修正 | `if p.Moving { SetPosition(TX,TY) }` | 已覆盖 | 需确认该行为不会绕过地图阻挡或瞬移校验。 |
| 35 | Talk 打开商店 | 当前 NPC 目标保存 | `lpObj->TargetNumber` | `p.TargetNumber` | 已覆盖 | Buy/Sell 依赖该字段确认商店目标。 |
| 36 | Talk 打开商店 | NpcTypeShop 路由 | `NpcTalk` 分发 | `switch tobj.NpcType` | 部分覆盖 | 已路由商店、仓库等；商店特殊类型仍未细分。 |
| 37 | Talk 打开商店 | TalkReply 结果码 | `PMSG_NPC_TALK_SEND` | `MsgTalkReply{Result:34}` | 部分覆盖 | 需核对不同商店 NPC 的 result 编号和客户端 UI。 |
| 38 | Talk 打开商店 | 商品列表下发 | `CShop::SendItemData` | `MsgTypeItemListReply{Type:0, Items: inventory}` | 部分覆盖 | 当前每次根据 Items marshal；GameServer 有预编码 SendItemData，可考虑缓存。 |
| 39 | Talk 打开商店 | 关闭商店窗口 | NPC close talk handler | `CloseTalkWindow` 设置 `TargetNumber=-1` | 部分覆盖 | 需同步清理 NPC 商店接口状态和交易状态；个人商店会话清理由 `20-personal-shops.md` 承担。 |
| 40 | Talk 打开商店 | 接口状态限制 | `m_IfState`、ChaosBox/PShop 等冲突 | 当前缺少统一 interface state | 未覆盖 | 需防止打开 NPC 商店时同时交易、合成、仓库或个人商店。 |
| 41 | NPC 购买流程 | 购买请求解析 | `CGShopBuyRecv` | `MsgBuyItem.Unmarshal` | 已覆盖 | 已解析商店位置。 |
| 42 | NPC 购买流程 | 购买位置范围校验 | `MAX_ITEM_IN_SHOP` | `msg.Position < 0 || >= MaxShopItemCount` | 已覆盖 | 返回失败结果 `-1`。 |
| 43 | NPC 购买流程 | 购买目标存在校验 | `TargetNumber`、NPC object | `TargetNumber`、`MapManager.GetObject` | 部分覆盖 | 已校验 target；需处理目标变更、消失、跨图。 |
| 44 | NPC 购买流程 | 购买距离校验 | GameServer NPC 距离限制 | `Distance > 5` | 已覆盖 | 需确认与 Talk 距离一致。 |
| 45 | NPC 购买流程 | 购买 NPC 类型校验 | `IsShopNpc` | `tobj.NpcType != NpcTypeShop` | 部分覆盖 | 只检查类型，未验证 shopTable 中确实存在配置。 |
| 46 | NPC 购买流程 | 商品复制 | `CItem` 从商店库存复制 | `ShopManager.GetShopItem` 返回 `Copy()` | 已覆盖 | 购买不会直接复用商店模板实例。 |
| 47 | NPC 购买流程 | 背包空位查找 | `gObjShopBuyInventoryInsertItem` | `Inventory.FindFreePositionForItem` | 部分覆盖 | 已查找空位/堆叠；需对齐 MU 不同物品堆叠规则。 |
| 48 | NPC 购买流程 | 余额校验 | `lpObj->Money` 比较价格 | `obj.Money < sit.BuyMoney` | 已覆盖 | 需补税率、折扣、特殊货币前置后重新计算价格。 |
| 49 | NPC 购买流程 | 扣钱与推送 | `MoneySend` | `obj.Money -= sit.BuyMoney`、`PushMoney` | 部分覆盖 | 需保证扣钱、加物品、回复包失败时有一致回滚策略。 |
| 50 | NPC 购买流程 | 购买成功回复 | `PMSG_BUY_RESULT` | `MsgBuyItemReply` | 部分覆盖 | 已返回位置和物品；需核对失败码与客户端兼容。 |
| 51 | NPC 出售流程 | 出售请求解析 | `CGShopSellRecv` | `MsgSellItem.Unmarshal` | 已覆盖 | 已解析背包位置。 |
| 52 | NPC 出售流程 | 出售位置范围校验 | 背包范围检查 | `position < 0 || >= Inventory.Size` | 已覆盖 | 越界直接失败。 |
| 53 | NPC 出售流程 | 出售目标存在校验 | `TargetNumber`、NPC object | `TargetNumber`、`MapManager.GetObject` | 部分覆盖 | 已检查目标；需处理目标消失后的状态清理。 |
| 54 | NPC 出售流程 | 出售距离校验 | GameServer NPC 距离限制 | `Distance > 5` | 已覆盖 | 与购买保持一致。 |
| 55 | NPC 出售流程 | 出售 NPC 类型校验 | `IsShopNpc` | `NpcTypeShop` | 部分覆盖 | 未进一步验证商店配置或修理/卖出权限。 |
| 56 | NPC 出售流程 | 背包物品存在校验 | `pInventory[pos].IsItem()` | `Inventory.Items[msg.Position]` | 已覆盖 | 空格不能出售。 |
| 57 | NPC 出售流程 | 不可出售物品限制 | `SellToNPC`、事件/绑定限制 | `ItemBase.SellToNPC` 存在但未校验 | 未覆盖 | 需禁止不可卖、绑定、期限、任务物品等。 |
| 58 | NPC 出售流程 | 出售价与金币上限 | `ItemValue`、Zen 上限 | `it.SellMoney`、`MaxZen` | 部分覆盖 | 已校验金币上限；需对齐价格公式和税率。 |
| 59 | NPC 出售流程 | 加钱与移除物品 | GameServer 卖出事务 | `obj.Money += it.SellMoney`、`Inventory.RemoveItem` | 部分覆盖 | 当前先加钱再删物品，需明确异常回滚和持久化时序。 |
| 60 | NPC 出售流程 | 出售回复 | `PMSG_SELL_RESULT` | `MsgSellItemReply` | 部分覆盖 | 已返回 result 和 money；需对齐所有失败码。 |
| 61 | 修理流程 | 修理请求解析 | `CGShopRepairRecv` | `MsgRepairItem.Unmarshal` | 已覆盖 | 已解析位置，`0xFF` 表示全部修理。 |
| 62 | 修理流程 | 单件修理入口 | `RepairItem` | `repairItem(it,false)` | 部分覆盖 | 已支持单格修理，但未要求玩家正与修理 NPC 交互。 |
| 63 | 修理流程 | 全部修理入口 | GameServer 全部修理 | `Position == 0xFF` 遍历装备/背包 | 部分覆盖 | 已遍历所有物品；需确认是否包含扩展背包、宠物、特殊栏位。 |
| 64 | 修理流程 | 修理价格计算 | `GetRepairMoney` | `it.CalculateRepairMoney()` | 部分覆盖 | 公式在道具侧，需和 GameServer 价格完全对齐。 |
| 65 | 修理流程 | fast 修理倍率 | GameServer repair type | `repairItem(it, fast bool)` | 部分覆盖 | 函数支持 fast，但协议调用目前固定 false。 |
| 66 | 修理流程 | 修理金币扣除 | GameServer 金币事务 | `obj.Money -= money` | 部分覆盖 | 需补扣钱失败、并发请求和持久化回滚。 |
| 67 | 修理流程 | 耐久恢复 | `m_Durability = max` | `it.Durability = it.DefaultDurability` | 部分覆盖 | 需确认套装、宠物、期限、耐久上限变化的处理。 |
| 68 | 修理流程 | 耐久同步 | `GCItemDurSend` 等 | `pushDurability` | 部分覆盖 | 已推送耐久；需确认全部修理时客户端 UI 同步完整。 |
| 69 | 修理流程 | 修理 NPC/接口状态校验 | `m_IfState`、NPC talk 状态 | 当前缺少目标商店校验 | 未覆盖 | 需禁止客户端绕过 NPC 直接发修理包。 |
| 70 | 修理流程 | 修理失败回滚 | GameServer 失败返回 | 当前局部返回失败 | 未覆盖 | 扣钱后任何失败都应恢复金币和耐久。 |
| 71 | 价格、税率与特殊商店 | 买价来源 | `CItem::m_BuyMoney` | `Item.BuyMoney` | 部分覆盖 | 基础价格由物品计算生成；需确认商店配置是否可覆盖价格。 |
| 72 | 价格、税率与特殊商店 | 卖价来源 | `CItem::m_SellMoney` | `Item.SellMoney` | 部分覆盖 | 需核对耐久、优秀、套装、Harmony、Socket 对卖价影响。 |
| 73 | 价格、税率与特殊商店 | ItemPrice 配置 | `ItemValue`/价格配置 | `conf.ItemPrice` | 未覆盖 | 配置结构存在，需要接入 NPC 买卖价；个人商店估值归 `20-personal-shops.md`。 |
| 74 | 价格、税率与特殊商店 | 城堡税率 | Castle Siege tax shop | 当前未接入 | 未覆盖 | 需决定是否在 NPC 买入/修理/出售中使用城堡税率。 |
| 75 | 价格、税率与特殊商店 | PK 使用商店限制 | `PKCanUseshops` | `conf.Common.PKCanUseshops` | 未覆盖 | 配置存在但 Talk/Buy/Sell/Repair 未校验。 |
| 76 | 价格、税率与特殊商店 | GM 商店权限 | `btOnlyForGameMaster` | `Shop.GMShop` | 未覆盖 | 需限制非 GM 打开和购买。 |
| 77 | 价格、税率与特殊商店 | VIP 商店权限 | `btReqVipLevel` | `Shop.VipType/VIPType` | 未覆盖 | 需统一字段并按账号/角色 VIP 校验。 |
| 78 | 价格、税率与特殊商店 | Moss 商店特殊规则 | Moss Merchant 逻辑 | `Shop.MossMerchant` | 未覆盖 | 需明确是否关联随机购买、特殊货币或 UI result。 |
| 79 | 价格、税率与特殊商店 | BattleCore 商店限制 | BattleCore shop flag | `Shop.BattleCore` | 未覆盖 | 需限制跨服、竞技场或 BattleCore 状态。 |
| 80 | 价格、税率与特殊商店 | 商店交易审计日志 | `LogAdd`、DB 日志 | 当前缺少统一日志 | 未覆盖 | NPC 购买、出售、修理、回购都应输出可审计日志；个人商店交易日志归 `20-personal-shops.md`。 |
| 81 | 售出回购/撤销出售 | 售出列表请求入口 | `CCancelItemSale::CGReqSoldItemList` | `handle` 0x6F00 `itemSoldList` | 未覆盖 | handler 有协议入口痕迹，业务逻辑未落地。 |
| 82 | 售出回购/撤销出售 | 售出列表回复 | `CCancelItemSale::GCAnsSoldItemList` | 当前无完整消息结构 | 未覆盖 | 需定义 Go 侧 sold item list reply 并对齐客户端格式。 |
| 83 | 售出回购/撤销出售 | 卖出物品入回购列表 | `GDReqAddItemToList` | `SellItem` 未写入回购记录 | 未覆盖 | NPC 出售成功后应按配置写入可回购列表。 |
| 84 | 售出回购/撤销出售 | 回购请求解析 | `CGReqReBuyItem` | `handle` 0x6F02 `itemsoldReBuy` | 未覆盖 | 需解析回购列表索引、物品标识和价格。 |
| 85 | 售出回购/撤销出售 | 回购 DB 获取 | `GDReqGetReBuyItem`、`DGAnsGetReBuyItem` | 当前无 DB 交互 | 未覆盖 | 需设计 sold item 持久化模型或复用角色数据存储。 |
| 86 | 售出回购/撤销出售 | 回购背包空间校验 | GameServer 回购前 inventory check | 当前无实现 | 未覆盖 | 回购物品前必须找空位并校验占格。 |
| 87 | 售出回购/撤销出售 | 回购扣钱事务 | GameServer rebuy money check | 当前无实现 | 未覆盖 | 扣钱、加物品、删除回购记录应作为一个事务处理。 |
| 88 | 售出回购/撤销出售 | 删除售出记录 | `GDReqDeleteSoldItem` | `itemSoldCancelSale` 未实现 | 未覆盖 | 取消/回购/过期都需要删除对应记录。 |
| 89 | 售出回购/撤销出售 | 售出记录过期 | `USER_SHOP_REBUY_ITEM` 时间字段 | 当前无实现 | 未覆盖 | 需实现过期时间、列表上限和清理策略。 |
| 90 | 售出回购/撤销出售 | 回购错误码 | `PMSG_ANS_REBUY_ITEM` | 当前无实现 | 未覆盖 | 需对齐余额不足、背包满、物品不存在、过期等失败码。 |
| 91 | 协议、边界与测试 | 商店 opcode 映射 | GameServer protocol dispatch | `handle/c1c2.go` 0x30-0x34、0x6F、0xD202 | 部分覆盖 | 需把 NPC 商店、售出回购和 CashShop 边界入口连接到明确业务函数；个人商店 0x3F/0xEC31/0xEC33 归 `20-personal-shops.md`。 |
| 92 | 协议、边界与测试 | 商品列表协议 | `CShop::SendItemData` | `MsgTypeItemListReply`、`MsgItemListReply` | 部分覆盖 | 需核对 count、position、item bytes 与客户端版本兼容。 |
| 93 | 协议、边界与测试 | 购买回复协议 | `PMSG_BUY_RESULT` | `MsgBuyItemReply.Marshal` | 部分覆盖 | 需补全失败码并测试成功/失败包长度。 |
| 94 | 协议、边界与测试 | 出售回复协议 | `PMSG_SELL_RESULT` | `MsgSellItemReply.Marshal` | 部分覆盖 | 需核对 result byte 和 money 字段端序。 |
| 95 | 协议、边界与测试 | 修理回复协议 | `PMSG_REPAIR_RESULT` | `MsgRepairItemReply` | 部分覆盖 | 当前复用出售回复结构，需确认协议是否完全一致。 |
| 96 | 协议、边界与测试 | CashShop 边界 | `CashShop.cpp/.h` | `handle` 0xD202 `cashShopOpen` | 未覆盖 | 本模块只记录现金商城入口，不展开 WCoin、GoblinPoint、商品包等完整现金商城。 |
| 97 | 协议、边界与测试 | NPC 商店回归测试 | `CShop`、`CShopMng` | `game/shop`、`Object.Talk` | 未覆盖 | 覆盖加载、NPC 打开、商品列表、未知商店不 panic。 |
| 98 | 协议、边界与测试 | 买卖修理事务测试 | GameServer shop buy/sell/repair | `BuyItem`、`SellItem`、`RepairItem` | 未覆盖 | 覆盖余额不足、背包满、金币上限、距离过远、失败回滚。 |
| 99 | 协议、边界与测试 | 售出回购回归测试 | `CCancelItemSale` | 待实现 0x6F 业务 | 未覆盖 | 覆盖卖出入列表、列表查询、回购成功、过期、删除、DB 失败。 |
