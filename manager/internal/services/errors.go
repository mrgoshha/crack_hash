package services

import "errors"

var (
	ErrorConnectionRefused = errors.New("services error: rabbitMQ connection refused")
)
