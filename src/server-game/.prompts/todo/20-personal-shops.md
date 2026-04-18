# 20. 个人商店系统

本模块覆盖玩家个人商店 PShop 的开店、关店、商品定价、商品列表、购买事务、宝石货币、视野广播、搜索、日志、价格持久化、离线交易和跨系统状态互斥。NPC 商店买卖、修理和售出回购归 `08-shops.md`；面对面玩家交易归 `19-trade.md`。个人商店禁用处罚、异常购买审计、搜索/开店限流归 `27-security.md`；价格持久化和跨服/离线交易外部通道归 `28-external-comm.md`，本模块负责执行业务事务。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 模块边界与入口 | PersonalStore 独立管理器 | `CPersonalStore` | 暂无独立 personal shop manager | 未覆盖 | 建立独立模块，避免继续塞在 NPC 商店系统中。 |
| 2 | 模块边界与入口 | PShop 协议分发 | `protocol.cpp` 0x3F、0xEC31、0xEC33 | `handle/c1c2.go` 有路由占位 | 部分覆盖 | 将 0x3F01/02/03/05/06/07、0xEC31、0xEC33 统一归本模块。 |
| 3 | 模块边界与入口 | 个人商店总开关 | `gDoPShopOpen` | `PersonalShopEnable` | 部分覆盖 | 配置存在，需接入所有入口。 |
| 4 | 模块边界与入口 | OffTrade 配置 | `g_ConfigRead.offtrade` | `conf.OffTrade` | 部分覆盖 | 配置已加载，需定义离线商店对象和断线保持开店语义。 |
| 5 | 对象状态与数据 | 个人商店容量 | `PSHOP_SIZE`、`PSHOP_MAP_SIZE` | `PShopSize=32`、`PShopMapSize=32` | 部分覆盖 | 常量已存在，需落地业务结构。 |
| 6 | 对象状态与数据 | 个人商店背包区域 | `PSHOP_RANGE` | `PShopRangeStart=204`、`PShopRangeEnd=236` | 部分覆盖 | 确认 204-236 与客户端个人商店区域完全对齐。 |
| 7 | 对象状态与数据 | 开店状态 | `m_bPShopOpen` | 注释字段 `isPShopOpen` | 未覆盖 | 恢复或重设对象开店状态字段。 |
| 8 | 对象状态与数据 | 交易中状态 | `m_bPShopTransaction` | 注释字段 `isPShopTransaction` | 未覆盖 | 防止并发购买同一格商品。 |
| 9 | 对象状态与数据 | 商品变更状态 | `m_bPShopItemChange` | 注释字段 `isPShopItemChange` | 未覆盖 | 卖家商品变化时刷新浏览者列表。 |
| 10 | 对象状态与数据 | 视野重绘状态 | `m_bPShopRedrawAbs` | 注释字段 `isPShopRedrawABS` | 未覆盖 | 控制个人商店视野缓存重建。 |
| 11 | 对象状态与数据 | 店名文本 | `m_szPShopText`、`szPShopText[36/37]` | 注释字段 `PShopText` | 未覆盖 | 需处理长度、编码、空字符串和广播。 |
| 12 | 对象状态与数据 | 浏览会话状态 | `m_bPShopWantDeal` | 注释字段 `isPShopWantDeal` | 未覆盖 | 买家正在浏览某个个人商店时记录状态。 |
| 13 | 对象状态与数据 | 浏览卖家编号 | `m_iPShopDealerIndex` | 注释字段 `PShopDealerIndex` | 未覆盖 | 浏览列表、购买、关闭交易时使用。 |
| 14 | 对象状态与数据 | 浏览卖家名称 | `m_szPShopDealerName` | 注释字段 `PShopDealerName` | 未覆盖 | 用于日志、验证和断线清理。 |
| 15 | 对象状态与数据 | 个人商店互斥锁 | 临界区和交易标记 | 注释字段 `muPShopTrade` | 未覆盖 | Go 侧需要 mutex 或游戏单协程串行事务。 |
| 16 | 价格设置 | 设置价格请求 | `CGPShopReqSetItemPrice` | `handle` 0x3F01 `shopItemSetPrice` | 未覆盖 | 解析物品位置、Zen、Bless、Soul、Chaos 价格。 |
| 17 | 价格设置 | 设置价格回复 | `CGPShopAnsSetItemPrice` | 当前无完整消息结构 | 未覆盖 | 对齐成功、关闭、位置错误、物品不存在、价格错误等结果码。 |
| 18 | 价格设置 | 等级限制 | `Level is Under 6` | 当前无 | 未覆盖 | 低等级不能设置个人商店价格。 |
| 19 | 价格设置 | 位置范围校验 | `Item Position Out of Bound` | `PShopRangeStart/End` | 未覆盖 | 只能给个人商店区域物品定价。 |
| 20 | 价格设置 | 物品存在校验 | `Item Does Not Exist` | 当前无 | 未覆盖 | 空格不能设置价格。 |
| 21 | 价格设置 | 价格正数校验 | `m_iPShopValue`、宝石价格 | 当前无 | 未覆盖 | Zen 和宝石价格不能全部为 0 或负数。 |
| 22 | 价格设置 | 0 serial 物品限制 | Anti-hack Serial 0 Item | `EnableTrade0SerialItem` | 未覆盖 | 禁止异常序列物品上架或记录审计。 |
| 23 | 价格设置 | 全卓越出售限制 | `CanSellInStoreFullExcItem` | `EnableSellFullExcItemInPShop` | 未覆盖 | 配置存在，需限制全卓越物品上架。 |
| 24 | 价格设置 | Item Block 限制 | `m_cAccountItemBlock` | 暂无账号物品锁 | 未覆盖 | 账号物品锁定时禁止改价。 |
| 25 | 价格设置 | 价格写入物品 | `m_iPShopValue`、`m_wPShop*Value` | 当前 Item 未完整承载 | 未覆盖 | 物品结构需承载个人商店价格。 |
| 26 | 价格设置 | 改价日志 | `Changing Item Price` | 当前无 | 未覆盖 | 记录账号、角色、物品、价格和 serial。 |
| 27 | 开关店 | 开店请求 | `CGPShopReqOpen` | `handle` 0x3F02 `shopOpen` | 未覆盖 | 校验状态后开启个人商店。 |
| 28 | 开关店 | 开店回复 | `CGPShopAnsOpen` | 当前无完整消息结构 | 未覆盖 | 对齐成功、失败、等级不足、状态冲突等结果码。 |
| 29 | 开关店 | 关店请求 | `CGPShopReqClose` | `handle` 0x3F03 `shopClose` | 未覆盖 | 关闭个人商店并清理状态。 |
| 30 | 开关店 | 关店回复 | `CGPShopAnsClose` | 当前无完整消息结构 | 未覆盖 | 通知客户端关店结果。 |
| 31 | 开关店 | 空店检查 | `PShop_CheckInventoryEmpty` | 当前无 | 未覆盖 | 没有上架商品时不能开店或应自动关店。 |
| 32 | 开关店 | 店名修改 | 已开店时更新 `m_szPShopText` | 当前无 | 未覆盖 | 已开店重复 open 可改店名并触发视野更新。 |
| 33 | 开关店 | 接口状态限制 | `m_IfState` 检查 | InterfaceState 未完整接入 | 未覆盖 | 交易、仓库、合成、NPC 商店等状态下禁止开店。 |
| 34 | 开关店 | Transaction 限制 | `Transaction == 1` | 当前无 | 未覆盖 | 正在事务中不能开关店。 |
| 35 | 开关店 | 地图限制 | `PShopOpen MapNumber` 等配置 | OffTrade Map 配置 | 未覆盖 | 定义哪些地图允许开店或离线开店。 |
| 36 | 开关店 | 开店保存价格 | `GDAllSavePShopItemValue` | 当前无 | 未覆盖 | 开店时保存所有上架物品价格。 |
| 37 | 开关店 | 关店清理店名 | `memset m_szPShopText` | 当前无 | 未覆盖 | 关店后清空店名和浏览状态。 |
| 38 | 商品列表 | 请求商品列表 | `CGPShopReqBuyList` | `handle` 0x3F05 `shopItemList` | 未覆盖 | 买家请求卖家个人商店商品。 |
| 39 | 商品列表 | 商品列表回复 | `CGPShopAnsBuyList` | 当前无完整消息结构 | 未覆盖 | 返回店名、商品、价格和结果码。 |
| 40 | 商品列表 | 卖家在线检查 | Seller connected check | 当前无 | 未覆盖 | 卖家不存在或不是角色时返回失败。 |
| 41 | 商品列表 | 买家在线检查 | Buyer character check | 当前无 | 未覆盖 | 买家非 Playing 状态不能浏览。 |
| 42 | 商品列表 | 自购限制 | Requested to Him/HerSelf | 当前无 | 未覆盖 | 不允许浏览自己的个人商店。 |
| 43 | 商品列表 | 卖家开店检查 | `m_bPShopOpen == false` | 当前无 | 未覆盖 | 卖家未开店时返回店铺关闭。 |
| 44 | 商品列表 | Item Block 检查 | seller item block | 当前无 | 未覆盖 | 卖家物品锁时不能被浏览或购买。 |
| 45 | 商品列表 | 浏览者事务检查 | requester Transaction | 当前无 | 未覆盖 | 买家处于交易/合成等事务时不能浏览。 |
| 46 | 商品列表 | 浏览会话绑定 | `m_bPShopWantDeal`、`m_iPShopDealerIndex` | 当前无 | 未覆盖 | 商品列表成功后记录买家浏览状态。 |
| 47 | 商品列表 | 列表重发 | `bResend` | 当前无 | 未覆盖 | 卖家商品变化时向浏览者刷新列表。 |
| 48 | 购买事务 | 购买请求 | `CGPShopReqBuyItem` | `handle` 0x3F06 `shopItemBuy` | 未覆盖 | 解析卖家、商品位置和购买信息。 |
| 49 | 购买事务 | 购买回复 | `CGPShopAnsBuyItem` | 当前无完整消息结构 | 未覆盖 | 对齐成功、店铺关闭、背包满、货币不足等结果码。 |
| 50 | 购买事务 | 商品位置校验 | `btItemPos` | 当前无 | 未覆盖 | 只能购买卖家 PShop 区域有效物品。 |
| 51 | 购买事务 | 卖家交易锁 | `m_bPShopTransaction` | 当前无 | 未覆盖 | 购买开始前锁定卖家，结束后释放。 |
| 52 | 购买事务 | 价格重验 | `m_iPShopValue`、宝石价格 | 当前无 | 未覆盖 | 使用服务器端物品价格，不能信任客户端。 |
| 53 | 购买事务 | 买家 Zen 校验 | Zen price check | `Object.Money` | 未覆盖 | 买家 Zen 不足时失败。 |
| 54 | 购买事务 | 卖家 Zen 上限 | Exceeding Zen of the Host | `MaxZen` | 未覆盖 | 卖家收款后不能超过 Zen 上限。 |
| 55 | 购买事务 | 买家背包空间 | No Room to Buy Item | `Inventory.FindFreePositionForItem` | 未覆盖 | 购买前确认买家背包可放入物品。 |
| 56 | 购买事务 | 卖家物品移除 | `gObjInventoryItemSet_PShop(...,-1)` | 当前无 | 未覆盖 | 成交后从卖家 PShop 区域移除物品和占位。 |
| 57 | 购买事务 | 买家获得物品 | `gObjInventoryInsertItem` | 当前无完整 PShop 流程 | 未覆盖 | 成交后将商品插入买家背包。 |
| 58 | 购买事务 | 卖家收款 | Zen/Jewel 收款 | 当前无 | 未覆盖 | 卖家获得 Zen 或宝石货币。 |
| 59 | 购买事务 | 自动关店 | `PShop_CheckInventoryEmpty` 后 close | 当前无 | 未覆盖 | 商品售空后自动关闭个人商店。 |
| 60 | 购买事务 | 成交通知卖家 | `CGPShopAnsSoldItem` | 当前无 | 未覆盖 | 通知卖家某格商品已售出。 |
| 61 | 购买事务 | 交易失败回滚 | 多分支失败返回 | 当前无 | 未覆盖 | 任一步失败必须恢复物品、货币和交易锁。 |
| 62 | 宝石货币 | Bless 价格 | `m_wPShopBlessValue` | 当前无 | 未覆盖 | 支持祝福宝石计价。 |
| 63 | 宝石货币 | Soul 价格 | `m_wPShopSoulValue` | 当前无 | 未覆盖 | 支持灵魂宝石计价。 |
| 64 | 宝石货币 | Chaos 价格 | `m_wPShopChaosValue` | 当前无 | 未覆盖 | 支持玛雅宝石计价。 |
| 65 | 宝石货币 | 宝石不足失败 | Lack of Bless/Soul/Chaos | 当前无 | 未覆盖 | 买家对应宝石不足时返回明确结果码。 |
| 66 | 宝石货币 | 宝石不可拆分检查 | Not Share Bless/Soul/Chaos | 当前无 | 未覆盖 | 无法拆分或找零时返回失败。 |
| 67 | 宝石货币 | 卖家宝石背包空间 | Seller Inventory is Full | 当前无 | 未覆盖 | 卖家收宝石前确认背包空间。 |
| 68 | 宝石货币 | 宝石手续费 | `BJ/SJ/CJ Commision` 日志 | 当前无 | 未覆盖 | 如保留手续费，需要统一配置和日志。 |
| 69 | 视野广播 | 视野列表重建 | `PShop_ViewportListRegenarate` | 当前无 | 未覆盖 | 开店、关店、改名、进入视野时同步店铺状态。 |
| 70 | 视野广播 | 视野 PShop 缓存 | `m_iVpPShopPlayer`、`m_wVpPShopPlayerCount` | 注释字段 `VPPShopPlayer` | 未覆盖 | 记录当前视野内已开店玩家，避免重复广播。 |
| 71 | 视野广播 | 开店出现消息 | `PMSG_ANS_PSHOP_VIEWPORT_NOTIFY` | 当前无 | 未覆盖 | 向视野内玩家通知店名和开店状态。 |
| 72 | 视野广播 | 关店清理消息 | DealerClosedShop/viewport update | 当前无 | 未覆盖 | 关店时通知浏览者和视野玩家。 |
| 73 | 视野广播 | 浏览者卖家关闭检测 | `m_bPShopWantDeal` 分支 | 当前无 | 未覆盖 | 浏览中卖家关闭、离线、商品变动时通知买家。 |
| 74 | 搜索与日志 | 搜索个人商店请求 | `CGReqSearchItemInPShop` | `handle` 0xEC31 `reqSearchItemInPShop` | 未覆盖 | 处理搜索入口。 |
| 75 | 搜索与日志 | 全部商店搜索 | `GCPShop_AllInfo` | 当前无 | 未覆盖 | 分页返回所有开店玩家。 |
| 76 | 搜索与日志 | 按物品搜索 | `GCPShop_SearchItem` | 当前无 | 未覆盖 | 根据物品 type 搜索包含该物品的店铺。 |
| 77 | 搜索与日志 | 搜索分页 | `iLastUserCount`、`iPShopCnt` | 当前无 | 未覆盖 | 支持 LastCount 分页和结果为空。 |
| 78 | 搜索与日志 | 搜索命中检查 | `PShop_CheckExistItemInInventory` | 当前无 | 未覆盖 | 检查卖家 PShop 区域是否存在指定物品。 |
| 79 | 搜索与日志 | PShopLog 请求 | `CGReqPShopLog` | `handle` 0xEC33 `reqPShopLog` | 未覆盖 | 支持客户端查询或触发店铺日志动作。 |
| 80 | 搜索与日志 | PShopLog 私聊卖家 | `PShopLog Whisper To` | `Object.Whisper` 可辅助 | 未覆盖 | 搜索结果中向卖家私聊的日志/入口。 |
| 81 | 搜索与日志 | PShopLog 邮件卖家 | `PShopLog Mail To` | `18-friends-mail.md` 待接入 | 未覆盖 | 搜索结果中给卖家发邮件的日志/入口。 |
| 82 | 搜索与日志 | 成交审计日志 | `PShop Item Buy Request Succeed` | 当前无 | 未覆盖 | 记录买卖双方、价格、手续费、物品、serial。 |
| 83 | 价格持久化 | 请求价格持久化 | `GDRequestPShopItemValue` | 当前无 | 未覆盖 | 登录或加载角色后通过 `28-external-comm.md` 请求个人商店物品价格。 |
| 84 | 价格持久化 | 更新价格持久化 | `GDUpdatePShopItemValue` | 当前无 | 未覆盖 | 单个物品价格变化时保存。 |
| 85 | 价格持久化 | 保存全部价格 | `GDAllSavePShopItemValue` | 当前无 | 未覆盖 | 开店时保存全部上架物品价格。 |
| 86 | 价格持久化 | 删除价格记录 | `GDDelPShopItemValue` | 当前无 | 未覆盖 | 物品下架、售出或删除时移除价格记录。 |
| 87 | 价格持久化 | 移动价格记录 | `GDMovePShopItem` | 当前无 | 未覆盖 | PShop 区域物品移动时同步旧位置和新位置。 |
| 88 | 价格持久化 | 接收价格响应 | `GDAnsPShopItemValue` | 当前无 | 未覆盖 | 根据 serial 和位置恢复物品价格。 |
| 89 | 价格持久化 | 价格信息下发 | `GCPShopItemValueInfo` | 当前无 | 未覆盖 | 向客户端下发每格商品价格。 |
| 90 | 离线交易 | 离线开店地图 | OffTrade Map | `conf.OffTrade.Map` | 未覆盖 | 定义允许离线交易的地图。 |
| 91 | 离线交易 | 离线货币类型 | OffTrade CoinType | `conf.OffTrade` | 未覆盖 | 影响个人商店浏览时 `GCAlterPShopVault`。 |
| 92 | 离线交易 | 断线保持开店 | OffTrade flow | 当前无 | 未覆盖 | 角色断线后是否保留对象和店铺状态。 |
| 93 | 离线交易 | Bot 商店货币 | BotSystem shop coin type | Bot 系统未实现 | 未覆盖 | GameServer 有 Bot 商店货币分支，Go 侧先记录边界。 |
| 94 | 跨系统限制 | 与玩家交易互斥 | `m_bPShopOpen` trade checks | `19-trade.md` 待接入 | 未覆盖 | 开店或浏览个人商店时不能面对面交易。 |
| 95 | 跨系统限制 | 与合成互斥 | ChaosBox PShop checks | `12-mix.md` 待接入 | 未覆盖 | 开店时不能使用 ChaosBox，使用 ChaosBox 时不能开店。 |
| 96 | 跨系统限制 | 与 NPC 商店互斥 | interface state | `08-shops.md` 边界 | 未覆盖 | NPC 商店和个人商店状态互斥。 |
| 97 | 跨系统限制 | 与仓库互斥 | vault/interface state | `07-items.md`/仓库边界 | 未覆盖 | 仓库打开时禁止开店或购买。 |
| 98 | 跨系统限制 | 地图移动限制 | MapServerMove checks | `06-maps.md` 待接入 | 未覆盖 | 地图移动、传送、换线时应关闭或禁止个人商店。 |
| 99 | 跨系统限制 | 角色下线清理 | close flow calls `CGPShopReqClose` | `03-characters.md`、`04-objects.md` 待接入 | 未覆盖 | 下线时关闭店铺或进入离线交易。 |
| 100 | 跨系统限制 | 物品移动限制 | `GDMovePShopItem`、占位图 | `07-items.md` 待接入 | 未覆盖 | 个人商店区域物品移动必须同步价格和占位。 |
| 101 | 协议与测试 | PShop opcode 覆盖 | 0x3F01/02/03/05/06/07、0xEC31/33 | `handle/c1c2.go` | 部分覆盖 | 确认每个 opcode 有明确处理函数和稳定失败结果。 |
| 102 | 协议与测试 | 设置价格测试 | `CGPShopReqSetItemPrice` | 待实现 | 未覆盖 | 覆盖位置越界、空格、价格为 0、禁售物品、成功。 |
| 103 | 协议与测试 | 开关店测试 | `CGPShopReqOpen/Close` | 待实现 | 未覆盖 | 覆盖等级不足、空店、接口冲突、改店名、关店。 |
| 104 | 协议与测试 | 商品列表测试 | `CGPShopReqBuyList` | 待实现 | 未覆盖 | 覆盖卖家离线、未开店、自购、浏览状态、刷新。 |
| 105 | 协议与测试 | 购买事务测试 | `CGPShopReqBuyItem` | 待实现 | 未覆盖 | 覆盖并发购买、背包满、Zen 不足、卖家上限、失败回滚。 |
| 106 | 协议与测试 | 宝石交易测试 | Bless/Soul/Chaos price | 待实现 | 未覆盖 | 覆盖宝石不足、不可拆分、卖家背包满、手续费。 |
| 107 | 协议与测试 | 搜索分页测试 | `GCPShop_AllInfo/SearchItem` | 待实现 | 未覆盖 | 覆盖全部搜索、按物品搜索、分页、空结果。 |
| 108 | 协议与测试 | 持久化测试 | GD PShop item value flow | 待实现 | 未覆盖 | 覆盖价格保存、移动、删除、重登恢复。 |
| 109 | 协议与测试 | 离线交易测试 | OffTrade | 待实现 | 未覆盖 | 覆盖断线保持、地图限制、货币类型、重新登录。 |
| 110 | 协议与测试 | 互斥状态测试 | Trade/Shop/Warehouse/Chaos/Move | 待实现 | 未覆盖 | 覆盖与交易、NPC 商店、仓库、合成、地图移动的互斥。 |
