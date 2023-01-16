package CounterLimit

import (
	"fmt"
	"sync"
	"time"
)

type counter struct {
	// 上一次访问的时间
	lastAccessTime time.Time
	// 计数器清零的时间间隔
	duration time.Duration
	// 请求计数器
	sum int64
	// 请求的最大阀值
	Threshold int64
	sync.Mutex
}

func Newcounter(duration time.Duration, threshold int64) *counter {

	return &counter{
		lastAccessTime: time.Now(),
		sum:            int64(0),
		Threshold:      threshold,
		duration:       duration,
	}
}

func (c *counter) Allow(num int64) bool {

	duration := time.Now().Sub(c.lastAccessTime)
	fmt.Println(duration)
	c.Lock()
	defer c.Unlock()
	if duration >= c.duration {
		// 超时重置
		fmt.Println("超时重置")
		c.sum = 0
		c.lastAccessTime = time.Now()
		if num > c.Threshold {
			return false
		} else {
			c.sum = num
			return true
		}
	} else {
		fmt.Println("未超时")

		if c.sum+num >= c.Threshold {
			return false
		} else {
			c.sum = c.sum + num
		}
	}

	return true
}
