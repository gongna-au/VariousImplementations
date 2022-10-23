package Random

//正式的 random LB 的代码要比上面的核心部分还复杂一点点。原因在于我们还需要达成另外两个设计目标：
import (
	"fmt"
	"github.com/VariousImplementations/LoadBalancingAlgorithm/Random/Balancer"
	mrand "math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var seedRand = mrand.New(mrand.NewSource(time.Now().Unix()))
var seedMutex sync.Mutex

func InRange(min, max int64) int64 {
	seedMutex.Lock()
	defer seedMutex.Unlock()
	return seedRand.Int63n(max-min) + min
}

// New 使用 Round-Robin 创建一个新的负载均衡器实例
func New(opts ...Balancer.Opt) Balancer.Balancer {
	return (&randomS{}).Init(opts...)
}

type randomS struct {
	peers []Balancer.Peer
	count int64
	rw    sync.RWMutex
}

func (s *randomS) Init(opts ...Balancer.Opt) *randomS {
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// 实现了Balancer.NexT()方法
func (s *randomS) Next(factor Balancer.Factor) (next Balancer.Peer, c Balancer.Constrainable) {
	next = s.miniNext()

	if fc, ok := factor.(Balancer.FactorComparable); ok {
		next, c, ok = fc.ConstrainedBy(next)
	} else if nested, ok := next.(Balancer.BalancerLite); ok {
		next, c = nested.Next(factor)
	}

	return
}

// 实现了Balancer.Count()方法
func (s *randomS) Count() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.peers)
}

// 实现了Balancer.Add()方法
func (s *randomS) Add(peers ...Balancer.Peer) {
	for _, p := range peers {
		// 判断要添加的元素是否存在，并且在添加元素的时候为s.peers 加锁
		s.AddOne(p)
	}
}

// 实现了Balancer.Remove()方法
// 如果 s.peers 中间有和传入的peer相等的函数就那么就删除这个元素
// 在删除这个元素的时候，
func (s *randomS) Remove(peer Balancer.Peer) {
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

// 实现了Balancer.Clear()方法
func (s *randomS) Clear() {
	// 加写锁
	// 对于Set() ,Delete(),Update()这类操作就一般都是加写锁
	// 对于Get() 这类操作我们往往是加读锁，阻塞对同一变量的更改操作，但是读操作将不会受到影响
	s.rw.Lock()
	defer s.rw.Unlock()
	s.peers = nil
}

// 我们希望s在返回后端peers 节点的时候，在同一个时刻只能被一个线程拿到。
// 所以需要对 s.peers进行加锁
func (s *randomS) miniNext() (next Balancer.Peer) {
	// 读锁定 写将被阻塞，读不会被锁定
	s.rw.RLock()
	defer s.rw.RUnlock()
	l := int64(len(s.peers))
	ni := atomic.AddInt64(&s.count, InRange(0, l)) % l
	next = s.peers[ni]
	fmt.Printf("s.peers[%d] is be returned\n", ni)
	return
}

func (s *randomS) AddOne(peer Balancer.Peer) {
	if s.find(peer) {
		return
	}
	// 加了写锁
	// 在更改s.peers的时候，其他的线程将不可以调用s.miniNext()读出和获得peer，其他的线程也不可以调用s.AddOne()对s.peers 进行添加操作
	s.rw.Lock()
	defer s.rw.Unlock()
	s.peers = append(s.peers, peer)
	fmt.Printf(peer.String() + "is be appended!\n")
}

func (s *randomS) find(peer Balancer.Peer) (found bool) {
	// 加读锁
	s.rw.RLock()
	defer s.rw.RUnlock()
	for _, p := range s.peers {
		if Balancer.DeepEqual(p, peer) {
			return true
		}
	}
	fmt.Printf("peer in s.peers is be found!\n")
	return
}

func Client() {
	// wg让主进程进行等待我所有的goroutinue 完成
	wg := sync.WaitGroup{}
	// 假设我们有20个不同的客户端（goroutinue）去调用我们的服务
	wg.Add(20)
	lb := &randomS{
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
