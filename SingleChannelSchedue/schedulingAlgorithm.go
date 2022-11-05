package SingleChannelSchedue

import (
	"errors"
	/* "fmt"
	"sort"
	"strconv"
	"time" */

	"github.com/VariousImplementations/SecondarySchedue/pkg"
)

var AlgorithmFactory map[string]Algorithm

func init() {
	AlgorithmFactory = map[string]Algorithm{
		"FCFS": NewFCFS(),
		"HPF":  NewHPF(),
		"HRRF": NewHRRF(),
		"SJF":  NewShortJobFirst(),
	}
}

// 典型的迭代器模式
type ICollection interface {
	CreateAlgorithm(string) (Algorithm, error)
}

// 集合可以创建出某一种具体的算法
// 然后调用该算法进行
type WorkCollection struct {
	Works []*pkg.Work
}

func NewWorkCollection(w ...*pkg.Work) *WorkCollection {
	return &WorkCollection{
		Works: w,
	}
}

func (w *WorkCollection) CreateAlgorithm(key string) (Algorithm, error) {
	v, ok := AlgorithmFactory[key]
	if ok {
		return v, nil
	} else {
		return nil, errors.New("This Algorithm is not defined!")
	}
}

// 算法实现的调度接口
type Algorithm interface {
	SetCollection(i ICollection)
	Schedule()
}

//先来先服务（First-Come First-Served，FCFS）调度算法
type FCFS struct {
	Collection      ICollection
	WorkInMemory    []*pkg.Work
	WorkNotInMemory []*pkg.Work
	WorkHasFinished []*pkg.Work
	MaxNumInMemory  int
}

func NewFCFS(i ICollection, max int) *FCFS {
	w, ok := i.(*WorkCollection)
	if ok {
	}
	return &FCFS{
		Collection:      i,
		WorkInMemory:    make([]*pkg.Work, max),
		WorkNotInMemory: make([]*pkg.Work, len(w.Works)-1),
		MaxNumInMemory:  max,
	}
}

func (f *FCFS) Init() {
	w, ok := f.Collection.(*WorkCollection)
	if ok {

	}

}

// 最重要的：先把无序的切片变成有序的切片
func (f *FCFS) Schedule() {
	w, _ := f.Collection.(*WorkCollection)

}

//优先权高者优先（Highest-Priority-First，HPF）调度算法
type HPF struct {
	Collection ICollection
}

func NewHPF() *HPF {
	return &HPF{}
}

func (h *HPF) SetCollection(i ICollection) {
	h.Collection = i
}

func (h *HPF) Schedule() {

}

func HighestPrivligeNext(current int, works []*pkg.Work, isscheduled map[int]int) int {
	// 在当前时间点有两种类型的作业，一种是已经等待了一会的作业和一分钟都还没等待的作业

}

// 响应比高者优先（HRRF）调度算法
type HRRF struct {
	Collection ICollection
}

func NewHRRF() *HRRF {
	return &HRRF{}
}

func (h *HRRF) SetCollection(i ICollection) {
	h.Collection = i
}

func (h *HRRF) Schedule() {

}

//  作业周转时间 /作业处理时间
func HighestResponseRatioNext(current int, works []*pkg.Work, isscheduled map[int]int) int {

}

// 短作业优先
type ShortJobFirst struct {
	Collection ICollection
}

func NewShortJobFirst() *ShortJobFirst {
	return &ShortJobFirst{}
}

func (s *ShortJobFirst) SetCollection(i ICollection) {
	s.Collection = i
}

func (s *ShortJobFirst) Schedule() {

}

// 原始的works 切片里面
func ShortestWork(current int, works []*pkg.Work, isscheduled map[int]int) int {

}

func PrintWorks(works []*pkg.Work) {

}

func InItscheduledMap(status map[int]int, length int) {
	for i := 0; i < length; i++ {
		status[i] = 0
	}
}

func ClientHRRF() {

}

func ClientFCFS() {

}
func ClientSJF() {

}
func ClientHPF() {

}
