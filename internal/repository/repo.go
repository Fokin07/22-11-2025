package repository

import "LinksChecker/internal/models"

// Repo interface for working with storage
type Repo interface {
	Save(links map[string]string) int
	Get(id int) (*models.Task, error)
	GetAll(ids []int) []*models.Task
	SaveState()
	RestoreState()
}
