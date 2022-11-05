package SecondarySchedue

import (
	"errors"
	"fmt"
	"github.com/VariousImplementations/SecondarySchedue/pkg"
	"sort"
	"time"
)

// 两道批处理系统中基于优先数的调度算法计算结果
// 作业调度采用最短作业优先 进程调度采用优先数优先的抢占式调度算法

type HRRF struct {
	SJFStartTime     time.Time
	Works            []*pkg.Work
	QueueNotInMemory []*pkg.Work
	QueueInMemory    []*pkg.Work
	QueueHasFinish   []*pkg.Work
	// 进度时间用来记录系统现在运行的时间点
	ProgressTime time.Time
}

func NewHRRF(w ...*pkg.Work) *HRRF {
	var arrive time.Time
	if len(w) == 0 {
		arrive = time.Now()
	} else {
		arrive = w[0].ArriveTime
		w[0].ArriveMemoryTime = w[0].ArriveTime
	}

	return &HRRF{
		SJFStartTime:     arrive,
		Works:            append([]*pkg.Work{}, w...),
		QueueNotInMemory: append([]*pkg.Work{}, w[1:]...),
		QueueInMemory:    append([]*pkg.Work{}, w[0]),
		QueueHasFinish:   []*pkg.Work{},
		ProgressTime:     arrive,
	}
}

func (s *HRRF) DeleteWorkNotInMemoryById(id int) {
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

// 在内存中没有作业的时候调用
func (s *HRRF) JudgeWorkhasCome() (*pkg.Work, bool) {
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

	if len(hasArrrice) > 1 {
		// 在已经到达的任务下选择
		sort.Slice(hasArrrice, func(i, j int) bool {
			if hasArrrice[i].Level < hasArrrice[j].Level {
				return true
			} else {
				return false
			}
		})
		temp := hasArrrice[0]
		// 然后遍历删除该任务
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
func (s *HRRF) WorkInMemoryExcuteByLevel() {
	sort.Slice(s.QueueInMemory, func(i, j int) bool {
		if s.QueueInMemory[i].Level < s.QueueInMemory[j].Level {
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
		s.ProgressTime = pkg.Add(s.ProgressTime, "5m")

	}

}

// 只有一个任务在内存中
func (s *HRRF) WorkInMemoryExcute() {
	if pkg.RemainTimeLessThanZero(s.QueueInMemory[0].RemainingExecuteTime) == true {
		// 放入完成的队列// 放入完成的队列
		s.QueueInMemory[0].OverTime = s.ProgressTime
		s.QueueHasFinish = append(s.QueueHasFinish, s.QueueInMemory[0])
		s.QueueInMemory = []*pkg.Work{}
		return
	}
	s.QueueInMemory[0].RemainingExecuteTime = pkg.Sub(s.QueueInMemory[0].RemainingExecuteTime, "5m")
	s.ProgressTime = pkg.Add(s.ProgressTime, "5m")
}

func (s *HRRF) UpdateProgressTime() {
	s.ProgressTime = pkg.Add(s.ProgressTime, "5m")
}

func (s *HRRF) Schedue() {
	for i := 0; len(s.QueueHasFinish) < 4; i++ {
		if len(s.QueueInMemory) == 0 {
			return
		} else if len(s.QueueInMemory) == 1 {
			// 调度进入一个作业 然后根据优先级别判断后执行
			w, isCome := s.JudgeWorkhasCome()

			if isCome {
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

func (s *HRRF) InitRoundTime() {
	for _, v := range s.QueueHasFinish {
		v.RoundTime = pkg.TimeSub(v.OverTime, v.ArriveTime)

	}
}

func (s *HRRF) InitWeights() {
	for _, v := range s.QueueHasFinish {
		var float64_Minutes1 float64 = v.RoundTime.Minutes()
		var float64_Minutes2 float64 = v.ExcuteTime.Minutes()
		v.Weights = float64_Minutes1 / float64_Minutes2
	}
}

// 计算平均周转时间  完成时间-开始时间
func (s *HRRF) GetAverageTurnaRoundTime() time.Duration {
	var d time.Duration
	for _, v := range s.QueueHasFinish {
		d = pkg.DurationAdd(d, v.RoundTime)
	}
	return d
}
func (s *HRRF) GetAverageTurnaRoundTimeByWeight() float64 {
	var d float64
	for _, v := range s.QueueHasFinish {
		d = d + v.Weights
	}
	return d
}

func HRRFClient() {
	s := NewHRRF(
		pkg.NewWork(1, "10:00", "40m", 5),
		pkg.NewWork(2, "10:05", "30m", 3),
		pkg.NewWork(3, "10:30", "50m", 4),
		pkg.NewWork(4, "10:50", "20m", 6),
	)
	/*
		10:00


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
