package services

import (
	"errors"
	"fmt"
	"plexify-test/models"
	"plexify-test/repos"
	"plexify-test/utils"
	"sync"
	"time"
)

var (
	ErrNotFound  = errors.New("NOT FOUND")
	ErrQueueFull = errors.New("QUEUE IS FULL")
)

const (
	numWorkers     = 5
	chanBufferSize = 1000
)

type JobService interface {
	GetJobStatus(id int64) (models.JobStatusDto, error)
	JobCreate(newJob models.JobCreateDto) (models.JobIDDto, error)
	StartWorkers()
	StopWorkers()
}

type jobService struct {
	jobRepo      repos.JobRepo
	jobProcessor utils.JobProcessor
	jobChan      chan models.JobTask
	wg           *sync.WaitGroup
}

func NewJobService(jobRepo repos.JobRepo, jobProcessor utils.JobProcessor) JobService {
	return jobService{
		jobRepo:      jobRepo,
		jobProcessor: jobProcessor,
		jobChan:      make(chan models.JobTask, chanBufferSize),
		wg:           &sync.WaitGroup{},
	}
}

func (s jobService) StartWorkers() {
	s.wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go s.worker(s.wg, s.jobChan)
	}
}

func (s jobService) worker(wg *sync.WaitGroup, jobs <-chan models.JobTask) {
	defer wg.Done()
	for job := range jobs {

		err := s.jobProcessor.Process(models.Job{job.Payload})
		if err != nil {
			fmt.Println("job failed:", err.Error())
		}

		err = s.jobRepo.SetStatus(job.JobID, models.Completed)
		if err != nil {
			fmt.Println("tried to set status of mom existent job")
		}
	}
}

func (s jobService) StopWorkers() {
	close(s.jobChan)
	s.wg.Wait()
}

func (s jobService) JobCreate(newJob models.JobCreateDto) (models.JobIDDto, error) {
	newID := s.jobRepo.AddNew()

	select {
	case s.jobChan <- models.JobTask{newID, newJob.Payload}:
	case <-time.After(2 * time.Second):
		// if couldnt add to queue because buffer full

		s.jobRepo.DeleteJob(newID)

		return models.JobIDDto{}, ErrQueueFull
	}

	return models.JobIDDto{newID}, nil
}

func (s jobService) GetJobStatus(id int64) (models.JobStatusDto, error) {
	status, err := s.jobRepo.GetStatus(id)
	if err != nil {
		if err == repos.ErrNotFound {
			return models.JobStatusDto{}, ErrNotFound
		}
		return models.JobStatusDto{}, err
	}
	return models.JobStatusDto{id, status}, nil
}
