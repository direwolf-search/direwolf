package scheduler

import (
	"context"

	"direwolf/internal/domain/model/task"
)

type Scheduler interface {
	ScheduleTask(task *task.Task, jobFunc func())
	GetScheduled(ctx context.Context) ([]*task.Task, error)
	RemoveTask(taskID int64)
	GetTasks(ctx context.Context) ([]*task.Task, error)
	Start()
}
