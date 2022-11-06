package ProcessControlBlock

import (
	"fmt"
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
	ProgressTime time.Time
}

func NewMultiLevelFeedbackQueue(p *PCBdPool) *MultiLevelFeedbackQueue {
	return &MultiLevelFeedbackQueue{
		Pool: p,
		PCBChains: []*PCBLinkedListChain{
			// 数字越小优先级越高
			NewPCBLinkedListChain(0),
			NewPCBLinkedListChain(1),
			NewPCBLinkedListChain(2),
			NewPCBLinkedListChain(3),
			NewPCBLinkedListChain(4),
			NewPCBLinkedListChain(5),
		},
		WorksNotInMemory: []*pkg.Work{},
	}
}

func (m *MultiLevelFeedbackQueue) SetWorks(w ...*pkg.Work) {
	m.WorksNotInMemory = append(m.WorksNotInMemory, w...)
}

func (m *MultiLevelFeedbackQueue) Schedule() {
	if len(m.WorksNotInMemory)<=0{
		return
	}
	// 得到一个进程
	process := m.Pool.QueuesCreate.Pop()
	firstWork := m.WorksNotInMemory[0]
	process.SetWork(firstWork)
	// 所有先来的任务都会被先放在优先级最高的队列
	m.PCBChains[0].Push(process)
	if len(m.WorksNotInMemory)>1{
		m.WorksNotInMemory = m.WorksNotInMemory[1:]
	}
	
}

func (m *MultiLevelFeedbackQueue) WorkRun() {
	for _,v := m.PCBChains{
		if !v.IsEmptyQueue(){
			
		}
	}
	
}

func (m *MultiLevelFeedbackQueue) Run() {
	for _,v := m.PCBChains{
		if !v.IsEmptyQueue(){
			
		}
	}
	
}

