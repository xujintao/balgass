module github.com/xujintao/balgass/src/server_game

go 1.20

replace github.com/xujintao/balgass/src/c1c2 => ../c1c2

replace github.com/xujintao/balgass/src/utils => ../utils

require (
	github.com/gorilla/websocket v1.5.0
	github.com/xujintao/balgass/src/c1c2 v0.0.0-00010101000000-000000000000
	github.com/xujintao/balgass/src/utils v0.0.0-00010101000000-000000000000
	gopkg.in/ini.v1 v1.67.0
	gorm.io/driver/postgres v1.5.3
	gorm.io/gorm v1.25.5
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/stretchr/testify v1.8.2 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)
