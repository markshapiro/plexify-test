package app

import "errors"

const (
	Pending    = "pending"
	Processing = "processing"
	Completed  = "completed"

	numWorkers = 5

	ChanBufferSize = 1000

	minJobDurationSeconds = 5
	maxJobDurationSeconds = 30
)

var (
	ErrNotFound  = errors.New("NOT FOUND")
	ErrQueueFull = errors.New("QUEUE IS FULL")
)
