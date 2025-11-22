package inmemory

import (
	"LinksChecker/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

// Storage in-memory implementation of storage
type Storage struct {
	tasks    map[int]*models.Task
	count    int
	saveFile string
	mu       sync.RWMutex
}

// New init in-memory storage
func New(fileName string) *Storage {
	return &Storage{
		tasks:    make(map[int]*models.Task),
		count:    0,
		saveFile: fileName,
	}
}

// Save saves a new link verification task
func (s *Storage) Save(links map[string]string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.count++
	s.tasks[s.count] = &models.Task{ID: s.count, Links: links}

	return s.count
}

// Get returns a task by ID
func (s *Storage) Get(id int) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, errors.New("the task does not exist")
	}

	return task, nil
}

// GetAll returns all tasks
func (s *Storage) GetAll(ids []int) []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*models.Task
	for _, id := range ids {
		if task, exists := s.tasks[id]; exists {
			result = append(result, task)
		}
	}

	return result
}

// SaveState saves the state only when called shutdown/restart
func (s *Storage) SaveState() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := struct {
		Tasks  map[int]*models.Task `json:"tasks"`
		NextID int                  `json:"next_id"`
	}{
		Tasks:  s.tasks,
		NextID: s.count,
	}

	data, err := json.Marshal(state)
	if err != nil {
		fmt.Printf("Error marshaling state: %v\n", err)
		return
	}

	if err := os.WriteFile(s.saveFile, data, 0644); err != nil {
		fmt.Printf("Error writing state file: %v\n", err)
	}

	fmt.Printf("State saved: %d tasks, next ID: %d\n", len(s.tasks), s.count)
}

// RestoreState restores the state after shutdown/restart
func (s *Storage) RestoreState() {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.saveFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No previous state found, starting fresh")
			return
		}
		fmt.Printf("Error reading state file: %v\n", err)
		return
	}

	var state struct {
		Tasks  map[int]*models.Task `json:"tasks"`
		NextID int                  `json:"next_id"`
	}

	if err := json.Unmarshal(data, &state); err != nil {
		fmt.Printf("Error unmarshaling state: %v\n", err)
		return
	}

	s.tasks = state.Tasks
	s.count = state.NextID

	fmt.Printf("State restored: %d tasks, next ID: %d\n", len(s.tasks), s.count)
}

// Cleanup deletes the temporary status file
func (s *Storage) Cleanup() {
	if err := os.Remove(s.saveFile); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error cleaning up state file: %v\n", err)
	} else {
		fmt.Println("State file cleaned up")
	}
}
