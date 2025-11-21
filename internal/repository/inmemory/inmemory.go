package inmemory

import (
	"LinksChecker/internal/models"
	"maps"
	"sync"
	"sync/atomic"
)

// Storage in-memory implementation of storage
type Storage struct {
	tasks map[int]*models.Task
	count int32
	mu    sync.RWMutex
}

// New init in-memory storage
func New() *Storage {
	return &Storage{
		tasks: make(map[int]*models.Task),
		count: 0,
	}
}

// Save saves a new link verification task
func (s *Storage) Save(links []string) int {
	taskID := int(atomic.AddInt32(&s.count, 1))

	linksMap := make(map[string]string)
	for _, link := range links {
		linksMap[link] = "checking"
	}

	s.mu.Lock()
	s.tasks[taskID] = &models.Task{ID: taskID, Links: linksMap}
	s.mu.Unlock()

	return taskID
}

// Get returns a task by ID
func (s *Storage) Get(id int) *models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tasks[id]
}

// Update updates the verification status
func (s *Storage) Update(taskID int, link string, status string) {
	s.mu.RLock()
	task, exists := s.tasks[taskID]
	s.mu.RUnlock()

	if exists {
		task.Mu.Lock()
		task.Links[link] = status
		task.Mu.Unlock()
	}
}

// GetAll returns a copy of all tasks
func (s *Storage) GetAll() map[int]*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasksCopy := make(map[int]*models.Task, len(s.tasks))
	for id, task := range s.tasks {
		task.Mu.RLock()
		links := make(map[string]string, len(task.Links))
		maps.Copy(links, task.Links)
		task.Mu.RUnlock()

		tasksCopy[id] = &models.Task{
			ID:    task.ID,
			Links: links,
		}
	}

	return tasksCopy
}
