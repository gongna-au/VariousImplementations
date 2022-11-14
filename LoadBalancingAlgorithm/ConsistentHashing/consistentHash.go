package ConsistentHashing

import (
	"fmt"
	"hash/crc32"

	"github.com/golang/protobuf/ptypes/any"
)

type Hash func(data []byte) int32
type UInt32Slice []uint32

// 为了解决数据倾斜的问题，为节点添加虚拟节点

type ConsistentHashBalance struct {
	hash Hash
	// 虚拟节点的个数
	replicas int
	// 已经排序的节点的hash 切片
	keys    UInt32Slice
	hashMap map[uint32]string
}

//利用复制因子和哈希函数创建一个一致性哈希算法
func NewConsistentHashBalance(replica int)
