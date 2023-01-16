package CounterLimit

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestCounter(t *testing.T) {
	counter := NewCounter()
	// 在time2.Sub(time1)这个时间内只能请求三次
	counter.Set(time.Second, 3)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			if counter.Allow() {
				log.Println("响应请求：", i)
			}
			wg.Done()
			time.Sleep(200 * time.Millisecond)
		}(i)
	}

	wg.Wait()

}
