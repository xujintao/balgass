# 28. 外部通信系统

本模块覆盖 GameServer 与外部服务之间的通信基础设施：ConnectServer 注册、JoinServer 登录态、DataServer 持久化通道、ExDB 扩展通道、MapServer 跨服通信、协议编解码、请求响应、超时重连、服务状态和业务适配接口。外部通信系统不拥有账号、角色、地图、战盟、好友、邮件、活动、个人商店等业务语义；业务模块只通过本模块发送领域请求并接收结果。运营管理系统需要跨服公告、踢线、封禁、中心服务同步时，通过本模块投递。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 模块边界与总入口 | ExternalCommManager 总入口 | `wsJServerCli`、`wsDataCli`、`wsExDbCli` 分散全局对象 | 暂无 `game/external` | 未覆盖 | 建立统一外部通信系统，管理 ConnectServer、JoinServer、DataServer、ExDB、MapServer 通道。 |
| 2 | 模块边界与总入口 | 外部通信初始化 | `GameMainInit`、`GMJoinServerConnect`、`GMDataServerConnect` | `game.Start` 内联 ConnectServer UDP goroutine | 部分覆盖 | 从运行时剥离外部通信初始化，形成独立组件。 |
| 3 | 模块边界与总入口 | 外部通信关闭 | `wsJServerCli.Close`、`wsDataCli.Close`、`wsExDbCli.Close` | `game.Close` 未统一处理 | 未覆盖 | 关闭时断开外部连接、停止重连、刷出待发送请求。 |
| 4 | 模块边界与总入口 | 服务通道枚举 | Join/Data/ExDB/Connect | 配置字段存在 | 部分覆盖 | 定义 connect、join、data、exdb、mapserver 类型和状态。 |
| 5 | 模块边界与总入口 | 业务请求接口 | `GJ*`、`GD*`、`DG*` 函数族 | 业务直接 DB 或暂无 | 未覆盖 | 业务模块调用领域接口，不直接操作 socket 和协议帧。 |
| 6 | 模块边界与总入口 | 外部事件回调 | `DataServerProtocolCore`、Join protocol core | 暂无 | 未覆盖 | 外部响应进入游戏主循环，避免并发修改游戏状态。 |
| 7 | 模块边界与总入口 | 单服降级模式 | GameServer 依赖外部服 | server-game 当前本地 DB | 设计替代 | 允许本地 DB 模式绕过 DataServer，但接口边界保持一致。 |
| 8 | 模块边界与总入口 | 多服模式开关 | `ServerType`、MapServerInfo | `ServerType`、`MapServers` 已解析 | 部分覆盖 | 配置决定是否启用 Join/Data/ExDB/MapServer 通道。 |
| 9 | 模块边界与总入口 | 外部请求上下文 | `aIndex`、account、name、server code | 分散字段 | 未覆盖 | 每次请求携带玩家索引、账号、角色、server code、trace id。 |
| 10 | 模块边界与总入口 | 统一错误模型 | connect fail、auth fail、timeout | 暂无 | 未覆盖 | 将连接失败、协议失败、超时、拒绝、业务失败区分。 |
| 11 | 模块边界与总入口 | 外部通信测试夹具 | GameServer 靠运行验证 | 暂无 | 未覆盖 | 提供 fake Join/Data/ExDB 服务用于单元和集成测试。 |
| 12 | 外部服务配置 | ServerCode 读取 | `GetGameServerCode` | `conf.Server.GameServerInfo.Code` | 已覆盖 | 统一 server code 来源，供注册、在线态、跨服移动使用。 |
| 13 | 外部服务配置 | ServerType 读取 | `GetServerType` | `GameServerInfo.Type` | 部分覆盖 | BattleCore、普通服、活动服等类型影响连接策略。 |
| 14 | 外部服务配置 | ConnectServer 地址 | `GetConnectServerIP/Port` | `ConnectServerIP/Port` | 已覆盖 | 用于服务器列表/在线率注册。 |
| 15 | 外部服务配置 | JoinServer 地址 | `GetJoinServerIP/Port` | `JoinServerIP/Port` | 已覆盖 | 用于登录态、在线唯一性、跨服认证。 |
| 16 | 外部服务配置 | DataServer 地址 | `GetDataServerIP/Port` | `DataServerIP/Port` | 已覆盖 | 用于角色、仓库、战盟、活动等持久化通道。 |
| 17 | 外部服务配置 | ExDB 地址 | `GetExDBIP/Port` | `ExDBIP/Port` | 已覆盖 | 用于好友、邮件等扩展服务通道。 |
| 18 | 外部服务配置 | DB 直连配置 | DataServer 内部使用 DB | server-game `DBHost/DBPort` tag 当前可疑 | 需修正 | 单服本地 DB 模式保留，但不应污染外部通信接口。 |
| 19 | 外部服务配置 | MapServerInfo 读取 | `CMapServerManager::LoadData` | `conf.MapServers` 已解析 | 部分覆盖 | 建立运行态 server list、mapping 和 map ownership 索引。 |
| 20 | 外部服务配置 | 配置合法性校验 | `ReadServerInfo` | 当前缺少统一校验 | 未覆盖 | 启动时校验 IP、端口、server code、map server group。 |
| 21 | 外部服务配置 | 服务启用策略 | GameServer 强依赖部分服务 | 暂无 | 未覆盖 | 定义哪些外部服务为必需，哪些可降级。 |
| 22 | 外部服务配置 | BattleCore 例外 | `ExDataServerConnect` 跳过 ExDB | `ServerType` 未接入 | 未覆盖 | BattleCore 模式下按 GameServer 语义调整连接。 |
| 23 | ConnectServer 注册 | ConnectServer 客户端 | ConnectServer/WzUdp 注册 | `net.Dial("udp")` | 部分覆盖 | 封装 UDP 客户端，不在 `game.Start` 内直接拼帧发送。 |
| 24 | ConnectServer 注册 | 注册消息结构 | server register packet | `model.MsgServerRegister` | 已覆盖 | 保留当前模型，纳入外部通信协议包。 |
| 25 | ConnectServer 注册 | 注册包 Marshal | `DataSend` packet build | `MsgServerRegister.Marshal` | 已覆盖 | 增加 marshal 错误测试和字段范围测试。 |
| 26 | ConnectServer 注册 | C1 0x01 帧封装 | GameServer/ConnectServer 协议 | `[]byte{0xC1,0x08,0x01}` | 部分覆盖 | 将硬编码帧头封装为协议编码器。 |
| 27 | ConnectServer 注册 | 在线率计算输入 | ServerInfoDisplayer 在线率 | `ObjectManager.GetPlayerPercent` | 已覆盖 | 外部通信系统只接收在线率，不直接扫描对象池。 |
| 28 | ConnectServer 注册 | NonPVP 类型上报 | ServerType/NonPK | `NonPVP` 转 `Type_` | 已覆盖 | 保持当前非 PVP 标记语义。 |
| 29 | ConnectServer 注册 | 周期注册 Tick | `GameServerInfoSend` | `Process1000ms` 入队 | 部分覆盖 | Tick 可由游戏主循环触发，发送由外部通信系统执行。 |
| 30 | ConnectServer 注册 | 注册发送失败处理 | `DataSend` 返回值 | 当前只日志 | 未覆盖 | 发送失败后标记服务异常并重试。 |
| 31 | ConnectServer 注册 | ConnectServer 重连 | UDP/连接重建 | 当前无 | 未覆盖 | 地址不可用或写失败时重建连接。 |
| 32 | ConnectServer 注册 | 关闭下线注册 | `GameServerInfoSendStop` | 当前无 | 未覆盖 | 进程关闭时上报不可用或停止注册。 |
| 33 | ConnectServer 注册 | ConnectServer 状态监控 | ServerInfoDisplayer 状态 | 暂无 | 未覆盖 | 暴露注册成功、失败次数、最近发送时间。 |
| 34 | JoinServer 连接与登录态 | JoinServer TCP 客户端 | `wsJServerCli` | 暂无 | 未覆盖 | 建立 JoinServer 连接、读写循环和状态机。 |
| 35 | JoinServer 连接与登录态 | GMJoinServerConnect | `GMJoinServerConnect` | 暂无 | 未覆盖 | 连接成功后设置 protocol core 并发送登录。 |
| 36 | JoinServer 连接与登录态 | JoinServer protocol core | `SProtocolCore` | 暂无 | 未覆盖 | 解码 JoinServer 返回并派发到游戏主循环。 |
| 37 | JoinServer 连接与登录态 | GJServerLogin | `GJServerLogin` | 暂无 | 未覆盖 | GameServer 启动后向 JoinServer 登录。 |
| 38 | JoinServer 连接与登录态 | JoinServer 登录结果 | JoinServer auth result | 暂无 | 未覆盖 | 登录失败应停止或降级，成功后标记可用。 |
| 39 | JoinServer 连接与登录态 | JoinServer 断线处理 | `JoinServerConnected`、`JoinServerDCTime` | 暂无 | 未覆盖 | 断线后禁止新登录或进入降级策略。 |
| 40 | JoinServer 连接与登录态 | JoinServer 重连 | GameMain reconnect flow | 暂无 | 未覆盖 | 定时重连并避免并发重连。 |
| 41 | JoinServer 连接与登录态 | 账号登录请求 | `CSPJoinIdPassRequest` 到 JoinServer | `Player.Login` 直查 DB | 未覆盖 | 多服模式下账号认证走 JoinServer。 |
| 42 | JoinServer 连接与登录态 | 登录结果回调 | join result packet | `MsgLoginReply` 本地生成 | 未覆盖 | JoinServer 结果回到账号系统，账号系统决定回复客户端。 |
| 43 | JoinServer 连接与登录态 | 账号在线唯一性 | JoinServer online table | 当前无 | 未覆盖 | 查询或通知中心在线态，防止跨服重复登录。 |
| 44 | JoinServer 连接与登录态 | GJPUserClose | `GJPUserClose` | 当前只本地 Offline | 未覆盖 | 下线时通知 JoinServer 清理在线状态。 |
| 45 | JoinServer 连接与登录态 | GJSetCharacterInfo 触发 | `GJSetCharacterInfo` 保存角色 | `saveCharacter` 直写 DB | 部分覆盖 | 角色保存请求进入 DataServer，在线态通知可走 JoinServer。 |
| 46 | JoinServer 连接与登录态 | 断线踢线请求 | JoinServer close/kick | 暂无 | 未覆盖 | 中心服务要求踢线时通过对象系统关闭玩家。 |
| 47 | JoinServer 连接与登录态 | 禁重连通知 | disable reconnect | 安全系统待实现 | 未覆盖 | 安全系统触发后通过 JoinServer 同步中心状态。 |
| 48 | JoinServer 连接与登录态 | MapServerAuth 请求 | `GJReqMapSvrAuth` | `0xB101 mapServerAuth` 路由占位 | 未覆盖 | 迁入地图服时向 JoinServer 验证跨服认证信息。 |
| 49 | JoinServer 连接与登录态 | MapServerMove 请求 | `GJReqMapSvrMove` | 当前无 | 未覆盖 | 迁出时请求 JoinServer 分配目标服认证。 |
| 50 | JoinServer 连接与登录态 | 角色所在服更新 | `GJUpdateMatchDBUserCharacters` 等 | 暂无 | 未覆盖 | 多服下维护角色当前 server code。 |
| 51 | DataServer 连接与持久化通道 | DataServer TCP 客户端 | `wsDataCli` | 暂无 | 未覆盖 | 建立 DataServer 连接、读写循环和发送队列。 |
| 52 | DataServer 连接与持久化通道 | GMDataServerConnect | `GMDataServerConnect` | 暂无 | 未覆盖 | 连接成功后设置 `DataServerProtocolCore` 并登录。 |
| 53 | DataServer 连接与持久化通道 | DataServerLogin | `DataServerLogin` | 暂无 | 未覆盖 | 启动后向 DataServer 登录并上报 server code。 |
| 54 | DataServer 连接与持久化通道 | DataServerLoginResult | `DataServerLoginResult` | 暂无 | 未覆盖 | 登录结果失败应停止服务或切换本地 DB 降级。 |
| 55 | DataServer 连接与持久化通道 | DataServerProtocolCore | `DataServerProtocolCore` | 暂无 | 未覆盖 | 按 protoNum 分发 DataServer 响应。 |
| 56 | DataServer 连接与持久化通道 | 角色列表请求 | `DataServerGetCharListRequest` | `DB.GetCharacterList` | 设计替代 | 多服模式通过 DataServer，请求结果回角色系统。 |
| 57 | DataServer 连接与持久化通道 | 角色列表响应 | `JGPGetCharList` | `GetCharacterList` 本地组包 | 设计替代 | 响应需校验账号一致性后转给角色系统。 |
| 58 | DataServer 连接与持久化通道 | 创建角色请求 | `SDHP_CREATECHAR` | `DB.CreateCharacter` | 设计替代 | DataServer 模式由角色系统发请求，外部通信负责投递。 |
| 59 | DataServer 连接与持久化通道 | 创建角色响应 | `JGCharacterCreateRequest` | 本地返回 | 设计替代 | 响应映射为角色系统创建结果。 |
| 60 | DataServer 连接与持久化通道 | 删除角色请求 | `SDHP_CHARDELETE` | `DB.DeleteCharacterByName` | 设计替代 | 删除请求投递和超时处理归本模块。 |
| 61 | DataServer 连接与持久化通道 | 删除角色响应 | `JGCharDelRequest` | 本地返回 | 设计替代 | 响应需校验 account/character 一致性。 |
| 62 | DataServer 连接与持久化通道 | 加载角色请求 | DB char info request | `DB.GetCharacterByName` | 设计替代 | DataServer 模式下加载角色通过请求响应完成。 |
| 63 | DataServer 连接与持久化通道 | 加载角色响应 | `JGGetCharacterInfo` | `LoadCharacter` 本地赋值 | 设计替代 | 响应交给角色系统灌入 Player。 |
| 64 | DataServer 连接与持久化通道 | 保存角色请求 | `GJSetCharacterInfo`、`SDHP_DBCHAR_INFOSAVE` | `Player.saveCharacter` | 部分覆盖 | DataServer 模式下角色保存走外部通道。 |
| 65 | DataServer 连接与持久化通道 | 保存失败处理 | `wsDataCli.DataSend` false | 当前本地 DB error 日志 | 未覆盖 | 保存失败要返回业务可处理错误并记录重试策略。 |
| 66 | DataServer 连接与持久化通道 | 仓库读取请求 | `GDGetWarehouseList` | `DB.GetAccountByID` | 设计替代 | 仓库业务归道具/账号，DataServer 投递归本模块。 |
| 67 | DataServer 连接与持久化通道 | 仓库保存请求 | `GDSetWarehouseList`、`GDSetWarehouseMoney` | `UpdateAccountWarehouse` | 设计替代 | 关闭仓库时可通过 DataServer 保存。 |
| 68 | DataServer 连接与持久化通道 | OptionData 请求 | `DGOptionDataSend/Recv` | `UpdateCharacterMuKey` | 部分覆盖 | MuKey/MuBot 可通过 DataServer 通道抽象。 |
| 69 | DataServer 连接与持久化通道 | 个人商店价格请求 | `GDRequestPShopItemValue` | 暂无 | 未覆盖 | 个人商店模块通过本模块请求价格持久化。 |
| 70 | DataServer 连接与持久化通道 | 个人商店价格移动 | `GDMovePShopItem` | 暂无 | 未覆盖 | 物品移动导致价格位置变化时投递 DataServer。 |
| 71 | DataServer 连接与持久化通道 | 战盟创建持久化 | DataServer guild create | 暂无 | 未覆盖 | 战盟系统发领域请求，外部通信投递。 |
| 72 | DataServer 连接与持久化通道 | 战盟成员持久化 | guild member DB calls | 暂无 | 未覆盖 | 战盟加入、退出、职位变更通过 DataServer。 |
| 73 | DataServer 连接与持久化通道 | Gens 持久化 | Gens rank/point DB calls | 暂无 | 未覆盖 | 家族贡献和排名保存通道。 |
| 74 | DataServer 连接与持久化通道 | 活动积分保存 | ITL/Kanturu/Event save calls | 暂无 | 未覆盖 | 副本和活动模块保存积分、胜负和奖励记录。 |
| 75 | DataServer 连接与持久化通道 | 匹配系统请求 | party/guild matching GDReq | 暂无 | 未覆盖 | 组队/战盟匹配的数据服请求通道。 |
| 76 | DataServer 连接与持久化通道 | 封禁/处罚同步 | ban/check packets | 安全系统待实现 | 未覆盖 | `27-security.md` 和 `29-ops.md` 通过 DataServer 同步封禁状态。 |
| 77 | ExDB 连接与扩展通道 | ExDB TCP 客户端 | `wsExDbCli` | 暂无 | 未覆盖 | 建立 ExDB 连接、读写循环和状态机。 |
| 78 | ExDB 连接与扩展通道 | ExDataServerConnect | `ExDataServerConnect` | 暂无 | 未覆盖 | 连接 ExDB 并设置 `ExDataServerProtocolCore`。 |
| 79 | ExDB 连接与扩展通道 | ExDataServerLogin | `ExDataServerLogin` | 暂无 | 未覆盖 | 启动后向 ExDB 登录。 |
| 80 | ExDB 连接与扩展通道 | BattleCore 跳过 ExDB | `SERVER_BATTLECORE` 分支 | 暂无 | 未覆盖 | 按 server type 决定是否连接 ExDB。 |
| 81 | ExDB 连接与扩展通道 | ExDB 协议分发 | `ExDataClientMsgProc` | 暂无 | 未覆盖 | 接收 ExDB 响应并转发给对应业务。 |
| 82 | ExDB 连接与扩展通道 | 好友状态通道 | Friend server data | 好友系统待实现 | 未覆盖 | 好友在线、换线、上线下线通知通道。 |
| 83 | ExDB 连接与扩展通道 | 邮件 Memo 通道 | Memo/ExDB data | 邮件系统待实现 | 未覆盖 | Memo 发送、读取、删除、列表请求通道。 |
| 84 | ExDB 连接与扩展通道 | 聊天室邀请通道 | friend chat room | 暂无 | 未覆盖 | 好友聊天室或跨服邀请的外部服务适配。 |
| 85 | ExDB 连接与扩展通道 | 扩展服务降级 | ExDB disconnected state | 暂无 | 未覆盖 | ExDB 不可用时好友/邮件返回明确错误。 |
| 86 | MapServer 跨服通信 | MapServerManager 运行表 | `CMapServerManager` | `conf.MapServers` 仅解析 | 未覆盖 | 构建 server code、group、map ownership、dest server 索引。 |
| 87 | MapServer 跨服通信 | ServerList 映射 | `_MAPSVR_DATA` server list | `MapServers.ServerList` | 部分覆盖 | 解析 code、group、IP、port、name。 |
| 88 | MapServer 跨服通信 | ServerMapping 映射 | map to server mapping | `MapServers.ServerMapping` | 部分覆盖 | 解析 MapNumber、MoveAble、DestServerCode。 |
| 89 | MapServer 跨服通信 | CheckMapCanMove | `CheckMapCanMove` | `06-maps.md` 待实现 | 未覆盖 | 查询目标地图是否归当前服或允许跨服。 |
| 90 | MapServer 跨服通信 | CheckMoveMapSvr | `CheckMoveMapSvr` | `06-maps.md` 待实现 | 未覆盖 | 根据目标地图和当前 server code 解析目标服。 |
| 91 | MapServer 跨服通信 | GetSvrCodeData | `GetSvrCodeData` | 暂无 | 未覆盖 | 查询目标服 IP、端口和状态。 |
| 92 | MapServer 跨服通信 | MapServerAuthInfo | `m_MapServerAuthInfo` | 暂无 | 未覆盖 | 保存跨服迁入认证临时字段。 |
| 93 | MapServer 跨服通信 | 迁出状态设置 | `m_bMapSvrMoveReq` | 暂无 | 未覆盖 | 迁出请求中禁止交易、移动、背包等敏感操作。 |
| 94 | MapServer 跨服通信 | 迁出完成状态 | `m_bMapSvrMoveQuit` | 暂无 | 未覆盖 | 等待目标服接管期间保持临时用户状态。 |
| 95 | MapServer 跨服通信 | 跨服移动请求发送 | `GJReqMapSvrMove` | 暂无 | 未覆盖 | 向 JoinServer 请求目标服移动认证。 |
| 96 | MapServer 跨服通信 | 跨服认证请求发送 | `GJReqMapSvrAuth` | `0xB101` 入口占位 | 未覆盖 | 目标服收到客户端认证后向 JoinServer 验证。 |
| 97 | MapServer 跨服通信 | 跨服保存前置 | `GJSetCharacterInfo(..., TRUE)` | `saveCharacter` 无 map move 参数 | 未覆盖 | 迁出前保存角色并标记 MapServerMove。 |
| 98 | MapServer 跨服通信 | 跨服失败回滚 | GameServer logs/fail branch | 暂无 | 未覆盖 | 目标服不可用、认证失败、保存失败时回滚状态。 |
| 99 | MapServer 跨服通信 | 地图服组在线数 | `GDReqMapSrvGroupServerCount` | 暂无 | 未覆盖 | 向 DataServer 请求同组在线 server 数。 |
| 100 | MapServer 跨服通信 | 地图服组播 | `SendMapServerGroupMsg` | 暂无 | 未覆盖 | 世界事件需要跨地图服组播状态。 |
| 101 | 外部协议编解码 | 外部协议帧头 | C1/C2/C3/C4 server packets | 仅部分 c1c2 客户端协议 | 未覆盖 | 封装外部服务协议帧，与客户端协议隔离。 |
| 102 | 外部协议编解码 | JoinServer packet 编码 | `GJ*` structs | 暂无 | 未覆盖 | 定义 JoinServer 请求结构和 marshal。 |
| 103 | 外部协议编解码 | JoinServer packet 解码 | Join protocol core structs | 暂无 | 未覆盖 | 解码登录、踢线、跨服认证等响应。 |
| 104 | 外部协议编解码 | DataServer packet 编码 | `SDHP_*` structs | 暂无 | 未覆盖 | 定义角色、仓库、战盟、活动等 DataServer 请求结构。 |
| 105 | 外部协议编解码 | DataServer packet 解码 | `DataServerProtocolCore` | 暂无 | 未覆盖 | 按 protoNum 解码 DataServer 响应。 |
| 106 | 外部协议编解码 | ExDB packet 编码 | ExDB structs | 暂无 | 未覆盖 | 定义好友、邮件等扩展请求结构。 |
| 107 | 外部协议编解码 | ExDB packet 解码 | ExDB protocol core | 暂无 | 未覆盖 | 解码好友/邮件响应。 |
| 108 | 外部协议编解码 | 端序处理 | `memcpy` binary structs | Go `binary.LittleEndian` 局部使用 | 部分覆盖 | 统一外部协议端序和定长字符串处理。 |
| 109 | 外部协议编解码 | 定长账号字段 | `MAX_ACCOUNT_LEN` | 多处 trim/xor | 部分覆盖 | 统一账号、角色名、战盟名字段编码。 |
| 110 | 外部协议编解码 | GBK/编码边界 | GameServer char arrays | 客户端模型已有部分 | 部分覆盖 | 外部服务和客户端编码转换应集中处理。 |
| 111 | 请求响应与超时 | 请求 ID 生成 | aIndex/sequence/proto | 暂无 | 未覆盖 | 为异步请求生成可追踪 request id。 |
| 112 | 请求响应与超时 | 玩家索引关联 | `aIndex` in packet | 当前业务本地直接调用 | 未覆盖 | 响应回来时找到原玩家或处理玩家已下线。 |
| 113 | 请求响应与超时 | 超时表 | GameServer 分散等待 | 暂无 | 未覆盖 | 请求发出后登记超时，避免业务永久等待。 |
| 114 | 请求响应与超时 | 超时回调 | 分散错误处理 | 暂无 | 未覆盖 | 超时后通知账号/角色/业务模块返回失败。 |
| 115 | 请求响应与超时 | 请求取消 | disconnect path | 暂无 | 未覆盖 | 玩家断线时取消未完成的外部请求。 |
| 116 | 请求响应与超时 | 重复响应处理 | GameServer index check | 暂无 | 未覆盖 | 响应已处理或玩家已离线时只记录日志。 |
| 117 | 请求响应与超时 | 响应账号一致性校验 | `gObjIsAccontConnect` | 角色 DB 按 accountID 查询 | 部分覆盖 | 外部响应必须校验 account/character 与当前连接一致。 |
| 118 | 请求响应与超时 | 响应角色一致性校验 | `JGGetCharacterInfo` duplicate checks | 当前缺少在线唯一性 | 未覆盖 | 防止响应灌入错误对象或角色复制。 |
| 119 | 请求响应与超时 | 发送队列 | `FDWRITE_MsgDataSend` | 暂无 | 未覆盖 | 外部连接写入统一走队列，避免阻塞游戏主循环。 |
| 120 | 请求响应与超时 | 背压策略 | socket send buffer | 暂无 | 未覆盖 | 队列满时拒绝新请求或触发服务不可用。 |
| 121 | 请求响应与超时 | 幂等请求策略 | save/move repeated calls | 暂无 | 未覆盖 | 角色保存、跨服移动等请求需要防重复发送。 |
| 122 | 请求响应与超时 | 请求审计日志 | `g_Log.Add` scattered | 暂无 | 未覆盖 | 记录请求类型、耗时、结果、玩家上下文。 |
| 123 | 重连与服务状态 | 服务状态枚举 | Connected/Disconnected globals | 暂无 | 未覆盖 | 定义 disconnected、connecting、authenticating、ready、degraded。 |
| 124 | 重连与服务状态 | JoinServerConnected | `JoinServerConnected` | 暂无 | 未覆盖 | 维护 JoinServer 当前可用性。 |
| 125 | 重连与服务状态 | DataServerConnected | `DataServerConnected` | 暂无 | 未覆盖 | 维护 DataServer 当前可用性。 |
| 126 | 重连与服务状态 | ExDBConnected | ExDB state | 暂无 | 未覆盖 | 维护 ExDB 当前可用性。 |
| 127 | 重连与服务状态 | 断线时间记录 | `JoinServerDCTime` | 暂无 | 未覆盖 | 记录断线时间和持续时长。 |
| 128 | 重连与服务状态 | 重连间隔 | GameMain reconnect timer | 暂无 | 未覆盖 | 配置化重连间隔，避免快速重试。 |
| 129 | 重连与服务状态 | 重连抖动 | 无统一 | 暂无 | 未覆盖 | 多服同时启动时避免重连风暴。 |
| 130 | 重连与服务状态 | 认证后恢复 | login after reconnect | 暂无 | 未覆盖 | 重连成功后重新登录外部服务并恢复状态。 |
| 131 | 重连与服务状态 | 服务不可用策略 | auth closed stops server | 暂无 | 未覆盖 | 必需服务不可用时拒绝登录或关闭进程。 |
| 132 | 重连与服务状态 | DataServer 降级策略 | GameServer 多数强依赖 | 本地 DB 可降级 | 设计替代 | 明确本地 DB 模式和 DataServer 模式不能混用混乱。 |
| 133 | 重连与服务状态 | ExDB 降级策略 | Friend server offline | 暂无 | 未覆盖 | 好友/邮件不可用时业务返回明确错误。 |
| 134 | 重连与服务状态 | ConnectServer 降级策略 | server list unavailable | 当前发送失败只日志 | 未覆盖 | 注册失败不应影响游戏主循环，但要暴露健康状态。 |
| 135 | 重连与服务状态 | 健康检查接口 | ServerInfoDisplayer | 暂无 | 未覆盖 | HTTP/WS 管理入口可查看外部服务状态。 |
| 136 | 业务模块适配接口 | AccountService 外部适配 | JoinServer auth | `02-accounts.md` | 未覆盖 | 提供登录、下线、在线唯一性、中心踢线接口。 |
| 137 | 业务模块适配接口 | CharacterService 外部适配 | DataServer char list/info/save | `03-characters.md` | 未覆盖 | 提供角色列表、创建、删除、加载、保存接口。 |
| 138 | 业务模块适配接口 | WarehouseService 外部适配 | GD warehouse calls | 账号/道具模块待拆 | 未覆盖 | 提供仓库读取、保存、金币保存接口。 |
| 139 | 业务模块适配接口 | GuildService 外部适配 | guild DataServer calls | `16-guild.md` | 未覆盖 | 提供战盟创建、成员变更、公告、联盟等持久化接口。 |
| 140 | 业务模块适配接口 | FriendMailService 外部适配 | ExDB friend/memo | `18-friends-mail.md` | 未覆盖 | 提供好友在线、邮件发送读取删除接口。 |
| 141 | 业务模块适配接口 | MapMoveService 外部适配 | `GJReqMapSvrMove/Auth` | `06-maps.md` | 未覆盖 | 提供跨服目标解析后的请求、认证和回调接口。 |
| 142 | 业务模块适配接口 | PersonalShopService 外部适配 | PShop item value GD calls | `20-personal-shops.md` | 未覆盖 | 提供个人商店价格读取、移动、删除、保存接口。 |
| 143 | 业务模块适配接口 | EventService 外部适配 | ITL/Kanturu/Event save | `21/22/23` | 未覆盖 | 提供副本、普通活动、世界事件积分和状态保存接口。 |
| 144 | 业务模块适配接口 | SecurityService 外部适配 | ban/kick/disable reconnect | `27-security.md` | 未覆盖 | 提供跨服踢线、封禁、禁重连中心通知接口。 |
| 145 | 业务模块适配接口 | MatchingService 外部适配 | party/guild matching GDReq | `15-party.md`、`16-guild.md` | 未覆盖 | 提供组队匹配、战盟匹配请求通道。 |
| 146 | 业务模块适配接口 | LocalRepository 适配器 | 不适用 | `model.DB` 当前直连 | 部分覆盖 | 单服模式用本地 repository 实现同一接口。 |
| 147 | 业务模块适配接口 | RemoteDataServer 适配器 | DataServer requests | 暂无 | 未覆盖 | 多服模式用外部请求实现同一接口。 |
| 148 | 日志、监控与测试 | 外部通信日志 | `g_Log.Add` connect/send logs | `slog` 分散 | 部分覆盖 | 统一记录 service、proto、request id、account、result、duration。 |
| 149 | 日志、监控与测试 | 连接状态指标 | ServerInfoDisplayer | 暂无 | 未覆盖 | 暴露连接状态、重连次数、发送失败、超时数。 |
| 150 | 日志、监控与测试 | 请求耗时指标 | 无统一 | 暂无 | 未覆盖 | 统计各请求类型 p50/p95/p99 耗时。 |
| 151 | 日志、监控与测试 | 队列长度指标 | socket send queue | 暂无 | 未覆盖 | 监控发送队列和 pending 请求数量。 |
| 152 | 日志、监控与测试 | 协议解码错误日志 | protocol core error | 暂无 | 未覆盖 | 非法 proto、长度错误、字段越界都要记录。 |
| 153 | 日志、监控与测试 | fake JoinServer 测试 | 无统一 | 暂无 | 未覆盖 | 覆盖登录成功、重复登录、踢线、断线重连。 |
| 154 | 日志、监控与测试 | fake DataServer 测试 | 无统一 | 暂无 | 未覆盖 | 覆盖角色列表、创建、删除、加载、保存和超时。 |
| 155 | 日志、监控与测试 | fake ExDB 测试 | 无统一 | 暂无 | 未覆盖 | 覆盖好友在线、邮件列表、发送、读取、删除。 |
| 156 | 日志、监控与测试 | MapServer 测试 | GameServer 运行验证 | 暂无 | 未覆盖 | 覆盖目标服解析、迁出请求、迁入认证、失败回滚。 |
| 157 | 日志、监控与测试 | ConnectServer 注册测试 | 当前无 | 暂无 | 未覆盖 | 覆盖注册包编码、在线率、发送失败和关闭注册。 |
| 158 | 日志、监控与测试 | 降级模式测试 | 本地 DB 模式 | 当前无显式测试 | 未覆盖 | 覆盖外部服务关闭时本地 DB 模式仍可登录。 |
| 159 | 日志、监控与测试 | 并发请求测试 | IOCP/socket 并发 | Go channel 模式 | 未覆盖 | 覆盖多个玩家同时登录、加载角色、保存角色。 |
| 160 | 日志、监控与测试 | 回归测试总集 | GameServer 靠运行验证 | 暂无 | 未覆盖 | 覆盖连接、登录、协议、请求响应、超时、重连、关闭和跨模块适配。 |
