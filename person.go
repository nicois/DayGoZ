///usr/bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"fmt"
	"math"
	"nicois/bmon/ntfy"
	"sync"
	"time"
)

type NumericAttribute interface {
	Initialise()
	Listen(chan float64)
	Add(float64) float64
}

type Interval struct {
	min float64
	max float64
}

type Scalar struct {
	min       float64
	max       float64
	current   float64
	listeners []chan float64
}

func (a *Scalar) limits() Interval {
	return Interval{min: a.min - a.current, max: a.max - a.current}
}

//func (a *Scalar) String() string {
//	return fmt.Sprintf("[%v-%v] %v", a.min, a.max, a.current)
//}

func (a *Scalar) Add(v float64) float64 {
	original := a.current
	a.current += v
	if a.current < a.min {
		a.current = a.min
	}
	if a.current > a.max {
		a.current = a.max
	}
	for _, listener := range a.listeners {
		listener <- a.current
	}
	return a.current - original
}

// the state of A (e.g. food) can affect B (e.g. health)
func (a *Scalar) Link(b *Scalar) {
	last_update := time.Now()
	for {
		time.Sleep(10 * time.Millisecond)
		if a.current > a.min {
			now := time.Now()
			duration := now.Sub(last_update)
			added := b.Add(duration.Seconds())
			a.Add(-added)
			last_update = now
		}
	}
}

type Feedback interface {
	Link()
}

type Person struct {
	mu                  *sync.Mutex
	name                string
	health              Scalar
	food                Scalar
	water               Scalar
	blood               Scalar
	stamina             Scalar
	temperature         Scalar
	ambient_temperature Scalar
	insulation          Scalar
	sender              Sender

	activity          Activity
	activityMutex     *sync.Mutex
	activityTicker    *chan float64
	activityCanceller *chan bool
}

type Activity interface {
	GetChannels() (canceller *chan bool, ticker *chan float64)
	Begin() error
	Cancel()
}

func (p *Person) StartActivity(activity Activity) error {
	if !p.activityMutex.TryLock() {
		return fmt.Errorf("Could not get a lock; an activity is probably in progress.")
	}
	activity.Cancel()
	p.activityCanceller, p.activityTicker = activity.GetChannels()
	p.activity = activity
	activity.Begin()
	return nil
}

func (p *Person) Initialise() {
	p.mu = new(sync.Mutex)
	p.activityMutex = new(sync.Mutex)
	p.stamina = Scalar{max: 1, current: 0.8}
	p.temperature = Scalar{max: 1, current: 0.5}
	p.ambient_temperature = Scalar{min: -0.5, max: 1.5, current: 0.4}
	p.insulation = Scalar{max: 1, current: 0.1}
	p.blood = Scalar{max: 1, current: 1}
	p.health = Scalar{max: 1, current: 1}
	p.water = Scalar{max: 1, current: 0.9}
	p.food = Scalar{max: 1, current: 0.9}
	// p.sender = ntfy.Create("foobarbaz")

	last_time := time.Now()
	summary := fmt.Sprint(p)
	p.notify(summary, "min")
	// fixme:  define channel to allow activities
	// to force a resync of the time-based effects,
	// taking place as the action passes warmup
	for tick := range time.Tick(100 * time.Millisecond) {
		if p.health.current == 0 {
			p.notify(summary, "max")
			return
		}
		p.calculate_time_based_effects(tick.Sub(last_time).Minutes())
		last_time = tick
	}
}

type Sender interface {
	Send(ntfy.Message) error
}

func (p *Person) notify(text string, priority string) {
	if sender := p.sender; sender != nil {
		headers := map[string]string{"Priority": priority}
		err := sender.Send(ntfy.Message{Text: text, Headers: headers})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (p *Person) calculate_time_based_effects(minutes float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.activity != nil {
	}

	// blood: consume food and water, goes up faster when healthier
	p.adjust(map[*Scalar]float64{&p.food: -0.01, &p.water: -0.01, &p.blood: p.health.current * 0.01}, minutes)

	// stamina:
	rate := math.Pow(3*(1-math.Abs(p.stamina.current-0.5)), 4) / 10
	p.adjust(map[*Scalar]float64{&p.food: -0.01 * rate, &p.water: -0.01 * rate, &p.stamina: rate, &p.temperature: 0.01 * rate}, minutes)

	// body temp:
	temperature_difference := p.temperature.current - 0.5
	if temperature_difference > 0.01 {
		rate := math.Pow(10*temperature_difference, 2)
		p.adjust(map[*Scalar]float64{&p.water: -0.01 * rate, &p.temperature: -0.01 * rate * (1 - p.insulation.current)}, minutes)
	}

	if temperature_difference < -0.01 {
		rate := math.Pow(10*temperature_difference, 2)
		p.adjust(map[*Scalar]float64{&p.food: -0.01 * rate, &p.temperature: 0.01 * rate}, minutes)
	}

	// ambient temp effects:
	rate = (p.ambient_temperature.current - p.temperature.current) * (1 - p.insulation.current)
	p.adjust(map[*Scalar]float64{&p.temperature: 0.1 * rate}, minutes)

	// health (temperature effects)
	temperature_happiness := 0.1 - math.Abs(p.temperature.current-0.5)
	p.adjust(map[*Scalar]float64{&p.health: temperature_happiness / 10, &p.food: -0.01, &p.water: -0.01}, minutes)

	// health (starvation effects)
	if deficit := 0.2 - min(p.food.current, p.water.current); deficit > 0 {
		p.adjust(map[*Scalar]float64{&p.health: -deficit / 2}, minutes)
	}

	// metabolism
	p.adjust(map[*Scalar]float64{&p.food: -0.0005, &p.water: -0.0005, &p.temperature: 0.0005}, minutes)
}

// Adjust the attributes of the person respecting the configured ratios. Scale by rate (per minute)
// If necessary scale down the adjustment to not exceed any absolute constraints given in  max_adjust.
func (p *Person) adjust(amounts map[*Scalar]float64, rate float64) {
	rescale := 1.0
	effect := make(map[*Scalar]float64)
	for scalar, ratio := range amounts {
		amount_to_add := ratio * rate
		effect[scalar] = amount_to_add
		limit := scalar.limits()
		if amount_to_add > 0 {
			if limit.max == 0 {
				return
			}
			if exceeded_ratio := (amount_to_add / limit.max); exceeded_ratio > rescale {
				rescale = exceeded_ratio
			}
		}
		if amount_to_add < 0 {
			if limit.min == 0 {
				return
			}
			if exceeded_ratio := (amount_to_add / limit.min); exceeded_ratio > rescale {
				rescale = exceeded_ratio
			}
		}
	}
	if rescale > 1.0 {
		for scalar := range effect {
			effect[scalar] *= 1 / rescale
		}
	}
	for scalar, amount := range effect {
		scalar.Add(amount)
	}
}

func (p Person) String() string {
	return fmt.Sprintf("%v: F%0.3f; W%0.3f; B%0.3f; H:%0.3f; S:%0.3f; T:%0.3f; A:%0.3f; I:%0.3f", p.name, p.food.current, p.water.current, p.blood.current, p.health.current, p.stamina.current, p.temperature.current, p.ambient_temperature.current, p.insulation.current)
}
