///usr/bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"fmt"
	"time"
)

type Consumable interface {
	GetRemainingVolume() float64
	Consume(person *Person, minutes float64)
}

type Food struct {
	name            string
	volume          float64
	volumePerMinute float64
	foodPerVolume   float64
	waterPerVolume  float64
	healthPerVolume float64
}

func Digester(p *Person, mouth chan Consumable) {
	chewsPerMinute := 60.0
	capacity := 500.0
	warned := false
	var stomach []Consumable
	for {
		if len(stomach) == 0 {
			warned = false
			stomach = append(stomach, <-mouth)
		}
		digestionTimer := time.NewTimer(time.Minute / time.Duration(chewsPerMinute))
		select {
		case consumable := <-mouth:
			stomach = append(stomach, consumable)
		case <-digestionTimer.C:
			totalVolume := 0.0
			for _, consumable := range stomach {
				totalVolume += consumable.GetRemainingVolume()
			}
			remainingCapacity := capacity - totalVolume
			p.stomach = fmt.Sprintf("Stomach remaining capacity: %.0f", remainingCapacity)
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
				warned = true
			}
			for _, consumable := range stomach {
				consumable.Consume(p, consumable.GetRemainingVolume()/(totalVolume*chewsPerMinute))
			}
		}
	}
}

func (f *Food) GetRemainingVolume() float64 {
	return f.volume
}

func (f *Food) Consume(person *Person, minutes float64) {
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

func Water(volume float64) Consumable {
	return &Food{name: "water", volume: volume, volumePerMinute: 100, foodPerVolume: 0, waterPerVolume: 0.01}
}

func Chocolate() Consumable {
	return &Food{name: "chocolate", volume: 100, volumePerMinute: 100, foodPerVolume: 0.001, waterPerVolume: 0}
}

func Peaches() Consumable {
	return &Food{name: "a can of juicy peaches", volume: 300, volumePerMinute: 10, foodPerVolume: 0.005, waterPerVolume: 0.002}
}

func Apple() Consumable {
	return &Food{name: "apple", volume: 80, volumePerMinute: 10, foodPerVolume: 0.001, waterPerVolume: 0.001}
}

func Banana() Consumable {
	return &Food{name: "banana", volume: 100, volumePerMinute: 1, foodPerVolume: 0.003, waterPerVolume: 0.0001}
}

func Ryvita() Consumable {
	return &Food{name: "ryvita", volume: 100, volumePerMinute: 25, foodPerVolume: 0.001, waterPerVolume: -0.001}
}
func MagicPotion() Consumable {
	return &Food{name: "magic potion", volume: 100, volumePerMinute: 100, healthPerVolume: 0.001, waterPerVolume: 0.001}
}
