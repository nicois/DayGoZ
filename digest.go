///usr/bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
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
	var stomach []Consumable
	for {
		if len(stomach) == 0 {
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
			if totalVolume <= 0.0000001 {
				p.Log("finished digesting for now")
				stomach = nil
				break
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
	return &Food{name: "water", volume: volume, volumePerMinute: 10, foodPerVolume: 0, waterPerVolume: 0.01}
}

func Chocolate() Consumable {
	return &Food{name: "chocolate", volume: 100, volumePerMinute: 100, foodPerVolume: 0.001, waterPerVolume: 0}
}

func Peaches() Consumable {
	return &Food{name: "a can of juicy peaches", volume: 300, volumePerMinute: 10, foodPerVolume: 0.005, waterPerVolume: 0.002}
}

func Banana() Consumable {
	return &Food{name: "banana", volume: 100, volumePerMinute: 1, foodPerVolume: 0.005, waterPerVolume: 0.0001}
}
