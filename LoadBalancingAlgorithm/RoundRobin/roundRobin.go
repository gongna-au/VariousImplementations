package RoundRobin

import (
	"sync"
	"sync/atomic"

	Balancer "github.com/VariousImplementations/LoadBalancingAlgorithm"
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

// 核心的算法 ni 对后端节点的列表长度取余
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
