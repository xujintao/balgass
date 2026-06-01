# Project Context

你现在位于 Go 项目 server-game 的仓库根目录：$HOME/balgass/src/server-game

另有 C++ 参考项目 GameServer：$HOME/balgass-igc/igc/9.5.1.15/source/GameServer

# 新业务配置文件加载规范

新增带配置文件的业务时，先判断配置归属：

- IGC 原始静态配置放在 `IGCData` 下，使用 `conf.PathCommon`。
- 项目自有配置放在 `PATH_COMMON` 下，使用 `conf.ServerEnv.PathCommon`。

使用包级 `init()` 调用业务管理器的 `init()` 方法，在服务启动时一次性加载配置。在加载方法内部定义贴合 XML 或 INI 文件结构的局部 DTO，通过 `conf.XML` 或 `conf.INI` 读取。DTO 只用于解析配置；业务代码使用转换后的 `map`、`slice`、固定数组或领域模型。

转换阶段补齐派生字段，并优先复用现有反序列化、绑定和计算逻辑。新配置必须校验索引边界、重复项、非法枚举、缺失引用和必填字段。核心静态配置错误应阻止启动；可选开发配置可以记录错误并跳过无效项。

标准示例：

```go
func init() {
    Manager.init()
}

func (m *manager) init() {
    type config struct {
        // XML or INI mapping fields
    }

    var cfg config
    conf.XML(basePath, "relative/path.xml", &cfg)

    // Validate and convert cfg into runtime lookup structures.
}
```

# 测试环境变量加载

执行任意测试前必须参考 `test.sh` 加载环境变量文件，否则 `conf` 包初始化会从当前目录查找 `GameServer.ini` 并失败。适用范围包括 `test.sh`、`go test`、`go test -race`、包测试、函数级测试和通配符测试。

必须先执行：

```bash
set -a
. "${HOME}/balgass/config/server-game/.env"
set +a
```

单条测试命令推荐写法：

```bash
set -a
. "${HOME}/balgass/config/server-game/.env"
set +a
<test command> > /tmp/server-game-test.log 2>&1
tail -n 120 /tmp/server-game-test.log
```

# 测试输出控制

执行测试时必须控制输出，避免大量配置加载日志消耗上下文 token。无论是 `test.sh`、race 测试、包测试，还是通配符/函数级测试，都不要直接把完整输出打印到对话上下文。

推荐做法：

```bash
<test command> > /tmp/server-game-test.log 2>&1
tail -n 120 /tmp/server-game-test.log
```

如果测试失败，再只提取关键错误信息：

```bash
rg -n "FAIL|panic|error|Error|fatal|Fatal|DATA RACE" /tmp/server-game-test.log
```

需要更完整诊断时，优先读取相关失败包、失败测试名、panic/race 附近的局部日志，而不是整份日志。
