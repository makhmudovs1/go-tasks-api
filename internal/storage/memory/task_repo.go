package memory

import (
	"context"
	"errors"
	"github.com/makhmudovs1/go-tasks-api/internal/task"
	"sync"
	"sync/atomic"
	"time"
)

type TaskRepo struct {
	mu   sync.RWMutex
	data map[int64]task.Task
	seq  int64
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{
		data: make(map[int64]task.Task),
	}
}

func (r *TaskRepo) Create(ctx context.Context, t task.Task) (task.Task, error) {
	id := atomic.AddInt64(&r.seq, 1)

	t.ID = id
	t.CreatedAt = time.Now().UTC()
	t.UpdatedAt = t.CreatedAt

	r.mu.Lock()
	r.data[id] = t
	r.mu.Unlock()

	return t, nil
}

var ErrNotFound = errors.New("task not found")

func (r *TaskRepo) GetByID(ctx context.Context, id int64) (task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res, ok := r.data[id]
	if !ok {
		return task.Task{}, ErrNotFound
	}
	return res, nil
}

func (r *TaskRepo) List(ctx context.Context, status *task.Status) ([]task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]task.Task, 0, len(r.data))
	for _, t := range r.data {
		if status != nil && t.Status != *status {
			continue
		}
		res = append(res, t)
	}
	return res, nil
}
