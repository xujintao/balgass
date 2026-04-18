# 29. 运营管理系统

本模块覆盖 GM 权限、GM 命令、HTTP/WS 后台管理入口、公告与系统消息、维护控制、在线统计、人工踢人/封禁/禁言入口、运营日志和审计。运营管理系统不拥有账号认证、角色生命周期、对象行为、安全自动判定、活动状态机、外部通信、VIP 权益效果或 CashShop 商城交易；它只提供后台入口、权限控制、操作编排和审计记录，并调用对应业务系统执行。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 模块边界与总入口 | OpsManager 总入口 | `CGMMng`、`TNotice`、`TServerInfoDisplayer`、日志对象分散 | 暂无 `game/ops` | 未覆盖 | 建立统一运营管理入口，管理 GM、后台 API、公告、统计、维护和审计。 |
| 2 | 模块边界与总入口 | 运营系统初始化 | `CGMMng::Init`、`ManagerInit` | `HTTPHandle.init`、`Command` 分散 | 部分覆盖 | 启动时加载 GM、命令、公告配置、审计 logger 和后台权限策略。 |
| 3 | 模块边界与总入口 | 运营系统关闭 | GameMainFree 释放日志/状态 | 暂无 | 未覆盖 | 关闭时停止维护倒计时、刷出审计日志、拒绝新后台操作。 |
| 4 | 模块边界与总入口 | 运营动作上下文 | `lpObj`、GM 名称、AuthLevel | HTTP/Command 无统一上下文 | 未覆盖 | 每次后台/GM 操作携带操作者、来源、目标、理由、trace id。 |
| 5 | 模块边界与总入口 | 运营结果模型 | 命令返回/MsgOutput 分散 | Command 返回 any/error | 部分覆盖 | 统一成功、失败、拒绝、部分成功和异步已受理结果。 |
| 6 | 模块边界与总入口 | 运营领域接口 | `ManagementProc` 调业务函数 | `game.Command` 反射 | 部分覆盖 | 运营系统只编排调用，不直接修改业务内部状态。 |
| 7 | 模块边界与总入口 | 单服/多服运营边界 | GameServer 多服命令分散 | 暂无 | 未覆盖 | 单服直接执行，多服通过 `28-external-comm.md` 广播/投递。 |
| 8 | 模块边界与总入口 | VIP 边界引用 | `VipSys` | VIP 字段/配置零散 | 未覆盖 | 只保留 VIP 配置和授权入口，VIP 效果后续独立盘。 |
| 9 | 模块边界与总入口 | CashShop 边界引用 | `CashShop`、InGameShop GDReq | 暂无完整商城 | 未覆盖 | 只保留商城管理和审计入口，购买/货币/回滚后续独立盘。 |
| 10 | 模块边界与总入口 | 运营回归测试夹具 | GameServer 靠运行验证 | 暂无 | 未覆盖 | 提供 fake GM、fake HTTP、fake 对象和审计断言。 |
| 11 | GM 权限与身份 | GM 权限等级枚举 | `AuthLevel`、GM 文件 | `handle.AuthLevel` 有 `GM/Admin` | 部分覆盖 | 统一 Guest/Player/GM/Admin/Owner 等权限等级。 |
| 12 | GM 权限与身份 | GM 列表加载 | `LoadGMFile`、`GM_DATA` | 暂无 | 未覆盖 | 加载 GM 名称、过期时间、权限等级。 |
| 13 | GM 权限与身份 | GM 过期时间 | `GM_DATA::ExpiryDate` | 暂无 | 未覆盖 | 过期 GM 自动降权并审计。 |
| 14 | GM 权限与身份 | GM 名称匹配 | `ManagerAdd`、GM name | 暂无 | 未覆盖 | 角色名和 GM 配置匹配时授予权限。 |
| 15 | GM 权限与身份 | 在线 GM 注册 | `ManagerAdd` | 暂无 | 未覆盖 | GM 上线后登记在线 GM 列表。 |
| 16 | GM 权限与身份 | 在线 GM 移除 | `ManagerDel` | 暂无 | 未覆盖 | GM 下线后清理在线 GM 列表。 |
| 17 | GM 权限与身份 | GM 权限缓存 | `GameMaster` 字段 | `Object.GameMaster` 注释字段 | 未覆盖 | 玩家对象保存 GM 权限快照，避免每次查配置。 |
| 18 | GM 权限与身份 | GM 操作来源 | chat command / backend | HTTP 无鉴权 | 未覆盖 | 区分游戏内 GM 命令、HTTP 后台、WS 调试、控制台。 |
| 19 | GM 权限与身份 | HTTP 管理鉴权 | GameServer 无 HTTP | 当前 `/api/accounts` 无鉴权 | 需修正 | 后台 API 必须接入 token/session/内网策略。 |
| 20 | GM 权限与身份 | WS 管理鉴权 | GameServer 无 WS | WS `CheckOrigin` 永远 true | 需修正 | WS 观察/管理通道必须限制来源和权限。 |
| 21 | GM 权限与身份 | 命令权限矩阵 | `command.level` | 暂无 | 未覆盖 | 每条 GM/后台命令声明最低权限。 |
| 22 | GM 权限与身份 | 目标权限保护 | GM 互相操作限制 | 暂无 | 未覆盖 | 低权限 GM 不能操作高权限 GM 或 Owner。 |
| 23 | GM 权限与身份 | GM 隐身状态 | Hide/UnHide 命令 | 暂无 | 未覆盖 | 隐身只影响客户端可见性，不绕过审计。 |
| 24 | GM 权限与身份 | GM Chat 颜色 | `GMChat` 配置 | `conf.Common.General.GMChat` | 部分覆盖 | GM 发言使用配置颜色和标记。 |
| 25 | GM 权限与身份 | GM 权限测试 | GameServer 无统一 | 暂无 | 未覆盖 | 覆盖权限不足、过期、目标保护和后台鉴权。 |
| 26 | GM 命令体系 | 命令管理器 | `CGMCommand`、`CGMMng` | `game/cmd.Command` | 部分覆盖 | 建立命令注册表，替代裸反射方法名。 |
| 27 | GM 命令体系 | 命令配置加载 | `LoadCommandFile`、`LoadCommandXML` | `conf` 读 Post/PKClear 等 | 部分覆盖 | 加载命令名、权限、用法、开关。 |
| 28 | GM 命令体系 | 命令查找 | `GetCmd`、`findCommand` | 反射 MethodByName | 部分覆盖 | 命令不存在返回明确错误，不 panic。 |
| 29 | GM 命令体系 | 命令参数解析 | `GetTokenString/Number` | 暂无通用解析 | 未覆盖 | 支持字符串、数字、玩家名、坐标、地图、理由。 |
| 30 | GM 命令体系 | 命令用法帮助 | `command.usage` | 暂无 | 未覆盖 | 权限不足或参数错误时返回用法。 |
| 31 | GM 命令体系 | Post 命令 | `PostSend` | `PostCMD` 配置已读 | 未覆盖 | 全服/频道发言命令，聊天业务执行发送。 |
| 32 | GM 命令体系 | GlobalPost 命令 | `GlobalPostSend`、`GS_GDReqGlobalPostMultiCast` | `GlobalPostCMD` 配置已读 | 未覆盖 | 跨服全局发言通过外部通信系统投递。 |
| 33 | GM 命令体系 | Item 命令 | `ItemCMD` | 暂无 | 未覆盖 | GM 发物品入口，实际创建归道具系统。 |
| 34 | GM 命令体系 | ServerInfo 命令 | `SvInfoCMD` | `GetOnlineObjectsNumber` | 部分覆盖 | 查询在线、地图人数、服务状态。 |
| 35 | GM 命令体系 | CharacterInfo 命令 | `CharInfoCMD` | 暂无 | 未覆盖 | 查询角色基础状态、地图、账号、背包摘要。 |
| 36 | GM 命令体系 | Hide 命令 | `HideCMD` | 暂无 | 未覆盖 | GM 隐身，调用对象视野系统。 |
| 37 | GM 命令体系 | UnHide 命令 | `UnHideCMD` | 暂无 | 未覆盖 | 取消 GM 隐身并刷新视野。 |
| 38 | GM 命令体系 | ClearInventory 命令 | `ClearInvCMD` | 暂无 | 未覆盖 | 清背包入口，实际删除归道具系统。 |
| 39 | GM 命令体系 | AddStat 命令 | `AddSTR/AGI/VIT/ENE/CMD` | 暂无 | 未覆盖 | 增加属性点入口，角色/公式系统执行重算。 |
| 40 | GM 命令体系 | Online 命令 | `ONLINECMD` | `GetOnlineObjectsNumber` | 部分覆盖 | 查询全服在线和分类在线。 |
| 41 | GM 命令体系 | Warehouse 命令 | `WARECMD` | 暂无 | 未覆盖 | 打开/查看仓库入口，仓库业务另属道具/账号。 |
| 42 | GM 命令体系 | PKClear 命令 | `PKCLEARCMD` | `PKClearCommand` 配置已读 | 未覆盖 | 清 PK 入口，PK/Gens 业务执行。 |
| 43 | GM 命令体系 | PKSet 命令 | `PKSETCMD` | 暂无 | 未覆盖 | 设置 PK 状态入口，PK 系统执行。 |
| 44 | GM 命令体系 | Skin 命令 | `SKINCMD` | 暂无 | 未覆盖 | 改变外观入口，对象系统执行广播。 |
| 45 | GM 命令体系 | Watch 命令 | `Watch`、`WatchTargetIndex` | 暂无 | 未覆盖 | 观察目标玩家状态和位置。 |
| 46 | GM 命令体系 | Trace 命令 | `Trace` | 暂无 | 未覆盖 | 传送到目标玩家位置，对象/地图系统执行。 |
| 47 | GM 命令体系 | DC 命令 | `DC` | `OfflineAllObjects` 只有全体 | 未覆盖 | 踢指定玩家，账号/对象系统执行断线。 |
| 48 | GM 命令体系 | GuildDC 命令 | `GuildDC` | 暂无 | 未覆盖 | 踢某战盟成员，战盟系统提供成员列表。 |
| 49 | GM 命令体系 | Move 命令 | `MoveCMD` | 暂无 | 未覆盖 | 移动指定玩家，对象/地图系统执行。 |
| 50 | GM 命令体系 | GlobalMove 命令 | `gMoveCMD` | 暂无 | 未覆盖 | 批量移动玩家，地图系统执行目标校验。 |
| 51 | GM 命令体系 | ChatBan 命令 | `ChatBan` | `ChatBlockTime` 注释字段 | 未覆盖 | 人工禁言入口，安全/聊天系统执行限制。 |
| 52 | GM 命令体系 | ChatUnban 命令 | `ChatUnban` | 暂无 | 未覆盖 | 解除禁言入口。 |
| 53 | GM 命令体系 | BanAcc 命令 | `BanAccCMD`、`GDReqBanUser` | 暂无 | 未覆盖 | 封账号入口，账号/外部通信系统执行。 |
| 54 | GM 命令体系 | UnBanAcc 命令 | `UnBanAccCMD` | 暂无 | 未覆盖 | 解封账号入口。 |
| 55 | GM 命令体系 | BanChar 命令 | `BanCharCMD` | 暂无 | 未覆盖 | 封角色入口，角色系统执行。 |
| 56 | GM 命令体系 | UnBanChar 命令 | `UnBanCharCMD` | 暂无 | 未覆盖 | 解封角色入口。 |
| 57 | GM 命令体系 | EventStart 命令 | `BCStart/CCStart/DSStart/ITStart` | 暂无 | 未覆盖 | 副本活动手动启动入口，副本系统执行。 |
| 58 | GM 命令体系 | WorldEvent 命令 | `CWSetState/CSSetOwner/KTSet*` | 暂无 | 未覆盖 | 世界事件运营命令入口，世界事件系统执行。 |
| 59 | GM 命令体系 | GremoryGift 命令 | `GremoryGiftCMD` | 暂无 | 未覆盖 | 奖励发放入口，奖励/道具系统执行。 |
| 60 | GM 命令体系 | 命令审计包装 | `GMLog` | 暂无 | 未覆盖 | 所有命令执行前后写审计，含参数和结果。 |
| 61 | 后台管理 API | HTTP Engine 初始化 | GameServer 无 HTTP | `HTTPHandle.init` | 已覆盖 | 运营系统管理后台路由注册和鉴权。 |
| 62 | 后台管理 API | 账号创建 API | GameServer 后台另行实现 | `POST /api/accounts` | 已覆盖 | API 入口归运营系统，账号创建业务归 `02-accounts.md`。 |
| 63 | 后台管理 API | 账号列表 API | GameServer 后台另行实现 | `GET /api/accounts` | 已覆盖 | API 入口归运营系统，账号查询业务归账号系统。 |
| 64 | 后台管理 API | 账号删除 API | GameServer 后台另行实现 | `DELETE /api/accounts/:id` | 已覆盖 | API 入口归运营系统，删除约束归账号系统。 |
| 65 | 后台管理 API | 游戏状态 API | ServerInfoDisplayer | `/api/game` WS | 部分覆盖 | 提供在线、地图、对象、服务状态查看。 |
| 66 | 后台管理 API | 后台错误模型 | 无统一 HTTP | `ConfigError` | 部分覆盖 | 统一后台错误码、HTTP 状态码和审计字段。 |
| 67 | 后台管理 API | 参数校验 | 后台/GM 分散 | `validator.New()` | 部分覆盖 | 所有后台 API 统一参数校验。 |
| 68 | 后台管理 API | 后台分页查询 | 无 | 当前账号列表无分页 | 未覆盖 | 账号、角色、日志、在线列表需要分页。 |
| 69 | 后台管理 API | 后台搜索过滤 | 无 | 账号按 email 简单过滤 | 部分覆盖 | 支持账号、角色、IP、地图、状态等条件。 |
| 70 | 后台管理 API | 后台操作 reason | GM 命令日志 | 暂无 | 未覆盖 | 封禁、删除、踢线等敏感操作必须带理由。 |
| 71 | 后台管理 API | 后台 CSRF/Origin | 无 HTTP | WS CheckOrigin true | 需修正 | 生产环境限制 Origin、鉴权和内网访问。 |
| 72 | 后台管理 API | 后台审计中间件 | GMLog | 暂无 | 未覆盖 | 所有后台 API 统一记录操作者、IP、请求体摘要和结果。 |
| 73 | 后台管理 API | 后台权限中间件 | GM level | 暂无 | 未覆盖 | 根据 API 和命令权限矩阵判断访问。 |
| 74 | 后台管理 API | 后台只读模式 | maintenance/readonly | 暂无 | 未覆盖 | 维护时允许查询，禁止写操作。 |
| 75 | 后台管理 API | 后台测试 | 无 | 当前少量 handler 测试 | 未覆盖 | 覆盖鉴权、参数、错误、审计和业务调用。 |
| 76 | 公告与系统消息 | Notice 消息结构 | `PMSG_NOTICE` | 当前无统一公告模型 | 未覆盖 | 定义公告类型、颜色、次数、延迟、速度和文本。 |
| 77 | 公告与系统消息 | MakeNoticeMsg | `TNotice::MakeNoticeMsg` | 暂无 | 未覆盖 | 构造普通公告包。 |
| 78 | 公告与系统消息 | MakeNoticeMsgEx | `TNotice::MakeNoticeMsgEx` | 暂无 | 未覆盖 | 支持格式化公告文本。 |
| 79 | 公告与系统消息 | SetNoticeProperty | `TNotice::SetNoticeProperty` | 暂无 | 未覆盖 | 设置公告颜色、速度、次数和延迟。 |
| 80 | 公告与系统消息 | SendNoticeToAllUser | `TNotice::SendNoticeToAllUser` | 暂无 | 未覆盖 | 全服公告，推送由对象系统执行。 |
| 81 | 公告与系统消息 | SendNoticeToUser | `TNotice::SendNoticeToUser` | 暂无 | 未覆盖 | 指定玩家公告。 |
| 82 | 公告与系统消息 | AllSendServerMsg | `TNotice::AllSendServerMsg` | 暂无 | 未覆盖 | 全服系统消息入口。 |
| 83 | 公告与系统消息 | GCServerMsgStringSend | `TNotice::GCServerMsgStringSend` | `PushSystemMsg` 部分能力 | 部分覆盖 | 系统消息复用对象推送能力。 |
| 84 | 公告与系统消息 | 维护倒计时公告 | GameServer close countdown | 暂无 | 未覆盖 | 维护关闭前周期广播剩余时间。 |
| 85 | 公告与系统消息 | 活动运营公告 | event notice | 事件模块分散 | 未覆盖 | GM 开关活动时统一公告。 |
| 86 | 公告与系统消息 | 跨服公告 | `GS_GDReqGlobalPostMultiCast` | 暂无 | 未覆盖 | 多服公告通过 `28-external-comm.md` 投递。 |
| 87 | 公告与系统消息 | 公告模板 | Lang/TNotice text | 暂无 | 未覆盖 | 支持可配置模板和参数填充。 |
| 88 | 公告与系统消息 | 公告限流 | 防刷屏策略 | 暂无 | 未覆盖 | 防止 GM 或后台误发高频公告。 |
| 89 | 公告与系统消息 | 公告审计 | GMLog/MsgLog | 暂无 | 未覆盖 | 记录公告内容、范围、操作者和发送结果。 |
| 90 | 公告与系统消息 | 公告测试 | 无 | 暂无 | 未覆盖 | 覆盖单人、全服、跨服、颜色和长度限制。 |
| 91 | 在线统计与服务状态 | 在线数统计 | `TServerInfoDisplayer` | `GetOnlineObjectsNumber` | 部分覆盖 | 统计玩家、user、怪物、总在线。 |
| 92 | 在线统计与服务状态 | 在线百分比 | ServerInfoDisplayer percent | `GetPlayerPercent` | 已覆盖 | 提供后台展示和 ConnectServer 注册输入。 |
| 93 | 在线统计与服务状态 | 地图人数统计 | Map user count | 暂无完整 API | 未覆盖 | 按地图统计玩家、怪物、NPC、物品。 |
| 94 | 在线统计与服务状态 | 职业分布统计 | 无统一 | 暂无 | 未覆盖 | 按职业、等级段、地图统计在线玩家。 |
| 95 | 在线统计与服务状态 | 服务连接状态 | Join/Data/ExDB display | `28-external-comm.md` 状态待实现 | 未覆盖 | 后台展示外部服务连接状态。 |
| 96 | 在线统计与服务状态 | 运行时间统计 | server uptime | 暂无 | 未覆盖 | 展示启动时间、运行时长、版本和配置路径。 |
| 97 | 在线统计与服务状态 | Tick 延迟统计 | 无统一 | 暂无 | 未覆盖 | 展示主循环 tick 耗时和积压。 |
| 98 | 在线统计与服务状态 | channel 积压统计 | 无 | 暂无 | 未覆盖 | 展示 player/user/command channel 长度。 |
| 99 | 在线统计与服务状态 | goroutine 统计 | 无 | runtime 可读 | 未覆盖 | 后台展示 goroutine 数和内存。 |
| 100 | 在线统计与服务状态 | DB 状态统计 | DataServer/DB | `model.DB` | 未覆盖 | 展示 DB 连通性、错误数和慢查询。 |
| 101 | 在线统计与服务状态 | 后台订阅推送 | ServerInfoDisplayer UI | WS debug page | 部分覆盖 | WS 推送地图和对象状态。 |
| 102 | 在线统计与服务状态 | 状态快照 API | ServerInfoDisplayer | 暂无统一 JSON | 未覆盖 | 提供只读状态快照接口。 |
| 103 | 在线统计与服务状态 | 异常状态告警 | log/msgbox | 暂无 | 未覆盖 | 外部服务断开、在线异常、队列积压时告警。 |
| 104 | 在线统计与服务状态 | 统计权限 | GM/Admin | 暂无 | 未覆盖 | 部分统计只允许管理员查看。 |
| 105 | 在线统计与服务状态 | 统计测试 | 无 | 暂无 | 未覆盖 | 覆盖空服、满服、多地图和外部服务状态。 |
| 106 | 运营处罚入口 | 踢玩家入口 | `DC`、CloseClient | `OfflineAllObjects` 只有全体 | 未覆盖 | 踢指定玩家，执行归对象/账号系统。 |
| 107 | 运营处罚入口 | 批量踢人入口 | `GuildDC`、维护关闭 | `OfflineAllObjects` | 部分覆盖 | 按战盟、地图、账号列表或全服批量踢人。 |
| 108 | 运营处罚入口 | 封账号入口 | `BanAccCMD`、`GDReqBanUser` | 暂无 | 未覆盖 | 后台/GM 发起封账号，账号和外部通信系统执行。 |
| 109 | 运营处罚入口 | 解封账号入口 | `UnBanAccCMD` | 暂无 | 未覆盖 | 解除账号封禁。 |
| 110 | 运营处罚入口 | 封角色入口 | `BanCharCMD` | 暂无 | 未覆盖 | 后台/GM 发起封角色，角色系统执行。 |
| 111 | 运营处罚入口 | 解封角色入口 | `UnBanCharCMD` | 暂无 | 未覆盖 | 解除角色封禁。 |
| 112 | 运营处罚入口 | 禁言入口 | `ChatBan`、`ChatBlockTime` | 注释字段 | 未覆盖 | 写入禁言状态，聊天/安全系统执行限制。 |
| 113 | 运营处罚入口 | 解禁言入口 | `ChatUnban` | 暂无 | 未覆盖 | 清除禁言状态。 |
| 114 | 运营处罚入口 | 账号物品锁入口 | `AccountItemBlock` | 注释字段 | 未覆盖 | 锁定敏感物品操作，交易/道具系统执行。 |
| 115 | 运营处罚入口 | 安全码锁入口 | `GDReqSecLock` | 暂无 | 未覆盖 | 发起安全锁/二级密码相关后台操作。 |
| 116 | 运营处罚入口 | 处罚 reason | GM logs | 暂无 | 未覆盖 | 所有处罚必须记录原因、期限和操作者。 |
| 117 | 运营处罚入口 | 临时处罚过期 | Penalty time | `27-security.md` 待实现 | 未覆盖 | 运营入口设置期限，安全/账号系统执行过期恢复。 |
| 118 | 运营处罚入口 | 跨服处罚同步 | `GDReqBanUser` | `28-external-comm.md` | 未覆盖 | 多服封禁/踢线通过外部通信投递。 |
| 119 | 运营处罚入口 | 处罚查询 API | 后台查询 | 暂无 | 未覆盖 | 查询账号/角色当前处罚状态和历史。 |
| 120 | 运营处罚入口 | 处罚测试 | 无 | 暂无 | 未覆盖 | 覆盖踢人、封禁、禁言、过期、权限不足和审计。 |
| 121 | 维护与服务器控制 | 维护模式开关 | GameServer close control | 暂无 | 未覆盖 | 进入维护后拒绝新登录，可保留在线玩家。 |
| 122 | 维护与服务器控制 | 关闭倒计时 | `GameServerInfoSendStop`/close countdown | 暂无 | 未覆盖 | 发起全服关闭倒计时并公告。 |
| 123 | 维护与服务器控制 | 优雅踢下线 | `gObjAllLogOut` | `OfflineAllObjects` | 部分覆盖 | 倒计时结束后保存并踢下线。 |
| 124 | 维护与服务器控制 | 禁止建角开关 | `gCreateCharacter` | 暂无 | 未覆盖 | 维护期间禁止创建角色。 |
| 125 | 维护与服务器控制 | 禁止交易开关 | TradeBlock config | 交易系统待实现 | 未覆盖 | 运营入口改变交易可用状态。 |
| 126 | 维护与服务器控制 | 禁止商店开关 | shop enable config | 商店系统待实现 | 未覆盖 | 运营入口改变商店可用状态。 |
| 127 | 维护与服务器控制 | 配置热重载入口 | `ReloadEvent`、LoadFile | 暂无 | 未覆盖 | 发起配置重载，具体模块执行解析。 |
| 128 | 维护与服务器控制 | 怪物重载入口 | `GameMonsterAllCloseAndReLoad` | monster spawn 当前 init | 未覆盖 | 运营入口触发怪物重载，对象/AI 系统执行。 |
| 129 | 维护与服务器控制 | 地图重载入口 | `LoadMapFile` | maps init | 未覆盖 | 运营入口触发地图配置重载。 |
| 130 | 维护与服务器控制 | 脚本重载入口 | Lua/Bag reload | `26-script.md` | 未覆盖 | 运营入口触发脚本系统重载。 |
| 131 | 维护与服务器控制 | 服务状态切换审计 | GMLog | 暂无 | 未覆盖 | 维护、关闭、重载都写审计。 |
| 132 | 维护与服务器控制 | 控制命令幂等 | 无统一 | 暂无 | 未覆盖 | 重复维护/关闭/重载请求返回一致结果。 |
| 133 | 维护与服务器控制 | 控制失败回滚 | 无统一 | 暂无 | 未覆盖 | 某模块重载失败时保留旧配置。 |
| 134 | 维护与服务器控制 | 控制权限 | Admin/Owner | 暂无 | 未覆盖 | 维护和关闭必须高权限。 |
| 135 | 维护与服务器控制 | 维护测试 | 无 | 暂无 | 未覆盖 | 覆盖维护开关、倒计时、拒绝登录和审计。 |
| 136 | 活动运营控制 | 普通活动启停入口 | `Start_Menual/End_Menual` | `22-events.md` | 未覆盖 | GM 发起，普通活动系统执行。 |
| 137 | 活动运营控制 | 副本启停入口 | `BCStart/CCStart/DSStart/ITStart` | `21-dungeons.md` | 未覆盖 | GM 发起，副本系统执行。 |
| 138 | 活动运营控制 | CastleSiege 控制入口 | `CSSet*` 命令 | `23-world-events.md` | 未覆盖 | 设置攻城状态、城主、注册、开始结束。 |
| 139 | 活动运营控制 | Crywolf 控制入口 | `CWSetStateCMD` | `23-world-events.md` | 未覆盖 | 设置狼魂要塞状态。 |
| 140 | 活动运营控制 | Kanturu 控制入口 | `KTSetStandby/Maya/Tower` | `23-world-events.md` | 未覆盖 | 设置坎特鲁阶段。 |
| 141 | 活动运营控制 | 活动刷怪入口 | event monster spawn | 对象/AI 系统 | 未覆盖 | GM 触发刷怪，活动系统返回规则，对象系统生成。 |
| 142 | 活动运营控制 | 活动清场入口 | ClearMonster | 对象系统 | 未覆盖 | GM 触发清理活动怪和临时物品。 |
| 143 | 活动运营控制 | 活动奖励补发入口 | event reward logs | 暂无 | 未覆盖 | 后台补发奖励，奖励/道具系统执行。 |
| 144 | 活动运营控制 | 活动状态查询 | event state | `22/23` 待实现 | 未覆盖 | 后台展示各活动状态、时间和参与人数。 |
| 145 | 活动运营控制 | 活动操作审计 | GMLog | 暂无 | 未覆盖 | 记录手动开关、跳阶段、补发奖励。 |
| 146 | 日志与审计 | GMLog | `GMLog` | 暂无 | 未覆盖 | 记录 GM 命令和后台敏感操作。 |
| 147 | 日志与审计 | MsgLog | `MsgLog` | 暂无 | 未覆盖 | 记录公告、系统消息和重要广播。 |
| 148 | 日志与审计 | TradeLog 引用 | `TradeLog` | `19-trade.md` | 未覆盖 | 交易日志业务归交易系统，运营系统提供查看/导出入口。 |
| 149 | 日志与审计 | AntiHackLog 引用 | `AntiHackLog` | `27-security.md` | 未覆盖 | 安全日志归安全系统，运营系统提供查询入口。 |
| 150 | 日志与审计 | BotShopLog 引用 | `BotShopLog` | Bot/商店待拆 | 未覆盖 | Bot 商店日志后续归对应系统，运营系统提供查询入口。 |
| 151 | 日志与审计 | SerialCheck 引用 | `SerialCheck` | 道具系统待实现 | 未覆盖 | 序列号审计归道具/安全，运营系统提供查询入口。 |
| 152 | 日志与审计 | 账号操作日志 | 后台/GM | 当前无 | 未覆盖 | 创建、删除、封禁、解封账号写日志。 |
| 153 | 日志与审计 | 角色操作日志 | 后台/GM | 当前无 | 未覆盖 | 踢人、封禁、移动、改状态、补偿写日志。 |
| 154 | 日志与审计 | 日志查询 API | 无 | 暂无 | 未覆盖 | 后台按时间、操作者、目标、类型查询。 |
| 155 | 日志与审计 | 日志导出 | 无 | 暂无 | 未覆盖 | 支持导出审计日志。 |
| 156 | 跨系统接口 | AccountOps 接口 | BanAcc/UnBanAcc | `02-accounts.md` | 未覆盖 | 创建/查询/删除/封禁账号调用账号系统。 |
| 157 | 跨系统接口 | CharacterOps 接口 | CharInfo/BanChar | `03-characters.md` | 未覆盖 | 查角色、封角色、解封角色调用角色系统。 |
| 158 | 跨系统接口 | ObjectOps 接口 | Move/DC/Trace/Hide | `04-objects.md` | 未覆盖 | 踢人、移动、召唤、隐身调用对象系统。 |
| 159 | 跨系统接口 | SecurityOps 接口 | ChatBan/Ban | `27-security.md` | 未覆盖 | 人工处罚入口调用安全系统写限制和审计。 |
| 160 | 测试与回归 | 运营回归测试总集 | GameServer 靠运行验证 | 暂无 | 未覆盖 | 覆盖 GM 权限、命令、HTTP 管理、公告、统计、处罚、维护、活动控制和审计。 |
