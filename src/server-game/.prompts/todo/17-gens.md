# 17. Gens 系统

本模块覆盖 Gens 阵营加入/退出、阵营信息、BattleZone、贡献点、排名、奖励、PK 惩罚、阵营聊天，以及 Gens 对组队、战盟、联盟、掉落、移动和私聊的限制。Gens 不是战盟子模块，而是独立阵营/PVP/排行系统。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 现有入口 | 加入 Gens 协议 | `CGReqRegGensMember` | `handle/c1c2.go` 0xF801 | 部分覆盖 | 路由存在，业务未实现。 |
| 2 | 现有入口 | 退出 Gens 协议 | `CGReqSecedeGensMember` | `handle` 0xF803 | 部分覆盖 | 路由存在，业务未实现。 |
| 3 | 现有入口 | 领取奖励协议 | `CGReqGensReward` | `handle` 0xF809 | 部分覆盖 | 路由存在，业务未实现。 |
| 4 | 现有入口 | 请求成员信息协议 | `CGReqGensMemberInfo` | `handle` 0xF80B | 部分覆盖 | 路由存在，业务未实现。 |
| 5 | 现有入口 | 玩家 Gens 字段 | `m_PlayerData->m_GensInfluence` 等 | `Player` 注释字段 | 未覆盖 | 恢复或重设阵营、贡献、排名、奖励状态。 |
| 6 | 现有入口 | Gens 聊天入口 | `PChatProc` 中 `$` | `Object.Chat` 有分支但空 | 未覆盖 | 对象入口解析后交给 Gens 模块广播。 |
| 7 | 初始化 | GensSystem 加载 | `CGensSystem::LoadData` | 暂无 Gens manager | 未覆盖 | 加载 BattleZone、排行、奖励、限制等配置。 |
| 8 | 初始化 | Gens 成员数量请求 | `ReqExDBGensMemberCount` | 暂无 | 未覆盖 | 统计两阵营成员数。 |
| 9 | 初始化 | Gens 成员数量设置 | `SetGensMemberCount` | 暂无 | 未覆盖 | 维护阵营人数用于显示和限制。 |
| 10 | 初始化 | 获取成员数量 | `GetGensMemberCount` | 暂无 | 未覆盖 | NPC 显示和限制判断需要查询。 |
| 11 | 阵营加入 | 加入入口 | `ReqRegGensMember` | 暂无 | 未覆盖 | 执行加入 Gens 的业务流程。 |
| 12 | 阵营加入 | 加入 NPC 检查 | `IsInfluenceNPC` | 暂无 | 未覆盖 | 只有指定 NPC 可办理加入。 |
| 13 | 阵营加入 | 已加入限制 | `GetGensInfluence > 0` | 暂无 | 未覆盖 | 已加入阵营不能重复加入。 |
| 14 | 阵营加入 | 等级限制 | `Below 50Lv` 检查 | 暂无 | 未覆盖 | 低等级不能加入。 |
| 15 | 阵营加入 | 战盟盟主限制 | `Guild Leader` Gens 检查 | `16-guild.md` 待实现 | 未覆盖 | 盟主加入 Gens 会影响战盟限制。 |
| 16 | 阵营加入 | 组队限制 | `Already Partymember` 检查 | `15-party.md` 待实现 | 未覆盖 | 加入 Gens 时可能要求不在队伍中。 |
| 17 | 阵营加入 | 联盟限制 | `Union GuildMaster` 检查 | `16-guild.md` 待实现 | 未覆盖 | 联盟主盟与成员阵营需要一致。 |
| 18 | 阵营加入 | DB 注册请求 | `GDReqRegGensMember` | 暂无 | 未覆盖 | server-game 可先实现本地持久化，再抽象外部 DB。 |
| 19 | 阵营加入 | DB 注册响应 | `AnsRegGensMember` | 暂无 | 未覆盖 | 成功后设置玩家阵营并下发信息。 |
| 20 | 阵营退出 | 退出入口 | `ReqSecedeGensMember` | 暂无 | 未覆盖 | 执行退出阵营流程。 |
| 21 | 阵营退出 | 未加入限制 | `!GetGensInfluence` | 暂无 | 未覆盖 | 未加入不能退出。 |
| 22 | 阵营退出 | 退出 NPC 检查 | `IsInfluenceNPC` | 暂无 | 未覆盖 | 只有指定 NPC 可办理退出。 |
| 23 | 阵营退出 | 退出冷却 | leaving time 检查 | 暂无 | 未覆盖 | 防止频繁切换阵营。 |
| 24 | 阵营退出 | DB 退出请求 | `ReqSecedeGensMember` / ExDB | 暂无 | 未覆盖 | 持久化清除阵营、贡献和排行状态。 |
| 25 | 阵营退出 | DB 退出响应 | `AnsSecedeGensMember` | 暂无 | 未覆盖 | 成功后刷新客户端状态。 |
| 26 | 信息同步 | 请求 Gens 信息 | `ReqExDBGensInfo` | 暂无 | 未覆盖 | 登录或客户端请求时加载阵营数据。 |
| 27 | 信息同步 | 发送 Gens 信息 | `SendGensInfo` | 暂无 | 未覆盖 | 下发阵营、等级、贡献、排名、奖励状态。 |
| 28 | 信息同步 | 视野 Gens 信息 | `GensViewportListProtocol` | 暂无 | 未覆盖 | 视野内显示玩家阵营标识。 |
| 29 | 信息同步 | BattleZone 数据 | `SendBattleZoneData` | 暂无 | 未覆盖 | 下发 BattleZone 开关、状态和阵营数据。 |
| 30 | 影响力接口 | 设置阵营 | `SetGensInfluence` | 暂无 | 未覆盖 | 设置玩家 Duprian/Vanert/无阵营。 |
| 31 | 影响力接口 | 获取阵营 | `GetGensInfluence` | 暂无 | 未覆盖 | 供组队、战盟、私聊、掉落、地图限制调用。 |
| 32 | 影响力接口 | 获取阵营名称 | `GetGensInfluenceName` | 暂无 | 未覆盖 | 日志、系统消息和 UI 显示使用。 |
| 33 | 影响力接口 | 是否已注册阵营 | `IsRegGensInfluence` | 暂无 | 未覆盖 | 简化各种前置判断。 |
| 34 | BattleZone | 地图是否 BattleZone | `IsMapBattleZone` | 暂无 | 未覆盖 | 判断地图是否开启 Gens BattleZone。 |
| 35 | BattleZone | 移动地图是否 BattleZone | `IsMoveMapBattleZone` | `06-maps.md` 待接入 | 未覆盖 | MoveCommand/Gate 判断是否进入 BattleZone。 |
| 36 | BattleZone | 用户 BattleZone 状态 | `SetUserBattleZoneEnable/IsUserBattleZoneEnable` | 暂无 | 未覆盖 | 玩家进入/离开 BattleZone 时维护状态。 |
| 37 | BattleZone | BattleZone PVP 开关 | `IsPkEnable` | `04-objects.md` 战斗待接入 | 未覆盖 | 战斗系统询问 Gens 是否允许 PVP。 |
| 38 | BattleZone | BattleZone 聊天 | `BattleZoneChatMsgSend` | `Object.Chat` 分支未实现 | 未覆盖 | BattleZone 内按阵营或规则广播聊天。 |
| 39 | BattleZone | BattleZone 掉落加成 | `GetBattleZoneDropBonus` | `14-drops.md` 引用 | 未覆盖 | 掉落系统向 Gens 查询加成。 |
| 40 | BattleZone | BattleZone 卓越加成 | `GetBattleZoneExcDropBonus` | `14-drops.md` 引用 | 未覆盖 | 卓越掉落倍率由 Gens 提供。 |
| 41 | BattleZone | BattleZone 入场类型 | `GetEntryAllowType` | 暂无 | 未覆盖 | 地图系统判断哪些阵营/状态允许进入。 |
| 42 | 贡献点 | 设置贡献点 | `SetContributePoint` | 暂无 | 未覆盖 | 设置当前贡献值。 |
| 43 | 贡献点 | 增加贡献点 | `AddContributePoint` | 暂无 | 未覆盖 | 击杀、任务或活动增加贡献。 |
| 44 | 贡献点 | 扣减贡献点 | `SubContributePoint` | 暂无 | 未覆盖 | 惩罚或消耗扣减贡献。 |
| 45 | 贡献点 | 获取贡献点 | `GetContributePoint` | 暂无 | 未覆盖 | 排名、等级、奖励判断需要查询。 |
| 46 | 贡献点 | 贡献点计算 | `CalcContributePoint` | 暂无 | 未覆盖 | 按双方等级、阵营、PK 状态计算获得贡献。 |
| 47 | 贡献点 | 保存贡献点 | `GDReqSaveContributePoint` | 暂无 | 未覆盖 | 持久化贡献变化。 |
| 48 | 排名等级 | 计算 Gens 等级 | `CalGensClass` | 暂无 | 未覆盖 | 根据贡献和排名计算 Gens Class。 |
| 49 | 排名等级 | 设置 Gens 等级 | `SetGensClass` | 暂无 | 未覆盖 | 更新玩家当前 Gens 等级。 |
| 50 | 排名等级 | 获取 Gens 等级 | `GetGensClass` | 暂无 | 未覆盖 | UI、奖励、属性或限制查询。 |
| 51 | 排名等级 | 设置 Gens 排名 | `SetGensRanking` | 暂无 | 未覆盖 | 从 DB/排行计算写入玩家排名。 |
| 52 | 排名等级 | 下一等级贡献 | `GetNextContributePoint` | 暂无 | 未覆盖 | 客户端显示下一阶所需贡献。 |
| 53 | 排名等级 | 排名保存 | `ReqExDBSetGensRanking` | 暂无 | 未覆盖 | 周期性保存排名。 |
| 54 | 奖励 | 奖励检查 | `ReqExDBGensRewardCheck` | 暂无 | 未覆盖 | 领取前向持久层检查资格。 |
| 55 | 奖励 | 奖励完成 | `ReqExDBGensRewardComplete` | 暂无 | 未覆盖 | 成功领取后标记完成。 |
| 56 | 奖励 | 奖励日设置 | `ReqExDBSetGensRewardDay` | 暂无 | 未覆盖 | 维护每日/周期奖励状态。 |
| 57 | 奖励 | 请求奖励日 | `ReqGensRewardDay` | 暂无 | 未覆盖 | 查询当前可领奖周期。 |
| 58 | 奖励 | 发送奖励 | `SendGensReward` | 暂无 | 未覆盖 | 下发奖励结果。 |
| 59 | 奖励 | 发放奖励物品 | `SendGensRewardItem` | `07-items.md` 待接入 | 未覆盖 | 检查背包后创建奖励物品。 |
| 60 | 奖励 | 奖励背包检查 | `GensRewardInventoryCheck` | `Inventory.FindFreePositionForItem` | 未覆盖 | 发奖前确认背包空间。 |
| 61 | PK 惩罚 | PK Party Level | `GetPKPartyLevel` | `15-party.md` 待接入 | 未覆盖 | 组队内 PK 等级或惩罚等级查询。 |
| 62 | PK 惩罚 | 地图移动 Zen 惩罚 | `PkPenaltyAddNeedZenMapMove` | `06-maps.md` 待接入 | 未覆盖 | PK 状态移动地图可能额外扣 Zen。 |
| 63 | PK 惩罚 | 掉落背包物品 | `PkPenaltyDropInvenItem` | `14-drops.md`/`07-items.md` 待接入 | 未覆盖 | PK 惩罚掉落玩家物品。 |
| 64 | PK 惩罚 | 掉落 Zen | `PkPenaltyDropZen` | 暂无 | 未覆盖 | PK 惩罚掉金币。 |
| 65 | PK 惩罚 | PK 惩罚 Debuff | `SendPKPenaltyDebuff` | `13-buffs.md` 待接入 | 未覆盖 | 由 Gens 触发 Debuff，状态由 Buff 系统承载。 |
| 66 | 滥杀防刷 | 检查击杀用户名 | `ChkKillUserName` | 暂无 | 未覆盖 | 防止重复刷同一目标贡献。 |
| 67 | 滥杀防刷 | 插入击杀用户名 | `InsertKillUserName` | 暂无 | 未覆盖 | 记录近期击杀目标。 |
| 68 | 滥杀防刷 | 保存滥杀记录 | `DBSaveAbusingKillUserName` | 暂无 | 未覆盖 | 持久化防刷记录。 |
| 69 | 滥杀防刷 | 请求滥杀信息 | `GDReqAbusingInfo` | 暂无 | 未覆盖 | 登录时加载防刷数据。 |
| 70 | 滥杀防刷 | 滥杀惩罚 | `AbusingPenalty` | 暂无 | 未覆盖 | 触发刷分惩罚。 |
| 71 | 滥杀防刷 | 重置滥杀信息 | `AbusingInfoReset` | 暂无 | 未覆盖 | 周期或条件清理防刷记录。 |
| 72 | 组队/战盟限制 | 不同阵营组队限制 | `CanGensJoinPartyWhileOppositeGens` | `15-party.md` 待调用 | 未覆盖 | 组队模块调用 Gens 判断。 |
| 73 | 组队/战盟限制 | BattleZone 创建队伍限制 | `CanGensCreatePartyOnBattleZone` | `15-party.md` 待调用 | 未覆盖 | BattleZone 内组队创建限制。 |
| 74 | 组队/战盟限制 | BattleZone 队伍拆分 | `MoveInBattleZonePartySplit` | `15-party.md` 待调用 | 未覆盖 | 进入 BattleZone 时按阵营拆分队伍。 |
| 75 | 组队/战盟限制 | 不同阵营入战盟限制 | `CanGensJoinGuildWhileOppositeGens` | `16-guild.md` 待调用 | 未覆盖 | 战盟加入流程调用 Gens 判断。 |
| 76 | 组队/战盟限制 | 不同阵营联盟限制 | `CanGensJoinUnionWhileOppositeGens` | `16-guild.md` 待调用 | 未覆盖 | 联盟关系建立前调用 Gens 判断。 |
| 77 | 私聊限制 | BattleZone 私聊限制 | `CGChatWhisperRecv` 中 Gens 判断 | `Object.Whisper` 未接入 | 未覆盖 | BattleZone 内可能只允许同阵营私聊。 |
| 78 | 测试与校验 | 加入/退出测试 | Gens register/secede | 暂无 | 未覆盖 | 覆盖 NPC、等级、已有阵营、战盟/队伍限制。 |
| 79 | 测试与校验 | BattleZone 测试 | BattleZone APIs | 暂无 | 未覆盖 | 覆盖地图进入、PVP、聊天、掉落加成。 |
| 80 | 测试与校验 | 贡献/奖励测试 | contribution/ranking/reward | 暂无 | 未覆盖 | 覆盖贡献变更、排行等级、奖励领取和重复领取。 |
