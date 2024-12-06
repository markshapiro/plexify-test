package app

import "errors"

const (
	Pending    = "pending"
	Processing = "processing"
	Completed  = "completed"

	numWorkers = 5
)

var (
	ErrNotFound  = errors.New("NOT FOUND")
	ErrQueueFull = errors.New("QUEUE IS FULL")
)
