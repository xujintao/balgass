## Delay msg

| code | time: ms   | player: gObjStateProc    | monster: gObjMonsterStateProc |
| ---- | ---------- | ------------------------ | ----------------------------- |
| 0    | 100        |                          | base action                   |
| 1    | 800        |                          | die drop item                 |
| 2    | 150        | gObjBackSpring knockback | gObjBackSpring knockback      |
| 3    | 200        |                          | Emotion = 0                   |
| 3    | 2000       | gObjMonsterDieLifePlus   |                               |
| 4    | 100        |                          | Emotion = 3                   |
| 4    | 2000       | gObjGuildWarEnd          |                               |
| 5    | 10000      | BattleSoccerGoalStart    | gObjMemFree                   |
| 7    | 10000      | GCManagerGuildWarEnd     | gObjBackSpring2               |
| 10   | 10         | DamageReflect gObjAttack |                               |
| 12   | 10         | gObjAttack               |                               |
| 13   | 100        | GCReFillSend             |                               |
| 14   | 100        | GCManaSend               |                               |
| 15   | 100        | GCReFillSend             |                               |
| 16   | 200\*(n+1) | GCReFillSend             |                               |
| 50   |            |                          | gObjAttack                    |
| 54   | 100        | gObjAttack               |                               |
| 55   | 400        |                          | gObjAttack                    |
| 55   | 1000       |                          |                               |
| 55   |            |                          |                               |
| 56   |            |                          | gObjAddBuffEffect             |
| 57   |            |                          | gObjBackSpring2               |
| 58   | 300        | gObjAttack               |                               |
| 59   | 1000       | NewSkillProc             |                               |
| 62   |            |                          | gObjAttack                    |
| 65   |            | NewSkillProc             | gObjAttack                    |
| 1000 | 100        | gObjBillRequest          |                               |
| 1001 | 5000       | gObjReqMapSvrAuth        |                               |
