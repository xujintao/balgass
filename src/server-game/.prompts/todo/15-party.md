# 15. 组队系统

本模块覆盖组队邀请、响应、队伍生命周期、成员列表、队长、退队/踢出、队伍聊天、队友坐标/血条/Buff 同步、组队匹配、事件入场授权，以及经验、掉落、任务、地图、Gens 对组队的联动边界。普通聊天和私聊仍归对象系统入口；组队聊天只记录队伍路由能力。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 现有入口 | 组队协议路由 | `protocol.cpp::CGPartyRequestRecv/CGPartyRequestResultRecv/CGPartyList/CGPartyDelUser` | `handle/c1c2.go` 0x40/0x41/0x42 及相关占位 | 部分覆盖 | Go 侧已有协议表入口，但缺实际业务处理。 |
| 2 | 现有入口 | 队友坐标开关 | `GCPartyMemberPosition` | `Player.StartPartyNumberPosition/StopPartyNumberPosition` 空实现 | 未覆盖 | 实现客户端请求队友坐标同步的开关和定时下发。 |
| 3 | 现有入口 | 组队经验配置 | `gPartyExp*`、`gSetPartyExp*` | `conf.CommonServer.GameServerInfo.Party*ExpBonus` | 部分覆盖 | 配置已存在，需由组队模块提供成员遍历给经验系统。 |
| 4 | 现有入口 | 对象队伍字段 | `OBJECTSTRUCT::PartyNumber/PartyTargetUser` | `object.Object` 注释字段 | 未覆盖 | 恢复或重设 Go 侧对象队伍状态字段。 |
| 5 | 队伍管理器 | PartyClass 初始化 | `CPartyClass::CPartyClass/Init` | 暂无 party manager | 未覆盖 | 新建队伍管理器，负责队伍数组、ID 分配和并发安全。 |
| 6 | 队伍管理器 | 队伍槽清理 | `CPartyClass::Delete` | 暂无 | 未覆盖 | 清空队伍成员、等级、事件授权、队长等状态。 |
| 7 | 队伍管理器 | 是否有效队伍 | `CPartyClass::IsParty` | 暂无 | 未覆盖 | 判断 party number 是否有效且成员数大于 0。 |
| 8 | 队伍管理器 | 获取队伍人数 | `CPartyClass::GetCount` | 暂无 | 未覆盖 | 为经验、掉落、任务、事件提供统一人数查询。 |
| 9 | 队伍管理器 | 获取成员对象编号 | `CPartyClass::GetIndexUser/GetIndex` | 暂无 | 未覆盖 | 按队伍位置或对象编号查询成员。 |
| 10 | 队伍管理器 | 修正成员索引 | `CPartyClass::RevisionIndexUser` | 暂无 | 未覆盖 | 处理断线、重连、对象编号变化后的队伍成员修正。 |
| 11 | 队伍创建 | 创建队伍 | `CPartyClass::Create` | 暂无 | 未覆盖 | 接受邀请成功时创建队伍并设置两名成员。 |
| 12 | 队伍创建 | 设置队伍等级 | `CPartyClass::SetLevel` | 暂无 | 未覆盖 | 维护队伍平均/最高等级信息，供匹配和经验使用。 |
| 13 | 队伍创建 | 获取队伍等级 | `CPartyClass::GetLevel` | 暂无 | 未覆盖 | 为组队匹配、活动入场、经验修正提供队伍等级。 |
| 14 | 队伍创建 | 组队人数上限 | `MAX_USER_IN_PARTY` | 暂无 | 未覆盖 | 固定最大 5 人队伍，并处理扩展配置可能性。 |
| 15 | 邀请流程 | 发起组队邀请 | `CGPartyRequestRecv` | `handle` 有 0x40 | 未覆盖 | 校验目标、距离、状态、等级、地图、Gens 限制后发送邀请。 |
| 16 | 邀请流程 | 邀请目标解析 | `lpMsg->NumberH/NumberL` | 当前无业务解析 | 未覆盖 | 从客户端对象编号解析目标玩家并防止伪造。 |
| 17 | 邀请流程 | 自己邀请自己限制 | `CGPartyRequestRecv` 检查 | 暂无 | 未覆盖 | 禁止向自己发起组队。 |
| 18 | 邀请流程 | 目标在线检查 | `gObjIsConnected` | `ObjectManager.GetObject` | 未覆盖 | 目标不存在或非玩家时返回失败。 |
| 19 | 邀请流程 | 邀请距离检查 | `gObjCalDistance` | `Object.CalcDistance` | 未覆盖 | 队伍邀请通常要求同地图近距离。 |
| 20 | 邀请流程 | 交易/商店/仓库冲突 | `m_IfState` 检查 | InterfaceState 尚未完整使用 | 未覆盖 | 队伍邀请不应打断交易、商店、合成等关键状态。 |
| 21 | 邀请流程 | 已有队伍检查 | `lpObj->PartyNumber` | 暂无 | 未覆盖 | 邀请者和目标已有队伍时按队长/人数规则处理。 |
| 22 | 邀请流程 | 队伍满员检查 | `GetCount >= MAX_USER_IN_PARTY` | 暂无 | 未覆盖 | 满员时返回失败。 |
| 23 | 邀请流程 | 队长权限检查 | `Isleader` | 暂无 | 未覆盖 | 非队长不能邀请新成员。 |
| 24 | 邀请流程 | 设置邀请目标 | `PartyTargetUser` | 注释字段 | 未覆盖 | 记录邀请双方，避免响应伪造或过期邀请。 |
| 25 | 邀请流程 | 邀请包下发 | `GCPartyRequestSend` 相关逻辑 | 当前无响应结构 | 未覆盖 | 向目标玩家发送组队请求。 |
| 26 | 响应流程 | 组队响应入口 | `CGPartyRequestResultRecv` | `handle` 有 0x41 | 未覆盖 | 处理同意、拒绝、超时和目标不一致。 |
| 27 | 响应流程 | 拒绝邀请 | `CGPartyRequestResultRecv` Response 分支 | 暂无 | 未覆盖 | 通知邀请者并清理邀请状态。 |
| 28 | 响应流程 | 同意后创建队伍 | `Create/Add` | 暂无 | 未覆盖 | 双方都无队伍时创建新队伍。 |
| 29 | 响应流程 | 同意后加入队伍 | `Add` | 暂无 | 未覆盖 | 邀请者已有队伍且是队长时加入现有队伍。 |
| 30 | 响应流程 | 响应失败回滚 | `CGPartyRequestResultRecv` 多分支 | 暂无 | 未覆盖 | 任意校验失败要清除 PartyTargetUser 并返回客户端结果。 |
| 31 | 成员维护 | 添加成员 | `CPartyClass::Add` | 暂无 | 未覆盖 | 维护队伍槽、对象 PartyNumber、成员顺序。 |
| 32 | 成员维护 | 删除成员 | `CPartyClass::Delete` | 暂无 | 未覆盖 | 支持按队伍位置和对象编号删除成员。 |
| 33 | 成员维护 | 队伍解散 | `CPartyClass::Destroy` | 暂无 | 未覆盖 | 队伍人数不足或队长退出时按规则解散或转移队长。 |
| 34 | 成员维护 | 队长判断 | `CPartyClass::Isleader` | 暂无 | 未覆盖 | 第一个成员或显式 leader 字段作为队长。 |
| 35 | 成员维护 | 更换队长 | `CPartyClass::ChangeLeader` | 暂无 | 未覆盖 | 支持队长离线、主动转让或事件规则导致的队长变更。 |
| 36 | 成员维护 | 断线退队 | `ProtocolCore` 关闭流程调用 `CGPartyDelUser` | `ObjectManager.DeletePlayer` 未接入 | 未覆盖 | 玩家下线时自动退出队伍并通知成员。 |
| 37 | 成员维护 | 地图服移动退队/保留 | `MapServerMove` 相关检查 | 暂无 | 未覆盖 | 明确换线、跨服、地图服移动时队伍状态。 |
| 38 | 成员维护 | 队伍状态并发保护 | GameServer 临界区语义 | Go 暂无 | 未覆盖 | Go 侧需要锁或单协程串行处理队伍修改。 |
| 39 | 列表与通知 | 请求队伍列表 | `CGPartyList` | `handle` 有 0x42 | 未覆盖 | 返回当前队伍成员名称、地图、血量、状态。 |
| 40 | 列表与通知 | 广播队伍列表 | `CGPartyListAll` | 暂无 | 未覆盖 | 成员增删、等级变化、地图变化时广播。 |
| 41 | 列表与通知 | 删除成员通知 | `GCPartyDelUserSend` | 暂无 | 未覆盖 | 退队/踢出/解散时通知客户端。 |
| 42 | 列表与通知 | 无消息删除通知 | `GCPartyDelUserSendNoMessage` | 暂无 | 未覆盖 | 某些断线或事件场景下静默移除。 |
| 43 | 列表与通知 | 队伍血量同步 | `PartyMemberLifeSend` | 暂无 | 未覆盖 | 定时向队友下发 HP/MP/SD 等。 |
| 44 | 列表与通知 | 队友 Buff 同步 | `GCDisplayBuffeffectPartyMember` | Buff 系统未接队伍 | 未覆盖 | 队友 Buff 状态变化时下发给队伍成员。 |
| 45 | 列表与通知 | 队友坐标同步 | `GCPartyMemberPosition` | 空实现 | 未覆盖 | 按开关定时发送队友地图和坐标。 |
| 46 | 队伍聊天 | 组队聊天前缀 | `PChatProc` 中 `~`/`]` | `Object.Chat` 有分支但空 | 未覆盖 | 对象入口解析到队伍聊天后调用组队模块广播。 |
| 47 | 队伍聊天 | 队伍聊天目标 | `PartyNumber` 成员遍历 | 暂无 | 未覆盖 | 只发送给同队在线成员。 |
| 48 | 队伍聊天 | 队伍聊天颜色 | `ChatColors.Party` | `conf.Common.ChatColor.Party` | 部分覆盖 | 配置已存在，需在协议下发或消息构造中使用。 |
| 49 | 队伍聊天 | 跨服组队聊天 | ServerGroup 相关聊天 | 暂无 | 未覆盖 | 如果后续有跨服通信，组队聊天需支持跨服务器成员。 |
| 50 | 组队匹配 | 注册招募 | `CGReqRegWantedPartyMember` | `handle` 0xEF00 | 未覆盖 | 实现队伍招募信息发布。 |
| 51 | 组队匹配 | 获取匹配列表 | `CGReqGetPartyMatchingList` | `handle` 0xEF01 | 未覆盖 | 返回可申请队伍列表。 |
| 52 | 组队匹配 | 申请加入 | `CGReqJoinMemberPartyMatching` | `handle` 0xEF02 | 未覆盖 | 玩家向招募队伍申请加入。 |
| 53 | 组队匹配 | 接受申请 | `CGReqAcceptJoinMemberPartyMatching` | `handle` 0xEF03 | 未覆盖 | 队长接受申请并执行加入流程。 |
| 54 | 组队匹配 | 取消匹配 | `CGReqCancelPartyMatching` | `handle` 0xEF04/0xEF06 | 未覆盖 | 招募者或申请者取消匹配状态。 |
| 55 | 事件授权 | IllusionTemple 队伍授权 | `EnterITL_PartyAuth/AllAgreeEnterITL` | 活动未实现 | 未覆盖 | 组队模块提供事件入场确认状态，不实现事件本体。 |
| 56 | 事件授权 | ITR 队伍授权 | `EnterITR_PartyAuth/AllAgreeEnterITR` | 活动未实现 | 未覆盖 | 记录 ITR 入场确认、同意人数和清理。 |
| 57 | 事件授权 | DSF 队伍授权 | `EnterDSF_PartyAuth/AllAgreeEnterDSF` | `handle` 有 DSF 协议占位 | 未覆盖 | 记录 DSF 入场确认、同意人数和清理。 |
| 58 | 跨系统联动 | 组队经验成员过滤 | `gObjExpParty` | `09-exp.md` 待接入 | 未覆盖 | 组队模块提供同地图、距离、成员等级等基础数据。 |
| 59 | 跨系统联动 | 组队任务掉落 | `QuestMonsterItemDropParty` | `11-quests.md` 待接入 | 未覆盖 | 提供队伍成员遍历和任务资格查询入口。 |
| 60 | 跨系统联动 | 队伍掉落归属 | Party drop ownership | `14-drops.md` 待接入 | 未覆盖 | 掉落系统使用队伍关系判断共享或归属。 |
| 61 | 跨系统联动 | Gens 组队限制 | `CanGensJoinPartyWhileOppositeGens` | `17-gens.md` 待实现 | 未覆盖 | Gens 模块决定不同阵营是否允许组队。 |
| 62 | 跨系统联动 | BattleZone 组队限制 | `CanGensCreatePartyOnBattleZone` | `17-gens.md` 待实现 | 未覆盖 | BattleZone 内创建/加入队伍按 Gens 配置限制。 |
| 63 | 测试与校验 | 邀请/同意流程测试 | Party request/result | 暂无 | 未覆盖 | 覆盖创建队伍、加入队伍、拒绝、目标失效、满员。 |
| 64 | 测试与校验 | 退队/踢出/断线测试 | Party delete | 暂无 | 未覆盖 | 覆盖队长退出、普通成员退出、断线、解散。 |
| 65 | 测试与校验 | 跨系统回归测试 | Exp/Drop/Quest/Gens hooks | 暂无 | 未覆盖 | 验证组队模块只提供关系和遍历，不直接实现经验、掉落、任务规则。 |
