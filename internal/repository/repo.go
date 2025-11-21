package repository

import "LinksChecker/internal/models"

// Repo interface for working with storage
type Repo interface {
	Save(links []string) int
	Get(id int) *models.Task
	Update(taskID int, link string, status string)
	GetAll() map[int]*models.Task
}
