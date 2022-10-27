package SecondarySchedue

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/VariousImplementations/SecondarySchedue/pkg"
)

// 具体的work
type Work struct {
	Id int
	// 到达时间
	ArriveTime time.Time
	// 进入内存的时间
	ArriveMemoryTime time.Time
	// 第一次开始执行的时间
	StartTime time.Time
	// 上一次执行结束的时间 , 也就是被暂停的时间点
	StopTime time.Time
	// 结束时间
	OverTime time.Time
	// 需要的执行时间
	ExcuteTime time.Duration
	// 剩余执行时间
	RemainingExecuteTime time.Duration
	// 等待时间
	WaitTime time.Duration
	// 周转时间
	RoundTime time.Duration
	// 级别
	Level int
	// 响应比
	Excellent float64
	//权重
	Weights float64
}

func NewWork(id int, arriveTime string, excuteTime string, level int) *Work {
	arrive, erra := pkg.TimeFormat(arriveTime)
	excute, erre := pkg.TimeDurationFormat(excuteTime)
	if erra != nil || erre != nil {
		return nil
	}

	return &Work{
		Id:                   id,
		ArriveTime:           arrive,
		ExcuteTime:           excute,
		Level:                level,
		RemainingExecuteTime: excute,
		Weights:              0.00000000,
	}
}

func OutPutWorksArriveTimeAndOverTime(w []*Work) {
	for _, v := range w {
		fmt.Printf("Id:%-2d ArriveTime:%-8v ArriveMemoryTime:%-8v OverTime:%-8v RoundTime %-8v Weights %-8.4f\n", v.Id, v.ArriveTime, v.ArriveMemoryTime, v.OverTime, v.RoundTime, v.Weights)
	}
}

type SJF struct {
	SJFStartTime     time.Time
	Works            []*Work
	QueueNotInMemory []*Work
	QueueInMemory    []*Work
	QueueHasFinish   []*Work
	// 进度时间用来记录系统现在运行的时间点
	ProgressTime time.Time
}

func NewSJF(w ...*Work) *SJF {
	var arrive time.Time
	if len(w) == 0 {
		arrive = time.Now()
	} else {
		arrive = w[0].ArriveTime
		w[0].ArriveMemoryTime = w[0].ArriveTime
	}

	return &SJF{
		SJFStartTime:     arrive,
		Works:            append([]*Work{}, w...),
		QueueNotInMemory: append([]*Work{}, w[1:]...),
		QueueInMemory:    append([]*Work{}, w[0]),
		QueueHasFinish:   []*Work{},
		ProgressTime:     arrive,
	}
}

// 在内存中没有作业的时候调用
func (s *SJF) JudgeWorkhasCome() (*Work, bool) {
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
	hasArrrice := []*Work{}
	// 可以先排序再找
	for _, v := range s.QueueNotInMemory {
		if pkg.TimeCompare(v.ArriveTime, s.ProgressTime) == true {
			hasArrrice = append(hasArrrice, v)
		}
	}

	if len(hasArrrice) > 1 {
		sort.Slice(s.QueueNotInMemory, func(i, j int) bool {
			if int(s.QueueNotInMemory[i].RemainingExecuteTime) < int(s.QueueNotInMemory[j].RemainingExecuteTime) {
				return true
			} else {
				return false
			}
		})
		temp := s.QueueNotInMemory[0]
		s.QueueNotInMemory = s.QueueNotInMemory[1:]
		return temp, true
	} else if len(hasArrrice) == 1 {

		if len(s.QueueNotInMemory) > 1 {
			temp := s.QueueNotInMemory[0]
			s.QueueNotInMemory = s.QueueNotInMemory[1:]
			return temp, true
		} else {
			temp := s.QueueNotInMemory[0]
			s.QueueNotInMemory = []*Work{}
			return temp, true
		}

	} else {
		return nil, false
	}

}

// 现在有两个任务
func (s *SJF) WorkInMemoryExcuteByLevel() {
	sort.Slice(s.QueueInMemory, func(i, j int) bool {
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
		s.ProgressTime = pkg.Add(s.ProgressTime, "5m")

	}

}

// 只有一个任务在内存中
func (s *SJF) WorkInMemoryExcute() {
	if pkg.RemainTimeLessThanZero(s.QueueInMemory[0].RemainingExecuteTime) == true {
		// 放入完成的队列// 放入完成的队列
		s.QueueInMemory[0].OverTime = s.ProgressTime
		s.QueueHasFinish = append(s.QueueHasFinish, s.QueueInMemory[0])
		s.QueueInMemory = []*Work{}
		return
	}
	s.QueueInMemory[0].RemainingExecuteTime = pkg.Sub(s.QueueInMemory[0].RemainingExecuteTime, "5m")
	s.ProgressTime = pkg.Add(s.ProgressTime, "5m")
}

func (s *SJF) UpdateProgressTime() {
	s.ProgressTime = pkg.Add(s.ProgressTime, "5m")
}

func (s *SJF) Schedue() {
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

func (s *SJF) InitRoundTime() {
	for _, v := range s.QueueHasFinish {
		v.RoundTime = pkg.TimeSub(v.OverTime, v.ArriveTime)

	}
}

func (s *SJF) InitWeights() {
	for _, v := range s.QueueHasFinish {
		var float64_Minutes1 float64 = v.RoundTime.Minutes()
		var float64_Minutes2 float64 = v.ExcuteTime.Minutes()
		v.Weights = float64_Minutes1 / float64_Minutes2
	}
}

// 计算平均周转时间  完成时间-开始时间
func (s *SJF) GetAverageTurnaRoundTime() time.Duration {
	var d time.Duration
	for _, v := range s.QueueHasFinish {
		d = pkg.DurationAdd(d, v.RoundTime)
	}
	return d
}
func (s *SJF) GetAverageTurnaRoundTimeByWeight() float64 {
	var d float64
	for _, v := range s.QueueHasFinish {
		d = d + v.Weights
	}
	return d
}

func Client() {
	s := NewSJF(
		NewWork(1, "10:00", "30m", 5),
		NewWork(2, "10:05", "20m", 3),
		NewWork(3, "10:10", "5m", 4),
		NewWork(4, "10:20", "10m", 6),
	)
	s.Schedue()
	s.InitRoundTime()
	s.InitWeights()
	OutPutWorksArriveTimeAndOverTime(s.QueueHasFinish)
	r := s.GetAverageTurnaRoundTime() / 4
	fmt.Printf("Round Time is %v \n", r)
	t := s.GetAverageTurnaRoundTimeByWeight() / 4
	fmt.Printf("Round Time By Weights is %v \n", t)

}
