package delivery

import (
	"LinksChecker/internal/service/checker"
	"net/http"
	"sync/atomic"
)

// Handler provides HTTP handlers
type Handler struct {
	checker      *checker.Service
	currentTasks *atomic.Int32
	isReady      *atomic.Value
}

// NewHandler creates a new Handler
func NewHandler(checker *checker.Service) *Handler {
	isReady := &atomic.Value{}
	isReady.Store(true)

	currentTasks := &atomic.Int32{}
	currentTasks.Store(0)

	return &Handler{
		checker:      checker,
		isReady:      isReady,
		currentTasks: currentTasks,
	}
}

// CheckLinks processes requests to check the availability of links
func (h *Handler) CheckLinks(w http.ResponseWriter, r *http.Request) {
	if !h.isReady.Load().(bool) {
		http.Error(w, "Service is shutting down", http.StatusServiceUnavailable)
		return
	}

	h.currentTasks.Add(1)
	defer h.currentTasks.Add(-1)
}

// GenerateReport processes requests to generate a PDF report
func (h *Handler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	if !h.isReady.Load().(bool) {
		http.Error(w, "Service is shutting down", http.StatusServiceUnavailable)
		return
	}

	h.currentTasks.Add(1)
	defer h.currentTasks.Add(-1)
}

// WaitForActiveTasks awaiting completion of active tasks
func (h *Handler) WaitForActiveTasks(timeout int) bool {

	return false
}

// SetReady sets the status of the service
func (h *Handler) SetReady(ready bool) {
	h.isReady.Store(ready)
}
