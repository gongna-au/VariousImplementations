package ChannelApplicationScenario

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()
	for {
		select {
		default:
			fmt.Println("Hello")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func ChannelCancel() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(ctx, &wg)
	}

	cancel()
	wg.Wait()
}
