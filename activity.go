///usr/bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"fmt"
	"time"
)

type Run struct {
	person    *Person
	warmedUp  bool
	lastTick  time.Time
	canceller chan bool
	ticker    chan float64
}

func (activity *Run) GetChannels() (canceller *chan bool, ticker *chan float64) {
	return &activity.canceller, &activity.ticker
}

func (activity *Run) Cancel() {
	activity.canceller <- true
}

func (activity *Run) Begin() error {
	defer activity.person.activityMutex.Unlock()
	warmupInterval := time.Millisecond * 500
	cooldownInterval := time.Millisecond * 200
	startTime := time.Now()
	// warmup
	warmupTimer := time.NewTimer(warmupInterval)
	fmt.Println("preparing to run")
warmingUp:
	for {
		select {
		case <-activity.canceller:
			// did not complete warmup. Cooldown starts, but is not longer
			// than the amount of time we warmed up for
			cooldown := minDuration(cooldownInterval, time.Since(startTime))
			fmt.Println("cooling down for ", cooldown)
			time.Sleep(cooldown)
			return nil
		case <-warmupTimer.C:
			break warmingUp
		}
	}
	<-activity.ticker
	// flush any additional calctime events, in case of a race with a previous
	// activity
flush:
	for {
		select {
		case <-activity.ticker:
		default:
			break flush
		}
	}
	for {
		select {
		case duration := <-activity.ticker:
			fmt.Println("ran for ", duration)
		case <-activity.canceller:
			fmt.Println("cooling down for ", cooldownInterval)
			time.Sleep(cooldownInterval)
			return nil
		}
	}
}
