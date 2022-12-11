///usr/bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"time"
)

func min(fs ...float64) float64 {
	min := fs[0]
	for _, i := range fs[1:] {
		if i < min {
			min = i
		}
	}
	return min
}

func minDuration(a, b time.Duration) time.Duration {
	if a <= b {
		return a
	}
	return b
}

func remove[Thing interface{}](s []Thing, i int) []Thing {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
