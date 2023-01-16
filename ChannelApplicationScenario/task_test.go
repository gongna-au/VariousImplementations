package ChannelApplicationScenario

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	InitTaskReceiver(runtime.NumCPU())
	for i := 0; i < 5; i++ {
		temp := i
		go func() {
			fmt.Println("Produce Task", temp)
			AddTask(
				NewTask(
					map[string]interface{}{
						"name": "Anna",
						"age":  23,
					},
					[]FactFunc{
						func(uuid string, params map[string]interface{}) (string, error) {
							fmt.Println(uuid)
							fmt.Println(params["name"])
							return uuid, nil
						},
						func(uuid string, params map[string]interface{}) (string, error) {
							fmt.Println(uuid)
							fmt.Println(params["age"])
							return uuid, nil
						},
					},
					-1,
				),
			)
		}()
	}
	time.Sleep(6 * time.Second)
}
