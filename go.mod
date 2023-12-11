module github.com/gizmo-ds/misstodon

go 1.20

require (
	github.com/dop251/goja v0.0.0-20231027120936-b396bb4c349d
	github.com/duke-git/lancet/v2 v2.2.7
	github.com/go-resty/resty/v2 v2.10.0
	github.com/gorilla/websocket v1.5.1
	github.com/jinzhu/configor v1.2.2
	github.com/joho/godotenv v1.5.1
	github.com/labstack/echo/v4 v4.11.3
	github.com/pkg/errors v0.9.1
	github.com/rs/xid v1.5.0
	github.com/rs/zerolog v1.31.0
	github.com/stretchr/testify v1.8.4
	github.com/tidwall/buntdb v1.3.0
	github.com/urfave/cli/v2 v2.26.0
	golang.org/x/crypto v0.16.0
	golang.org/x/net v0.19.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dlclark/regexp2 v1.7.0 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/pprof v0.0.0-20231205033806-a5a03c77bf08 // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/tidwall/btree v1.4.2 // indirect
	github.com/tidwall/gjson v1.14.3 // indirect
	github.com/tidwall/grect v0.1.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/rtred v0.1.2 // indirect
	github.com/tidwall/tinyqueue v0.1.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/exp v0.0.0-20231206192017-f3f8817b8deb // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gorilla/websocket v1.5.1 => github.com/gizmo-ds/gorilla-websocket v0.0.0-20230212044710-0f26ab2a978a
