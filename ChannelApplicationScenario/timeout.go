package ChannelApplicationScenario

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func DoWithTimeOut(timeout time.Duration) (int, error) {
	select {
	case <-Do():
		for v := range Do() {
			fmt.Println(v)
		}
		return 999, nil
	case <-time.After(timeout):
		return -1, errors.New("time out")
	}
}

func Do() <-chan int {
	var wg sync.WaitGroup
	outCh := make(chan int, 30)
	wg.Add(10)
	for i := 0; i < 10; i++ {
		temp := i
		go func() {
			outCh <- temp
			wg.Done()
		}()
	}
	wg.Wait()
	close(outCh)
	return outCh
}
