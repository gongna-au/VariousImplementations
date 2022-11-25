package RandomSchedule

import (
	"container/list"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

type PCBsControl struct {
	PCBs []*PCB
}

func randomData() []string {
	result := []string{}
	ran, _ := rand.Int(rand.Reader, big.NewInt(10))
	datanum := ran.Int64()
	for i := 0; int64(i) <= datanum; i++ {
		result = append(result, "DataContent:", fmt.Sprint(i))
	}
	return result
}

// 在这里模拟的是每个进程是可以多次请求打印的
func NewPCBsControl(num int) *PCBsControl {
	pcbs := []*PCB{}
	for i := 0; i < num; i++ {
		pcb := NewPCB(i, 0, "userProcess-"+fmt.Sprint(i))
		pcb.SetContent(randomData())
		pcbs = append(pcbs, pcb)
	}
	return &PCBsControl{
		PCBs: pcbs,
	}
}

type PCB struct {
	// 进程标志ID
	ID int
	// 进程的状态
	Status int
	// 标明是用户进程还是Spooling 进程
	Name string
	// 起始地址
	StartAddress int
	// 长度
	Legth int
	// 请求输出的内容
	ContentBuffer []string
}

func NewPCB(id int, status int, Name string) *PCB {
	return &PCB{
		ID:            id,
		Status:        status,
		Name:          Name,
		ContentBuffer: []string{},
	}
}

func (p *PCB) SetStartAddress(start int) *PCB {
	p.StartAddress = start
	return p
}

func (p *PCB) SetLength(length int) *PCB {
	p.Legth = length
	return p
}

func (p *PCB) SetContent(content []string) *PCB {
	p.ContentBuffer = append(p.ContentBuffer, content...)
	return p
}

// 若没有可用请求块时，调用进程进入"等待状态3
func (p *PCB) SetStatus(status int) {
	p.Status = status
}

// 请求输出
func (p *PCB) RequestPrintOut() *RequestOutPutBlock {
	return NewRequestOutPutBlock(p.ID, p.ContentBuffer...)
}

// 模拟输出井
type OutputFilesQueue struct {
	// 当容量超过限制就先等待
	Capacity int
	List     *list.List
}

// list内部存储的元素
// 把每个进程要求输出的内容包装成请求输出块
type RequestOutPutBlock struct {
	// 要求输出的进程的ID
	PcbID int
	// 要求输出的内容
	ContentBuffer []string
}

func NewRequestOutPutBlock(id int, content ...string) *RequestOutPutBlock {
	fmt.Println("PCB:", id, " request to print ", content)
	buffer := []string{}
	return &RequestOutPutBlock{
		PcbID:         id,
		ContentBuffer: append(buffer, content...),
	}
}

// 模拟输出请求块
func NewOutputFilesQueue(cap int) *OutputFilesQueue {
	return &OutputFilesQueue{
		List:     list.New(),
		Capacity: cap,
	}
}

func (o *OutputFilesQueue) Add(request *RequestOutPutBlock) {
	if o.List.Len() <= o.Capacity {
		o.List.PushBack(request)
	}
}

func (o *OutputFilesQueue) Pop() *RequestOutPutBlock {
	head := o.List.Front()
	if head == nil {
		return nil
	}
	result, ok := head.Value.(*RequestOutPutBlock)
	if !ok {
		return nil
	}
	o.List.Remove(head)
	return result
}

type SpoolingSystemControl struct {
	OutputFilesQueue *OutputFilesQueue
}

func NewSpoolingSystemControl(cap int) *SpoolingSystemControl {
	return &SpoolingSystemControl{
		OutputFilesQueue: NewOutputFilesQueue(cap),
	}
}

func (s *SpoolingSystemControl) Spooling() {
	fmt.Println("Spooling is be scheduled")
	block := s.Pop()
	if block != nil {
		fmt.Println("Output data content is : user-", block.PcbID, "data", block.ContentBuffer, "")
		fmt.Println("输出数据后输出井的大小为：", s.OutputFilesQueue.List.Len())
	}
	fmt.Println("There is not content to Print")
}

func (s *SpoolingSystemControl) Add(request *RequestOutPutBlock) {
	s.OutputFilesQueue.Add(request)
}

func (s *SpoolingSystemControl) Pop() *RequestOutPutBlock {
	return s.OutputFilesQueue.Pop()
}

type CPU struct {
	PCBsControl           *PCBsControl
	SpoolingSystemControl *SpoolingSystemControl
}

func NewCPU() *CPU {
	return &CPU{}
}

func (c *CPU) SetPCBsControl(p *PCBsControl) *CPU {
	c.PCBsControl = p
	return c
}
func (c *CPU) SetSpoolingSystemControl(s *SpoolingSystemControl) *CPU {
	c.SpoolingSystemControl = s
	return c
}

func (c *CPU) Schedele() {
	for i := 0; i < 10; i++ {
		fmt.Println(time.Now())
		// 真随机
		// ran, _ := rand.Int(rand.Reader, big.NewInt(100))
		// 伪随机
		judge := GetRandom(100)
		if judge >= 0 && judge <= 45 {
			if len(c.PCBsControl.PCBs) >= 1 && (c.PCBsControl.PCBs[0].Status == 0) {
				request := c.PCBsControl.PCBs[0].RequestPrintOut()
				c.SpoolingSystemControl.OutputFilesQueue.Add(request)
				fmt.Println("输出井的大小为：", c.SpoolingSystemControl.OutputFilesQueue.List.Len())
				fmt.Println()

			}
		} else if judge > 45 && judge <= 90 {
			if len(c.PCBsControl.PCBs) >= 2 && (c.PCBsControl.PCBs[1].Status == 0) {
				request := c.PCBsControl.PCBs[1].RequestPrintOut()
				c.SpoolingSystemControl.OutputFilesQueue.Add(request)
				fmt.Println("输出井的大小为：", c.SpoolingSystemControl.OutputFilesQueue.List.Len())
				fmt.Println()

			}
		} else if judge > 90 && judge <= 100 {
			c.SpoolingSystemControl.Spooling()
			fmt.Println()
		}
	}
}

func Schedele() {
	// 2个用户进程
	// 2个用户进程分别请求的次数是3次
	// 表示有两个进程同时请求打印三次
	cpu := NewCPU().SetSpoolingSystemControl(NewSpoolingSystemControl(100)).SetPCBsControl(NewPCBsControl(2))
	cpu.Schedele()
}

func GetRandom(bigest int64) int64 {
	ran, _ := rand.Int(rand.Reader, big.NewInt(bigest))
	result := ran.Int64()
	fmt.Println("produce:", result)
	return result
}
