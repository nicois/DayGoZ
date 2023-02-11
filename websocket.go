package main

import (
	"fmt"
	"log"
	"regexp"
	"sync"

	"github.com/olahol/melody"
)

type app struct {
	player        *Person
	sessions      []*melody.Session
	sessionsMutex *sync.Mutex
	stats         map[string]string
	logs          []string
}

var players map[string]*app = make(map[string]*app)

func (a *app) sync(session *melody.Session) {
	fmt.Printf("Syncing...\n")
	for k, v := range a.stats {
		err := session.Write([]byte(fmt.Sprintf("STAT %v: %v", k, v)))
		if err != nil {
			fmt.Println(err)
		}
	}
	for _, l := range a.logs {
		err := session.Write([]byte(fmt.Sprintf("LOG: %v", l)))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Connect(username string, session *melody.Session) *app {
	reInfo := regexp.MustCompile(`^(\w+)(| \w+): (.+)`)
	a, found := players[username]
	fmt.Printf("Found player %v? %v\n", username, found)
	fmt.Printf("Connected to player %v from %v\n", username, session.RemoteAddr())
	if !found {
		player := &Person{name: username}
		player.Initialise()
		a = &app{player: player, stats: make(map[string]string), sessionsMutex: new(sync.Mutex)}
		go func(a *app) {
			for {
				message := <-a.player.infoStream.information
				matches := reInfo.FindStringSubmatch(message)
				if matches == nil {
					fmt.Printf("got unparseable message for player %v: %v\n", player.name, message)
				} else {
					switch matches[1] {
					case "LOG":
						if len(a.logs) > 10 {
							a.logs = a.logs[1:]
						}
						a.logs = append(a.logs, matches[3])
					case "STAT":
						stat := matches[2][1:]
						fmt.Printf("%v: stat is %v, value is %v\n", player.name, stat, matches[3])
						if old_value, found := a.stats[stat]; found && old_value == matches[3] {
							fmt.Printf("same value, not re-sending: %v\n", message)
							continue
						}
						a.stats[stat] = matches[3]
					default:
						fmt.Printf("got matches: %q\n", matches)
						fmt.Printf("got unparseable type for player %v: %v\n", player.name, matches[1])
						continue
					}
					a.sessionsMutex.Lock()
					atLeastOneSessionIsClosed := false
					for _, s := range a.sessions {
						if s.IsClosed() {
							atLeastOneSessionIsClosed = true
						} else {
							err := s.Write([]byte(message))
							if err != nil {
								fmt.Println(err)
							}
						}
					}
					if atLeastOneSessionIsClosed {
						for index, s := range a.sessions {
							if s.IsClosed() {
								a.sessions = append(a.sessions[:index], a.sessions[index+1:]...)
								fmt.Println("removed a closed session.")
								break
							}
						}
					}
					a.sessionsMutex.Unlock()
				}
			}
		}(a)
		players[username] = a
	}
	a.sync(session)
	a.sessionsMutex.Lock()
	defer a.sessionsMutex.Unlock()
	a.sessions = append(a.sessions, session)
	return a
}

func (a *app) Command(message string) error {
	a.player.infoStream.command <- message
	return nil
}

type App interface {
	Command(message string) error
	Disconnect()
}

func register_websockets(m *melody.Melody) {
	m.HandleConnect(func(s *melody.Session) {
		log.Println("client connected")
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Println("client disconnected")
	})

	reAuth := regexp.MustCompile(`AUTH (.+)`)

	m.HandleMessage(func(s *melody.Session, messageBytes []byte) {
		message := string(messageBytes)
		fmt.Printf("got message %v\n", message)
		value, exists := s.Get("app")
		if !exists {
			if matches := reAuth.FindStringSubmatch(message); matches != nil {
				fmt.Printf("got matches: %q\n", matches)
				app := Connect(matches[1], s)
				s.Set("app", app)
				s.Write([]byte("Logged in successfully"))
			} else {
				s.Write([]byte("unrecognised command (while not logged in)"))
			}
			return
		}
		a := value.(*app)
		fmt.Printf("sending command %v\n", message)
		a.Command(message)
	})
}
