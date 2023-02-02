package main

import (
	"fmt"
	"sync"

	"go.etcd.io/etcd/client/v3/concurrency"
	"go.etcd.io/etcd/clientv3"
)

var wg sync.WaitGroup

func NewCounter() *Counter {
	return &Counter{}
}

type Counter struct {
	count int
	mu    sync.Mutex
}

func (m *Counter) Inc() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.count++
}

func (m *Counter) Get() int {
	return m.count
}

func main() {
	endpoints := []string{"http://127.0.0.1:12379", "http://127.0.0.1:22379", "http://127.0.0.1:32379"}
	// 初始化etcd客户端
	client, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()
	counter := &Counter{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			// 这里会生成租约，默认是60秒
			session, err := concurrency.NewSession(client)
			if err != nil {
				panic(err)
			}
			defer session.Close()

			locker := concurrency.NewLocker(session, "/my-test-lock")
			locker.Lock()
			counter.Inc()
			locker.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println("count:", counter.Get())
}

// 实现了阻塞接口的锁，即当前有获取到锁的请
