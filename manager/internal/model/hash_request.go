package model

import "time"

const (
	Created    status = "CREATED"
	InProgress status = "IN_PROGRESS"
	Ready      status = "READY"
	Error      status = " ERROR"
)

type status string

type HashRequest struct {
	ID        string `bson:"_id,omitempty"`
	Hash      string
	MaxLength int
	Data      []string
	Status    status
	DateTime  time.Time
}
