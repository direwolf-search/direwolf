package crawler

import (
	"direwolf/internal/domain/model/task"
)

type TaskPool interface {
	AddTask(task *task.CrawlerTask)
	AddJobForTask(task *task.CrawlerTask, links []string)
	GetTaskList() []*task.CrawlerTask
	RemoveJob(jobID int)
	Start()
}
