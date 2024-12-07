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
	time.Sleep(time.Duration(rand.IntN(maxJobDurationSeconds-minJobDurationSeconds+1)+minJobDurationSeconds) * time.Second)
	return nil
}
