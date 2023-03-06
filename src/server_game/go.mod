module github.com/xujintao/balgass/src/server_game

go 1.20

replace github.com/xujintao/balgass/src/c1c2 => ../c1c2

require (
	github.com/xujintao/balgass/src/c1c2 v0.0.0-00010101000000-000000000000
	gopkg.in/ini.v1 v1.67.0
)

require github.com/stretchr/testify v1.8.2 // indirect
