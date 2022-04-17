// Package scheduler implements cron scheduling and tasks management
package scheduler

import (
	"context"

	"github.com/robfig/cron/v3"

	"direwolf/internal/domain"
	"direwolf/internal/domain/model/task"
	"direwolf/internal/domain/repository"
	"direwolf/internal/domain/service/scheduler"
)

type service struct {
	scheduler  *cron.Cron
	of         string // other service maintained by the scheduler
	logger     domain.Logger
	taskList   map[int]int64
	repository repository.SchedulerRepository
}

// NewService creates new scheduler service
func NewService(of string, l domain.Logger, r repository.SchedulerRepository) scheduler.Scheduler {
	return &service{
		scheduler:  cron.New(cron.WithLogger(l)),
		taskList:   make(map[int]int64),
		of:         of,
		logger:     l,
		repository: r,
	}
}

// GetTasks gets a list of tasks to execute from the repository
func (s *service) GetTasks(ctx context.Context) ([]*task.Task, error) {
	return s.repository.ByType(ctx, s.of)
}

// ScheduleTask sets a task to run with its schedule
func (s *service) ScheduleTask(task *task.Task, jobFunc func()) {
	var (
		wrapper cron.JobWrapper
	)

	// behaviour when next run is happens
	if task.SkipNext() {
		wrapper = cron.SkipIfStillRunning(s.logger)
	} else {
		wrapper = cron.DelayIfStillRunning(s.logger)
	}

	funcJob := cron.FuncJob(jobFunc)

	// register job with its wrappers in cron
	if cronEntryID, err := s.scheduler.AddJob(
		task.Schedule(),
		cron.NewChain(
			wrapper,
			cron.Recover(s.logger),
		).Then(funcJob),
	); err != nil {
		s.logger.Critical(err, "cannot schedule task with id ", task.ID())
	} else {
		s.logger.Info("Successfully scheduled task with id", task.ID())
		s.taskList[int(cronEntryID)] = task.ID()
	}
}

// GetScheduled returns list of scheduled tasks
func (s *service) GetScheduled(ctx context.Context) ([]*task.Task, error) {
	var (
		tasks = make([]*task.Task, 0, len(s.taskList))
	)

	for _, taskID := range s.taskList {
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
	var (
		jobID int
	)

	for jid, tid := range s.taskList {
		if tid == taskID {
			jobID = jid
			s.scheduler.Remove(cron.EntryID(jid))
		}
	}

	delete(s.taskList, jobID)
}

// Start starts cron
func (s *service) Start() {
	s.scheduler.Start()
}
