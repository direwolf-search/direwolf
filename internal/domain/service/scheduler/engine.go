package scheduler

type Engine interface {
	Schedule(task map[string]interface{}, jobFunc func())
	Remove(taskID int64)
	Start()
	TaskList() map[int]int64
}
