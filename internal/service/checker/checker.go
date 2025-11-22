package checker

import (
	"LinksChecker/internal/repository"
)

// Service provides methods for check the links
type Service struct {
	repo repository.Repo
}

// New creates a new checker service
func New(repo repository.Repo) *Service {
	return &Service{
		repo: repo,
	}
}
