package store

import (
	"manager/api/hash"
	"manager/internal/model"
)

type RequestRepository interface {
	Create(Hash string, MaxLength int) (string, error)
	GetRequestById(id string) (*model.HashRequest, error)
	SetStatus(id, status string) (*model.HashRequest, error)
	SetResults(r *hash.CrackHashWorkerResponse) error
	GetRequestsByStatus(status string) ([]*model.HashRequest, error)
}
