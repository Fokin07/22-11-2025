package models

import "sync"

// Task the task's model
type Task struct {
	ID    int
	Links map[string]string
	mu    sync.RWMutex
}
