package sheduler

import (
	"github.com/robfig/cron/v3"

	"direwolf/internal/domain/model/task"
	"direwolf/internal/domain/service/scheduler"
)

type service struct {
	scheduler cron.Cron
}

func NewService() scheduler.Scheduler {
	return &service{}
}

func (s *service) ScheduleTask(task *task.CrawlerTask, jobFunc func()) {

}

func (s *service) FillTask(task *task.CrawlerTask) *task.CrawlerTask {
	return nil
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
