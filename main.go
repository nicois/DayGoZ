///usr/bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	tcell "github.com/gdamore/tcell/v2"
)

func (c *Cell) Shutdown() {
	c.screen.Fini()
}

func CreateCell() *Cell {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Set default text style
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	// Clear screen
	s.Clear()
	return &Cell{style: defStyle, screen: s}
}

func (c *Cell) DrawTextWithStyle(x1, y1, x2, y2 int, text string, style tcell.Style) {
	row := y1
	col := x1
runes:
	for _, r := range []rune(text) {
		if row > y2 {
			break
		}
		if r == '\n' {
			for x := col; x < x2; x++ {
				c.screen.SetContent(x, row, ' ', nil, c.style)
			}
			col = x1
			row++
			continue runes
		}
		if col < x2 {
			c.screen.SetContent(col, row, r, nil, style)
			col++
		}
	}
	for y := row; y < y2; y++ {
		for x := col; x < x2; x++ {
			c.screen.SetContent(x, y, ' ', nil, c.style)
		}
		col = x1
	}
	c.screen.Sync()
}

func (c *Cell) DrawText(x1, y1, x2, y2 int, text string) {
	c.DrawTextWithStyle(x1, y1, x2, y2, text, c.style)
}

/*
type Screen interface {
	DrawText(x1, y1, x2, y2 int, text string)
	DrawTextWithStyle(x1, y1, x2, y2 int, text string, style tcell.Style)
	Shutdown()
}
*/

type Cell struct {
	style  tcell.Style
	screen tcell.Screen
	lines  []string
	w      int
	h      int
}

var (
	dangerStyle  tcell.Style = tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorReset)
	warningStyle tcell.Style = tcell.StyleDefault.Foreground(tcell.ColorOrange).Background(tcell.ColorReset)
	okStyle      tcell.Style = tcell.StyleDefault.Foreground(tcell.ColorReset).Background(tcell.ColorReset)
	goodStyle    tcell.Style = tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorReset)
)

func (c *Cell) Log(s string) {
	if len(c.lines) > 5 {
		c.lines = c.lines[1:]
	}
	c.lines = append(c.lines, fmt.Sprintf("%v: %v", time.Now().Format(time.UnixDate), strings.TrimSpace(s)))

	c.DrawText(1, c.h-6, c.w, c.h, strings.Join(c.lines, "\n"))
}

func kb(p *Person, cell *Cell) {
	for {
		// Update screen
		cell.screen.Show()

		// Poll event
		ev := cell.screen.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			cell.screen.Sync()
			cell.w, cell.h = ev.Size()
			cell.screen.Clear()
		case *tcell.EventKey:

			switch ev.Key() {
			case tcell.KeyEscape:
				fmt.Println("ESC pressed")
				cell.Shutdown()
				os.Exit(0)
			case tcell.KeyRune:
				switch ev.Rune() {
				case 't':
					cell.Log("t-shirt!")
					p.mu.Lock()
					p.insulation.Add(-0.1)
					p.mu.Unlock()
				case 'j':
					cell.Log("toasty jumper!")
					p.mu.Lock()
					p.insulation.Add(0.1)
					p.mu.Unlock()
				case 'w':
					cell.Log("glugging on some water!")
					p.mouth <- Water(10)
				case 'b':
					cell.Log("munching on a banana!")
					p.mouth <- Banana()
				case 'C':
					cell.Log("brrrr!")
					p.mu.Lock()
					p.ambient_temperature.Add(-0.1)
					p.mu.Unlock()
				case 72: // H
					cell.Log("phew!")
					p.mu.Lock()
					p.ambient_temperature.Add(0.1)
					p.mu.Unlock()
				case 'd': // d
					cell.Log("gluuuuug!")
					p.mu.Lock()
					p.water.Add(0.3)
					p.mu.Unlock()
				case 'e': // e
					// fmt.Println("yuuuuunm!")
					cell.Log("yuuuuunm!")
					p.mu.Lock()
					p.food.Add(0.3)
					p.mu.Unlock()
				case 'r':
					cell.Log("ruuuuunnn!")
					// run := Run{}
					// p.StartActivity(&run)
					p.mu.Lock()
					p.stamina.Add(-0.3)
					p.mu.Unlock()

				case 'z':
					cell.Log("owch!")
					p.mu.Lock()
					damage := math.Max(0, 0.001+rand.NormFloat64()/30)
					if damage > 0 {
						cell.Log(fmt.Sprintf("You are hurt, losing %.3f health\n", damage))
						p.health.Add(-damage)
						p.notify(fmt.Sprintf("You are hurt, losing %.3f health\n", damage), "default")
						p.notify(fmt.Sprint(p), "low")
					}
					if rand.Intn(10) >= 8 {
						p.blood.Add(-0.2)
					}
					p.mu.Unlock()
				default:
					cell.Log(fmt.Sprintf("You pressed the %v key.\n", ev.Rune()))
				}
			default:
				cell.Log("unknown event")
			}
		}
	}
}

func showValue(c *Cell, title string, attr *Scalar, x int, y int) {
	var style tcell.Style
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
		style = dangerStyle
	} else if value < 0.5 {
		style = warningStyle
	} else if value < 0.8 {
		style = okStyle
	} else {
		style = goodStyle
	}

	c.DrawTextWithStyle(x, y, x+9, y, title, style)
	c.DrawTextWithStyle(x, y+1, x+9, y+1, fmt.Sprintf("%.03f", attr.current), style)
}

func showPlayer(p *Person, c *Cell) {
	// return fmt.Sprintf("%v: F%0.3f; W%0.3f; B%0.3f; H:%0.3f; S:%0.3f; T:%0.3f; A:%0.3f; I:%0.3f", p.name, p.food.current, p.water.current, p.blood.current, p.health.current, p.stamina.current, p.temperature.current, p.ambient_temperature.current, p.insulation.current)
	showValue(c, "Food", &p.food, 10, 2)
	showValue(c, "Water", &p.water, 20, 2)
	showValue(c, "Blood", &p.blood, 30, 2)
	showValue(c, "Health", &p.health, 40, 2)
	showValue(c, "Stamina", &p.stamina, 50, 2)
	showValue(c, "Temp.", &p.temperature, 60, 2)
	showValue(c, "Amb.Temp.", &p.ambient_temperature, 70, 2)
	showValue(c, "Insul.", &p.insulation, 80, 2)
}

func main() {
	p1 := Person{name: "Alex"}
	cell := CreateCell()
	go p1.Initialise()
	go kb(&p1, cell)
	time.Sleep(100 * time.Millisecond)
	for {
		showPlayer(&p1, cell)
		time.Sleep(100 * time.Millisecond)
	}
}
