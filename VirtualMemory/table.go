package VirtualMemory

import "fmt"

type PageTableItem struct {
	// 页号
	PageId int
	// 物理⻚号
	PhysicalPageId int
	// 是否在内存
	IsInMemory bool
	// 在硬盘中的地址
	LocationInHardDisk int
	// 被调入内存使用的次数
	UsedNum int
	// 被调入内存是否被修改
	IsChanged bool
}

func ProgramWantLoadPage() {
	// 模拟程序要访问100这个页面 load 100
	var item = PageTableItem{

		PageId:             50,
		PhysicalPageId:     100,
		IsInMemory:         false,
		LocationInHardDisk: 100,
		UsedNum:            0,
		IsChanged:          false,
	}
	var tableLength int = 100
	// 模拟保存页面的页表 保存了100个页面 TLB 中间的快表
	table := make(map[int]int, 100)
	// 假设我们要查找的这个表就在快表里面
	table[3] = 50
	if item.PageId > tableLength {
		fmt.Println("内存越界异常")
	} else {
		k, ok := table[page]
		if ok {

		}
	}
}
