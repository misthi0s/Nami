module NamiC2

go 1.17

require (
	NamiDatabase v1.0.0
	NamiPrompt v1.0.0
	NamiUtilities v1.0.0
	NamiWorker v1.0.0
	github.com/elastic/go-windows v1.0.1
	github.com/go-ole/go-ole v1.3.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871
	golang.org/x/sys v0.6.0
)

require (
	NamiCommands v1.0.0 // indirect
	github.com/desertbit/closer/v3 v3.1.2 // indirect
	github.com/desertbit/columnize v2.1.0+incompatible // indirect
	github.com/desertbit/go-shlex v0.1.1 // indirect
	github.com/desertbit/grumble v1.1.3 // indirect
	github.com/desertbit/readline v1.5.1 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	github.com/rodaine/table v1.1.0 // indirect
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/tools v0.0.0-20201124115921-2c860bdd6e78 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	lukechampine.com/uint128 v1.1.1 // indirect
	modernc.org/cc/v3 v3.35.17 // indirect
	modernc.org/ccgo/v3 v3.12.65 // indirect
	modernc.org/libc v1.11.71 // indirect
	modernc.org/mathutil v1.4.1 // indirect
	modernc.org/memory v1.0.5 // indirect
	modernc.org/opt v0.1.1 // indirect
	modernc.org/sqlite v1.14.1 // indirect
	modernc.org/strutil v1.1.1 // indirect
	modernc.org/token v1.0.0 // indirect
)

replace NamiDatabase v1.0.0 => ./database/

replace NamiPrompt v1.0.0 => ./prompt/

replace NamiWorker v1.0.0 => ./workers/

replace NamiCommands v1.0.0 => ./commands/

replace NamiUtilities v1.0.0 => ./utilities/
