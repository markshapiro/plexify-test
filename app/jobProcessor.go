package app

import (
	"math/rand/v2"
	"time"
)

type JobProcessor interface {
	Process(job Job) error
}

type StringJobProcessor struct {
}

func NewStringJobProcessor() JobProcessor {
	return &StringJobProcessor{}
}

func (StringJobProcessor) Process(job Job) error {
	time.Sleep(time.Duration(rand.IntN(maxJobDurationSeconds-minJobDurationSeconds+1)+minJobDurationSeconds) * time.Second)
	return nil
}
