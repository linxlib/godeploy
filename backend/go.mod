module github.com/linxlib/godeploy

go 1.22.0

replace (
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.5.0
	github.com/linxlib/astp => ../../astp
	github.com/linxlib/fw => ../../../repos/fw
)

require (
	github.com/coreos/go-systemd v0.0.0-00010101000000-000000000000
	github.com/fasthttp/session/v2 v2.5.6
	github.com/glebarez/sqlite v1.11.0
	github.com/linxlib/conv v1.1.1
	github.com/linxlib/fw v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.3
	github.com/valyala/fasthttp v1.55.0
	gorm.io/gorm v1.25.11
)

require (
	atomicgo.dev/cursor v0.2.0 // indirect
	atomicgo.dev/keyboard v0.2.9 // indirect
	atomicgo.dev/schedule v0.1.0 // indirect
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fasthttp/router v1.5.2 // indirect
	github.com/fasthttp/websocket v1.5.10 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/godbus/dbus/v5 v5.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/gookit/filter v1.2.1 // indirect
	github.com/gookit/goutil v0.6.15 // indirect
	github.com/gookit/validate v1.5.2 // indirect
	github.com/jinzhu/configor v1.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/linxlib/astp v0.0.0-20240725151659-eeff57fd3522 // indirect
	github.com/linxlib/inject v0.1.3 // indirect
	github.com/lithammer/fuzzysearch v1.1.8 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/philhofer/fwd v1.1.3-0.20240612014219-fbbf4953d986 // indirect
	github.com/pterm/pterm v0.12.79 // indirect
	github.com/redis/go-redis/v9 v9.5.3 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/savsgio/gotils v0.0.0-20240704082632-aef3928b8a38 // indirect
	github.com/tinylib/msgp v1.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/term v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.23.1 // indirect
)
