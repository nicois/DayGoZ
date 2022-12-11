///usr/bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"fmt"
	"log"
	"time"
)

type Consumable interface {
	GetRemainingVolume() float64
	Consume(person *Person, minutes float64)
	String() string
	Merge(consumable Consumable) bool
}

type blood struct {
	volume float64
}

func (f *blood) GetRemainingVolume() float64 {
	return f.volume / 100
}

func (f *blood) String() string {
	return fmt.Sprintf("🩸: %.1f", f.volume)
}

func (f *blood) Consume(person *Person, minutes float64) {
	volumePerMinute := 1000.0
	bloodPerVolume := 0.002
	volume := min(f.volume, volumePerMinute*minutes)
	if volume <= 0 {
		return
	}
	f.volume -= volume
	person.mu.Lock()
	defer person.mu.Unlock()
	person.blood.Add(bloodPerVolume * volume)
}

func (f *blood) Merge(c Consumable) bool {
	switch v := c.(type) {
	case *blood:
		f.volume += v.volume
		return true
	}
	return false
}

type food struct {
	name            string
	icon            rune
	volume          float64
	volumePerMinute float64
	foodPerVolume   float64
	waterPerVolume  float64
	healthPerVolume float64
	height          int
	width           int
}

func (f *food) Use(p *Person) bool {
	p.mouth <- f
	return true
}

func (f *food) Combine(other Item) Item {
	return nil
}

func (f *food) Icon() rune {
	return f.icon
}

func (f *food) GetSize() (int, int) {
	if f.height > 0 && f.width > 0 {
		return f.width, f.height
	}
	if f.volume > 200 {
		return 2, 2
	}
	if f.volume > 100 {
		return 2, 1
	}
	return 1, 1
}

func (f *food) Merge(c Consumable) bool {
	switch v := c.(type) {
	case *food:
		if v.name == f.name {
			f.volume += v.volume
			return true
		}
		return false
	}
	return false
}

func (f *food) String() string {
	if f.icon == 0 {
		return fmt.Sprintf("%v: %.1f", f.name, f.volume)
	} else {
		return fmt.Sprintf("%v: %.1f", string(f.icon), f.volume)
	}
}

func Digester(p *Person, mouth chan Consumable) {
	digestPerMinute := 60.0
	capacity := 500.0
	warned := false
	var stomach []Consumable
	log.Println("Starting to digest")
	defer log.Println("Finished digesting.")
loopy:
	for {
		if len(stomach) == 0 {
			warned = false
			stomach = append(stomach, <-mouth)
		}
		digestionTimer := time.NewTimer(time.Minute / time.Duration(digestPerMinute))
		select {
		case consumable := <-mouth:
			for _, digesting := range stomach {
				if digesting.Merge(consumable) {
					continue loopy
				}
			}
			stomach = append(stomach, consumable)
		case <-digestionTimer.C:
			totalVolume := 0.0
			for _, consumable := range stomach {
				totalVolume += consumable.GetRemainingVolume()
			}
			remainingCapacity := capacity - totalVolume
			p.stomach = fmt.Sprintf("Stomach remaining capacity: %.0f", remainingCapacity)
			for _, digesting := range stomach {
				p.stomach += fmt.Sprintf("; %v", digesting)
			}
			if totalVolume <= 0.0000001 {
				p.Log("finished digesting for now")
				stomach = nil
				break
			}
			if totalVolume > 500 {
				p.Log("oh dear. You should eat more slowly! I hope there is no vomit on your shoes!")
				stomach = nil
				p.health.Add(-0.2)
				p.water.Add(-0.2)
				break
			}
			if totalVolume <= 400 && warned {
				p.Log("Your stomach is less bloated.")
				warned = false
			}

			if totalVolume > 400 && !warned {
				p.Log("Your stomach is getting pretty full. You probably shouldn't eat for a while.")
				p.Log(p.stomach)
				warned = true
			}
			for _, consumable := range stomach {
				consumable.Consume(p, consumable.GetRemainingVolume()/(totalVolume*digestPerMinute))
			}
		}
	}
}

func (f *food) GetRemainingVolume() float64 {
	return f.volume
}

func (f *food) Consume(person *Person, minutes float64) {
	volume := min(f.volume, f.volumePerMinute*minutes)
	if volume <= 0 {
		return
	}
	f.volume -= volume
	person.mu.Lock()
	defer person.mu.Unlock()
	person.food.Add(f.foodPerVolume * volume)
	person.water.Add(f.waterPerVolume * volume)
	person.health.Add(f.healthPerVolume * volume)
}

func Blood(volume float64) *blood {
	return &blood{volume: volume}
}

func Water(volume float64) *food {
	return &food{name: "water", volume: volume, volumePerMinute: 100, foodPerVolume: 0, waterPerVolume: 0.01, icon: '🚰'}
}

func Chocolate() *food {
	return &food{name: "chocolate", volume: 100, volumePerMinute: 100, foodPerVolume: 0.001, waterPerVolume: 0, healthPerVolume: -0.0005, icon: '🍫'}
}

func MilkPowder() *food {
	return &food{name: "a packet of milk powder", height: 2, width: 1, volume: 100, volumePerMinute: 20, foodPerVolume: 0.002, waterPerVolume: -0.001, icon: '🥛'}
}

func Peaches() *food {
	return &food{name: "a can of juicy peaches", height: 2, width: 1, volume: 300, volumePerMinute: 30, foodPerVolume: 0.001, waterPerVolume: 0.001, icon: '🍑'}
}

func Apple() *food {
	return &food{name: "apple", volume: 80, volumePerMinute: 10, foodPerVolume: 0.001, waterPerVolume: 0.001, icon: '🍏'}
}

func Banana() *food {
	return &food{name: "banana", volume: 100, volumePerMinute: 1, foodPerVolume: 0.003, waterPerVolume: 0.0001, icon: '🍌'}
}

func Jaffle() *food {
	return &food{name: "jaffle", volume: 120, volumePerMinute: 4, foodPerVolume: 0.003, waterPerVolume: 0.00001, icon: '🫓'}
}

func Ryvita() *food {
	return &food{name: "ryvita", volume: 50, volumePerMinute: 10, foodPerVolume: 0.01, waterPerVolume: -0.001, icon: '🍘'}
}

func MagicPotion() *food {
	return &food{name: "magic potion", volume: 100, volumePerMinute: 100, healthPerVolume: 0.001, waterPerVolume: 0.001, icon: '🏺'}
}
