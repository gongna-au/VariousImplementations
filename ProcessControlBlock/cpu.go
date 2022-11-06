package ProcessControlBlock

import "fmt"

// 模拟CPU从池中拿出一个任务，然后执行5min 后变为就绪状态的情况，以及执行结束的情况 ，以及被阻塞的情况

type Cpu struct {
	Pool *PCBdPool
}

func NewCpu(p *PCBdPool) *Cpu {
	return &Cpu{
		Pool: p,
	}
}

func (c *Cpu) Schedule() {
	fmt.Println("PCB in pool is:")
	c.Pool.TraverseCreated()
	var i int
	for {
		if i > 35 {
			break
		}
		if i == 0 {
			// 唤醒一个pool 已经被创建的进程
			fmt.Println("Time is :", i, "唤醒一个pool 已经被创建的进程")
			process := c.Pool.QueuesCreate.Pop()
			process.SetStatus("ready")
			c.Pool.QueuesReady.Push(process)
			c.Pool.TraverseCreated()
			c.Pool.TraverseReady()
		}
		if i == 5 {
			fmt.Println("Time is :", i, "选取一个已经就绪的PCB 让它运行")
			//把就绪的让他运行
			process := c.Pool.QueuesReady.Pop()
			process.SetStatus("running")
			c.Pool.QueuesRun.Push(process)
			c.Pool.TraverseReady()
			c.Pool.TraverseRun()
		}
		if i == 10 {
			fmt.Println("Time is :", i, "一个已经运行的PCB用完CPU 分配给它的时间片后它变成就绪状态")
			//把运行的让他时间片用完变成就绪
			process := c.Pool.QueuesRun.Pop()
			process.SetStatus("ready")
			c.Pool.QueuesReady.Push(process)
			c.Pool.TraverseRun()
			c.Pool.TraverseReady()

		}

		if i == 15 {
			fmt.Println("Time is :", i, "一个已经就绪的PCB得到CPU 分配给它的时间片后它变成运行状态")
			//把就绪的进程再次让他运行
			process := c.Pool.QueuesReady.Pop()
			process.SetStatus("running")
			c.Pool.QueuesRun.Push(process)
			c.Pool.TraverseReady()
			c.Pool.TraverseRun()
		}
		if i == 20 {
			//把运行的进程再次让他阻塞
			fmt.Println("Time is :", i, "一个已经运行的PCB被I/O设备阻塞,变成阻塞状态")
			process := c.Pool.QueuesRun.Pop()
			process.SetStatus("block")
			c.Pool.QueuesBlocked.Push(process)
			c.Pool.TraverseRun()
			c.Pool.TraverseBlock()
		}
		if i == 25 {
			//把阻塞的进程再次让他就绪
			fmt.Println("Time is :", i, "一个已经阻塞的PCB等待的I/O事件完成,变为就绪状态")
			process := c.Pool.QueuesBlocked.Pop()
			process.SetStatus("ready")
			c.Pool.QueuesReady.Push(process)
			c.Pool.TraverseBlock()
			c.Pool.TraverseReady()

		}
		if i == 30 {
			fmt.Println("Time is :", i, "一个已经就绪的PCB得到CPU 分配给它的时间片后它变成运行状态")
			//把就绪的进程再次让他运行
			process := c.Pool.QueuesReady.Pop()
			process.SetStatus("running")
			c.Pool.QueuesRun.Push(process)
			c.Pool.TraverseReady()
			c.Pool.TraverseRun()
		}
		if i == 35 {
			fmt.Println("Time is :", i, "一个已经运行的PCB执行完任务后它变成结束状态")
			//把就绪的进程再次让他运行
			process := c.Pool.QueuesRun.Pop()
			process.SetStatus("block")
			c.Pool.QueuesTerminated.Push(process)
			c.Pool.TraverseRun()
			c.Pool.TraverseTerminated()
		}
		i = i + 5
	}

}
