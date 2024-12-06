package app

import (
	"fmt"
	"sync"
	"time"
)

var (
	jobRepo   = NewJobRepo()
	jobChan   = make(chan JobTask, 1000)
	processor JobProcessor
	wg        sync.WaitGroup
)

func Start() {

	processor = NewStringJobProcessor()

	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go worker(&wg, jobChan)
	}
}

func Stop() {
	close(jobChan)
	wg.Wait()
}

func worker(wg *sync.WaitGroup, jobs <-chan JobTask) {
	defer wg.Done()
	for job := range jobs {

		processor.Process(Job{job.Payload})

		err := jobRepo.setStatus(job.JobID, Completed)
		if err != nil {
			fmt.Println("tried to set status of mom existent job")
		}
	}
}

func JobCreate(newJob JobCreateDto) (JobIDDto, error) {

	newID := jobRepo.addNew()

	select {
	case jobChan <- JobTask{newID, newJob.Payload}:
	case <-time.After(2 * time.Second):
		return JobIDDto{}, ErrQueueFull
	}

	return JobIDDto{newID}, nil
}

func GetJobStatus(id int64) (JobStatusDto, error) {
	s, err := jobRepo.getStatus(id)
	if err != nil {
		return JobStatusDto{}, err
	}
	return JobStatusDto{id, s}, nil
}
