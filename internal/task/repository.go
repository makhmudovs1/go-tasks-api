package task

import "context"

type Repository interface {
	Create(ctx context.Context, task Task) (Task, error)
	GetByID(ctx context.Context, id int64) (Task, error)
	List(ctx context.Context, status *Status) ([]Task, error)
}
