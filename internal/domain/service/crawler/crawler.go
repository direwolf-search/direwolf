package crawler

import (
	"direwolf/internal/domain/model/task"
)

type Crawler interface {
	DoTasks()
	GetTask() *task.CrawlerTask
}
