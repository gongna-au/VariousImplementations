package SmoothWeightRoundRobin

import (
	"fmt"
	"testing"
)

func TestSmoothWeight(t *testing.T) {
	nodes := []*Node{
		{"a", 0, 5},
		{"b", 0, 1},
		{"c", 0, 1},
	}
	for i := 0; i < 7; i++ {
		best := SmoothWeightRoundRobin(nodes)
		if best != nil {
			fmt.Println(best.Name)
		}
	}

}
