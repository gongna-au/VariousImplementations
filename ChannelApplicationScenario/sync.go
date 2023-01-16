package ChannelApplicationScenario

import (
	"fmt"
	"sync"
	"time"
)

func MutexWait() {
	var mutex sync.Mutex
	mutex.Lock()
	go func() {
		fmt.Println("User1 is coming")
		mutex.Unlock()
	}()
	mutex.Lock()
	fmt.Println("User1 has came,User2 can go")
}

func ChanWait() {
	ch := make(chan int)
	go func() {
		fmt.Println("User1 is coming")
		ch <- 1
	}()
	<-ch
	fmt.Println("User1 has came,User2 can go")
}

func ManyChanWait() {
	done := make(chan int, 10)
	for i := 0; i < cap(done); i++ {
		temp := i
		go func() {
			fmt.Println("User", fmt.Sprint(temp), " is coming")
			done <- 1
		}()
	}
	for i := 0; i < cap(done); i++ {
		<-done
	}
	fmt.Println("All user has came, can go")

}

func WaitGroupWait() {
	var wg sync.WaitGroup
	done := make(chan int, 10)
	for i := 0; i < cap(done); i++ {
		temp := i
		wg.Add(1)
		go func() {
			fmt.Println("User", fmt.Sprint(temp), " is coming")
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("All user has came, can go")
}

func Producer1(out chan<- int) {
	for i := 0; i < 5; i++ {
		out <- i
		temp := i
		fmt.Println(fmt.Sprint(temp), "is be producedd by Producer1 ")
	}
}

func Producer2(out chan<- int) {
	for i := 5; i < 10; i++ {
		out <- i
		temp := i
		fmt.Println(fmt.Sprint(temp), "is be producedd by Producer2 ")
	}
}

func Consume(out <-chan int) {
	for v := range out {
		temp := v
		fmt.Println(fmt.Sprint(temp), "is be consumed ")
	}
}

func ConsumeProduceModel() {
	ch := make(chan int, 10)
	go Producer1(ch)
	go Producer2(ch)
	go Consume(ch)
	time.Sleep(5 * time.Second)

}

// 传统生产者和消费者模型中，是将消息发送到一个队列中，而发布订阅模型则是将消息发布给一个主题。
func PublishSubscribe() {
	ch := make(chan int, 10)
	go Producer1(ch)
	go Producer2(ch)
	go Consume(ch)
	time.Sleep(5 * time.Second)
}

type (
	// 订阅者是一个
	subscriber chan interface{}
	// 主题和订阅者之间是1对多的关系
	topicFunc func(v interface{}) bool
)

type Publisher struct {
	subscribers map[subscriber]topicFunc
}
