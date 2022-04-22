// Package task defines a model of a certain task intended for a certain service
package task

import (
	"errors"
	"fmt"
)

var ErrInvalidTask = errors.New("error of invalid task")

type typeOfTask struct {
	// service that owns the task
	of string
	// rule is an execution rule
	rule string
}

func (t *typeOfTask) String() string {
	return fmt.Sprintf("%s.%s", t.of, t.rule)
}

// Task is a scheduled task for some service.
type Task struct {
	id       int64
	taskType *typeOfTask
	schedule string
	skipNext bool
}

// NewTask creates new *Task
func NewTask(id int64, of, rule, schedule string, skipNext bool) *Task {
	return &Task{
		id:       id,
		taskType: &typeOfTask{of: of, rule: rule},
		schedule: schedule,
		skipNext: skipNext,
	}
}

// NewTaskFromMap creates new *Task from map[string]interface{}
func NewTaskFromMap(m map[string]interface{}) (*Task, error) {
	var (
		of, rule, schedule string
		skipNext           bool
	)

	if v, ok := m["of"]; ok {
		if stringVal, ok := v.(string); ok {
			of = stringVal
		}
	}

	if v, ok := m["rule"]; ok {
		if stringVal, ok := v.(string); ok {
			rule = stringVal
		}
	}

	if v, ok := m["schedule"]; ok {
		if stringVal, ok := v.(string); ok {
			schedule = stringVal
		}
	}

	if v, ok := m["skip_next"]; ok {
		if boolVal, ok := v.(bool); ok {
			skipNext = boolVal
		}
	}

	t := &Task{
		taskType: &typeOfTask{
			of:   of,
			rule: rule,
		},
		schedule: schedule,
		skipNext: skipNext,
	}

	err := t.Validate()
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Map creates map[string]interface{} from *Task
func (t *Task) Map() map[string]interface{} {
	var (
		m = make(map[string]interface{})
		//t = make(map[string]interface{})
	)

	if t.id != 0 {
		m["id"] = t.id
	}
	if t.taskType != nil {
		m["of"] = t.taskType.of
		m["rule"] = t.taskType.rule
	}
	if t.schedule != "" {
		m["schedule"] = t.schedule
	}
	if !t.skipNext {
		m["skip_next"] = t.skipNext
	}

	return m
}

// Of returns `of` field of the taskType field
func (t *Task) Of() string {
	return t.taskType.of
}

// Rule returns rule field of the taskType field
func (t *Task) Rule() string {
	return t.taskType.rule
}

// Schedule returns task's schedule field
func (t *Task) Schedule() string {
	return t.schedule
}

// SkipNext returns task's skipNext field
func (t *Task) SkipNext() bool {
	return t.skipNext
}

// ID returns task's id field
func (t *Task) ID() int64 {
	return t.id
}

func (t *Task) Validate() error {
	if t.schedule == "" {
		return ErrInvalidTask
	}

	return nil
}

/*
	Select [ALL] or [FIELD] [RELATION] [VALUE]
*/
