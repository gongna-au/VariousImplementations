package ProcessControlBlock

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var ErrPoolClosed = errors.New("Pool has been closed.")
var ErrInvalidPoolCap = errors.New("invalid pool cap")
var ErrPoolOverflow = errors.New("pool can not contain too many task")
var ErrQueueEmpty = errors.New("queue is empty")
var ErrPoolEmpty = errors.New("pool is empty")

const (
	RUNNING = 1
	STOPED  = 0
)

type PCBdPool struct {
	// 池的容量
	capacity uint64
	// 已经运行的PCB
	runningPCB uint64
	// 任务池的状态
	status int64
	// PCD队列池
	QueuesCreate     *PCBLinkedListChain
	QueuesReady      *PCBLinkedListChain
	QueuesRun        *PCBLinkedListChain
	QueuesTerminated *PCBLinkedListChain
	QueuesBlocked    *PCBLinkedListChain
	QueuesNum        uint64
}

func NewPool(cap uint64, num uint64) (*PCBdPool, error) {

	if cap <= 0 {
		return nil, ErrInvalidPoolCap
	}
	queuesCreate := &PCBLinkedListChain{}
	queuesReady := &PCBLinkedListChain{}
	queuesRun := &PCBLinkedListChain{}
	queuesTerminated := &PCBLinkedListChain{}
	queuesBlocked := &PCBLinkedListChain{}

	for i := 0; uint64(i) < num; i++ {
		queuesCreate.Push(NewProcessControlBlock().SetStatus("created"))
	}
	return &PCBdPool{
		capacity:         cap,
		status:           RUNNING,
		QueuesCreate:     queuesCreate,
		QueuesReady:      queuesReady,
		QueuesRun:        queuesRun,
		QueuesTerminated: queuesTerminated,
		QueuesBlocked:    queuesBlocked,
		QueuesNum:        uint64(num),
	}, nil
}

func (p *PCBdPool) incRunning() { // runningWorkers + 1
	atomic.AddUint64(&p.runningPCB, 1)
}

func (p *PCBdPool) decRunning() { // runningWorkers - 1
	atomic.AddUint64(&p.runningPCB, ^uint64(0))
}

func (p *PCBdPool) PutCreated(block *ProcessControlBlock) {
	if p.runningPCB < p.capacity {
		p.QueuesCreate.Push(block)
	} else {
		fmt.Println("error:", ErrPoolOverflow)
	}
}

func (p *PCBdPool) PopCreated() *ProcessControlBlock {
	if p.runningPCB < p.capacity {
		return p.QueuesCreate.Pop()
	} else {
		fmt.Println("error:", ErrQueueEmpty)
		return nil
	}
}

func (p *PCBdPool) PutRun(block *ProcessControlBlock) {
	if p.runningPCB < p.capacity {
		p.QueuesRun.Push(block)
	} else {
		fmt.Println("error:", ErrPoolOverflow)
	}
}

func (p *PCBdPool) PopRun() *ProcessControlBlock {
	if p.runningPCB < p.capacity {
		return p.QueuesRun.Pop()
	} else {
		fmt.Println("error:", ErrQueueEmpty)
		return nil
	}
}

func (p *PCBdPool) PutTerminated(block *ProcessControlBlock) {
	if p.runningPCB < p.capacity {
		p.QueuesTerminated.Push(block)
	} else {
		fmt.Println("error:", ErrPoolOverflow)
	}
}

func (p *PCBdPool) PopTerminated() *ProcessControlBlock {
	if p.runningPCB < p.capacity {
		return p.QueuesTerminated.Pop()
	} else {
		fmt.Println("error:", ErrQueueEmpty)
		return nil
	}
}

func (p *PCBdPool) PutReady(block *ProcessControlBlock) {
	if p.runningPCB < p.capacity {
		p.QueuesReady.Push(block)
	} else {
		fmt.Println("error:", ErrPoolOverflow)
	}
}

func (p *PCBdPool) PopReady() *ProcessControlBlock {
	if p.runningPCB < p.capacity {
		return p.QueuesReady.Pop()
	} else {
		fmt.Println("error:", ErrQueueEmpty)
		return nil
	}
}

func (p *PCBdPool) PutBlock(block *ProcessControlBlock) {
	if p.runningPCB < p.capacity {
		p.QueuesBlocked.Push(block)
	} else {
		fmt.Println("error:", ErrPoolOverflow)
	}
}

func (p *PCBdPool) PopBlock() *ProcessControlBlock {
	if p.runningPCB < p.capacity {
		return p.QueuesBlocked.Pop()
	} else {
		fmt.Println("error:", ErrQueueEmpty)
		return nil
	}

}

func (p *PCBdPool) TraverseCreated() {
	current := p.QueuesCreate.HeadBlock
	fmt.Println("process in created queue")
	for {
		if current != nil {
			fmt.Println("Id:", current.Id, "Status:", current.Status)
			current = current.NextBlock
		} else {
			break
		}
	}
	fmt.Println()
}

func (p *PCBdPool) TraverseReady() {
	current := p.QueuesReady.HeadBlock
	fmt.Println("process in ready queue")
	for {
		if current != nil {
			fmt.Println("Id:", current.Id, "Status:", current.Status)
			current = current.NextBlock
		} else {
			break
		}
	}
	fmt.Println()
}

func (p *PCBdPool) TraverseRun() {
	current := p.QueuesRun.HeadBlock
	fmt.Println("process in run queue")
	for {
		if current != nil {
			fmt.Println("Id:", current.Id, "Status:", current.Status)
			current = current.NextBlock
		} else {
			break
		}
	}
	fmt.Println()
}

func (p *PCBdPool) TraverseTerminated() {
	current := p.QueuesTerminated.HeadBlock
	fmt.Println("process in terminated queue")
	for {
		if current != nil {
			fmt.Println("Id:", current.Id, "Status:", current.Status)
			current = current.NextBlock
		} else {
			break
		}
	}
	fmt.Println()
}

func (p *PCBdPool) TraverseBlock() {
	current := p.QueuesBlocked.HeadBlock
	fmt.Println("process in block queue")
	for {
		if current != nil {
			fmt.Println("Id:", current.Id, "Status:", current.Status)
			current = current.NextBlock
		} else {
			break
		}
	}
	fmt.Println()
}
