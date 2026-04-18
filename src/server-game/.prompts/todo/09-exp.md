# 9. 经验系统

本模块覆盖经验表、经验倍率、怪物击杀经验、普通升级、大师经验、大师升级、组队经验、经验附加来源、经验扣减、协议下发和测试。具体经验公式归公式系统，本模块只记录经验流程、状态变化、倍率应用和协议边界。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 经验表与基础配置 | 普通经验表初始化 | `gLevelExperience`、`gObjNextExpCal` | `ExperienceTable`、`formula.SetExpTable_Normal` | 部分覆盖 | Go 侧启动时按最大普通等级生成经验表，需补表长度、0 级、最大等级边界测试。 |
| 2 | 经验表与基础配置 | 大师经验表初始化 | `g_MasterLevelSkillTreeSystem.gMasterExperience` | `MasterExperienceTable`、`formula.SetExpTable_Master` | 部分覆盖 | Go 侧按最大大师等级生成表，需确认索引从 0/1 开始的语义。 |
| 3 | 经验表与基础配置 | Lua 经验公式调用 | GameServer 经验表脚本/公式 | `Misc/ExpCalc.lua`、`game/formula/exp_calc.go` | 部分覆盖 | 公式调用已封装，需补失败加载、返回负数、溢出等保护。 |
| 4 | 经验表与基础配置 | 最大普通等级配置 | `UserMaxLevel` | `conf.Common.General.MaxLevelNormal` | 部分覆盖 | Go 侧用于建表，但升级路径未完整阻断超过上限。 |
| 5 | 经验表与基础配置 | 最大大师等级配置 | Master level max config | `conf.Common.General.MaxLevelMaster` | 部分覆盖 | Go 侧用于建表和协议下发，升级路径需防止越界。 |
| 6 | 经验表与基础配置 | 普通 NextExp 查询 | `gObjNextExpCal` | `exp.ExperienceTable[p.Level]` | 部分覆盖 | 需集中封装 next exp 查询，避免登录和升级路径重复索引。 |
| 7 | 经验表与基础配置 | 大师 NextExp 查询 | `MasterNextExp` | `exp.MasterExperienceTable[p.masterLevel]` | 部分覆盖 | 需确认当前等级对应的是“本级门槛”还是“下级门槛”。 |
| 8 | 经验表与基础配置 | 经验表越界保护 | GameServer 等级上限判断 | 当前多处直接索引表 | 未覆盖 | 登录、升级、大师数据下发都需要防止等级坏数据导致 panic。 |
| 9 | 经验表与基础配置 | 经验字段持久化 | `Experience`、`MasterExperience`、`MasterNextExp` | `model.Character.Experience/MasterExperience` | 部分覆盖 | Go 侧字段已持久化，需确认升级后保存时机和异常恢复。 |
| 10 | 经验表与基础配置 | 经验表加载诊断 | GameServer load log | 当前无专门诊断 | 未覆盖 | 需要输出最大等级、表长度、关键等级经验值和异常配置日志。 |
| 11 | ExpManager 倍率配置 | ExpSystem XML 加载 | `CExpManager::LoadScript` | `ExpManager.init` 读取 `IGC_ExpSystem.xml` | 部分覆盖 | Go 侧已读取 XML，但只使用静态 Normal/Master/Event。 |
| 12 | ExpManager 倍率配置 | StaticExp 普通倍率 | `m_fStaticExp` | `ExpManager.Normal` | 已覆盖 | 普通击杀经验已乘基础普通倍率。 |
| 13 | ExpManager 倍率配置 | StaticExp 大师倍率 | `m_fStaticMLExp` | `ExpManager.Master` | 已覆盖 | 大师击杀经验已乘基础大师倍率。 |
| 14 | ExpManager 倍率配置 | StaticExp 活动倍率 | `m_fEventExp` | `ExpManager.Event` | 未覆盖 | Go 侧读取 Event 但击杀经验未使用。 |
| 15 | ExpManager 倍率配置 | StaticExp 任务倍率 | `m_fQuestExp` | XML 结构读取但未保存 | 未覆盖 | 任务经验倍率应记录到 ExpManager 并给任务系统调用。 |
| 16 | ExpManager 倍率配置 | CalcType 静态模式 | `EXP_CALC_STATIC_ONLY` | `CalcType` 字段读取为 string 后未使用 | 未覆盖 | 需实现静态、静态+动态、静态*动态三类计算。 |
| 17 | ExpManager 倍率配置 | DynamicExpRangeList 解析 | `m_vtExpRanges` | `DynamicExpRangeList.Range` 临时结构 | 未覆盖 | Go 侧未保存 range，需建结构保存 reset/level/normal/master。 |
| 18 | ExpManager 倍率配置 | reset 范围匹配 | `iMinReset/iMaxReset` | 当前无 reset 倍率匹配 | 未覆盖 | 需和账号角色 reset 字段接入，或记录 reset 系统未完成依赖。 |
| 19 | ExpManager 倍率配置 | 动态等级范围匹配 | `iMinLevel/iMaxLevel` | 当前无动态匹配 | 未覆盖 | 需按当前普通等级或普通+大师等级选择倍率。 |
| 20 | ExpManager 倍率配置 | Exp debug 日志 | `m_bDebugMode`、`[EXP DEBUG]` | 当前无 debug 输出 | 未覆盖 | 需在倍率计算时输出等级、reset、ML、CalcType、最终倍率。 |
| 21 | 击杀经验入口 | 怪物死亡给经验入口 | `gObjMonsterExpDivision` | `Monster.Die -> DieGiveExperience` | 部分覆盖 | Go 侧怪物死亡会给击杀者经验，但未走完整伤害列表分配。 |
| 22 | 击杀经验入口 | 伤害来源记录 | `sHD[n].HitDamage` | 当前 `damage` 参数只代表本次伤害 | 未覆盖 | 需要记录怪物死亡前所有攻击者累计伤害。 |
| 23 | 击杀经验入口 | 单人击杀经验 | `gObjMonsterExpSingle` | `Object.DieGiveExperience` | 部分覆盖 | Go 侧实现基础单人经验公式和倍率。 |
| 24 | 击杀经验入口 | 伤害占比分配 | `exp = dmg * exp / tot_dmg` | 当前未按总伤害分配 | 未覆盖 | 多人打同一怪时应按贡献分配经验。 |
| 25 | 击杀经验入口 | BattleCore 禁经验 | `SERVER_BATTLECORE` 返回 0 | 当前无 server type 判断 | 未覆盖 | BattleCore/跨服场景应禁止或特殊处理经验。 |
| 26 | 击杀经验入口 | DevilSquare 经验分流 | `g_DevilSquare.gObjMonsterExpSingle` | 当前无 DS 分流 | 未覆盖 | DS 地图经验应交给事件系统处理。 |
| 27 | 击杀经验入口 | 怪物等级基础经验 | `(Level+25)*Level/3` | `level := (obj.Level+25)*obj.Level/3` | 已覆盖 | 基础公式已对齐旧版单人公式。 |
| 28 | 击杀经验入口 | 玩家高等级经验衰减 | `(monster+10)<userlevel` | `if obj.Level+10 < targetLevel` | 已覆盖 | 需补大师等级参与 userlevel 的测试。 |
| 29 | 击杀经验入口 | 65 级以上怪物补偿 | `(Level-64)*(Level/4)` | `if obj.Level >= 65` | 部分覆盖 | Go 写法为 `(obj.Level-64)*obj.Level/4`，需核对整数优先级和 GameServer 一致性。 |
| 30 | 击杀经验入口 | 经验为 0 处理 | `if exp <= 0` 分支 | `if baseExp <= 0 { return }` | 已覆盖 | 需补负数、总伤害为 0、怪物等级异常时的保护。 |
| 31 | 击杀经验倍率叠加 | 地图普通经验加成 | `g_MapAttr.GetExpBonus` | `maps.MapManager.GetExpBonus` | 已覆盖 | Go 侧普通经验会乘 `(1+mapBonus)`。 |
| 32 | 击杀经验倍率叠加 | 地图大师经验加成 | `g_MapAttr.GetMasterExpBonus` | `maps.MapManager.GetMasterExpBonus` | 已覆盖 | Go 侧大师经验会使用大师地图倍率。 |
| 33 | 击杀经验倍率叠加 | 基础倍率选择 | `g_ExpManager.GetExpMultiplier` | `ExpManager.Normal/Master` | 部分覆盖 | Go 侧只用静态倍率，未按等级/reset/CalcType 动态计算。 |
| 34 | 击杀经验倍率叠加 | VIP 经验倍率 | `g_VipSystem.GetExpBonus` | 当前无 VIP 经验加成 | 未覆盖 | 需接入账号/角色 VIP，并明确与基础倍率相加。 |
| 35 | 击杀经验倍率叠加 | BonusEvent 普通经验 | `g_BonusEvent.GetAddExp` | 当前无 `BonusEvent` | 未覆盖 | 需实现活动经验倍率并参与普通经验叠加。 |
| 36 | 击杀经验倍率叠加 | BonusEvent 大师经验 | `g_BonusEvent.GetAddMLExp` | 当前无 `BonusEvent` | 未覆盖 | 大师经验应使用独立 MasterExpMultiplier。 |
| 37 | 击杀经验倍率叠加 | Crywolf 经验惩罚 | `g_CrywolfSync.GetGettingExpPenaltyRate` | 当前无 Crywolf 惩罚 | 未覆盖 | 事件状态为占领惩罚时应按百分比降低经验。 |
| 38 | 击杀经验倍率叠加 | DoppelGanger 经验倍率 | `g_DoppelGanger.GetDoppelGangerExpRate/MLExpRate` | 当前无 DG 倍率 | 未覆盖 | DG 地图需要区分普通经验和大师经验倍率。 |
| 39 | 击杀经验倍率叠加 | BloodCastle 经验折减 | `BC_MAP_RANGE` 经验 *50% | 当前无 BC 经验修正 | 未覆盖 | BC 地图经验应按 GameServer 特殊规则修正。 |
| 40 | 击杀经验倍率叠加 | ChaosCastle 经验加成 | `g_ChaosCastle.GetExperienceBonus` | 当前无 CC 经验修正 | 未覆盖 | CC 地图需按场次/索引应用经验加成。 |
| 41 | 普通经验增加与升级 | 普通升级统一入口 | `gObjLevelUp` | `Player.LevelUp` | 部分覆盖 | Go 侧经验增加和升级在一个方法中，需拆出更明确的 AddExperience 语义。 |
| 42 | 普通经验增加与升级 | 普通经验字段累加 | `Experience += addexp` | `p.experience += addexp` | 已覆盖 | 需改用 64 位或防止大经验溢出。 |
| 43 | 普通经验增加与升级 | 普通等级上限阻断 | `UserMaxLevel` 判断 | 当前未判断 MaxLevelNormal | 未覆盖 | 达到最大普通等级后应停止普通升级并切换或等待大师条件。 |
| 44 | 普通经验增加与升级 | NextExp 判定 | `Experience < NextExp` | `p.experience < levelUpExp` | 已覆盖 | 需统一 NextExp 语义，避免当前等级门槛和下级门槛混淆。 |
| 45 | 普通经验增加与升级 | 经验截断到门槛 | `Experience = NextExp` | `p.experience = levelUpExp` | 已覆盖 | 目前会截断溢出经验，需确认是否符合目标版本。 |
| 46 | 普通经验增加与升级 | 多级连升 | GameServer 单次通常升一级 | 当前单次最多升一级 | 未覆盖 | 如果保留 GameServer 行为则写明；如支持连升需循环并限制次数。 |
| 47 | 普通经验增加与升级 | 溢出经验处理 | GameServer 截断 | 当前截断 | 部分覆盖 | 需明确是否丢弃超出本级的经验，避免后续实现分歧。 |
| 48 | 普通经验增加与升级 | 升级返回值语义 | `gObjLevelUp` 返回是否发送经验包 | `LevelUp` 返回是否已升级 | 部分覆盖 | 当前 `DieGiveExperience` 根据 false 下发经验包，需固定语义并测试。 |
| 49 | 普通经验增加与升级 | 宠物经验联动 | `gObjSetExpPetItem` | 当前无宠物经验 | 未覆盖 | 普通经验增加时需同步宠物/DarkSpirit 等子系统。 |
| 50 | 普通经验增加与升级 | 升级后保存时机 | CharInfoSave/DB save | `Player.Save` 保存字段 | 部分覆盖 | 需明确升级后立即保存、定时保存或退出保存。 |
| 51 | 普通升级奖励 | 角色等级加一 | `lpObj->Level++` | `p.Level++` | 已覆盖 | 需配合最大等级和转大师条件。 |
| 52 | 普通升级奖励 | 普通职业升级点 | `gLevelUpPointNormal` | `LevelPoint5` | 已覆盖 | 普通职业升级加 5 点。 |
| 53 | 普通升级奖励 | 特殊职业升级点 | `gLevelUpPointMGDL` | `LevelPoint7` | 部分覆盖 | Go 包含 MG/DL/RF/GL，需确认职业枚举和 GameServer 完全一致。 |
| 54 | 普通升级奖励 | reset 后点数阻断 | `iBlockLevelUpPointAfterResets` | 当前无 reset 阻断 | 未覆盖 | reset 次数达到配置时升级不发普通点数。 |
| 55 | 普通升级奖励 | PlusStatQuest 额外点 | `PlusStatQuestClear` | 当前无额外升级点 | 未覆盖 | 完成额外点任务后升级应多给 1 点。 |
| 56 | 普通升级奖励 | Muun 升级条件检查 | `CheckMuunItemConditionLevelUp` | 当前无 Muun 联动 | 未覆盖 | 升级后需要触发 Muun 条件刷新。 |
| 57 | 普通升级奖励 | 升级后角色重算 | `gObjCalCharacter.CalcCharacter` | `p.calc()` | 已覆盖 | Go 侧升级后重算角色属性。 |
| 58 | 普通升级奖励 | 升级后 HP/MP/SD/AG 回满 | `GCReFillSend`、`GCManaSend` | 设置 HP/SD/MP/AG 并 `PushHPSD/PushMPAG` | 已覆盖 | 需测试升级时客户端资源条同步。 |
| 59 | 普通升级奖励 | 无限箭技能自动学习 | `InfinityArrowUseLevel` | 当前无升级自动加技能 | 未覆盖 | 妖精达到等级和转职条件时应自动学习无限箭。 |
| 60 | 普通升级奖励 | 普通升级协议 | `GCLevelUpMsgSend` | `MsgLevelUpReply` | 部分覆盖 | Go 侧已下发等级、点数和资源上限，需核对包结构和效果广播。 |
| 61 | 大师经验与大师升级 | 大师用户判定 | `IsMasterLevelUser` | `Player.IsMasterLevel` | 部分覆盖 | Go 侧判定逻辑需确认是否仅看普通等级上限，还是还需转职/任务条件。 |
| 62 | 大师经验与大师升级 | 大师经验资格检查 | `CheckMLGetExp` | 当前无完整 ML monster 条件 | 未覆盖 | 需按玩家大师状态、怪物最低等级、地图限制决定是否给大师经验。 |
| 63 | 大师经验与大师升级 | 大师经验字段累加 | `MasterExperience += addexp` | `p.masterExperience += addexp` | 已覆盖 | 需使用 64 位并防止配置表越界。 |
| 64 | 大师经验与大师升级 | MasterNextExp 判定 | `MasterNextExp` | `MasterExperienceTable[p.masterLevel]` | 部分覆盖 | 需明确当前索引对应下一大师等级需求。 |
| 65 | 大师经验与大师升级 | 大师等级上限 | `MaxMasterLevel` | `conf.Common.General.MaxLevelMaster` | 部分覆盖 | 协议下发上限，但升级未阻断越界。 |
| 66 | 大师经验与大师升级 | 大师升级点发放 | `MasterPoint` | `p.masterPoint += MasterPointPerLevel` | 已覆盖 | 需加入 reset 后大师点数阻断配置。 |
| 67 | 大师经验与大师升级 | 大师升级角色重算 | `MasterLevelUp` 后重算 | `p.calc()` | 已覆盖 | 大师等级参与生命魔法和战斗属性计算。 |
| 68 | 大师经验与大师升级 | 大师升级资源回满 | Master level up refill | `PushHPSD`、`PushMPAG` | 已覆盖 | 需测试大师升级时客户端资源刷新。 |
| 69 | 大师经验与大师升级 | 大师技能树联动边界 | `MasterLevelSkillTreeSystem` | `LearnMasterSkill`、`pushMasterSkillList` | 部分覆盖 | 本模块只记录大师点数和经验，具体技能学习归技能系统。 |
| 70 | 大师经验与大师升级 | 大师升级协议 | `GCMasterLevelUpSend` | `MsgMasterLevelUpReply` | 部分覆盖 | Go 侧已下发 MasterLevel、MasterPoint、MaxMasterLevel 和资源上限，需核对 opcode 命名。 |
| 71 | 组队经验 | 组队经验入口 | `gObjExpParty` | 当前无完整组队经验入口 | 未覆盖 | 怪物死亡时应按是否组队选择单人或组队经验。 |
| 72 | 组队经验 | 队伍成员遍历 | `gParty.m_PartyS[partynum].Number` | 当前无经验侧 party 遍历 | 未覆盖 | 需接入 Go 侧 party 模块或先保留接口边界。 |
| 73 | 组队经验 | 同地图过滤 | `lpTargetObj->MapNumber == lpPartyObj->MapNumber` | 当前无实现 | 未覆盖 | 只有同地图成员参与经验分配。 |
| 74 | 组队经验 | 距离 10 格过滤 | `gObjCalDistance < 10` | 当前无实现 | 未覆盖 | 经验分配以怪物和成员距离为准。 |
| 75 | 组队经验 | 队伍最高等级统计 | `toplevel` | 当前无实现 | 未覆盖 | 最高普通+大师等级用于等级差修正。 |
| 76 | 组队经验 | 队伍总等级统计 | `totallevel` | 当前无实现 | 未覆盖 | 只统计符合地图和距离条件的成员。 |
| 77 | 组队经验 | 等级差 +200 限制 | `toplevel >= member+200` | 当前无实现 | 未覆盖 | 等级差过大时低等级成员按 `level+200` 参与总等级。 |
| 78 | 组队经验 | set party 职业判定 | `bCheckSetParty` 职业组合 | 当前无实现 | 未覆盖 | 需按职业组合决定是否使用 set party 经验倍率。 |
| 79 | 组队经验 | 组队人数倍率 | `gPartyExp2..5`、`gSetPartyExp2..5` | `conf.CommonServer.GameServerInfo.Party*ExpBonus` | 未覆盖 | 配置字段存在，需接入人数和 set party 分支。 |
| 80 | 组队经验 | 个人经验分配 | `totalexp*viewpercent*memberLevel/totallevel/100` | 当前无实现 | 未覆盖 | 每个成员按等级权重获得个人经验。 |
| 81 | 组队经验倍率与通知 | 组队普通经验倍率 | `gPartyExp*` | `Party2ExpBonus` 等配置 | 未覆盖 | 普通队伍按人数取倍率。 |
| 82 | 组队经验倍率与通知 | 组队 set party 倍率 | `gSetPartyExp*` | `Party2ExpBonusSet` 等配置 | 未覆盖 | 满足职业组合后使用 set party 倍率。 |
| 83 | 组队经验倍率与通知 | 组队大师经验 | `IsMasterLevelUser` 分支 | 当前无实现 | 未覆盖 | 成员分别按普通/大师状态计算倍率和上限。 |
| 84 | 组队经验倍率与通知 | 组队地图倍率 | `g_MapAttr.GetExpBonus/GetMasterExpBonus` | 当前无组队路径 | 未覆盖 | 每个成员按自身地图倍率计算。 |
| 85 | 组队经验倍率与通知 | 组队 VIP/Event 倍率 | `g_VipSystem`、`g_BonusEvent` | 当前无实现 | 未覆盖 | 每个成员的 VIP 和活动倍率独立叠加。 |
| 86 | 组队经验倍率与通知 | 组队 Crywolf/DG 修正 | Crywolf/DG 分支 | 当前无实现 | 未覆盖 | 每个成员按所在地图和事件状态修正经验。 |
| 87 | 组队经验倍率与通知 | 组队道具加成 | `CheckItemOptForGetExpExRenewal` | 当前无实现 | 未覆盖 | 组队经验也应应用经验道具、buff、印章等加成。 |
| 88 | 组队经验倍率与通知 | 组队 Master 检查 | `CheckMLGetExp` | 当前无实现 | 未覆盖 | 不满足大师经验条件时该成员经验为 0。 |
| 89 | 组队经验倍率与通知 | 组队经验下发 | `GCKillPlayerMasterExpSend` | 当前无组队经验包 | 未覆盖 | 需要按成员分别下发普通或大师经验获得包。 |
| 90 | 组队经验倍率与通知 | 组队经验异常处理 | GameServer 多处 continue/return | 当前无实现 | 未覆盖 | 队伍解散、成员离线、总等级为 0、怪物销毁都需安全处理。 |
| 91 | 经验附加来源 | QuestExp 奖励入口 | `QuestExp*`、`g_ExpManager.m_fQuestExp` | `handle` 0xF630 `questExp` | 未覆盖 | 本模块只定义任务经验入口和倍率，任务流程归任务系统。 |
| 92 | 经验附加来源 | 活动经验入口 | Event 系统经验奖励 | 当前无通用活动经验接口 | 未覆盖 | 需要统一活动模块调用 AddExperience 的入口。 |
| 93 | 经验附加来源 | GM 经验天数 | `GmExpDays`、登录提示 | 当前无 GM Exp | 未覆盖 | 需记录 GM/运营临时经验加成的账号状态和到期。 |
| 94 | 经验附加来源 | 离线经验 | `OfflineLevelling` | 当前无离线练级经验 | 未覆盖 | 需设计离线对象、技能、拾取、经验获得和安全限制。 |
| 95 | 经验附加来源 | 经验道具加成 | `CheckItemOptForGetExpExRenewal` | 当前无经验道具加成 | 未覆盖 | 需覆盖印章、戒指、buff、期限道具对经验的影响。 |
| 96 | 经验附加来源 | ExpUpCharm 位置 | `m_btExpUpCharmPos` | 当前无 charm 位置跟踪 | 未覆盖 | 穿戴/背包变化时需维护经验符咒位置。 |
| 97 | 经验附加来源 | 宠物经验 | `gObjSetExpPetItem` | 当前无宠物经验 | 未覆盖 | DarkHorse、DarkRaven、宠物升级经验需要独立联动。 |
| 98 | 经验附加来源 | DarkSpirit 经验倍率 | `DarkSpiritAddExperience` 配置 | `conf` 有 `ItemDarkSpiritAddExperience` | 未覆盖 | 配置存在但未接入宠物经验。 |
| 99 | 经验附加来源 | 事件副本经验 | DS/BC/CC/DG 专用经验 | 当前只普通击杀路径 | 未覆盖 | 事件副本应有自己的经验入口或修正流程。 |
| 100 | 经验附加来源 | 经验来源枚举 | `szEventType` 参数 | 当前无来源枚举 | 未覆盖 | AddExperience 应携带 Single/Party/Quest/Event/GM/Offline 等来源。 |
| 101 | 经验扣减与限制 | 死亡扣经验入口 | GameServer death exp deduction | `DieDropExperience` placeholder | 未覆盖 | 玩家死亡时应按配置扣普通或大师经验。 |
| 102 | 经验扣减与限制 | PK 扣经验配置 | `PK.ExpDeduction` | `conf.PK.ExpDeduction` | 未覆盖 | 配置已存在，需按 PK 等级和角色等级取扣减比例。 |
| 103 | 经验扣减与限制 | 普通/大师扣减区分 | Master user 分支 | 当前无实现 | 未覆盖 | 大师状态扣 MasterExperience，普通状态扣 Experience。 |
| 104 | 经验扣减与限制 | 最低经验保护 | `minexp/maxexp` | 当前无实现 | 未覆盖 | 扣经验不得低于当前等级最低经验。 |
| 105 | 经验扣减与限制 | 等级上限阻断提示 | `GCServerMsgStringSend` | 当前无上限提示 | 未覆盖 | 达最大等级后应停止升级并可下发提示。 |
| 106 | 经验扣减与限制 | Master 怪物最低等级 | `MonsterMinLevelForMLExp` | `conf.Common.General.MinMonsterLevelForMasterExp` | 未覆盖 | 大师经验应要求怪物等级达到配置。 |
| 107 | 经验扣减与限制 | PVP 不给经验 | `lpObj USER && target USER => exp=0` | 当前攻击死亡路径需确认 | 未覆盖 | 玩家击杀玩家不应走怪物经验。 |
| 108 | 经验扣减与限制 | 地图禁经验 | BattleCore/事件地图限制 | 当前无统一地图禁经验 | 未覆盖 | 需给地图或服务器类型提供禁止经验的判断。 |
| 109 | 经验扣减与限制 | 异常值防护 | GameServer 多处 return/日志 | 当前缺少统一保护 | 未覆盖 | 经验、等级、倍率、伤害、总伤害都需防负数和溢出。 |
| 110 | 经验扣减与限制 | 经验审计日志 | `g_Log.Add` 经验/扣经验日志 | 当前缺少经验审计 | 未覆盖 | 需要记录来源、怪物、基础经验、倍率链、最终经验、升级结果。 |
| 111 | 协议与测试 | 经验获得协议 | `GCKillPlayerExpSend`、`GCKillPlayerMasterExpSend` | `MsgExperienceReply` | 部分覆盖 | Go 侧已有经验包，需区分普通/大师或确认同包兼容。 |
| 112 | 协议与测试 | 普通升级协议 | `GCLevelUpMsgSend` | `MsgLevelUpReply` | 部分覆盖 | 需核对 opcode、字段端序和客户端升级特效。 |
| 113 | 协议与测试 | 大师数据协议 | Master data send | `MsgMasterDataReply` | 部分覆盖 | 登录下发大师等级、经验、next exp、点数，需修正 next exp 计算。 |
| 114 | 协议与测试 | 大师升级协议 | Master level up send | `MsgMasterLevelUpReply` | 部分覆盖 | 需核对 0xF351 命名、包头和字段顺序。 |
| 115 | 协议与测试 | 登录经验下发 | `GCJoinResult`/char load exp | `MsgLoadCharacterReply` | 需修正 | 当前大师状态下 experience/nextExperience 赋值疑似错误，应列为修复点。 |
| 116 | 协议与测试 | 经验包端序 | GameServer protocol structs | `MsgExperienceReply.Marshal` | 部分覆盖 | 当前 Experience 拆高低 16 位小端写入，需用抓包/客户端测试锁定。 |
| 117 | 协议与测试 | AddLevelPoint 联动 | `gObjLevelUpPointAdd` | `AddLevelPoint` | 部分覆盖 | 升级点数发放后，手动加点消耗需和经验系统点数一致。 |
| 118 | 协议与测试 | 单人经验回归测试 | `gObjMonsterExpSingle` | `DieGiveExperience`、`LevelUp` | 未覆盖 | 覆盖低等级、高等级衰减、地图倍率、升级、不升级、上限。 |
| 119 | 协议与测试 | 组队经验回归测试 | `gObjExpParty` | 待实现组队经验 | 未覆盖 | 覆盖距离、同地图、set party、等级差、成员大师状态。 |
| 120 | 协议与测试 | 经验边界和压测 | GameServer 经验日志/保护 | 待实现测试集 | 未覆盖 | 覆盖大经验、负倍率、坏等级、并发击杀、重复死亡、经验表越界。 |
