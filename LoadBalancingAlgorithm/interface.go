package LoadBalancingAlgorithm

import "reflect"

// -----------------------------------------------------
// 随机负载均衡
// Peer 是一个后端节点。
type Peer interface {
	// String 将返回一个后端节点的主要身份字符串。
	// 所有的后端需要实现这一个函数即可插入到负载均衡的架构当中
	String() string
}

// 实现Peer 接口
type ExP string

func (s ExP) String() string {
	return string(s)
}

// -----------------------------------------------------
// 加权负载均衡接口
// WeightedPeer 是一个具有权重属性的 Peer。
// 该接口i要求你既要实现String() string方法
// 又要实现Weight() int 方法
type WeightedPeer interface {
	Peer
	Weighted
}

// Weighted 是一个权重接口
type Weighted interface {
	Weight() int
}

// BalancerLite 代表一个通用的负载均衡器
type BalancerLite interface {
	Next(factor Factor) (next Peer, c Constrainable)
}

// -----------------------------------------------------
// 现实的负载均衡器接口
// Balancer 代表一个通用的负载均衡器。
// 对于现实世界，Balancer 是一个有用的接口，而不是 BalancerLite。
type Balancer interface {
	BalancerLite
	Count() int
	Add(peers ...Peer)
	Remove(peer Peer)
	Clear()
}

// Opt 是 New Balancer 的类型原型
type Opt func(balancer Balancer)

// Factor 是 BalancerLite.Next 的因子参数。
// 如果你不知道应该传入什么
// BalancerLite.Next，发送 DummyFactor 作为它。
// 但是当你在构建一个复杂的 lb 系统时
// 或者一个多级的，例如加权的
// 版本控制负载均衡器。看看
// version.New 和 version.VersioningBackendFactor
type Factor interface {
	Factor() string
}

// 实现Factor 给Factor接口类型提供一个默认的实现类
// 默认类
type FactorString string

func (s FactorString) Factor() string {
	return string(s)
}

// 当你身为调度者时，想要调用 Next，却没有什么合适的“因素”提供的话，就提供 DummyFactor 好了。
// 传入给BalancerLite接口的Next(factor Factor)方法进行调用
// 默认类 FactorString的实例
const DummyFactor FactorString = ""

//FactorComparable是组装Factor和约束比较的复合接口
type FactorComparable interface {
	Factor
	ConstrainedBy(constraints interface{}) (peer Peer, c Constrainable, satisfied bool)
}

// FactorHashable是一个组装Factor和散列计算的复合接口。
type FactorHashable interface {
	Factor
	HashCode() uint32
}

// Constrainable 目前将其视而不见就足够了。其实是一个可以应用于 BalancerLite.Next(factor)
type Constrainable interface {
	CanConstrain(o interface{}) (yes bool)
	Check(o interface{}) (satisfied bool)
	Peer
}

// 其实是一个可以应用于 BalancerLite.Next(factor)的类
type WeightedConstrainable interface {
	Constrainable
	Weighted
}

// DeepEqualAware 可以由 Peer 具体化，这样你就可以
// 自定义如何比较两个节点，避免使用 reflect.DeepEqual旁路。
type DeepEqualAware interface {
	DeepEqual(b Peer) bool
}

func DeepEqual(a, b Peer) (yes bool) {
	if a == b {
		return true
	}

	if e, ok := a.(DeepEqualAware); ok {
		return e.DeepEqual(b)
	}

	return reflect.DeepEqual(a, b)
}
