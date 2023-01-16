package LeakyBucket

import (
	"math"
	"sync"
	"time"
)

type leakyBucket struct {
	// 最大可以容纳的水的容量
	cap float64
	// 出水速率
	rate float64
	//当前桶内部水的含量
	water float64
	//上次露水的时间
	lastLeakMs int64
	// 控制并发
	lock sync.Mutex
}

func (l *leakyBucket) Allow() bool {
	l.lock.Lock()
	l.lock.Unlock()
	now := time.Now().UnixNano() / 1e6
	// 先执行漏水
	eclipse := float64(now-l.lastLeakMs) * float64(l.rate) / 1000
	// 更新当前水的含量
	l.water = l.water - eclipse
	l.water = math.Max(float64(0), float64(l.water))
	l.lastLeakMs = now
	if (l.water + 1) < l.cap {
		l.water++
		return true
	} else {
		return false
	}
}

func (l *leakyBucket) Set(rate float64, cap float64) {
	l.rate = rate
	l.cap = cap
	l.water = 0
	l.lastLeakMs = time.Now().UnixNano() / 1e6
}
