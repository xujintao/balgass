你现在位于 Go 项目 server-game 的仓库根目录：

/home/pi/balgass/src/server-game

另有 C++ 参考项目 GameServer：

/home/pi/balgass-igc/igc/9.5.1.15/source/GameServer

目标：重新生成 `.prompts/todo/` 目录及其中全部 Markdown 文件，用于后续逐模块补功能、写 prompt 和推进开发。这个生成过程必须可重复：先阅读两个项目的真实代码，再根据 GameServer 的业务逻辑与 server-game 的当前覆盖情况生成 TODO 文档，而不是机械复制旧文档。

## 一、必须生成的文件

创建或重建 `.prompts/todo/`，并生成以下文件：

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

同时生成 `.prompts/todo/00-todo.md`，内容必须是以上 33 个正式模块总表。默认不要生成候选系统区；除非重新阅读两个项目后发现新的高价值候选系统，并且能明确说明它为什么暂不转正。

不要生成或更新 `00-todo.html`，除非用户另行明确要求。

## 二、代码阅读要求

生成文档前必须先阅读代码。不要凭印象写。

必须阅读 server-game 中与模块相关的 Go 代码，例如：

- `game/game.go`
- `game/object`
- `game/object/player`
- `game/object/monster`
- `game/maps`
- `game/item`
- `game/model`
- `game/bot`
- `game/lang`
- `game/random`
- `game/math2`
- `handle`
- `conf`

必须阅读 GameServer 中与模块相关的 C++ 代码，例如：

- `GameMain.cpp`、`GameServer.cpp`
- `user.h`、`user.cpp`
- `ObjAttack.cpp`、`ObjUseSkill.cpp`
- `gObjMonster.cpp`
- `protocol.cpp`、`DSProtocol.cpp`
- `Map*`、`Gate*`
- `Item*`、`ChaosBox*`、`JewelMix*`
- `Quest*`
- `BuffEffectSlot*`
- `BotSystem.*`
- `OfflineLevelling.*`
- `Lang.*`、`TRandomPoolMgr.*`、`ReadScript.*`、`WzMemScript*`
- 各活动、副本、世界事件、社交、交易、战盟、Gens、安全和运营相关文件

如果某个模块在 GameServer 中由多个文件共同实现，要在该模块文档里写出主要参考点，不要只引用一个入口文件。

## 三、文档格式

每个模块文件必须采用统一结构：

```md
# 编号. 模块名

本模块覆盖……。本模块不拥有……；它只……。

明确排除：……归 `xx.md`；……归 `yy.md`。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | ... | ... | ... | ... | ... | ... |
```

表格规则：

- `序号` 从 1 连续编号。
- `子模块分类` 写功能分组，例如“模块边界与总入口”“对象管理器”“MuBot 协议”“BotShop”等。
- `子模块` 使用功能级或函数级颗粒度，优先对应真实函数、真实结构、真实协议、真实配置或明确业务动作。
- `GameServer` 写 C++ 侧参考类、函数、结构或文件。
- `server-game` 写 Go 侧对应包、函数、结构、协议或“暂无”。
- `状态` 只能使用：`已覆盖`、`部分覆盖`、`未覆盖`、`需修正`、`已区分`。
- `TODO 说明` 必须是可执行工程任务，避免“完善系统”“参考实现”这类空话。

每个模块建议生成约 100 到 160 个子模块。不要为了凑数量强行拆无意义条目；如果某模块天然较小，可以少于 100，但必须解释粒度来自真实代码和业务边界。

## 四、模块边界规则

必须严格避免重复归属。

对象系统：

- 只负责对象生命周期、状态承载、移动/视野/攻击/技能/道具/交互入口。
- 不吞完整地图规则、道具规则、公式、经验、技能定义、任务、合成、掉落表、社交系统。
- 战斗、死亡、PVP、掉落触发可以有对象入口，但完整规则归对应模块。

地图系统：

- 只负责地图数据、属性、阻挡、安全区、Gate、MoveCommand、小地图、地图对象容器。
- 不拥有对象生命周期、道具规则、活动状态机。

道具系统：

- 负责物品定义、背包、仓库、装备栏、耐久、套装、卓越、Socket、Pentagram 等物品能力。
- 不负责商店交易、合成事务、掉落表、奖励发放状态机。

商店系统：

- 负责 NPC 商店买卖、修理、回购、特殊商店边界。
- 个人商店归 `20-personal-shops.md`。
- BotShop 归 `32-service-bots.md` 编排，并调用个人商店/道具/经济能力。

公式系统：

- 只负责数值公式：基础属性、攻击、防御、攻速、命中、装备/技能/经验/经济公式。
- 不负责脚本运行时，也不负责业务流程状态机。

脚本系统：

- 负责 Lua 运行时、WZ/SMD 解析、脚本调用、Go/Lua 绑定和脚本化业务基础设施。
- 不把公式系统、基础工具系统和具体活动业务吞进去。

基础工具系统：

- 只负责语言包、随机数、数学工具、路径/配置辅助、通用解析辅助和基础诊断 helper。
- `MuHelper`、`MuBot`、`OfflineLevelling` 不属于基础工具。
- `BotSystem` 不属于基础工具。

助手挂机系统：

- 只负责玩家自己的自动化挂机能力：MuHelper、MuBot、OfflineLevelling、自动战斗、自动拾取、自动修理、挂机计费、挂机限制。
- 不拥有对象、地图、技能、经验、掉落、道具、经济、安全或运营规则，只调用这些系统。

服务Bot系统：

- 只负责 GameServer `BotSystem` 的业务服务 Bot：Buffer Bot、Trade/Alchemist Bot、BotShop、HideSeek Bot、Bot 对象、Bot 外观装备、Bot 交互事务。
- 不包含玩家 MuBot/OfflineLevelling。
- 不包含 server-game `game/bot` 的 fake connection。

测试Bot/压测系统：

- 只负责 server-game `game/bot` 这类开发/压测工具：fake connection、模拟登录、模拟选角、模拟在线、断线重连、批量压测、调试命令。
- 不属于游戏业务 Bot。
- 生产环境必须有关闭或权限控制说明。

运营管理系统：

- 负责 GM、后台 API、公告、维护控制、在线统计、踢人/封禁/禁言入口、运营审计。
- 不直接拥有账号、角色、对象、安全、活动、经济等业务规则；只做入口、编排和审计。

外部通信系统：

- 负责 ConnectServer、JoinServer、DataServer、ExDB、MapServer、跨服通信、请求响应、重连与服务状态。
- 不负责业务规则本身。

安全风控系统：

- 负责封包校验、CRC、AntiHack 心跳、速度检测、移动检测、攻击检测、聊天限流、处罚与审计。
- 不实现正常业务动作，只判定、限制、处罚和记录。

## 五、00-todo.md 生成要求

`00-todo.md` 必须包含：

```md
# GameServer 到 server-game 模块 TODO 总表

这个作为后续逐模块补功能、写 prompt 和推进开发的 backlog。

## 正式模块

| 顺序 | 模块 | 文件 | 说明 |
|---:|---|---|---|
...
```

正式模块必须从 1 到 33 连续编号，文件名必须与第一节文件清单完全一致。

默认不要添加“推荐 prompt 顺序”。

默认不要添加状态含义表格。

默认不要添加 HTML 文件。

## 六、状态判断原则

`已覆盖`：

- server-game 已有对应实际业务逻辑，不只是协议登记、空方法或注释。

`部分覆盖`：

- server-game 有结构、协议、配置、数据模型或部分流程，但缺关键校验、事务、状态机或边界处理。

`未覆盖`：

- server-game 基本没有对应实现。

`需修正`：

- server-game 有实现，但存在明显 bug、边界错误、安全风险、并发风险或行为与 GameServer 语义冲突。

`已区分`：

- 该条目主要用于边界说明，表示已经明确不归本模块，或与另一个模块区分清楚。

不要把协议表中出现过某个消息就判断为已覆盖。只有 handler 调到真实业务逻辑，并完成主要状态变化，才算已覆盖或部分覆盖。

## 七、输出质量要求

- 使用中文。
- 使用 Markdown。
- 文件名全部使用小写英文和连字符。
- 表格中不要写超长段落，TODO 说明要短而可执行。
- 每个模块开头必须说明“覆盖什么”和“不拥有什么”。
- 跨系统关系写清楚，但不要重复实现归属。
- 如果 GameServer 的模块是 C++ 全局对象或散落函数，要在 GameServer 列写主要类/函数/文件名。
- 如果 server-game 只有空方法、TODO、注释或协议注册，要写“部分覆盖”或“未覆盖”，不能写“已覆盖”。
- 如果 server-game 暂无对应代码，写“暂无”。
- 如果某个能力当前应该归另一个模块，状态写“已区分”，TODO 说明写“归 xx.md，本模块只调用/只排除”。

## 八、验收检查

生成后必须执行只读检查：

- 确认 `.prompts/todo/00-todo.md` 存在。
- 确认 `.prompts/todo/01-runtime.md` 到 `.prompts/todo/33-test-bots.md` 都存在。
- 确认总表正式模块编号连续到 33。
- 确认每个模块文件都有标题、边界说明和统一表格。
- 确认不生成 `00-todo.html`。
- 使用 `rg` 检查关键边界词：`MuBot`、`OfflineLevelling`、`BotSystem`、`game/bot`、`00-todo.md`。

最终回复用户时，只需要简要说明生成了哪些文件、做了哪些检查、是否有未完成事项。
