package VirtualMemory

import (
	"fmt"
	"strconv"
	"sync"
)

type ActualCallPages struct {
	Pages []string
}

func Create(pages ...string) *ActualCallPages {
	return &ActualCallPages{
		Pages: append([]string{}, pages...),
	}
}

func (a *ActualCallPages) Travese() {
	fmt.Println("Actually System Call Page Order:")
	for _, v := range a.Pages {
		fmt.Print(v, "  ")
	}
}

type HardDisk struct {
	maxPageNum int64
	Data       map[string]*ByteData
}

type ByteData struct {
	content any
}

func NewByteData(content any) *ByteData {
	return &ByteData{
		content: content,
	}
}

var hardDisk *HardDisk
var once sync.Once

// 用最大的页数来初始化磁盘
// 每一页都有对应的数据
func HardDiskInit(maxPageNum int64, data ...*ByteData) *HardDisk {
	once.Do(func() {
		result := make(map[string]*ByteData, maxPageNum)
		for k, v := range data {
			result[strconv.Itoa(k)] = v
		}
		hardDisk = &HardDisk{
			maxPageNum: maxPageNum,
			Data:       result,
		}
	})
	return hardDisk
}

func GetHardDisk() *HardDisk {
	return hardDisk
}

func (h *HardDisk) GetData(key string) string {
	if v, ok := h.Data[key]; ok {
		return fmt.Sprintf("%s", v.content)
	} else {
		return "0"
	}
}

// 每个页面对应的具体数据被放入
type Node struct {
	Key   string //每个节点的唯一标识，作为key储存到lru的cache里
	Value []byte //携带的数据
}

func NewNode(key string, node []byte) Node {
	return Node{
		Key:   key,
		Value: node,
	}
}

func (n *Node) PrintlnNode() {
	if n != nil {
		fmt.Println("Node key: ", n.Key, " Node value:", string(n.Value))
	}
}
