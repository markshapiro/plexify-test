package utils

import (
	"math/rand/v2"
	"plexify-test/models"
	"time"
)

const (
	minJobDurationSeconds = 5
	maxJobDurationSeconds = 30
)

type JobProcessor interface {
	Process(job models.Job) error
}

type stringJobProcessor struct {
}

func NewStringJobProcessor() JobProcessor {
	return &stringJobProcessor{}
}

func (stringJobProcessor) Process(job models.Job) error {

	var randomTimeMS = (rand.Int64N(maxJobDurationSeconds-minJobDurationSeconds+1) + minJobDurationSeconds) * 1000
	start := makeTimestamp()
	for makeTimestamp()-start < randomTimeMS {
	}

	return nil
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}
