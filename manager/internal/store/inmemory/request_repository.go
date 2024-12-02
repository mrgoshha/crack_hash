package inmemory

import (
	"github.com/google/uuid"
	"manager/api/hash"
	"manager/internal/model"
	"manager/internal/store"
	"sync"
	"time"
)

type RequestRepository struct {
	Requests sync.Map
}

func NewRequestRepository() store.RequestRepository {
	return &RequestRepository{}
}

func (rr *RequestRepository) Create(Hash string, MaxLength int) (string, error) {
	id := uuid.New().String()

	req := &model.HashRequest{
		ID:        id,
		Hash:      Hash,
		MaxLength: MaxLength,
		Data:      []string{},
		Status:    model.InProgress,
		DateTime:  time.Now(),
	}
	rr.Requests.Store(req.ID, req)
	return id, nil
}

func (rr *RequestRepository) GetRequestById(id string) (*model.HashRequest, error) {
	req, ok := rr.Requests.Load(id)
	if !ok {
		return nil, store.ErrorRecordNotFound
	}

	return req.(*model.HashRequest), nil
}

func (rr *RequestRepository) GetRequestsByStatus(status string) ([]*model.HashRequest, error) {
	//TODO implement me
	return nil, nil
}

func (rr *RequestRepository) SetStatus(id, status string) (*model.HashRequest, error) {
	req, _ := rr.Requests.Load(id)
	req.(*model.HashRequest).Status = model.Error
	return req.(*model.HashRequest), nil
}

func (rr *RequestRepository) SetResults(r *hash.CrackHashWorkerResponse) error {
	if len(r.Answers.Words) == 0 {
		return nil
	}
	result, err := rr.GetRequestById(r.RequestId)
	if err != nil {
		return err
	}

	result.Data = append(result.Data, r.Answers.Words...)

	_, _ = rr.SetStatus(r.RequestId, string(model.Ready))

	return nil

}
