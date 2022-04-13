package scheduler

import (
	"direwolf/internal/domain/model/task"
)

type Scheduler interface {
	ScheduleTask(task *task.Task, jobFunc func())
	GetTaskList(taskType string) []*task.Task
	RemoveJob(jobID int)
	Start()
}
