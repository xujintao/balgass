## Skill

`game/skill` owns skill static data, runtime skill records, skill settings, and master skill lookup.

### Package files

- `skill_table.go`: skill indexes, `SkillBase`, `MasterSkillBase`, `SkillManager`, IGC skill table loading, master skill tree loading, and hit box lookup.
- `skill.go`: runtime `Skill` and `Skills`, JSON/database serialization, normal skill learning, master skill learning, and skill marshal format.
- `settings.go`: runtime skill settings loaded from `IGC_SkillSettings.ini`.

### Config files

- `Skills/IGC_SkillList.xml`: base skill table.
- `Skills/IGC_SkillSettings.ini`: skill runtime settings.
- `Skills/IGC_SkillSpear.hit`: spear skill hit box.
- `Skills/IGC_SkillElect.hit`: electric skill hit box.
- `IGC_MasterSkillTree.xml`: master skill tree.

### Runtime flow

`SkillManager` loads static skill data during package initialization. Objects keep learned skills in `skill.Skills`, learn or forget entries through object methods, and use the selected runtime skill when handling object skill actions.
