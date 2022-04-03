package crawler_task_pool

import (
	"github.com/robfig/cron/v3"

	"direwolf/internal/domain/model/task"
	"direwolf/internal/domain/service/crawler"
)

type service struct {
	scheduler cron.Cron
}

func NewService() crawler.TaskPool {
	return &service{}
}

func (s *service) AddTask(task *task.CrawlerTask) {

}

func (s *service) AddJobForTask(task *task.CrawlerTask, links []string) {

}

func (s *service) GetTaskList() []*task.CrawlerTask {
	//TODO implement me
	panic("implement me")
}

func (s *service) RemoveJob(jobID int) {
	//TODO implement me
	panic("implement me")
}

func (s *service) Start() {
	//TODO implement me
	panic("implement me")
}
