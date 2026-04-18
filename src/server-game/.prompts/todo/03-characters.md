# 3. 角色系统

本模块覆盖角色列表、创建角色、删除角色、检查角色名、加载角色、角色 DB 模型、角色外观帧、进入游戏初始化、角色保存、回选角与角色在线唯一性。账号认证、登录安全、账号管理 API 和账号态清理归账号系统；对象移动、视野、攻击、技能、道具交互归对象系统；DataServer/JoinServer 请求响应、跨服迁移认证和中心保存通道归 `28-external-comm.md`。GM 查角色、踢角色、锁定/解锁、后台删除入口和审计归 `29-ops.md`。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 角色协议入口 | 获取角色列表路由 | `ProtocolCore` 分发 `0xF3:0x00` | `0xF300 -> GetCharacterList` | 已覆盖 | 要求仅已登录账号可请求角色列表。 |
| 2 | 角色协议入口 | 创建角色路由 | `ProtocolCore` 分发 `0xF3:0x01` | `0xF301 -> CreateCharacter` | 已覆盖 | 要求创建角色只能在已登录未进游戏状态执行。 |
| 3 | 角色协议入口 | 删除角色路由 | `ProtocolCore` 分发 `0xF3:0x02` | `0xF302 -> DeleteCharacter` | 已覆盖 | 删除角色必须拒绝 Playing 状态请求。 |
| 4 | 角色协议入口 | 加载角色路由 | `ProtocolCore` 分发 `0xF3:0x03` | `0xF303 -> LoadCharacter` | 已覆盖 | 加载角色后进入 Playing，需要防止重复加载同一连接。 |
| 5 | 角色协议入口 | 检查角色名路由 | `ProtocolCore` 扩展角色名检查语义 | `0xF315 -> CheckCharacter` | 部分覆盖 | 当前固定返回成功，需要接入 DB 名称占用和命名规则。 |
| 6 | 角色 DB 模型 | 角色 DB 模型 | DataServer Character 表 | `model.Character` | 已覆盖 | 当前字段覆盖基础角色、背包、技能、快捷键、MuBot。 |
| 7 | 角色 DB 模型 | 账号归属 | DataServer `AccountID` | `Character.AccountID` | 已覆盖 | 所有角色查询必须带账号 ID。 |
| 8 | 角色 DB 模型 | 角色名唯一性 | DataServer `Name` 唯一 | `Character.Name gorm:"unique"` | 已覆盖 | 需要明确跨账号是否全服唯一，当前是全局唯一。 |
| 9 | 角色 DB 模型 | 角色槽位 | DataServer `Index/Pos` | `Character.Position` | 已覆盖 | 创建角色使用 0-4 槽位，需测试中间空洞填补。 |
| 10 | 角色 DB 模型 | 职业与转职 | DataServer `Class/ClassSkin` | `Class/ChangeUp` | 部分覆盖 | 职业编码需要和客户端 frame 以及公式系统统一。 |
| 11 | 角色 DB 模型 | 基础属性 | DataServer `Str/Dex/Vit/Energy/Leadership` | `Strength/Dexterity/Vitality/Energy/Leadership` | 已覆盖 | 初始值来自 `CharacterTable`，后续由公式系统校验。 |
| 12 | 角色 DB 模型 | 等级经验字段 | DataServer `Level/Exp/MasterExp` | `Level/Experience/MasterLevel/MasterExperience` | 已覆盖 | MasterLevel 逻辑在经验/公式模块继续补齐。 |
| 13 | 角色 DB 模型 | 地图坐标字段 | DataServer `MapNumber/MapX/MapY/Dir` | `MapNumber/X/Y/Dir` | 已覆盖 | 加载时需验证地图合法和阻挡点修正。 |
| 14 | 角色 DB 模型 | 背包与技能 JSONB | DataServer 二进制 item/magic | `Inventory`、`Skills` JSONB | 设计替代 | Go 版结构化持久化合理，但需兼容客户端二进制编码。 |
| 15 | 角色 DB 模型 | MuKey/MuBot 字段 | GameServer OptionData/MuBot 语义 | `MuKey`、`MuBot` JSONB | 部分覆盖 | 已有保存入口，进入游戏后加载/推送还需补全。 |
| 16 | 角色列表 | 角色列表请求 | `DataServerGetCharListRequest` | `Player.GetCharacterList` | 已覆盖 | 单服可本地 DB 查询；DataServer 模式应通过 `28-external-comm.md` 发起请求。 |
| 17 | 角色列表 | 账号合法性检查 | `gObjGetAccountId`、`gObjIsAccontConnect` | 当前依赖 `p.accountID` | 部分覆盖 | 未登录账号请求角色列表应拒绝。 |
| 18 | 角色列表 | DB 列表查询 | `JGPGetCharList` 接收 DataServer 列表 | `DB.GetCharacterList` | 已覆盖 | 按 `position ASC` 排序。 |
| 19 | 角色列表 | 开放职业回复 | `PMSG_CHARLIST_ENABLE_CREATION` | `MsgEnableCharacterClassReply` | 部分覆盖 | 当前固定 `0xFF`，需按账号/等级/VIP/配置计算。 |
| 20 | 角色列表 | Reset 信息回复 | `PMSG_RESET_INFO_CHARLIST` | `MsgResetCharacterReply` | 部分覆盖 | 当前固定字符串，占位明显，需要接入角色 reset 数据。 |
| 21 | 角色列表 | 仓库扩展字段 | `WhExpansion` | `WarehouseExpansion` | 已覆盖 | 与账号仓库扩展保持一致。 |
| 22 | 角色列表 | 角色外观帧 | `JGPGetCharList` 组装 `CharSet` | `MakeCharacterFrame` | 部分覆盖 | 已实现主要装备外观，但需和 SeasonX 编码逐项对齐。 |
| 23 | 角色列表 | GuildStatus 字段 | `btGuildStatus` | `MsgCharacter.GuildStatus` 固定 `0xFF` | 未覆盖 | 公会系统完成后应从角色公会状态生成。 |
| 24 | 角色列表 | PKLevel 字段 | `btPkLevel` | `MsgCharacter.PKLevel` 固定 `0` | 未覆盖 | PK 系统完成后应展示真实 PK 等级。 |
| 25 | 角色列表 | 角色列表编码测试 | `PMSG_CHARLISTCOUNT/PMSG_CHARLIST_S9` | `MsgGetCharacterListReply.Marshal` | 未覆盖 | 覆盖空列表、多槽位、装备外观、GBK 名字。 |
| 26 | 创建角色 | 创建请求解析 | `PMSG_CHARCREATE` | `MsgCreateCharacter.Unmarshal` | 已覆盖 | 解析 GBK 名字和 class skin 高 4 位。 |
| 27 | 创建角色 | 创建状态校验 | `CGPCharacterCreate` 检查 `PLAYER_LOGGED/PLAYING` | `Player.CreateCharacter` 当前只校验 Playing 间接不足 | 部分覆盖 | 必须拒绝未登录和游戏中创建角色。 |
| 28 | 创建角色 | 全局建角开关 | `gCreateCharacter` | 当前无全局开关 | 未覆盖 | 维护期禁止建角的运营入口归 `29-ops.md`，角色系统执行创建限制。 |
| 29 | 创建角色 | BattleCore 禁止建角 | `SERVER_BATTLECORE` 分支 | 当前无 BattleCore 限制 | 未覆盖 | 若实现 BattleCore，需要禁止或转发建角。 |
| 30 | 创建角色 | 安全码检查 | `m_bSecurityCheck` | 当前无账号安全码 | 未覆盖 | 有安全码账号需先通过验证才能建删角色。 |
| 31 | 创建角色 | 职业合法性 | `ClassSkin` 白名单 `0x00..0x70` | `msg.Class > 6` 校验 | 部分覆盖 | GrowLancer/SeasonX 职业边界需明确，不能硬编码遗漏。 |
| 32 | 创建角色 | 职业开放校验 | `EnableCharacterCreate` 位检查 | 当前固定开放 | 未覆盖 | MG/DL/Summoner/RF/GL 等职业按账号条件开放。 |
| 33 | 创建角色 | 名字非法字符校验 | `g_prohibitedSymbols.Validate` | 当前仅非空和模型 validate 间接 | 未覆盖 | 补齐长度、ASCII/GBK、特殊字符、控制字符校验。 |
| 34 | 创建角色 | 脏词/保留名校验 | `SwearFilter.CompareText`、`[A]` 禁止 | 当前无 | 未覆盖 | 建角前检查脏词、GM 标记、系统保留名。 |
| 35 | 创建角色 | 创建 DB 与回复 | `SDHP_CREATECHAR`、`JGCharacterCreateRequest` | `DB.CreateCharacter`、`MsgCreateCharacterReply` | 已覆盖 | 补齐重复名、DB 错误、槽位满时的明确结果码。 |
| 36 | 删除与检查角色 | 删除请求解析 | `PMSG_CHARDELETE` | `MsgDeleteCharacter.Unmarshal` | 已覆盖 | 解析角色名和 7 字节删除密码。 |
| 37 | 删除与检查角色 | 删除登录态校验 | `CGPCharDel` 检查已登录 | `DeleteCharacter` 未显式检查未登录 | 部分覆盖 | 未登录请求应拒绝并记录异常。 |
| 38 | 删除与检查角色 | 禁止游戏中删除 | `Connected == PLAYER_PLAYING` 拒绝 | `p.ConnectState == ConnectStatePlaying` 返回 | 已覆盖 | 应返回明确失败码或断线策略，而不是静默失败。 |
| 39 | 删除与检查角色 | 防删除时间 | `bEnableDelCharacter` | 当前无 | 未覆盖 | 加载/进入游戏后一定时间内禁止删除角色。 |
| 40 | 删除与检查角色 | 公会角色限制 | `GuildNumber/lpGuild` 检查 | 当前无公会限制 | 未覆盖 | 公会系统完成后禁止删除有公会角色或转交到公会模块判断。 |
| 41 | 删除与检查角色 | 账号物品锁限制 | `m_cAccountItemBlock` | 当前无账号物品锁 | 未覆盖 | 账号锁定时禁止删除角色和敏感操作。 |
| 42 | 删除与检查角色 | 删除密码校验 | `gObjPasswordCheck` | 当前硬编码 `"1234567"` | 需修正 | 应使用账号密码或安全码策略，不能保留硬编码。 |
| 43 | 删除与检查角色 | 删除 DB 请求 | `SDHP_CHARDELETE`、`JGCharDelRequest` | `DB.DeleteCharacterByName` | 已覆盖 | 需要区分不存在、非本账号、DB 错误。 |
| 44 | 删除与检查角色 | CheckCharacter 查询 | GameServer 名称检查语义 | `Player.CheckCharacter` 固定 `Result:0` | 未覆盖 | 实现名字是否可用、非法字符、保留名和重复名检查。 |
| 45 | 删除与检查角色 | 删除/检查测试 | GameServer 多防护分支 | 当前无专门测试 | 未覆盖 | 覆盖密码错误、游戏中删除、非本账号、重复名检查。 |
| 46 | 加载角色进入游戏 | 加载请求解析 | `PMSG_CHARMAPJOIN` | `MsgLoadCharacter.Unmarshal` | 已覆盖 | 解析 GBK 名字和槽位。 |
| 47 | 加载角色进入游戏 | 加载状态校验 | `CGPCharacterMapJoinRequest` 检查已登录未 Playing | `Player.LoadCharacter` 部分校验 | 部分覆盖 | 未登录或已 Playing 时应拒绝加载。 |
| 48 | 加载角色进入游戏 | 账号一致性 | `gObjIsAccontConnect` | `DB.GetCharacterByName(accountID,name)` | 已覆盖 | DB 查询带账号 ID，能防非本账号角色加载。 |
| 49 | 加载角色进入游戏 | 防重复在线角色 | `JGGetCharacterInfo` 扫描同名/同账号 Playing | 当前无 | 未覆盖 | 加载前检查同账号或同角色是否已在线。 |
| 50 | 加载角色进入游戏 | CtlCode 封禁检查 | `lpMsg->CtlCode & 1` | 当前无角色封禁字段 | 未覆盖 | 角色封禁/解封入口归 `29-ops.md`，加载时由角色系统执行限制。 |
| 51 | 加载角色进入游戏 | DB 字段灌入 Player | `gObjSetCharacter` | `Player.LoadCharacter` 字段赋值 | 已覆盖 | 当前已覆盖基础属性、经验、HP/MP、技能、背包、坐标。 |
| 52 | 加载角色进入游戏 | 地图合法性检查 | `MapNumberCheck`、PK 红名回城 | `SpawnPosition` 仅低等级修正 | 部分覆盖 | 加载后应校验地图存在、坐标可站立、PK 强制位置。 |
| 53 | 加载角色进入游戏 | 地图服迁移 | `CheckMoveMapSvr/GJReqMapSvrMove` | 当前单服，无地图服迁移 | 设计替代 | 单服阶段可不实现；多服迁移请求、认证和响应归 `28-external-comm.md`。 |
| 54 | 加载角色进入游戏 | 加载成功回复 | `PMSG_CHARMAPJOINRESULT` | `MsgLoadCharacterReply` | 部分覆盖 | 回复字段已覆盖核心值，但 MaxHP/MaxMP/SD/AG 等依赖 calc 完整性。 |
| 55 | 加载角色进入游戏 | 加载失败处理 | `JGGetCharacterInfo` 失败关闭连接 | 当前只日志 return | 需修正 | DB 不存在、封禁、重复在线应返回失败或断线，不能静默无响应。 |
| 56 | 进入游戏后初始化 | 设置站立地图属性 | `MapC.SetStandAttr` 语义 | `MapManager.SetMapAttrStand` | 已覆盖 | 进入游戏时占用地图格，登出时清理。 |
| 57 | 进入游戏后初始化 | 创建视野范围 | `gObjViewportListProtocolCreate` | `p.CreateFrustum()` | 已覆盖 | 视野细节归对象系统，本模块记录进入游戏调用点。 |
| 58 | 进入游戏后初始化 | 设置移动/恢复参数 | GameServer 对象初始化 | `MoveSpeed`、`MaxRegenTime` | 已覆盖 | 这些默认值后续应配置或由职业/公式计算。 |
| 59 | 进入游戏后初始化 | 生命状态 | `Connected=PLAYER_PLAYING`、对象存活 | `Live=true`、`State=1` | 已覆盖 | 需要和死亡/复活模块统一状态语义。 |
| 60 | 进入游戏后初始化 | 技能数据填充 | `MagicByteConvert/GCMagicListMultiSend` | `Skills.FillSkillData`、后续 push skill | 部分覆盖 | 进入游戏后应推送技能列表，避免客户端技能栏缺数据。 |
| 61 | 进入游戏后初始化 | 物品列表推送 | `GCItemListSend` | `MsgItemListReply` 相关入口 | 部分覆盖 | LoadCharacter 后应稳定推送背包/装备列表。 |
| 62 | 进入游戏后初始化 | Master 数据推送 | `SendMLData/CGReqGetMasterLevelSkillTree` | `pushMasterSkillList` 等入口 | 部分覆盖 | MasterSkillTree 需接入加载和推送流程。 |
| 63 | 进入游戏后初始化 | 任务/社交数据加载 | Quest/Friend/Guild/Gens 请求 | 当前大多未接入 | 未覆盖 | 任务、好友、公会、Gens 等归后续模块，本处保留加载入口。 |
| 64 | 进入游戏后初始化 | Buff/宠物/活动数据加载 | PeriodBuff、Muun、Pentagram、EventInventory | 当前部分有字段/入口 | 未覆盖 | 后续对应模块实现后接入 LoadCharacter 成功后初始化链。 |
| 65 | 进入游戏后初始化 | HP/MP/SD/AG 推送 | `GCReFillSend/GCManaSend` | `MsgLoadCharacterReply` 含部分字段 | 部分覆盖 | 需要在公式系统完整后统一推送当前值和最大值。 |
| 66 | 保存登出与临时用户 | 角色保存入口 | `GJSetCharacterInfo` | `Player.saveCharacter` | 部分覆盖 | 单服直接保存基础字段；DataServer 保存通道、跨服保存标记归 `28-external-comm.md`。 |
| 67 | 保存登出与临时用户 | DB 更新字段 | `SDHP_DBCHAR_INFOSAVE` | `DB.UpdateCharacter` | 已覆盖 | 当前 Omit 账号、槽位、名字、职业、MuKey、MuBot，符合基础保存边界。 |
| 68 | 保存登出与临时用户 | 背包保存 | `ItemByteConvert32` | `Inventory` JSONB | 已覆盖 | 道具模块需保证 JSONB 和客户端编码一致。 |
| 69 | 保存登出与临时用户 | 技能保存 | `MagicByteConvert` | `Skills` JSONB | 已覆盖 | 技能学习/删除后要及时保存或统一下线保存。 |
| 70 | 保存登出与临时用户 | MuKey/MuBot 保存 | `DGOptionDataSend/DGOptionDataRecv` 等 | `UpdateCharacterMuKey/UpdateCharacterMuBot` | 部分覆盖 | 已有单独保存函数，但加载和推送需要补齐。 |
| 71 | 保存登出与临时用户 | 仓库关闭保存 | `GDSetWarehouseList`、`GDSetWarehouseMoney` | `CloseWarehouseWindow`、`UpdateAccountWarehouse` | 部分覆盖 | 仓库模块需补打开、锁定、金币、异常回滚。 |
| 72 | 保存登出与临时用户 | 登出回选角 | `gObjGameClose` 后 `PLAYER_LOGGED` | `Logout Flag=1` 保存并 `Reset` | 部分覆盖 | 需要确认 `Reset` 后保留账号登录态且清理地图/视野/队伍。 |
| 73 | 保存登出与临时用户 | 断线删除对象 | `gObjDel`、`GJPUserClose` | `Player.Offline`、`ObjectManager.DeletePlayer` | 部分覆盖 | 断线时应保存角色、清理地图、清理账号在线索引。 |
| 74 | 保存登出与临时用户 | 异步保存队列 | `DbSave::Add/ThreadProc` | 当前同步 `DB.UpdateCharacter` | 设计替代 | Go 版可先同步保存；高并发后再评估异步保存队列。 |
| 75 | 保存登出与临时用户 | 临时用户/跨图保留 | `TemporaryUserManager` | 当前无 | 未覆盖 | 多地图服/活动重连前可暂不实现，但需作为未来扩展点记录。 |
