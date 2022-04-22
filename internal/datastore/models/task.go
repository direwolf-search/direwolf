package models

import (
	"github.com/uptrace/bun"

	"direwolf/internal/domain/model/task"
)

// Task is a scheduled task for some service.
type Task struct {
	bun.BaseModel `bun:"links"`
	ID            int64  `bun:"id"`
	Of            string `bun:"of"`
	Rule          string `bun:"rule"`
	Schedule      string `bun:"schedule"`
	SkipNext      bool   `bun:"skip_next"`
}

func NewTaskFromModel(t *task.Task) *Task {
	return &Task{
		ID:       t.ID(),
		Of:       t.Of(),
		Rule:     t.Rule(),
		Schedule: t.Schedule(),
		SkipNext: t.SkipNext(),
	}
}

func (t *Task) ToModel() *task.Task {
	return task.NewTask(t.ID, t.Of, t.Rule, t.Schedule, t.SkipNext)
}
