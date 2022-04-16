// Package task defines a model of a certain task intended for a certain service
package task

import "fmt"

type typeOfTask struct {
	// service that owns the task
	of string
	// link selection rule
	rule string
}

func (t *typeOfTask) String() string {
	return fmt.Sprintf("%s.%s", t.of, t.rule)
}

// Task is a schedulable task for some service.
type Task struct {
	id       int64
	taskType *typeOfTask
	schedule string
	skipNext bool
}

// NewTask creates new task
func NewTask(taskType *typeOfTask, schedule string, skipNext bool) *Task {
	return &Task{
		taskType: taskType,
		schedule: schedule,
		skipNext: skipNext,
	}
}

// Of returns `of` field of the taskType field
func (ct *Task) Of() string {
	return ct.taskType.of
}

// Rule returns rule field of the taskType field
func (ct *Task) Rule() string {
	return ct.taskType.rule
}

// Schedule returns task's schedule field
func (ct *Task) Schedule() string {
	return ct.schedule
}

// SkipNext returns task's skipNext field
func (ct *Task) SkipNext() bool {
	return ct.skipNext
}

// ID returns task's id field
func (ct *Task) ID() int64 {
	return ct.id
}

/*
	Select [ALL] or [FIELD] [RELATION] [VALUE]
*/
