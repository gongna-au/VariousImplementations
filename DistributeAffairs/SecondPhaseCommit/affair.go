package SecondPhaseCommit

import (
	"fmt"
)

// 2PC（phase-commit） 实现
// 事务发起者
type Affairs struct {
	Coordinator ICoordination
}

func NewAffairs() *Affairs {
	return &Affairs{}
}

func (a *Affairs) SetCoordinator(c ICoordination) {
	a.Coordinator = c
}

// 事务发起者首先向协调者发起事务请求
func (a *Affairs) Begin(str string) {
	// 构建事务内容
	content := "content has build over"
	affaircontent := NewAffairsContent(str, content)
	// 协调者会给所有参与者发送 prepare 请求（其中包括事务内容)
	a.Coordinator.Prepare(affaircontent)
}

func (a *Affairs) Commit() {
	a.Coordinator.Commit()
}

// 定义事务内容
type AffairsContent struct {
	Name    string
	Content string
}

func NewAffairsContent(name string, content string) *AffairsContent {
	return &AffairsContent{
		Name:    name,
		Content: content,
	}
}

type ICoordination interface {
	Prepare(a *AffairsContent)
	Commit()
}

// 协调者
type Coordinator struct {
	Participants []IParticipate
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		Participants: []IParticipate{},
	}
}

func (c *Coordinator) SetParticipants(i ...IParticipate) {
	c.Participants = append(c.Participants, i...)
}

// 把具体的回滚任务交给了History结构体
func (c *Coordinator) Commit() {
	history := History{
		Participants: make([]IParticipate, len(c.Participants)),
	}
	for _, v := range c.Participants {
		isSuccess := v.Excute()
		if isSuccess {
			// 这里保证了前面一个成功才能执行下面一个
			history.Add(v)
			continue
		} else {
			history.Rollback()
		}
	}
}

func (c *Coordinator) Prepare(a *AffairsContent) {
	for _, v := range c.Participants {
		v.Prepare(a)
	}
}

type History struct {
	Participants []IParticipate
}

func (h *History) Add(i IParticipate) {
	h.Participants = append(h.Participants, i)
}

// 对已经添加的命令执行撤销操作 封装了一下
func (c *History) Rollback() {
	for i := len(c.Participants) - 1; i >= 0; i-- {
		c.Participants[i].Undo()
	}
}

type IParticipate interface {
	// 准备
	Prepare(a *AffairsContent)
	//执行
	Excute() bool
	// 撤销执行
	Undo()
}

// 参与者
type Participant struct {
	Integral int
}

func NewParticipant() *Participant {
	return &Participant{}
}

func (c *Participant) Prepare(a *AffairsContent) {
	fmt.Printf("Participant has get affair name: %s\n", a.Name)
	fmt.Printf("Participant has get affair content: %s\n", a.Content)
	fmt.Println("Participant Prepared to excute affair")
}

func (c *Participant) Excute() bool {
	// 给积分加10
	c.Integral = c.Integral + 10
	fmt.Println("Participant excute successfally!")
	return true
}

// 否则撤销对积分的操作
func (c *Participant) Undo() {
	fmt.Println("Participant excute undo!")
	c.Integral = c.Integral - 10
}
