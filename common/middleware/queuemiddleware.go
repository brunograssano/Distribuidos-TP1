package middleware

import (
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type QueueMiddleware struct {
	channel *amqp.Channel
	conn    *amqp.Connection
}

func NewQueueMiddleware(address string) *QueueMiddleware {
	conn, err := amqp.Dial(address)
	FailOnError(err, "Failed to connect via Dial to RabbitMQ.")
	log.Infof("Connected to RabbitMQ")
	ch, err := conn.Channel()
	FailOnError(err, "Failed to create RabbitMQ Channel.")
	log.Infof("Created RabbitMQ Channel")
	return &QueueMiddleware{
		channel: ch,
		conn:    conn,
	}
}

func (qm *QueueMiddleware) CreateConsumer(name string, durable bool) ConsumerInterface {
	return NewConsumer(qm.channel, name, durable)
}

func (qm *QueueMiddleware) CreateProducer(name string, durable bool) ProducerInterface {
	return NewProducer(qm.channel, name, durable)
}

func (qm *QueueMiddleware) Close() {
	err := qm.channel.Close()
	FailOnError(err, "Error closing QueueMiddleware Channel")
	err = qm.conn.Close()
	FailOnError(err, "Error closing QueueMiddleware Connection")
}