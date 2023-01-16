package CounterLimit

import (
	"fmt"
	"testing"
	"time"
)

func TestCommonCounter(t *testing.T) {
	now1 := time.Now()
	time.Sleep(1 * time.Second)
	now2 := time.Now()

	counter := Newcounter(now2.Sub(now1), 30)
	result := counter.Allow(3)
	fmt.Println(result)
	result = counter.Allow(20)
	fmt.Println(result)
	result = counter.Allow(10)
	fmt.Println(result)
	time.Sleep(1 * time.Second)
	result = counter.Allow(10)
	fmt.Println(result)
	time.Sleep(1 * time.Second)
	result = counter.Allow(40)
	fmt.Println(result)

}
