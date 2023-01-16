package DeviceManage

import "fmt"

func Run() {
	array := make([]int, 4)
	array[0] = 0
	array[1] = 1
	array[2] = 2
	array[3] = 3
	for k, v := range array {
		if k == 0 {
			array = append(array, 4)
		}
		if k == 1 {
			array[2] = 7
		}
		fmt.Println(v)
	}
	fmt.Println(array)
}

func AppendChangeAndChange(array []int) {
	array[0] = 1
	array = append(array, 4)
}

func AppendChange(array []int) {
	array = append(array, 4)
}

func Change(array []int) {
	array[0] = 1
}
