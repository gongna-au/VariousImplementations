package ThirdPhaseCommit

import (
	"fmt"
)

// 事务就是一连串的命令
type Affairs struct {
	commands []Command
}

func NewAffairs(c ...Command) *Affairs {
	return &Affairs{
		commands: append([]Command{}, c...),
	}
}

// 真正的执行者执行调用的执行接口
type Command interface {
	// Exec 执行insert、update、delete命令
	Exec() error
	// Undo 回滚命令
	Undo()
	// SetDb 设置关联的数据库
	setDb(db string)
}

type UpdateCommand struct {
	DB string
}

func NewUpdateCommand() *UpdateCommand {
	return &UpdateCommand{}
}

func (u *UpdateCommand) Exec() error {
	fmt.Print(u.DB)
	fmt.Println(" update command exec successfully!")
	return nil
}
func (u *UpdateCommand) Undo() {
	fmt.Print(u.DB)
	fmt.Println("Update command exec undo successfully!")
}
func (u *UpdateCommand) setDb(db string) {
	u.DB = db
	fmt.Print(u.DB)
	fmt.Println("Update command set DB successfully!")
}

// History 结构体保存已经执行过的Command
type History struct {
	Commands []Command
}

func (h *History) Add(i ...Command) {
	h.Commands = append(h.Commands, i...)
}

// 对已经添加的命令执行撤销操作 封装了一下
func (c *History) Rollback() {
	for i := len(c.Commands) - 1; i >= 0; i-- {
		c.Commands[i].Undo()
	}
}

type ICoordination interface {
	// 协调者向参与者询问是否准备好了
	CanCommit(a *Affairs) bool
	// 预被提交
	PreCommit() bool
	// 真正的提交
	DoCommit() bool
}

// 协调者
type Coordinator struct {
	Participants []IParticipate
	Affairs      *Affairs
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		Participants: []IParticipate{},
	}
}

func (c *Coordinator) SetParticipants(i ...IParticipate) {
	c.Participants = append(c.Participants, i...)
}

// CanCommit 实现ICoordination接口
func (c *Coordinator) CanCommit(a *Affairs) bool {
	c.Affairs = a
	for _, v := range c.Participants {
		success := v.CanCommit(a)
		if !success {
			return false
		}
	}
	return true
}

// PreCommit 实现ICoordination接口
func (c *Coordinator) PreCommit() bool {
	array := []IParticipate{}
	for _, v := range c.Participants {
		success := v.PreCommit()
		if success {
			// 标记它已经准备好了
			array = append(array, v)
		} else {
			// 让所有的参与者终端中断执行事务
			for _, v2 := range array {
				v2.StopAffair()
			}
			return false
		}
	}
	return true
}

//  DoCommit 实现ICoordination接口
func (c *Coordinator) DoCommit() bool {
	// 进入第三阶段参与者都会进行事务的提交操作。
	for _, v := range c.Participants {
		v.DoCommit()
	}
	return true
}

type IParticipate interface {
	CanCommit(a *Affairs) bool
	PreCommit() bool
	DoCommit()
	// 接收到协调者告诉它的中断的消息
	StopAffair()
}

// 参与者
type Participant struct {
	Integral int
	Affair   *Affairs
	StatusOK bool
}

func NewParticipant() *Participant {
	return &Participant{}
}

// CanCommit 实现ICoordination接口
func (c *Participant) CanCommit(a *Affairs) bool {

	c.Affair = a
	fmt.Println("Participant can commit")
	if c.Integral == 0 {
		return true
	} else {
		return false
	}
}

// PreCommit 实现ICoordination接口 参与者收到预提交请求后
// 会进行事务的执行操作并将 `Undo` 和 `Redo` 信息写入事务日志中
func (c *Participant) PreCommit() bool {
	fmt.Println("Participant prepare commit")
	for _, v := range c.Affair.commands {
		v.Exec()
	}
	c.StatusOK = true
	return true
}

func (c *Participant) StopAffair() {
	fmt.Println("Participant stop affair")
	c.StatusOK = false
}

//  DoCommit 实现ICoordination接口
func (c *Participant) DoCommit() {
	fmt.Println("Participant do commit")

}
