package rabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"time"
)

type Rabbit struct {
	Connection *amqp.Connection
	Producer   *Producer
	Consumer   *Consumer
	logger     *logrus.Logger
	url        string
}

func New(url string, logger *logrus.Logger) (*Rabbit, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Infof("rabbit error: fialed to connect due to error %v", err)
		return nil, err
	}

	prod, err := NewProducer(conn, logger)
	if err != nil {
		logger.Infof("rabbit error: fialed to crate a producer due to error %v", err)
		return nil, err
	}
	cons, err := NewConsumer(conn, logger)
	if err != nil {
		logger.Infof("rabbit error: fialed to crate a consumer due to error %v", err)
		return nil, err
	}

	return &Rabbit{
		Connection: conn,
		logger:     logger,
		Producer:   prod,
		Consumer:   cons,
		url:        url,
	}, nil
}

func (rb *Rabbit) Reconnect() (*Rabbit, error) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			r, err := New(rb.url, rb.logger)
			rb.logger.Infof("rabbit: reconnection attempt, %s", time.Now())
			if err == nil {
				ticker.Stop()
				r.logger.Infof("rabbit: successful reconnection, %s", time.Now())
				return r, nil
			}
		}
	}
}
