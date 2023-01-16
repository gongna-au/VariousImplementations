package VirtualMemory

import (
	"fmt"
	"runtime"
	"testing"
)

func TestTask(t *testing.T) {
	InitTaskReceiver(runtime.NumCPU())
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

}
