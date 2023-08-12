module nicois/dayz

go 1.18

replace nicois/bmon => ../../battery_monitor

require nicois/bmon v0.0.0-0000-0000

require (
	github.com/olahol/melody v1.1.1
	github.com/sirupsen/logrus v1.9.2
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.5.0 // indirect
)
