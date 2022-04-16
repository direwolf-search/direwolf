package repository

import (
	"context"
	"direwolf/internal/domain/model/task"
)

type SchedulerRepository interface {
	ByType(ctx context.Context, taskType string) ([]*task.Task, error)
	ByID(ctx context.Context, id int64) (*task.Task, error)
}
