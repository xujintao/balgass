# 22. 普通活动系统

普通活动系统拥有 EventChip、Rena、Lotto、LuckyCoin、BonusEvent、节日掉落和普通地图入侵活动的业务规则；活动积分、抽奖记录、跨服状态同步等外部服务通信归 `28-external-comm.md`；GM 手动开关活动、重载、跳阶段和运营公告入口归 `29-ops.md`。

本模块覆盖非实例副本类活动：EventChip、Rena、Lotto、LuckyCoin、BonusEvent、节日掉落、普通地图入侵活动、活动 NPC、MuRummy/卡牌类活动等。普通活动系统是活动编排层，不拥有对象、道具、掉落、经验、Buff 的底层能力；它负责判断活动是否开启、何时触发、触发什么，再把对象创建交给对象系统，把奖励生成交给掉落/道具系统，把经验倍率交给经验系统，把活动 Buff 交给 Buff 系统。BloodCastle、DevilSquare、ChaosCastle、IllusionTemple、ImperialGuardian、DoppelGanger 归 `21-dungeons.md`；CastleSiege、Crywolf、Kanturu、Raklion、Arca、Acheron 归世界事件候选。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 模块边界与总入口 | EventManager 总管理器 | `CEventManagement`、`g_EventManager` | 暂无 `game/events` | 未覆盖 | 建立普通活动服务，统一管理非副本、非世界事件活动。 |
| 2 | 模块边界与总入口 | 活动类型枚举 | `EVENT_ID_*`、各活动类 | 暂无 | 未覆盖 | 定义 EventChip、BonusEvent、Dragon、Eledorado、RingAttack、XMas、MuRummy 等类型。 |
| 3 | 模块边界与总入口 | 活动注册表 | `RegEvent` | 暂无 | 未覆盖 | 将各活动实现注册到普通活动管理器。 |
| 4 | 模块边界与总入口 | 活动启停总入口 | `StartEvent`、`Run` | 暂无 | 未覆盖 | 提供 Start/Stop/Run 接口，供运行时或 GM 调用。 |
| 5 | 模块边界与总入口 | 活动状态查询 | 各活动 `GetState/IsEventEnable` | 暂无 | 未覆盖 | 给掉落、经验、对象、NPC 查询活动是否开启。 |
| 6 | 模块边界与总入口 | 手动活动控制 | `Start_Menual/End_Menual` | 暂无 | 未覆盖 | GM 手动入口归 `29-ops.md`，普通活动系统只执行启停和状态切换。 |
| 7 | 配置与时间表 | EventManagement 时间表 | `EVENT_ID_TIME` | `conf.Events` 有 server 节点 | 部分覆盖 | 加载活动日、小时、分钟和是否已启动。 |
| 8 | 配置与时间表 | 普通活动 XML/INI 配置 | `Load`、`LoadFile`、各活动 `Load` | `conf/config.go` 已读大量字段 | 部分覆盖 | 将配置字段归并到普通活动服务使用。 |
| 9 | 配置与时间表 | 每日状态重置 | today year/month/day | 暂无 | 未覆盖 | 每日清理活动启动标记、兑换次数和掉落计数。 |
| 10 | 配置与时间表 | 活动开关缓存 | `m_bEventEnable`、`m_bDoEvent` | 配置字段散落 | 未覆盖 | 每个活动统一启用/禁用判断。 |
| 11 | 配置与时间表 | 活动通知开关 | `m_bEventNotice` | 暂无 | 未覆盖 | BonusEvent、入侵活动开始结束时可广播。 |
| 12 | 配置与时间表 | 活动重载 | `LoadFile/LoadScript` | 暂无 | 未覆盖 | 支持重载配置后刷新活动规则。 |
| 13 | BonusEvent 全局加成 | BonusEvent 管理器 | `CBonusEvent`、`g_BonusEvent` | 当前无 `BonusEvent` | 未覆盖 | 实现全局经验、ML 经验、掉落、卓越掉落加成服务。 |
| 14 | BonusEvent 全局加成 | BonusEvent 配置加载 | `CBonusEvent::LoadFile` | 暂无 | 未覆盖 | 加载星期、开始小时、结束小时和加成值。 |
| 15 | BonusEvent 全局加成 | BonusEvent 当前事件选择 | `m_curEvent_ptr` | 暂无 | 未覆盖 | 按当前时间选择生效的 BonusEvent。 |
| 16 | BonusEvent 全局加成 | 普通经验加成 | `GetAddExp` | `09-exp.md` 记录缺口 | 未覆盖 | 经验系统查询并叠加普通经验倍率。 |
| 17 | BonusEvent 全局加成 | 大师经验加成 | `GetAddMLExp` | `09-exp.md` 记录缺口 | 未覆盖 | 经验系统查询并叠加大师经验倍率。 |
| 18 | BonusEvent 全局加成 | 普通掉落加成 | `GetAddDrop` | `14-drops.md` 记录缺口 | 未覆盖 | 掉落系统查询并叠加普通掉落概率。 |
| 19 | BonusEvent 全局加成 | 卓越掉落加成 | `GetAddExcDrop` | `14-drops.md` 记录缺口 | 未覆盖 | 掉落系统查询并叠加卓越掉落概率。 |
| 20 | BonusEvent 全局加成 | BonusEvent 并发保护 | critical section | Go 单协程/锁待设计 | 未覆盖 | 当前活动指针更新和读取需要安全。 |
| 21 | 收集兑换协议 | EventChip 协议核心 | `EventChipEventProtocolCore` | `handle` 0x95/0x96/0x97 占位 | 部分覆盖 | 将 EventChip 请求委托到普通活动服务。 |
| 22 | 收集兑换协议 | EventChip 注册请求 | `PMSG_ANS_REGISTER_EVENTCHIP` | `registerEventChip` | 未覆盖 | 解析物品位置，扣除 EventChip 并累计账号/角色计数。 |
| 23 | 收集兑换协议 | EventChip 信息查询 | `PMSG_EVENTCHIPINFO` | 暂无完整响应 | 未覆盖 | 返回当前 EventChip 数量和 MutoNumber。 |
| 24 | 收集兑换协议 | MutoNumber 查询 | `PMSG_GETMUTONUMBER_RESULT` | `getMutoNum` | 未覆盖 | 查询并返回玩家的 MutoNumber 信息。 |
| 25 | 收集兑换协议 | EventChip 结束/重置 | `PMSG_ANS_RESET_EVENTCHIP` | `endEventChip` | 未覆盖 | 活动结束或重置玩家 EventChip 计数。 |
| 26 | 收集兑换协议 | EventChip DB 查询 | `EGRecvEventChipInfo` | 暂无 | 未覆盖 | 对接账号/角色活动计数持久化。 |
| 27 | 收集兑换协议 | EventChip DB 注册结果 | `EGResultRegEventChip` | 暂无 | 未覆盖 | 根据 DB 结果返回客户端成功或失败。 |
| 28 | 收集兑换协议 | EventChip 对象字段 | `EventChipCount` | `object.Object` 有注释字段 | 未覆盖 | 恢复或重设 EventChip 运行时字段。 |
| 29 | Rena/LordMark/Lotto | Rena 兑换 Zen | `EGRecvChangeRena` | `useRenaChangeZen` | 部分覆盖 | 使用 Rena/瑞娜兑换 Zen，处理物品扣除和金额上限。 |
| 30 | Rena/LordMark/Lotto | Rena 掉落率 | `MarkOfTheLord`/Rena 配置 | `RenaDropRate` | 部分覆盖 | 活动系统提供 Rena/LordMark 掉落触发条件。 |
| 31 | Rena/LordMark/Lotto | Lotto 注册 | `reqLottoRegister` | `handle` 0x9D 占位 | 部分覆盖 | 处理抽奖/乐透注册、资格、消耗和结果。 |
| 32 | Rena/LordMark/Lotto | Lotto DB 记录 | lotto register flow | 暂无 | 未覆盖 | 持久化抽奖注册状态和发奖结果。 |
| 33 | LuckyCoin | LuckyCoin 信息查询 | `PMSG_ANS_LUCKYCOIN` | `reqLuckyCoinInfo` 0xBF0B | 部分覆盖 | 返回玩家 LuckyCoin 数量。 |
| 34 | LuckyCoin | LuckyCoin 注册 | `PMSG_ANS_REG_LUCKYCOIN` | `reqLuckyCoinRegister` 0xBF0C | 部分覆盖 | 扣除 LuckyCoin 道具并累计数量。 |
| 35 | LuckyCoin | LuckyCoin 兑换 | LuckyCoin trade | `reqLuckyCoinTrade` 0xBF0D | 部分覆盖 | 按数量和配置兑换奖励。 |
| 36 | LuckyCoin | LuckyCoin 掉落率 | `LuckyCoinDrop` | `LuckyCoinDropRate` | 部分覆盖 | 掉落系统向活动系统查询 LuckyCoin 是否可掉。 |
| 37 | LuckyCoin | LuckyCoin 对象字段 | `LuckyCoinCount` | `object.Object` 有注释字段 | 未覆盖 | 恢复或重设 LuckyCoin 运行时字段。 |
| 38 | 节日掉落活动 | FireCracker 活动 | FireCracker event | `FireCrackerEventEnable/DropRate` | 部分覆盖 | 判断爆竹活动是否开启并提供掉落概率。 |
| 39 | 节日掉落活动 | Medal 活动 | Medal event | `MedalEventEnable/Silver/GoldDropRate` | 部分覆盖 | 提供银牌、金牌掉落触发。 |
| 40 | 节日掉落活动 | Halloween 活动 | Halloween event | `HalloweenEventEnable/LuckyPumpkinDropRate` | 部分覆盖 | 提供南瓜类道具掉落和使用触发。 |
| 41 | 节日掉落活动 | HeartOfLove 活动 | Heart of Love | `HeartOfLoveEventEnable/DropRate` | 部分覆盖 | 提供爱心掉落触发。 |
| 42 | 节日掉落活动 | CherryBlossom 活动 | CherryBlossom | `CherryBlossomEventEnable/BoxDropRate` | 部分覆盖 | 提供樱花箱掉落触发；合成归 `12-mix.md`。 |
| 43 | 节日掉落活动 | CandyBox 活动 | CandyBox | `CandyBoxEventEnable` 和三色配置 | 部分覆盖 | 按等级范围和概率掉落不同糖果盒。 |
| 44 | 节日掉落活动 | RibbonBox 活动 | RibbonBox | `RibbonBoxEventEnable` 和三色配置 | 部分覆盖 | 按等级范围和概率掉落圣诞丝带盒。 |
| 45 | 节日掉落活动 | Silver/Gold Box 活动 | Box drop rate | `BoxSilverDropRate/BoxGoldDropRate` | 部分覆盖 | 提供银箱、金箱掉落触发。 |
| 46 | 节日掉落活动 | HiddenTreasureBox 活动 | HiddenTreasureBox | `HiddenTreasureBoxDropRate` | 部分覆盖 | 隐藏宝箱掉落触发，合成/开启奖励归相关模块。 |
| 47 | 节日掉落活动 | MysteryBead 活动 | Mystery Bead | `SecretGemDropRate1/2` | 部分覆盖 | 神秘珠/宝石类活动掉落触发。 |
| 48 | 节日掉落活动 | 变身戒指掉落 | Ring of Transform | `ItemRingTransformDropEnable/Rate` | 部分覆盖 | 变身戒指掉落由活动判断，效果归 `13-buffs.md`/宠物候选。 |
| 49 | 节日掉落活动 | CondorFlame 掉落 | DarkLord heart event | `CondorFlameDropRate` | 部分覆盖 | 神鹰火种掉落开关和概率接入。 |
| 50 | 节日掉落活动 | 活动掉落统一入口 | `gEventMonsterItemDrop` 多分支 | `14-drops.md` 记录缺口 | 未覆盖 | 怪物死亡时活动系统选择是否触发特殊掉落。 |
| 51 | 活动 NPC | Christmas NPC | `MerryXMasTalkNpc` | `NPCChristmasEnable` | 部分覆盖 | 圣诞 NPC 对话、奖励和 Buff 触发。 |
| 52 | 活动 NPC | NewYear NPC | `HappyNewYearTalkNpc` | `NPCNewYearEnable` | 部分覆盖 | 新年 NPC 对话和奖励触发。 |
| 53 | 活动 NPC | SantaVillage 次数限制 | SantaCheck/visit count | `SantaVillage` 配置存在 | 部分覆盖 | 圣诞老人访问次数、领奖限制和 DB 记录。 |
| 54 | 活动 NPC | Santa Buff | `XMasAttackEvent` buff duration | `13-buffs.md` 记录缺口 | 未覆盖 | 活动 NPC 触发 Buff，状态归 Buff 系统。 |
| 55 | 活动 NPC | 活动 NPC Talk 路由 | NPC talk event branch | `Player.Talk` 基础存在 | 未覆盖 | 对象/NPC 对话时委托普通活动系统判断。 |
| 56 | 普通地图入侵活动 | EventManagement 调度 | `CEventManagement::Run` | 暂无 | 未覆盖 | 按时间表启动入侵活动。 |
| 57 | 普通地图入侵活动 | DragonEvent 管理器 | `CDragonEvent` | 暂无 | 未覆盖 | 红龙/龙群活动的生成、清理、状态和奖励。 |
| 58 | 普通地图入侵活动 | DragonEvent 配置 | `LoadScript`、`DRAGON_EVENT_INFO` | 暂无 | 未覆盖 | 加载地图、坐标、怪物类型、数量。 |
| 59 | 普通地图入侵活动 | DragonEvent 怪物识别 | `IsDragonEventMonster` | 暂无 | 未覆盖 | 怪物死亡时识别是否是 Dragon 活动怪。 |
| 60 | 普通地图入侵活动 | DragonEvent 清怪 | `ClearMonster` | 对象删除能力存在 | 未覆盖 | 活动结束时清理活动怪。 |
| 61 | 普通地图入侵活动 | Eledorado 管理器 | `CEledoradoEvent` | 暂无 | 未覆盖 | 黄金怪/黄金军团活动管理。 |
| 62 | 普通地图入侵活动 | Eledorado 怪物刷新 | `RegenGoldGoblen/RegenTitan/...` | 暂无 | 未覆盖 | 按怪物类型、地图、坐标和数量刷新黄金怪。 |
| 63 | 普通地图入侵活动 | Eledorado 怪物识别 | `IsEledoradoMonster` | 暂无 | 未覆盖 | 怪物死亡时识别黄金怪并触发奖励。 |
| 64 | 普通地图入侵活动 | Eledorado Boss 地图记录 | `m_Boss*MapNumber` | 暂无 | 未覆盖 | 记录黄金龙、泰坦等 Boss 刷新地图和坐标。 |
| 65 | 普通地图入侵活动 | RingAttack 管理器 | `CRingAttackEvent` | `conf.Events.RingAttack` 有 server 节点 | 部分覆盖 | 管理指环攻击/怪群活动。 |
| 66 | 普通地图入侵活动 | RingAttack MonsterHerd | `CRingMonsterHerd` | 怪群 AI 未实现 | 未覆盖 | 生成、移动、攻击和掉落 Ring 活动怪群。 |
| 67 | 普通地图入侵活动 | XMasAttack 管理器 | `CXMasAttackEvent` | `conf.Events.ChristmasAttack` 有 server 节点 | 部分覆盖 | 管理圣诞怪群活动。 |
| 68 | 普通地图入侵活动 | XMasAttack MonsterHerd | `CXMasMonsterHerd` | 怪群 AI 未实现 | 未覆盖 | 生成、移动、攻击和掉落圣诞活动怪群。 |
| 69 | 普通地图入侵活动 | AttackEvent 通用怪群 | `AttackEvent` | 暂无 | 未覆盖 | 白法师/侵略军等通用入侵活动边界。 |
| 70 | 普通地图入侵活动 | 怪群移动 | `Move`、MonsterHerd | `25-monster-ai.md` | 未覆盖 | 活动系统触发怪群移动，AI 细节归怪物 AI 系统。 |
| 71 | 普通地图入侵活动 | 活动怪死亡奖励 | `MonsterHerdItemDrop` | `14-drops.md` 待接入 | 未覆盖 | 活动怪死亡时触发对应 EventBag 或特殊掉落。 |
| 72 | 普通地图入侵活动 | 活动怪清理 | `ClearMonster/StopEvent` | 对象管理器可删除对象 | 未覆盖 | 活动结束、重载或异常时清理活动怪。 |
| 73 | EventBag 边界 | EventBag 注册 | `EventBag.h/cpp` | `14-drops.md` 记录缺口 | 未覆盖 | 普通活动只选择 Bag，具体生成归掉落系统。 |
| 74 | EventBag 边界 | CommonBag 活动触发 | `CCommonBag` | `14-drops.md` 待实现 | 未覆盖 | 使用活动道具或 NPC 奖励时触发 CommonBag。 |
| 75 | EventBag 边界 | MonsterBag 活动触发 | `CMonsterBag` | `14-drops.md` 待实现 | 未覆盖 | 活动怪死亡时触发 MonsterBag。 |
| 76 | EventBag 边界 | EventBag 活动触发 | `CEventBag` | `14-drops.md` 待实现 | 未覆盖 | 活动上下文触发 EventBag。 |
| 77 | EventBag 边界 | Gremory/CashCoin 出口 | `UseBag_GremoryCase/AddCashCoin` | 暂无 | 未覆盖 | 只记录边界，具体奖励系统后续独立。 |
| 78 | MuRummy/卡牌活动 | MuRummy 管理器 | `CMuRummyMng` | 暂无 | 未覆盖 | 管理 MuRummy 活动状态、卡牌、分数和奖励。 |
| 79 | MuRummy/卡牌活动 | MuRummy 配置 | `MuRummyInfo`、`MuRummyCardInfo` | 暂无 | 未覆盖 | 加载卡牌定义、组合、奖励和掉落概率。 |
| 80 | MuRummy/卡牌活动 | MuRummy 掉落 | `GetMuRummyEventItemDropRate` | `14-drops.md` 记录缺口 | 未覆盖 | 活动开启时掉落卡牌或活动物品。 |
| 81 | MuRummy/卡牌活动 | MuRummy 玩家状态 | card state/session | 暂无 | 未覆盖 | 保存玩家卡牌、分数、组合和领奖状态。 |
| 82 | MuRummy/卡牌活动 | MuRummy 奖励 | reward table | 暂无 | 未覆盖 | 根据卡牌组合发放奖励。 |
| 83 | 活动道具使用 | 活动道具使用入口 | Event item use branches | 道具使用系统未完整 | 未覆盖 | 道具系统使用活动物品时委托普通活动系统。 |
| 84 | 活动道具使用 | 万圣节道具 Buff | Halloween buff types | `13-buffs.md` 已记录 | 未覆盖 | 使用南瓜类道具触发 Buff。 |
| 85 | 活动道具使用 | 樱花道具 Buff | CherryBlossom buff types | `13-buffs.md` 已记录 | 未覆盖 | 使用樱花道具触发 Buff。 |
| 86 | 活动道具使用 | 活动箱子开启 | Box/Medal/Heart/EventItem | `14-drops.md` 待实现 | 未覆盖 | 开启活动箱子时触发 Bag 或奖励表。 |
| 87 | 活动道具使用 | 活动道具每日限制 | event item daily limit | 暂无 | 未覆盖 | 控制活动道具使用次数或领奖次数。 |
| 88 | 活动库存/事件背包 | EventInventory 开关 | `EventInventory` | `conf.EventInventory` 已读取 | 部分覆盖 | 判断事件背包 UI/功能是否开启。 |
| 89 | 活动库存/事件背包 | EventInventory 日期 | event inventory date | `conf.EventInventory.Date` | 部分覆盖 | 控制活动背包生效周期。 |
| 90 | 活动库存/事件背包 | EventInventory 保存 | event inventory save | 暂无 | 未覆盖 | 活动物品背包持久化边界。 |
| 91 | 活动库存/事件背包 | EventInventory 过期清理 | expire flow | 暂无 | 未覆盖 | 活动结束或日期过期后清理事件背包物品。 |
| 92 | 数据与持久化 | 活动账号数据 | EventChip/LuckyCoin/Santa | 暂无统一表/仓库 | 未覆盖 | 账号级活动计数、次数、注册状态持久化。 |
| 93 | 数据与持久化 | 活动角色数据 | MuRummy/EventInventory | 暂无 | 未覆盖 | 角色级活动状态、卡牌、事件背包持久化。 |
| 94 | 数据与持久化 | 活动每日计数 | Event1 today max、visit count | 部分配置已读 | 未覆盖 | 活动掉落每日上限和领奖次数计数。 |
| 95 | 数据与持久化 | 活动日志 | LogAdd event logs | 暂无 | 未覆盖 | 记录兑换、掉落、刷怪、发奖和异常。 |
| 96 | 数据与持久化 | 活动审计 | anti-abuse event logs | 暂无 | 未覆盖 | 非法兑换、伪造物品、重复领奖需要审计。 |
| 97 | 跨系统接口 | 与对象系统联动 | MonsterHerd/create/delete | `04-objects.md` 需接入 | 未覆盖 | 活动系统调用对象系统生成/清理活动怪。 |
| 98 | 跨系统接口 | 与对象死亡联动 | `MonsterHerdItemDrop`、event monster death | 怪物死亡入口存在 | 未覆盖 | 对象系统通知 `OnMonsterDie`，活动系统判断活动奖励。 |
| 99 | 跨系统接口 | 与 NPC 对话联动 | Christmas/NewYear/Santa NPC | `Player.Talk` 基础存在 | 未覆盖 | NPC 对话委托活动系统判断奖励或 Buff。 |
| 100 | 跨系统接口 | 与道具系统联动 | event item consume/use | `07-items.md` 需接入 | 未覆盖 | 活动兑换、使用、扣除、发奖都通过道具系统。 |
| 101 | 跨系统接口 | 与经验系统联动 | `g_BonusEvent.GetAddExp` | `09-exp.md` 已记录 | 未覆盖 | 经验系统查询普通活动全局加成。 |
| 102 | 跨系统接口 | 与 Buff 系统联动 | Santa/Halloween/Cherry buffs | `13-buffs.md` 已记录 | 未覆盖 | 普通活动触发 Buff，Buff 生命周期由 Buff 系统承载。 |
| 103 | 跨系统接口 | 与掉落系统联动 | `gEventMonsterItemDrop`、EventBag | `14-drops.md` 已记录 | 未覆盖 | 掉落系统查询活动触发和 Bag 选择。 |
| 104 | 跨系统接口 | 与副本系统边界 | BC/DS/CC/IT/IG/DG | `21-dungeons.md` 已转正 | 未覆盖 | 普通活动不得实现副本状态机。 |
| 105 | 跨系统接口 | 与世界事件边界 | CastleSiege/Crywolf/Kanturu/Raklion/Arca/Acheron | C04 候选 | 未覆盖 | 普通活动不得实现大型世界事件规则。 |
| 106 | 协议与测试 | EventChip 协议测试 | 0x95/0x96/0x97 | 待实现 | 未覆盖 | 覆盖注册、查询、结束、DB 失败、物品位置错误。 |
| 107 | 协议与测试 | Rena/Lotto 测试 | 0x98/0x9D | 待实现 | 未覆盖 | 覆盖兑换、金额上限、物品不足、重复注册。 |
| 108 | 协议与测试 | LuckyCoin 测试 | 0xBF0B/0C/0D | 待实现 | 未覆盖 | 覆盖查询、注册、兑换、奖励不足和持久化。 |
| 109 | 协议与测试 | BonusEvent 测试 | add exp/drop | 待实现 | 未覆盖 | 覆盖时间窗口、普通经验、大师经验、普通掉落、卓越掉落。 |
| 110 | 协议与测试 | 节日掉落测试 | FireCracker/Medal/Halloween 等 | 待实现 | 未覆盖 | 覆盖开关关闭、等级范围、概率命中和每日上限。 |
| 111 | 协议与测试 | 入侵活动测试 | Dragon/Eledorado/Ring/XMas | 待实现 | 未覆盖 | 覆盖刷怪、移动、死亡、奖励和清理。 |
| 112 | 协议与测试 | 活动 NPC 测试 | Christmas/NewYear/Santa | 待实现 | 未覆盖 | 覆盖对话、领奖次数、Buff 触发和 DB 失败。 |
| 113 | 协议与测试 | MuRummy 测试 | MuRummy card flow | 待实现 | 未覆盖 | 覆盖卡牌掉落、组合、积分、领奖和保存。 |
| 114 | 协议与测试 | 跨系统边界测试 | events vs dungeons/world events | 待实现 | 未覆盖 | 确认副本和世界事件不误路由到普通活动系统。 |
