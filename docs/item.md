## Item

### Item frame

```
pack(1)
[[12]byte]
```

```
field[0]:  item index 0~255

field[1]:
bit0~bit1: addition attack/defense
bit2:      lucky flag
bit3~bit6: level 0~15
bit7:      skill flag

field[2]:  durability 0~255

field[3]:
bit0~bit5: excellent option
bit0:      Increases the amount of Mana received for hunting monsters
           Increases the amount of Zen received for hunting monsters
bit1:      Increases the amount of Life received for hunting monsters
           Defense Success Rate 10%
bit2:      Increase Attack (Wizardry) Speed 7
           Reflect Damage 5%
bit3:      Increases (Magic)Attack 2%
           Decreases Damage 4%
bit4:      Increase (Magic)Attack/Level =20
           Increase Maximum Mana 4%
bit5:      Excellent damage rate 10%
           Increase Maximum Life 4%
bit6:      addition attack/defense extension
           Now field[3].bit6 with field[1].bit0~bit1 may range as follow:
           000: addition 0
           001: addition 4
           010: addition 8
           011: addition 12
           100: addition 16
bit7:      item index extension
           Now field[3].bit7 with field[0] may range as 0~511

field[4]:  set

field[5]:
bit0:      period
bit1:      period expire
bit3:      option380 flag
bit4~bit7: item section 0~15

field[6]~field[11]: for socket/pentagram/muun
```

### Inventory

```
inventory position 0~11 and 236:
```

|                   |              |             |            |                     |
| ----------------- | ------------ | ----------- | ---------- | ------------------- |
| [8]: Imp          |              | [2]: Helmet |            | [7]: Wing           |
|                   | [9]: Pendant |             |            |                     |
| [0]: Primary hand |              | [3]: Armor  |            | [1]: Secondary hand |
|                   | [10]: Ring   |             | [11]: Ring |                     |
| [5]: Glove        |              | [4]: Pant   |            | [6]: Boot           |
|                   |              |             |            | [236]: Pentagram    |

```
inventory position 12~235:
```

#### Request

empty

#### Reply

```
pack(1)
[C4 [2]byte F3 10 n [13n]byte]
```

| Index   | Element       | Description                       |
| ------- | ------------- | --------------------------------- |
| 0       | 0xC4          | c1c2 frame flag                   |
| 1~2     | [2]byte       | c1c2 frame size: BE               |
| 3~4     | 0xF310        | c1c2 frame code: BE               |
| 5       | n             | inventory item count              |
| 5~5+13n | [(1+12)n]byte | inventory position and item frame |

### Move inventory item

#### Request

```
pack(1)
[C1 13 24 byte byte [12]byte byte byte]
```

| Index | Element  | Description          |
| ----- | -------- | -------------------- |
| 0     | 0xC1     | c1c2 frame flag      |
| 1     | 0x13     | c1c2 frame size      |
| 2     | 0x24     | c1c2 frame code      |
| 3     | byte     | source flag          |
| 4     | byte     | source position      |
| 5~16  | [12]byte | item frame           |
| 17    | byte     | destination flag     |
| 18    | byte     | destination position |

#### Reply

```
pack(1)
[C3 11 24 byte byte [12]byte]
```

| Index | Element  | Description                 |
| ----- | -------- | --------------------------- |
| 0     | 0xC1     | c1c2 frame flag             |
| 1     | 0x11     | c1c2 frame size             |
| 2     | 0x24     | c1c2 frame code             |
| 3     | byte     | result: 0=success -1=failed |
| 4     | byte     | destination position        |
| 5~16  | [12]byte | item frame                  |
