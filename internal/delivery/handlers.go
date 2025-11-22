package delivery

import (
	"LinksChecker/internal/delivery/dto"
	"LinksChecker/internal/pkg/pdf"
	"LinksChecker/internal/service/checker"
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
	"time"
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
		return
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

	var req dto.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.currentTasks.Add(1)
	defer h.currentTasks.Add(-1)

	tasks, err := h.checker.GetAll(req.LinksList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pdfReport := pdf.GenerateReport(tasks)

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=report.pdf")
	w.Write(pdfReport)
}

// WaitForActiveTasks awaiting completion of active tasks
func (h *Handler) WaitForActiveTasks(timeout int) bool {
	log.Printf("Waiting for completion of %d tasks...", h.currentTasks.Load())

	for range timeout {
		currentTasks := h.currentTasks.Load()
		if currentTasks == 0 {
			log.Println("All tasks completed")
			return true
		}

		log.Printf("Waiting for %d tasks...", currentTasks)
		time.Sleep(1 * time.Second)
	}

	log.Printf("Timeout, %d tasks is active", h.currentTasks.Load())
	return false
}

// SetReady sets the status of the service
func (h *Handler) SetReady(ready bool) {
	h.isReady.Store(ready)
}
