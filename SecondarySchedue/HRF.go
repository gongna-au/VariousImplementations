package SecondarySchedue

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/VariousImplementations/SecondarySchedue/pkg"
)

// 两道批处理系统中基于优先数的调度算法计算结果
// 作业调度采用最短作业优先 进程调度采用高响应比优先的抢占式调度算法
type HRF struct {
	SJFStartTime     time.Time
	Works            []*pkg.Work
	QueueNotInMemory []*pkg.Work
	QueueInMemory    []*pkg.Work
	QueueHasFinish   []*pkg.Work
	// 进度时间用来记录系统现在运行的时间点
	ProgressTime time.Time
}

func NewHRF(w ...*pkg.Work) *HRF {
	var arrive time.Time
	if len(w) == 0 {
		arrive = time.Now()
	} else {
		arrive = w[0].ArriveTime
		w[0].ArriveMemoryTime = w[0].ArriveTime
	}

	return &HRF{
		SJFStartTime:     arrive,
		Works:            append([]*pkg.Work{}, w...),
		QueueNotInMemory: append([]*pkg.Work{}, w[1:]...),
		QueueInMemory:    append([]*pkg.Work{}, w[0]),
		QueueHasFinish:   []*pkg.Work{},
		ProgressTime:     arrive,
	}
}

func (s *HRF) MarkHasArrivedWorkWaitTime(w []*pkg.Work) {
	for _, v := range w {
		v.WaitTime = pkg.TimeDurationAdd(v.WaitTime, "5m")
	}
}

func (s *HRF) DeleteWorkNotInMemoryById(id int) {
	var index int
	for k, v := range s.QueueNotInMemory {
		if v.Id == id {
			index = k
		}
	}
	if index == len(s.QueueNotInMemory)-1 {
		s.QueueNotInMemory = s.QueueNotInMemory[:index]
	} else {
		s.QueueNotInMemory = append(s.QueueNotInMemory[:index], s.QueueNotInMemory[index+1:]...)
	}
}

func (s *HRF) CalculateResponseRatio(w []*pkg.Work) {
	for _, v := range w {
		var minutes1 = v.WaitTime.Minutes()
		var minutes2 = v.ExcuteTime.Minutes()
		v.Excellent = float64(minutes1+minutes2) / float64(minutes2)
	}
}

// 对于后续队列中已经到来的任务增加等待时间
func (s *HRF) UpdateWaitTime() {
	// 判断有没有作业到来
	hasArrrice := []*pkg.Work{}
	// 可以先排序再找
	for _, v := range s.QueueNotInMemory {
		if pkg.TimeCompare(v.ArriveTime, s.ProgressTime) == true {
			hasArrrice = append(hasArrrice, v)
		}
	}
	for _, v := range hasArrrice {
		v.WaitTime = pkg.TimeDurationAdd(v.WaitTime, "5m")
	}
}

// 在内存中没有作业的时候调用
func (s *HRF) JudgeWorkhasCome() (*pkg.Work, bool) {
	// 先排序
	if len(s.QueueNotInMemory) > 1 {
		sort.Slice(s.QueueNotInMemory, func(i, j int) bool {
			if pkg.TimeCompare(s.QueueNotInMemory[i].ArriveTime, s.QueueNotInMemory[j].ArriveTime) {
				return true
			} else {
				return false
			}
		})
	}
	// 判断后续有没有作业到来
	hasArrrice := []*pkg.Work{}
	// 可以先排序再找
	for _, v := range s.QueueNotInMemory {
		if pkg.TimeCompare(v.ArriveTime, s.ProgressTime) == true {
			hasArrrice = append(hasArrrice, v)
		}
	}
	// 对于已经到达的任务，却没有被调度进入内存的任务需要等待
	if len(hasArrrice) > 1 {
		// 先计算出每个任务的响应比
		s.CalculateResponseRatio(hasArrrice)
		// 然后排序
		sort.Slice(hasArrrice, func(i, j int) bool {
			if hasArrrice[i].Excellent > hasArrrice[j].Excellent {
				return true
			} else {
				return false
			}
		})
		temp := hasArrrice[0]
		// 本次循环也只是调度，并没有实际运行任何一个任务
		//s.MarkHasArrivedWorkWaitTime(hasArrrice[1:])
		s.DeleteWorkNotInMemoryById(temp.Id)
		return temp, true
	} else if len(hasArrrice) == 1 {
		if len(s.QueueNotInMemory) > 1 {
			temp := hasArrrice[0]
			s.DeleteWorkNotInMemoryById(temp.Id)
			return temp, true
		} else {
			temp := s.QueueNotInMemory[0]
			s.QueueNotInMemory = []*pkg.Work{}
			return temp, true
		}

	} else {
		return nil, false
	}

}

// 现在有两个任务
func (s *HRF) WorkInMemoryExcuteByLevel() {
	sort.Slice(s.QueueInMemory, func(i, j int) bool {
		// 进程调度都是基于--最短执行时间的可抢占式调度册罗
		if int(s.QueueInMemory[i].RemainingExecuteTime) < int(s.QueueInMemory[j].RemainingExecuteTime) {
			return true
		} else {
			return false
		}
	})

	if pkg.RemainTimeLessThanZero(s.QueueInMemory[0].RemainingExecuteTime) == true {
		// 放入完成的队列
		s.QueueInMemory[0].OverTime = s.ProgressTime
		s.QueueHasFinish = append(s.QueueHasFinish, s.QueueInMemory[0])
		s.QueueInMemory = s.QueueInMemory[1:]

	} else {
		s.QueueInMemory[0].RemainingExecuteTime = pkg.Sub(s.QueueInMemory[0].RemainingExecuteTime, "5m")
		// 每次拖动进度条，已经到来的任务就一定会等待
		s.UpdateWaitTime()
		s.ProgressTime = pkg.Add(s.ProgressTime, "5m")

	}

}

// 只有一个任务在内存中
func (s *HRF) WorkInMemoryExcute() {
	if pkg.RemainTimeLessThanZero(s.QueueInMemory[0].RemainingExecuteTime) == true {
		// 放入完成的队列// 放入完成的队列
		s.QueueInMemory[0].OverTime = s.ProgressTime
		s.QueueHasFinish = append(s.QueueHasFinish, s.QueueInMemory[0])
		s.QueueInMemory = []*pkg.Work{}
		return
	}
	s.QueueInMemory[0].RemainingExecuteTime = pkg.Sub(s.QueueInMemory[0].RemainingExecuteTime, "5m")

	s.UpdateWaitTime()
	s.ProgressTime = pkg.Add(s.ProgressTime, "5m")
}

func (s *HRF) UpdateProgressTime() {
	// 只要是更新进度条的地方，就需要给已经到来的任务增加等待时间
	s.UpdateWaitTime()
	s.ProgressTime = pkg.Add(s.ProgressTime, "5m")
}

func (s *HRF) Schedue() {
	for i := 0; len(s.QueueHasFinish) < 4; i++ {
		if len(s.QueueInMemory) == 0 {
			return
		} else if len(s.QueueInMemory) == 1 {
			// 调度进入一个作业 然后根据优先级别判断后执行
			w, isCome := s.JudgeWorkhasCome()
			if isCome {
				// 意味着有任务
				w.ArriveMemoryTime = s.ProgressTime
				s.QueueInMemory = append(s.QueueInMemory, w)
			} else {
				s.WorkInMemoryExcute()
			}
		} else if len(s.QueueInMemory) == 2 {
			s.WorkInMemoryExcuteByLevel()
		} else {
			fmt.Println(errors.New("Work Queue In Memory length can not more than 2!"))
		}
	}

	/* fmt.Println("Queue in memory")
	OutPutWorksArriveTimeAndOverTime(s.QueueInMemory)
	fmt.Println("Queue not in memory")
	OutPutWorksArriveTimeAndOverTime(s.QueueNotInMemory) */
}

func (s *HRF) InitRoundTime() {
	for _, v := range s.QueueHasFinish {
		v.RoundTime = pkg.TimeSub(v.OverTime, v.ArriveTime)

	}
}

func (s *HRF) InitWeights() {
	for _, v := range s.QueueHasFinish {
		var float64_Minutes1 float64 = v.RoundTime.Minutes()
		var float64_Minutes2 float64 = v.ExcuteTime.Minutes()
		v.Weights = float64_Minutes1 / float64_Minutes2
	}
}

// 计算平均周转时间  完成时间-开始时间
func (s *HRF) GetAverageTurnaRoundTime() time.Duration {
	var d time.Duration
	for _, v := range s.QueueHasFinish {
		d = pkg.DurationAdd(d, v.RoundTime)
	}
	return d
}
func (s *HRF) GetAverageTurnaRoundTimeByWeight() float64 {
	var d float64
	for _, v := range s.QueueHasFinish {
		d = d + v.Weights
	}
	return d
}

func HRFClient() {
	s := NewSJF(
		pkg.NewWork(1, "10:00", "30m", 5),
		pkg.NewWork(2, "10:05", "20m", 3),
		pkg.NewWork(3, "10:10", "5m", 4),
		pkg.NewWork(4, "10:20", "10m", 6),
	)
	/*
		1: 10:00～10:05 剩余25
		2: 10:05～10：25 2被调度进入内存 并且2执行结束 3等待15m 4 等待5m  3的响应比例15+5)/5=4   4的响应比例 (5+10)/10=1.5
		3：10:25~10:30  3被调度进入内存 并且3执行结束
		4：10:30~10:40  4被调度进入内存 并且4执行结束
		5：10:40~11:05  1被调度进入内存 并且1执行结束
	*/

	fmt.Println("Input information:")
	pkg.OutPutWorksArriveTimeAndOverTime(s.QueueInMemory)
	pkg.OutPutWorksArriveTimeAndOverTime(s.QueueNotInMemory)
	s.Schedue()
	s.InitRoundTime()
	s.InitWeights()
	fmt.Println("Information after schedued:")
	pkg.OutPutWorksArriveTimeAndOverTime(s.QueueHasFinish)
	r := s.GetAverageTurnaRoundTime() / 4
	fmt.Printf("Round Time is %v \n", r)
	t := s.GetAverageTurnaRoundTimeByWeight() / 4
	fmt.Printf("Round Time By Weights is %v \n", t)
}
