package services

import (
	"manager/api/hash"
	"manager/internal/model"
)

func toCrackHashManagerRequest(r *model.HashRequest, partNumber int, partCount int) *hash.CrackHashManagerRequest {
	return &hash.CrackHashManagerRequest{
		RequestId:  r.ID,
		PartNumber: partNumber,
		PartCount:  partCount,
		Hash:       r.Hash,
		MaxLength:  r.MaxLength,
		Alphabet: &hash.Symbols{
			Symbols: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
				"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
				"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		},
	}
}
