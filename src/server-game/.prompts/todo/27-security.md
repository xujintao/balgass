# 27. 安全风控系统

本模块覆盖封包与协议安全、CheckSum/CRC、AntiHack 心跳、攻击速度检测、多段攻击检测、技能距离检测、移动异常检测、聊天/行为限流、Penalty/Ban/Kick 和安全审计。安全风控系统不拥有战斗伤害、技能效果、地图规则、聊天业务、交易业务或账号登录业务；这些业务模块在入口处调用安全系统，安全系统返回通过、拒绝、记录、限制、踢线或封禁等结果。人工封禁/解封、后台处罚入口和审计查看归 `29-ops.md`；跨服踢线、中心封禁、禁重连通知等外部服务调用归 `28-external-comm.md`。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 模块边界与总入口 | SecurityManager 总入口 | `GameSecurity`、分散全局安全对象 | 暂无 `game/security` | 未覆盖 | 建立统一安全风控入口，承接封包、移动、攻击、聊天和处罚检查。 |
| 2 | 模块边界与总入口 | 安全系统初始化 | `GameMainInit` 初始化安全配置 | `conf` 已加载配置 | 部分覆盖 | 启动时加载策略、阈值、开关、审计 logger 和检查器实例。 |
| 3 | 模块边界与总入口 | 玩家安全状态挂载 | `OBJECTSTRUCT` 安全字段 | `Object`/`Player` 缺少集中状态 | 未覆盖 | 为每个玩家保存 checksum、CRC、速度、距离、聊天、处罚等运行态。 |
| 4 | 模块边界与总入口 | 安全检查结果模型 | 多处分支 return/CloseClient | 暂无 | 未覆盖 | 统一返回 Allow、Deny、Warn、Kick、Ban、DisableReconnect 等结果。 |
| 5 | 模块边界与总入口 | 业务调用接口 | `gCheckSkillDistance` 等函数 | 业务直接处理 | 未覆盖 | 提供给对象、技能、地图、交易、聊天等模块调用的稳定接口。 |
| 6 | 模块边界与总入口 | 安全上下文 | `aIndex`、account、name、IP 日志 | 分散字段 | 未覆盖 | 每次检查携带玩家 ID、账号、角色名、IP、opcode、地图、坐标等上下文。 |
| 7 | 模块边界与总入口 | 安全配置快照 | `g_ConfigRead` 安全项 | `conf.CommonServer.GameServerInfo` | 部分覆盖 | 把安全相关配置聚合成只读快照，避免业务到处读全局配置。 |
| 8 | 模块边界与总入口 | 检查器注册表 | 分散全局对象 | 暂无 | 未覆盖 | 注册封包、CRC、速度、移动、技能距离、多段攻击等检查器。 |
| 9 | 模块边界与总入口 | 检查开关统一判断 | 多处 `if enable` | 配置存在但未统一 | 部分覆盖 | 所有检查先走统一开关，避免业务层重复判断。 |
| 10 | 模块边界与总入口 | 安全系统 Tick | `gObjSecondProc` 周期检测 | `Process1000ms` 可接入 | 未覆盖 | 周期检查 AntiHack、CheckSum、CRC、聊天窗口、处罚过期等。 |
| 11 | 模块边界与总入口 | 连接关闭接入 | `CloseClient`、`GCCloseClientSend` | `DeletePlayer`/连接关闭 | 部分覆盖 | 安全系统只输出关闭原因，运行时/对象系统执行关闭。 |
| 12 | 模块边界与总入口 | 禁止重连接入 | `GCSendDisableReconnect` | 暂无 | 未覆盖 | 严重安全命中时下发禁重连或写入账号状态。 |
| 13 | 模块边界与总入口 | 安全事件枚举 | 日志字符串分散 | 暂无 | 未覆盖 | 定义 packet_checksum、crc_mismatch、speed_hack、move_hack 等事件类型。 |
| 14 | 模块边界与总入口 | 安全严重级别 | 日志和断线分支 | 暂无 | 未覆盖 | 区分 info、warn、deny、kick、ban，便于配置化处理。 |
| 15 | 模块边界与总入口 | 安全测试夹具 | 无统一 | 暂无 | 未覆盖 | 为移动、速度、checksum、限频等检查提供可重复单测夹具。 |
| 16 | 封包与协议安全 | 加密标记校验 | GameServer 加密协议检查 | `api.enc && !req.Encrypt` 只 warn | 需修正 | 运行时应调用安全策略决定拒绝、记录或踢线。 |
| 17 | 封包与协议安全 | opcode 白名单 | `ProtocolCore` switch | `handle/c1c2.go` apiMap | 部分覆盖 | 未注册 opcode 应记录安全事件，并按阈值处理。 |
| 18 | 封包与协议安全 | C1/C2 长度校验 | 协议头长度检查 | `c1c2` 解析已有基础 | 部分覆盖 | 校验包长与协议体结构一致，异常包进入安全审计。 |
| 19 | 封包与协议安全 | C3/C4 加密包边界 | `PacketEncrypt` | 当前未完整接入 | 未覆盖 | 区分明文和密文协议，错误类型必须可审计。 |
| 20 | 封包与协议安全 | PacketEncrypt 加密 | `CPacketEncrypt::Encrypt` | 暂无 | 未覆盖 | 如需要兼容客户端加密包，封装加密输出。 |
| 21 | 封包与协议安全 | PacketEncrypt 解密 | `CPacketEncrypt::Decrypt` | 暂无 | 未覆盖 | 解密失败应计入协议安全事件。 |
| 22 | 封包与协议安全 | Rijndael/RSA 兼容 | CryptoPP Rijndael/RSA | 暂无 | 未覆盖 | 明确是否完全兼容 GameServer 加密协议。 |
| 23 | 封包与协议安全 | 包频率限制 | Packet time check | `PacketLimit`、`PacketTimeMin` 配置 | 未覆盖 | 按玩家和 opcode 统计收包间隔，超过阈值拒绝或踢线。 |
| 24 | 封包与协议安全 | 包速率滑动窗口 | GameServer packet time 字段 | 暂无 | 未覆盖 | 用滑动窗口替代单点时间差，减少网络抖动误判。 |
| 25 | 封包与协议安全 | 未登录协议限制 | 协议状态检查 | handle auth 注释 | 需修正 | 未登录阶段只允许登录/握手/安全协议。 |
| 26 | 封包与协议安全 | 角色未进入游戏限制 | Object 状态检查 | 部分业务未校验 | 未覆盖 | 进入游戏前禁止移动、攻击、交易、商店等业务包。 |
| 27 | 封包与协议安全 | 协议阶段状态机 | `Connected/Logged/Playing` 语义 | 当前分散 | 未覆盖 | 安全系统提供阶段校验，业务不重复写状态判断。 |
| 28 | 封包与协议安全 | 协议参数范围审计 | 多处 bound check | 分散 | 未覆盖 | 负数索引、越界坐标、非法 slot、非法目标统一记录。 |
| 29 | 封包与协议安全 | 重放包检测 | 序列/时间检查分散 | 暂无 | 未覆盖 | 对关键协议记录近期摘要，发现重复确认、重复交易等异常。 |
| 30 | 封包与协议安全 | 协议异常计数 | HackUserKickCount | 配置存在 | 未覆盖 | 按安全事件累计，达到阈值触发踢线或封禁。 |
| 31 | CheckSum / CRC | PacketCheckSum 初始化 | `CPacketCheckSum::Init` | 暂无 | 未覆盖 | 初始化 checksum 表、学习状态、清理函数。 |
| 32 | CheckSum / CRC | PacketCheckSum Check | `CPacketCheckSum::Check` | `0x72 packetCheckSum` 路由 | 未覆盖 | 定期检查客户端上报的函数 checksum 是否完整有效。 |
| 33 | CheckSum / CRC | PacketCheckSum Add | `CPacketCheckSum::Add` | 暂无 | 未覆盖 | 接收客户端 checksum 上报并与期望值比较。 |
| 34 | CheckSum / CRC | PacketCheckSum Clear | `ClearCheckSum` | 暂无 | 未覆盖 | 周期清理玩家 checksum 缓存，避免旧数据误用。 |
| 35 | CheckSum / CRC | CheckSum 函数索引上限 | `MAX_PACKET_CHECKSUM_FUNCTION_INDEX=22` | 暂无 | 未覆盖 | 固定客户端函数索引范围，越界视为异常。 |
| 36 | CheckSum / CRC | CheckSum 平均表 | `MAX_CHECKSUM_TABLE_AVG=100` | 暂无 | 未覆盖 | 支持 GameServer 的学习/平均 checksum 表语义。 |
| 37 | CheckSum / CRC | CheckSum 学习锁定 | `SetClearChecksumFunc` | 暂无 | 未覆盖 | 明确是否迁移学习模式；线上应使用固定表。 |
| 38 | CheckSum / CRC | CheckSum 超时检测 | `CheckSumTime` 超时 | 配置存在 | 未覆盖 | 超时未上报时记录并按配置关闭连接。 |
| 39 | CheckSum / CRC | CheckSum 失败踢线 | `gDisconnectHackUser` | `HackUserKickEnable` | 未覆盖 | 对 checksum mismatch 使用统一处罚策略。 |
| 40 | CheckSum / CRC | MultiCheckSum 文件加载 | `CMultiCheckSum::LoadFile` | 暂无 | 未覆盖 | 加载多文件 CRC/checksum 配置。 |
| 41 | CheckSum / CRC | MultiCheckSum 路径设置 | `SetFilePath` | 暂无 | 未覆盖 | 支持按环境配置 checksum 文件路径。 |
| 42 | CheckSum / CRC | MultiCheckSum 比较 | `CompareCheckSum` | 暂无 | 未覆盖 | 接收客户端文件校验并比较。 |
| 43 | CheckSum / CRC | MultiCheckSum 数量上限 | `MAX_MULTICHECKSUM=10` | 暂无 | 未覆盖 | 限制单次多文件校验数量。 |
| 44 | CheckSum / CRC | FileCRC 协议入口 | `fileCRC` 对应协议 | `0xFA0D fileCRC` 路由 | 未覆盖 | 实现客户端文件 CRC 上报处理。 |
| 45 | CheckSum / CRC | DLL 版本匹配 | Kick unmatched DLL version | `EnableKickUnmatchedDLLVersion` | 未覆盖 | DLL 版本不匹配时按配置踢线。 |
| 46 | CheckSum / CRC | CRC 超时检测 | `CrcCheckTime` 超时 | 配置存在 | 未覆盖 | 超时未上报时触发审计和处罚。 |
| 47 | CheckSum / CRC | CRC 失败计数 | 多处日志/断线 | 暂无 | 未覆盖 | CRC 异常累计到玩家安全状态。 |
| 48 | CheckSum / CRC | 校验响应包下发 | GameServer 请求客户端校验 | 暂无 | 未覆盖 | 定期请求客户端上报 checksum/CRC。 |
| 49 | CheckSum / CRC | 校验结果日志 | `g_Log.Add` | `slog` 未统一 | 未覆盖 | 记录账号、角色、IP、索引、校验类型、期望值、实际值。 |
| 50 | CheckSum / CRC | 校验配置降级 | check disable 配置 | `PacketHackCheckDisable` | 部分覆盖 | 支持开发环境关闭 packet hack check。 |
| 51 | AntiHack 心跳 | AntiCheat 协议入口 | AntiCheat 消息 | `0xFA08 antiCheat` 路由 | 未覆盖 | 接收客户端安全心跳或状态包。 |
| 52 | AntiHack 心跳 | AntiHackBreach 入口 | AntiHack breach 消息 | `0xFA0A antiHackBreach` 路由 | 未覆盖 | 客户端报告反外挂 breach 时记录并处理。 |
| 53 | AntiHack 心跳 | AntiHackCheck 入口 | AntiHack check 消息 | `0xFA11 antiHackCheck` 路由 | 未覆盖 | 处理周期性 anti-hack check 响应。 |
| 54 | AntiHack 心跳 | Hack 消息处理 | Hack detect message | `Player.Hack` 空实现 | 未覆盖 | 实现 `Player.Hack` 或迁移到安全系统入口。 |
| 55 | AntiHack 心跳 | 心跳最后响应时间 | `AntiHackCheckTime` | 暂无 | 未覆盖 | 保存最后一次 anti-hack 响应 tick。 |
| 56 | AntiHack 心跳 | 心跳超时 180s | user.cpp 周期检查 | 暂无 | 未覆盖 | 超时未响应时按配置断开或禁重连。 |
| 57 | AntiHack 心跳 | AntiRef 检查 | `EnableAntiRefCheck` | 配置存在 | 未覆盖 | 迁移反反射/客户端保护相关检查配置。 |
| 58 | AntiHack 心跳 | RecvHook 保护 | `EnableRecvHookProtection` | 配置存在 | 未覆盖 | 接入客户端 hook 检测结果。 |
| 59 | AntiHack 心跳 | HackDetectMessage | `HackDetectMessage` | 配置存在 | 未覆盖 | 命中时是否下发提示由配置决定。 |
| 60 | AntiHack 心跳 | AutoBanHackUser | `EnableAutoBanHackUser` | 配置存在 | 未覆盖 | 严重 breach 时自动封禁账号或写入处罚状态。 |
| 61 | AntiHack 心跳 | KickAntiHackBreach | `EnableKickAntiHackBreach` | 配置存在 | 未覆盖 | 客户端 breach 时是否立即踢线。 |
| 62 | AntiHack 心跳 | AntiHack 请求调度 | GameServer 周期请求 | 暂无 | 未覆盖 | 由安全 Tick 定时向客户端发起检查。 |
| 63 | AntiHack 心跳 | 心跳异常审计 | `LogAdd` | 暂无统一安全日志 | 未覆盖 | 记录心跳丢失、breach、伪造响应、超时。 |
| 64 | AntiHack 心跳 | 心跳误判保护 | GameServer 阈值 | 暂无 | 未覆盖 | 网络抖动下不应一次超时直接封禁，踢线和封禁分层。 |
| 65 | 攻击速度检测 | SpeedHackCheck 总入口 | `CObjUseSkill::SpeedHackCheck` | 暂无 | 未覆盖 | 技能/普攻入口调用速度检测。 |
| 66 | 攻击速度检测 | LastAttackTime 记录 | `m_LastAttackTime` | 暂无集中字段 | 未覆盖 | 每次攻击记录服务端 tick。 |
| 67 | 攻击速度检测 | DetectSpeedHackTime | `m_DetectSpeedHackTime` | 暂无 | 未覆盖 | 按窗口统计攻击间隔。 |
| 68 | 攻击速度检测 | DetectCount 统计 | `m_DetectCount` | 暂无 | 未覆盖 | 统计窗口内异常次数。 |
| 69 | 攻击速度检测 | SumLastAttackTime | `m_SumLastAttackTime` | 暂无 | 未覆盖 | 计算平均攻击间隔。 |
| 70 | 攻击速度检测 | AttackSpeedHackDetectedCount | `m_AttackSpeedHackDetectedCount` | 暂无 | 未覆盖 | 累计攻击速度异常次数。 |
| 71 | 攻击速度检测 | DetectedHackKickCount | `m_DetectedHackKickCount` | `SpeedHackKickCount` 配置 | 未覆盖 | 达到阈值后执行踢线。 |
| 72 | 攻击速度检测 | SpeedHackPenalty | `m_SpeedHackPenalty` | `SpeedHackPenaltyEnable` | 未覆盖 | 配置开启时先限制攻击而不是立即踢线。 |
| 73 | 攻击速度检测 | 物理攻击速度阈值 | 攻速公式参与 | 公式系统待接入 | 未覆盖 | 按职业、装备、buff 计算允许最小攻击间隔。 |
| 74 | 攻击速度检测 | 魔法攻击速度阈值 | 魔攻速参与 | 公式系统待接入 | 未覆盖 | 技能类型不同使用不同速度阈值。 |
| 75 | 攻击速度检测 | PotionDelayTime | 药水延迟配置 | `PotionDelayTime` | 未覆盖 | 药水使用频率也纳入行为限速。 |
| 76 | 攻击速度检测 | SpeedHackCheckEnable | `SpeedHackCheckEnable` | 配置存在 | 部分覆盖 | 配置关闭时跳过速度检查但仍可记录统计。 |
| 77 | 攻击速度检测 | 攻击速度日志 | DebugInfo ASData | `GameSecurity::m_ASData` | 未覆盖 | 记录攻速计算值、实际间隔、技能号、目标。 |
| 78 | 攻击速度检测 | 速度检测误差容忍 | SpeedHack temp/threshold | 暂无 | 未覆盖 | 保留网络延迟和服务器 tick 抖动容忍。 |
| 79 | 攻击速度检测 | 速度异常返回 | `return false` | 暂无 | 未覆盖 | 命中后技能/攻击入口应中止业务执行。 |
| 80 | 攻击速度检测 | 速度处罚恢复 | Penalty 时间字段 | 暂无 | 未覆盖 | 临时处罚到期后恢复攻击能力。 |
| 81 | 多段攻击检测 | MultiAttackHackCheck 总入口 | `CMultiAttackHackCheck` | 暂无 | 未覆盖 | 范围技能、多段技能、穿透技能调用多段攻击检测。 |
| 82 | 多段攻击检测 | Insert 记录 | `CMultiAttackHackCheck::Insert` | 暂无 | 未覆盖 | 记录技能、目标、序列、时间。 |
| 83 | 多段攻击检测 | PenetrationSkill 检查 | `CheckPenetrationSkill` | 暂无 | 未覆盖 | 限制穿透技能单次合法命中数量。 |
| 84 | 多段攻击检测 | FireScreamSkill 检查 | `CheckFireScreamSkill` | 暂无 | 未覆盖 | 限制 FireScream 类技能异常多目标。 |
| 85 | 多段攻击检测 | Serial 计数 | serial counters | 暂无 | 未覆盖 | 使用客户端/服务端序列防止重复命中。 |
| 86 | 多段攻击检测 | 多目标数量限制 | GameServer 多目标限制 | 暂无 | 未覆盖 | 单次技能命中目标数不能超过技能定义和视野范围。 |
| 87 | 多段攻击检测 | 同目标重复命中限制 | GameServer serial 语义 | 暂无 | 未覆盖 | 同一技能窗口内不能对同目标重复结算异常次数。 |
| 88 | 多段攻击检测 | 范围技能时间窗 | MultiAttack 时间窗口 | 暂无 | 未覆盖 | 多段命中必须落在允许窗口内。 |
| 89 | 多段攻击检测 | 多段异常处罚 | GameServer log/return false | 暂无 | 未覆盖 | 异常时中止技能结算并累计安全事件。 |
| 90 | 多段攻击检测 | 技能白名单 | 特殊技能分支 | 暂无 | 未覆盖 | 特殊召唤、持续伤害、陷阱类技能需排除或单独规则。 |
| 91 | 技能距离检测 | SkillDistance 总入口 | `gCheckSkillDistance` | 暂无 | 未覆盖 | 技能系统调用安全距离检查。 |
| 92 | 技能距离检测 | SkillDistanceCheckEnable | `g_iSkillDistanceCheck` | 配置存在 | 部分覆盖 | 配置关闭时不阻断但可保留日志。 |
| 93 | 技能距离检测 | 技能距离表读取 | `MagicDamageC.GetSkillDistance` | 技能定义待完善 | 未覆盖 | 从技能表获取基础距离。 |
| 94 | 技能距离检测 | DistanceTemp 容忍 | `g_iSkillDistanceCheckTemp` | 配置存在 | 部分覆盖 | 加入距离容忍值，避免坐标同步误差。 |
| 95 | 技能距离检测 | dx/dy 比较 | `abs(dx/dy)` | 暂无 | 未覆盖 | 对照 GameServer 使用格子差而非欧式距离。 |
| 96 | 技能距离检测 | 技能例外 40 | `skill 40` exception | 暂无 | 未覆盖 | 迁移 GameServer 特殊技能例外。 |
| 97 | 技能距离检测 | 技能例外 392 | `skill 392` exception | 暂无 | 未覆盖 | 迁移 GameServer 特殊技能例外。 |
| 98 | 技能距离检测 | 目标坐标来源 | target X/Y/TX/TY | `Object` 坐标字段 | 部分覆盖 | 明确移动中对象使用当前位置还是目标点。 |
| 99 | 技能距离检测 | 距离错误计数 | `m_iSkillDistanceErrorCount` | 暂无 | 未覆盖 | 连续距离异常才踢线，降低误判。 |
| 100 | 技能距离检测 | 距离错误 tick | `m_dwSkillDistanceErrorTick` | 暂无 | 未覆盖 | 错误计数按时间窗口衰减。 |
| 101 | 技能距离检测 | SkillDistanceKick | `SkillDistanceKickEnable` | 配置存在 | 部分覆盖 | 达到阈值时按配置踢线。 |
| 102 | 技能距离检测 | SkillDistanceKickCount | `SkillDistanceKickCount` | 配置存在 | 部分覆盖 | 配置距离异常次数阈值。 |
| 103 | 技能距离检测 | SkillDistanceKickCheckTime | `SkillDistanceKickCheckTime` | 配置存在 | 部分覆盖 | 配置距离异常统计窗口。 |
| 104 | 技能距离检测 | 距离异常日志 | `LogAdd` invalid distance | 暂无 | 未覆盖 | 记录技能、距离、双方坐标、地图、容忍值。 |
| 105 | 技能距离检测 | 距离检查返回策略 | `return false` | 暂无 | 未覆盖 | 命中后技能系统不得继续结算。 |
| 106 | 移动检测 | MoveCheck 总入口 | `CMoveCheck` | 暂无 | 未覆盖 | 对象移动入口调用移动安全检查。 |
| 107 | 移动检测 | 最近坐标环 | `CMoveCheck` 保存 5 个坐标 | 暂无 | 未覆盖 | 保存最近移动点，用于检测瞬移和异常回退。 |
| 108 | 移动检测 | MoveCheck Insert | `CMoveCheck::Insert` | 暂无 | 未覆盖 | 每次合法移动后写入轨迹。 |
| 109 | 移动检测 | MoveCheck Check | `CMoveCheck::Check` | 暂无 | 未覆盖 | 客户端上报路径前先校验起点和轨迹。 |
| 110 | 移动检测 | 起点漂移检查 | GameServer 校验客户端坐标 | `Object.Move` 信任起点 | 需修正 | 客户端起点与服务端位置差距过大时拒绝。 |
| 111 | 移动检测 | 目标越界检查 | `GetAttr` 越界阻挡 | `_map.valid` | 部分覆盖 | 越界移动记安全事件。 |
| 112 | 移动检测 | 阻挡路径检查 | MapClass wall/attr | `checkNoWall`、`processMove` | 部分覆盖 | 穿墙、越障和不可站点移动进入安全审计。 |
| 113 | 移动检测 | 移动速度检查 | Tick/路径耗时 | `processMove` 400ms | 部分覆盖 | 检测路径步进过快、频繁改向、瞬移。 |
| 114 | 移动检测 | 斜向耗时检查 | 斜向 1.3 倍 | `pathDir%2==0` | 部分覆盖 | 移动安全检查应复用地图移动耗时规则。 |
| 115 | 移动检测 | 地图切换期间移动限制 | `m_bMapSvrMoveQuit` | 暂无 | 未覆盖 | 跨图/跨服移动过程中拒绝普通移动包。 |
| 116 | 移动检测 | 死亡移动限制 | Live/State check | 分散 | 未覆盖 | 死亡状态移动包记录异常并拒绝。 |
| 117 | 移动检测 | 交易移动限制 | GameServer interface check | `19-trade.md` 待实现 | 未覆盖 | 交易中移动由交易系统处理业务取消，安全系统记录异常。 |
| 118 | 移动检测 | 个人商店移动限制 | PShop checks | `20-personal-shops.md` 待实现 | 未覆盖 | 开店/离线交易场景移动规则需审计。 |
| 119 | 移动检测 | GM 强制移动例外 | GM move/teleport | 暂无 | 未覆盖 | 服务端主动传送不应被判定为作弊。 |
| 120 | 移动检测 | 移动异常处罚 | MoveCheck return | 暂无 | 未覆盖 | 多次异常才踢线，单次异常可纠偏。 |
| 121 | 安全区与场景限制 | 安全区攻击限制 | safe zone attack block | `EnableBlockAttackInSafeZone` | 未覆盖 | 攻击入口调用安全区检查，命中后拒绝攻击。 |
| 122 | 安全区与场景限制 | 安全区技能限制 | `SkillSafeZoneUse` | 暂无 | 未覆盖 | 区分 Buff、攻击、召唤、传送技能。 |
| 123 | 安全区与场景限制 | 安全区击退限制 | 技能击退分支 | `Knockback` 部分检查 | 未覆盖 | 安全区内外击退需要明确策略。 |
| 124 | 安全区与场景限制 | 地图 PVP 限制审计 | `PvPConfig` | 地图属性未完整保存 | 未覆盖 | 地图系统给出规则，安全系统记录越权攻击。 |
| 125 | 安全区与场景限制 | 活动地图限制审计 | 事件地图分支 | 活动系统待实现 | 未覆盖 | 非法进入/非法攻击/非法交易等进入审计。 |
| 126 | 安全区与场景限制 | 传送状态限制 | teleport state | `Teleport` 待完善 | 未覆盖 | 传送中攻击、交易、移动包均应被拒绝。 |
| 127 | 安全区与场景限制 | 界面状态限制 | `m_IfState` | `InterfaceState` 待完善 | 未覆盖 | 仓库、合成、商店等互斥状态下的异常行为记录。 |
| 128 | 安全区与场景限制 | PK 惩罚限制接口 | Penalty/PK state | `17-gens.md`/对象系统 | 未覆盖 | 业务提供 PK 状态，安全系统只执行限制和审计。 |
| 129 | 安全区与场景限制 | 召唤物安全继承 | summon owner checks | `24-pets-summons.md` 待实现 | 未覆盖 | 召唤物攻击限制继承主人安全区和地图限制。 |
| 130 | 安全区与场景限制 | 异常场景审计 | 分散日志 | 暂无 | 未覆盖 | 记录场景、地图、坐标、接口状态、业务入口。 |
| 131 | 聊天与行为限制 | ChatLimitTime | `ChatLimitTime` | 暂无 | 未覆盖 | 聊天入口按时间限制发言频率。 |
| 132 | 聊天与行为限制 | ChatLimitTimeSec | `ChatLimitTimeSec` | 暂无 | 未覆盖 | 配置聊天限制窗口长度。 |
| 133 | 聊天与行为限制 | ChatFloodTime | `m_ChatFloodTime` | 暂无 | 未覆盖 | 记录刷屏窗口起始时间。 |
| 134 | 聊天与行为限制 | ChatFloodCount | `m_ChatFloodCount` | 暂无 | 未覆盖 | 统计窗口内聊天次数。 |
| 135 | 聊天与行为限制 | 聊天刷屏处罚 | chat flood penalty | 聊天系统未独立 | 未覆盖 | 达到阈值时禁言、拒绝消息或踢线。 |
| 136 | 聊天与行为限制 | 聊天内容过滤接口 | GameServer 过滤分散 | 暂无 | 未覆盖 | 安全系统提供过滤/审计接口，聊天系统决定频道业务。 |
| 137 | 聊天与行为限制 | 重复消息检测 | flood 语义 | 暂无 | 未覆盖 | 短时间重复文本计入刷屏。 |
| 138 | 聊天与行为限制 | 行为冷却通用接口 | potion/trade/chat delay | 暂无 | 未覆盖 | 为药水、聊天、交易请求等提供通用 cooldown。 |
| 139 | 聊天与行为限制 | 请求骚扰限制 | trade/party/guild invite flood | 暂无 | 未覆盖 | 高频邀请、交易、组队请求进入风控。 |
| 140 | 聊天与行为限制 | 行为限流日志 | 分散日志 | 暂无 | 未覆盖 | 记录行为类型、窗口次数、阈值和处理结果。 |
| 141 | Penalty / Ban / Kick | Penalty 字段模型 | `Penalty` | 暂无 | 未覆盖 | 建立临时处罚类型，如禁言、禁交易、禁攻击、禁移动。 |
| 142 | Penalty / Ban / Kick | PenaltyMask | `PenaltyMask` | 暂无 | 未覆盖 | 支持多个处罚并存和位掩码查询。 |
| 143 | Penalty / Ban / Kick | 处罚过期时间 | GameServer penalty 时间 | 暂无 | 未覆盖 | 临时处罚需要过期自动恢复。 |
| 144 | Penalty / Ban / Kick | 禁言处罚 | GMMng/user checks | 暂无 | 未覆盖 | 聊天系统查询是否禁言。 |
| 145 | Penalty / Ban / Kick | 禁交易处罚 | PersonalStore/user checks | 暂无 | 未覆盖 | 交易、个人商店、NPC 商店入口查询是否被限制。 |
| 146 | Penalty / Ban / Kick | 禁攻击处罚 | SpeedHack penalty | 暂无 | 未覆盖 | 攻击/技能入口查询是否被限制。 |
| 147 | Penalty / Ban / Kick | Kick 策略 | `CloseClient` | 连接关闭流程 | 部分覆盖 | 安全系统给出原因，由运行时执行断开。 |
| 148 | Penalty / Ban / Kick | Ban 策略 | AutoBanHackUser | 暂无 | 未覆盖 | 自动判定归安全系统，人工封禁/解封入口归 `29-ops.md`，跨服/中心通知通过 `28-external-comm.md` 发起。 |
| 149 | Penalty / Ban / Kick | DisableReconnect | `GCSendDisableReconnect` | 暂无 | 未覆盖 | 对严重外挂下发禁重连或写入会话状态。 |
| 150 | Penalty / Ban / Kick | HackUserKickCount | `HackUserKickCount` | 配置存在 | 部分覆盖 | 多次异常后踢线。 |
| 151 | Penalty / Ban / Kick | HackUserKickEnable | `HackUserKickEnable` | 配置存在 | 部分覆盖 | 开关控制是否踢线。 |
| 152 | Penalty / Ban / Kick | 处罚审计 | `LogAdd` | 暂无统一 | 未覆盖 | 每次处罚记录触发规则、证据、处理结果。 |
| 153 | 跨系统接口 | 运行时封包接口 | ProtocolCore/PacketEncrypt | `01-runtime.md` | 未覆盖 | 运行时负责收包发包，安全系统负责校验结果。 |
| 154 | 跨系统接口 | 对象移动接口 | `CGMoveRecv`/MoveCheck | `04-objects.md`、`06-maps.md` | 未覆盖 | 对象系统调用移动安全检查后再更新坐标。 |
| 155 | 跨系统接口 | 对象攻击接口 | `CGAttackRecv`/SpeedHack | `04-objects.md` | 未覆盖 | 攻击入口调用速度、安全区、距离检查。 |
| 156 | 跨系统接口 | 技能接口 | `ObjUseSkill` | `10-skills.md` | 未覆盖 | 技能系统调用距离、多段、速度检查。 |
| 157 | 跨系统接口 | 交易接口 | Trade/PShop penalty checks | `19-trade.md`、`20-personal-shops.md` | 未覆盖 | 交易系统查询处罚状态并上报异常交易。 |
| 158 | 配置与审计日志 | 安全配置聚合 | `CommonServer.cfg` 安全项 | `conf/config.go` | 部分覆盖 | 把已加载配置映射到安全策略结构。 |
| 159 | 配置与审计日志 | 安全审计日志 | `g_Log.Add`/hack logs | `slog` | 未覆盖 | 独立安全日志字段，方便排查和运营处理。 |
| 160 | 配置与审计日志 | 安全回归测试 | GameServer 无统一 | 暂无 | 未覆盖 | 覆盖封包异常、CRC mismatch、速度异常、距离异常、移动异常和处罚恢复。 |
