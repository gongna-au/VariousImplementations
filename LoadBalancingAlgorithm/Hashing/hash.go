package Hashing

import (
	"fmt"
	Balancer "github.com/VariousImplementations/LoadBalancingAlgorithm"
	"hash/crc32"
	"sort"
	"sync"
)

// HashKetama 是一个带有 ketama 组合哈希算法的 impl
type HashKetama struct {
	// default is crc32.ChecksumIEEE
	hasher Hasher
	// 负载均衡领域中的一致性 Hash 算法加入了 Replica 因子，计算 Peer 的 hash 值时为 peer 的主机名增加一个索引号的后缀，索引号增量 replica 次
	// 也就是说一个 peer 的 拥有replica 个副本，n 台 peers 的规模扩展为 n x Replica 的规模，有助于进一步提高选取时的平滑度。
	replica int
	// 通过每调用一次Next()函数 ，往hashRing中添加一个计算出的哈希数值
	// 从哈希列表中得到一个哈希值，然后立即得到该哈希值对应的后端的节点
	hashRing []uint32
	// 每个节点都拥有一个属于自己的hash值
	// 每往hashRing 添加一个元素就，就往map中添加一个元素
	keys map[uint32]Balancer.Peer
	// 得到的节点状态是否可用
	peers map[Balancer.Peer]bool
	rw    sync.RWMutex
}

// Hasher 代表可选策略
type Hasher func(data []byte) uint32

//  New 使用 HashKetama 创建一个新的负载均衡器实例
func New(opts ...Balancer.Opt) Balancer.Balancer {

	return (&HashKetama{
		hasher:  crc32.ChecksumIEEE,
		replica: 32,
		keys:    make(map[uint32]Balancer.Peer),
		peers:   make(map[Balancer.Peer]bool),
	}).init(opts...)
}

// 典型的 “把不同参数类型的函数包装成为相同参数类型的函数”

// WithHashFunc allows a custom hash function to be specified.
// The default Hasher hash func is crc32.ChecksumIEEE.
func WithHashFunc(hashFunc Hasher) Balancer.Opt {
	return func(balancer Balancer.Balancer) {
		if l, ok := balancer.(*HashKetama); ok {
			l.hasher = hashFunc
		}
	}
}

// WithReplica allows a custom replica number to be specified.
// The default replica number is 32.
func WithReplica(replica int) Balancer.Opt {
	return func(balancer Balancer.Balancer) {
		if l, ok := balancer.(*HashKetama); ok {
			l.replica = replica
		}
	}
}

// 让 HashKetama 指针穿过一系列的Opt函数
func (s *HashKetama) init(opts ...Balancer.Opt) *HashKetama {
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Balancer.Factor本质是 string 类型
// 调用Factor()转化为string 类型
// 让 HashKetama实现了Balancer.Balancer接口是一个具体的负载均衡器
// 所有的HashKetama都会接收类型为Balancer.Factor的实例，Balancer.Factor的实例
func (s *HashKetama) Next(factor Balancer.Factor) (next Balancer.Peer, c Balancer.Constrainable) {
	var hash uint32
	// 生成哈希code
	if h, ok := factor.(Balancer.FactorHashable); ok {
		// 如果传入的是具体的实现了Balancer.FactorHashable接口的类
		// 那么肯定实现了具体的HashCode()函数，调用就ok了
		hash = h.HashCode()
	} else {
		// 如果只是传入了实现了父类接口的类的实例
		// 调用hasher 处理父类实例
		// factor.Factor() 把请求"https://abc.local/user/profile"
		hash = s.hasher([]byte(factor.Factor()))
		fmt.Printf("Hash code is %d\n", hash)
		// s.hasher([]byte(factor.Factor()))本质是 crc32.ChecksumIEEE()函数处理得到的[]byte类型的string
		// 所以重点是crc32.ChecksumIEEE()如何把[]byte转化wei hash code 的
		// 哈希Hash，就是把任意长度的输入，通过散列算法，变换成固定长度的输出，该输出就是散列值。
		// 不定长输入-->哈希函数-->定长的散列值
		// 哈希算法的本质是对原数据的有损压缩
		/* CRC检验原理实际上就是在一个p位二进制数据序列之后附加一个r位二进制检验码(序列)，
		从而构成一个总长为n＝p＋r位的二进制序列；附加在数据序列之后的这个检验码与数据序列的内容之间存在着某种特定的关系。
		如果因干扰等原因使数据序列中的某一位或某些位发生错误，这种特定关系就会被破坏。因此，通过检查这一关系，就可以实现对数据正确性的检验
		注：仅用循环冗余检验 CRC 差错检测技术只能做到无差错接受（只是非常近似的认为是无差错的），并不能保证可靠传输
		*/
	}

	// 根据具体的策略得到下标
	next = s.miniNext(hash)
	if next != nil {
		if fc, ok := factor.(Balancer.FactorComparable); ok {
			next, c, _ = fc.ConstrainedBy(next)
		} else if nested, ok := next.(Balancer.BalancerLite); ok {
			next, c = nested.Next(factor)
		}
	}

	return
}

// 已经有存储着一些哈希数值的切片
// 产生哈希数值
// 在切片中找到大于等于得到的哈希数值的元素
// 该元素作为map的key一定可以找到一个节点

func (s *HashKetama) miniNext(hash uint32) (next Balancer.Peer) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	// 得到的hashcode 去和 hashRing[i]比较
	// sort.Search()二分查找 本质: 找到满足条件的最小的索引
	/*
		//golang 官方的二分写法 (学习一波)

		func Search(n int, f func(int) bool) int {
			// Define f(-1) == false and f(n) == true.
			// Invariant: f(i-1) == false, f(j) == true.
			i, j := 0, n
			for i < j {
				// avoid overflow when computing h
				// 右移一位 相当于除以2
				h := int(uint(i+j) >> 1)
				// i ≤ h < j
				if !f(h) {
					i = h + 1 // preserves f(i-1) == false
				} else {
					j = h // preserves f(j) == true
				}
			}
			// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
			return i
		}
	*/

	// 在s.hashRing找到大于等于hash的hashRing的下标
	ix := sort.Search(len(s.hashRing), func(i int) bool {
		return s.hashRing[i] >= hash
	})

	// 当这个下标是最后一个下标时，相当于没有找到
	if ix == len(s.hashRing) {
		ix = 0
	}

	// 如果没有找到就返回s.hashRing的第一个元素
	hashValue := s.hashRing[ix]

	// s.keys 存储 peers 每一个 peers 都有一个hashValue 对应
	// hashcode 对应 hashValue （被Slice存储）
	// hashValue 对应节点 peer  (被Map存储)
	if p, ok := s.keys[hashValue]; ok {
		if _, ok = s.peers[p]; ok {
			next = p
		}
	}

	return
}

/*
在 Add 实现中建立了 hashRing 结构，
它虽然是环形，但是是以数组和下标取模的方式来达成的。
此外，keys 这个 map 解决从 peer 的 hash 值到 peer 的映射关系，今后（在 Next 中）就可以通过从 hashRing 上 pick 出一个 point 之后立即地获得相应的 peer.
在 Next 中主要是在做 factor 的 hash 值计算，计算的结果在 hashRing 上映射为一个点 pt，如果不是恰好有一个 peer 被命中的话，就向后扫描离 pt 最近的 peer。

*/
func (s *HashKetama) Count() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.peers)
}

func (s *HashKetama) Add(peers ...Balancer.Peer) {
	s.rw.Lock()
	defer s.rw.Unlock()

	for _, p := range peers {
		s.peers[p] = true
		for i := 0; i < s.replica; i++ {
			hash := s.hasher(s.peerToBinaryID(p, i))
			s.hashRing = append(s.hashRing, hash)
			s.keys[hash] = p
		}
	}

	sort.Slice(s.hashRing, func(i, j int) bool {
		return s.hashRing[i] < s.hashRing[j]
	})
}

func (s *HashKetama) peerToBinaryID(p Balancer.Peer, replica int) []byte {
	str := fmt.Sprintf("%v-%05d", p, replica)
	return []byte(str)
}

func (s *HashKetama) Remove(peer Balancer.Peer) {
	s.rw.Lock()
	defer s.rw.Unlock()

	if _, ok := s.peers[peer]; ok {
		delete(s.peers, peer)
	}

	var keys []uint32
	var km = make(map[uint32]bool)
	for i, p := range s.keys {
		if p == peer {
			keys = append(keys, i)
			km[i] = true
		}
	}

	for _, key := range keys {
		delete(s.keys, key)
	}

	var vn []uint32
	for _, x := range s.hashRing {
		if _, ok := km[x]; !ok {
			vn = append(vn, x)
		}
	}
	s.hashRing = vn
}

func (s *HashKetama) Clear() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.hashRing = nil
	s.keys = make(map[uint32]Balancer.Peer)
	s.peers = make(map[Balancer.Peer]bool)
}
