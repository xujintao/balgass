# 10. 技能系统

本模块覆盖普通技能、大师技能、技能学习、技能列表、技能使用、资源消耗、目标校验、范围命中、冷却防挂、Buff/状态、Combo、协议与测试。战斗伤害主体归对象系统和公式系统，本模块只记录技能选择、技能规则、技能效果入口和协议边界。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 技能表与基础定义 | 技能表 XML 加载 | `CMagicDamage::Set`、SkillList 配置 | `SkillManager.init` 读取 `Skills/IGC_SkillList.xml` | 已覆盖 | Go 侧已加载技能表并构建 `skillTable`。 |
| 2 | 技能表与基础定义 | SkillBase 字段建模 | `MAGIC_DAMAGE` 字段族 | `SkillBase` | 部分覆盖 | 字段已包含需求、消耗、距离、延迟、类型、Buff 等，但使用路径未完整校验。 |
| 3 | 技能表与基础定义 | 技能 ID 映射 | `MagicDamageC.Get/SkillGet` | `skillTable map[int]*SkillBase` | 已覆盖 | Go 侧按技能 Index 查询技能基础表。 |
| 4 | 技能表与基础定义 | 职业需求映射 | `SkillGetRequireClass` | `ReqClass`、职业字段映射 | 部分覆盖 | Go 侧映射到 7 个职业，GrowLancer 字段仍未从技能表映射。 |
| 5 | 技能表与基础定义 | 等级需求字段 | `SkillGetRequireLevel` | `ReqLevel` | 未覆盖 | 字段已加载，但技能学习和使用时未校验角色等级。 |
| 6 | 技能表与基础定义 | 属性需求字段 | `GetRequireStrength/Dexterity/Energy/Leadership` | `ReqStrength`、`ReqDexterity`、`ReqEnergy`、`ReqCommand` | 未覆盖 | 字段已加载，但学习和使用路径未校验。 |
| 7 | 技能表与基础定义 | 技能类型字段 | `GetSkillType` | `Type` | 部分覆盖 | 字段存在，但 Go 侧技能分发主要靠硬编码技能 ID。 |
| 8 | 技能表与基础定义 | 使用类型字段 | `GetSkillUseType` | `UseType` | 部分覆盖 | 登录技能列表会过滤 `UseType==3`，其余使用语义未完整落地。 |
| 9 | 技能表与基础定义 | Brand 技能链字段 | `GetBrandOfSkill`、`CheckBrandOfSkill` | `Brand`、`getMasterSkillDamage` | 部分覆盖 | Go 侧仅在大师技能伤害推导中使用，普通 brand 规则未完整实现。 |
| 10 | 技能表与基础定义 | BuffIndex 字段 | `byBufIndex` | `BuffIndex` | 未覆盖 | 字段已加载，但 Buff 系统和状态下发未实现。 |
| 11 | 技能实例与持久化 | 技能实例结构 | `CMagicInf` | `skill.Skill` | 部分覆盖 | Go 侧实例包含 Index、Level、DamageMin/Max、UIIndex、CurValue、NextValue。 |
| 12 | 技能实例与持久化 | 空技能实例 | 普攻/空 magic 处理 | `Skill0` | 已覆盖 | 普攻时用 `Skill0` 作为空技能。 |
| 13 | 技能实例与持久化 | 技能协议编码 | `MagicByteConvert` | `Skill.Marshal` | 部分覆盖 | Go 侧编码 Index 和 Level，需要协议 round-trip 测试。 |
| 14 | 技能实例与持久化 | 技能排序 | GameServer 固定 Magic 数组顺序 | `SortedSkillSlice` | 已覆盖 | Go 侧 JSON 和列表下发按 Index 排序。 |
| 15 | 技能实例与持久化 | 技能 JSON 保存 | DB magic byte array | `Skills.MarshalJSON`、`Value` | 部分覆盖 | Go 侧用 JSON 保存技能，需考虑坏数据和版本迁移。 |
| 16 | 技能实例与持久化 | 技能 JSON 加载 | DB magic decode | `Skills.UnmarshalJSON`、`Scan` | 部分覆盖 | 未知技能会记录错误并跳过，需补测试。 |
| 17 | 技能实例与持久化 | 技能基础表回填 | `CMagicInf::Set` | `FillSkillData` | 部分覆盖 | 登录后给技能回填 `SkillBase`、Damage、UIIndex，需覆盖所有大师技能。 |
| 18 | 技能实例与持久化 | 普通技能和大师技能区分 | `CheckMasterLevelSkill` | `Index < 300` 判断 | 部分覆盖 | Go 侧用索引粗分，需确认所有版本技能 ID 边界。 |
| 19 | 技能实例与持久化 | 技能等级字段 | `CMagicInf.m_Level` | `Skill.Level` | 部分覆盖 | 普通技能多为 0，大师技能用 Level 表示投入点数。 |
| 20 | 技能实例与持久化 | 技能坏数据修复 | GameServer magic decode 校验 | 当前仅跳过未知技能 | 未覆盖 | 需要清理重复技能、越界等级、缺基础表和非法大师技能。 |
| 21 | 技能学习与遗忘 | 技能书使用入口 | `gObjMagicAdd`、技能书使用 | `Object.UseItem` KindASkill 分支 | 部分覆盖 | Go 侧技能书可触发学习，但规则校验不完整。 |
| 22 | 技能学习与遗忘 | 普通技能学习 | `gObjMagicAdd` | `Object.LearnSkill`、`Skills.Get` | 部分覆盖 | 只防重复和查基础表，未校验职业、等级、属性。 |
| 23 | 技能学习与遗忘 | 重复学习处理 | `MagicSearch` | `LearnSkill` 检查已有技能 | 已覆盖 | 重复学习会失败并记录日志。 |
| 24 | 技能学习与遗忘 | 技能书消耗 | GameServer 使用成功后删物品 | 当前物品使用路径需核对 | 部分覆盖 | 需确保学习成功才消耗技能书，失败不消耗。 |
| 25 | 技能学习与遗忘 | 职业校验 | `SkillGetRequireClass` | 当前未校验 `ReqClass` | 未覆盖 | 技能书学习和技能使用都需要校验职业/转职等级。 |
| 26 | 技能学习与遗忘 | 等级校验 | `SkillGetRequireLevel` | 当前未校验 `ReqLevel` | 未覆盖 | 未达等级不能学习和使用技能。 |
| 27 | 技能学习与遗忘 | 属性校验 | `SkillGetRequireEnergy` 等 | 当前未校验属性需求 | 未覆盖 | 力量、敏捷、体力、智力、统率需求需要接入角色属性。 |
| 28 | 技能学习与遗忘 | 装备技能自动学习 | `MLS_WeaponSkillAdd`、武器技能 | `EquipmentChanged` 学习主副手 `SkillIndex` | 部分覆盖 | Go 侧已有主副手装备技能增删，但需核对所有武器技能和大师武器技能。 |
| 29 | 技能学习与遗忘 | 装备技能自动遗忘 | `MLS_WeaponSkillDel` | `EquipmentChanged` 遗忘旧 ItemSkill | 部分覆盖 | 需要防止遗忘玩家通过技能书永久学习的同 ID 技能。 |
| 30 | 技能学习与遗忘 | 技能学习保存 | Char save magic list | `Player.Save` 保存 `Skills` | 部分覆盖 | Go 侧会保存技能 JSON，需明确临时装备技能是否持久化。 |
| 31 | 技能列表与协议 | 登录技能加载 | Char load magic list | `p.Skills = c.Skills`、`FillSkillData` | 已覆盖 | 登录时从角色数据加载并回填技能数据。 |
| 32 | 技能列表与协议 | 普通技能列表下发 | `GCMagicListMultiSend` | `pushSkillList`、`MsgSkillListReply` | 部分覆盖 | Go 侧下发 active 技能，需校验客户端字段和排序。 |
| 33 | 技能列表与协议 | 技能列表过滤 | GameServer active/passive 区分 | `ForEachActiveSkill` 跳过 `UseType==3` | 部分覆盖 | 需确认所有被动/隐藏技能不应出现在普通技能栏。 |
| 34 | 技能列表与协议 | 单技能新增协议 | `GCMagicListOneSend` | `MsgSkillOneReply{Flag:-2}` | 部分覆盖 | 装备技能和技能书学习都应下发新增。 |
| 35 | 技能列表与协议 | 单技能删除协议 | `GCMagicListOneDelSend` | `MsgSkillOneReply{Flag:-1}` | 部分覆盖 | 装备卸下时下发删除，需验证客户端行为。 |
| 36 | 技能列表与协议 | 技能列表 Marshal | `MagicByteConvert` | `MsgSkillListReply.Marshal` | 部分覆盖 | 当前写 count、padding、slot 和技能 bytes，需要抓包/客户端测试。 |
| 37 | 技能列表与协议 | 单技能 Marshal | 单技能协议 | `MsgSkillOneReply.Marshal` | 部分覆盖 | Flag 使用 -2/-1，需确认 Season 9 客户端兼容。 |
| 38 | 技能列表与协议 | 大师技能列表下发 | `CGReqGetMasterLevelSkillTree` | `pushMasterSkillList`、`MsgMasterSkillListReply` | 部分覆盖 | Go 侧下发 UIIndex/Level/CurValue/NextValue，但数值多为 0。 |
| 39 | 技能列表与协议 | 大师技能学习协议 | `CGReqGetMasterLevelSkill` | `LearnMasterSkill`、`MsgLearnMasterSkillReply` | 部分覆盖 | Go 侧已有请求和回复结构，需补完整错误码。 |
| 40 | 技能列表与协议 | 技能协议 opcode 映射 | GameServer protocol dispatch | `handle/c1c2.go` 0x19、0xF311、0xF352、0xF353 | 部分覆盖 | 需把所有技能协议入口连接到稳定业务和失败回复。 |
| 41 | 技能使用入口与消耗 | UseSkill 请求解析 | `CGMagicAttack` | `MsgUseSkill.Unmarshal` | 已覆盖 | Go 侧解析 target 和 skill 的高低字节。 |
| 42 | 技能使用入口与消耗 | UseSkill 回复协议 | `GCMagicAttackNumberSend` | `MsgUseSkillReply` | 部分覆盖 | 成功时 target 置 `0x8000`，需核对失败语义。 |
| 43 | 技能使用入口与消耗 | 目标对象获取 | `gObj[aTargetIndex]` | `ObjectManager.objects[msg.Target]` | 部分覆盖 | 目标不存在只记录日志，没有失败包。 |
| 44 | 技能使用入口与消耗 | 已学技能校验 | `MagicSearch` | `obj.Skills[msg.Skill]` | 已覆盖 | 未学习技能直接返回。 |
| 45 | 技能使用入口与消耗 | MP 消耗计算 | `CObjUseSkill::GetUseMana` | `GetSkillMPAG` 返回 `ManaUsage` | 部分覆盖 | 未考虑装备、Buff、大师技能、特殊技能附加消耗。 |
| 46 | 技能使用入口与消耗 | AG/BP 消耗计算 | `CObjUseSkill::GetUseBP` | `GetSkillMPAG` 返回 `BPUsage` | 部分覆盖 | 只读取基础表 BPUsage，未处理特殊技能规则。 |
| 47 | 技能使用入口与消耗 | 资源不足处理 | GameServer 消耗失败返回 | `if obj.MP < mp || obj.AG < ag { return }` | 部分覆盖 | 当前无失败回复和审计日志。 |
| 48 | 技能使用入口与消耗 | 消耗扣除时机 | `UseSkill` 消耗逻辑 | 当前先执行 `canUseSkill` 后扣 MP/AG | 需修正 | 如果技能执行失败或部分命中，资源扣除时机需要明确。 |
| 49 | 技能使用入口与消耗 | 资源同步 | `GCManaSend` | `PushMPAG` | 已覆盖 | 技能使用后会同步 MP/AG。 |
| 50 | 技能使用入口与消耗 | 道具/箭矢消耗 | `DecreaseArrow`、特殊消耗 | 当前无技能道具消耗 | 未覆盖 | 弓箭、召唤、特殊技能需要道具/箭矢消耗规则。 |
| 51 | 目标校验与安全限制 | 目标存在校验 | `RunningSkill` 目标检查 | `UseSkill` 目标 nil 检查 | 部分覆盖 | 需返回失败包并处理目标已销毁。 |
| 52 | 目标校验与安全限制 | 同地图校验 | GameServer 地图检查 | 当前未显式校验 | 未覆盖 | 跨地图对象不能被技能命中。 |
| 53 | 目标校验与安全限制 | 距离校验 | `GetSkillDistance` | 当前未按 `Distance` 校验 | 未覆盖 | 技能使用必须校验技能表距离。 |
| 54 | 目标校验与安全限制 | 视野校验 | viewport/target check | 当前未校验视野 | 未覆盖 | 防止客户端攻击不可见目标。 |
| 55 | 目标校验与安全限制 | 死亡状态校验 | `lpObj->Live` 等 | 当前缺少统一校验 | 未覆盖 | 施法者或目标死亡时不能执行多数技能。 |
| 56 | 目标校验与安全限制 | 安全区限制 | `SkillSafeZoneUse` | 当前无安全区技能规则 | 未覆盖 | 安全区可用技能、攻击技能、Buff 技能需要分开判断。 |
| 57 | 目标校验与安全限制 | PVP/PK 限制 | `PkCheck`、地图 PVP 规则 | 当前主要复用 attack 路径 | 部分覆盖 | 技能使用前应判断是否允许攻击玩家。 |
| 58 | 目标校验与安全限制 | 友方/敌方限制 | Buff/攻击目标类型 | 当前未区分友敌技能 | 未覆盖 | 治疗、Buff、召唤、攻击技能需要不同目标类型。 |
| 59 | 目标校验与安全限制 | 施法状态限制 | stun/sleep/cloak/interface state | 当前无统一状态限制 | 未覆盖 | 交易、死亡、眩晕、睡眠、隐身等状态影响技能使用。 |
| 60 | 目标校验与安全限制 | 技能启用开关 | `CObjUseSkill::EnableSkill` | 当前无禁用表 | 未覆盖 | 需要支持配置禁用技能或地图禁用技能。 |
| 61 | 范围与命中形状 | 单体技能命中 | `UseSkill(aIndex,aTargetIndex,...)` | `canUseSkill` 单体 attack | 部分覆盖 | 多个基础攻击技能已按单体处理。 |
| 62 | 范围与命中形状 | 圆形范围技能 | `SkillAreaMonsterAttack` | 当前无通用圆形范围 | 未覆盖 | Evil Spirit、HellFire 等范围技能需要统一圆形/距离筛选。 |
| 63 | 范围与命中形状 | 扇形范围技能 | `SkillFrustrum` | `CreateSkillFrustum`、`CheckSkillFrustum` | 部分覆盖 | DeathStab 已有扇形雏形，需扩展到其他技能。 |
| 64 | 范围与命中形状 | 直线技能路径 | `GetTargetLinePath` | 当前无直线路径 | 未覆盖 | PowerSlash、穿透类技能需要路径和阻挡校验。 |
| 65 | 范围与命中形状 | 矩形/HitBox 技能 | `SkillHitBox::HitCheck` | 当前无 SkillHitBox 数据 | 未覆盖 | 需要加载 hitbox 配置并按方向判断命中。 |
| 66 | 范围与命中形状 | 技能角度计算 | `GetAngle`、`SkillAreaCheck` | `getAngle` | 部分覆盖 | Go 侧可算目标角度，需统一成通用范围工具。 |
| 67 | 范围与命中形状 | 范围对象筛选 | `Viewport` 遍历 | `ForEachViewportObject` | 部分覆盖 | DeathStab 已遍历视野对象，需补玩家/PVP/友方过滤。 |
| 68 | 范围与命中形状 | 地图阻挡校验 | path/map attr checks | 当前无技能阻挡校验 | 未覆盖 | 范围和直线技能应考虑墙体、安全区和不可达格。 |
| 69 | 范围与命中形状 | 多目标命中上限 | GameServer 多目标限制 | 当前无上限 | 未覆盖 | 范围技能需限制最大命中数和重复命中。 |
| 70 | 范围与命中形状 | 范围技能广播 | GameServer 技能效果广播 | 当前主要 `UseSkillReply` 和 attack 包 | 部分覆盖 | 需要按技能类型广播施法、命中和效果。 |
| 71 | 冷却、防挂与使用频率 | 技能 Delay 字段 | `MagicDamage.GetDelayTime` | `SkillBase.Delay` | 未覆盖 | 字段已加载，但 UseSkill 未使用。 |
| 72 | 冷却、防挂与使用频率 | SkillDelay 检查 | `CSkillDelay::Check` | 当前无 SkillDelay | 未覆盖 | 需要判断技能是否走延迟控制。 |
| 73 | 冷却、防挂与使用频率 | SkillUseTime 配置 | `CSkillUseTime::LoadFile` | 当前无使用时间配置 | 未覆盖 | 需要加载技能最小使用间隔配置。 |
| 74 | 冷却、防挂与使用频率 | 单技能冷却 | `CheckSkillTime` | 当前无 per-skill 时间戳 | 未覆盖 | 每个技能需要记录上次使用时间。 |
| 75 | 冷却、防挂与使用频率 | 公共冷却 | GameServer 速度限制 | 当前无公共 CD | 未覆盖 | 普攻、技能、魔法应受攻击速度/魔法速度限制。 |
| 76 | 冷却、防挂与使用频率 | SpeedHackCheck | `CObjUseSkill::SpeedHackCheck` | 当前无技能防加速 | 未覆盖 | 调用 `27-security.md` 基于攻击速度、魔法速度和包频率检测异常。 |
| 77 | 冷却、防挂与使用频率 | MultiAttackHackCheck | `MultiAttackHackCheck` | 当前无多重攻击检测 | 未覆盖 | 调用 `27-security.md` 限制范围和多段技能的异常多目标/高频攻击。 |
| 78 | 冷却、防挂与使用频率 | 延迟攻击消息 | GameServer delayed skill hit | `AddDelayMsg` 用于击退/延迟 | 部分覆盖 | Go 侧已有少量延迟消息，需形成技能延迟体系。 |
| 79 | 冷却、防挂与使用频率 | 冷却失败回复 | GameServer 拒绝/日志 | 当前直接 return | 未覆盖 | CD、速度、防挂失败应稳定返回或记录。 |
| 80 | 冷却、防挂与使用频率 | 技能异常审计 | hack log | 当前无技能审计 | 未覆盖 | 记录高频技能、非法目标、非法距离、非法状态。 |
| 81 | 普通技能效果分发 | 技能总分发入口 | `CObjUseSkill::RunningSkill` | `canUseSkill` | 部分覆盖 | Go 侧只有少量硬编码分支，需要扩展为可维护分发。 |
| 82 | 普通技能效果分发 | 魔法攻击技能 | `ObjUseSkill` 魔法技能族 | Poison/Meteorite/Lightning/FireBall 等 | 部分覆盖 | 当前这些技能多数按普通 attack 处理，缺元素/Buff/Debuff。 |
| 83 | 普通技能效果分发 | 物理武器技能 | Cyclone/Slash/Lunge 等 | weapon skill 分支 | 部分覆盖 | 已按物理攻击和技能加成处理，缺技能特效和特殊规则。 |
| 84 | 普通技能效果分发 | 防御姿态技能 | `SkillDefense` | `SkillIndexDefense` 下发 Action | 部分覆盖 | 当前只广播动作，未实现防御状态/减伤。 |
| 85 | 普通技能效果分发 | 治疗技能 | `SkillHealing` | 当前无 Heal 效果 | 未覆盖 | 需要治疗目标、消耗、上限、友方校验和回复协议。 |
| 86 | 普通技能效果分发 | 攻击/防御 Buff | `SkillAttack`、`WizardMagicDefense` | 当前无 Buff 效果 | 未覆盖 | GreaterAttack、GreaterDefense、SoulBarrier 等需进入 Buff 系统。 |
| 87 | 普通技能效果分发 | Debuff 技能 | Poison/Ice/Sleep/Weakness 等 | 当前无持续 Debuff | 未覆盖 | 中毒、冰冻、睡眠、虚弱等需要状态和持续时间。 |
| 88 | 普通技能效果分发 | 召唤技能 | `SkillMonsterCall`、`SkillSummon` | 当前无召唤技能 | 未覆盖 | 召唤怪、星云召唤、召唤队友等需要对象和地图联动。 |
| 89 | 普通技能效果分发 | 位移技能 | Teleport、Rush、Charge 等 | 当前无位移技能 | 未覆盖 | 位移技能需校验地图格、距离、阻挡和广播。 |
| 90 | 普通技能效果分发 | 特殊技能后处理 | `SpecificSkillAdditionTreat` | 当前无统一后处理 | 未覆盖 | 火舞旋风爆炸、连锁、电击、反射等特殊效果需统一挂接。 |
| 91 | 重点普通技能族 | Twisting/Rageful/DeathStab | `SkillWheel`、`SkillBlowOfFury`、`SkillKnightBlow` | Twisting/Rageful 普通 attack，DeathStab 扇形 | 部分覆盖 | 需补完整范围、连击、技能倍率和特效。 |
| 92 | 重点普通技能族 | PowerSlash/FireSlash | `SkillPowerSlash`、`SkillMaGum...` | `SkillIndexPowerSlash` 基础处理 | 部分覆盖 | 需补直线/角度、多目标和魔剑士特殊规则。 |
| 93 | 重点普通技能族 | FireScream | `SkillFireScream`、`FireScreamExplosionAttack` | 当前无完整实现 | 未覆盖 | 需要主目标、爆炸范围、附加伤害和多段处理。 |
| 94 | 重点普通技能族 | InfinityArrow | `SkillInfinityArrow`、升级自动学习 | 当前无完整实现 | 未覆盖 | 需要 Buff、箭矢消耗减免和自动学习。 |
| 95 | 重点普通技能族 | Summoner 技能族 | Sleep/DrainLife/Weakness/Innovation | 当前无完整实现 | 未覆盖 | 召唤师诅咒、睡眠、吸血、弱化需独立效果。 |
| 96 | 重点普通技能族 | RageFighter 技能族 | MonkBuff、DarkSide、Barrage | 当前无完整实现 | 未覆盖 | RF 技能需要 Buff、连击、范围和特殊动作。 |
| 97 | 重点普通技能族 | GrowLancer 技能族 | Breche、ShiningPeak、SpinStep | 当前无完整实现 | 未覆盖 | GL 技能表映射和主动技能均需补齐。 |
| 98 | 重点普通技能族 | Fenrir 技能 | `SkillFenrirAttack` | 当前无完整实现 | 未覆盖 | 炎狼兽技能需处理坐骑、范围和耐久。 |
| 99 | 重点普通技能族 | DarkHorse/DarkRaven 技能 | `SkillDarkHorseAttack` | 当前无完整实现 | 未覆盖 | 黑王马/黑鹰技能需要宠物系统配合。 |
| 100 | 重点普通技能族 | Castle/事件技能 | 攻城、IllusionTemple 等特殊技能 | handler 有事件技能入口痕迹 | 未覆盖 | 事件/攻城技能只记录边界，具体事件规则归事件模块。 |
| 101 | 大师技能系统 | 大师技能树加载 | `MasterLevelSkillTreeSystem::Load` | `IGC_MasterSkillTree.xml` 加载到数组 | 部分覆盖 | Go 侧已加载树，但未完整校验 rank/tree。 |
| 102 | 大师技能系统 | 大师技能值表 | `AddToValueTable`、Lua value table | `masterSkillValueTable` 未填充 | 未覆盖 | 需要从 Lua/配置填充 CurValue/NextValue 和伤害值。 |
| 103 | 大师技能系统 | 大师技能查找 | `GetMasterSkillUIIndex`、`GetRequireMLPoint` | `getMasterSkillBase` | 部分覆盖 | Go 侧按 class/tree/rank/pos 扫描 MagicNumber。 |
| 104 | 大师技能系统 | 大师技能学习入口 | `CGReqGetMasterLevelSkill` | `Player.LearnMasterSkill` | 部分覆盖 | Go 侧有入口，但错误码和校验不完整。 |
| 105 | 大师技能系统 | 大师等级/点数前置 | `CheckSkillCondition` | `masterLevel > 0 && masterPoint > 0` | 部分覆盖 | 需校验大师等级、剩余点数、职业和技能树条件。 |
| 106 | 大师技能系统 | 父技能前置校验 | `CheckPreviousRankSkill`、ParentSkill1/2 | `GetMaster` 只检查父技能存在于表 | 未覆盖 | Go 侧未检查玩家是否已学习父技能和父技能等级。 |
| 107 | 大师技能系统 | 最大点数校验 | `GetMaxPointOfMasterSkill` | hardcoded `ss.Level+ReqMinPoint > 20` | 部分覆盖 | 需要使用每个技能 `MaxPoint`，不是固定 20。 |
| 108 | 大师技能系统 | 点数扣除 | `MasterPoint -= ReqPoint` | `p.masterPoint -= point` | 已覆盖 | 需保证失败不扣点，并保存角色。 |
| 109 | 大师技能系统 | 大师被动效果 | `CalcPassiveSkillData` | 当前无完整被动应用 | 未覆盖 | 需要在角色重算时应用大师被动属性。 |
| 110 | 大师技能系统 | 大师主动技能分发 | `RunningSkill_MLS`、`UseMasterSkill` | 当前无大师主动技能分发 | 未覆盖 | 大师技能使用时应路由到普通技能强化或专属效果。 |
| 111 | Combo、Buff 与测试 | Combo 技能位置 | `ComboAttack::GetSkillPos` | 当前无 Combo 结构 | 未覆盖 | 需要识别 Combo 起手、第二段、终结技能。 |
| 112 | Combo、Buff 与测试 | Combo 序列检查 | `ComboAttack::CheckCombo` | 当前无 Combo 检查 | 未覆盖 | 需要记录技能序列和目标，判定 Combo 成功。 |
| 113 | Combo、Buff 与测试 | Combo 时间窗口 | GameServer combo time | 当前无时间窗口 | 未覆盖 | Combo 技能必须在限定时间内连续释放。 |
| 114 | Combo、Buff 与测试 | Combo 伤害入口 | `ObjAttack` combo damage | 当前无 Combo 伤害 | 未覆盖 | Combo 成功后触发额外伤害和协议广播。 |
| 115 | Combo、Buff 与测试 | BuffEffect 加载 | `BuffScriptLoader`、`BuffEffect` | 当前无 Buff 表 | 未覆盖 | 需要加载 Buff 类型、持续时间、效果值和图标。 |
| 116 | Combo、Buff 与测试 | BuffEffectSlot 状态 | `BuffEffectSlot` | 当前无 Buff slot | 未覆盖 | 对象需要维护当前 Buff/Debuff 状态、过期和叠加。 |
| 117 | Combo、Buff 与测试 | SkillAdditionInfo | `SkillAdditionInfo::Load` | 当前无技能附加配置 | 未覆盖 | 需要加载生命之光、无限箭等技能附加参数。 |
| 118 | Combo、Buff 与测试 | 技能协议回归 | GameServer 技能协议 | `UseSkillReply`、`SkillListReply`、`MasterSkillListReply` | 未覆盖 | 覆盖技能列表、学习、遗忘、施法成功/失败协议。 |
| 119 | Combo、Buff 与测试 | 战斗联动回归 | `ObjAttack`、`ObjUseSkill` | `attack`、`getDamage`、`canUseSkill` | 未覆盖 | 覆盖技能伤害、范围、死亡、经验、掉落联动。 |
| 120 | Combo、Buff 与测试 | 技能边界和压测 | 防挂/异常日志 | 待实现测试集 | 未覆盖 | 覆盖非法技能、非法目标、高频施法、多目标范围、并发死亡。 |
