# 19. 玩家交易系统

本模块覆盖玩家面对面交易请求、响应、交易栏物品、交易金币、确认按钮、取消、成交、失败回滚、交易状态冲突和事务安全。NPC 商店归 `08-shops.md`，个人商店归 `20-personal-shops.md`；本模块只负责玩家与玩家之间的直接交易。交易禁用处罚、异常交易审计和请求限流归 `27-security.md`，本模块在入口查询安全状态并上报异常证据；跨服移动状态和中心服务通知归 `28-external-comm.md`。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 现有入口 | 交易请求协议 | `CGTradeRequestSend` | `handle/c1c2.go` 0x36 | 部分覆盖 | 路由存在，业务未实现。 |
| 2 | 现有入口 | 交易响应协议 | `CGTradeResponseRecv` | `handle` 0x37 | 部分覆盖 | 路由存在，业务未实现。 |
| 3 | 现有入口 | 交易金币协议 | `CGTradeMoneyRecv` | `handle` 0x3A | 部分覆盖 | 路由存在，业务未实现。 |
| 4 | 现有入口 | 交易确认协议 | `CGTradeOkButtonRecv` | `handle` 0x3C | 部分覆盖 | 路由存在，业务未实现。 |
| 5 | 现有入口 | 交易取消协议 | `CGTradeCancelButtonRecv` | `handle` 0x3D | 部分覆盖 | 路由存在，业务未实现。 |
| 6 | 现有入口 | 对象交易字段 | `Trade/TradeMoney/TradeOK` | `object.Object` 注释字段 | 未覆盖 | 恢复或重设交易状态字段。 |
| 7 | 配置 | 交易总开关 | `Trade` | `conf.CommonServer.GameServerInfo.TradeEnable`、`conf.Common.General.EnableTrade` | 部分覆盖 | 配置已加载，需要统一语义并接入请求入口。 |
| 8 | 配置 | Harmony 物品交易 | `CanTradeHarmonyItem` | `EnableTradeHarmonyItem` | 部分覆盖 | 物品系统判断是否可交易。 |
| 9 | 配置 | 全卓越物品交易 | `CanTradeFullExcItem` | `EnableTradeFullExcItem` | 部分覆盖 | 限制全卓越物品交易。 |
| 10 | 配置 | 0 序列物品交易 | `CanTradeFFFFFFFFSerialItem` | `EnableTrade0SerialItem` | 部分覆盖 | 限制异常序列物品交易。 |
| 11 | 交易管理 | 交易状态模型 | `m_IfState` + `Trade*` | `InterfaceState` 未完整接入 | 未覆盖 | 设计双人交易会话、目标编号、锁定状态。 |
| 12 | 交易管理 | 交易会话创建 | `CGTradeRequestSend` 成功分支 | 暂无 | 未覆盖 | 双方进入交易状态并记录目标。 |
| 13 | 交易管理 | 交易会话关闭 | `CGTradeResult/Cancel` | 暂无 | 未覆盖 | 成交、取消、断线、失败都必须清理状态。 |
| 14 | 交易管理 | 双方互斥锁 | GameServer 事务状态 | 暂无 | 未覆盖 | Go 侧需要防止双方并发修改背包和交易栏。 |
| 15 | 交易请求 | 请求目标解析 | `lpMsg->NumberH/NumberL` | 暂无 | 未覆盖 | 从客户端编号解析目标玩家。 |
| 16 | 交易请求 | 目标在线检查 | `gObjIsConnected` | `ObjectManager.GetObject` | 未覆盖 | 目标不存在或非玩家时失败。 |
| 17 | 交易请求 | 自己交易自己限制 | `CGTradeRequestSend` 检查 | 暂无 | 未覆盖 | 禁止自己和自己交易。 |
| 18 | 交易请求 | 距离检查 | `gObjCalDistance` | `Object.CalcDistance` | 未覆盖 | 玩家直接交易要求同地图近距离。 |
| 19 | 交易请求 | 等级/状态限制 | GameServer trade checks | 暂无 | 未覆盖 | 低等级、死亡、传送、地图移动等状态禁止交易。 |
| 20 | 交易请求 | MapServerMove 限制 | `m_bMapSvrMoveQuit` 等 | 暂无 | 未覆盖 | 地图服移动过程中禁止交易，跨服移动状态来源归 `28-external-comm.md`。 |
| 21 | 交易请求 | 个人商店冲突 | `m_bPShopOpen` 检查 | `20-personal-shops.md` 待实现 | 未覆盖 | 开个人商店或浏览个人商店时不能直接交易。 |
| 22 | 交易请求 | 其他接口冲突 | `m_IfState` | `InterfaceState` | 未覆盖 | 仓库、商店、合成、NPC 对话等状态禁止交易。 |
| 23 | 交易请求 | 目标忙碌检查 | target transaction/interface | 暂无 | 未覆盖 | 目标已交易或处于互斥界面时失败。 |
| 24 | 交易请求 | 自动接受预留 | GameServer auto response path | 暂无 | 未覆盖 | 如后续有 GM/脚本自动交易，统一入口处理。 |
| 25 | 交易响应 | 响应包下发 | `GCTradeResponseSend` | 暂无 | 未覆盖 | 向双方发送请求、接受、拒绝、失败结果。 |
| 26 | 交易响应 | 拒绝交易 | `CGTradeResponseRecv` Response false | 暂无 | 未覆盖 | 清理请求状态并通知请求者。 |
| 27 | 交易响应 | 接受交易 | `CGTradeResponseRecv` success | 暂无 | 未覆盖 | 双方进入交易界面。 |
| 28 | 交易响应 | 响应目标校验 | `TargetNumber` 检查 | 暂无 | 未覆盖 | 防止非邀请对象伪造接受。 |
| 29 | 交易响应 | 接受时二次校验 | `CGTradeResponseRecv` 重复检查 | 暂无 | 未覆盖 | 接受瞬间重新检查距离、状态、个人商店等。 |
| 30 | 交易栏物品 | 交易栏容量 | `TRADE_BOX_SIZE` | `object.TradeBoxSize` | 部分覆盖 | 常量存在，业务未实现。 |
| 31 | 交易栏物品 | 交易栏占位图 | `TradeMap` | 注释字段 | 未覆盖 | 维护 8x4 交易栏占位。 |
| 32 | 交易栏物品 | 放入交易栏 | item move to trade | `Object.MoveItem` 相关待接入 | 未覆盖 | 从背包移动到交易栏，锁定物品。 |
| 33 | 交易栏物品 | 移出交易栏 | trade item remove | 暂无 | 未覆盖 | 从交易栏移回背包，释放占位。 |
| 34 | 交易栏物品 | 交易物品下发给对方 | `GCTradeOtherAdd` | 暂无 | 未覆盖 | 一方放入物品时向对方显示。 |
| 35 | 交易栏物品 | 删除对方显示物品 | `GCTradeOtherDel` | 暂无 | 未覆盖 | 一方移除物品时通知对方。 |
| 36 | 交易栏物品 | 物品可交易检查 | `IsEnableToTrade` 等 | `item.ItemBase` 有交易字段 | 未覆盖 | 检查绑定、任务、Harmony、期限、宠物、锁定等规则。 |
| 37 | 交易栏物品 | 背包占位检查 | inventory flags | `item.Inventory` | 部分覆盖 | 交易完成前后都要确认背包空间。 |
| 38 | 交易栏物品 | 修改后取消确认 | `GCTradeOkButtonSend(...,2/0)` | 暂无 | 未覆盖 | 任意物品变更后双方确认状态失效。 |
| 39 | 交易金币 | 设置交易金币 | `CGTradeMoneyRecv` | `handle` 0x3A | 未覆盖 | 校验金币数量并记录到交易会话。 |
| 40 | 交易金币 | 金币余额检查 | `TradeMoney <= Money` | `Object.Money` | 未覆盖 | 不能设置超过自身持有金币。 |
| 41 | 交易金币 | 金币上限检查 | Max money checks | 暂无 | 未覆盖 | 成交后双方金币不能超过上限。 |
| 42 | 交易金币 | 金币变更通知对方 | `GCTradeMoneyOther` | 暂无 | 未覆盖 | 一方设置金币后通知另一方。 |
| 43 | 交易金币 | 修改后取消确认 | `GCTradeOkButtonSend` | 暂无 | 未覆盖 | 金币变化后重置确认状态。 |
| 44 | 确认与成交 | 确认按钮 | `CGTradeOkButtonRecv` | `handle` 0x3C | 未覆盖 | 设置当前玩家确认状态。 |
| 45 | 确认与成交 | 确认状态同步 | `GCTradeOkButtonSend` | 暂无 | 未覆盖 | 将一方确认状态告知另一方。 |
| 46 | 确认与成交 | 双方确认检测 | `TradeOK` both true | 暂无 | 未覆盖 | 双方确认后进入成交事务。 |
| 47 | 确认与成交 | 成交前最终校验 | `CGTradeOkButtonRecv` final checks | 暂无 | 未覆盖 | 重新检查在线、状态、背包空间、金币、物品合法性。 |
| 48 | 确认与成交 | 物品交换 | trade commit item swap | 暂无 | 未覆盖 | 原子交换双方交易栏物品。 |
| 49 | 确认与成交 | 金币交换 | trade commit money | 暂无 | 未覆盖 | 原子扣加双方交易金币。 |
| 50 | 确认与成交 | 成交结果下发 | `CGTradeResult(...,1)` | 暂无 | 未覆盖 | 成功后通知双方并刷新背包/金币。 |
| 51 | 取消与失败 | 主动取消 | `CGTradeCancelButtonRecv` | `handle` 0x3D | 未覆盖 | 玩家点取消时回滚交易栏物品和金币状态。 |
| 52 | 取消与失败 | 对方取消通知 | `CGTradeResult(...,0)` | 暂无 | 未覆盖 | 一方取消，双方都收到失败/取消结果。 |
| 53 | 取消与失败 | 断线取消 | close flow calls cancel | `ObjectManager.DeletePlayer` 未接入 | 未覆盖 | 下线时取消交易并归还物品。 |
| 54 | 取消与失败 | 移动/死亡取消 | GameServer state checks | 暂无 | 未覆盖 | 进入不允许状态时取消交易。 |
| 55 | 取消与失败 | 失败回滚 | trade rollback | 暂无 | 未覆盖 | 成交中任意失败必须保证双方背包和金币一致。 |
| 56 | 取消与失败 | 状态清理 | `CGTradeResult` cleanup | 暂无 | 未覆盖 | 清理 TargetNumber、InterfaceState、TradeOK、TradeMoney、交易栏。 |
| 57 | 事务与日志 | 交易事务日志 | `LogAdd` trade logs | 暂无 | 未覆盖 | 记录双方账号、角色、物品、金币、结果。 |
| 58 | 事务与日志 | 异常交易审计 | Anti-hack trade logs | 暂无 | 未覆盖 | 非法物品、伪造目标、越权移动等上报 `27-security.md` 记录安全日志。 |
| 59 | 事务与日志 | 数据保存 | character save after trade | `Player.saveCharacter` 可复用 | 未覆盖 | 成交后保存双方角色背包和金币。 |
| 60 | 跨系统联动 | 交易与合成互斥 | ChaosBox check | `12-mix.md` 已记录 | 未覆盖 | 合成系统调用接口状态判断。 |
| 61 | 跨系统联动 | 交易与商店互斥 | Shop/PShop checks | `08-shops.md`、`20-personal-shops.md` 已记录 | 未覆盖 | NPC 商店、个人商店和交易互斥。 |
| 62 | 跨系统联动 | 交易与地图移动互斥 | Move checks | `06-maps.md` 已记录 | 未覆盖 | 地图移动前取消或禁止交易。 |
| 63 | 测试与校验 | 请求/响应测试 | Trade request/response | 暂无 | 未覆盖 | 覆盖接受、拒绝、距离不足、目标忙碌。 |
| 64 | 测试与校验 | 物品/金币测试 | Trade item/money | 暂无 | 未覆盖 | 覆盖物品放入移除、金币变更、确认重置。 |
| 65 | 测试与校验 | 成交/回滚测试 | Trade commit/rollback | 暂无 | 未覆盖 | 覆盖双方确认、背包满、金币上限、断线取消、并发确认。 |
