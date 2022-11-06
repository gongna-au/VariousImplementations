package ProcessControlBlock

import (
	"fmt"
	"testing"
)

func TestCPUSchedule(t *testing.T) {
	pool, err := NewPool(100, 7)
	if err != nil {
		fmt.Println("error:", err)
	}
	NewCpu(pool).Schedule()
}
