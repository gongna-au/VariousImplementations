package ChannelApplicationScenario

import (
	"fmt"
	"time"
)

// 往桶里填充数据
func TokenBucketFill(ch chan int) {
	for {
		select {
		case <-time.After(1 * time.Second):
			select {
			case ch <- 1:
			default:
			}

		}
	}
}

// 往桶里填充数据
func TokenBucketGet(ch chan int, workFunc t) {
	for {
		select {
		case <-ch:
			workFunc()
		}
	}
}

func TestTokenBucket() {
	ch := make(chan int, 10)
	go TokenBucketFill(ch)
	go TokenBucketGet(ch, w)
	time.Sleep(3 * time.Second)

}

type t func()

func w() {
	fmt.Println("ok")
}
