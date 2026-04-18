# 18. 好友与邮件系统

本模块覆盖好友列表、好友申请、申请等待、删除好友、在线状态、Memo 邮件写入/读取/删除/列表、好友聊天室邀请，以及登录/下线时的状态同步。聊天系统不单列；好友聊天室只记录邀请和入口，不展开成通用聊天基础设施。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 现有入口 | 好友列表协议 | Friend list protocol | `handle/c1c2.go` 0xC0 | 部分覆盖 | 路由存在，业务未实现。 |
| 2 | 现有入口 | 添加好友协议 | Friend add protocol | `handle` 0xC1 | 部分覆盖 | 路由存在，业务未实现。 |
| 3 | 现有入口 | 好友申请等待协议 | Friend add wait protocol | `handle` 0xC2 | 部分覆盖 | 路由存在，业务未实现。 |
| 4 | 现有入口 | 删除好友协议 | Friend delete protocol | `handle` 0xC3 | 部分覆盖 | 路由存在，业务未实现。 |
| 5 | 现有入口 | 好友状态协议 | Friend state protocol | `handle` 0xC4 | 部分覆盖 | 路由存在，业务未实现。 |
| 6 | 现有入口 | Memo 写入协议 | Memo send/write protocol | `handle` 0xC5 | 部分覆盖 | 路由存在，业务未实现。 |
| 7 | 现有入口 | Memo 读取协议 | Memo read protocol | `handle` 0xC7 | 部分覆盖 | 路由存在，业务未实现。 |
| 8 | 现有入口 | Memo 删除协议 | Memo delete protocol | `handle` 0xC8 | 部分覆盖 | 路由存在，业务未实现。 |
| 9 | 现有入口 | Memo 列表协议 | Memo list protocol | `handle` 0xC9 | 部分覆盖 | 路由存在，业务未实现。 |
| 10 | 现有入口 | 好友聊天室创建 | Friend chat room create | `handle` 0xCA | 部分覆盖 | 路由存在，业务未实现。 |
| 11 | 现有入口 | 聊天室邀请 | Friend room invitation | `handle` 0xCB | 部分覆盖 | 路由存在，业务未实现。 |
| 12 | 数据模型 | 好友关系实体 | Friend table / ExDB friend data | 暂无 | 未覆盖 | 设计账号/角色维度的好友关系；ExDB 通信通道归 `28-external-comm.md`。 |
| 13 | 数据模型 | 好友申请实体 | Friend wait/request data | 暂无 | 未覆盖 | 存储待确认申请和方向。 |
| 14 | 数据模型 | 好友在线状态 | Friend state | `ObjectManager.GetPlayerByName` 可辅助 | 未覆盖 | 维护在线、离线、服务器、角色名状态。 |
| 15 | 数据模型 | Memo 实体 | Memo table / ExDB memo data | 暂无 | 未覆盖 | 存储邮件标题、正文、发送者、时间、已读状态；ExDB 通信通道归 `28-external-comm.md`。 |
| 16 | 数据模型 | Memo 编号 | Memo GUID/number | 暂无 | 未覆盖 | 每封邮件要有稳定 ID 用于读取和删除。 |
| 17 | 数据模型 | Memo 容量限制 | GameServer memo limit | 暂无 | 未覆盖 | 限制收件箱数量，避免无限增长。 |
| 18 | 管理器 | 好友管理器初始化 | Friend system init | 暂无 | 未覆盖 | 初始化好友、申请、Memo 存储访问。 |
| 19 | 管理器 | 登录加载好友 | Login friend load | `Player.LoadCharacter` 未接入 | 未覆盖 | 进入游戏后加载好友列表和待处理状态。 |
| 20 | 管理器 | 下线状态广播 | Friend state offline | `Player.Offline` 未接入 | 未覆盖 | 玩家下线时通知好友。 |
| 21 | 管理器 | 查询在线玩家 | Friend online lookup | `ObjectManager.GetPlayerByName` | 部分覆盖 | 复用对象管理器查找在线角色。 |
| 22 | 好友列表 | 请求好友列表 | Friend list request | `handle` 0xC0 | 未覆盖 | 返回好友名、服务器、在线状态。 |
| 23 | 好友列表 | 好友列表分页/容量 | Friend list response count | 暂无 | 未覆盖 | 按客户端包容量限制返回。 |
| 24 | 好友列表 | 双向好友关系 | Friend relation pair | 暂无 | 未覆盖 | 好友应双向可见，删除需同步双方。 |
| 25 | 添加好友 | 发起添加好友 | Friend add request | `handle` 0xC1 | 未覆盖 | 校验目标存在、非自己、未重复、未满。 |
| 26 | 添加好友 | 目标在线申请 | Friend add online | 暂无 | 未覆盖 | 目标在线时下发确认请求。 |
| 27 | 添加好友 | 目标离线申请 | Friend add offline/wait | 暂无 | 未覆盖 | 目标离线时保存等待申请。 |
| 28 | 添加好友 | 申请等待列表 | Friend add wait | `handle` 0xC2 | 未覆盖 | 目标上线或请求时返回待处理申请。 |
| 29 | 添加好友 | 接受申请 | Friend add accept | 暂无 | 未覆盖 | 建立双向好友关系并通知双方。 |
| 30 | 添加好友 | 拒绝申请 | Friend add reject | 暂无 | 未覆盖 | 删除申请并通知发起者。 |
| 31 | 添加好友 | 重复申请处理 | Friend duplicated request | 暂无 | 未覆盖 | 防止重复插入申请和重复好友。 |
| 32 | 添加好友 | 黑名单/限制预留 | Friend block checks | 暂无 | 未覆盖 | 如后续有黑名单或禁社交状态，应在此校验。 |
| 33 | 删除好友 | 删除好友入口 | Friend delete | `handle` 0xC3 | 未覆盖 | 删除双方好友关系。 |
| 34 | 删除好友 | 删除结果通知 | Friend delete result | 暂无 | 未覆盖 | 通知请求方删除成功或失败。 |
| 35 | 删除好友 | 被删方状态同步 | Friend delete counterpart | 暂无 | 未覆盖 | 被删除方在线时刷新好友列表或状态。 |
| 36 | 状态同步 | 好友上线通知 | Friend state online | `handle` 0xC4 | 未覆盖 | 玩家上线后通知在线好友。 |
| 37 | 状态同步 | 好友下线通知 | Friend state offline | `handle` 0xC4 | 未覆盖 | 玩家下线后通知在线好友。 |
| 38 | 状态同步 | 好友换线通知 | Friend server move | 暂无 | 未覆盖 | 多服/跨服时通过 `28-external-comm.md` 同步服务器编号。 |
| 39 | 状态同步 | 状态查询兜底 | Friend state request | 暂无 | 未覆盖 | 客户端请求状态时返回当前状态。 |
| 40 | Memo 写入 | 发送 Memo | Memo write/send | `handle` 0xC5 | 未覆盖 | 校验收件人、标题/正文长度、容量。 |
| 41 | Memo 写入 | Memo 标题编码 | Memo subject | 暂无 | 未覆盖 | 兼容客户端字符集和长度。 |
| 42 | Memo 写入 | Memo 正文编码 | Memo body | 暂无 | 未覆盖 | 兼容客户端字符集、换行和长度。 |
| 43 | Memo 写入 | Memo 保存 | Memo DB insert | 暂无 | 未覆盖 | 写入持久层。 |
| 44 | Memo 写入 | 新邮件通知 | Memo notify | 暂无 | 未覆盖 | 收件人在线时通知有新邮件。 |
| 45 | Memo 列表 | 请求 Memo 列表 | Memo list | `handle` 0xC9 | 未覆盖 | 返回邮件摘要列表。 |
| 46 | Memo 列表 | Memo 摘要字段 | Memo list response | 暂无 | 未覆盖 | 包含编号、发送者、标题、时间、已读。 |
| 47 | Memo 列表 | Memo 排序 | Memo list sort | 暂无 | 未覆盖 | 按时间倒序或 GameServer 兼容顺序。 |
| 48 | Memo 读取 | 读取 Memo | Memo read | `handle` 0xC7 | 未覆盖 | 按编号读取正文并标记已读。 |
| 49 | Memo 读取 | 读取权限校验 | Memo owner check | 暂无 | 未覆盖 | 只能读取自己的邮件。 |
| 50 | Memo 读取 | 已读状态 | Memo read flag | 暂无 | 未覆盖 | 读取后更新已读标记。 |
| 51 | Memo 删除 | 删除 Memo | Memo delete | `handle` 0xC8 | 未覆盖 | 按编号删除邮件。 |
| 52 | Memo 删除 | 批量删除预留 | Memo delete list | 暂无 | 未覆盖 | 如果客户端支持多封删除，应统一接口。 |
| 53 | Memo 删除 | 删除权限校验 | Memo owner check | 暂无 | 未覆盖 | 只能删除自己的邮件。 |
| 54 | 聊天室 | 创建好友聊天室 | Friend chat room create | `handle` 0xCA | 未覆盖 | 记录客户端入口和邀请结果，具体聊天房可后续外部化。 |
| 55 | 聊天室 | 好友聊天室邀请 | Friend room invitation | `handle` 0xCB | 未覆盖 | 邀请好友进入聊天室。 |
| 56 | 聊天室 | 邀请目标校验 | Friend relation check | 暂无 | 未覆盖 | 只能邀请好友或符合规则的目标。 |
| 57 | 聊天室 | 聊天室状态清理 | Friend chat room cleanup | 暂无 | 未覆盖 | 下线或退出时清理临时邀请状态。 |
| 58 | 跨系统联动 | 登录后数据加载 | Friend/Guild/Gens requests | `03-characters.md` 已记录入口 | 未覆盖 | 角色系统只触发加载，本模块负责好友邮件数据。 |
| 59 | 跨系统联动 | 私聊目标辅助 | Whisper lookup | `Object.Whisper` 已有在线查找 | 部分覆盖 | 好友系统可补充离线状态显示，不替代私聊入口。 |
| 60 | 测试与校验 | 好友申请流程测试 | Friend add/wait/delete | 暂无 | 未覆盖 | 覆盖在线申请、离线申请、接受、拒绝、删除。 |
| 61 | 测试与校验 | 好友状态同步测试 | Friend state | 暂无 | 未覆盖 | 覆盖上线、下线、换线、重复通知。 |
| 62 | 测试与校验 | Memo 流程测试 | Memo write/list/read/delete | 暂无 | 未覆盖 | 覆盖写入、容量、读取、删除、权限校验。 |
