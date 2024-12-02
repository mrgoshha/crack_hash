package services

import "worker/api/hash"

func toCrackHashWorkerResponse(res []string, r *hash.CrackHashManagerRequest) *hash.CrackHashWorkerResponse {
	return &hash.CrackHashWorkerResponse{
		RequestId:  r.RequestId,
		PartNumber: r.PartNumber,
		Answers: &hash.Words{
			Words: res,
		},
	}
}
