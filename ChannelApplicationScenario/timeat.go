package ChannelApplicationScenario

import (
	"fmt"
	"time"
)

func TaskTime(t time.Duration) {
	i := 0
	for {
		select {
		case <-time.After(t * time.Second):
			fmt.Println("task", i)
			i++
		case <-time.After(10 * time.Second):
			return
		}
	}
}
