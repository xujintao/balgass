# Project Context

你现在位于 Go 项目 server-game 的仓库根目录：$HOME/balgass/src/server-game

另有 C++ 参考项目 GameServer：$HOME/balgass-igc/igc/9.5.1.15/source/GameServer

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
