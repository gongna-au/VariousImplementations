package SmoothWeightRoundRobin

import (
	"strings"
)

type Node struct {
	Name    string
	Current int
	Weight  int
}

// 一次负载均衡的选择 找到最大的节点，把最大的节点减去权重量和
// 算法的核心是current 记录找到权重最大的节点，这个节点的权重-总权重
// 然后在这个基础上的切片 他们的状态是 现在的权重状态+最初的权重状态
func SmoothWeightRoundRobin(nodes []*Node) (best *Node) {
	if len(nodes) == 0 {
		return nil
	}
	weightnum := 0
	for k, v := range nodes {
		weightnum = weightnum + v.Weight
		if k == 0 {
			best = v
		}
		if v.Current > best.Current {
			best = v
		}
	}
	for _, v := range nodes {
		if strings.Compare(v.Name, best.Name) == 0 {
			v.Current = v.Current - weightnum + v.Weight
		} else {
			v.Current = v.Current + v.Weight
		}
	}

	return best
}
