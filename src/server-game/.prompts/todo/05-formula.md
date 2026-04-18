# 5. 公式系统

本模块覆盖 server-game 中所有“怎么算”的规则，包括 Lua 公式桥接、职业基础攻击、防御/命中/攻速、角色重算、HP/MP/SD/AG、StatSpec、装备/翅膀/宠物/技能修正、经验等级与战斗结算。流程模块只调用公式，不在各自模块内重复公式实现。GameServer 的 `ObjCalCharacter`、`ExpManager`、`MagicDamage`、`MasterLevelSkillTreeSystem` 与战斗经验函数作为业务参考；Go 版当前保留 Lua wrapper 与 `Player.calc()` 的实现方向，后续通用 Lua 运行时、签名调用、脚本加载和热重载基础设施归 `26-script.md`。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 公式引擎与 Lua 桥接 | 公式包初始化 | `CObjCalCharacter::Init`、各系统 LoadScript | `game/formula/formula.go::init` | 已覆盖 | 保持公式系统在进程启动时加载；测试注入、热重载和通用加载器归 `26-script.md`。 |
| 2 | 公式引擎与 Lua 桥接 | Lua state 管理 | `CObjCalCharacter::m_Lua` | `formula` 结构内多个 `*lua.LState` | 已覆盖 | 当前每类脚本一个 Lua state，合理；通用 VM 生命周期和并发访问边界归 `26-script.md`。 |
| 3 | 公式引擎与 Lua 桥接 | CalcCharacter 脚本加载 | `Scripts\\Character\\CalcCharacter.lua` | `load("Character/CalcCharacter.lua")` | 已覆盖 | 需要启动期校验关键函数存在，避免运行时才发现。 |
| 4 | 公式引擎与 Lua 桥接 | StatSpec 脚本加载 | `StatSpecialize` 脚本/配置语义 | `load("Specialization/StatSpec.lua")` | 已覆盖 | StatSpec 作为角色重算的一部分，需覆盖所有客户端显示项。 |
| 5 | 公式引擎与 Lua 桥接 | ItemCalc 脚本加载 | 道具/翅膀公式脚本语义 | `load("Misc/ItemCalc.lua")` | 已覆盖 | 当前只包装翅膀公式，后续扩展装备价格/耐久等函数。 |
| 6 | 公式引擎与 Lua 桥接 | ExpCalc 脚本加载 | `MasterSkillSystem::SetExpTable` 等经验表语义 | `load("Misc/ExpCalc.lua")` | 已覆盖 | 普通/大师经验表已通过 Lua 初始化。 |
| 7 | 公式引擎与 Lua 桥接 | RegularSkillCalc 脚本加载 | `ObjUseSkill`、技能倍率脚本语义 | `load("Skills/RegularSkillCalc.lua")` | 已覆盖 | 当前仅包装骑士/魔剑技能倍率与穿刺倍率。 |
| 8 | 公式引擎与 Lua 桥接 | Lua 调用签名解析 | `Generic_Call("iiii>iiii")` | `call(ls, method, sig, args...)` | 已覆盖 | 签名字符串设计接近 GameServer；通用签名解析器归 `26-script.md`，公式系统只定义公式调用点。 |
| 9 | 公式引擎与 Lua 桥接 | Lua 返回值写回 | `Generic_Call` 输出指针 | `call` 中 `*int/*float64` 写回 | 已覆盖 | 需补返回类型不匹配测试；通用返回值校验归 `26-script.md`。 |
| 10 | 公式引擎与 Lua 桥接 | Lua 调用错误处理 | GameServer 启动/运行日志 | `call` 仅 `slog.Error` 后返回 | 需修正 | 公式调用失败是高风险；通用错误处理机制归 `26-script.md`，公式系统负责业务降级策略。 |
| 11 | 职业基础攻击公式 | 法师物理攻击 | `WizardDamageCalc` | `formula.WizardDamageCalc` | 已覆盖 | 和 Lua 脚本保持同名映射，需数值对齐 GameServer。 |
| 12 | 职业基础攻击公式 | 法师魔法攻击 | `WizardMagicDamageCalc` | `formula.WizardMagicDamageCalc` | 需修正 | 当前调用 `"WizardDamageCalc"`，疑似应改为 `"WizardMagicDamageCalc"`。 |
| 13 | 职业基础攻击公式 | 战士物理攻击 | `KnightDamageCalc` | `formula.KnightDamageCalc` | 已覆盖 | 用于 `Player.calc` 的战士基础左右手攻击。 |
| 14 | 职业基础攻击公式 | 战士魔法攻击 | `KnightMagicDamageCalc` | `formula.KnightMagicDamageCalc` | 已覆盖 | 虽然战士魔攻较弱，仍需和客户端显示一致。 |
| 15 | 职业基础攻击公式 | 弓箭手有弓攻击 | `ElfWithBowDamageCalc` | `formula.ElfWithBowDamageCalc` | 已覆盖 | `Player.calc` 已按弓/弩识别选择该公式。 |
| 16 | 职业基础攻击公式 | 弓箭手无弓攻击 | `ElfWithoutBowDamageCalc` | `formula.ElfWithoutBowDamageCalc` | 已覆盖 | 需覆盖空手/非弓武器测试。 |
| 17 | 职业基础攻击公式 | 魔剑士攻击 | `GladiatorDamageCalc`、`GladiatorMagicDamageCalc` | `formula.GladiatorDamageCalc/GladiatorMagicDamageCalc` | 已覆盖 | 同时影响物理和魔法攻击。 |
| 18 | 职业基础攻击公式 | 圣导师攻击 | `LordDamageCalc`、`LordMagicDamageCalc` | `formula.LordDamageCalc/LordMagicDamageCalc` | 已覆盖 | 物理攻击需要 leadership 参与。 |
| 19 | 职业基础攻击公式 | 召唤术士攻击 | `SummonerDamageCalc`、`SummonerMagicDamageCalc` | `formula.SummonerDamageCalc/SummonerMagicDamageCalc` | 已覆盖 | 召唤包含 curse attack，需测试诅咒上下限。 |
| 20 | 职业基础攻击公式 | 格斗/枪手职业缺口 | `RageFighterDamageCalc`、`GrowLancerDamageCalc` | `RageFighterDamageCalc`，无 GrowLancer wrapper | 部分覆盖 | RageFighter 已有，GrowLancer 缺 wrapper 和 `Player.calc` 分支。 |
| 21 | 防御/命中/速度公式 | 基础防御 | `CalcDefense` | `formula.CalcDefense` | 已覆盖 | `Player.calc` 已调用。 |
| 22 | 防御/命中/速度公式 | PVM 攻击成功率 | `CalcAttackSuccessRate_PvM` | `formula.CalcAttackSuccessRate_PvM` | 已覆盖 | 需要确认 level 使用普通等级还是普通+大师等级，当前 Go 用普通等级。 |
| 23 | 防御/命中/速度公式 | PVP 攻击成功率 | `CalcAttackSuccessRate_PvP` | `formula.CalcAttackSuccessRate_PvP` | 已覆盖 | GameServer 使用普通+大师等级，Go 当前用普通等级，需对齐。 |
| 24 | 防御/命中/速度公式 | PVM 防御成功率 | `CalcDefenseSuccessRate_PvM` | `formula.CalcDefenseSuccessRate_PvM` | 已覆盖 | `Player.calc` 已写入 `DefenseRate`。 |
| 25 | 防御/命中/速度公式 | PVP 防御成功率 | `CalcDefenseSuccessRate_PvP` | `formula.CalcDefenseSuccessRate_PvP` | 已覆盖 | 需要和 `Object.CheckMiss` 的 PVP 判定联动测试。 |
| 26 | 防御/命中/速度公式 | 攻击速度 | `CalcAttackSpeed` | `formula.CalcAttackSpeed` | 已覆盖 | `Player.calc` 已计算基础 AttackSpeed。 |
| 27 | 防御/命中/速度公式 | 魔法速度 | `CalcAttackSpeed` 第二返回值 | `formula.CalcAttackSpeed` 写入 `magicSpeed` | 已覆盖 | 武器/手套/宠物修正会继续叠加。 |
| 28 | 防御/命中/速度公式 | PVE miss 判定 | `ObjBaseAttack` 命中判定语义 | `Object.CheckMiss` PVE 分支 | 部分覆盖 | 当前公式较简化，需和 GameServer 攻防率随机规则对齐。 |
| 29 | 防御/命中/速度公式 | PVP miss 判定 | `ObjBaseAttack` PVP 命中语义 | `Object.CheckMiss` PVP 分支 | 部分覆盖 | 当前等级差 switch 顺序存在覆盖风险，需测试 100/200/300 级差。 |
| 30 | 防御/命中/速度公式 | 速度异常检测 | `SpeedHackCheck`、攻速校验 | `conf.CommonServer` 有配置，公式层未使用 | 未覆盖 | 后续安全模块应基于公式攻速提供阈值。 |
| 31 | 双持与武器修正 | 左右手武器识别 | `CalcCharacter` 中 Right/Left 有效性 | `Player.calc` 的 `left/right` bool | 已覆盖 | 排除箭/弩箭，识别武器槽。 |
| 32 | 双持与武器修正 | 弓弩特殊识别 | `ITEM_CROSSBOW/ITEM_BOW` | `KindBCrossbow/KindBBow` | 已覆盖 | 弓箭手公式选择依赖该判断。 |
| 33 | 双持与武器修正 | 相同双持加成 | `CalcTwoSameWeaponBonus` | `formula.CalcTwoSameWeaponBonus` | 已覆盖 | 战士/魔剑士双持相同武器时调用。 |
| 34 | 双持与武器修正 | 不同双持加成 | `CalcTwoDifferentWeaponBonus` | `formula.CalcTwoDifferentWeaponBonus` | 已覆盖 | 战士/魔剑士双持不同武器时调用。 |
| 35 | 双持与武器修正 | 格斗双持加成 | `CalcRageFighterTwoWeaponBonus` | `formula.CalcRageFighterTwoWeaponBonus` | 已覆盖 | RageFighter 使用独立双持公式。 |
| 36 | 双持与武器修正 | 单左手攻击转换 | `CalcCharacter` 左手攻击聚合 | `Player.calc` `case left` | 已覆盖 | 左手有效时将 leftAttack 写入面板攻击。 |
| 37 | 双持与武器修正 | 单右手攻击转换 | `CalcCharacter` 右手攻击聚合 | `Player.calc` `case right` | 已覆盖 | 右手有效时将 rightAttack 写入面板攻击。 |
| 38 | 双持与武器修正 | 空手攻击转换 | `HaveWeaponInHand=false` 语义 | `Player.calc` 默认左右均值 | 部分覆盖 | 空手/盾牌场景需和 GameServer 细节对齐。 |
| 39 | 双持与武器修正 | 武器攻速聚合 | `CalcCharacter` 武器攻速加成 | `Player.calc` 双持平均/单手加成 | 已覆盖 | 需补弓、弩、箭、盾牌组合测试。 |
| 40 | 双持与武器修正 | 魔法/诅咒百分比转换 | `m_Magic/m_Curse` 影响最终伤害 | `Player.calc` magic/curse 转换 | 需修正 | Go 当前 curse 计算使用 `curseAttackMin * curseAttackMax / 100` 可疑，应单独验证。 |
| 41 | Player.calc 重算流程 | 重算入口 | `CObjCalCharacter::CalcCharacter` | `Player.calc` | 已覆盖 | Go 版角色数值核心入口，装备变化和升级后会调用。 |
| 42 | Player.calc 重算流程 | 装备槽读取 | `Right/Left/Gloves/Boots/Wing/Helper` | `Player.calc` 读取 inventory 槽位 | 已覆盖 | 当前覆盖核心装备槽，Pentagram/耳环等后续扩展。 |
| 43 | Player.calc 重算流程 | 清理旧属性加成 | `ClearPrevEffectAll` 后清理 Add* 字段 | `Player.calc` 将 add/bonus 字段归零 | 已覆盖 | 需保证新增加成字段也在重算前清理。 |
| 44 | Player.calc 重算流程 | Buff 前置应用 | `g_BuffEffect.ApplyPrevEffectStat` | 当前无 Buff 前置公式 | 未覆盖 | Buff 系统完成后应在基础属性聚合前应用。 |
| 45 | Player.calc 重算流程 | Master stat 前置应用 | `m_MPSkillOpt` 加到属性 | 当前仅注释 master skill contribution | 未覆盖 | Master 被动属性应在基础公式前聚合。 |
| 46 | Player.calc 重算流程 | 套装 stat 前置计算 | `CalcSetItemStat` | `p.calcSetItem(true)` | 部分覆盖 | 已有入口，需确认所有 SetEffectType 对齐 GameServer。 |
| 47 | Player.calc 重算流程 | 基础属性聚合 | `Strength + AddStrength` 等 | `strength/dexterity/vitality/energy/leadership` | 已覆盖 | 这是后续所有公式输入。 |
| 48 | Player.calc 重算流程 | 基础攻防速计算顺序 | `CalcCharacter` 固定顺序 | `Player.calc` 职业攻击 -> 防御 -> 命中 -> 速度 | 已覆盖 | 保持顺序，避免装备/StatSpec 前后叠加错误。 |
| 49 | Player.calc 重算流程 | 重算后数值裁剪 | `CalcCharacter` 裁剪当前 HP/MP/SD/AG | `Player.calc` 裁剪 HP/MP/SD/AG | 已覆盖 | 当前值不能超过最大值。 |
| 50 | Player.calc 重算流程 | 重算后推送 | GameServer `GCReFillSend/GCManaSend` 等 | `MsgStatSpecReply`、攻速、MaxHP/MP/SD/AG 推送 | 部分覆盖 | 需补客户端需要的完整数值刷新包。 |
| 51 | HP/MP/SD/AG 公式 | 最大 HP 基础公式 | `CalcCharacter` Life/MaxLife | `Player.calc` `MaxHP = base + level + vitality` | 已覆盖 | 使用 `CharacterTable` 中职业成长参数。 |
| 52 | HP/MP/SD/AG 公式 | 最大 MP 基础公式 | `CalcCharacter` Mana/MaxMana | `Player.calc` `MaxMP = base + level + energy` | 已覆盖 | 使用 `CharacterTable` 中职业成长参数。 |
| 53 | HP/MP/SD/AG 公式 | 等级 HP 成长 | GameServer `LevelToLife` 语义 | `c.LevelHP` | 已覆盖 | 角色模板负责职业差异。 |
| 54 | HP/MP/SD/AG 公式 | 等级 MP 成长 | GameServer `LevelToMana` 语义 | `c.LevelMP` | 已覆盖 | 角色模板负责职业差异。 |
| 55 | HP/MP/SD/AG 公式 | 体力转 HP | `VitalityToLife` | `c.VitalityHP` | 已覆盖 | 属性点增加体力后需触发重算。 |
| 56 | HP/MP/SD/AG 公式 | 智力转 MP | `EnergyToMana` | `c.EnergyMP` | 已覆盖 | 属性点增加智力后需触发重算。 |
| 57 | HP/MP/SD/AG 公式 | MaxSD 公式 | `CObjCalCharacter::CalcShieldPoint` | `Player.calc` 使用 `SDGageConstA/B` | 部分覆盖 | 当前实现了基础公式，SD 伤害分配和恢复仍未完整。 |
| 58 | HP/MP/SD/AG 公式 | DarkLord SD 统率加成 | `CalcShieldPoint` 含 Leadership | `Player.calc` DarkLord 加 leadership | 已覆盖 | 需测试圣导师与普通职业差异。 |
| 59 | HP/MP/SD/AG 公式 | MaxAG 职业权重 | `gObjSetBP`/AG 语义 | `Player.calc` 按职业权重计算 `MaxAG` | 部分覆盖 | 需和 GameServer 各职业 BP/AG 公式精确对齐。 |
| 60 | HP/MP/SD/AG 公式 | MaxLifePower | `user.cpp::gObjCalcMaxLifePower` | 当前无独立字段 | 未覆盖 | 生命之光等技能需要 MaxLifePower 语义时再补。 |
| 61 | StatSpec 公式 | StatSpec 百分比入口 | `g_StatSpec.CalcStatOption` | `formula.StatSpec_GetPercent` | 已覆盖 | Go 通过 Lua 获取百分比，再在 `Player.calc` 应用。 |
| 62 | StatSpec 公式 | 攻击力专精 | `STAT_OPTION_INC_ATTACK_POWER` | `MsgStatSpec ID=1` | 已覆盖 | 当前按基础攻击百分比增加左右手攻击。 |
| 63 | StatSpec 公式 | 魔法攻击专精 | `STAT_OPTION_INC_MAGIC_DAMAGE` | `MsgStatSpec ID=9` | 已覆盖 | 当前按 magicAttack 百分比增加。 |
| 64 | StatSpec 公式 | 诅咒攻击专精 | `STAT_OPTION_INC_CURSE_DAMAGE` | `MsgStatSpec ID=10` | 已覆盖 | 召唤职业需重点测试。 |
| 65 | StatSpec 公式 | 防御专精 | `STAT_OPTION_INC_DEFENSE` | `MsgStatSpec ID=4` | 已覆盖 | 当前按基础防御百分比增加。 |
| 66 | StatSpec 公式 | PVM 攻击率专精 | `STAT_OPTION_INC_ATTACK_RATE` | `MsgStatSpec ID=2` | 已覆盖 | 与 `AttackRate` 联动。 |
| 67 | StatSpec 公式 | PVP 攻击率专精 | `STAT_OPTION_INC_ATTACK_RATE_PVP` | `MsgStatSpec ID=3` | 已覆盖 | 与 `attackRatePVP` 联动。 |
| 68 | StatSpec 公式 | PVM 防御率专精 | `STAT_OPTION_INC_DEFENSE_RATE` | `MsgStatSpec ID=6` | 已覆盖 | 与 `DefenseRate` 联动。 |
| 69 | StatSpec 公式 | PVP 防御率专精 | `STAT_OPTION_INC_DEFENSE_RATE_PVP` | `MsgStatSpec ID=7` | 已覆盖 | 与 `defenseRatePVP` 联动。 |
| 70 | StatSpec 公式 | 攻速专精 | `STAT_OPTION_INC_ATTACK_SPEED` | 当前未调用对应 ID | 未覆盖 | GameServer 会处理攻速专精，Go 侧应补对应 StatSpec 项。 |
| 71 | 装备基础修正 | 武器基础攻击 | `Right/Left m_DamageMin/Max` | `it.AttackMin/AttackMax` | 已覆盖 | 当前叠加到左右手攻击。 |
| 72 | 装备基础修正 | 武器追加攻击 | `PlusSpecial 80` | `it.AdditionAttack` | 已覆盖 | 需确认追加来源覆盖普通/卓越/套装。 |
| 73 | 装备基础修正 | 武器追加魔攻 | `PlusSpecial 81` | `it.AdditionMagicAttack` | 已覆盖 | 叠加到 magicAttackMin/Max。 |
| 74 | 装备基础修正 | 武器追加诅咒 | `PlusSpecial 113` | `it.AdditionCurseAttack` | 已覆盖 | 召唤书/诅咒装备需测试。 |
| 75 | 装备基础修正 | 防具基础防御 | `ItemDefense()` | `it.Defense` | 已覆盖 | 遍历防具/盾/翅膀槽位叠加。 |
| 76 | 装备基础修正 | 防具追加防御 | `PlusSpecial 83` | `it.AdditionDefense` | 已覆盖 | 叠加到 `Defense`。 |
| 77 | 装备基础修正 | 盾牌防御率 | `m_DefenseRate`、`PlusSpecial 82` | `it.SuccessfulBlocking/AdditionDefenseRate` | 已覆盖 | 当前所有防具槽都可能贡献，需确认盾牌槽规则。 |
| 78 | 装备基础修正 | 手套攻速 | `Gloves->m_AttackSpeed` | `glove.AttackSpeed` | 已覆盖 | 叠加到物理和魔法速度。 |
| 79 | 装备基础修正 | 装备耐久影响 | `m_CurrentDurabilityState` 影响攻击/魔攻 | 当前未体现耐久折损 | 未覆盖 | 需在道具耐久系统完成后把耐久折损接入公式。 |
| 80 | 装备基础修正 | 防具同类型等级套防 | `sameTypeCount/Level10..15` | `Player.calc` 同类型防御加成 | 部分覆盖 | Go 实现已有，但 Magumsa/RageFighter 槽位补偿逻辑需精确对齐。 |
| 81 | 卓越/幸运/套装修正 | 幸运一击率 | Lucky option | `it.Lucky -> criticalAttackRate += 5` | 已覆盖 | 和客户端幸运一击表现保持一致。 |
| 82 | 卓越/幸运/套装修正 | 卓越一击率 | ExcellentDamageSuccessRate | `ExcellentAttackRate -> +10` | 已覆盖 | 需和多件叠加上限对齐。 |
| 83 | 卓越/幸运/套装修正 | 卓越等级伤害 | Excellent damage by level | `ExcellentAttackLevel` | 已覆盖 | 当前按 `level/20` 增加攻击和魔攻。 |
| 84 | 卓越/幸运/套装修正 | 卓越百分比伤害 | ExcellentAttackPercent | `+2%` 攻击/魔攻 | 已覆盖 | 需确认 GameServer 百分比是否固定 2%。 |
| 85 | 卓越/幸运/套装修正 | 卓越攻速 | ExcellentAttackSpeed | `AttackSpeed/magicSpeed += 7` | 已覆盖 | 需确认叠加规则和上限。 |
| 86 | 卓越/幸运/套装修正 | 击杀恢复 HP/MP | MonsterDieGetLife/Mana | `monsterDieGetHP/MP` | 已覆盖 | `DieRecoverHPMP` 已读取该值。 |
| 87 | 卓越/幸运/套装修正 | 卓越防具 HP/MP | ExcellentDefenseHP/MP | `MaxHP/MaxMP += 4%` | 已覆盖 | 需测试多件叠加。 |
| 88 | 卓越/幸运/套装修正 | 卓越减伤/反伤 | DamageMinus/DamageReflect | `armorReduceDamage/armorReflectDamage` | 部分覆盖 | 减伤已用于战斗，反伤字段尚未参与结算。 |
| 89 | 卓越/幸运/套装修正 | 套装属性聚合 | `CalcSetItemOption/SetItemApply` | `calcSetItem(false)`、`addSetEffect` | 部分覆盖 | 需逐项对齐 SetEffectType 与 GameServer SetItem option。 |
| 90 | 卓越/幸运/套装修正 | 无视防御/双倍伤害 | SetOpIgnoreDefense/SetOpDoubleDamage | `ignoreDefenseRate/doubleDamageRate` | 部分覆盖 | 无视和双倍已用于战斗，来源项需要完整接入。 |
| 91 | 翅膀/宠物/戒指修正 | 翅膀增伤公式 | `Wings_CalcIncAttack` | `formula.Wings_CalcIncAttack` | 已覆盖 | `Player.calc` 将返回值转换为增伤百分比。 |
| 92 | 翅膀/宠物/戒指修正 | 翅膀吸收公式 | `Wings_CalcAbsorb` | `formula.Wings_CalcAbsorb` | 已覆盖 | `Player.calc` 将返回值转换为减伤百分比。 |
| 93 | 翅膀/宠物/戒指修正 | 翅膀追加属性 | `Wing->PlusSpecial` | `wing.AdditionAttack/Magic/Curse` | 已覆盖 | 基础追加已接入。 |
| 94 | 翅膀/宠物/戒指修正 | 翅膀卓越属性 | `m_WingExcOption` | `ExcellentWing2/3*` 字段 | 部分覆盖 | Go 已处理部分二代/三代翅膀卓越属性，需补完整列表。 |
| 95 | 翅膀/宠物/戒指修正 | 小恶魔/强化恶魔 | `Imp/Demon` 宠物配置 | `helper Code(13,1/64)` | 已覆盖 | 增伤和强化恶魔攻速已接入。 |
| 96 | 翅膀/宠物/戒指修正 | 守护天使/强化天使 | `Angel/SpiritAngel` | `helper Code(13,0/65)` | 已覆盖 | HP 加成和减伤已接入。 |
| 97 | 翅膀/宠物/戒指修正 | 黑王马/彩云兽/炎狼兽 | `gObjDarkHorse/gObjFenrir` | `p.pet` 对应 Code(13,3/4/37) | 部分覆盖 | 基础增减伤已接入，黑王马防御/技能等仍缺。 |
| 98 | 翅膀/宠物/戒指修正 | 特殊戒指攻击/防御 | WizardRing、PandaRing、SkeletonRing 等 | 当前仅部分 PetRing config 字段 | 未覆盖 | GameServer 特殊戒指修正较多，Go 侧未完整接入。 |
| 99 | 翅膀/宠物/戒指修正 | Muun/PeriodItem 修正 | `MuunSystem`、`PeriodItemEx` | 当前无公式接入 | 未覆盖 | 后续宠物/召唤系统完成后接入。 |
| 100 | 翅膀/宠物/戒指修正 | Pentagram 元素修正 | `PentagramSystem` | 当前无公式接入 | 未覆盖 | 元素伤害/防御独立成后续增强点。 |
| 101 | 技能与 Master 公式 | 骑士/魔剑技能倍率 | `RegularSkillCalc.lua` | `Knight_Gladiator_CalcSkillBonus` | 已覆盖 | 用于 `Object.getDamage` 的武器技能倍率。 |
| 102 | 技能与 Master 公式 | 穿刺技能倍率 | `ImpaleSkillCalc` | `formula.ImpaleSkillCalc` | 已覆盖 | 用于 `SkillIndexImpale`。 |
| 103 | 技能与 Master 公式 | 技能基础伤害 | `CMagicDamage::Get/SkillGet` | `Skill.DamageMin/DamageMax` | 部分覆盖 | 普通技能初始伤害由 SkillBase.Damage 派生。 |
| 104 | 技能与 Master 公式 | 技能消耗公式 | `CMagicDamage::SkillGetMana/SkillGetBP`、`gObjMagicManaUse/BPUse` | `Player.GetSkillMPAG` | 部分覆盖 | 消耗读取已有入口，但技能模块需完整验证。 |
| 105 | 技能与 Master 公式 | 技能距离公式 | `CMagicDamage::GetSkillDistance`、`gCheckSkillDistance` | `SkillBase.Distance` 等字段 | 部分覆盖 | 距离校验应在技能模块落地。 |
| 106 | 技能与 Master 公式 | 技能需求公式 | `SkillGetRequireEnergy/Class/Level` | `SkillBase.Req*`、`limitUseItem` 等 | 部分覆盖 | 需要统一技能学习/使用条件校验。 |
| 107 | 技能与 Master 公式 | Master 技能获得 | `CGReqGetMasterLevelSkill` | `Skills.GetMaster` | 部分覆盖 | Go 有基础加点逻辑，但前置、点数、树位置需完整测试。 |
| 108 | 技能与 Master 公式 | Master 技能伤害 | `GetSkillAttackDamage` | `Skills.getMasterSkillDamage` | 部分覆盖 | 当前函数存在但未启用，需要接入技能伤害计算。 |
| 109 | 技能与 Master 公式 | Master 被动加成 | `CalcPassiveSkillData`、`ApplyMLSkillItemOption` | 当前基本未接入 `Player.calc` | 未覆盖 | 被动攻击、防御、属性、装备类型加成是重要缺口。 |
| 110 | 技能与 Master 公式 | Combo/特殊技能倍率 | `ComboAttack`、大量 `MLS_*` | 当前仅少数技能走倍率 | 未覆盖 | Combo、范围技能、特殊技能应归技能系统逐步接入公式层。 |
| 111 | 经验/等级/战斗结算公式 | 普通经验表 | `gLevelExperience`、`ExpCalc` | `SetExpTable_Normal`、`ExperienceTable` | 已覆盖 | 启动时按最大普通等级生成。 |
| 112 | 经验/等级/战斗结算公式 | 大师经验表 | `CMasterLevelSystem::SetExpTable` | `SetExpTable_Master`、`MasterExperienceTable` | 已覆盖 | 启动时按最大大师等级生成。 |
| 113 | 经验/等级/战斗结算公式 | 经验倍率配置 | `CExpManager::LoadScript` | `ExpManager.init` 读取 `IGC_ExpSystem.xml` | 部分覆盖 | 当前只使用 StaticExp，动态区间未实现。 |
| 114 | 经验/等级/战斗结算公式 | 动态经验倍率 | `CExpManager::GetExpMultiplier` 动态范围 | 当前无 DynamicExpRangeList 计算 | 未覆盖 | 需要按 reset/level/master 选择动态倍率。 |
| 115 | 经验/等级/战斗结算公式 | 单人怪物经验 | `gObjMonsterExpSingle` | `Object.DieGiveExperience` | 部分覆盖 | Go 已有基础怪物等级公式，但未按伤害占比分配。 |
| 116 | 经验/等级/战斗结算公式 | 伤害占比经验 | `exp = dmg * exp / tot_dmg` | 当前击杀者获得全部经验 | 未覆盖 | 需要记录怪物受击列表后按伤害分配。 |
| 117 | 经验/等级/战斗结算公式 | 组队经验 | `gObjExpParty/gObjExpPartyRenewal` | 当前无组队经验 | 未覆盖 | 社交/队伍系统完成后接入。 |
| 118 | 经验/等级/战斗结算公式 | 地图/活动/VIP 经验加成 | `MapAttr/BonusEvent/Vip/Crywolf/Doppelganger` | `MapManager.GetExpBonus`、`ExpManager.Normal/Master` | 部分覆盖 | 当前有地图和静态倍率，缺 VIP/活动/惩罚/副本特化。 |
| 119 | 经验/等级/战斗结算公式 | 升级点数公式 | `gObjLevelUpPointAdd`、`LevelPoint5/7` | `Player.LevelUp` | 部分覆盖 | 普通升级点已实现，属性点上限/任务点/水果点仍缺。 |
| 120 | 经验/等级/战斗结算公式 | 战斗伤害结算公式 | `ObjBaseAttack`、`gObjAttack/gObjLifeCheck` | `Object.attack/getDamage/getDefense` | 部分覆盖 | 已有 miss、防御、暴击、卓越、减伤、双倍；缺 SD 分伤、反射、元素、最低伤害等。 |
