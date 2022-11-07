package ProcessControlBlock

import (
	"fmt"
	"testing"

	"github.com/VariousImplementations/SecondarySchedue/pkg"
)

func TestMLFQSchedule(t *testing.T) {
	pool, err := NewPool(100, 7)
	if err != nil {
		fmt.Println("error:", err)
	}
	// 5个优先级
	// 每个优先级的时间片为 1 2 3 4 5
	multiQaueue := NewMultiLevelFeedbackQueue(pool, 5)
	multiQaueue.SetWorks(
		pkg.NewWork(1, "10:00", "40m", 5),
		pkg.NewWork(2, "10:20", "30m", 3),
		pkg.NewWork(3, "10:30", "50m", 4),
		pkg.NewWork(4, "10:50", "20m", 6),
	)
	multiQaueue.OutputQueueLevelAndAllottedTime()
	pkg.OutPutWorksArriveTime(multiQaueue.WorksNotInMemory)
	fmt.Println()
	multiQaueue.SetProgressTimeBar("10:00")
	fmt.Println()
	multiQaueue.Run()
	/*
		10:00~10:20 工作1 还剩20 在Queue4
		10:20~10:30 工作2 运行10min 还剩20min 在Queue4
		10:30~10:40 工作3 运行10min 还剩40min 在Queue4
		10:40~10:50 工作1 运行10min 还剩10min 在Queue4
		10:50~11:00 工作4 运行10min 还剩10min 在Queue4
		11:00~11:10 工作1 运行10min 还剩0 从Queue4删除
		11:10~11:30 工作2 运行20min 还剩0 从Queue4删除
		11:30~12:10 工作3 运行40min 还剩0 从Queue4删除
		12:10~12:20 工作4 运行10min 还剩0 从Queue4删除

	*/
}
