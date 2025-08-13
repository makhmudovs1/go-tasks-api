package service

import (
	"context"
	"errors"
	"github.com/makhmudovs1/go-tasks-api/internal/task"
	"strings"
)

type TaskService struct {
	repo task.Repository
}

func NewTaskService(repo task.Repository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

// error for incorrect status
var ErrInvalidTask = errors.New("invalid status")

// error if the title is empty
var ErrEmptyTitle = errors.New("title cannot be empty")

func isValidStatus(s task.Status) bool {
	switch s {
	case task.StatusTodo, task.StatusDone, task.StatusInProgress:
		return true
	default:
		return false
	}
}

func (s *TaskService) CreateTask(ctx context.Context, t task.Task) (task.Task, error) {
	if strings.TrimSpace(t.Title) == "" {
		return task.Task{}, ErrInvalidTask
	}
	if t.Status == "" {
		t.Status = task.StatusTodo
	}
	if !isValidStatus(t.Status) {
		return task.Task{}, ErrInvalidTask
	}

	return s.repo.Create(ctx, t)
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (task.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) ListTasks(ctx context.Context, statusStr *string) ([]task.Task, error) {
	var status *task.Status

	if statusStr != nil {
		st := task.Status(strings.ToLower(strings.TrimSpace(*statusStr)))
		if !isValidStatus(st) {
			return nil, ErrInvalidTask
		}
		status = &st
	}
	return s.repo.List(ctx, status)
}
