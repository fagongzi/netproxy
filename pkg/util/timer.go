package util

import (
	"time"

	"github.com/fagongzi/goetty"
)

var (
	tw = goetty.NewTimeoutWheel(goetty.WithTickInterval(time.Millisecond * 100))
)

// GetTimeWheel get time wheel
func GetTimeWheel() *goetty.TimeoutWheel {
	return tw
}
