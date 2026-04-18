# 16. 战盟系统

本模块覆盖战盟创建/解散、加入/退出、成员列表、职位、公告、战盟标识、战盟视野、联盟/敌对关系、战盟匹配、战盟战，以及 Gens、攻城、活动等系统对战盟的联动边界。个人商店不归本模块；战盟聊天只作为战盟成员广播能力记录。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 现有入口 | 战盟协议路由 | `CGGuildRequestRecv/CGGuildRequestResultRecv/CGGuildListAll` | `handle/c1c2.go` 0x50/0x51/0x52 | 部分覆盖 | 协议表存在，业务未实现。 |
| 2 | 现有入口 | 战盟成员删除入口 | `CGGuildDelUser` | `handle` 0x53 | 未覆盖 | 实现退出、踢出、解散时的统一入口。 |
| 3 | 现有入口 | 战盟创建答题 | `GCGuildMasterQuestionSend/CGGuildMasterAnswerRecv` | `handle` 0x54 | 未覆盖 | NPC 创建战盟前的客户端问答流程。 |
| 4 | 现有入口 | 战盟信息保存 | `CGGuildMasterInfoSave` | `handle` 0x55 | 未覆盖 | 保存战盟名和 Mark。 |
| 5 | 现有入口 | 战盟创建取消 | `CGuildMasterCreateCancel` | `handle` 0x57 | 未覆盖 | 取消创建流程并清理接口状态。 |
| 6 | 现有入口 | 战盟视野信息 | `GCGuildViewportInfo` | `handle` 0x66 | 未覆盖 | 返回视野内对象的战盟标识信息。 |
| 7 | 现有入口 | Go 战盟结构 | `_GUILD_INFO_STRUCT` | `game/guild/guild_class.go::GuildInfo{}` | 未覆盖 | 当前只有空结构，需要完整战盟实体和管理器。 |
| 8 | 配置 | 战盟创建开关 | `GuildCreate` | `conf.Common.Guild.EnableCreate` | 部分覆盖 | 配置已加载，需接入创建流程。 |
| 9 | 配置 | 战盟解散开关 | `GuildDestroy` | `conf.Common.Guild.EnableDestroy` | 部分覆盖 | 解散前校验配置、身份和城主限制。 |
| 10 | 配置 | 创建等级限制 | `GuildCreateLevel` | `conf.Common.Guild.CreateLevel` | 部分覆盖 | 创建战盟要求角色等级达标。 |
| 11 | 配置 | 最大成员数 | `MaxGuildMember` | `conf.Common.Guild.MaxMember` | 部分覆盖 | 成员上限和 Dark Lord 统率扩展规则需确定。 |
| 12 | 配置 | 联盟最小成员 | `AllianceMinGuildMember` | `conf.Common.Guild.AllianceGuildMinMember` | 部分覆盖 | 创建联盟时检查成员数。 |
| 13 | 配置 | 联盟最大数量 | `AllianceMaxGuilds` | `conf.Common.Guild.AllianceMaxGuildCount` | 部分覆盖 | 限制联盟内战盟数量。 |
| 14 | 战盟管理器 | GuildClass 初始化 | `CGuildClass::Init` | 暂无 guild manager | 未覆盖 | 新建管理器维护战盟表、成员索引和名称索引。 |
| 15 | 战盟管理器 | 添加战盟 | `CGuildClass::AddGuild` | 暂无 | 未覆盖 | 从 DB 或创建流程构建战盟对象。 |
| 16 | 战盟管理器 | 链表插入 | `CGuildClass::AddTail` | Go 可用 map/list | 未覆盖 | Go 侧可用 map 替代链表，但语义需覆盖。 |
| 17 | 战盟管理器 | 删除全部战盟 | `CGuildClass::AllDelete` | 暂无 | 未覆盖 | 支持重载、停服或测试清理。 |
| 18 | 战盟管理器 | 删除战盟 | `CGuildClass::DeleteGuild` | 暂无 | 未覆盖 | 删除战盟对象、成员引用、联盟关系和视野信息。 |
| 19 | 战盟管理器 | 按名称查找 | `CGuildClass::SearchGuild` | 暂无 | 未覆盖 | 按战盟名查找。 |
| 20 | 战盟管理器 | 按编号查找 | `SearchGuild_Number` | 暂无 | 未覆盖 | 按 GuildNumber 查找。 |
| 21 | 战盟管理器 | 编号+成员查找 | `SearchGuild_NumberAndId` | 暂无 | 未覆盖 | 校验玩家是否属于指定战盟。 |
| 22 | 成员管理 | 成员上线绑定 | `CGuildClass::ConnectUser` | `Player.LoadCharacter` 未接入 | 未覆盖 | 角色进入游戏后绑定战盟对象和成员状态。 |
| 23 | 成员管理 | 添加成员 | `CGuildClass::AddMember` | 暂无 | 未覆盖 | 加入战盟后写入成员列表并同步在线状态。 |
| 24 | 成员管理 | 删除成员 | `CGuildClass::DelMember` | 暂无 | 未覆盖 | 玩家退出或被踢时移除成员。 |
| 25 | 成员管理 | 成员下线 | `CGuildClass::CloseMember` | 暂无 | 未覆盖 | 下线只清理在线对象编号，不删除成员。 |
| 26 | 成员管理 | 设置成员服务器 | `CGuildClass::SetServer` | 暂无 | 未覆盖 | 跨服/多服场景维护成员所在服务器。 |
| 27 | 成员管理 | 重建成员总数 | `BuildMemberTotal` | 暂无 | 未覆盖 | 根据成员列表统计总人数和在线人数。 |
| 28 | 成员管理 | 成员状态查询 | `GetGuildMemberStatus` | 暂无 | 未覆盖 | 获取普通成员、副盟主、战斗队长等职位。 |
| 29 | 成员管理 | 成员状态设置 | `SetGuildMemberStatus` | `handle` 0xE1 | 未覆盖 | 盟主设置成员职位。 |
| 30 | 战盟类型 | 获取战盟类型 | `GetGuildType` | 暂无 | 未覆盖 | 获取普通/联盟主盟等类型。 |
| 31 | 战盟类型 | 设置战盟类型 | `SetGuildType` | `handle` 0xE2 | 未覆盖 | 设置战盟类型并同步 DB。 |
| 32 | 创建流程 | 创建 NPC 入口 | `GCGuildMasterManagerRun` | `Player.Talk` 未接入 | 未覆盖 | 对战盟创建 NPC 打开创建流程。 |
| 33 | 创建流程 | 战盟名校验 | `CGuildMasterInfoSave` | 暂无 | 未覆盖 | 长度、字符、重复名、敏感词校验。 |
| 34 | 创建流程 | 战盟 Mark 校验 | `GUILDMARK` | 暂无 | 未覆盖 | 校验客户端提交的战盟标志数据。 |
| 35 | 创建流程 | 创建费用校验 | GameServer 创建流程 | 暂无 | 未覆盖 | 扣除 Zen 或创建道具条件。 |
| 36 | 创建流程 | 创建 DB 请求 | DataServer guild create | 暂无 | 未覆盖 | 战盟业务定义创建语义，DataServer 请求通道归 `28-external-comm.md`。 |
| 37 | 创建流程 | 创建结果下发 | `GCGuildMasterInfoSave` 结果包 | 暂无 | 未覆盖 | 成功后刷新角色战盟状态和视野。 |
| 38 | 加入流程 | 战盟加入请求 | `CGGuildRequestRecv` | `handle` 0x50 | 未覆盖 | 玩家向盟主申请加入。 |
| 39 | 加入流程 | 加入响应 | `CGGuildRequestResultRecv` | `handle` 0x51 | 未覆盖 | 盟主同意/拒绝加入请求。 |
| 40 | 加入流程 | 加入距离/目标检查 | `CGGuildRequestRecv` 检查 | 暂无 | 未覆盖 | 校验盟主、申请者在线、同地图距离、状态。 |
| 41 | 加入流程 | 满员检查 | `MaxGuildMember` | 暂无 | 未覆盖 | 成员满时拒绝加入。 |
| 42 | 加入流程 | Gens 加入限制 | `CanGensJoinGuildWhileOppositeGens` | `17-gens.md` 待实现 | 未覆盖 | 不同阵营是否允许加入同战盟由 Gens 决定。 |
| 43 | 退出/解散 | 主动退出 | `CGGuildDelUser` | `handle` 0x53 | 未覆盖 | 普通成员主动退出。 |
| 44 | 退出/解散 | 盟主踢人 | `CGGuildDelUser` | 暂无 | 未覆盖 | 盟主或有权限职位踢出成员。 |
| 45 | 退出/解散 | 战盟解散 | `DeleteGuild` | 暂无 | 未覆盖 | 盟主解散战盟，清理所有成员和联盟关系。 |
| 46 | 退出/解散 | 城主战盟限制 | `CastleOwnerGuildDestroyLimit` | 配置已加载 | 未覆盖 | 罗兰城主战盟可能禁止解散。 |
| 47 | 列表与公告 | 战盟列表下发 | `CGGuildListAll` | `handle` 0x52 | 未覆盖 | 下发成员、职位、在线状态、服务器。 |
| 48 | 列表与公告 | 战盟公告校验 | `gGuildNoticeStringCheck` | 暂无 | 未覆盖 | 公告长度、字符、过滤规则。 |
| 49 | 列表与公告 | 战盟公告保存 | Guild notice DB | 暂无 | 未覆盖 | 保存公告并广播给成员。 |
| 50 | 聊天与广播 | 战盟聊天 | `PChatProc` 中 `@` | `Object.Chat` 有分支但空 | 未覆盖 | 对象入口解析后调用战盟成员广播。 |
| 51 | 聊天与广播 | 联盟聊天 | `GCServerMsgStringSendGuild` / Union chat | `conf` 有 ServerGroupUnionChatEnable | 未覆盖 | 向同联盟战盟成员广播。 |
| 52 | 聊天与广播 | 战盟广播消息 | `GCServerMsgStringSendGuild` | 暂无 | 未覆盖 | 系统消息、公告、战盟事件通知。 |
| 53 | 视野 | 战盟标识出现 | `GCGuildViewportNowPaint` | 暂无 | 未覆盖 | 玩家进入视野时广播战盟名、Mark、盟主标识。 |
| 54 | 视野 | 战盟标识删除 | `GCGuildViewportDelNow` | 暂无 | 未覆盖 | 离开视野或退出战盟时清理客户端显示。 |
| 55 | 视野 | 视野战盟缓存 | `CViewportGuild::Init/Add` | 暂无 | 未覆盖 | 避免重复发送战盟视野信息。 |
| 56 | 联盟/敌对 | Union 管理器 | `TUnion` | 暂无 | 未覆盖 | 维护联盟主盟、成员战盟、敌对关系。 |
| 57 | 联盟/敌对 | 添加联盟 | `TUnion::AddUnion` | 暂无 | 未覆盖 | 创建联盟关系。 |
| 58 | 联盟/敌对 | 删除联盟 | `TUnion::DelUnion/DelAllUnion` | 暂无 | 未覆盖 | 解散联盟或清空全部关系。 |
| 59 | 联盟/敌对 | 关系查询 | `GetGuildRelationShip` | 暂无 | 未覆盖 | 查询两个战盟是联盟、敌对还是无关系。 |
| 60 | 联盟/敌对 | 关系统计 | `GetGuildRelationShipCount` | 暂无 | 未覆盖 | 统计联盟或敌对数量。 |
| 61 | 联盟/敌对 | 设置联盟成员列表 | `SetGuildUnionMemberList` | 暂无 | 未覆盖 | 从 DB 或协议更新联盟成员。 |
| 62 | 联盟/敌对 | 设置敌对成员列表 | `SetGuildRivalMemberList` | 暂无 | 未覆盖 | 从 DB 或协议更新敌对成员。 |
| 63 | 联盟/敌对 | 获取联盟列表 | `CGUnionList` | `handle` 0xEB??/0xE9 相关 | 未覆盖 | 下发联盟战盟列表。 |
| 64 | 联盟/敌对 | 踢出联盟成员 | `CGRelationShipReqKickOutUnionMember` | `handle` 0xEB01 | 未覆盖 | 联盟主盟踢出成员战盟。 |
| 65 | 战盟战 | 战盟战请求 | `GCGuildWarRequestSendRecv` | `handle` 0x61 | 未覆盖 | 处理战盟战挑战、接受和拒绝。 |
| 66 | 战盟战 | 战盟战宣告 | `GCGuildWarDeclare` | 暂无 | 未覆盖 | 战盟战开始时广播。 |
| 67 | 战盟战 | 战盟战结束 | `GCGuildWarEnd` | 暂无 | 未覆盖 | 结束战盟战并清理双方状态。 |
| 68 | 战盟战 | 战盟战比分 | `GCGuildWarScore` | 暂无 | 未覆盖 | 同步战盟战分数。 |
| 69 | 战盟匹配 | 匹配列表 | `CGReqGuildMatchingList` | `handle` 0xED00 | 未覆盖 | 返回战盟招募列表。 |
| 70 | 战盟匹配 | 搜索匹配 | `CGReqGuildMatchingListSearchWord` | `handle` 0xED01 | 未覆盖 | 按关键词搜索战盟。 |
| 71 | 战盟匹配 | 注册匹配 | `CGReqRegGuildMatchingList` | `handle` 0xED02 | 未覆盖 | 战盟发布招募信息。 |
| 72 | 战盟匹配 | 取消匹配 | `CGReqCancelGuildMatchingList` | `handle` 0xED03 | 未覆盖 | 取消招募。 |
| 73 | 战盟匹配 | 申请加入匹配 | `CGReqJoinGuildMatching` | `handle` 0xED04 | 未覆盖 | 玩家申请加入招募战盟。 |
| 74 | 战盟匹配 | 审批申请 | `CGReqAllowJoinGuildMatching` | `handle` 0xED06 | 未覆盖 | 盟主审批申请。 |
| 75 | 战盟匹配 | 等待列表 | `CGReqGetWaitStateListGuildMatching` | `handle` 0xED08 | 未覆盖 | 查询申请等待状态。 |
| 76 | 跨系统联动 | 攻城战盟报名 | CastleSiege guild registration | `handle` 0xB4/0xB5/B902 | 未覆盖 | 攻城系统使用战盟身份，本模块提供战盟基础能力。 |
| 77 | 跨系统联动 | Arca 战盟报名 | ArcaBattle guild join | `handle` 0xF830/0xF832 | 未覆盖 | 活动系统使用战盟成员/盟主校验。 |
| 78 | 跨系统联动 | 战盟 Buff 触发 | Guild Period Buff | `13-buffs.md` 待接入 | 未覆盖 | 战盟系统只触发，Buff 状态归 Buff 系统。 |
| 79 | 测试与校验 | 创建/加入/退出测试 | Guild protocol flows | 暂无 | 未覆盖 | 覆盖创建、重复名、加入、踢出、解散。 |
| 80 | 测试与校验 | 联盟/敌对测试 | TUnion/TUnionInfo | 暂无 | 未覆盖 | 覆盖联盟创建、踢出、敌对、关系查询。 |
