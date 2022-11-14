package LRU

import (
	"fmt"
	"testing"
)

func TestGetNode(t *testing.T) {

	k1, k2, k3, k4, k5, k6, k7, k8, k9, k10, k11, k12 := "4", "3", "2", "1", "4", "3", "5", "4", "3", "2", "1", "5"
	lru := New(3)
	fmt.Println("After Schedele page fail num is: ", lru.Schedule(k1, k2, k3, k4, k5, k6, k7, k8, k9, k10, k11, k12))
}
func TestLRUSchedule(t *testing.T) {
	// 4            1次缺页中断
	// 3 4			1次缺页中断
	// 2 3 4 		1次缺页中断
	// 1 2 3        1次缺页中断
	// 4 1 2        1次缺页中断
	// 3 4 1  		1次缺页中断
	// 5 3 4 		1次缺页中断
	// 4 5 3
	// 3 4 5
	// 2 3 4  		1次缺页中断
	// 1 2 3		1次缺页中断
	// 5 1 2        1次缺页中断
	// 一共10
	k1, k2, k3, k4, k5, k6, k7, k8, k9, k10, k11, k12 := "4", "3", "2", "1", "4", "3", "5", "4", "3", "2", "1", "5"
	lru := New(3)
	fmt.Println("After Schedele page fail num is: ", lru.Schedule(k1, k2, k3, k4, k5, k6, k7, k8, k9, k10, k11, k12))
}
