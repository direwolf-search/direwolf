package repository

import (
	"direwolf/internal/domain/model/task"
)

type SchedulerRepository interface {
	GetTasks(taskType string) []*task.Task
}
