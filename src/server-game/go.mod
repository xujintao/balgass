module github.com/xujintao/balgass/src/server-game

go 1.21

replace github.com/xujintao/balgass/src/c1c2 => ../c1c2

replace github.com/xujintao/balgass/src/utils => ../utils

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/go-playground/validator/v10 v10.15.5
	github.com/gorilla/websocket v1.5.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/xujintao/balgass/src/c1c2 v0.0.0-00010101000000-000000000000
	github.com/xujintao/balgass/src/utils v0.0.0-00010101000000-000000000000
	github.com/yuin/gopher-lua v1.1.1
	golang.org/x/text v0.14.0
	gopkg.in/ini.v1 v1.67.0
	gorm.io/driver/postgres v1.5.3
	gorm.io/gorm v1.25.5
)

require (
	github.com/bytedance/sonic v1.10.2 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.5.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
