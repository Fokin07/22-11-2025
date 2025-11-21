package inmemory

import (
	"LinksChecker/internal/models"
	"sync"
)

// Storage in-memory implementation of storage
type Storage struct {
	tasks  map[int]*models.Task
	taskID int32
	mu     sync.RWMutex
}

// New init in-memory storage
func New() *Storage {
	return &Storage{
		tasks:  make(map[int]*models.Task),
		taskID: 0,
	}
}
