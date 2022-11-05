package SecondPhaseCommit

import "testing"

func TestAffair(t *testing.T) {
	affair := NewAffairs()
	participant := NewParticipant()
	coordinator := NewCoordinator()
	coordinator.SetParticipants(participant)
	affair.SetCoordinator(coordinator)
	// 开始执行
	affair.Begin("start")
	// 提交
	affair.Commit()
}
