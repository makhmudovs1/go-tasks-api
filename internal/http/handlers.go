package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/makhmudovs1/go-tasks-api/internal/logging"
	"github.com/makhmudovs1/go-tasks-api/internal/service"
	"github.com/makhmudovs1/go-tasks-api/internal/task"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	svc     *service.TaskService
	logChan chan<- logging.LogEvent
}

func NewHandler(svc *service.TaskService, logChan chan<- logging.LogEvent) *Handler {
	return &Handler{
		svc:     svc,
		logChan: logChan,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /tasks", h.listTasks)
	mux.HandleFunc("POST /tasks", h.createTask)

	mux.HandleFunc("GET /tasks/{id}", h.handleTaskByID)
}

// To not duplicate in every handler

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type errResp struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errResp{Error: msg})
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	var statusPtr *string
	if s := strings.TrimSpace(r.URL.Query().Get("status")); s != "" {
		statusPtr = &s
	}

	t, err := h.svc.ListTasks(r.Context(), statusPtr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	select {
	case h.logChan <- logging.LogEvent{
		Time:    time.Now(),
		Action:  "LIST_TASKS",
		Details: fmt.Sprintf("Count=%d, StatusFilter=%q", len(t), *statusPtr),
	}:
	default:
	}
	writeJSON(w, http.StatusOK, t)
}

type createTaskReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req createTaskReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	t := task.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      task.Status(strings.ToLower(strings.TrimSpace(req.Status))),
	}
	t, err := h.svc.CreateTask(r.Context(), t)
	if err != nil {
		writeError(w, http.StatusBadRequest, "internal error")
		return
	}
	select {
	case h.logChan <- logging.LogEvent{
		Time:    time.Now(),
		Action:  "CREATE_TASK",
		Details: fmt.Sprintf("Task ID=%d, Title=%q", t.ID, t.Title),
	}:
	default:
	}

	writeJSON(w, http.StatusCreated, t)
}

func (h *Handler) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	t, err := h.svc.GetTask(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	select {
	case h.logChan <- logging.LogEvent{
		Time:    time.Now(),
		Action:  "GET_TASK",
		Details: fmt.Sprintf("Task ID=%d", id),
	}:
	default:
	}

	writeJSON(w, http.StatusOK, t)
}
