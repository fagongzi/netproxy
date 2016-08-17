package util

import (
	"github.com/fagongzi/goetty"
	"time"
)

var (
	TimeWheel = goetty.NewHashedTimeWheel(time.Millisecond, 60, 3)
)

func Init() {
	TimeWheel.Start()
}

func GetTimeWheel() *goetty.HashedTimeWheel {
	return TimeWheel
}
