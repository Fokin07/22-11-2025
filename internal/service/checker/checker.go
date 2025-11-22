package checker

import (
	"LinksChecker/internal/models"
	"LinksChecker/internal/repository"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Service provides methods for check the links
type Service struct {
	repo   repository.Repo
	client *http.Client
}

// New creates a new checker service
func New(repo repository.Repo) *Service {
	return &Service{
		repo: repo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CheckLink checks link
func (s *Service) CheckLink(link string) string {
	if !strings.HasPrefix(link, "http") {
		link = "https://" + link
	}

	resp, err := s.client.Get(link)
	if err != nil {
		return "not available"
	}

	defer resp.Body.Close()

	if resp.StatusCode < 400 {
		return "available"
	}

	return "not available"
}

// CheckLinks checks links
func (s *Service) CheckLinks(links []string) (map[string]string, int, error) {
	if len(links) == 0 {
		return nil, 0, errors.New("No links provided")
	}

	result := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, link := range links {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			status := s.CheckLink(l)
			mu.Lock()
			result[l] = status
			mu.Unlock()
		}(link)
	}

	wg.Wait()

	id := s.repo.Save(result)

	return result, id, nil
}

// GetAll returns all tasks
func (s *Service) GetAll(ids []int) ([]*models.Task, error) {
	if len(ids) == 0 {
		return nil, errors.New("No links list provided")
	}

	return s.repo.GetAll(ids), nil
}
