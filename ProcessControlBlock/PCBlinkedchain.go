package ProcessControlBlock

import (
	"errors"
	"fmt"
	"strings"

	"github.com/VariousImplementations/SecondarySchedue/pkg"
	"github.com/google/uuid"
)

type Status string

// created ready  running  waiting  terminated
var ProcessStatus = []Status{
	"created",
	"ready",
	"running",
	"block",
	"terminated",
}

func NewStatus(str string) Status {
	if strings.Compare(str, "created") == 0 {
		return ProcessStatus[0]
	}
	if strings.Compare(str, "ready") == 0 {
		return ProcessStatus[1]
	}
	if strings.Compare(str, "running") == 0 {
		return ProcessStatus[2]
	}
	if strings.Compare(str, "block") == 0 {
		return ProcessStatus[3]
	}
	if strings.Compare(str, "terminated") == 0 {
		return ProcessStatus[4]
	}

	return Status("")
}

type ProcessControlBlock struct {
	// 进程标志符号
	Id string
	// 用户标志符
	UserId string

	Status Status
	// 资源
	Resource string
	// CPU的状态
	CPUStatus string
	// 优先级
	Priority int
	// 上一个进程块
	LastBlock *ProcessControlBlock
	// 下一个进程块
	NextBlock *ProcessControlBlock
	//抽象的任务
	Work *pkg.Work
}

func NewProcessControlBlock() *ProcessControlBlock {
	return &ProcessControlBlock{
		Id: uuid.New().String(),
	}
}

func (p *ProcessControlBlock) SetUserId(userId string) *ProcessControlBlock {
	p.UserId = userId
	return p
}

// created ready  running  waiting  terminated
func (p *ProcessControlBlock) SetResource(resource string) *ProcessControlBlock {
	p.Resource = resource
	return p
}

// created ready  running  waiting  terminated
func (p *ProcessControlBlock) SetCPU(cpu string) *ProcessControlBlock {
	p.CPUStatus = cpu
	return p
}

// created ready  running  waiting  terminated
func (p *ProcessControlBlock) SetStatus(status Status) *ProcessControlBlock {
	p.Status = status
	return p
}

func (p *ProcessControlBlock) SetPriority(priority int) *ProcessControlBlock {
	p.Priority = priority
	return p
}

func (p *ProcessControlBlock) SetWork(work *pkg.Work) *ProcessControlBlock {
	p.Work = work
	return p
}

func (p *ProcessControlBlock) WorkRun() {
	if p.Work != nil {

	}
}

func (p *ProcessControlBlock) SetNext(next *ProcessControlBlock) {
	if p.NextBlock != nil {
		left := p
		right := p.NextBlock
		mid := next
		left.NextBlock = mid
		mid.NextBlock = right
		right.LastBlock = mid
		mid.LastBlock = left
	} else {
		p.NextBlock = next
		next.LastBlock = p
	}
}

func (p *ProcessControlBlock) SetLast(last *ProcessControlBlock) {
	if p.LastBlock != nil {
		left := p.LastBlock
		mid := last
		right := p
		right.LastBlock = mid
		mid.LastBlock = left
		left.NextBlock = mid
		mid.NextBlock = right
	} else {
		fmt.Println("Error:", errors.New("The queue header not allow to insert elem before it"))
	}
}

func (p *ProcessControlBlock) PriorityReduce() {
	if p.Priority != 0 {
		p.Priority = p.Priority - 1
	}
}

// 进程阻塞
func (p *ProcessControlBlock) SetBlock() {
	if p != nil {
		p.Status = NewStatus("block")
	}
}

// 进程唤醒
func (p *ProcessControlBlock) SetReady() {
	if p != nil {
		p.Status = NewStatus("ready")
	}
}

// 进程结束
func (p *ProcessControlBlock) SetTreminate() {
	if p != nil {
		p.Status = NewStatus("terminated")
	}
}

// 进程运行
func (p *ProcessControlBlock) SetRunning() {
	if p != nil {
		p.Status = NewStatus("running")
	}
}

// 进程撤销删除
func (p *ProcessControlBlock) Delte() {
	if p != nil {
		p = nil
	}
}

type PCBLinkedListChain struct {

	// 优先级
	Priority int
	// 保存一个队列头部的数据
	HeadBlock *ProcessControlBlock
	// 保存队列尾部的数据
	TailBlock *ProcessControlBlock
	// 被分配的时间份额
	AllottedTime string
}

func NewPCBLinkedListChain(pri int) *PCBLinkedListChain {
	return &PCBLinkedListChain{
		Priority:  pri,
		HeadBlock: nil,
		TailBlock: nil,
	}
}

func (p *PCBLinkedListChain) SetPCBLinkedListChainAllottedTime(time string) *PCBLinkedListChain {
	p.AllottedTime = time
	return p
}

func (p *PCBLinkedListChain) SetAllottedTime(time string) {
	p.AllottedTime = time
}

func (q *PCBLinkedListChain) GetLastProcessControlElem() (*ProcessControlBlock, error) {
	if q.TailBlock != nil {
		return q.TailBlock, nil
	}
	return nil, errors.New("Queue has not any process control block")
}

func (q *PCBLinkedListChain) GetHeadProcessControlElem() (*ProcessControlBlock, error) {
	if q.HeadBlock != nil {
		return q.HeadBlock, nil
	}
	return nil, errors.New("Queue has not any process control block")
}

func (q *PCBLinkedListChain) IsEmptyQueue() bool {
	currency := q.HeadBlock
	if currency != nil {
		return false
	}
	return true
}

func (q *PCBLinkedListChain) Push(next *ProcessControlBlock) error {
	currency := q.HeadBlock
	tail := q.TailBlock
	if currency != nil {
		q.TailBlock = next
		tail.NextBlock = next
		next.LastBlock = tail
		return nil
	} else {
		q.HeadBlock = next
		q.TailBlock = next
		return nil
	}
}

// 删除队头的元素
func (q *PCBLinkedListChain) Pop() *ProcessControlBlock {
	var newHead *ProcessControlBlock
	var oldHead *ProcessControlBlock
	if q.HeadBlock.NextBlock != nil {
		newHead = q.HeadBlock.NextBlock
		oldHead = q.HeadBlock
		newHead.LastBlock = nil
		oldHead.NextBlock = nil
		q.HeadBlock = newHead
		return oldHead
	} else if q.HeadBlock != nil {
		oldHead = q.HeadBlock
		q.HeadBlock = nil
		return oldHead
	}
	fmt.Println("Error:", ErrQueueEmpty)
	return nil
}

func (q *PCBLinkedListChain) Traverse() {
	current := q.HeadBlock
	//fmt.Println("process in  queue")
	for {
		if current != nil {
			fmt.Println("WorkId:", current.Work.Id, "ProcessId:", current.Id)
			current = current.NextBlock
		} else {
			break
		}
	}
	//fmt.Println()
}
