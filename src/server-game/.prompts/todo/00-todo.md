# GameServer 到 server-game 模块 TODO 总表

这个作为后续逐模块补功能、写 prompt 和推进开发的 backlog。

## 正式模块

| 顺序 | 模块 | 文件 | 说明 |
|---:|---|---|---|
| 1 | 运行时 | `01-runtime.md` | 启动、主循环、网络、协议、定时器、配置、日志 |
| 2 | 账号系统 | `02-accounts.md` | 登录、登出、登录态、账号认证、安全校验、账号在线唯一性、账号模型、账号管理 API、账号仓库归属 |
| 3 | 角色系统 | `03-characters.md` | 角色列表、创建/删除/检查角色、加载角色、角色 DB 模型、进入游戏初始化、角色保存、角色在线唯一性 |
| 4 | 对象系统 | `04-objects.md` | 对象池、玩家/怪物/NPC、对象移动、对象视野、对象攻击、对象行为、对象道具、对象状态 |
| 5 | 公式系统 | `05-formula.md` | 职业基础数值、属性点、角色重算、攻击/防御/攻速/命中公式、装备/技能/经验/经济公式 |
| 6 | 地图系统 | `06-maps.md` | 地图数据、地图属性、阻挡、安全区、Gate、MoveCommand、小地图、地图对象容器 |
| 7 | 道具系统 | `07-items.md` | 物品定义、背包、仓库、装备栏、耐久、套装、卓越、Socket、Pentagram |
| 8 | 商店系统 | `08-shops.md` | NPC 商店、买卖、修理、售出回购、特殊商店、CashShop 边界、商店事务 |
| 9 | 经验系统 | `09-exp.md` | 经验来源、经验分配、等级提升、大师等级、经验倍率应用、经验结果下发 |
| 10 | 技能系统 | `10-skills.md` | 普通技能、大师技能、技能学习、技能消耗、技能范围、技能 CD、Combo |
| 11 | 任务系统 | `11-quests.md` | 老任务、新任务、任务条件、任务进度、任务奖励、职业任务 |
| 12 | 合成系统 | `12-mix.md` | Chaos Box、Jewel Mix、Socket Mix、Pentagram Mix、翅膀合成、合成事务 |
| 13 | Buff系统 | `13-buffs.md` | Buff 定义、BuffSlot、增删查清、持续时间、数值效果、Debuff、状态限制、期限 Buff、协议同步 |
| 14 | 掉落系统 | `14-drops.md` | 怪物基础掉落、金币、卓越、Bag、事件掉落、指定掉落、套装掉落、掉落归属、掉落倍率 |
| 15 | 组队系统 | `15-party.md` | 组队邀请、响应、队伍列表、退队/踢出、队长、队友坐标、队伍血条、组队匹配、事件入场授权 |
| 16 | 战盟系统 | `16-guild.md` | 战盟创建/解散、加入/退出、成员列表、职位、公告、战盟标识、战盟视野、联盟/敌对、战盟匹配、战盟战 |
| 17 | 家族系统 | `17-gens.md` | Gens 加入/退出、阵营信息、BattleZone、贡献点、排名、奖励、PK 惩罚、Gens 跨系统限制 |
| 18 | 好友与邮件系统 | `18-friends-mail.md` | 好友列表、好友申请、删除好友、在线状态、Memo 邮件、邮件读取/删除/列表、好友聊天室邀请 |
| 19 | 玩家交易系统 | `19-trade.md` | 玩家交易请求、响应、交易栏物品、金币、确认、取消、成交、失败回滚、交易状态冲突 |
| 20 | 个人商店系统 | `20-personal-shops.md` | 玩家个人商店、商品定价、开关店、商品列表、购买事务、搜索、日志、离线交易 |
| 21 | 副本系统 | `21-dungeons.md` | BloodCastle、DevilSquare、ChaosCastle、IllusionTemple、ImperialGuardian、DoppelGanger |
| 22 | 普通活动系统 | `22-events.md` | EventChip、Rena、Lotto、LuckyCoin、BonusEvent、节日掉落、普通地图入侵活动 |
| 23 | 世界事件系统 | `23-world-events.md` | CastleSiege、Crywolf、Kanturu、Raklion、ArcaBattle、AcheronGuardian |
| 24 | 宠物与召唤系统 | `24-pets-summons.md` | Helper 宠物、坐骑、DarkHorse、DarkSpirit、Fenrir、Muun、召唤技能、召唤物 |
| 25 | 怪物 AI 系统 | `25-monster-ai.md` | 怪物 AI、仇恨、追击、巡逻、怪物技能、怪群、特殊怪物行为 |
| 26 | 脚本系统 | `26-script.md` | Lua 运行时、WZ/SMD 解析、脚本调用、Go/Lua 绑定、脚本化 Bag、任务/副本/Buff 脚本基础设施 |
| 27 | 安全风控系统 | `27-security.md` | 封包校验、CRC、AntiHack 心跳、速度检测、移动检测、攻击检测、聊天限流、处罚与审计 |
| 28 | 外部通信系统 | `28-external-comm.md` | ConnectServer、JoinServer、DataServer、ExDB、MapServer、跨服通信、请求响应、重连与服务状态 |
| 29 | 运营管理系统 | `29-ops.md` | GM、后台 API、公告、维护控制、在线统计、踢人/封禁/禁言入口、运营审计 |
| 30 | 基础工具系统 | `30-utils.md` | 语言包、随机数、数学工具、路径/配置辅助、通用解析辅助、基础诊断 helper |
| 31 | 助手挂机系统 | `31-helper.md` | MuHelper、MuBot、OfflineLevelling、自动战斗、自动拾取、自动修理、挂机计费、挂机限制 |
| 32 | 服务Bot系统 | `32-service-bots.md` | Buffer Bot、Trade Bot、BotShop、HideSeek Bot、Bot对象、Bot外观装备、Bot交互事务 |
| 33 | 测试Bot/压测系统 | `33-test-bots.md` | 模拟连接、模拟登录、模拟选角、模拟在线、断线重连、压测、开发调试命令 |
