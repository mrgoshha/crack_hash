package app

import (
	"context"
	"github.com/sirupsen/logrus"
	"log"
	"manager/internal/manager"
	"manager/internal/rabbit"
	"manager/internal/services"
	"manager/internal/store/mongodb"
)

func Start(config *manager.Config) error {
	logger := logrus.New()

	if err := configureLogger(config.LogLevel); err != nil {
		logger.Infof("app init error: %v", err)
		return err
	}

	client, db, err := NewClient(context.Background(),
		"crackHash",
		config.MongoUrl,
		logger)
	if err != nil {
		logger.Infof("app init error: %v", err)
		return err
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	requestRepository := mongodb.NewRequestRepository(db, "requests", logger)

	rabbit, err := rabbit.New(
		config.RabbitMQUrl,
		logger)
	if err != nil {
		logger.Infof("app init error: %v", err)
		return err
	}
	defer rabbit.Connection.Close()
	defer rabbit.Producer.Channel.Close()
	defer rabbit.Consumer.Channel.Close()

	service := services.NewManagerService(requestRepository, rabbit, logger)

	manager := manager.NewManager(config, rabbit, service, logger)

	err = manager.Start()
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
