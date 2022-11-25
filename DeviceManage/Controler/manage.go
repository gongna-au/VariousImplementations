package Controler

import (
	"container/list"
	"fmt"
	"strings"
)

// 输入控制器
type InputControler struct {
	InputBuffer     *InputBuffer
	InputFilesQueue *InputFilesQueue
}

func NewInputControler() *InputControler {
	return &InputControler{
		InputBuffer:     NewInputBuffer(10),
		InputFilesQueue: NewInputFilesQueue(),
	}
}

func (i *InputControler) Schedule() {
	fmt.Println("InputControler start schedule")
	request := i.InputBuffer.Pop()
	if request == nil {
		return
	}
	file := ConvertRequestToFILE(request)
	i.InputFilesQueue.Add(file)
	fmt.Println("InputControler schedule successfully")
}

// 输出控制器
type OutputControler struct {
	OutputBuffer     *OutputBuffer
	OutputFilesQueue *OutputFilesQueue
}

func NewOutputControler() *OutputControler {
	return &OutputControler{
		OutputBuffer:     NewOutputBuffer(10),
		OutputFilesQueue: NewOutputFilesQueue(),
	}
}

func (o *OutputControler) Schedule() *Request {
	fmt.Println("OutputControler start schedule")
	// 从输出井中取出数据
	file := o.OutputFilesQueue.Pop()
	if file == nil {
		return nil
	}
	// 把取出的文件转化为数据结构存放在输出的缓冲区
	data := ConvertFILEToRequest(file)
	o.OutputBuffer.Add(data)
	return o.OutputBuffer.Pop()
}

type Request struct {
	Content string
}

func NewRequest(content string) *Request {
	return &Request{
		Content: "Request:" + content,
	}
}

// 输入缓冲
type InputBuffer struct {
	num  int
	list *list.List
}

func NewInputBuffer(num int) *InputBuffer {
	return &InputBuffer{
		num:  num,
		list: list.New(),
	}
}

func (i *InputBuffer) Add(request *Request) {
	if i.list.Len()+1 <= i.num {
		i.list.PushBack(request)
		fmt.Println(request.Content, " has been added in inputBuffer!")
	}
}

func (i *InputBuffer) Pop() *Request {

	head := i.list.Front()
	if head == nil {
		fmt.Println("InputBuffer has storing any request element")
		return nil
	}
	i.list.Remove(head)
	if result, ok := head.Value.(*Request); ok {
		//fmt.Println("InputBuffer 断言成功")
		fmt.Println(result.Content, " has been poped from inputBuffer!")
		return result
	}

	fmt.Print("InputBuffer 断言失败")
	return nil
}

// 输出缓冲
type OutputBuffer struct {
	num  int
	list *list.List
}

func NewOutputBuffer(num int) *OutputBuffer {
	return &OutputBuffer{
		num:  num,
		list: list.New(),
	}
}

func (i *OutputBuffer) Add(request *Request) {
	if i.list.Len()+1 <= i.num {
		i.list.PushBack(request)
		fmt.Println(request.Content, " has been added in outputBuffer!")
	}
}

func (i *OutputBuffer) Pop() *Request {
	head := i.list.Front()
	if head == nil {
		fmt.Println("OutputBuffer has storing any request element")
		return nil
	}
	i.list.Remove(head)
	if result, ok := head.Value.(*Request); ok {
		//fmt.Println("OutputBuffer 断言成功")
		fmt.Println(result.Content, " has been poped from outputBuffer!")
		return result
	}

	fmt.Println("OutputBufferPop断言失败")
	return nil
}

type FILE struct {
	Data string
}

func NewFILE(data string) *FILE {
	return &FILE{
		Data: "FILE:" + data,
	}
}

// 输入井（输入文件队列）
type InputFilesQueue struct {
	list *list.List
}

func NewInputFilesQueue() *InputFilesQueue {
	return &InputFilesQueue{
		list: list.New(),
	}
}

func (i *InputFilesQueue) Add(file *FILE) {
	fmt.Println(file.Data, " has been added in InputFilesQueue(输入井)!")
	i.list.PushBack(file)
}

func (i *InputFilesQueue) Pop() *FILE {
	head := i.list.Front()
	if head == nil {
		fmt.Println("InputFilesQueue has no file")
		return nil
	}
	i.list.Remove(head)
	if result, ok := head.Value.(*FILE); ok {
		//fmt.Println("InputFilesQueue 断言成功")
		fmt.Println(result.Data, " has been poped from InputFilesQueue(输入井)!")
		return result
	}
	fmt.Println("InputFilesQueue断言失败")
	return nil
}

// 输出井（输出文件队列）
type OutputFilesQueue struct {
	list *list.List
}

func NewOutputFilesQueue() *OutputFilesQueue {
	return &OutputFilesQueue{
		list: list.New(),
	}
}

func (o *OutputFilesQueue) Add(file *FILE) {
	fmt.Println(file.Data, " has been added in OutputFilesQueue(输出井)!")
	o.list.PushBack(file)
}

func (o *OutputFilesQueue) Pop() *FILE {
	head := o.list.Front()
	if head == nil {
		fmt.Println("OutputFilesQueue has no file")
		return nil
	}
	o.list.Remove(head)
	if result, ok := head.Value.(*FILE); ok {
		//fmt.Println("OutputFilesQueue 断言成功")
		fmt.Println(result.Data, " has been poped from OutputFilesQueue(输出井)!")
		return result
	}
	fmt.Println("OutputFilesQueue断言失败")
	return nil
}

func ConvertRequestToFILE(request *Request) *FILE {
	if request != nil {
		return NewFILE(SplitString(request.Content))
	}
	fmt.Println("No Request can be Converted to FILE ")
	return nil
}

func ConvertFILEToRequest(file *FILE) *Request {
	if file != nil {
		return NewRequest(SplitString(file.Data))
	}
	fmt.Println("No FILE can be Converted to Request ")
	return nil
}

func MoveHeadToTail(InputQueue *InputFilesQueue, OutputQueue *OutputFilesQueue) {
	head := InputQueue.list.Front()
	if head == nil {
		fmt.Println("InputFilesQueue has not FIFE data")
		return
	}
	InputQueue.list.Remove(head)
	if v, ok := head.Value.(*FILE); ok {
		OutputQueue.list.PushBack(v)
		fmt.Println(v.Data, "in 输入井 (head) has been moved to 输出井 (tail) successfully")
	}
}

type SpoolingSystem struct {
	InputControler  *InputControler
	OutputControler *OutputControler
}

func NewSpoolingSystem() *SpoolingSystem {
	return &SpoolingSystem{
		InputControler:  NewInputControler(),
		OutputControler: NewOutputControler(),
	}
}

func (s *SpoolingSystem) Schedele() {
	// 随机设置一些假数据模拟磁盘和缓冲区已经有数据
	// 假设输入井中已经有一些数据
	// 输入输出井中的数据加上前缀"FIFE"表示是文件格式
	s.InputControler.InputFilesQueue.Add(NewFILE("data1"))
	s.InputControler.InputFilesQueue.Add(NewFILE("data2"))
	s.InputControler.InputFilesQueue.Add(NewFILE("data3"))
	// 继续往输入缓冲区放入请求
	s.InputControler.InputBuffer.Add(NewRequest("data4"))
	s.InputControler.InputBuffer.Add(NewRequest("data5"))
	s.InputControler.InputBuffer.Add(NewRequest("data6"))

	for s.InputControler.InputFilesQueue.list.Len() != 0 {
		fmt.Println()
		s.InputControler.Schedule()
		fmt.Println()
		//fmt.Println(s.InputControler.InputBuffer.list.Len())
		MoveHeadToTail(s.InputControler.InputFilesQueue, s.OutputControler.OutputFilesQueue)
		fmt.Println()

		request := s.OutputControler.Schedule()
		fmt.Println()

		PrintRequest(request)
	}
	/* s.InputControler.Schedule()
	MoveHeadToTail(s.InputControler.InputFilesQueue, s.OutputControler.OutputFilesQueue)
	s.OutputControler.Schedule()

	s.InputControler.Schedule()
	MoveHeadToTail(s.InputControler.InputFilesQueue, s.OutputControler.OutputFilesQueue)
	s.OutputControler.Schedule()

	s.InputControler.Schedule()
	MoveHeadToTail(s.InputControler.InputFilesQueue, s.OutputControler.OutputFilesQueue)
	s.OutputControler.Schedule() */

}

func SplitString(str string) string {
	result := strings.Split(str, ":")
	if len(result) == 2 {
		return result[1]
	}
	return "wrongdata"
}

func PrintRequest(request *Request) {
	if request != nil {
		fmt.Println("Content:", request.Content, " has been put to PrintDevice(打印机设备)")
	}
}
