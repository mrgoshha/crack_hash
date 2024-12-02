package app

import (
	"github.com/sirupsen/logrus"
	rabbit2 "worker/internal/rabbit"
	"worker/internal/services"
	"worker/internal/worker"
)

func Start(config *worker.Config) error {
	logger := logrus.New()

	if err := configureLogger(config.LogLevel); err != nil {
		logger.Infof("app init error: %v", err)
		return err
	}

	rabbit, err := rabbit2.New(config.RabbitMqUrl, logger)
	if err != nil {
		logger.Infof("app init error: %v", err)
		return err
	}
	defer rabbit.Connection.Close()
	defer rabbit.Producer.Channel.Close()
	defer rabbit.Consumer.Channel.Close()

	service := services.NewWorkerService(rabbit, logger)

	worker := worker.NewWorker(config, rabbit, service, logger)

	err = worker.Start()
	if err != nil {
		logger.Infof("app init error: %v", err)
		return err
	}

	return nil
}

func configureLogger(logLevel string) error {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(level)

	return nil
}
