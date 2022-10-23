package Hashing

import (
	"fmt"
	Balancer "github.com/VariousImplementations/LoadBalancingAlgorithm"
	"hash/crc32"
	"strings"
	"sync"
	"testing"
)

func TestHash(t *testing.T) {

	h := int(uint(0+3) >> 1)
	fmt.Print(h)

}

type ConcretePeer string

func (s ConcretePeer) String() string {
	return string(s)
}

var factors = []Balancer.FactorString{
	"https://abc.local/user/profile",
	"https://abc.local/admin/",
	"https://abc.local/shop/item/1",
	"https://abc.local/post/35719",
}

func TestHash1(t *testing.T) {
	lb := New()
	lb.Add(
		ConcretePeer("172.16.0.7:3500"),
		ConcretePeer("172.16.0.8:3500"),
		ConcretePeer("172.16.0.9:3500"),
	)
	// 记录某个节点被调用的次数
	sum := make(map[Balancer.Peer]int)
	// 记录某个具体的节点被哪些ip地址访问过
	hits := make(map[Balancer.Peer]map[Balancer.Factor]bool)
	// 模拟不同时间三个ip 地址对服务端发起多次的请求
	for i := 0; i < 300; i++ {
		// ip 地址依次对服务端发起多次的请求
		factor := factors[i%len(factors)]
		// 把 ip 地址传进去得到具体的节点
		peer, _ := lb.Next(factor)

		sum[peer]++

		if ps, ok := hits[peer]; ok {
			// 判断该ip 地址是否之前访问过该节点
			if _, ok := ps[factor]; !ok {
				// 如果没有访问过则标志为访问过
				ps[factor] = true
			}
		} else {
			// 如过该节点对应的 (访问过该节点的map不存在)证明该节点一次都没有被访问过
			// 那么创建map来 存储该ip地址已经被访问过
			hits[peer] = make(map[Balancer.Factor]bool)
			hits[peer][factor] = true
		}
	}

	// results
	total := 0
	for _, v := range sum {
		total += v
	}

	for p, v := range sum {
		var keys []string
		// p为节点
		for fs := range hits[p] {
			// 打印出每个节点被哪些ip地址访问过
			if kk, ok := fs.(interface{ String() string }); ok {
				keys = append(keys, kk.String())
			} else {
				keys = append(keys, fs.Factor())
			}
		}
		fmt.Printf("%v\nis be invoked %v nums\nis be accessed by these [%v]\n", p, v, strings.Join(keys, ","))
	}

	lb.Clear()
}

func TestHash_M1(t *testing.T) {
	lb := New()
	lb.Add(
		ConcretePeer("172.16.0.7:3500"),
		ConcretePeer("172.16.0.8:3500"),
		ConcretePeer("172.16.0.9:3500"),
	)

	var wg sync.WaitGroup
	var rw sync.RWMutex
	sum := make(map[Balancer.Peer]int)

	const threads = 8
	wg.Add(threads)

	// 这个是最接近业务场景的因为是并发的请求
	for x := 0; x < threads; x++ {
		go func(xi int) {
			defer wg.Done()
			for i := 0; i < 600; i++ {
				p, c := lb.Next(factors[i%3])
				adder(p, c, sum, &rw)
			}
		}(x)
	}
	wg.Wait()
	// results
	for k, v := range sum {
		fmt.Printf("Peer:%v InvokeNum:%v\n", k, v)
	}
}

func TestHash2(t *testing.T) {
	lb := New(
		WithHashFunc(crc32.ChecksumIEEE),
		WithReplica(16),
	)
	lb.Add(
		ConcretePeer("172.16.0.7:3500"),
		ConcretePeer("172.16.0.8:3500"),
		ConcretePeer("172.16.0.9:3500"),
	)
	sum := make(map[Balancer.Peer]int)
	hits := make(map[Balancer.Peer]map[Balancer.Factor]bool)

	for i := 0; i < 300; i++ {
		factor := factors[i%len(factors)]
		peer, _ := lb.Next(factor)

		sum[peer]++
		if ps, ok := hits[peer]; ok {
			if _, ok := ps[factor]; !ok {
				ps[factor] = true
			}
		} else {
			hits[peer] = make(map[Balancer.Factor]bool)
			hits[peer][factor] = true
		}
	}
	lb.Clear()
}

func adder(key Balancer.Peer, c Balancer.Constrainable, sum map[Balancer.Peer]int, rw *sync.RWMutex) {
	rw.Lock()
	defer rw.Unlock()
	sum[key]++
}
