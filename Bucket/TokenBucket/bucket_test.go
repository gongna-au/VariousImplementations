package TokenBucket

import (
	"math"
	"sync"
	"time"
)

type tokenBucket struct {
	// 放置的速率
	rate int64
	// 桶最大的限制
	threshold int64
	// 桶当前的令牌的总数量
	sum int64
	// 上一次防置Token 令牌的时间
	lastTokenSet int64
	sync.Mutex
}

func NewTokenBucket() *tokenBucket {
	return &tokenBucket{
		sum: 0,
	}
}

func (t *tokenBucket) Allow() bool {
	t.Lock()
	defer t.Unlock()
	now := time.Now().Unix()
	//先添加令牌
	t.sum = int64(math.Min(float64((now-t.lastTokenSet)*t.rate+t.sum), float64(t.threshold)))
	if t.sum-1 > 0 {
		// 消耗令牌桶
		t.sum = t.sum - 1
		return true
	}
	return false
}

func (t *tokenBucket) Set(rate int64, threshold int64) {
	t.rate = rate
	t.threshold = threshold
	t.lastTokenSet = time.Now().Unix()
	t.sum = 0
}
