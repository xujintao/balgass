## F303-LoadCharacter

Sequence Diagram

| dir | code | handle                     | handle                          | code | dir | Description |
| --- | ---- | -------------------------- | ------------------------------- | ---- | --- | ----------- |
| ->  | F303 | CGPCharacterMapJoinRequest |                                 |      |     |             |
|     |      | wsDataCli.DataSend         |                                 | 06   | ->  |             |
|     |      |                            | JGGetCharacterInfo              | 06   | <-  |             |
| <-  | F303 |                            | IOCP.DataSend                   |      |     |             |
| <-  | F311 |                            | GSProtocol.GCMagicListMultiSend |      |     |             |
|     |      |                            | DGOptionDataRecv                | 60   | <-  |             |
|     |      |                            | GCSkillKeySend                  |      |     |             |
| <-  | F330 |                            | IOCP.DataSend                   |      |     |             |
|     |      |                            | DGAnsMuBotData                  | AE   |     |             |
|     |      |                            | GSProtocol.GCAnsMuBotData       |      |     |             |
| <-  | AE   |                            | IOCP.DataSend                   |      |     |             |
