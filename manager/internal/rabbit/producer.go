package rabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
	logger     *logrus.Logger
}

func NewProducer(conn *amqp.Connection, logger *logrus.Logger) (*Producer, error) {
	channel, err := conn.Channel()
	if err != nil {
		logger.Infof("producer error: failed to open a Channel due to error %v", err)
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		"task_queue",
		true, // durable - чтобы очередь выдержала перезапуск узла RabbitMQ нужно объявить ее как долговечную.
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Infof("producer error: failed to declare a Queue due to error %v", err)
		return nil, err
	}

	return &Producer{
		connection: conn,
		Channel:    channel,
		Queue:      queue,
	}, nil
}
