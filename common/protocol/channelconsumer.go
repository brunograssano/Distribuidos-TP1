package protocol

import (
	dataStructures "github.com/brunograssano/Distribuidos-TP1/common/data_structures"
)

type ConsumerChannel struct {
	consumerChan chan *dataStructures.Message
	recvCount    int
}

func NewConsumerChannel(consumerChan chan *dataStructures.Message) *ConsumerChannel {
	return &ConsumerChannel{
		consumerChan: consumerChan,
		recvCount:    0,
	}
}

func (c *ConsumerChannel) Pop() (*dataStructures.Message, bool) {
	msg, ok := <-c.consumerChan
	if ok {
		if msg.TypeMessage == dataStructures.FlightRows {
			c.recvCount += len(msg.DynMaps)
		}
	}
	return msg, ok
}

func (c *ConsumerChannel) GetReceivedMessages() int {
	return c.recvCount
}

func (c *ConsumerChannel) ClearData() {
	c.recvCount = 0
}