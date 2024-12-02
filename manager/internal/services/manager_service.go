package services

import (
	"context"
	"encoding/xml"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"manager/api/hash"
	"manager/internal/model"
	"manager/internal/rabbit"
	"manager/internal/store"
	"slices"
	"time"
)

type ManagerService struct {
	requestRepository store.RequestRepository
	Rabbit            *rabbit.Rabbit
	partCount         int
	logger            *logrus.Logger
}

func NewManagerService(requestRepository store.RequestRepository, rabbit *rabbit.Rabbit, logger *logrus.Logger) *ManagerService {
	ms := &ManagerService{
		requestRepository: requestRepository,
		Rabbit:            rabbit,
		partCount:         4,
		logger:            logger,
	}
	go ms.checkRequests()
	return ms
}

func (ms *ManagerService) CreateCrackHashTask(req *hash.CrackHashRequest) (*model.HashRequest, error) {
	taskId, err := ms.requestRepository.Create(req.Hash, req.MaxLength)
	if err != nil {
		return nil, err
	}
	task, err := ms.requestRepository.GetRequestById(taskId)
	if err != nil {
		return nil, err
	}
	if err := ms.sendTask(task); err != nil {
		ms.logger.Infof("services error: %v", err)
	} else {
		ms.requestRepository.SetStatus(taskId, string(model.InProgress))
	}
	return task, nil
}

func (ms *ManagerService) GetRequest(id string) (*model.HashRequest, error) {
	return ms.requestRepository.GetRequestById(id)
}

func (ms *ManagerService) SetResultCrackHashTask(r *hash.CrackHashWorkerResponse) error {
	if r.Answers.Words == nil && len(r.Answers.Words) == 0 {
		return nil
	}
	req, err := ms.requestRepository.GetRequestById(r.RequestId)
	if err != nil {
		ms.logger.Infof("services error: %v", err)
		return err
	}
	// проверка на то что у нас не будет записано два одинаковых ответа
	if slices.Contains(req.Data, r.Answers.Words[0]) {
		return nil
	}

	err = ms.requestRepository.SetResults(r)
	if err != nil {
		ms.logger.Infof("services error: %v", err)
	}
	return nil
}

func (ms *ManagerService) sendTask(hashRequest *model.HashRequest) error {
	for i := 0; i < ms.partCount; i++ {
		task := toCrackHashManagerRequest(hashRequest, i+1, ms.partCount)

		xmlTask, err := xml.Marshal(&task)
		if err != nil {
			ms.logger.Infof("producer error: failed to marshal task due to error %v", err)
			return err
		}

		// Отправляем задачу в очередь
		err = ms.Rabbit.Producer.Channel.PublishWithContext(context.Background(),
			"",                            // exchange
			ms.Rabbit.Producer.Queue.Name, // routing key
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, // - чтобы сообщения выдержали перезапуск узла RabbitMQ нужно объявить их как долговечные
				ContentType:  "application/xml",
				Body:         xmlTask,
			})
		if err != nil {
			ms.logger.Infof("producer error: failed to publish a message due to error %v", err)
			return ErrorConnectionRefused
		}
		ms.logger.Infof("producer: send a message")
	}
	return nil
}

func (ms *ManagerService) Resend() {
	req, err := ms.requestRepository.GetRequestsByStatus(string(model.Created))
	if errors.Is(err, store.ErrorRecordNotFound) {
		return
	}
	for _, r := range req {
		err := ms.sendTask(r)
		if err == nil {
			ms.requestRepository.SetStatus(r.ID, string(model.InProgress))
		}
	}
}

func (ms *ManagerService) checkRequests() {
	// если запрос сделан больше 10 минут назад и данных все ещё нет, ставим статус error
	ticker := time.NewTicker(10 * time.Minute)
	for _ = range ticker.C {
		req, err := ms.requestRepository.GetRequestsByStatus(string(model.InProgress))
		if errors.Is(err, store.ErrorRecordNotFound) {
			return
		}
		for _, r := range req {
			if r.Data == nil && time.Now().Sub(r.DateTime) > time.Minute*10 {
				ms.requestRepository.SetStatus(r.ID, string(model.Error))
			}

		}
	}
}
