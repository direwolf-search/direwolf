// Package scheduler implements cron scheduling and tasks management
package scheduler

import (
	"context"
	scheduler2 "direwolf/internal/domain/repository/scheduler"

	"direwolf/internal/domain"
	"direwolf/internal/domain/model/task"
	"direwolf/internal/domain/service/scheduler"
)

type service struct {
	scheduler  scheduler.Engine
	of         string // other service maintained by the scheduler
	logger     domain.Logger
	repository scheduler2.Repository
}

// NewService creates new scheduler service
func NewService(se scheduler.Engine, of string, l domain.Logger, r scheduler2.Repository) scheduler.Scheduler {
	return &service{
		scheduler:  se,
		of:         of,
		logger:     l,
		repository: r,
	}
}

// Maintain sets name of service maintained by s
func (s *service) Maintain(serviceName string) {
	s.of = serviceName
}

// GetTasks gets a list of tasks to execute from the repository
func (s *service) GetTasks(ctx context.Context) ([]*task.Task, error) {
	return s.repository.ByType(ctx, s.of)
}

// ScheduleTask sets a task to run with its schedule
func (s *service) ScheduleTask(task *task.Task, jobFunc func()) {
	taskmap := task.Map()
	s.scheduler.Schedule(taskmap, jobFunc)
}

// GetScheduled returns list of scheduled tasks
func (s *service) GetScheduled(ctx context.Context) ([]*task.Task, error) {
	var (
		list  = s.scheduler.TaskList()
		tasks = make([]*task.Task, 0, len(list))
	)

	for _, taskID := range list {
		t, err := s.repository.ByID(ctx, taskID)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

// RemoveTask removes task with given ID from cron
func (s *service) RemoveTask(taskID int64) {
	s.scheduler.Remove(taskID)
}

// Start starts cron
func (s *service) Start() {
	s.scheduler.Start()
}
