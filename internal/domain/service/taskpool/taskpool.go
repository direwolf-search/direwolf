package taskpool

import (
	"direwolf/internal/domain/model/task"
)

type TaskPool interface {
	ScheduleTask(task *task.CrawlerTask, jobFunc func())
	FillTask(task *task.CrawlerTask) *task.CrawlerTask
	GetTaskList() []*task.CrawlerTask
	RemoveJob(jobID int)
	Start()
}
