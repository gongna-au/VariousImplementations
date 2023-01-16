package ChannelApplicationScenario

import (
	"fmt"
	"time"
)

func ChanCacel(ch chan bool) {
	for {
		select {
		case <-ch:
			return
		default:
			fmt.Println("work ok")
		}
	}
}

func ChanReturnTest() {
	ch := make(chan bool)
	go ChanCacel(ch)
	time.Sleep(1 * time.Second)
	close(ch)

}
