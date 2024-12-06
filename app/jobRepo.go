package app

import (
	"sync"
	"sync/atomic"
)

type JobRepo struct {
	statuses sync.Map
	nextID   int64
}

func NewJobRepo() *JobRepo {
	return &JobRepo{}
}

func (repo *JobRepo) addNew() int64 {
	newID := atomic.AddInt64(&repo.nextID, 1)
	repo.statuses.Store(newID, Pending)
	return newID
}

func (repo *JobRepo) setStatus(key int64, status string) error {

	// waning: not atomic
	_, ok := repo.statuses.Load(key)
	if !ok {
		return ErrNotFound
	}

	repo.statuses.Store(key, status)

	return nil
}

func (repo *JobRepo) getStatus(key int64) (string, error) {
	result, ok := repo.statuses.Load(key)
	if !ok {
		return "", ErrNotFound
	}
	return result.(string), nil
}
