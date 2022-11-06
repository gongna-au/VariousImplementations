package ProcessControlBlock

import (
	"fmt"
	"testing"
)

func TestNewPool(t *testing.T) {
	pool, err := NewPool(100, 7)
	if err != nil {
		fmt.Println("error:", err)
	}
	pool.PutCreated(NewProcessControlBlock().SetStatus("created"))
	pool.TraverseCreated()
	fmt.Println("has poped:")
	p := pool.PopCreated()
	fmt.Println("ID", p.Id)
	fmt.Println("pool after poped:")
	pool.TraverseCreated()
}
