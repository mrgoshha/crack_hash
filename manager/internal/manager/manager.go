package manager

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"manager/api/hash"
	"manager/internal/rabbit"
	"manager/internal/services"
	"net/http"
)

type manager struct {
	router  *mux.Router
	rabbit  *rabbit.Rabbit
	logger  *logrus.Logger
	config  *Config
	service *services.ManagerService
}

func NewManager(config *Config, rabbit *rabbit.Rabbit, service *services.ManagerService, logger *logrus.Logger) *manager {
	m := &manager{
		router:  mux.NewRouter(),
		logger:  logger,
		config:  config,
		rabbit:  rabbit,
		service: service,
	}
	return m
}

func (m *manager) Start() error {
	m.configureRouter()

	m.logger.Infof("manager start")

	return http.ListenAndServe(m.config.BaseUrl, m.router)
}

func (m *manager) configureRouter() {
	m.router.Use(m.logRequest)
	m.router.NotFoundHandler = http.HandlerFunc(m.notFoundHandler)
	m.router.HandleFunc("/hello-manager", m.handleHello())
	m.router.HandleFunc("/api/hash/crack", m.createCrackRequest()).Methods(http.MethodPost)
	m.router.HandleFunc("/api/hash/status", m.getResultCrackRequest()).Methods(http.MethodGet)
	m.Serve()
}

func (m *manager) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		m.logger.Infof("[%s] %s", r.Method, r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

func (m *manager) Serve() {
	m.ConsumeResults()
}

func (m *manager) handleHello() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, "hello manager")
	}
}

func (m *manager) createCrackRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &hash.CrackHashRequest{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			m.error(w, r, http.StatusBadRequest, err)
			return
		}

		hashResult, err := m.service.CreateCrackHashTask(req)
		if err != nil {
			m.error(w, r, http.StatusBadRequest, err)
		}

		res := toResponseID(hashResult)

		m.response(w, r, http.StatusCreated, res)
	}
}

func (m *manager) getResultCrackRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("requestId")

		statusRequest, err := m.service.GetRequest(id)
		if err != nil {
			m.error(w, r, http.StatusNotFound, err)
			return
		}

		res := toResponseResult(statusRequest)

		m.response(w, r, http.StatusOK, res)
	}
}

func (m *manager) ConsumeResults() {
	tasks, err := m.rabbit.Consumer.Channel.Consume(
		m.rabbit.Consumer.Queue.Name, // queue
		"",                           // consumer
		true,                         // auto-ack
		false,                        // exclusive
		false,                        // no-local
		false,                        // no-wait
		nil,                          // args
	)
	if err != nil {
		m.logger.Infof("consumer error: failed to register a consumer. Error: %v", err)
	}

	go func() {
		for d := range tasks {
			m.logger.Infof("consumer: received a message")
			req := &hash.CrackHashWorkerResponse{}

			if err := xml.Unmarshal(d.Body, &req); err != nil {
				m.logger.Infof("consumer error: failed to unmarshal request. Error: %v", err)
				return
			}
			m.service.SetResultCrackHashTask(req)
		}
		go m.reconnectRabbitMq()
	}()
}

func (m *manager) reconnectRabbitMq() {
	rabbit, _ := m.rabbit.Reconnect()
	m.rabbit = rabbit
	m.service.Rabbit = rabbit
	m.service.Resend()
	m.Serve()
}

func (m *manager) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	m.error(w, r, http.StatusNotFound, fmt.Errorf("not Found %s", r.URL))
}

func (m *manager) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	m.response(w, r, code, map[string]string{"error": err.Error()})
}

func (m *manager) response(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
