package CounterLimit

import (
	"sync"
	"time"
)

// 标准版实现
type Counter struct {
	// 计数周期的开始时间
	begin time.Time
	// 计数周期
	cycle time.Duration
	// 请求计数器
	sum int64
	// 请求周期允许的最大的请求的数量
	Threshold int64

	sync.Mutex
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Allow() bool {
	c.Lock()
	defer c.Unlock()
	if time.Now().Sub(c.begin) > c.cycle {
		c.Reset()
	}
	if c.sum+1 >= c.Threshold {
		return false
	} else {
		c.sum = c.sum + 1
	}

	return true
}

func (c *Counter) Reset() {
	c.begin = time.Now()
	c.sum = 0
}

// 多少时间内最多请求多少次
func (c *Counter) Set(cycle time.Duration, threshold int64) {
	c.begin = time.Now()
	c.sum = 0
	c.Threshold = threshold
	c.cycle = cycle

}
