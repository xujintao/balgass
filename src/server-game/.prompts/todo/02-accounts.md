# 2. 账号系统

本模块覆盖账号登录、登出、登录态、账号模型、登录包解析、账号认证、安全校验、账号在线唯一性、账号管理 API、账号仓库归属与账号态清理。角色列表、创建/删除/检查/加载角色、角色保存和进入游戏初始化归角色系统。JoinServer/DataServer/中心在线态通信归 `28-external-comm.md`，账号系统只定义登录业务语义和本地状态约束。后台账号创建/查询/删除、封号/解封入口和操作审计归 `29-ops.md`。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 协议入口与登录态 | 登录协议路由 | `protocol.cpp::ProtocolCore` 分发 `0xF1:0x01` | `handle/c1c2.go` `0xF101 -> Login` | 已覆盖 | 保持协议表到 `Player.Login` 的映射，补充未登录状态下的最小可调用协议集合测试。 |
| 2 | 协议入口与登录态 | 登出协议路由 | `ProtocolCore` 分发 `0xF1:0x02` | `0xF102 -> Logout`、`Player.Logout` | 已覆盖 | 明确 close game、返回选角、返回选服三种 flag 的状态流转。 |
| 3 | 协议入口与登录态 | Guest 权限等级 | `PLAYER_CONNECTED` 可登录 | `apiIn.level == Guest` | 部分覆盖 | `handle/c1c2.go` 权限校验被注释，需要恢复或重写。 |
| 4 | 协议入口与登录态 | Player 权限等级 | `PLAYER_LOGGED/PLAYER_PLAYING` 区分 | `apiIn.level == Player` | 需修正 | 登录后和进游戏后可调用协议不同，不能只用一个 Player 等级粗略处理。 |
| 5 | 协议入口与登录态 | 连接状态枚举 | `PLAYER_CONNECTED/LOGGED/PLAYING` | `object.ConnectStateLogged/Playing` | 部分覆盖 | 补齐初始连接态、已登录态、游戏中态的显式状态机文档和测试。 |
| 6 | 登录请求处理 | 登录包结构 | `PMSG_IDPASS` | `model.MsgLogin` | 已覆盖 | 保持账号、密码、HWID、TickCount、Version、Serial 字段完整解析。 |
| 7 | 登录请求处理 | 账号解码 | `BuxConvert` 处理账号 | `MsgLogin.Unmarshal` 中 `utils.Xor(account)` | 已覆盖 | 增加账号长度、空账号、非法字符测试。 |
| 8 | 登录请求处理 | 密码解码 | `BuxConvert` 处理密码 | `MsgLogin.Unmarshal` 中 `utils.Xor(password)` | 已覆盖 | 密码明文比较只是当前实现，后续应接入 hash 或兼容旧服加密策略。 |
| 9 | 登录请求处理 | HWID 解析 | `ProcessClientHWID` | `MsgLogin.Unmarshal` 读取 100 字节 HWID | 部分覆盖 | 当前只解析不校验，需补机器码限制和审计字段。 |
| 10 | 登录请求处理 | TickCount 解析 | `CSPJoinIdPassRequest` 记录 `TickCount` | `MsgLogin.TickCount` | 已覆盖 | 后续安全模块可使用该字段做速度/重放检查。 |
| 11 | 登录请求处理 | 客户端版本解析 | `CliVersion` 与 `szClientVersion` 比对 | `MsgLogin.Version` | 部分覆盖 | 当前未校验版本，需按配置决定是否拒绝。 |
| 12 | 登录请求处理 | 客户端 Serial 解析 | `CliSerial` 与 `szGameServerExeSerial` 比对 | `MsgLogin.Serial` | 部分覆盖 | 当前未校验 serial，需补配置项和失败结果码。 |
| 13 | 登录请求处理 | 登录处理函数 | `GameProtocol::CSPJoinIdPassRequest` | `Player.Login` | 部分覆盖 | Go 版直接查 DB，后续应通过 `28-external-comm.md` 接入 JoinServer 异步认证和中心在线态。 |
| 14 | 登录请求处理 | 登录回复结构 | `GCJoinResult` | `MsgLoginReply` | 已覆盖 | 结果码注释较完整，需保证所有失败路径按约定返回。 |
| 15 | 登录请求处理 | 登录回复编码 | `PMSG_RESULT` `F1:01` | `MsgLoginReply.Marshal` | 已覆盖 | 增加登录成功、密码错误、账号不存在、服务器满员等编码测试。 |
| 16 | 登录校验与安全 | 客户端版本校验 | `CSPJoinIdPassRequest` 比对 `szClientVersion` | 当前只解析 `Version` | 未覆盖 | 增加版本配置和 `Result=6` 的拒绝流程。 |
| 17 | 登录校验与安全 | 客户端 serial 校验 | `CSPJoinIdPassRequest` 比对 `szGameServerExeSerial` | 当前只解析 `Serial` | 未覆盖 | 增加 serial 配置和非法客户端断线策略。 |
| 18 | 登录校验与安全 | HWID 空值校验 | `CSPJoinIdPassRequest` 检查 HWID | 当前未校验 | 未覆盖 | HWID 为空应记录日志并拒绝或降级处理。 |
| 19 | 登录校验与安全 | ConnectMember 限制 | `ConMember.IsMember` | `conf.ConnectMember` 已加载但 Login 未使用 | 未覆盖 | `EnableConnectMember` 开启时限制非白名单账号登录。 |
| 20 | 登录校验与安全 | 测试服白名单 | `IsTestServer` + `ConMember.IsMember` | 当前无测试服登录限制 | 未覆盖 | 若保留测试服概念，需要配置化实现。 |
| 21 | 登录校验与安全 | 登录超时校验 | `PacketCheckTime` | 当前无登录包超时 | 未覆盖 | 连接后长时间不登录应关闭连接。 |
| 22 | 登录校验与安全 | 重复登录包校验 | `LoginMsgSnd/LoginMsgCount` | 当前无登录包发送状态 | 未覆盖 | 防止同一连接重复发送登录请求造成状态覆盖。 |
| 23 | 登录校验与安全 | 账号重复在线校验 | `gObjSetAccountLogin` 扫描在线账号、JoinServer 中心态 | 当前无在线账号唯一性检查 | 未覆盖 | 单服先查本地对象池，多服通过 `28-external-comm.md` 查询或通知中心在线态。 |
| 24 | 登录校验与安全 | 登录失败次数统计 | `CLoginCount::Add/Delete/Get` | 当前无失败计数 | 未覆盖 | 增加账号/IP/连接维度的失败计数和断线策略。 |
| 25 | 登录校验与安全 | 机器码连接限制 | GameServer HWID/机器码限制语义 | 当前只解析 HWID | 未覆盖 | 实现同 HWID 最大连接数、封禁、审计日志。 |
| 26 | 账号状态绑定 | 账号登录绑定 | `gObjSetAccountLogin` | `Player.Login` 设置 `accountID/accountName` | 已覆盖 | 保证账号绑定只在密码校验成功后发生。 |
| 27 | 账号状态绑定 | UserNumber/DBNumber 语义 | `gObjSetAccountLogin` 设置 `UserNumber/DBNumber` | `accountID` 兼作 DB 主键 | 设计替代 | 单服可用本地 DB 主键，多服或 DataServer 模式通过 `28-external-comm.md` 返回编号。 |
| 28 | 账号状态绑定 | 账号名一致性检查 | `gObjIsAccontConnect` | 当前角色请求依赖 `p.accountID` | 部分覆盖 | 角色 DB 查询已按账号 ID 过滤，但协议层仍应校验已登录。 |
| 29 | 账号状态绑定 | 密码缓存 | `m_PlayerData->Password` | `p.accountPassword` | 已覆盖 | 仅用于删除角色等二次校验时使用，避免日志泄露。 |
| 30 | 账号状态绑定 | 仓库扩展绑定 | `m_WarehouseExpansion` | `p.warehouseExpansion` | 已覆盖 | 角色列表和仓库逻辑都应使用同一字段。 |
| 31 | 账号状态绑定 | 登录态切换 | `Connected = PLAYER_LOGGED` | `p.ConnectState = ConnectStateLogged` | 已覆盖 | 切换后应允许角色列表/建删角色，拒绝移动/攻击。 |
| 32 | 账号状态绑定 | 游戏态切换 | `Connected = PLAYER_PLAYING` | `LoadCharacter` 设置 `ConnectStatePlaying` | 已覆盖 | 加载角色失败不能进入 Playing。 |
| 33 | 账号状态绑定 | 账号解绑 | `gObjDel` 清空 `AccountID/Password/HWID` | `ObjectManager.DeletePlayer` + `Player.Offline` | 部分覆盖 | 下线时应明确清空账号字段和角色字段，避免对象复用污染。 |
| 34 | 账号状态绑定 | 连接与角色名索引 | `gObjGetIndex`、`gObjUserIdConnectCheck` | `ObjectManager.GetPlayerByName` | 部分覆盖 | 加载角色时要防止同名角色重复在线。 |
| 35 | 账号状态绑定 | 账号状态测试 | GameServer 多状态防护 | 当前测试覆盖不足 | 未覆盖 | 覆盖未登录请求列表、已登录建角、Playing 删除、断线清理。 |
| 36 | 账号模型与管理 API | 账号 DB 模型 | DataServer 账号表语义 | `model.Account` | 已覆盖 | 现有字段可支撑基础登录和仓库。 |
| 37 | 账号模型与管理 API | 账号唯一名 | DataServer AccountId 唯一 | `Account.Name gorm:"unique"` | 已覆盖 | 需要处理大小写敏感策略，避免同名大小写绕过。 |
| 38 | 账号模型与管理 API | 密码字段 | DataServer 密码校验 | `Account.Password` | 部分覆盖 | 当前明文存储，不适合长期设计；需定义 hash/兼容策略。 |
| 39 | 账号模型与管理 API | 邮箱字段 | GameServer 原业务无 HTTP 邮箱管理 | `Account.UserEmail` | 设计替代 | HTTP 管理侧字段可保留，但不应影响游戏登录协议。 |
| 40 | 账号模型与管理 API | 仓库字段 | DataServer 仓库表 | `Account.Warehouse` JSONB | 部分覆盖 | 仓库深层行为归道具/仓库模块，本处记录账号归属。 |
| 41 | 账号模型与管理 API | 创建账号 HTTP 入口 | GameServer 无 HTTP 对照 | `handle/http.go::CreateAccount`、`POST /api/accounts` | 已覆盖 | HTTP 管理入口和鉴权审计归 `29-ops.md`，账号系统负责账号创建业务规则。 |
| 42 | 账号模型与管理 API | 创建账号 Command | 后台/GM 管理语义 | `game/cmd/cmd.go::CreateAccount` | 已覆盖 | Command 层将 HTTP 或其他管理入口转发到 DB，后续可统一鉴权和审计。 |
| 43 | 账号模型与管理 API | 创建账号 DB | DataServer 账号创建语义 | `game/model/db.go::CreateAccount` | 已覆盖 | 当前使用 `FirstOrCreate`，需明确重复账号错误映射和事务边界。 |
| 44 | 账号模型与管理 API | 创建账号参数校验 | GameServer 管理侧另行实现 | `validator`、`CreateAccountBind/CreateAccountValidate` | 部分覆盖 | 补充账号名、邮箱、密码策略、重复账号、HTTP 状态码和错误码测试。 |
| 45 | 账号模型与管理 API | 查询账号 DB | DataServer 账号查询 | `DB.GetAccountByName/GetAccountByID` | 已覆盖 | 登录和仓库入口已使用。 |
| 46 | 账号模型与管理 API | 获取账号列表 HTTP 入口 | GameServer 无 HTTP 对照 | `handle/http.go::GetAccountList`、`GET /api/accounts` | 已覆盖 | HTTP 管理入口和鉴权审计归 `29-ops.md`，账号系统负责查询语义。 |
| 47 | 账号模型与管理 API | 获取账号列表 Command | 后台/GM 管理语义 | `game/cmd/cmd.go::GetAccountList` | 已覆盖 | Command 层承接 HTTP 查询条件并调用 DB 查询。 |
| 48 | 账号模型与管理 API | 获取账号列表 DB | DataServer 账号查询语义 | `game/model/db.go::GetAccountList` | 已覆盖 | 当前按 `UserEmail` 过滤；后续可扩展账号名、ID、分页和排序。 |
| 49 | 账号模型与管理 API | 获取账号列表返回模型 | GameServer 管理侧另行实现 | `model.Account` JSON 输出 | 部分覆盖 | 明确是否暴露密码、仓库、角色关联等敏感或大字段。 |
| 50 | 账号模型与管理 API | 删除账号 HTTP 入口 | GameServer 无 HTTP 对照 | `handle/http.go::DeleteAccount`、`DELETE /api/accounts/:id` | 已覆盖 | HTTP 管理入口和鉴权审计归 `29-ops.md`，账号系统负责删除约束。 |
| 51 | 账号模型与管理 API | 删除账号 Command | 后台/GM 管理语义 | `game/cmd/cmd.go::DeleteAccount` | 已覆盖 | Command 层承接删除请求并调用 DB 删除。 |
| 52 | 账号模型与管理 API | 删除账号 DB | DataServer 账号删除语义 | `game/model/db.go::DeleteAccount` | 部分覆盖 | 删除账号应通过本地在线索引和 `28-external-comm.md` 中心在线态确认是否允许。 |
| 53 | 账号模型与管理 API | 删除账号级联策略 | GameServer 运营后台语义 | 当前未明确 | 未覆盖 | 明确是否级联删除角色、仓库、好友、邮件、战盟关系和日志数据。 |
