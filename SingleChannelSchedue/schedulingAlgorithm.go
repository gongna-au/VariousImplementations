package SingleChannelSchedue

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
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
	Works []*Work
}

func NewWorkCollection(w ...*Work) *WorkCollection {
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
	Collection ICollection
}

func NewFCFS() *FCFS {
	return &FCFS{}
}

func (f *FCFS) SetCollection(i ICollection) {
	f.Collection = i
}

// 最重要的：先把无序的切片变成有序的切片
func (f *FCFS) Schedule() {
	w, ok := f.Collection.(*WorkCollection)
	fmt.Printf("输入的信息：\n")
	for _, v := range w.Works {
		fmt.Printf("Id:%-8d ArriveTime:%-8d ExcuteTime:%-8d PrivilegeLevel:%-8d\n", v.Id, v.ArriveTime, v.ExcuteTime, v.Level)
	}

	if ok {
		sort.Slice(w.Works, func(i, j int) bool {
			if w.Works[i].ArriveTime < w.Works[j].ArriveTime {
				return true
			} else {
				return false
			}
		})
	} else {
		fmt.Println("ICollection Interface can not convert to *WorkCollection")
	}

	waitTime := 0
	roundTime := 0
	fmt.Printf("FCFS 作业调度算法：\n")
	for k, v := range w.Works {
		if k == 0 {
			v.StartTime = v.ArriveTime
			v.WaitTime = 0
			v.RoundTime = v.ExcuteTime
			v.OverTime = v.StartTime + v.ExcuteTime
			waitTime = waitTime + v.WaitTime
			roundTime = roundTime + v.RoundTime
			fmt.Printf("Id:%-8d ArriveTime:%-8d StartTime:%-8d WaitTime:%-8d RoundTime:%-8d\n", v.Id, v.ArriveTime, v.StartTime, v.WaitTime, v.RoundTime)
			continue
		}

		if v.ArriveTime >= w.Works[k-1].OverTime {
			v.StartTime = v.ArriveTime

			v.WaitTime = 0
			v.OverTime = v.StartTime + v.ExcuteTime
			v.RoundTime = v.OverTime - v.ArriveTime

		} else {
			v.StartTime = w.Works[k-1].OverTime
			v.WaitTime = v.StartTime - v.ArriveTime
			v.OverTime = v.StartTime + v.ExcuteTime
			v.RoundTime = v.OverTime - v.ArriveTime
		}

		waitTime = waitTime + v.WaitTime
		roundTime = roundTime + v.RoundTime
		fmt.Printf("Id:%-8d ArriveTime:%-8d StartTime:%-8d WaitTime:%-8d RoundTime:%-8d\n", v.Id, v.ArriveTime, v.StartTime, v.WaitTime, v.RoundTime)
	}
	fmt.Printf("总等待时间:%-8d 总周转时间:%-8d\n", waitTime, roundTime)
	fmt.Printf("平均等待时间: %4.2f 平均周转时间: %4.2f\n", float64(waitTime)/float64(len(w.Works)), float64(roundTime)/float64(len(w.Works)))

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

	w, ok := h.Collection.(*WorkCollection)
	fmt.Printf("输入的信息：\n")
	for _, v := range w.Works {
		fmt.Printf("Id:%-8d ArriveTime:%-8d ExcuteTime:%-8d PrivilegeLevel:%-8d\n", v.Id, v.ArriveTime, v.ExcuteTime, v.Level)
	}
	if !ok {
		fmt.Println("ICollection Interface can not convert to *WorkCollection")
	}
	fmt.Printf("HPF 作业调度算法：\n")
	currentTime := 0
	length := len(w.Works)
	isscheduled := make(map[int]int, length)
	InItscheduledMap(isscheduled, length)
	for i := 0; i < length; i++ {
		currentTime = HighestPrivligeNext(currentTime, w.Works, isscheduled)
	}
	sort.Slice(w.Works, func(i, j int) bool {
		if w.Works[i].StartTime < w.Works[j].StartTime {
			return true
		} else {
			return false
		}
	})

	waitTime := 0
	roundTime := 0
	for _, v := range w.Works {
		waitTime = waitTime + v.WaitTime
		roundTime = roundTime + v.RoundTime
		fmt.Printf("Id:%-8d ArriveTime:%-8d StartTime:%-8d WaitTime:%-8d RoundTime:%-8d\n", v.Id, v.ArriveTime, v.StartTime, v.WaitTime, v.RoundTime)
	}
	fmt.Printf("总等待时间:%-8d 总周转时间:%-8d\n", waitTime, roundTime)
	fmt.Printf("平均等待时间: %4.2f 平均周转时间: %4.2f\n", float64(waitTime)/float64(len(w.Works)), float64(roundTime)/float64(len(w.Works)))

}

func HighestPrivligeNext(current int, works []*Work, isscheduled map[int]int) int {
	// 在当前时间点有两种类型的作业，一种是已经等待了一会的作业和一分钟都还没等待的作业
	// 已经到达的任务
	beforeCurrentTime := []*Work{}
	beforeCurrentTimeIndexinworks := []int{}
	for k, v := range works {
		if isscheduled[k] == 1 {
			continue
		} else {
			if (current - v.ArriveTime) > 0 {
				beforeCurrentTime = append(beforeCurrentTime, v)
				beforeCurrentTimeIndexinworks = append(beforeCurrentTimeIndexinworks, k)
			}
		}
	}
	// 起始时间等于或者在当前时间之后的
	arriveFirst := 0
	if len(beforeCurrentTime) == 0 {
		// 找到达时间最早的
		for k, v := range works {
			if isscheduled[k] == 1 {
				continue
			}
			if v.ArriveTime < works[arriveFirst].ArriveTime {
				arriveFirst = k
			} else {
				continue
			}
		}
	}
	//fmt.Printf("arriveFirst ID is %-4d", arriveFirst)
	if len(beforeCurrentTime) == 0 {
		// 开始的时间为到达的时间
		works[arriveFirst].StartTime = works[arriveFirst].ArriveTime
		// 等待时间
		works[arriveFirst].WaitTime = 0
		//周转时间
		works[arriveFirst].RoundTime = works[arriveFirst].ExcuteTime
		// 结束时间
		works[arriveFirst].OverTime = works[arriveFirst].ArriveTime + works[arriveFirst].ExcuteTime
		// 最高响应比 周转时间/执行之间
		works[arriveFirst].Excellent = 1
		// 调度结束
		isscheduled[arriveFirst] = 1
		// 当前时间为结束时间
		current = works[arriveFirst].OverTime
		return current
	} else {
		// 已经有先来的任务，假设第一个响应比最大
		// tempResponseRatio
		highestLevelIndex := 0
		max := -1
		for k, v := range beforeCurrentTime {
			if v.Level > max {
				highestLevelIndex = k
				max = v.Level
			}
		}
		//fmt.Printf("BeforeWorksIndex:is %d\n", shortestTimeIndex)
		workIndex := beforeCurrentTimeIndexinworks[highestLevelIndex]
		//fmt.Printf("WorksIndex:is %d\n", workIndex)
		beforeCurrentTime[highestLevelIndex].StartTime = current
		beforeCurrentTime[highestLevelIndex].OverTime = beforeCurrentTime[highestLevelIndex].StartTime + beforeCurrentTime[highestLevelIndex].ExcuteTime
		beforeCurrentTime[highestLevelIndex].WaitTime = beforeCurrentTime[highestLevelIndex].StartTime - beforeCurrentTime[highestLevelIndex].ArriveTime
		beforeCurrentTime[highestLevelIndex].RoundTime = beforeCurrentTime[highestLevelIndex].OverTime - beforeCurrentTime[highestLevelIndex].ArriveTime
		excellent, _ := strconv.ParseFloat(fmt.Sprintf("%.5f", float64(beforeCurrentTime[highestLevelIndex].RoundTime)/float64(beforeCurrentTime[highestLevelIndex].ExcuteTime)), 64)
		beforeCurrentTime[highestLevelIndex].Excellent = excellent
		isscheduled[workIndex] = 1
		return beforeCurrentTime[highestLevelIndex].OverTime
	}

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
	w, ok := h.Collection.(*WorkCollection)
	fmt.Printf("输入的信息：\n")
	for _, v := range w.Works {
		fmt.Printf("Id:%-8d ArriveTime:%-8d ExcuteTime:%-8d PrivilegeLevel:%-8d\n", v.Id, v.ArriveTime, v.ExcuteTime, v.Level)
	}
	if !ok {
		fmt.Println("ICollection Interface can not convert to *WorkCollection")
	}
	fmt.Printf("HRRF 作业调度算法：\n")
	currentTime := 0
	length := len(w.Works)
	// 标记是否已经被调度
	// 细节之一 map和切片不同的是， map 在分配了大小之后，因为你没有给他key值，所以不会像切片那样 {0,0,0,0,0} 更像是一个有大小的空盒
	isscheduled := make(map[int]int, length)
	InItscheduledMap(isscheduled, length)
	for i := 0; i < length; i++ {
		currentTime = HighestResponseRatioNext(currentTime, w.Works, isscheduled)
	}
	sort.Slice(w.Works, func(i, j int) bool {
		if w.Works[i].StartTime < w.Works[j].StartTime {
			return true
		} else {
			return false
		}
	})

	waitTime := 0
	roundTime := 0
	for _, v := range w.Works {
		waitTime = waitTime + v.WaitTime
		roundTime = roundTime + v.RoundTime
		fmt.Printf("Id:%-8d ArriveTime:%-8d StartTime:%-8d WaitTime:%-8d RoundTime:%-8d\n", v.Id, v.ArriveTime, v.StartTime, v.WaitTime, v.RoundTime)
	}
	fmt.Printf("总等待时间:%-8d 总周转时间:%-8d\n", waitTime, roundTime)
	fmt.Printf("平均等待时间: %4.2f 平均周转时间: %4.2f\n", float64(waitTime)/float64(len(w.Works)), float64(roundTime)/float64(len(w.Works)))
}

//  作业周转时间 /作业处理时间
func HighestResponseRatioNext(current int, works []*Work, isscheduled map[int]int) int {
	// 在当前时间点有两种类型的作业，一种是已经等待了一会的作业和一分钟都还没等待的作业
	// 已经到达的任务
	beforeCurrentTime := []*Work{}
	beforeCurrentTimeIndexinworks := []int{}
	for k, v := range works {
		if isscheduled[k] == 1 {
			continue
		} else {
			if (current - v.ArriveTime) > 0 {
				beforeCurrentTime = append(beforeCurrentTime, v)
				beforeCurrentTimeIndexinworks = append(beforeCurrentTimeIndexinworks, k)
			}
		}
	}
	// 起始时间等于或者在当前时间之后的
	arriveFirst := 0
	if len(beforeCurrentTime) == 0 {
		// 找到达时间最早的
		for k, v := range works {
			if isscheduled[k] == 1 {
				continue
			}
			if v.ArriveTime < works[arriveFirst].ArriveTime {
				arriveFirst = k
			} else {
				continue
			}
		}
	}
	//fmt.Printf("arriveFirst ID is %-4d", arriveFirst)
	if len(beforeCurrentTime) == 0 {
		// 开始的时间为到达的时间
		works[arriveFirst].StartTime = works[arriveFirst].ArriveTime
		// 等待时间
		works[arriveFirst].WaitTime = 0
		//周转时间
		works[arriveFirst].RoundTime = works[arriveFirst].ExcuteTime
		// 结束时间
		works[arriveFirst].OverTime = works[arriveFirst].ArriveTime + works[arriveFirst].ExcuteTime
		// 最高响应比 周转时间/执行之间
		works[arriveFirst].Excellent = 1
		// 调度结束
		isscheduled[arriveFirst] = 1
		// 当前时间为结束时间
		current = works[arriveFirst].OverTime
		return current
	} else {
		// 已经有先来的任务，假设第一个响应比最大
		// tempResponseRatio
		highestResponseRatioIndex := 0
		max := 1.0
		for k, v := range beforeCurrentTime {
			//fmt.Printf("current time is %d", current)
			//fmt.Printf("arriveTime time is %d", v.ArriveTime)
			//fmt.Printf("excuteTime time is %d", v.ExcuteTime)
			//fmt.Printf("分子：%d", current-v.ArriveTime+v.ExcuteTime)
			//fmt.Printf("responseRation is:%4f\n", float64((current-v.ArriveTime+v.ExcuteTime)/v.ExcuteTime))
			value, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", float64(current-v.ArriveTime+v.ExcuteTime)/float64(v.ExcuteTime)), 64)
			//fmt.Printf("responseRation is:%4f\n", value)
			if value > max {
				highestResponseRatioIndex = k
				max = value
			}
		}
		//fmt.Printf("BeforeWorksIndex:is %d\n", shortestTimeIndex)
		workIndex := beforeCurrentTimeIndexinworks[highestResponseRatioIndex]
		//fmt.Printf("WorksIndex:is %d\n", workIndex)
		beforeCurrentTime[highestResponseRatioIndex].StartTime = current
		beforeCurrentTime[highestResponseRatioIndex].OverTime = beforeCurrentTime[highestResponseRatioIndex].StartTime + beforeCurrentTime[highestResponseRatioIndex].ExcuteTime
		beforeCurrentTime[highestResponseRatioIndex].WaitTime = beforeCurrentTime[highestResponseRatioIndex].StartTime - beforeCurrentTime[highestResponseRatioIndex].ArriveTime
		beforeCurrentTime[highestResponseRatioIndex].RoundTime = beforeCurrentTime[highestResponseRatioIndex].OverTime - beforeCurrentTime[highestResponseRatioIndex].ArriveTime
		//excellent, _ := strconv.ParseFloat(fmt.Sprintf("%.5f", float64(beforeCurrentTime[highestResponseRatioIndex].RoundTime)/float64(beforeCurrentTime[highestResponseRatioIndex].ExcuteTime)), 64)
		beforeCurrentTime[highestResponseRatioIndex].Excellent = max
		isscheduled[workIndex] = 1
		return beforeCurrentTime[highestResponseRatioIndex].OverTime
	}

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

	w, ok := s.Collection.(*WorkCollection)
	fmt.Printf("输入的信息：\n")
	for _, v := range w.Works {
		fmt.Printf("Id:%-8d ArriveTime:%-8d ExcuteTime:%-8d PrivilegeLevel:%-8d\n", v.Id, v.ArriveTime, v.ExcuteTime, v.Level)
	}
	if !ok {
		fmt.Println("ICollection Interface can not convert to *WorkCollection")
	}

	fmt.Printf("SJF 作业调度算法：\n")
	currentTime := 0
	length := len(w.Works)
	// 标记是否已经被调度
	// 细节之一 map和切片不同的是， map 在分配了大小之后，因为你没有给他key值，所以不会像切片那样 {0,0,0,0,0} 更像是一个有大小的空盒
	isscheduled := make(map[int]int, length)
	InItscheduledMap(isscheduled, length)
	for i := 0; i < length; i++ {
		currentTime = ShortestWork(currentTime, w.Works, isscheduled)
	}

	sort.Slice(w.Works, func(i, j int) bool {
		if w.Works[i].StartTime < w.Works[j].StartTime {
			return true
		} else {
			return false
		}
	})
	waitTime := 0
	roundTime := 0
	for _, v := range w.Works {
		waitTime = waitTime + v.WaitTime
		roundTime = roundTime + v.RoundTime
		fmt.Printf("Id:%-8d ArriveTime:%-8d StartTime:%-8d WaitTime:%-8d RoundTime:%-8d\n", v.Id, v.ArriveTime, v.StartTime, v.WaitTime, v.RoundTime)
	}
	fmt.Printf("总等待时间:%-8d 总周转时间:%-8d\n", waitTime, roundTime)
	fmt.Printf("平均等待时间: %4.2f 平均周转时间: %4.2f\n", float64(waitTime)/float64(len(w.Works)), float64(roundTime)/float64(len(w.Works)))
}

// 原始的works 切片里面
func ShortestWork(current int, works []*Work, isscheduled map[int]int) int {

	// 在当前时间点有两种类型的作业，一种是已经等待了一会的作业和一分钟都还没等待的作业
	// 已经到达的任务
	beforeCurrentTime := []*Work{}
	beforeCurrentTimeIndexinworks := []int{}
	for k, v := range works {
		if isscheduled[k] == 1 {
			continue
		} else {
			if (current - v.ArriveTime) > 0 {
				beforeCurrentTime = append(beforeCurrentTime, v)
				beforeCurrentTimeIndexinworks = append(beforeCurrentTimeIndexinworks, k)
			}
		}

	}
	// 起始时间等于或者在当前时间之后的
	arriveFirst := 0
	if len(beforeCurrentTime) == 0 {
		// 找到达时间最早的
		for k, v := range works {
			if isscheduled[k] == 1 {
				continue
			}
			if v.ArriveTime < works[arriveFirst].ArriveTime {
				arriveFirst = k
			} else {
				continue
			}
		}

	}
	//fmt.Printf("arriveFirst ID is %-4d", arriveFirst)
	if len(beforeCurrentTime) == 0 {

		// 开始的时间为到达的时间
		works[arriveFirst].StartTime = works[arriveFirst].ArriveTime
		// 等待时间
		works[arriveFirst].WaitTime = 0
		//周转时间
		works[arriveFirst].RoundTime = works[arriveFirst].ExcuteTime
		// 结束时间
		works[arriveFirst].OverTime = works[arriveFirst].ArriveTime + works[arriveFirst].ExcuteTime
		// 调度结束
		isscheduled[arriveFirst] = 1
		// 当前时间为结束时间
		current = works[arriveFirst].OverTime
		//fmt.Printf("current Time is %-4d", current)
		return current
	} else {
		// 已经有先来的任务,假设第一个时间最短
		shortestTimeIndex := 0
		/* fmt.Print("BeforeWorkis:\n")
		PrintWorks(beforeCurrentTime) */
		for k, v := range beforeCurrentTime {
			if v.ExcuteTime < beforeCurrentTime[shortestTimeIndex].ExcuteTime {
				shortestTimeIndex = k
			}
		}
		//fmt.Printf("BeforeWorksIndex:is %d\n", shortestTimeIndex)
		workIndex := beforeCurrentTimeIndexinworks[shortestTimeIndex]
		//fmt.Printf("WorksIndex:is %d\n", workIndex)
		beforeCurrentTime[shortestTimeIndex].StartTime = current
		beforeCurrentTime[shortestTimeIndex].OverTime = beforeCurrentTime[shortestTimeIndex].StartTime + beforeCurrentTime[shortestTimeIndex].ExcuteTime
		beforeCurrentTime[shortestTimeIndex].WaitTime = beforeCurrentTime[shortestTimeIndex].StartTime - beforeCurrentTime[shortestTimeIndex].ArriveTime
		beforeCurrentTime[shortestTimeIndex].RoundTime = beforeCurrentTime[shortestTimeIndex].OverTime - beforeCurrentTime[shortestTimeIndex].ArriveTime
		isscheduled[workIndex] = 1
		//fmt.Printf("changed is %v", *beforeCurrentTime[shortestTimeIndex])
		return beforeCurrentTime[shortestTimeIndex].OverTime
	}

}

func PrintWorks(works []*Work) {
	for _, v := range works {
		fmt.Print(*v)
		fmt.Print("  ")
	}
}

func InItscheduledMap(status map[int]int, length int) {
	for i := 0; i < length; i++ {
		status[i] = 0
	}
}

// 具体的work
type Work struct {
	Id int
	// 到达时间
	ArriveTime int
	// 执行时间
	ExcuteTime int
	// 级别
	Level int
	// 开始时间
	StartTime int
	// 等待时间
	WaitTime int
	// 周转时间
	RoundTime int
	// 响应比
	Excellent float64
	// 结束时间
	OverTime int
}

func NewWork(id int, arriveTime int, excuteTime int, level int) *Work {
	return &Work{
		Id:         id,
		ArriveTime: arriveTime,
		ExcuteTime: excuteTime,
		Level:      level,
	}
}

func (w *Work) OutPrintln() {
	fmt.Printf("")
}

func ClientHRRF() {
	work1 := NewWork(1, 800, 50, 0)
	work2 := NewWork(2, 815, 30, 1)
	work3 := NewWork(3, 830, 25, 2)
	work4 := NewWork(4, 835, 20, 2)
	work5 := NewWork(5, 845, 15, 2)
	work6 := NewWork(6, 700, 10, 1)

	//work7 := NewWork(7, 820, 5, 0)
	collection1 := NewWorkCollection(work1, work2, work3, work4, work5, work6)
	algorithm1, err := collection1.CreateAlgorithm("HRRF")
	algorithm1.SetCollection(collection1)
	if err != nil {
		fmt.Println(err)
	} else {
		algorithm1.Schedule()
	}

}

func ClientFCFS() {
	work1 := NewWork(1, 800, 120, 1)
	work2 := NewWork(2, 850, 50, 1)
	work3 := NewWork(3, 900, 10, 1)
	work4 := NewWork(4, 950, 20, 1)

	//work7 := NewWork(7, 820, 5, 0)
	collection1 := NewWorkCollection(work1, work2, work3, work4)
	algorithm1, err := collection1.CreateAlgorithm("FCFS")
	algorithm1.SetCollection(collection1)
	if err != nil {
		fmt.Println(err)
	} else {
		algorithm1.Schedule()
	}

}
func ClientSJF() {
	work1 := NewWork(1, 800, 120, 0)
	work2 := NewWork(2, 850, 50, 1)
	work3 := NewWork(3, 900, 10, 2)
	work4 := NewWork(4, 950, 20, 2)

	//work7 := NewWork(7, 820, 5, 0)
	collection1 := NewWorkCollection(work1, work2, work3, work4)
	algorithm1, err := collection1.CreateAlgorithm("SJF")
	algorithm1.SetCollection(collection1)
	if err != nil {
		fmt.Println(err)
	} else {
		algorithm1.Schedule()
	}

}

func ClientHPF() {
	work1 := NewWork(1, 800, 120, 1)
	work2 := NewWork(2, 850, 50, 1)
	work3 := NewWork(3, 900, 10, 1)
	work4 := NewWork(4, 950, 20, 1)

	//work7 := NewWork(7, 820, 5, 0)
	collection1 := NewWorkCollection(work1, work2, work3, work4)
	algorithm1, err := collection1.CreateAlgorithm("HPF")
	algorithm1.SetCollection(collection1)
	if err != nil {
		fmt.Println(err)
	} else {
		algorithm1.Schedule()
	}
}
