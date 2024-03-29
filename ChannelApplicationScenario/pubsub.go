package ChannelApplicationScenario

import (
	"fmt"
	"sync"
)

type Pubsub struct {
	mu     sync.Mutex
	subs   map[string][]chan string
	closed bool
}

func NewPubsub() *Pubsub {
	return &Pubsub{
		mu:   sync.Mutex{},
		subs: make(map[string][]chan string),
	}
}

func (ps *Pubsub) Subscribe(topic string) <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ch := make(chan string, 1)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch
}

func (ps *Pubsub) Publish(topic string, msg string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if ps.closed {
		return
	}

	for _, ch := range ps.subs[topic] {
		fmt.Println(msg, "has been send to", topic)
		ch <- msg
	}
}

func (ps *Pubsub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if !ps.closed {
		ps.closed = true
		for _, subs := range ps.subs {
			for _, ch := range subs {
				close(ch)
			}
		}
	}
}

func TestPubsub() {
	ps := NewPubsub()
	ps.Subscribe("name")
	ps.Subscribe("age")
	ps.Subscribe("school")
	ps.Publish("name", "golang")
	ps.Publish("age", "17")
	ps.Publish("school", "CCNU")

}
