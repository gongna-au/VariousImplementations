package ChannelApplicationScenario

import (
	"fmt"
	"sync"
)

func Work() {
	inCh := generator(20)
	outCh := make(chan int, 5)
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go do(inCh, outCh, &wg)
	}
	go func() {
		wg.Wait()
		close(outCh)
	}()
	for v := range outCh {
		fmt.Println(v)
	}
}

func generator(n int) <-chan int {
	outCh := make(chan int)
	go func() {
		for i := 0; i < n; i++ {
			outCh <- i
		}
		close(outCh)
	}()
	return outCh
}

func do(inCh <-chan int, outCh chan<- int, wg *sync.WaitGroup) {
	for v := range inCh {
		outCh <- v * v
	}
	fmt.Println("done")
	wg.Done()
}
