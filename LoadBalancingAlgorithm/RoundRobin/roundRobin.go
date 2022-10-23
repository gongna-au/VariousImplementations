package RoundRobin

import (
	"fmt"
	Balancer "github.com/VariousImplementations/LoadBalancingAlgorithm"
	mrand "math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

//  New 使用 Round-Robin 创建一个新的负载均衡器实例

type RoundRobin struct {
	peers []lbapi.Peer
	count int64
	rw    sync.RWMutex
}
