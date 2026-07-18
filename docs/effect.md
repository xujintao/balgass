## Effect

`game/effect` owns buff static data and runtime effect containers used by objects and combat.

### Package files

- `effect.go`: buff indexes, `BuffBase`, runtime `Effect`, `Effects`, `BuffTable` loading, conflict removal, and combat value aggregation.

### Config files

- `IGC_BuffEffectManager.xml`: buff base table.

### Runtime flow

`BuffTable` loads during package initialization. `Effects.Add` rejects nil or empty buff entries, fills category from `BuffTable`, removes effects with the same buff index or category, and keeps at most 32 active buffs.

Combat reads aggregate values from `Effects` for magic attack, curse attack, critical damage, ignore defense, attack reduction, defense reduction, and reflect.
