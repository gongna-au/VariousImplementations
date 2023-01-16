package ChannelApplicationScenario

import (
	"context"
	"fmt"
	"time"
)

func TestContext() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	go Watch(ctx, "worker1")
	go Watch(ctx, "worker2")
	go Watch(ctx, "worker3")
	index := 0
	for {
		index++
		time.Sleep(3 * time.Second)
		if index >= 1 {
			break
		}
	}
	cancel()
	fmt.Println("all work has done")
}

func Watch(ctx context.Context, name string) {

	//time.Sleep(4 * time.Second)

	select {
	case <-ctx.Done():
		return
	default:
		fmt.Println(name, " is watching")
	}

}
