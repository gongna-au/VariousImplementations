package ProcessControlBlock

import (
	"fmt"
	"strconv"
	"time"

	"github.com/VariousImplementations/SecondarySchedue/pkg"
)

/*
队列数量。
每个队列的调度算法。
用以确定何时升级到更高优先级队列的方法。
用以确定何时降级到更低优先级队列的方法。
用以确定进程在需要服务时将会进入哪个队列的方法。
*/
type MultiLevelFeedbackQueue struct {
	Pool             *PCBdPool
	PCBChains        []*PCBLinkedListChain
	WorksNotInMemory []*pkg.Work
	// 进度条 不是进程时间
	ProgressTimeBar     time.Time
	ProcessHasFinishsed []*ProcessControlBlock
	WorkNum             int
}

// 默认设置了5个队列
func NewMultiLevelFeedbackQueue(p *PCBdPool, queueLevel int) *MultiLevelFeedbackQueue {
	// 数字越小优先级越高
	chains := []*PCBLinkedListChain{}
	for i := 0; i < queueLevel; i++ {
		chains = append(chains, NewPCBLinkedListChain(i))
	}
	// 默认被分配的时间是优先级的+1
	for k, v := range chains {
		v.AllottedTime = strconv.Itoa(k+1) + "m"
	}
	return &MultiLevelFeedbackQueue{
		Pool:                p,
		PCBChains:           chains,
		WorksNotInMemory:    []*pkg.Work{},
		ProcessHasFinishsed: []*ProcessControlBlock{},
	}
}

func (m *MultiLevelFeedbackQueue) SetWorks(w ...*pkg.Work) {
	m.WorksNotInMemory = append(m.WorksNotInMemory, w...)
	m.WorkNum = len(w)
}

func (m *MultiLevelFeedbackQueue) SetAllottedTime(times ...string) {
	for k, v := range m.PCBChains {
		v.AllottedTime = times[k]
	}
}

func (m *MultiLevelFeedbackQueue) SetProgressTimeBar(time string) {
	m.ProgressTimeBar, _ = pkg.TimeFormat(time)
	fmt.Println("ProgressTimeBar init:", m.ProgressTimeBar)
}

func (m *MultiLevelFeedbackQueue) Run() {
	for {
		if len(m.ProcessHasFinishsed) == m.WorkNum {
			break
		}
		m.Schedule()
		m.WorkRun()
		m.Traverse()
		fmt.Println("Works has scheduled is:")
		pkg.OutPutWorksArriveTime(m.WorksNotInMemory)
		fmt.Println("Works has finished are:")
		m.TraverseFinishedProcess()

	}
}

func (m *MultiLevelFeedbackQueue) Schedule() {
	fmt.Println("---Before Schedule Time:", m.ProgressTimeBar)
	if len(m.WorksNotInMemory) <= 0 {
		return
	}
	// 进度条小于等于后面
	if pkg.TimeCompare(m.WorksNotInMemory[0].ArriveTime, m.ProgressTimeBar) {
		// 得到一个进程
		process := m.Pool.QueuesCreate.Pop()
		firstWork := m.WorksNotInMemory[0]
		process.SetWork(firstWork)
		// 所有先来的任务都会被先放在优先级最高的队列
		fmt.Println("Id:", process.Work.Id, " is push in Queue", m.PCBChains[0].Priority)
		m.PCBChains[0].Push(process)
		if len(m.WorksNotInMemory) > 1 {
			m.WorksNotInMemory = m.WorksNotInMemory[1:]
		} else {
			m.WorksNotInMemory = []*pkg.Work{}
		}
	}

}

func (m *MultiLevelFeedbackQueue) WorkRun() {
	for k, v := range m.PCBChains {
		if !v.IsEmptyQueue() {
			//fmt.Println("k", k, "level queue not empty")
			// 标志着是最后一级的队列
			if k == len(m.PCBChains)-1 {
				processs, _ := v.GetHeadProcessControlElem()
				if pkg.RemainTimeLessThanZero(processs.Work.RemainingExecuteTime) == true {
					// 放入完成的队列
					processs.Work.OverTime = m.ProgressTimeBar
					processs.Status = "terminated"
					m.ProcessHasFinishsed = append(m.ProcessHasFinishsed, processs)
					v.Pop()
				} else {
					// 执行时间,没有下一个队列了
					processs.Work.RemainingExecuteTime = pkg.Sub(processs.Work.RemainingExecuteTime, v.AllottedTime)
					m.ProgressTimeBar = pkg.Add(m.ProgressTimeBar, v.AllottedTime)
				}
				return
			}

			processs := v.Pop()
			if pkg.RemainTimeLessThanZero(processs.Work.RemainingExecuteTime) == true {
				// 放入完成的队列
				processs.Work.OverTime = m.ProgressTimeBar
				processs.Status = "terminated"
				m.ProcessHasFinishsed = append(m.ProcessHasFinishsed, processs)
			} else {
				// 执行时间放入下一个队列
				processs.Work.RemainingExecuteTime = pkg.Sub(processs.Work.RemainingExecuteTime, v.AllottedTime)
				//fmt.Printf("%d id be excutes 5m\n", s.QueueInMemory[0].Id)
				m.ProgressTimeBar = pkg.Add(m.ProgressTimeBar, v.AllottedTime)
				if (k + 1) <= len(m.PCBChains)-1 {
					m.PCBChains[k+1].Push(processs)
				}
			}
			break
		}
	}
}

func (m *MultiLevelFeedbackQueue) Traverse() {
	fmt.Println("---After Excute Time is:", m.ProgressTimeBar)
	fmt.Println("After Excute Queue is:")
	for k, v := range m.PCBChains {
		fmt.Println("Level", k, "Queue:")
		if !v.IsEmptyQueue() {
			v.Traverse()
		} else {
			fmt.Println("Null Queue")
		}
	}
}

func (m *MultiLevelFeedbackQueue) TraverseFinishedProcess() {
	for _, v := range m.ProcessHasFinishsed {
		fmt.Println("WorkId:", v.Work.Id, "ProcessIs:", v.Id)
	}
	fmt.Println()
}

func (m *MultiLevelFeedbackQueue) OutputQueueLevelAndAllottedTime() {
	for _, v := range m.PCBChains {
		fmt.Println("QueueLevel:", v.Priority, "AllottedTime:", v.AllottedTime)
	}
	fmt.Println()
}
