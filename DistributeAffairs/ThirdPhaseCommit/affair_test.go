package ThirdPhaseCommit

import (
	"fmt"
	"testing"
)

func TestAffair(t *testing.T) {
	// 创建一连串的命令
	command1 := NewUpdateCommand()
	command1.setDb("mysql")
	command2 := NewUpdateCommand()
	command2.setDb("mongoDB")
	// 创建一个事务
	affair := NewAffairs(command1, command2)
	// 创建一个协调者
	coordinator := NewCoordinator()
	participant := NewParticipant()
	coordinator.SetParticipants(participant)
	if coordinator.CanCommit(affair) {
		if coordinator.PreCommit() {
			if coordinator.DoCommit() {
				fmt.Println("coordinator can return successfal message to client")
			}
		}
	}

}
