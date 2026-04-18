# 6. 地图系统

本模块覆盖地图数据、地形属性、站位占用、随机坐标、路径与墙体、Gate、MoveCommand、小地图、地图物品、天气/经验/Zen、跨服地图归属与地图侧出生/刷怪点位。对象移动、视野、攻击、死亡归对象系统；怪物 AI 主体归 `25-monster-ai.md`；移动作弊、穿墙审计、异常坐标处罚归 `27-security.md`；跨服迁移请求、地图服认证和外部服状态归 `28-external-comm.md`；本模块只提供地图规则和路径能力。

| 序号 | 子模块分类 | 子模块 | GameServer | server-game | 状态 | TODO 说明 |
|---:|---|---|---|---|---|---|
| 1 | 地图注册与加载 | 地图编号常量 | `MapClass.h`、地图宏 | `game/maps/const.go` | 已覆盖 | 保持 Go 常量与客户端/配置地图编号一致，后续新增地图先补常量。 |
| 2 | 地图注册与加载 | 最大地图数量边界 | `MAX_NUMBER_MAP`、`MapNumberCheck` | `MAX_MAP_NUMBER` | 部分覆盖 | 需要统一所有入口的 map number 越界检查，避免直接索引 panic。 |
| 3 | 地图注册与加载 | MapManager 启动初始化 | `MapClass::LoadMapAttr` 调用链 | `MapManager.init` | 已覆盖 | 当前通过包 `init` 加载，后续测试可考虑注入式初始化。 |
| 4 | 地图注册与加载 | MapList 配置读取 | `MapClass::LoadMapAttr`、服务启动加载 | `IGC_MapList.xml` 读取 | 已覆盖 | 保持从 `DefaultMaps` 加载普通地图，事件地图属性另行处理。 |
| 5 | 地图注册与加载 | 地图文件名解析 | `MapClass::AttrLoad` 文件路径 | `regexp` 提取 `.att` 名称 | 已覆盖 | 当前从文件名生成地图名，需验证与语言包/客户端名称是否一致。 |
| 6 | 地图注册与加载 | 重复地图名处理 | GameServer 靠地图编号区分 | `mapNameCount` 追加序号 | 部分覆盖 | 重名地图如多层地图要有稳定命名规则，避免调试订阅混乱。 |
| 7 | 地图注册与加载 | 名称到编号索引 | Move/调试订阅地图语义 | `mapNameNumber`、`GetMapNumber` | 已覆盖 | `GetMapNumber` 当前无不存在判断，建议返回 `(int,bool)` 或调用端校验。 |
| 8 | 地图注册与加载 | 编号到名称索引 | GameServer 日志/消息地图名 | `mapNumberName`、`GetMapName` | 已覆盖 | 后续用于日志、GM 命令和 Web 调试界面。 |
| 9 | 地图注册与加载 | 地图文件路径拼接 | `Data/MapTerrains/*.att` | `path.Join(conf.PathCommon, "MapTerrains", v.File)` | 已覆盖 | 需要在启动日志中输出缺失文件路径，便于部署排错。 |
| 10 | 地图注册与加载 | 地图加载失败策略 | `g_Log.MsgBox`、`ExitProcess` | `slog.Error` 后 `os.Exit(1)` | 已覆盖 | 地图属于强依赖，启动失败退出合理；测试环境需可替换。 |
| 11 | 地形属性文件 | `.att` 文件头读取 | `MapClass::AttrLoad` | `_map.init` 读取 `buf[0]` | 部分覆盖 | 当前未显式校验文件长度和 head，需补坏文件测试。 |
| 12 | 地形属性文件 | 地图宽高读取 | `m_width/m_height` | `width=int(buf[1])+1`、`height=int(buf[2])+1` | 已覆盖 | 需确认 `+1` 与 GameServer 读取语义是否完全对齐。 |
| 13 | 地形属性文件 | 属性缓冲区初始化 | `m_attrbuf` | `_map.buf = buf[3:]` | 已覆盖 | Go 侧直接持有切片，站位 bit 会修改原 buffer。 |
| 14 | 地形属性文件 | 阻挡点缓存 | `GetAttr` 遍历语义 | `_map.pots` 收集 `attr&4` | 部分覆盖 | 主要用于 Web 调试地图，需明确不参与正式逻辑。 |
| 15 | 地形属性文件 | 坐标合法性检查 | `GetAttr` 边界分支 | `_map.valid` | 已覆盖 | 后续所有地图入口应先走 `valid` 或安全封装。 |
| 16 | 地形属性文件 | 坐标转索引 | `y*256+x`、`TERRAIN_INDEX_REPEAT` | `_map.pos2index` | 部分覆盖 | Go 侧用 `m.width`，GameServer 多处固定 256；需确认非 256 地图兼容。 |
| 17 | 地形属性文件 | 越界属性返回阻挡 | `GetAttr` 返回 `4` | `_map.getAttr` 返回 `4` | 已覆盖 | 这是安全默认值，保持所有越界为不可走。 |
| 18 | 地形属性文件 | 属性值读取接口 | `MapClass::GetAttr` | `MapManager.GetMapAttr` | 已覆盖 | 当前调用方较多，需避免 nil map 直接索引。 |
| 19 | 地形属性文件 | 调试地图点导出 | GameServer GM/日志辅助 | `GetMapPots`、`GetMapStands` | 部分覆盖 | 当前用于 Web 订阅调试，后续可扩展安全区/出生点可视化。 |
| 20 | 地形属性文件 | 地图属性加载测试 | GameServer 无单测 | `game/maps/map_test.go` 为空 | 未覆盖 | 需要补最小 `.att` fixture，覆盖宽高、阻挡、越界、站位 bit。 |
| 21 | 地图属性规则 | 安全区 bit 解析 | `attr&1` 相关判断 | `GetMapAttr` 调用方判断 `attr&1` | 部分覆盖 | 安全区影响攻击、击退、PK、掉落，应统一封装 `IsSafeZone`。 |
| 22 | 地图属性规则 | 站位 bit 解析 | `attr&2` | `CheckMapAttrStand`、`Set/ClearMapAttrStand` | 已覆盖 | 站位 bit 是运行态属性，和静态 `.att` 共用 buffer。 |
| 23 | 地图属性规则 | 障碍 bit 解析 | `attr&4` | `checkNoWall`、`getAttr` | 已覆盖 | 障碍用于寻路、视线、随机坐标和移动中断。 |
| 24 | 地图属性规则 | bit8 语义核对 | `attr&8` 在 `GetStandAttr` 中不可站 | 多处判 `attr&8` | 部分覆盖 | 需要确认 bit8 在当前版本是墙、水、禁区还是事件属性。 |
| 25 | 地图属性规则 | bit16 语义核对 | GameServer 部分特殊地图属性 | `Knockback` 判 `attr&16` | 部分覆盖 | Go 侧只有击退使用 bit16，需确认是否应进入通用站位/寻路规则。 |
| 26 | 地图属性规则 | PvPConfig 读取 | `CMapAttribute::LoadFile` | `IGC_MapAttribute.xml` 读取字段未保存 | 未覆盖 | 需要把 PvPConfig 存入 `_map`，供 PK/攻击规则查询。 |
| 27 | 地图属性规则 | ItemDropRateBonus 读取 | `CMapAttribute` | 字段读取但未保存 | 未覆盖 | 需要接入怪物掉落倍率，不应只读取经验倍率。 |
| 28 | 地图属性规则 | ExpBonus 读取 | `CMapAttribute` | `_map.expBonus` | 已覆盖 | 当前用于 `DieGiveExperience`，需与经验系统动态倍率合并。 |
| 29 | 地图属性规则 | MasterExpBonus 读取 | `CMapAttribute` | `_map.masterExpBonus` | 已覆盖 | 当前大师经验加成已接入基础杀怪经验。 |
| 30 | 地图属性规则 | BlockEntry/RegenOnSamePlace 等属性 | `CMapAttribute::LoadFile` | 字段读取但未保存 | 未覆盖 | 需要保存并提供查询接口，供传送、复活、事件地图进入限制使用。 |
| 31 | 站位占用 | 可站立检查 | `MapClass::GetStandAttr` | `_map.checkAttrStand` | 已覆盖 | Go 侧返回 false 覆盖站位、阻挡、bit8；需补安全区是否可站规则。 |
| 32 | 站位占用 | 设置站位 | `MapClass::SetStandAttr` | `_map.setAttrStand` | 已覆盖 | 进入地图、开始移动、怪物出生都依赖该接口。 |
| 33 | 站位占用 | 清除站位 | `MapClass::ClearStandAttr` | `_map.clearAttrStand` | 已覆盖 | 离开地图、死亡、传送和移动起点释放都要调用。 |
| 34 | 站位占用 | 移动目标占用 | `gObjMoveProc`/MapClass 站位语义 | `Object.Move` 设置 `TX/TY` 后占用 | 已覆盖 | 需要校验客户端路径目标不可站时拒绝，而不是先占位。 |
| 35 | 站位占用 | 移动起点释放 | GameServer 移动时清理旧站位 | `Object.Move` 清理 `X/Y` 与旧 `TX/TY` | 部分覆盖 | 当前连续清理两个坐标，需防止清掉其他对象占用。 |
| 36 | 站位占用 | 传送站位切换 | `gObjMoveGate` | `Object.gateMove` | 已覆盖 | 先清旧图站位，再切图设置新站位；需要失败回滚测试。 |
| 37 | 站位占用 | 死亡站位释放 | `MapClass::ClearStandAttr` 调用链 | `attack.go` 清理 `tobj.TX/TY` | 已覆盖 | 怪物死亡后已清理，玩家死亡/复活路径也要统一。 |
| 38 | 站位占用 | 调试 stands 表 | GameServer 无等价核心结构 | `conf.ServerEnv.Debug` 下 `stands` | 部分覆盖 | 只应作为调试输出，不能成为业务真相。 |
| 39 | 站位占用 | 站位并发边界 | GameServer 单线程对象循环 | Go 单游戏协程/调用方约束 | 部分覆盖 | `MapManager` 没锁，必须明确只能在游戏主循环时序内写。 |
| 40 | 站位占用 | 最近可站点搜索 | `MapClass::SearchStandAttr` | 当前无等价函数 | 未覆盖 | 传送/出生落到阻挡点时应搜索附近可站点，而不是直接回退原点。 |
| 41 | 随机坐标与出生点 | 区域随机坐标 | `CGate::GetGate`、`GetBoxPosition` | `_map.getRandomPos` | 已覆盖 | 当前重试 100 次，比 GameServer Gate 10 次更激进，可接受但需测试。 |
| 42 | 随机坐标与出生点 | 单点区域处理 | `CGate::GetGate` 起止相等处理 | `w/h <=0` 回退为 1 | 已覆盖 | 固定点 Gate/怪物出生可正常返回起点。 |
| 43 | 随机坐标与出生点 | 随机坐标阻挡过滤 | `attr&2/4/8` | `_map.getRandomPos` | 已覆盖 | 需补 bit16 是否也应过滤。 |
| 44 | 随机坐标与出生点 | 随机坐标失败回退 | `CGate::GetGate` 最后坐标 | `_map.getRandomPos` 返回 `x1,y1` | 部分覆盖 | 如果 `x1,y1` 不可站，后续应结合最近可站点搜索。 |
| 45 | 随机坐标与出生点 | 怪物区域出生点解析 | `CMonsterSetBase::LoadSetBase` | `SpawnMonster` 读取 `IGC_MonsterSpawn.xml` | 部分覆盖 | Go 已读区域，但缺 GameServer MonsterSetBase 的完整类型语义。 |
| 46 | 随机坐标与出生点 | 怪物固定点出生 | `CMonsterSetBase::GetPosition` | `newMonster`、`SpawnPosition` | 部分覆盖 | 固定点和区域点都应通过统一可站检查。 |
| 47 | 随机坐标与出生点 | 怪物区域随机出生 | `CMonsterSetBase::GetBoxPosition` | `Monster.RandPosition` 调用地图随机 | 部分覆盖 | 需确认 `RandPosition` 是否完全过滤站位、墙和事件属性。 |
| 48 | 随机坐标与出生点 | NPC/商店点位出生 | `MonsterSetBase` NPC 类型 | `shop.ForEachShop` 创建 NPC | 部分覆盖 | 商店 NPC 已生成，但一般 NPC 表和对话 NPC 需要扩展。 |
| 49 | 随机坐标与出生点 | Pentagram 属性出生参数 | `GetPentagramMainAttribute` | `newMonster(... element)` | 部分覆盖 | Go 已保留元素字段，需和怪物配置/地图点位来源对齐。 |
| 50 | 随机坐标与出生点 | 定时刷怪位置分配 | `CMonsterRegenSystem::SetPosMonster` | 当前无独立定时刷怪系统 | 未覆盖 | 后续事件/世界 Boss 需要基于地图区域分配出生点。 |
| 51 | 路径与墙体 | 八方向表 | `PATH_t` 方向语义 | `maps.Dirs` | 已覆盖 | 与客户端方向编码一致，是移动和怪物追击基础。 |
| 52 | 路径与墙体 | 方向计算 | GameServer 方向计算工具 | `CalcDir` | 已覆盖 | 远程攻击转向、怪物追击、移动包生成都依赖。 |
| 53 | 路径与墙体 | 客户端路径解析 | `CGMoveProc` 路径包语义 | `MsgMove.Unmarshal` | 已覆盖 | 当前按 nibble 解析方向，需要非法方向保护。 |
| 54 | 路径与墙体 | 路径长度上限 | GameServer 路径数组限制 | `Object.Move`、`findPath` 限制 15 | 已覆盖 | 保持服务端拒绝超长路径，防止异常包。 |
| 55 | 路径与墙体 | 路径连续性校验 | `MoveCheck`、`PathFinding*` | 当前主要信任客户端路径 | 未覆盖 | 应校验每步相邻、方向合法、终点可站。 |
| 56 | 路径与墙体 | AStar 寻路 | `MapClass::PathFinding4` | `_path.findPathAStar` | 已覆盖 | 主要用于怪物移动，已有单测/benchmark。 |
| 57 | 路径与墙体 | BFS 寻路保留 | `PathFinding2/3` 多版本 | `_path.findPathBFS` | 已覆盖 | 作为对照实现保留，不默认启用。 |
| 58 | 路径与墙体 | 贪心寻路保留 | GameServer path fallback | `_path.findPath` | 已覆盖 | 当前保留测试，可用于性能/路径差异对比。 |
| 59 | 路径与墙体 | 直线墙体检查 | `MapClass::CheckWall` | `_map.checkNoWall` | 已覆盖 | 用于怪物视线/远程攻击前置，需补边界越界测试。 |
| 60 | 路径与墙体 | CheckWall2 语义 | `MapClass::CheckWall2` | 当前无返回细分原因 | 未覆盖 | 如果技能需要区分墙体/越界/可见性，需补返回码版本。 |
| 61 | 对象地图移动入口 | 移动状态机 | `gObjMoveProc` | `Object.Move`、`processMove` | 部分覆盖 | 基础状态已覆盖，但需接入 `27-security.md` 的路径合法性和速度作弊校验。 |
| 62 | 对象地图移动入口 | 移动步进时间 | GameServer 对象 tick | `processMove` 400ms 基准 | 部分覆盖 | 需和客户端/攻速/移动状态对齐，避免位置漂移。 |
| 63 | 对象地图移动入口 | 斜向移动耗时 | GameServer 斜向步进语义 | `pathDir%2==0` 乘 1.3 | 已覆盖 | 保留当前规则，后续用实际客户端表现验证。 |
| 64 | 对象地图移动入口 | delayLevel 移动延迟 | `m_DelayLevel` 等状态 | `delayLevel` 加 300ms | 部分覆盖 | 需要明确来源和持续时间，避免状态永不恢复。 |
| 65 | 对象地图移动入口 | 移动中阻挡中断 | `MoveCheck`/MapClass | `processMove` 检查 `attr&4 && attr&8` | 需修正 | 这里使用 `&&` 可疑，通常任一阻挡属性就应中断。 |
| 66 | 对象地图移动入口 | 起点坐标同步 | GameServer 校验客户端坐标 | `Object.Move` 直接设置 `X/Y` | 需修正 | 当前过度信任客户端起点，需限制与服务端当前位置距离。 |
| 67 | 对象地图移动入口 | 目标坐标维护 | `TX/TY` 目标点 | `Object.Move` 更新 `TX/TY` | 已覆盖 | 目标点用于站位占用和广播。 |
| 68 | 对象地图移动入口 | 强制坐标设置 | `gObjSetPosition` 类语义 | `Object.SetPosition` | 已覆盖 | 用于击退/特殊技能，需统一可站校验。 |
| 69 | 对象地图移动入口 | 击退坐标校验 | Knockback 技能语义 | `Object.Knockback` | 部分覆盖 | 已检查多种 attr bit，需确认安全区是否允许击退。 |
| 70 | 对象地图移动入口 | 地图切换视锥刷新 | `gObjViewportListProtocolDestroy/Create` | `CreateFrustum`、viewport 处理 | 部分覆盖 | `gateMove` 更新视锥，但跨图时旧视野销毁/新视野创建需完整验证。 |
| 71 | Gate 传送 | Gate 配置加载 | `CGate::Load` | `GateMoveManager.init` | 已覆盖 | 读取 `Warps/IGC_GateSettings.xml`，字段基本齐全。 |
| 72 | Gate 传送 | Gate 索引表 | `m_This[]`、`m_*[]` | `gateMoveTable` | 已覆盖 | Map 结构更适合 Go，需对不存在 Gate 返回错误原因。 |
| 73 | Gate 传送 | Target Gate 解析 | `m_TargetGate` | `gateMove.Target` | 已覆盖 | `Target != 0` 时跳转到目标 Gate。 |
| 74 | Gate 传送 | 目标地图解析 | `m_MapNumber[gt]` | `gateMove.MapNumber` | 已覆盖 | 需要确认随机坐标属性检查使用的是目标地图而不是源地图。 |
| 75 | Gate 传送 | 目标方向解析 | `m_Dir[gt]` | `gateMove.Direction` | 已覆盖 | 下发 `MsgTeleportReply.Dir`。 |
| 76 | Gate 传送 | Gate 最小等级 | `CGate::GetLevel`、`CheckGateLevel` | 字段读取但 `gateMove` 未校验 | 未覆盖 | 直接 Teleport 应校验 Gate MinLevel，而不仅 MapMove 校验。 |
| 77 | Gate 传送 | Gate 是否存在检查 | `CGate::IsGate` | `GateMoveManager.Move` map lookup | 部分覆盖 | Go 当前静默 return，协议层应返回失败结果。 |
| 78 | Gate 传送 | Gate 区域内检查 | `CGate::IsInGate` | 当前无等价实现 | 未覆盖 | 走门/入口触发需要校验玩家是否在 Gate 起始区域附近。 |
| 79 | Gate 传送 | Gate 等级检查协议 | `CGate::CheckGateLevel` | `handle` 有 `checkGateLevel` 占位入口 | 未覆盖 | 需要实现 0xD00A 对应逻辑并返回客户端期望格式。 |
| 80 | Gate 传送 | Teleport 落地流程 | `gObjMoveGate` | `Object.Teleport`、`gateMove` | 部分覆盖 | 已完成基础切图坐标，缺事件限制、战斗状态限制、失败回复。 |
| 81 | MoveCommand | MoveReq 配置加载 | `CMoveCommand::Load` | `MapMoveManager.init` | 已覆盖 | 读取 `Warps/IGC_MoveReq.xml`，字段包括等级、钱、Gate。 |
| 82 | MoveCommand | MoveIndex 映射 | `FindIndex(int)` | `mapMoveTable[index]` | 已覆盖 | Go 直接用客户端 MoveIndex 查表。 |
| 83 | MoveCommand | 命令名映射 | `FindIndex(LPSTR)` | 当前无名称入口 | 未覆盖 | 如果实现聊天 `/move`，需要支持 ServerName/ClientName。 |
| 84 | MoveCommand | GateNumber 关联 | `GateNumber` | `mapMove.GateNumber` | 已覆盖 | MapMove 最终复用 Gate 传送。 |
| 85 | MoveCommand | 最小等级检查 | `NeedLevel` | `Object.MapMove` | 已覆盖 | 只覆盖普通等级，未覆盖职业折扣和大师/活动特殊规则。 |
| 86 | MoveCommand | 最大等级检查 | `MaxLevel` | 字段读取但未使用 | 未覆盖 | 需要用于部分活动/低级地图限制。 |
| 87 | MoveCommand | Zen 费用检查和扣除 | `NeedZen`、`GCMoneySend` | `Object.MapMove`、`PushMoney` | 已覆盖 | 成功后扣钱并推送钱，失败结果码需细化。 |
| 88 | MoveCommand | 职业等级折扣 | DL/MG/RF `2/3` 规则 | 当前无 | 未覆盖 | 对齐 `CMoveCommand::Move` 和 `GetMoveLevel`。 |
| 89 | MoveCommand | 装备限制 | `CheckEquipmentToMove` | 当前无 | 未覆盖 | Atlans/Icarus 等地图需要坐骑/翅膀限制。 |
| 90 | MoveCommand | 事件/接口限制 | `CheckMainToMove`、`CheckInterfaceToMove` | 当前无 | 未覆盖 | 交易、商店、活动地图、PK/Gens 状态都应限制移动。 |
| 91 | 小地图 | MiniMap 配置加载 | `CMinimapData::LoadFile` | `MiniManager.init` | 已覆盖 | 读取 `IGC_MiniMap.xml` 的 TypeOne/TypeTwo。 |
| 92 | 小地图 | NPC 点位表 | `MINIMAP_DATA` | `npcTable` | 已覆盖 | NPC 标记按 MapNumber 聚合。 |
| 93 | 小地图 | 入口点位表 | `MINIMAP_DATA` | `entranceTable` | 已覆盖 | 入口标记按 MapNumber 聚合。 |
| 94 | 小地图 | NPC 点位遍历 | `CMinimapData::SendData` | `ForEachMapNpc` | 已覆盖 | 当前 callback 下发，适合 Go 风格。 |
| 95 | 小地图 | 入口点位遍历 | `CMinimapData::SendData` | `ForEachMapEntrance` | 已覆盖 | 与 NPC 点位分离，便于不同图标类型。 |
| 96 | 小地图 | 小地图下发入口 | `SendData(aIndex)` | `Object.LoadMiniMap` | 已覆盖 | 地图切换时调用，登录初始地图也应确保调用。 |
| 97 | 小地图 | MiniMap 协议编码 | `PMSG_MINIMAP_DATA` | `MsgMiniMapReply.Marshal` | 已覆盖 | 已编码 ID、类型、坐标、名称。 |
| 98 | 小地图 | 小地图名称编码 | GameServer 多语言文本 | GBK 编码名称 | 已覆盖 | 需要处理名称过长和编码失败策略。 |
| 99 | 小地图 | 组队成员小地图 | `CMinimapData::SendPartyData` | 当前无等价实现 | 未覆盖 | 社交/组队模块落地后补队友坐标同步。 |
| 100 | 小地图 | 地图切换刷新小地图 | GameServer 地图切换后刷新 | `gateMove` 中 `LoadMiniMap` | 部分覆盖 | 需要避免同图传送重复发送大量小地图数据。 |
| 101 | 地图物品与 Zen | 地图物品槽初始化 | `MapClass::ItemInit` | `_map.inventory` | 已覆盖 | 槽数量来自 `MaxObjectItemCount`。 |
| 102 | 地图物品与 Zen | 地图物品创建 | `MapClass::ItemDrop`、`CMapItem::CreateItem` | `MapManager.AddItem`、`_map.addItem` | 部分覆盖 | 基础容器已完成，缺掉落归属、来源、可见时间。 |
| 103 | 地图物品与 Zen | 地图物品坐标校验 | `ItemDrop` 坐标规则 | `_map.addItem` 只校验 `valid` | 需修正 | 应过滤阻挡/站位/安全区限制，而不是只检查越界。 |
| 104 | 地图物品与 Zen | 地图物品查看 | `MapClass::m_cItem[item_num]` | `PeekItem`、`MapItem` | 已覆盖 | 拾取前可读取地图物品。 |
| 105 | 地图物品与 Zen | 地图物品删除 | `StateSetDestroy`、拾取删除 | `RemoveItem`、`removeItem` | 已覆盖 | 删除时还应广播视野内物品消失。 |
| 106 | 地图物品与 Zen | 地图物品过期 | `MapClass::StateSetDestroy` | `ExpireItem`、`expireItem` | 部分覆盖 | 当前固定 1 分钟，需读取配置或对齐 GameServer 状态时间。 |
| 107 | 地图物品与 Zen | 地图物品遍历 | `MapClass` 遍历物品槽 | `MapEachItem`、`eachItem` | 已覆盖 | 视野创建时用于下发地面物品。 |
| 108 | 地图物品与 Zen | 拾取/丢弃调用关系 | `ItemGive`、`ItemDrop` | `Object.PickItem`、`DropItem` | 部分覆盖 | 地图系统只提供容器；背包容量、权限、协议结果在道具系统处理。 |
| 109 | 地图物品与 Zen | ZenDrop 配置读取 | `MoneyItemDrop`、Zen 配置 | `IGC_ZenDrop.xml` | 已覆盖 | 已读取 Enable、Multiply、地图最小最大 Zen。 |
| 110 | 地图物品与 Zen | 地图 Zen 计算 | `MapClass::MoneyItemDrop` | `GetZen`、`_map.GetZen` | 部分覆盖 | 当前返回 `money*2` 可疑，需要按 GameServer/配置核对。 |
| 111 | 跨服/事件地图归属 | MapServerInfo 读取 | `CMapServerManager::LoadData` | `conf.MapServers` 已解析 | 部分覆盖 | 配置结构存在，但地图系统未提供业务查询。 |
| 112 | 跨服/事件地图归属 | 服务器组映射 | `_MAPSVR_DATA` | 当前无等价运行表 | 未覆盖 | 需要建立 server code、group、map list 的 Go 运行态索引。 |
| 113 | 跨服/事件地图归属 | 地图可移动检查 | `CheckMapCanMove` | 当前无 | 未覆盖 | MapMove/Gate 前应判断目标地图是否归当前服或可跨服。 |
| 114 | 跨服/事件地图归属 | 跨服目标解析 | `CheckMoveMapSvr` | 当前无 | 未覆盖 | 目标不在当前服时由地图系统解析目标服，再通过 `28-external-comm.md` 发起跨服移动协议。 |
| 115 | 跨服/事件地图归属 | 跨服 IP/端口查询 | `GetSvrCodeData` | 当前无 | 未覆盖 | 地图系统提供配置查询，跨服认证和外部通信由 `28-external-comm.md` 执行。 |
| 116 | 跨服/事件地图归属 | 活动地图进入限制 | `BC/CC/DS/IT/Kanturu` 范围判断 | 多数 handle 入口占位 | 未覆盖 | 先在地图系统提供通用限制查询，再由活动模块补具体状态。 |
| 117 | 跨服/事件地图归属 | Kalima 免费移动 | `MoveFree2Kalima`、`GetKalimaGateLevel` | 当前无 | 未覆盖 | 需要按职业等级计算 Kalima floor 并走对应 Gate。 |
| 118 | 跨服/事件地图归属 | MonsterRegen 配置加载 | `CMonsterRegenSystem::LoadScript` | 当前无 | 未覆盖 | 世界 Boss/事件怪定时刷新需要独立配置表。 |
| 119 | 跨服/事件地图归属 | MonsterRegen 定时运行 | `Run`、`IsRegenTime`、`RegenMonster` | 当前无 | 未覆盖 | 后续应挂到游戏 tick，按组刷新并记录存活状态。 |
| 120 | 跨服/事件地图归属 | 地图系统回归测试总集 | GameServer 靠运行验证 | 当前仅 path 测试 | 未覆盖 | 建议覆盖地图加载、属性、站位、随机点、Gate、MapMove、物品过期和 Zen。 |
