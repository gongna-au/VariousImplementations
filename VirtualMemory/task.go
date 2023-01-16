package VirtualMemory

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	StateWaiting   = "waiting"
	StateCompleted = "completed"
	StateError     = "failed"
	StateNone      = "none"
	StateOverdue   = "overdue"
)

// 一个全局的任务池——复用任务对象——如何复用？——更改任务实例的参数就可以辣
// 一个全局的map存储所有任务的状态，键是任务全局唯一的任务ID ，值是任务的状态
// 一个全局的锁，用来保护pool或者保护map,在并发的场景下，map是并发不安全的，所以需要加锁，尤其是判断某个key在map中间在不在
var (
	taskStateMutex sync.Mutex
	taskPool       sync.Pool
	taskChan       = make(chan *Task, runtime.NumCPU())
	// 每个任务都有一个唯一的任务ID，这个被存储
	taskState = make(map[string]string)
)

// 这里的Task 更像是任务的执行者，执行者拿到不同的参数和数据，化身不同的任务实例
// NewTask里面复用了pool里面的对象，只是更改了参数和UUID而已，一个UUID 才是一个任务实例
// taskState对应的是任务实例
type Task struct {
	Param      map[string]interface{}
	UUID       string
	Expiration int64
	Factory    []FactFunc
}

// 函数接收到任务的ID 和参数map
type FactFunc func(string, map[string]interface{}) (string, error)

func NewTask(params map[string]interface{}, factory []FactFunc, d time.Duration) *Task {
	var expiration int64
	if d > 0 {
		expiration = time.Now().Add(d).UnixNano()
	} else {
		expiration = -1
	}
	t := taskPool.Get()
	if t == nil {
		// 要么创建新的任务
		return &Task{
			Param:      params,
			Factory:    factory,
			UUID:       getUUID(20),
			Expiration: expiration,
		}
	} else {
		// 要么更改任务携带的数据
		task := t.(*Task)
		(*task).Param = params
		(*task).Factory = factory
		// 当任务的携带的数据也被改变，也要更改它的ID
		(*task).UUID = getUUID(20)
		(*task).Expiration = expiration
		return task
	}
}

// 任务的生产者
// 往全局的channel 写入任务实例
func AddTask(task *Task) string {
	go func() {
		taskChan <- task
	}()
	uuid := (*task).UUID
	UpdateTaskState(uuid, StateWaiting)
	return uuid
}

// 任务的消费者
// 从全局的channel 中间取出任务实例
func taskReceiver() {
	var taskUUID string
	var err error
	for {
		// 从全局的channel 中间读取出任务
		task := <-taskChan
		if ((*task).Expiration > 0 && time.Now().UnixNano() < (*task).Expiration) || (*task).Expiration < 0 {
			for _, f := range (*task).Factory {
				// 把任务实例让factory里面的函数都执行一遍
				taskUUID, err = f((*task).UUID, (*task).Param)
			}
			if err != nil {
				// 给ID为那个任务标记状态
				// 更改全局任务的Map里面任务实例的状态
				UpdateTaskState(taskUUID, StateError)
				//不断的往Pool里面添加
				taskPool.Put(task)
			} else {
				// 给ID为那个任务标记状态
				// 更改全局任务的Map里面任务实例的状态
				UpdateTaskState(taskUUID, StateCompleted)
				taskPool.Put(task)
			}
		} else {
			UpdateTaskState(taskUUID, StateOverdue)
			taskPool.Put(task)
		}
	}
}

func InitTaskReceiver(num int) {
	for i := 0; i < num; i++ {
		go taskReceiver()
	}
}

func UpdateTaskState(uuid, state string) {
	taskStateMutex.Lock()
	defer taskStateMutex.Unlock()
	// 新的任务的ID的状态改变
	taskState[uuid] = state
}

func GetTaskState(uuid string) (state string) {
	taskStateMutex.Lock()
	defer taskStateMutex.Unlock()

	resultState, exists := taskState[uuid]
	if !exists {
		state = StateNone
	} else {
		state = resultState
	}
	return
}

func random(strings []string) ([]string, error) {
	rand.Seed(time.Now().UnixNano())
	for i := len(strings) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		strings[i], strings[num] = strings[num], strings[i]
	}

	str := make([]string, 0)
	for i := 0; i < len(strings); i++ {
		str = append(str, strings[i])
	}
	return str, nil
}

func getUUID(length int64) string {
	ele := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "v", "k",
		"l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G",
		"H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	ele, _ = random(ele)
	uuid := ""
	rand.Seed(time.Now().UnixNano())
	var i int64
	for i = 0; i < length; i++ {
		uuid += ele[rand.Intn(59)]
	}
	return uuid
}
