package RoundRobin

import (
	"fmt"
	Balancer "github.com/VariousImplementations/LoadBalancingAlgorithm"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type RoundRobin struct {
	peers []Balancer.Peer
	count int64
	rw    sync.RWMutex
}

//  New 使用 Round-Robin 创建一个新的负载均衡器实例
func New(opts ...Balancer.Opt) Balancer.Balancer {
	return &RoundRobin{}
}

// RoundRobin  需要实现 Balancer接口下面的方法Balancer.Next()  Balancer.Count()  Balancer.Add() Balancer.Remove()  Balancer.Clear()
func (s *RoundRobin) Next(factor Balancer.Factor) (next Balancer.Peer, c Balancer.Constrainable) {
	next = s.miniNext()
	if fc, ok := factor.(Balancer.FactorComparable); ok {
		next, c, _ = fc.ConstrainedBy(next)
	} else if nested, ok := next.(Balancer.BalancerLite); ok {
		next, c = nested.Next(factor)
	}

	return
}

// s.count 会一直增量上去，并不会取模
// s.count 增量加1就是轮询的核心
// 这样做的用意在于如果 peers 数组发生了少量的增减变化时，最终发生选择时可能会更模棱两可。
// 但是！！！注意对于 Golang 来说，s.count 来到 int64.MaxValue 时继续加一会自动回绕到 0。
// 这一特性和多数主流编译型语言相同，都是 CPU 所提供的基本特性
// 核心的算法 s.count 对后端节点的列表长度取余
func (s *RoundRobin) miniNext() (next Balancer.Peer) {
	ni := atomic.AddInt64(&s.count, 1)
	ni--
	// 加入读锁
	s.rw.RLock()
	defer s.rw.RUnlock()
	if len(s.peers) > 0 {
		ni %= int64(len(s.peers))
		next = s.peers[ni]
	}
	fmt.Printf("s.peers[%d] is be returned\n", ni)
	return
}
func (s *RoundRobin) Count() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.peers)
}

func (s *RoundRobin) Add(peers ...Balancer.Peer) {
	for _, p := range peers {
		s.AddOne(p)
	}
}

func (s *RoundRobin) AddOne(peer Balancer.Peer) {
	if s.find(peer) {
		return
	}
	s.rw.Lock()
	defer s.rw.Unlock()
	s.peers = append(s.peers, peer)
}

func (s *RoundRobin) find(peer Balancer.Peer) (found bool) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	for _, p := range s.peers {
		if Balancer.DeepEqual(p, peer) {
			return true
		}
	}
	return
}

func (s *RoundRobin) Remove(peer Balancer.Peer) {
	// 加写锁
	s.rw.Lock()
	defer s.rw.Unlock()
	for i, p := range s.peers {
		if Balancer.DeepEqual(p, peer) {
			s.peers = append(s.peers[0:i], s.peers[i+1:]...)
			return
		}
	}
}

func (s *RoundRobin) Clear() {
	// 加写锁
	s.rw.Lock()
	defer s.rw.Unlock()
	s.peers = nil
}

func Client() {
	// wg让主进程进行等待我所有的goroutinue 完成
	wg := sync.WaitGroup{}
	// 假设我们有20个不同的客户端（goroutinue）去调用我们的服务
	wg.Add(20)
	lb := &RoundRobin{
		peers: []Balancer.Peer{
			Balancer.ExP("172.16.0.10:3500"), Balancer.ExP("172.16.0.11:3500"), Balancer.ExP("172.16.0.12:3500"),
		},
		count: 0,
	}
	for i := 0; i < 10; i++ {
		go func(t int) {
			lb.Next(Balancer.DummyFactor)
			wg.Done()
			time.Sleep(2 * time.Second)
			// 这句代码第一次运行后，读解锁。
			// 循环到第二个时，读锁定后，这个goroutine就没有阻塞，同时读成功。
		}(i)

		go func(t int) {
			str := "172.16.0." + strconv.Itoa(t) + ":3500"
			lb.Add(Balancer.ExP(str))
			fmt.Println(str + " is be added. ")
			wg.Done()
			// 这句代码让写锁的效果显示出来，写锁定下是需要解锁后才能写的。
			time.Sleep(2 * time.Second)
		}(i)
	}

	time.Sleep(5 * time.Second)
	wg.Wait()
}
