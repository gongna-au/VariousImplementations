package ChannelApplicationScenario

import (
	"fmt"
	"testing"
	"time"
)

func TestLimit(t *testing.T) {

	limit := NewChannelLimit()
	for i := 0; i < 30; i++ {
		temp := i

		if limit.Allow() {
			go func() {
				fmt.Println(temp)
			}()
			err := limit.Release()
			if err != true {
				fmt.Println("error")
			}
		}
	}
	time.Sleep(10 * time.Second)

}
