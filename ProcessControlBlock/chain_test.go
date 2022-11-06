package ProcessControlBlock

import (
	"fmt"
	"testing"
)

func TestChain(t *testing.T) {
	chain := &PCBLinkedListChain{}
	chain.Push(NewProcessControlBlock().SetStatus("created"))
	current := chain.HeadBlock
	current.NextBlock = NewProcessControlBlock().SetStatus("created")
	current.NextBlock.NextBlock = NewProcessControlBlock().SetStatus("created")
	for {
		if current == nil {
			break
		} else {
			fmt.Println("Id:", current.Id)
			current = current.NextBlock
		}
	}
}
