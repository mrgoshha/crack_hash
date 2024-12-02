package worker

import (
	"encoding/json"
	"encoding/xml"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"worker/api/hash"
	"worker/internal/rabbit"
	"worker/internal/services"
)

type worker struct {
	router  *mux.Router
	rabbit  *rabbit.Rabbit
	logger  *logrus.Logger
	config  *Config
	service *services.WorkerService
}

func NewWorker(config *Config, rabbit *rabbit.Rabbit, service *services.WorkerService, logger *logrus.Logger) *worker {
	m := &worker{
		router:  mux.NewRouter(),
		logger:  logger,
		config:  config,
		rabbit:  rabbit,
		service: service,
	}
	return m
}

func (wr *worker) Start() error {
	wr.configureRouter()

	wr.logger.Infof("worker start")

	return http.ListenAndServe(wr.config.BaseUrl, wr.router)
}

func (wr *worker) configureRouter() {
	wr.router.Use(wr.logRequest)
	wr.router.HandleFunc("/hello-worker", wr.handleHello())
	wr.Serve()
}

func (wr *worker) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		wr.logger.Infof("[%s] %s", r.Method, r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

func (wr *worker) handleHello() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, "hello worker")
	}
}

func (wr *worker) Serve() {
	wr.ConsumeTask()
}

func (wr *worker) ConsumeTask() {
	//Мы не хотим потерять ни одной задачи. Если рабочий умирает, мы бы хотели, чтобы задача была передана другому рабочему.
	//Чтобы гарантировать, что сообщение никогда не будет потеряно,
	//RabbitMQ поддерживает message acknowledgments.
	//Подтверждение (неподтверждение) отправляется обратно потребителем, чтобы сообщить RabbitMQ,
	//что конкретное сообщение было получено, обработано и что RabbitMQ может удалить его.
	//Если пользователь умирает (его канал закрыт, соединение закрыто или TCP-соединение потеряно) без отправки подтверждения,
	//RabbitMQ поймет, что сообщение не было обработано полностью, и переведет его в очередь повторно.

	//Для этого нужно использовать ручные подтверждения сообщений,
	//передавая false в качестве аргумента "auto-ack",
	//а затем отправим соответствующее подтверждение от воркера с d.Ack(false), как только мы закончим с задачей
	tasks, err := wr.rabbit.Consumer.Channel.Consume(
		wr.rabbit.Consumer.Queue.Name, // queue
		"",                            // consumer
		false,                         // auto-ack
		false,                         // exclusive
		false,                         // no-local
		false,                         // no-wait
		nil,                           // args
	)
	if err != nil {
		wr.logger.Infof("consumer error: failed to register a consumer. Error: %v", err)
	}

	go func() {
		for d := range tasks {
			wr.logger.Infof("consumer: received a message")
			req := &hash.CrackHashManagerRequest{}

			if err := xml.Unmarshal(d.Body, &req); err != nil {
				wr.logger.Infof("consumer error: failed to unmarshal request. Error: %v", err)
				return
			}
			wr.service.CrackHash(req)
			// отправка подтверждения
			err := d.Ack(false)
			if err != nil {
				wr.logger.Infof("consumer error: failed to ack message. Error: %v", err)
			}
		}
		go wr.reconnectRabbitMq()
	}()
}

func (wr *worker) reconnectRabbitMq() {
	rabbit, _ := wr.rabbit.Reconnect()
	wr.rabbit = rabbit
	wr.service.Rabbit = rabbit
	wr.Serve()
}

func (wr *worker) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	wr.response(w, r, code, map[string]string{"error": err.Error()})
}

func (wr *worker) response(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
