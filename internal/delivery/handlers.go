package delivery

import (
	"LinksChecker/internal/delivery/dto"
	"LinksChecker/internal/service/checker"
	"encoding/json"
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

func (h *Handler) writeData(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CheckLinks processes requests to check the availability of links
func (h *Handler) CheckLinks(w http.ResponseWriter, r *http.Request) {
	if !h.isReady.Load().(bool) {
		http.Error(w, "Service is shutting down", http.StatusServiceUnavailable)
		return
	}

	var req dto.CheckLinksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.currentTasks.Add(1)
	defer h.currentTasks.Add(-1)

	links, id, err := h.checker.CheckLinks(req.Links)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeData(w, dto.CheckLinksResponse{Links: links, LinksNum: id}, http.StatusOK)
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
