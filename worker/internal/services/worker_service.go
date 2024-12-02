package services

import (
	"context"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	_hash "hash"
	"io"
	"sort"
	"worker/api/hash"
	"worker/internal/math/combinatorics"
	"worker/internal/rabbit"
)

type WorkerService struct {
	Rabbit *rabbit.Rabbit
	logger *logrus.Logger
}

func NewWorkerService(rabbit *rabbit.Rabbit, logger *logrus.Logger) *WorkerService {
	return &WorkerService{
		Rabbit: rabbit,
		logger: logger,
	}
}

func (ws *WorkerService) CrackHash(req *hash.CrackHashManagerRequest) {
	if req.PartNumber > req.PartCount || req.PartNumber == 0 {
		ws.logger.Infof("services error: %v", ErrorInvalidHashOrPartNumber)
	}

	sort.Strings(req.Alphabet.Symbols)
	hasher := md5.New()

	var basePermutation, lastPermutation string
	alphabet := req.Alphabet.String()
	var results []string
	pi := &combinatorics.PermutationsIterator{
		Alphabet:    alphabet,
		AlphabetLen: len(req.Alphabet.String()),
		Last:        alphabet[len(alphabet)-1],
	}
	for i := 1; i <= req.MaxLength; i++ {
		basePermutation, lastPermutation = ws.getBounds(req.Alphabet, req.PartNumber, req.PartCount, i)
		pi.PermutationLength = i
		pi.CurrentPermutation = []byte(basePermutation)
		pi.LastPermutation = []byte(lastPermutation)
		for {
			stringHash := ws.calculateMD5(hasher, string(pi.CurrentPermutation))
			if stringHash == req.Hash {
				results = append(results, string(pi.CurrentPermutation))
			}
			if string(pi.CurrentPermutation) == string(pi.LastPermutation) {
				break
			}
			pi.NextPermutation()
		}
	}
	if err := ws.SendResults(results, req); err != nil {
		ws.logger.Infof("services error: %v", err)
	}
}

func (ws *WorkerService) SendResults(results []string, r *hash.CrackHashManagerRequest) error {

	res := toCrackHashWorkerResponse(results, r)

	xmlTask, err := xml.Marshal(&res)
	if err != nil {
		ws.logger.Infof("producer error: failed to marshal task due to error %v", err)
		return err
	}

	// Отправляем задачу в очередь
	err = ws.Rabbit.Producer.Channel.PublishWithContext(context.Background(),
		"",                            // exchange
		ws.Rabbit.Producer.Queue.Name, // routing key
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // - чтобы сообщения выдержали перезапуск узла RabbitMQ нужно объявить их как долговечные
			ContentType:  "application/xml",
			Body:         xmlTask,
		})
	if err != nil {
		ws.logger.Infof("producer error: failed to to publish a message due to error %v", err)
		return err
	}
	ws.logger.Infof("producer: send a message")

	return nil
}

func (ws *WorkerService) getBounds(alphabet *hash.Symbols, pn, pc, boundStringLen int) (string, string) {
	alphabetLength := len(alphabet.Symbols)
	partSize := alphabetLength / pc
	extra := alphabetLength % pc

	leftBound := (pn-1)*partSize + min(pn-1, extra)
	rightBound := (pn)*partSize + min(pn, extra)

	var leftBoundString, rightBoundString string
	for i := 0; i < boundStringLen; i++ {
		leftBoundString += alphabet.Symbols[leftBound]
		rightBoundString += alphabet.Symbols[rightBound-1]
	}

	return leftBoundString, rightBoundString
}

func (ws *WorkerService) calculateMD5(hasher _hash.Hash, value string) string {
	hasher.Reset()
	io.WriteString(hasher, value)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
