package ChannelApplicationScenario

type ChannelLimit struct {
	bufferChannel chan int
}

func NewChannelLimit() *ChannelLimit {
	return &ChannelLimit{
		bufferChannel: make(chan int, 10),
	}
}

func (c *ChannelLimit) Allow() bool {
	select {
	case c.bufferChannel <- 1:
		return true
	default:
		return false
	}
}

func (c *ChannelLimit) Release() bool {
	<-c.bufferChannel
	return true
}

func (c *ChannelLimit) Close() {
	close(c.bufferChannel)
}
