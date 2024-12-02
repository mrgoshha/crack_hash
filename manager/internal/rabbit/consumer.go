package rabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
	logger     *logrus.Logger
}

func NewConsumer(conn *amqp.Connection, logger *logrus.Logger) (*Consumer, error) {
	channel, err := conn.Channel()
	if err != nil {
		logger.Infof("consumer error: failed to open a Channel due to error %v", err)
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		"results_queue",
		true, // durable - чтобы очередь выдержала перезапуск узла RabbitMQ нужно объявить ее как долговечную.
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Infof("consumer error: failed to declare a Queue due to error %v", err)
		return nil, err
	}

	return &Consumer{
		connection: conn,
		Channel:    channel,
		Queue:      queue,
	}, nil
}
