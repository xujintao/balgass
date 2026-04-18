# 13. Buff 系统

本模块覆盖 `GameServer` 中 Buff 定义、BuffSlot、Effect、持续时间、增删查清、数值应用、Debuff、技能来源、道具/期限来源、视野同步、经验加成和协议下发。当前 `server-game` 只有 `EffectList`、`MaxBuffEffect`、部分宠物/戒指 buff 配置和攻击效果包体，尚未形成独立 Buff 生命周期和状态同步系统。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | Buff 定义与配置 | BuffScriptLoader 初始化 | `CBuffScriptLoader::Initialize` | 暂无 Buff 配置管理器 | 未覆盖 | 初始化 Buff 定义表、PeriodBuff 表和索引结构。 |
| 2 | Buff 定义与配置 | BuffScriptLoader 加载 | `CBuffScriptLoader::Load` | 暂无 Buff XML/TXT 加载 | 未覆盖 | 加载 BuffIndex、名称、类型、通知、清理类型和描述；通用脚本/配置加载基础设施归 `26-script.md`。 |
| 3 | Buff 定义与配置 | AddBuffEffectData | `CBuffScriptLoader::AddBuffEffectData` | 暂无 Buff 定义注册 | 未覆盖 | 将单个 Buff 定义注册到 Go 侧 Buff 管理器。 |
| 4 | Buff 定义与配置 | CheckVaildBuffEffect | `CBuffScriptLoader::CheckVaildBuffEffect` | 暂无 BuffIndex 合法性校验 | 未覆盖 | 校验 BuffIndex 是否存在，防止非法协议或脚本写入。 |
| 5 | Buff 定义与配置 | AddPeriodBuffEffectInfo | `CBuffScriptLoader::AddPeriodBuffEffectInfo` | 暂无期限 Buff 定义 | 未覆盖 | 注册期限道具 Buff 的效果类型、数值和默认持续时间。 |
| 6 | Buff 定义与配置 | BuffEffectType | `BuffEffectType` 字段 | 暂无 Buff 类型枚举 | 未覆盖 | 定义普通 Buff、Debuff、期限 Buff、事件 Buff 等类型。 |
| 7 | Buff 定义与配置 | BuffType | `btBuffType` | 暂无 Buff 作用域分类 | 未覆盖 | 区分正面、负面、状态、变身、图标类 Buff。 |
| 8 | Buff 定义与配置 | NoticeType | `btNoticeType` | 暂无 Buff 通知策略 | 未覆盖 | 决定 Buff 增删是否下发本人、视野或全局。 |
| 9 | Buff 定义与配置 | ClearType | `btClearType`、`eBuffClearType` | 暂无 Buff 清理类型 | 未覆盖 | 定义登出、死亡、地图切换、非 PCS、全部清理等策略。 |
| 10 | Buff 定义与配置 | ItemType/ItemIndex | `btItemType`、`wItemIndex` | 宠物/戒指配置存在但未统一 | 未覆盖 | 建立物品到 Buff 的映射，用于期限道具和消耗道具触发。 |
| 11 | Effect 基础模型 | CEffect 构造 | `CEffect::CEffect` | `EffectList` 只有结构体 | 部分覆盖 | 初始化单个效果实例的默认空状态。 |
| 12 | Effect 基础模型 | IsEffect | `CEffect::IsEffect` | 暂无 Effect 有效判断 | 未覆盖 | 判断效果槽是否被占用。 |
| 13 | Effect 基础模型 | Set 简化写入 | `CEffect::Set(int effect)` | 暂无 Effect 写入方法 | 未覆盖 | 支持只按 BuffIndex 标记一个效果。 |
| 14 | Effect 基础模型 | Set 完整写入 | `CEffect::Set(effect,item,option,value,time,tick)` | `EffectList` 字段不完整 | 未覆盖 | 写入 BuffIndex、物品、效果类型、效果值、持续时间和开始 tick。 |
| 15 | Effect 基础模型 | Clear | `CEffect::Clear` | 暂无 Effect 清理方法 | 未覆盖 | 清空单个 Buff 效果槽。 |
| 16 | Effect 基础模型 | EffectCategory | `EFFECTLIST::EffectCategory` | `EffectList.EffectCategory` 已有 | 部分覆盖 | 明确分类值语义并参与清理/叠加规则。 |
| 17 | Effect 基础模型 | EffectType1/2 | `EffectType1/EffectType2` | `EffectList.EffectType1/2` 已有 | 部分覆盖 | 支持一个 Buff 携带两个效果类型。 |
| 18 | Effect 基础模型 | EffectValue1/2 | `EffectValue1/EffectValue2` | `EffectList.EffectValue1/2` 已有 | 部分覆盖 | 支持两个效果数值并统一百分比/固定值语义。 |
| 19 | Effect 基础模型 | EffectSetTime | `EffectSetTime` | `EffectList.EffectSetTime` 已有 | 部分覆盖 | 记录设置时间，用于过期和剩余时间计算。 |
| 20 | Effect 基础模型 | EffectDuration | `EffectDuration` | `EffectList.EffectDuration` 已有 | 部分覆盖 | 支持正数限时、负数永久和 0 特殊语义。 |
| 21 | BuffSlot 生命周期 | CBuffEffectSlot 构造 | `CBuffEffectSlot::CBuffEffectSlot` | `Object` Buff 字段仍注释 | 未覆盖 | 初始化对象 BuffSlot 容器和计数。 |
| 22 | BuffSlot 生命周期 | SetEffect | `CBuffEffectSlot::SetEffect` | 暂无对象 Buff 设置 | 未覆盖 | 向对象 BuffSlot 写入 Buff，处理空槽、重复和上限。 |
| 23 | BuffSlot 生命周期 | RemoveEffect | `CBuffEffectSlot::RemoveEffect` | 暂无对象 Buff 移除 | 未覆盖 | 按 BuffIndex 删除对象身上的 Buff。 |
| 24 | BuffSlot 生命周期 | CheckUsedEffect | `CBuffEffectSlot::CheckUsedEffect` | 暂无对象 Buff 查询 | 未覆盖 | 查询对象是否已有指定 Buff。 |
| 25 | BuffSlot 生命周期 | RemoveBuffVariable | `CBuffEffectSlot::RemoveBuffVariable` | 暂无 Buff 变量清理 | 未覆盖 | Buff 删除时同步清理对象上的派生状态变量。 |
| 26 | BuffSlot 生命周期 | ClearEffect | `CBuffEffectSlot::ClearEffect` | 暂无按类型清 Buff | 未覆盖 | 按 ClearType 批量清理对象 Buff。 |
| 27 | BuffSlot 生命周期 | GetBuffClearType | `CBuffEffectSlot::GetBuffClearType` | 暂无 Buff 清理类型查询 | 未覆盖 | 根据 BuffIndex 获取清理策略。 |
| 28 | BuffSlot 生命周期 | Slot 上限 | `MAX_BUFFEFFECT` | `MaxBuffEffect = 32` 已有 | 部分覆盖 | 实现 32 槽上限检查和满槽失败行为。 |
| 29 | BuffSlot 生命周期 | 重复 Buff 策略 | `SetEffect` 重复分支 | 暂无重复策略 | 未覆盖 | 明确重复 Buff 是刷新时间、覆盖强度、拒绝还是叠加。 |
| 30 | BuffSlot 生命周期 | 强弱覆盖策略 | `gObjCheckPowerfulEffect` | 暂无强弱比较 | 未覆盖 | 防止低强度 Buff 覆盖高强度 Buff。 |
| 31 | Buff 增删查清 API | gObjAddBuffEffect 简化版 | `gObjAddBuffEffect(lpObj,iBuffIndex)` | 暂无全局 AddBuff API | 未覆盖 | 按 BuffIndex 添加默认定义 Buff。 |
| 32 | Buff 增删查清 API | gObjAddBuffEffect 完整版 | `gObjAddBuffEffect(...EffectType1...)` | 暂无完整 AddBuff API | 未覆盖 | 支持调用方传入两个效果类型、两个效果值和持续时间。 |
| 33 | Buff 增删查清 API | gObjAddBuffEffect Duration 版 | `gObjAddBuffEffect(lpObj,iBuffIndex,Duration)` | 暂无指定时长添加 | 未覆盖 | 使用定义默认效果，但覆盖持续时间。 |
| 34 | Buff 增删查清 API | gObjAddPeriodBuffEffect | `gObjAddPeriodBuffEffect` | 暂无期限 Buff 添加 | 未覆盖 | 按 PeriodBuff 定义添加期限 Buff。 |
| 35 | Buff 增删查清 API | gObjAddBuffEffectForInGameShop | `gObjAddBuffEffectForInGameShop` | 暂无商城 Buff 入口 | 未覆盖 | 支持商城/运营道具按 item code 添加 Buff。 |
| 36 | Buff 增删查清 API | gObjRemoveBuffEffect | `gObjRemoveBuffEffect` | 暂无 RemoveBuff API | 未覆盖 | 删除指定 Buff，并触发数值重算和协议通知。 |
| 37 | Buff 增删查清 API | gObjClearBuffEffect | `gObjClearBuffEffect` | 暂无 ClearBuff API | 未覆盖 | 按清理类型删除一组 Buff。 |
| 38 | Buff 增删查清 API | gObjCheckUsedBuffEffect | `gObjCheckUsedBuffEffect` | 暂无 CheckBuff API | 未覆盖 | 供技能、移动、攻击、任务、事件判断状态。 |
| 39 | Buff 增删查清 API | gObjRemoveOneDebuffEffect | `gObjRemoveOneDebuffEffect` | 暂无 Debuff 单个移除 | 未覆盖 | 支持净化类技能或道具移除一个负面状态。 |
| 40 | Buff 增删查清 API | gObjChangeBuffValidTime | `gObjChangeBuffValidTime` | 暂无修改剩余时间 | 未覆盖 | 支持延长、缩短或刷新 Buff 有效期。 |
| 41 | Buff 时间与过期 | gObjCheckBuffEffectList | `gObjCheckBuffEffectList` | `Process1000ms` 未处理 Buff | 未覆盖 | 每秒检查 Buff 过期、持续伤害和清理。 |
| 42 | Buff 时间与过期 | SetActiveEffectAtTick | `gObjSetActiveEffectAtTick` | 暂无活跃 Tick 标记 | 未覆盖 | 刷新 Buff 生效 tick，避免重复触发或漏触发。 |
| 43 | Buff 时间与过期 | 永久 Buff | Duration `-10` 等 | 暂无永久 Buff 语义 | 未覆盖 | 定义永久 Buff 的持续时间表达和清理条件。 |
| 44 | Buff 时间与过期 | 限时 Buff | `EffectDuration > 0` | 暂无限时 Buff 过期逻辑 | 未覆盖 | 到期自动移除并通知客户端。 |
| 45 | Buff 时间与过期 | 登出清理 | `CLEAR_TYPE_LOGOUT` | 登出未接入 Buff 清理 | 未覆盖 | 登出时按清理类型删除非持久 Buff。 |
| 46 | Buff 时间与过期 | 死亡清理 | 死亡相关 clear type | 死亡流程未接入 Buff 清理 | 未覆盖 | 玩家或怪物死亡时清理指定 Buff。 |
| 47 | Buff 时间与过期 | 地图切换清理 | 地图切换清理分支 | 地图移动未接入 Buff 清理 | 未覆盖 | 切地图时处理不能跨地图保留的状态。 |
| 48 | Buff 时间与过期 | 期限道具过期 | `PeriodItemEx` 过期逻辑 | 暂无 PeriodItem 生命周期 | 未覆盖 | 期限道具到期后移除对应 Buff 和物品效果。 |
| 49 | Buff 时间与过期 | 剩余时间计算 | `LeftTime` | 暂无剩余时间计算 | 未覆盖 | 协议下发时计算 Buff 剩余秒数。 |
| 50 | Buff 时间与过期 | 过期后角色重算 | `ClearBuffEffect` 后重算 | `Player.calc` 未由 Buff 触发 | 未覆盖 | Buff 过期后刷新角色属性、攻速、HP/MP/SD 上限。 |
| 51 | Buff 数值应用 | SetBuffEffect | `CBuffEffect::SetBuffEffect` | 暂无 Buff 数值写入 | 未覆盖 | 根据 EffectType 将数值加到对象对应字段。 |
| 52 | Buff 数值应用 | ClearBuffEffect | `CBuffEffect::ClearBuffEffect` | 暂无 Buff 数值移除 | 未覆盖 | 移除 Buff 时回退对应属性加成。 |
| 53 | Buff 数值应用 | SetActiveBuffEffect | `CBuffEffect::SetActiveBuffEffect` | 暂无主动效果触发 | 未覆盖 | 对部分 Buff 执行即时生效逻辑。 |
| 54 | Buff 数值应用 | ApplyPrevEffectStat | `CBuffEffect::ApplyPrevEffectStat` | `05-formula.md` 已标为缺口 | 未覆盖 | 在角色基础属性聚合前应用 Buff stat。 |
| 55 | Buff 数值应用 | ClearPrevEffectStat | `CBuffEffect::ClearPrevEffectStat` | 暂无 Buff stat 清理 | 未覆盖 | 重算前清理上次 Buff stat 加成。 |
| 56 | Buff 数值应用 | ApplyPrevEffectAll | `CBuffEffect::ApplyPrevEffectAll` | `Player.calc` 未接 Buff | 未覆盖 | 在装备、技能、套装等计算链中应用 Buff 全量效果。 |
| 57 | Buff 数值应用 | ClearPrevEffectAll | `CBuffEffect::ClearPrevEffectAll` | `Player.calc` 未清 Buff 字段 | 未覆盖 | 重算前清理所有 Buff 派生字段。 |
| 58 | Buff 数值应用 | gObjGetTotalValueOfEffect | `gObjGetTotalValueOfEffect` | 暂无按 EffectType 汇总 | 未覆盖 | 汇总所有 Buff 对某个 EffectType 的总加成。 |
| 59 | Buff 数值应用 | gObjGetValueOfBuffIndex | `gObjGetValueOfBuffIndex` | 暂无按 BuffIndex 取值 | 未覆盖 | 查询指定 Buff 的两个效果值，供技能/事件逻辑使用。 |
| 60 | Buff 数值应用 | Player.calc 接入 | `ObjCalCharacter.cpp` 调用 BuffEffect | `Player.calc` 未接 Buff | 未覆盖 | Buff 系统完成后必须纳入角色重算顺序。 |
| 61 | 持续伤害/恢复效果 | GiveDamageEffect | `CBuffEffect::GiveDamageEffect` | 暂无 Buff 持续伤害 | 未覆盖 | 对持续伤害 Buff 周期性扣血并处理死亡。 |
| 62 | 持续伤害/恢复效果 | PoisonEffect | `CBuffEffect::PoisonEffect` | `AttackEffectReply.PoisonEffect` 固定 0 | 未覆盖 | 实现中毒 Tick、伤害比例和状态表现。 |
| 63 | 持续伤害/恢复效果 | GiveDamageFillHPEffect | `CBuffEffect::GiveDamageFillHPEffect` | 暂无伤害吸血 Buff | 未覆盖 | 按造成伤害恢复 HP。 |
| 64 | 持续伤害/恢复效果 | HP Buff | `EFFECTTYPE_HP` | 暂无 HP 上限/即时 HP Buff | 未覆盖 | 支持生命之光、圣诞治疗等 HP 类效果。 |
| 65 | 持续伤害/恢复效果 | MP Buff | `EFFECTTYPE_MANA` | 暂无 MP Buff | 未覆盖 | 支持法力增加或保护类 Buff。 |
| 66 | 持续伤害/恢复效果 | SD/AG Buff | SD/AG 相关 EffectType | 暂无 SD/AG Buff | 未覆盖 | 支持护盾、AG 恢复和 SD 比例类效果。 |
| 67 | 持续伤害/恢复效果 | 自动恢复 Buff | 恢复类 EffectType | `recoverHPSD/recoverMPAG` 未读 Buff | 未覆盖 | 将恢复增强类 Buff 接入周期恢复。 |
| 68 | 持续伤害/恢复效果 | 持续效果死亡边界 | Buff Tick 死亡处理 | 死亡流程未接入 Buff | 未覆盖 | Buff 造成死亡时触发正确死亡、经验和掉落边界。 |
| 69 | 持续伤害/恢复效果 | 持续效果推送 | `GCUseBuffEffect`、HPSD | 暂无持续效果协议 | 未覆盖 | Tick 变化后推送 HP/MP/SD/AG 和 Buff 状态。 |
| 70 | 持续伤害/恢复效果 | 持续效果测试 | BuffEffect Tick 语义 | 暂无测试 | 未覆盖 | 覆盖毒伤、恢复、到期、死亡和断线场景。 |
| 71 | Debuff 与控制状态 | Poison | `BUFFTYPE_POISON` | 暂无中毒状态 | 未覆盖 | 中毒影响攻击效果显示和周期伤害。 |
| 72 | Debuff 与控制状态 | Freeze | `BUFFTYPE_FREEZE` | 暂无冰冻状态 | 未覆盖 | 冰冻影响移动速度或移动行为。 |
| 73 | Debuff 与控制状态 | Stone | `BUFFTYPE_STONE` | 暂无石化状态 | 未覆盖 | 石化限制移动、攻击和技能。 |
| 74 | Debuff 与控制状态 | Stun | `BUFFTYPE_STUN` | 暂无眩晕状态 | 未覆盖 | 眩晕阻止移动、攻击、技能和部分交互。 |
| 75 | Debuff 与控制状态 | Sleep | `BUFFTYPE_SLEEP` | 暂无睡眠状态 | 未覆盖 | 睡眠限制行动并在受击或到期时清理。 |
| 76 | Debuff 与控制状态 | Blind | `BUFFTYPE_BLIND_2` | 暂无致盲状态 | 未覆盖 | 致盲影响命中、目标或客户端表现。 |
| 77 | Debuff 与控制状态 | EarthBinds | `BUFFTYPE_EARTH_BINDS` | 暂无束缚状态 | 未覆盖 | 束缚限制移动并参与技能/事件校验。 |
| 78 | Debuff 与控制状态 | 移动限制 | `protocol.cpp` 移动前 Buff 检查 | 移动入口未查 Buff | 未覆盖 | 移动请求必须拒绝石化、眩晕、睡眠、束缚等状态。 |
| 79 | Debuff 与控制状态 | 攻击限制 | `ObjBaseAttack` Buff 检查 | 攻击入口未查控制 Buff | 未覆盖 | 攻击前校验控制、免疫、隐身等状态。 |
| 80 | Debuff 与控制状态 | 技能限制 | `ObjUseSkill` Buff 检查 | `canUseSkill` 未查 Buff | 未覆盖 | 施法前校验控制状态、安全状态和免疫状态。 |
| 81 | 技能/攻击 Buff 来源 | GreaterLife | `ObjUseSkill` 添加 `BUFFTYPE_HP_INC` | 技能系统未触发 Buff | 未覆盖 | 生命之光应添加 HP Buff，并支持队伍目标。 |
| 82 | 技能/攻击 Buff 来源 | GreaterDefense | `BUFFTYPE_DEFENSE_POWER_INC` | 技能系统未触发 Buff | 未覆盖 | 防御力增加技能应进入 Buff 系统。 |
| 83 | 技能/攻击 Buff 来源 | GreaterDamage | `BUFFTYPE_ATTACK_POWER_INC` | 技能系统未触发 Buff | 未覆盖 | 攻击力增加技能应进入 Buff 系统。 |
| 84 | 技能/攻击 Buff 来源 | SoulBarrier | SoulBarrier 相关 Buff | 技能系统未触发 Buff | 未覆盖 | 守护之魂应添加减伤或防御类 Buff。 |
| 85 | 技能/攻击 Buff 来源 | InfinityArrow | `BUFFTYPE_INFINITY_ARROW` | 技能系统未触发 Buff | 未覆盖 | 无限箭应添加长期 Buff 并影响箭矢消耗。 |
| 86 | 技能/攻击 Buff 来源 | DamageReflect | `BUFFTYPE_DAMAGE_REFLECT` | 暂无反伤 Buff | 未覆盖 | 反伤技能应添加 Buff 并在受击时触发。 |
| 87 | 技能/攻击 Buff 来源 | Berserker | `BUFFTYPE_BERSERKER*` | 暂无狂暴 Buff | 未覆盖 | 狂暴影响攻击、魔攻、防御等属性。 |
| 88 | 技能/攻击 Buff 来源 | Monk/RageFighter Buff | `BUFFTYPE_MONK_*` | RF 技能未实现 Buff | 未覆盖 | RF 的忽防、降防等技能进入 Buff 系统。 |
| 89 | 技能/攻击 Buff 来源 | Monster Immune | `BUFFTYPE_MONSTER_MAGIC_IMMUNE` 等 | 怪物免疫 Buff 未实现 | 未覆盖 | 怪物技能可添加魔法/物理免疫 Buff。 |
| 90 | 技能/攻击 Buff 来源 | 攻击触发 Poison/Freeze | `ObjBaseAttack` 中毒/冰冻分支 | 攻击效果固定 0 | 未覆盖 | 攻击命中时按技能和概率添加 Debuff。 |
| 91 | 道具/期限 Buff 来源 | PeriodItemEx | `PeriodItemEx.cpp/.h` | 暂无期限道具系统 | 未覆盖 | 管理期限道具 Buff 的保存、恢复、删除和到期。 |
| 92 | 道具/期限 Buff 来源 | InGameShop Buff | `gObjAddBuffEffectForInGameShop` | 暂无商城 Buff | 未覆盖 | 商城期限道具触发对应 Buff。 |
| 93 | 道具/期限 Buff 来源 | Halloween Buff | `BUFFTYPE_HALLOWEEN_*` | 使用道具未接 Buff | 未覆盖 | 万圣节类消耗道具触发 Buff。 |
| 94 | 道具/期限 Buff 来源 | CherryBlossom Buff | `BUFFTYPE_CHERRYBLOSSOM_*` | 樱花道具未接 Buff | 未覆盖 | 樱花饮料/花瓣等道具触发 Buff。 |
| 95 | 道具/期限 Buff 来源 | Light Blessing | `BUFFTYPE_LIGHT_BLESSING*` | 光之祝福未接 Buff | 未覆盖 | 光之祝福类经验道具触发 Buff。 |
| 96 | 道具/期限 Buff 来源 | 宠物 Buff | 宠物/Helper 相关 Buff | `conf.PetRing.Pets` 已读配置 | 部分覆盖 | 宠物装备后应转化为统一 Buff 或属性来源。 |
| 97 | 道具/期限 Buff 来源 | 戒指 Buff | Ring 相关 Buff | `conf.PetRing.Rings` 已读配置 | 部分覆盖 | 变身戒指和活动戒指应接入 Buff 与外观状态。 |
| 98 | 道具/期限 Buff 来源 | Panda/Skeleton/Unicorn | `ObjCalCharacter` 期限宠物判断 | 道具分类存在但未接统一 Buff | 未覆盖 | 熊猫、骷髅、兽角等经验或属性加成归 Buff 管理。 |
| 99 | 道具/期限 Buff 来源 | 期限 Buff 保存 | `RequestPeriodBuffInsert` | DB 无期限 Buff 状态 | 未覆盖 | 保存期限 Buff 到期时间，支持重登恢复。 |
| 100 | 道具/期限 Buff 来源 | 期限 Buff 删除 | `RequestPeriodBuffDelete` | 暂无期限 Buff 删除 | 未覆盖 | 到期、强制删除或物品删除时清理期限 Buff。 |
| 101 | 视野与协议同步 | GCUseBuffEffect | `GCUseBuffEffect` | 暂无 Buff 增删协议 | 未覆盖 | Buff 增加、移除、刷新时下发客户端图标和剩余时间。 |
| 102 | 视野与协议同步 | gObjSendBuffList | `gObjSendBuffList` | 暂无 Buff 列表响应 | 未覆盖 | 登录、请求或重同步时下发当前 Buff 列表。 |
| 103 | 视野与协议同步 | gObjMakeViewportState | `gObjMakeViewportState` | 视野创建未携带 Buff | 未覆盖 | 创建视野对象时携带 Buff 状态。 |
| 104 | 视野与协议同步 | AttackEffectReply Poison/Ice | `GCDamageSend` 中 Poison/Ice | `MsgAttackEffectReply` 字段固定 0 | 未覆盖 | 攻击效果包体应读取目标 Buff。 |
| 105 | 视野与协议同步 | 创建视野 BuffCount | `BuffEffectCount` | `MsgCreateViewport*` 无 Buff 列表 | 未覆盖 | 视野创建时包含 Buff 数量和 Buff 列表。 |
| 106 | 视野与协议同步 | Buff 增量通知 | Buff add/remove notice | 暂无增量通知 | 未覆盖 | Buff 变化时推送本人和视野对象。 |
| 107 | 视野与协议同步 | Buff 列表请求 | `protocol.cpp` Buff list | handle 未映射 Buff 列表请求 | 未覆盖 | 补齐客户端请求当前 Buff 列表的协议入口。 |
| 108 | 视野与协议同步 | 客户端图标 | BuffIndex/Icon | 暂无图标映射 | 未覆盖 | 确保 BuffIndex 与客户端图标一致。 |
| 109 | 视野与协议同步 | 剩余时间下发 | `nBuffDuration` | 暂无剩余时间字段 | 未覆盖 | Buff 列表和增量通知包含剩余时间。 |
| 110 | 视野与协议同步 | GM 隐身表现 | `BUFFTYPE_INVISABLE` | 暂无 GM 隐身 Buff | 未覆盖 | GM 隐身影响视野、攻击和协议表现。 |
| 111 | 经验/掉落/事件 Buff | GetPremiumExp | `GetPremiumExp` | 经验系统未读 Buff | 未覆盖 | Premium Buff 影响经验倍率。 |
| 112 | 经验/掉落/事件 Buff | CheckItemOptForGetExpExRenewal | `CheckItemOptForGetExpExRenewal` | `09-exp.md` 标为缺口 | 未覆盖 | 统一处理经验道具、Buff、印章和宠物经验加成。 |
| 113 | 经验/掉落/事件 Buff | Crywolf Altar Buff | `CrywolfAltar.cpp` Buff 逻辑 | 世界事件未实现 | 未覆盖 | Crywolf 祭坛 NPC/玩家状态通过 Buff 表达。 |
| 114 | 经验/掉落/事件 Buff | Santa/XMas Buff | `XMasAttackEvent.cpp` | 活动 Buff 未实现 | 未覆盖 | 圣诞 NPC 或事件区域添加群体 Buff。 |
| 115 | 经验/掉落/事件 Buff | Quest Buff | `Quests.cpp`、`QuestExp` 添加 Buff | 任务系统未接 Buff | 未覆盖 | 任务奖励或任务状态可添加 Buff。 |
| 116 | 经验/掉落/事件 Buff | MonsterSkill Buff | `TMonsterSkillElement.cpp` | 怪物技能未接 Buff | 未覆盖 | 怪物技能添加眩晕、中毒、免疫等 Buff。 |
| 117 | 经验/掉落/事件 Buff | Guild Period Buff | `RequestGuildPeriodBuffInsert/Delete` | 战盟系统未实现 | 未覆盖 | 战盟期限 Buff 后续归社交/战盟系统触发，状态归 Buff 系统。 |
| 118 | 经验/掉落/事件 Buff | Muun Buff | `m_MuunEffectList` | Muun 未实现 | 未覆盖 | Muun 效果后续可接入统一 Buff/Effect 查询。 |
| 119 | 经验/掉落/事件 Buff | 防作弊/安全状态 | 速度惩罚、控制状态 | `conf.Common` 有 SpeedHackPenalty | 部分覆盖 | 安全惩罚类状态可复用 Buff 或 Effect 表达。 |
| 120 | 经验/掉落/事件 Buff | 端到端回归测试 | BuffEffect/Protocol/User 流程 | 暂无 Buff 测试 | 未覆盖 | 覆盖添加、覆盖、过期、死亡清理、登出恢复、视野同步、经验加成和控制状态。 |
