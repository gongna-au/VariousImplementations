package TokenBucket

import (
	"fmt"
	"sync"
)

type commonBucket struct {
	ch       chan Request
	capacity int
}

type Request struct {
	cmd  int
	data []byte
}

func NewRequest(cmd int, data []byte) Request {
	return Request{
		cmd:  cmd,
		data: data,
	}
}

func (r Request) GetData() string {
	return string(r.data)
}

func NewcommonBucket(cap int) *commonBucket {
	return &commonBucket{
		ch:       make(chan Request, cap),
		capacity: cap,
	}
}

func (t *commonBucket) AddRequest(request Request, wg *sync.WaitGroup) {
	t.ch <- request
	wg.Done()

}

func (t *commonBucket) GetRequest() {
	for {
		select {
		case temp := <-t.ch:
			fmt.Println("Get ", temp.GetData(), " Request")
		default:
			break
		}
	}
}

func (t *commonBucket) Run() {
	wg := &sync.WaitGroup{}
	requests := []Request{
		Request{
			cmd:  1,
			data: []byte("sign"),
		},
		Request{
			cmd:  2,
			data: []byte("login"),
		},
		Request{
			cmd:  3,
			data: []byte("search"),
		},
		Request{
			cmd:  4,
			data: []byte("create"),
		},
		Request{
			cmd:  5,
			data: []byte("delete"),
		},
		Request{
			cmd:  6,
			data: []byte("update"),
		},
	}

	go t.GetRequest()
	for _, v := range requests {
		wg.Add(1)
		temp := v
		go t.AddRequest(temp, wg)
	}
	wg.Wait()
}
