///usr/bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/olahol/melody"
)

func handleCommands(p *Person) {
	fmt.Printf("Listening for commands for %v\n", p.name)
	for {
		command := <-p.infoStream.command
		fmt.Printf("%v: got command %v.\n", p.name, command)
		if len(command) == 0 {
			continue
		}
		switch command[0] {
		case 't':
			p.Log("t-shirt!")
			p.mu.Lock()
			p.insulation.Add(-0.1)
			p.mu.Unlock()
		case 'j':
			p.Log("toasty jumper!")
			p.mu.Lock()
			p.insulation.Add(0.1)
			p.mu.Unlock()
		case 'd':
			p.Log("glugging on a lot of 🚰!")
			p.mouth <- Water(100)
		case 'l':
			p.Log("Using a blood bag")
			p.mouth <- Blood(100)
		case 'w':
			p.Log("glugging on a bit of 🚰!")
			p.mouth <- Water(10)
		case 'c':
			p.Log("munching on some chocolate. Don't forget to clean your teeth!")
			p.mouth <- Chocolate()
		case 'P':
			p.Log("Crunching on some very dry 🥛 powder. It's making me thirsty!")
			p.mouth <- MilkPowder()
		case 'p':
			p.Log("munching and slurping on a container of juicy 🍑!")
			p.mouth <- Peaches()
		case 'm':
			p.Log("taking a swig of a magic potion!!!")
			p.mouth <- MagicPotion()
		case 'y':
			p.Log("first ryyyyyyyyyyyyyvita!")
			p.mouth <- Ryvita()
			p.Log("second ryyyyyyyyyyyyyvita!")
			p.mouth <- Ryvita()
		case 'e':
			p.Log("An ooey gooey yummy Jaffle. Just the way Simon likes them!")
			p.mouth <- Jaffle()
		case 'a':
			p.Log("ah, a crisp yummy apple. Just the thing I wanted!")
			p.mouth <- Apple()
		case 'b':
			p.Log("munching on a banana!")
			p.mouth <- Banana()
		case 'C':
			p.Log("brrrr!")
			p.mu.Lock()
			p.ambient_temperature.Add(-0.1)
			p.mu.Unlock()
		case 72: // H
			p.Log("phew!")
			p.mu.Lock()
			p.ambient_temperature.Add(0.1)
			p.mu.Unlock()
		case 'r':
			// run := Run{}
			// p.StartActivity(&run)
			p.mu.Lock()
			if p.stamina.TryAdd(-0.3) {
				p.Log("ruuuuunnn!")
			} else {
				p.Log("<pant!>")
			}
			p.mu.Unlock()
		case 's':
			// run := Run{}
			// p.StartActivity(&run)
			p.mu.Lock()
			if p.stamina.TryAdd(-0.1) {
				p.Log("sneaking quickly away from a zombie you spot in the distance.....!")
			} else {
				p.Log("time for a quick breather")
			}
			p.mu.Unlock()

		case 'z':
			p.Log("owch!")
			p.mu.Lock()
			damage := math.Max(0, 0.05+rand.NormFloat64()/30)
			if damage > 0 {
				p.Log(fmt.Sprintf("You are hurt, losing %.3f health\n", damage))
				p.health.Add(-damage)
				// p.notify(fmt.Sprintf("You are hurt, losing %.3f health\n", damage), "default")
				// p.notify(fmt.Sprint(p), "low")
			}
			if rand.Intn(10) >= 8 {
				p.blood.Add(-0.2)
			}
			p.mu.Unlock()
		}
	}
}

func getStyle(attr *Scalar) string {
	var value float64
	switch attr.profile {
	case OneGood:
		value = attr.current
	case HalfGood:
		value = 1 - 2*math.Abs(0.5-attr.current)
	case ZeroGood:
		value = 1 - attr.current
	case Undefined:
		value = 1
	}

	if value < 0.3 {
		return "danger"
	} else if value < 0.5 {
		return "warning"
	} else if value < 0.8 {
		return "ok"
	} else {
		return "good"
	}
}

type InfoStream struct {
	command     chan string
	information chan string
}

func start_websocket_server(domain string, port int) error {
	m := melody.New()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.css")
	})
	http.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.js")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	register_websockets(m)
	fmt.Printf("listening on port %v\n", port)
	// return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	httpsListener, err := makeHttpsListener(domain, port)
	if err != nil {
		return err
	}
	return http.Serve(httpsListener, nil)
}

func showPlayer(p *Person) {
	// p.infoStream.information <- fmt.Sprintf("%v: F%0.3f; W%0.3f; B%0.3f; H:%0.3f; S:%0.3f; T:%0.3f; A:%0.3f; I:%0.3f", p.name, p.food.current, p.water.current, p.blood.current, p.health.current, p.stamina.current, p.temperature.current, p.ambient_temperature.current, p.insulation.current)
	/*
		showValue(c, "Food", &p.food, 10, 2)
		showValue(c, "Water", &p.water, 20, 2)
		showValue(c, "Blood", &p.blood, 30, 2)
		showValue(c, "Health", &p.health, 40, 2)
		showValue(c, "Stamina", &p.stamina, 50, 2)
		showValue(c, "Temp.", &p.temperature, 60, 2)
		showValue(c, "Amb.Temp.", &p.ambient_temperature, 70, 2)
		showValue(c, "Insul.", &p.insulation, 80, 2)

		c.DrawText(1, 5, 30, 16, fmt.Sprint(p.pack))

		c.DrawText(40, 5, 120, 12, p.stomach)
	*/
	fmt.Println(p.stomach)
}

func update(name string, s *Scalar, c chan string) {
	previousValue := -9999.999
	for {
		newValue := s.current
		if math.Abs(previousValue-newValue) > 0.01 {
			// c <- fmt.Sprintf("STAT: %v %0.4f", name, newValue)
			c <- fmt.Sprintf("STAT %v: %v (%0.2f)", name, getStyle(s), newValue)
			previousValue = newValue
		}
		time.Sleep(time.Second)
	}
}

func main() {
	if err := start_websocket_server("dayz.dkyuf5dhjf.kozow.com", 8443); err != nil {
		fmt.Println(err)
	}
}
