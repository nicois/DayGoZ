module nicois/dayz

go 1.18

replace nicois/bmon => ../../battery_monitor

require nicois/bmon v0.0.0-0000-0000

require (
	github.com/gdamore/tcell/v2 v2.5.4
	github.com/olahol/melody v1.1.1
	golang.org/x/time v0.3.0
	nhooyr.io/websocket v1.8.7
)

require (
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/klauspost/compress v1.10.3 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.5.0 // indirect
)
