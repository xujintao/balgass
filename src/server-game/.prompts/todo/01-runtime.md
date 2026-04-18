# 1. 运行时

本模块覆盖 server-game 的进程生命周期、配置日志、网络入口、协议分发、请求串行化、Tick 调度、关闭流程与运行时诊断。GameServer 的 Win32/IOCP/窗口模型只作为运行语义参考；Go 版以 `main`、goroutine、channel、context、`net/http`、C1/C2 TCP server 和 WebSocket 为运行时主体。ConnectServer/JoinServer/DataServer/ExDB/MapServer 通信归 `28-external-comm.md`，运行时只负责启动和关闭外部通信组件；HTTP/WS 管理入口、维护关闭、运行状态展示和后台命令归 `29-ops.md`；封包加密、CheckSum、CRC、包频率和安全处罚策略归 `27-security.md`，运行时只负责在收发包边界调用安全系统。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 进程入口与服务生命周期 | 进程入口函数 | `GameServer.cpp::WinMain` | `main.go::main` | 已覆盖 | 保持 Go 入口只负责装配运行时组件，不把游戏业务初始化散落到入口函数。 |
| 2 | 进程入口与服务生命周期 | 窗口/控制台运行模型替代 | `GameServer.cpp::MyRegisterClass`、`InitInstance`、`WndProc` | `main.go` 无 GUI，使用日志和信号控制 | 设计替代 | Win32 窗口消息不迁移，Go 侧只保留运行状态、信号退出和日志观测。 |
| 3 | 进程入口与服务生命周期 | 全服启动编排 | `GameServer.cpp::AllServerStart` | `main.go` 顺序启动 `game`、TCP、HTTP | 部分覆盖 | 补齐启动失败后的统一回滚策略，避免某个服务失败后残留 goroutine。 |
| 4 | 进程入口与服务生命周期 | Game 运行时启动 | `GameMain.cpp::GameMainInit`、`GameMainServerCreate` | `game.Game.Start()` | 部分覆盖 | 明确 `Start` 是否可重复调用，并增加重复启动保护。 |
| 5 | 进程入口与服务生命周期 | TCP 服务启动 | `GameServer.cpp::GameServerStart`、`wsGameServer.cpp::CreateServer` | `c1c2.Server.ListenAndServe` | 部分覆盖 | 将监听失败、端口占用、关闭时返回值区分为可预期错误和致命错误。 |
| 6 | 进程入口与服务生命周期 | HTTP 服务启动 | GameServer 无对应 HTTP 管理入口 | `http.Server.ListenAndServe` | 已覆盖 | HTTP 是 Go 版调试/管理入口，具体后台 API 和权限归 `29-ops.md`。 |
| 7 | 进程入口与服务生命周期 | 启动错误汇聚 | `GameServer.cpp::GameServerStart` 返回状态 | `main.go::errChan` | 部分覆盖 | `errChan` 需要过滤 `http.ErrServerClosed` 等正常关闭错误。 |
| 8 | 进程入口与服务生命周期 | OS 信号监听 | `WndProc` 处理关闭消息 | `signal.Notify` 监听 `SIGINT/SIGTERM` | 已覆盖 | 保持信号只触发关闭编排，不直接操作对象状态。 |
| 9 | 进程入口与服务生命周期 | 主 goroutine 阻塞策略 | Win32 message loop | `select` 等待 signal 或 server error | 已覆盖 | 后续可增加 runtime health channel，但不要破坏当前单点退出语义。 |
| 10 | 进程入口与服务生命周期 | 启动/退出日志 | `TLog`、`LogToFile` | `slog.Info/Error` | 部分覆盖 | 补充服务名、端口、配置路径、退出原因等结构化字段。 |
| 11 | 配置与日志初始化 | 环境变量加载 | `configread.cpp::ReadConfig` | `conf.ENV` | 已覆盖 | 保持环境变量只处理运行环境，不混入游戏业务配置。 |
| 12 | 配置与日志初始化 | 日志 writer 初始化 | `TLog.cpp`、`LogToFile.cpp` | `conf.ENV` 中 `os.OpenFile` 和 `io.MultiWriter` | 部分覆盖 | 增加文件关闭/轮转策略，避免长期运行日志文件不可控。 |
| 13 | 配置与日志初始化 | 日志级别映射 | `TLog` 日志级别语义 | `conf.ENV` 中 `slog.Level` 映射 | 部分覆盖 | 非法日志级别应显式报错或回退，不应静默使用默认零值。 |
| 14 | 配置与日志初始化 | INI 配置加载 | `configread.cpp::ReadServerInfo` | `conf.INI` | 已覆盖 | 保持加载失败快速退出，但为测试提供可替换入口。 |
| 15 | 配置与日志初始化 | XML 配置加载 | `ReadCommonServerInfo`、各事件/系统 XML 加载 | `conf.XML` | 已覆盖 | 需要在业务模块中补齐具体 XML 字段语义校验。 |
| 16 | 配置与日志初始化 | JSON 配置加载 | GameServer 主要不用 JSON | `conf.JSON` | 已覆盖 | 当前是通用能力，保留但不强行用于 GameServer 对照。 |
| 17 | 配置与日志初始化 | 配置加载顺序 | `GameMainInit` 前后读取 common/event 配置 | `conf.init` 固定顺序加载 | 已覆盖 | 顺序依赖需要在文档里固定，避免后续模块初始化时读到零值。 |
| 18 | 配置与日志初始化 | CommonServer 配置加载 | `GameMain.cpp::ReadCommonServerInfo` | `conf.init` 加载 `CommonServer.cfg` | 部分覆盖 | 后续按业务模块补齐配置项使用点，不在运行时层展开业务含义。 |
| 19 | 配置与日志初始化 | GameServer.ini 字段绑定 | `SERVER_CONFIG::ReadServerInfo` | `configServer.GameServerInfo` | 需修正 | `DBName/DBUser/DBPassword/DBHost/DBPort` tag 写成 `int`，应改为 `ini`。 |
| 20 | 配置与日志初始化 | 配置失败退出策略 | `ReadConfig` 失败中断启动 | `conf.INI/XML/JSON` 中 `os.Exit(1)` | 部分覆盖 | 后续测试需要抽象加载函数，避免单元测试被 `os.Exit` 终止。 |
| 21 | Game 初始化与组件装载 | 全局 Game 实例 | GameServer 全局状态变量 | `var Game game` | 已覆盖 | 保持唯一实例，但要避免测试间全局状态污染。 |
| 22 | Game 初始化与组件装载 | 包 init 初始化入口 | `GameMainInit` | `func init(){ Game.init() }` | 已覆盖 | 注意 `init` 隐式执行，后续复杂依赖建议显式化。 |
| 23 | Game 初始化与组件装载 | channel 初始化 | GameServer 消息队列/事件队列 | `game.init` 创建各 channel | 已覆盖 | channel 容量应有配置或常量说明，便于压测调整。 |
| 24 | Game 初始化与组件装载 | 玩家连接队列 | IOCP accept 后进入对象流程 | `playerConnRequestChan` | 已覆盖 | 当前通过同步 response channel 返回 ID，语义清晰。 |
| 25 | Game 初始化与组件装载 | 玩家断线队列 | `CloseClient`、`GSDisconnect` | `playerCloseConnRequestChan` | 已覆盖 | 需要确保重复断线不会阻塞或二次删除异常。 |
| 26 | Game 初始化与组件装载 | 玩家 action 队列 | `GMClientMsgProc` 分发协议 | `playerActionChan` | 已覆盖 | 高并发下容量 1000 是否足够需要压测。 |
| 27 | Game 初始化与组件装载 | User/Web 管理连接队列 | GameServer 无直接 WS 对照 | `userConnRequestChan`、`userActionChan` | 已覆盖 | Go 版扩展入口合理，管理权限、审计和后台动作归 `29-ops.md`。 |
| 28 | Game 初始化与组件装载 | 命令队列 | GameServer 管理命令/控制台语义 | `commandRequestChan` | 已覆盖 | 命令执行进入主循环，GM/后台命令注册、权限和审计归 `29-ops.md`。 |
| 29 | Game 初始化与组件装载 | 怪物预生成 | `GameMainInit` 后加载怪物/地图对象 | `monster.SpawnMonster()` | 部分覆盖 | 运行时只负责调用入口，怪物数据和 AI 放入对象/地图模块。 |
| 30 | Game 初始化与组件装载 | Bot 注册 | GameServer bot/辅助逻辑语义 | `bot.BotManager.Register(g)` | 部分覆盖 | Bot 生命周期要和 `Game.Close` 对称，避免退出时残留。 |
| 31 | Game.Start 并发结构 | context 创建 | GameServer 全局关闭标志 | `context.WithCancel` | 已覆盖 | context 是 Go 版关闭信号来源，不需要移植 C++ 全局 flag。 |
| 32 | Game.Start 并发结构 | cancel 保存 | 全局 shutdown 状态 | `g.cancel = cancel` | 部分覆盖 | 需要处理 `Close` 在 `Start` 前调用或重复调用时的 nil/reuse 问题。 |
| 33 | Game.Start 并发结构 | 外部通信 goroutine 启动 | `GMJoinServerConnect`、`GMDataServerConnect`、`ExDataServerConnect` | `game.Start` 第一个 goroutine 只发 ConnectServer 注册 | 部分覆盖 | 后续应迁移到 `28-external-comm.md`，运行时只启动组件并接收退出信号。 |
| 34 | Game.Start 并发结构 | 主游戏循环 goroutine | `GameMain` 主循环/TimerProc | `game.Start` 第二个 goroutine | 已覆盖 | 这是 server-game 的核心串行化边界，业务状态修改应尽量进入此循环。 |
| 35 | Game.Start 并发结构 | goroutine 错误上报 | GameServer 启动状态返回 | 当前 goroutine 内部只记录日志或 `os.Exit` | 需修正 | 建议运行时 goroutine 错误进入统一 error channel，而不是局部退出进程。 |
| 36 | Game.Start 并发结构 | goroutine panic 防护 | SEH/日志语义 | 当前无 recover | 未覆盖 | 关键 goroutine 应有 recover 日志和退出上报策略。 |
| 37 | Game.Start 并发结构 | Start 重入保护 | GameServer 防重复启动语义 | 当前无 guard | 未覆盖 | 增加 started 状态或 sync.Once，避免重复创建 ticker/goroutine。 |
| 38 | Game.Start 并发结构 | Start 后可观测状态 | 信息栏/状态栏 | 当前仅日志 | 部分覆盖 | 可增加运行中、连接数、队列长度等诊断接口。 |
| 39 | Game.Start 并发结构 | goroutine 退出条件 | GameServer shutdown flag | `<-ctx.Done()` | 部分覆盖 | ticker 需要 `Stop`，UDP conn 需要 `Close`，避免资源泄漏。 |
| 40 | Game.Start 并发结构 | 启动单元测试边界 | GameServer 人工启动验证 | 当前无运行时测试 | 未覆盖 | 抽象网络依赖后补齐 Start/Close smoke test。 |
| 41 | ConnectServer 注册运行时 | UDP 地址拼接 | JoinServer/ConnectServer 配置 | `fmt.Sprintf("%s:%d", ConnectServerIP, ConnectServerPort)` | 已覆盖 | 该能力后续归 `28-external-comm.md`，运行时只保留配置装配和启动调用。 |
| 42 | ConnectServer 注册运行时 | UDP 连接建立 | `wsJoinServerCli::Connect`、`WzUdp::CreateSocket` | `net.Dial("udp", addr)` | 部分覆盖 | 需要迁移到外部通信系统并增加断线重连和启动失败回传。 |
| 43 | ConnectServer 注册运行时 | 注册请求队列 | JoinServer 注册包发送触发 | `serverRegisterChan` | 已覆盖 | 当前 1s Tick 入队在线百分比，后续由 `28-external-comm.md` 承接发送。 |
| 44 | ConnectServer 注册运行时 | 服务器编号注入 | `GameServerInfoSend` | `serverRegister.Code = conf.Server.GameServerInfo.Code` | 已覆盖 | 需要保证 Code 配置范围合法。 |
| 45 | ConnectServer 注册运行时 | PVP 类型注入 | GameServer server type/non-PK 字段 | `serverRegister.Type_` | 已覆盖 | 命名可后续统一为 `Type`，避免下划线泄漏。 |
| 46 | ConnectServer 注册运行时 | 注册消息序列化 | JoinServer protocol marshal | `serverRegister.Marshal()` | 已覆盖 | 序列化和协议帧归 `28-external-comm.md`，运行时不直接拼包。 |
| 47 | ConnectServer 注册运行时 | C1 注册帧编码 | `DataSend` C1/C2 包头 | `[]byte{0xC1,0x08,0x01}` | 部分覆盖 | 包长度写死，需要和 marshal 数据长度做一致性校验。 |
| 48 | ConnectServer 注册运行时 | UDP 发送 | `wsJoinServerCli::DataSend`、`WzUdp::SendData` | `c.Write(frame)` | 部分覆盖 | 发送失败、重连和状态降级归 `28-external-comm.md`。 |
| 49 | ConnectServer 注册运行时 | 注册响应读取 | `DataRecv`、`GMJoinClientMsgProc` | 当前无 UDP read loop | 未覆盖 | 如果 ConnectServer 有回包语义，需要补齐读取和协议处理。 |
| 50 | ConnectServer 注册运行时 | 注册关闭释放 | `WzUdp::Close` | 当前 context return 但不 close conn | 未覆盖 | goroutine 退出时应关闭 UDP conn。 |
| 51 | 主循环 channel 串行化 | 玩家连接入队 API | IOCP accept 后对象创建 | `Game.PlayerConn` | 已覆盖 | 同步等待 response channel 保证连接建立结果可返回到网络层。 |
| 52 | 主循环 channel 串行化 | 玩家连接处理 | `gObjAddSearch`/玩家对象创建语义 | `player.NewPlayer` | 已覆盖 | 玩家创建细节归账号/对象模块，运行时只保证串行调用。 |
| 53 | 主循环 channel 串行化 | 玩家断线入队 API | `CloseClient` | `Game.PlayerCloseConn` | 已覆盖 | 断线入队是正确边界，网络层不直接改对象表。 |
| 54 | 主循环 channel 串行化 | 玩家断线处理 | `GSDisconnect`、`gObjDel` | `ObjectManager.DeletePlayer` | 已覆盖 | 需要校验二次关闭和不存在 ID 的处理。 |
| 55 | 主循环 channel 串行化 | 玩家 action 入队 API | `GMClientMsgProc` | `Game.PlayerAction` | 已覆盖 | 当前无返回值，适合异步游戏指令。 |
| 56 | 主循环 channel 串行化 | 玩家 action 反射分发 | `ProtocolCore` 函数表 | `MethodByName(action).Call` | 部分覆盖 | 需要方法签名校验和 recover，避免反射 panic。 |
| 57 | 主循环 channel 串行化 | User 连接入队 API | GameServer 无 WS 管理连接 | `Game.UserConn` | 已覆盖 | 管理/观察用户和玩家连接分离是合理设计。 |
| 58 | 主循环 channel 串行化 | User action 反射分发 | 管理命令/观察入口语义 | `Game.UserAction` | 部分覆盖 | 当前缺少方法存在性检查，反射调用可能 panic。 |
| 59 | 主循环 channel 串行化 | 命令入队 API | 控制台/GM 命令语义 | `Game.Command` | 已覆盖 | 同步返回适合 HTTP 管理 API，但需防止长事务阻塞主循环。 |
| 60 | 主循环 channel 串行化 | channel 背压策略 | GameServer socket buffer/queue | 固定容量 100/1000 | 部分覆盖 | 需要定义满队列时阻塞、丢弃还是断线，避免网络 goroutine 无限卡住。 |
| 61 | Tick 调度与周期任务 | 100ms ticker 创建 | `CMuTimer::SetMuTimer` | `time.NewTicker(100ms)` | 已覆盖 | Go 版 ticker 替代 C++ timer queue 是合理设计。 |
| 62 | Tick 调度与周期任务 | 100ms 对象处理 | `TimerProcQueue` 驱动对象 Tick | `ObjectManager.Process100ms()` | 已覆盖 | 运行时只调度，具体对象行为放对象系统。 |
| 63 | Tick 调度与周期任务 | 1s 计数器 | GameServer 多定时器 | `cnt%10 == 0` | 部分覆盖 | 后续周期任务增多时应抽象调度器或明确注册表。 |
| 64 | Tick 调度与周期任务 | 天气周期处理 | 地图天气/状态 Tick | `MapManager.ProcessWeather` | 已覆盖 | 地图语义归地图模块，运行时保留调用顺序。 |
| 65 | Tick 调度与周期任务 | 1000ms 对象处理 | 对象状态周期更新 | `ObjectManager.Process1000ms()` | 已覆盖 | 注意耗时逻辑会阻塞全部 channel 请求。 |
| 66 | Tick 调度与周期任务 | 地图物品过期 | 地图道具清理 Timer | `MapManager.ExpireItem(time.Now())` | 已覆盖 | 过期策略归地图/道具模块，运行时只提供时间驱动。 |
| 67 | Tick 调度与周期任务 | 在线百分比注册 Tick | `GameServerInfoSend` 周期上报 | `serverRegisterChan <- MsgServerRegister{Percent: ...}` | 已覆盖 | 如果 channel 满，当前主循环会阻塞，需要非阻塞或超时策略。 |
| 68 | Tick 调度与周期任务 | ticker 资源释放 | `QueueTimer` 释放 | 当前未 `t100ms.Stop()` | 未覆盖 | goroutine 退出前应 `defer t100ms.Stop()`。 |
| 69 | Tick 调度与周期任务 | 慢 Tick 观测 | GameServer 运行状态日志 | 当前有注释掉的耗时打印 | 未覆盖 | 增加慢 Tick 阈值日志，定位对象/地图周期任务卡顿。 |
| 70 | Tick 调度与周期任务 | Tick 测试 | Timer 行为验证 | 当前无测试 | 未覆盖 | 用 fake ticker 或可注入 clock 测试周期任务顺序。 |
| 71 | TCP C1/C2 连接生命周期 | Handler 初始化 | `ProtocolCore` 注册 | `C1C2Handle.init` | 已覆盖 | 入口表构建放 init，运行前失败快速暴露。 |
| 72 | TCP C1/C2 连接生命周期 | ingress 重复 code 检查 | 协议号唯一性 | `apiIns duplicated code` 检查 | 已覆盖 | 保持启动期校验，避免运行期分发歧义。 |
| 73 | TCP C1/C2 连接生命周期 | egress 类型检查 | 发包函数表 | `apiOut msg field must be a pointer` | 已覆盖 | 还可检查重复输出类型和重复 code。 |
| 74 | TCP C1/C2 连接生命周期 | 连接 adapter | socket object wrapper | `type conn struct{ *c1c2.Conn; id int }` | 已覆盖 | adapter 隔离网络层和对象层，设计合理。 |
| 75 | TCP C1/C2 连接生命周期 | 获取远端地址 | socket remote addr | `conn.Addr()` | 已覆盖 | 可补充代理场景真实 IP 处理，但 TCP 通常不需要。 |
| 76 | TCP C1/C2 连接生命周期 | 连接写出 | `DataSocketSend` | `conn.Write` -> `C1C2Handle.marshal` | 已覆盖 | 写失败需要上层决定是否断线。 |
| 77 | TCP C1/C2 连接生命周期 | 连接关闭 | `CloseClient`、`CloseClientHard` | `conn.Close` | 已覆盖 | Close 语义依赖底层 c1c2 包，后续可补充幂等保证。 |
| 78 | TCP C1/C2 连接生命周期 | 新连接回调 | `AcceptSocket` | `C1C2Handle.OnConn` | 已覆盖 | OnConn 同步进入 game 主循环创建玩家对象。 |
| 79 | TCP C1/C2 连接生命周期 | 断线回调 | `GSDisconnect` | `C1C2Handle.OnClose` | 已覆盖 | context 中无 ID 时静默返回，合理但应记录异常诊断。 |
| 80 | TCP C1/C2 连接生命周期 | 底层 IOCP 替代 | `giocp.cpp::ServerWorkerThread` | Go TCP server + goroutine | 设计替代 | 不迁移 IOCP，只保留连接读写、关闭、缓冲和错误处理语义。 |
| 81 | C1/C2 协议路由与响应编码 | 请求 context 取玩家 ID | socket index/object index | `ctx.Value(c1c2.UserContextKey)` | 已覆盖 | ID 缺失直接返回，后续可加 debug 级日志。 |
| 82 | C1/C2 协议路由与响应编码 | 请求体空长度检查 | 协议包长度校验 | `Handle` 当前先读 `req.Body[0]` | 需修正 | 需要先检查 `len(req.Body)>0`，避免空包 panic。 |
| 83 | C1/C2 协议路由与响应编码 | 单字节协议号解析 | `ProtocolCore` head code | `code := int(req.Body[0])` | 已覆盖 | 保持常见 C1 单字节协议快速路径。 |
| 84 | C1/C2 协议路由与响应编码 | 双字节协议号解析 | 扩展协议号 | `binary.BigEndian.Uint16` | 已覆盖 | 当前 body 裁剪逻辑要配合双字节 code 测试验证。 |
| 85 | C1/C2 协议路由与响应编码 | 未知协议处理 | `GMClientMsgProc` default | `invalid api` 日志 | 已覆盖 | 可补充玩家 ID、包头、长度，便于追踪外挂/错误客户端。 |
| 86 | C1/C2 协议路由与响应编码 | 协议调试日志降噪 | GameServer selective log | `switch api.action` 跳过高频行为 | 已覆盖 | 保留高频动作降噪，必要时做采样。 |
| 87 | C1/C2 协议路由与响应编码 | 加密标记校验 | GameServer 加密协议检查 | `api.enc && !req.Encrypt` | 部分覆盖 | 当前只 warn 不 return，需要调用 `27-security.md` 决定拒绝、记录或踢线。 |
| 88 | C1/C2 协议路由与响应编码 | 登录态/权限校验 | 登录阶段协议限制 | 注释中的 `api.level` 校验 | 未覆盖 | 需要恢复或重写 auth level，避免未登录调用玩家协议。 |
| 89 | C1/C2 协议路由与响应编码 | 请求反序列化 | 协议结构 unpack | 反射调用 `Unmarshal` | 已覆盖 | 需要方法签名校验，避免错误模型导致 panic。 |
| 90 | C1/C2 协议路由与响应编码 | 响应编码 | `DataSocketSend` | `marshal` 写 code、body、head | 已覆盖 | 补齐响应表重复检查和 C1/C2 长度边界测试。 |
| 91 | HTTP/WS 运行时入口 | HTTP handler 初始化 | GameServer 无 HTTP 对照 | `HTTPHandle.init` | 已覆盖 | Go 版 HTTP 是管理/调试入口，路由语义和鉴权归 `29-ops.md`。 |
| 92 | HTTP/WS 运行时入口 | gin mode 设置 | 运行环境 debug/release | `gin.SetMode` | 已覆盖 | Debug 配置来自环境变量，合理。 |
| 93 | HTTP/WS 运行时入口 | HTTP validator 初始化 | 管理 API 参数校验 | `validator.New()` | 已覆盖 | 账号 API 校验归账号模块继续补齐。 |
| 94 | HTTP/WS 运行时入口 | HTTP 路由注册 | 管理入口/状态入口 | `/`、`/api/accounts`、`/api/game` | 已覆盖 | 运行时文档只描述入口，不展开账号业务。 |
| 95 | HTTP/WS 运行时入口 | HTTP 错误中间处理 | 统一错误响应 | `setErr`、`handleErr` | 已覆盖 | 错误码定义可在管理 API 模块扩展。 |
| 96 | HTTP/WS 运行时入口 | 调试首页 | GameServer 信息栏替代 | `handleHome` | 已覆盖 | 这是 Go 版可视化调试页，可后续替代 GameServer 状态窗口。 |
| 97 | HTTP/WS 运行时入口 | WebSocket 升级 | GameServer 无 WS 对照 | `handleGame` + `upgrader.Upgrade` | 已覆盖 | WS 作为观察/调试通道合理。 |
| 98 | HTTP/WS 运行时入口 | WebSocket Origin 策略 | 无直接对照 | `CheckOrigin` 永远 true | 需修正 | 生产环境应限制 Origin 或仅在 debug 开启。 |
| 99 | HTTP/WS 运行时入口 | WebSocket 请求路由 | 管理/订阅协议 | `WSHandle.Handle` | 部分覆盖 | 当前 action 表无重复检查，反射调用也需签名保护。 |
| 100 | HTTP/WS 运行时入口 | WebSocket 响应编码 | 管理/订阅响应 | `WSHandle.marshal` | 已覆盖 | 保持 `action/out` 包装格式，后续补错误响应规范。 |
| 101 | Shutdown 与资源释放 | 信号触发关闭 | `WndProc` 关闭分支 | `main.go` signal 分支 | 已覆盖 | 关闭入口集中在 main，结构清晰。 |
| 102 | Shutdown 与资源释放 | server error 触发退出 | 启动/运行错误返回 | `errChan` 分支 `os.Exit(1)` | 部分覆盖 | 直接 `os.Exit` 会跳过清理，应改为进入统一 shutdown。 |
| 103 | Shutdown 与资源释放 | TCP server 关闭 | `CloseClient`/server close | `tcpServer.Close()` | 已覆盖 | 需要确认底层 c1c2 Close 会唤醒 accept/read goroutine。 |
| 104 | Shutdown 与资源释放 | HTTP server 关闭 | 无直接对照 | `httpServer.Close()` | 已覆盖 | 后续可改 `Shutdown(ctx)` 实现优雅等待。 |
| 105 | Shutdown 与资源释放 | Game 关闭入口 | `GameMainFree` | `game.Game.Close()` | 已覆盖 | Close 顺序在 TCP/HTTP 之后，避免新请求进入。 |
| 106 | Shutdown 与资源释放 | 全对象下线命令 | `GSDisconnect` 全员处理 | `OfflineAllObjects` command | 已覆盖 | 使用 command 进入主循环，符合串行化原则。 |
| 107 | Shutdown 与资源释放 | 在线对象等待 | 退出前保存/断线等待 | `GetOnlineObjectsNumber` 轮询 | 部分覆盖 | 需要最大等待时间，避免关闭流程无限阻塞。 |
| 108 | Shutdown 与资源释放 | Bot 清理 | 辅助对象释放 | `BotManager.DeleteAllBots()` | 已覆盖 | 需要确认 Bot 删除不会反向阻塞主循环。 |
| 109 | Shutdown 与资源释放 | context cancel | shutdown flag | `g.cancel()` | 部分覆盖 | 应先释放 ticker/UDP conn，再确保 goroutine 已退出或可观测。 |
| 110 | Shutdown 与资源释放 | 退出等待 | GameServer 资源释放等待 | `time.Sleep(2s)` | 需修正 | 固定 sleep 不可靠，应使用 wait group 或 done channel。 |
| 111 | 运行时测试与诊断 | 启动 smoke test | 人工启动 GameServer | 当前无自动化 | 未覆盖 | 覆盖配置加载、Game.Start、TCP/HTTP 启动的最小路径。 |
| 112 | 运行时测试与诊断 | 配置缺失测试 | `ReadConfig` 错误路径 | 当前难测，因 `os.Exit` | 未覆盖 | 重构配置加载以返回 error，或用子进程测试 exit 行为。 |
| 113 | 运行时测试与诊断 | 端口占用测试 | socket bind 失败 | TCP/HTTP ListenAndServe | 未覆盖 | 验证端口占用进入统一错误处理。 |
| 114 | 运行时测试与诊断 | TCP 连接测试 | `AcceptSocket` 行为 | `C1C2Handle.OnConn/OnClose` | 未覆盖 | 使用 fake c1c2 conn 或集成测试验证连接生命周期。 |
| 115 | 运行时测试与诊断 | C1/C2 协议路由测试 | `ProtocolCore` 分发 | `C1C2Handle.Handle` | 未覆盖 | 覆盖空包、未知协议、双字节协议、Unmarshal 失败。 |
| 116 | 运行时测试与诊断 | WebSocket 路由测试 | 无直接对照 | `handleGame`、`WSHandle.Handle` | 未覆盖 | 覆盖缺 action、缺 in、非法 JSON、未知 action。 |
| 117 | 运行时测试与诊断 | channel 背压测试 | socket send/recv queue | 各 channel 固定容量 | 未覆盖 | 压测队列满时网络 goroutine 和主循环是否会死锁。 |
| 118 | 运行时测试与诊断 | 优雅关闭测试 | `GameMainFree` | `main` + `Game.Close` | 未覆盖 | 验证关闭期间不接收新连接，已有对象可下线。 |
| 119 | 运行时测试与诊断 | goroutine 泄漏测试 | 线程退出检查 | Start/Close goroutine | 未覆盖 | 检查 ticker、UDP、HTTP、TCP 相关 goroutine 是否退出。 |
| 120 | 运行时测试与诊断 | 运行时指标诊断 | GameServer 信息栏/日志 | 当前日志为主 | 未覆盖 | 后续增加在线数、队列长度、Tick 耗时、连接数等观测指标。 |
