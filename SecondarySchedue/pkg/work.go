package pkg

import (
	"fmt"
	"time"
)

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
	Weights       float64
	ResponseRatio float64
}

func NewWork(id int, arriveTime string, excuteTime string, level int) *Work {
	arrive, erra := TimeFormat(arriveTime)
	excute, erre := TimeDurationFormat(excuteTime)
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
		ResponseRatio:        0.00000000,
	}
}

func OutPutWorksArriveTimeAndOverTime(w []*Work) {
	for _, v := range w {
		fmt.Printf("\nId:%-2d ArriveTime:%-8v ArriveMemoryTime:%-8v OverTime:%-8v RoundTime %-8v Weights %-8.4f\n", v.Id, v.ArriveTime, v.ArriveMemoryTime, v.OverTime, v.RoundTime, v.Weights)
	}
}

func OutPutWorksWaitTime(w []*Work) {
	for _, v := range w {
		fmt.Printf("\nId:%-2d ArriveTime:%-8v ArriveMemoryTime:%-8v OverTime:%-8v WaitTimeTime %-8v RoundTime %-8v Weights %-8.4f \n", v.Id, v.ArriveTime, v.ArriveMemoryTime, v.OverTime, v.WaitTime, v.RoundTime, v.Weights)
	}
}

func OutPutWorksArriveTime(w []*Work) {
	for _, v := range w {
		fmt.Printf("\nId:%-2d ArriveTime:%-8v  ExcuteTime:%-8v", v.Id, v.ArriveTime, v.ExcuteTime)
	}
	fmt.Println("")
}

func OutPutWorksArriveAndOverTime(w []*Work) {
	for _, v := range w {
		fmt.Printf("\nId:%-2d ArriveTime:%-8v OverTime:%-8v  ", v.Id, v.ArriveTime, v.OverTime)
	}
	fmt.Println()
}
