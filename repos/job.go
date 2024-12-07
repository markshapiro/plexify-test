package repos

import (
	"errors"
	"sync"
	"sync/atomic"

	"plexify-test/models"
)

var (
	ErrNotFound = errors.New("NOT FOUND")
)

type JobRepo interface {
	AddNew() int64
	DeleteJob(key int64)
	SetStatus(key int64, status string) error
	GetStatus(key int64) (string, error)
}

type jobRepo struct {
	statuses *sync.Map
	nextID   *int64
}

func NewJobRepo() JobRepo {
	var nextID int64
	return &jobRepo{
		statuses: &sync.Map{},
		nextID:   &nextID,
	}
}

func (repo jobRepo) AddNew() int64 {
	newID := atomic.AddInt64(repo.nextID, 1)
	repo.statuses.Store(newID, models.Pending)
	return newID
}

func (repo jobRepo) DeleteJob(key int64) {
	repo.statuses.Delete(key)
}

func (repo jobRepo) SetStatus(key int64, status string) error {

	// waning: not atomic
	_, ok := repo.statuses.Load(key)
	if !ok {
		return ErrNotFound
	}

	repo.statuses.Store(key, status)

	return nil
}

func (repo jobRepo) GetStatus(key int64) (string, error) {
	result, ok := repo.statuses.Load(key)
	if !ok {
		return "", ErrNotFound
	}
	return result.(string), nil
}
