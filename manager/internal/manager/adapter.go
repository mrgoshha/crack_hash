package manager

import (
	"manager/api/hash"
	"manager/internal/model"
)

func toResponseID(r *model.HashRequest) *hash.ResponseID {
	return &hash.ResponseID{
		RequestID: r.ID,
	}
}

func toResponseResult(r *model.HashRequest) *hash.ResponseResult {
	return &hash.ResponseResult{
		Status: string(r.Status),
		Data:   r.Data,
	}
}
