package ChannelApplicationScenario

import (
	"fmt"
	"sync"
)

type Handler struct {
	reqCh  chan int
	stopCh chan int
	wg     *sync.WaitGroup
}

func (h *Handler) Stop() {
	h.wg.Wait()
	close(h.stopCh)
	// 可以使用WaitGroup等待所有协程退出
}

func NewHandler() *Handler {
	return &Handler{
		reqCh:  make(chan int, 10),
		stopCh: make(chan int, 1),
		wg:     &sync.WaitGroup{},
	}
}

func (h *Handler) GetReq() {
	for i := 0; i < 10; i++ {
		h.reqCh <- i
	}
}

func (h *Handler) loop() {
	for i := 0; i < 10; i++ {
		select {
		case v := <-h.reqCh:
			h.wg.Add(1)
			go func() {
				handler(v)
				h.wg.Done()
			}()
		case <-h.stopCh:
			fmt.Println("return")
			return
		default:
			fmt.Println("infor")

		}
	}
}

func handler(req int) {
	fmt.Println(req, "is handler")
}

func Test() {
	handler := NewHandler()
	handler.GetReq()
	handler.loop()
	handler.Stop()
}
