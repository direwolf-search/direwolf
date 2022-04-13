package sheduler

import (
	"github.com/robfig/cron/v3"

	"direwolf/internal/domain/model/task"
	"direwolf/internal/domain/service/scheduler"
)

type service struct {
	scheduler cron.Cron
	//service maintained by the scheduler
	of string
}

func NewService(of string) scheduler.Scheduler {
	return &service{
		of: of,
	}
}

func (s *service) ScheduleTask(task *task.Task, jobFunc func()) {

}

func (s *service) GetTaskList() []*task.Task {
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
