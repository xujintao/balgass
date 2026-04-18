你现在位于 Go 项目 server-game 的仓库根目录：

/home/pi/balgass/src/server-game

另有 C++ 参考项目 GameServer：

/home/pi/balgass-igc/igc/9.5.1.15/source/GameServer

当前仓库已经有 `.prompts/todo/` 目录。该目录中的文件是模块 TODO 和功能缺口清单，不是直接用于编码实现的提示词。

你的任务：根据 `.prompts/todo/` 中的 TODO 文档，生成 `.prompts/implement/` 目录及其中的模块实现提示词文件。每个 implement 文件都必须是可直接复制给编码模型执行的“模块实现任务书”，用于指导下一步 vibe coding。

不要实现 Go 代码。
不要修改 `.prompts/todo/`。
不要生成 HTML。
不要把 todo 表格机械复制成 implement 文件。

## 一、输入与输出

输入目录：

`.prompts/todo/`

必须读取：

- `.prompts/todo/00-todo.md`
- `.prompts/todo/01-runtime.md`
- `.prompts/todo/02-accounts.md`
- `.prompts/todo/03-characters.md`
- `.prompts/todo/04-objects.md`
- `.prompts/todo/05-formula.md`
- `.prompts/todo/06-maps.md`
- `.prompts/todo/07-items.md`
- `.prompts/todo/08-shops.md`
- `.prompts/todo/09-exp.md`
- `.prompts/todo/10-skills.md`
- `.prompts/todo/11-quests.md`
- `.prompts/todo/12-mix.md`
- `.prompts/todo/13-buffs.md`
- `.prompts/todo/14-drops.md`
- `.prompts/todo/15-party.md`
- `.prompts/todo/16-guild.md`
- `.prompts/todo/17-gens.md`
- `.prompts/todo/18-friends-mail.md`
- `.prompts/todo/19-trade.md`
- `.prompts/todo/20-personal-shops.md`
- `.prompts/todo/21-dungeons.md`
- `.prompts/todo/22-events.md`
- `.prompts/todo/23-world-events.md`
- `.prompts/todo/24-pets-summons.md`
- `.prompts/todo/25-monster-ai.md`
- `.prompts/todo/26-script.md`
- `.prompts/todo/27-security.md`
- `.prompts/todo/28-external-comm.md`
- `.prompts/todo/29-ops.md`
- `.prompts/todo/30-utils.md`
- `.prompts/todo/31-helper.md`
- `.prompts/todo/32-service-bots.md`
- `.prompts/todo/33-test-bots.md`

输出目录：

`.prompts/implement/`

必须生成：

- `.prompts/implement/01-runtime.md`
- `.prompts/implement/02-accounts.md`
- `.prompts/implement/03-characters.md`
- `.prompts/implement/04-objects.md`
- `.prompts/implement/05-formula.md`
- `.prompts/implement/06-maps.md`
- `.prompts/implement/07-items.md`
- `.prompts/implement/08-shops.md`
- `.prompts/implement/09-exp.md`
- `.prompts/implement/10-skills.md`
- `.prompts/implement/11-quests.md`
- `.prompts/implement/12-mix.md`
- `.prompts/implement/13-buffs.md`
- `.prompts/implement/14-drops.md`
- `.prompts/implement/15-party.md`
- `.prompts/implement/16-guild.md`
- `.prompts/implement/17-gens.md`
- `.prompts/implement/18-friends-mail.md`
- `.prompts/implement/19-trade.md`
- `.prompts/implement/20-personal-shops.md`
- `.prompts/implement/21-dungeons.md`
- `.prompts/implement/22-events.md`
- `.prompts/implement/23-world-events.md`
- `.prompts/implement/24-pets-summons.md`
- `.prompts/implement/25-monster-ai.md`
- `.prompts/implement/26-script.md`
- `.prompts/implement/27-security.md`
- `.prompts/implement/28-external-comm.md`
- `.prompts/implement/29-ops.md`
- `.prompts/implement/30-utils.md`
- `.prompts/implement/31-helper.md`
- `.prompts/implement/32-service-bots.md`
- `.prompts/implement/33-test-bots.md`

## 二、todo 与 implement 的区别

`.prompts/todo/*.md` 是 backlog 和缺口清单。它记录 GameServer 有什么、server-game 覆盖了什么、还缺什么。

`.prompts/implement/*.md` 是编码任务书。它必须告诉编码模型：

- 这次要实现什么。
- 先读哪些 server-game 文件。
- 参考哪些 GameServer 文件。
- 本模块边界是什么。
- 第一阶段目标是什么。
- 暂不实现什么。
- 推荐实现顺序是什么。
- 重点修哪些问题。
- 如何验收。
- 需要补哪些测试。

不要把 100 多行 TODO 表格直接复制进去。要把 TODO 转换成适合一次 vibe coding 会话执行的任务说明。

## 三、生成前必须阅读代码

为每个模块生成 implement 文件前，必须先阅读对应 TODO 文件，并阅读该模块相关的 server-game Go 代码与 GameServer C++ 代码。

server-game 常见参考路径包括：

- `game/game.go`
- `game/object`
- `game/object/player`
- `game/object/monster`
- `game/maps`
- `game/item`
- `game/skill`
- `game/formula`
- `game/model`
- `game/bot`
- `game/lang`
- `game/random`
- `game/math2`
- `handle`
- `conf`

GameServer 常见参考路径包括：

- `GameMain.cpp`、`GameServer.cpp`
- `user.h`、`user.cpp`
- `ObjAttack.cpp`、`ObjUseSkill.cpp`
- `gObjMonster.cpp`
- `protocol.cpp`、`DSProtocol.cpp`
- `BotSystem.*`
- `OfflineLevelling.*`
- `Lang.*`、`TRandomPoolMgr.*`、`ReadScript.*`、`WzMemScript*`
- 与当前模块相关的活动、副本、社交、交易、战盟、Gens、安全、运营文件

implement 文件中的“请先阅读 server-game 文件”和“请参考 GameServer 文件”必须来自真实代码路径或真实类/函数名，不要凭空编造。

## 四、每个 implement 文件的固定结构

每个 `.prompts/implement/*.md` 文件必须使用以下结构：

```md
# 模块名实现提示词

你现在要为 Go 项目 `server-game` 实现并补强“模块名”。

## 背景

说明该模块在 server-game 中的现状，以及它和 GameServer 的关系。

## 请先阅读 server-game 文件

列出必须阅读的 Go 文件、包、结构或函数。

## 请参考 GameServer 文件

列出必须参考的 C++ 文件、类、函数或结构。

## 模块边界

说明本模块负责什么，不负责什么。跨系统能力只能调用对应系统，不要顺手实现其他模块。

## 第一阶段目标

只选择最重要、最基础、最能支撑后续功能的部分作为第一阶段。不要试图一次实现 todo 文件中的全部条目。

## 实现要求

列出编码约束、架构约束、并发约束、事务约束、协议约束。

## 推荐实现顺序

用 5 到 10 步描述合理实现路径。

## 重点检查并修复的问题

从 todo 中挑出关键的 `需修正`、`部分覆盖`、`未覆盖` 条目，转换成检查项。

## 暂不实现

明确哪些复杂功能本轮不做，只保留扩展点。

## 验收标准

列出完成后应该满足的行为。

## 测试要求

列出必须新增或补充的测试场景。

## 输出要求

要求编码模型直接修改代码，并在最终回复中说明：
- 修改了哪些核心路径
- 哪些功能已完成
- 哪些功能保留为后续
- 跑了哪些测试
- 是否有未解决风险
```

## 五、状态转换规则

根据 todo 表格的 `状态` 列转换任务优先级：

- `已覆盖`：不要重复要求实现；可以要求补测试、补边界、补错误处理或补文档。
- `部分覆盖`：优先作为补强目标。
- `需修正`：优先作为修复目标，通常应进入第一阶段。
- `未覆盖`：不要全部一次性实现；只选择第一阶段最核心能力。
- `已区分`：只写进模块边界或排除项，不作为实现任务。

不要把协议表中出现过某个消息就当成已实现。只有 handler 能调到真实业务逻辑，并完成主要状态变化，才算实现。

## 六、粒度控制规则

每个 implement 文件必须适合一次 vibe coding 会话。

禁止生成这种提示词：

- “一次实现完整任务系统”
- “一次实现完整副本系统”
- “一次实现完整世界事件系统”
- “一次实现所有 120 个 TODO”

必须生成这种提示词：

- “第一阶段先实现任务状态模型、接取/完成最小闭环”
- “第一阶段先实现 BloodCastle 入场状态机和基础奖励”
- “第一阶段先实现 MuBot 配置保存/加载/启停，不做完整 OfflineLevelling”
- “第一阶段先实现 BotSystem 配置加载和 Bot 对象创建，不做 BotShop 全事务”

如果某个模块很大，第一阶段应该只覆盖 10 到 30 个最核心 TODO。其余写入“暂不实现”或“后续阶段”。

## 七、强制模块边界

对象系统：

- 只生成对象生命周期、状态承载、移动/视野/攻击/技能/道具/交互入口相关实现提示词。
- 不生成完整地图、道具、公式、经验、任务、合成、掉落表、社交实现任务。

公式系统：

- 只生成数值计算、角色重算、装备/技能/经验/经济公式相关实现提示词。
- 不生成脚本运行时基础设施实现任务。

脚本系统：

- 只生成 Lua/WZ/SMD 解析、脚本调用、Go/Lua 绑定、脚本化业务基础设施实现提示词。
- 不吞公式系统、基础工具系统和具体活动业务。

基础工具系统：

- 只生成语言包、随机数、数学工具、路径/配置辅助、通用解析辅助、基础诊断 helper 的实现提示词。
- 不生成 MuHelper、MuBot、OfflineLevelling 或 BotSystem 实现任务。

助手挂机系统：

- 只生成 MuHelper、MuBot、OfflineLevelling、自动战斗、自动拾取、自动修理、挂机计费、挂机限制相关实现提示词。
- 不生成 GameServer `BotSystem` 的 Buffer Bot、Trade Bot、BotShop、HideSeek Bot。
- 不生成 server-game `game/bot` fake connection/压测功能。

服务Bot系统：

- 只生成 GameServer `BotSystem` 的业务 Bot 实现提示词：Buffer Bot、Trade/Alchemist Bot、BotShop、HideSeek Bot、Bot 对象、Bot 外观装备、Bot 交互事务。
- 不生成玩家 MuBot/OfflineLevelling。
- 不生成 server-game `game/bot` fake connection。

测试Bot/压测系统：

- 只生成 server-game `game/bot` 的 fake connection、模拟登录、模拟选角、模拟在线、断线重连、批量压测、调试命令相关实现提示词。
- 不生成 GameServer `BotSystem` 业务 Bot。
- 不生成玩家 MuBot/OfflineLevelling。

运营管理系统：

- 只生成 GM、后台 API、公告、维护控制、在线统计、踢人/封禁/禁言入口、运营审计相关实现提示词。
- 不直接生成账号、角色、对象、安全、活动、经济等业务实现任务。

外部通信系统：

- 只生成 ConnectServer、JoinServer、DataServer、ExDB、MapServer、跨服通信、请求响应、重连与服务状态相关实现提示词。
- 不生成业务规则本身。

安全风控系统：

- 只生成封包校验、CRC、AntiHack 心跳、速度检测、移动检测、攻击检测、聊天限流、处罚与审计相关实现提示词。
- 不生成正常业务动作实现任务。

## 八、每个模块建议的第一阶段方向

以下只是默认方向。生成 implement 文件时仍要结合对应 todo 和当前代码重新判断。

- `01-runtime.md`：启动、主循环、网络连接生命周期、定时器和关闭流程稳定化。
- `02-accounts.md`：账号创建、列表、删除、登录认证、登录态和唯一在线。
- `03-characters.md`：角色列表、创建、删除、检查、加载、保存和在线唯一性。
- `04-objects.md`：对象池、玩家/怪物对象生命周期、视野和移动主链路稳定化。
- `05-formula.md`：统一 `Player.calc()` 幂等角色数值重算。
- `06-maps.md`：地图属性、阻挡、安全区、Gate、MoveCommand 和对象容器。
- `07-items.md`：背包、装备栏、仓库、物品移动、拾取、丢弃和耐久。
- `08-shops.md`：NPC 商店买卖、修理和基础事务回滚。
- `09-exp.md`：经验来源、杀怪经验、等级提升和经验结果下发。
- `10-skills.md`：技能学习/使用校验、技能列表协议、基础消耗和冷却。
- `11-quests.md`：任务状态模型、接取、进度、完成和奖励最小闭环。
- `12-mix.md`：Chaos Box 最小事务、材料校验、成功率和回滚。
- `13-buffs.md`：BuffSlot、增删查清、持续时间、协议同步。
- `14-drops.md`：怪物基础掉落、金币、掉落归属和地图投放。
- `15-party.md`：组队邀请、响应、成员列表、退队/踢出、队长。
- `16-guild.md`：战盟创建、加入、退出、成员列表和职位。
- `17-gens.md`：Gens 加入/退出、阵营信息、贡献点和基础限制。
- `18-friends-mail.md`：好友列表、申请、删除、在线状态、Memo 邮件最小闭环。
- `19-trade.md`：交易请求、响应、交易栏、金币、确认、取消和回滚。
- `20-personal-shops.md`：开店、定价、商品列表、购买事务和关店。
- `21-dungeons.md`：先选一个副本做入场、状态、怪物、奖励最小闭环。
- `22-events.md`：先选一个普通活动做开关、时间、奖励和掉落最小闭环。
- `23-world-events.md`：先选一个世界事件做状态机和运营控制骨架。
- `24-pets-summons.md`：宠物装备/卸下、耐久、基础属性修正和协议同步。
- `25-monster-ai.md`：仇恨、追击、巡逻、目标选择和怪物技能入口。
- `26-script.md`：Lua/WZ/SMD 解析与脚本调用基础设施，不做具体业务。
- `27-security.md`：协议校验、移动/攻击速度检测、聊天限流和审计。
- `28-external-comm.md`：Connect/Join/DataServer 通信骨架、请求响应和重连。
- `29-ops.md`：GM 权限、命令注册、后台鉴权、公告、踢线/封禁入口和审计。
- `30-utils.md`：语言包、随机池、数学 helper、路径/配置辅助和基础诊断。
- `31-helper.md`：MuBot 配置保存/加载/启停状态，不做完整 OfflineLevelling。
- `32-service-bots.md`：BotSystem 配置加载和 Bot 对象创建，不做 BotShop 全事务。
- `33-test-bots.md`：BotManager 并发安全、批量创建、状态观测和生产开关。

## 九、验收检查

生成完成后必须执行只读检查：

- 确认 `.prompts/implement/` 存在。
- 确认 33 个 implement 文件都存在。
- 确认每个文件都有这些章节：
  - `背景`
  - `请先阅读 server-game 文件`
  - `请参考 GameServer 文件`
  - `模块边界`
  - `第一阶段目标`
  - `实现要求`
  - `推荐实现顺序`
  - `重点检查并修复的问题`
  - `暂不实现`
  - `验收标准`
  - `测试要求`
  - `输出要求`
- 确认没有修改 `.prompts/todo/`。
- 使用 `rg` 检查关键边界词：
  - `MuBot`
  - `OfflineLevelling`
  - `BotSystem`
  - `game/bot`
  - `已覆盖`
  - `部分覆盖`
  - `需修正`
  - `未覆盖`
  - `已区分`

最终回复用户时，只需要简要说明生成了哪些 implement 提示词文件、做了哪些检查、是否有未完成事项。
