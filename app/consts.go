package app

import "errors"

const (
	Pending    = "Pending"
	Processing = "processing"
	Completed  = "completed"

	numWorkers = 3
)

var (
	ErrNotFound  = errors.New("NOT FOUND")
	ErrQueueFull = errors.New("QUEUE IS FULL")
)
