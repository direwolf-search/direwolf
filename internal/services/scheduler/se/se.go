package se

import (
	"github.com/robfig/cron/v3"
)

type schedulerLogger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(err error, msg string, keysAndValues ...interface{})
	Printf(format string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Warning(msg string, keysAndValues ...interface{})
	Critical(err error, msg string, keysAndValues ...interface{})
	Fatal(err error, msg string, keysAndValues ...interface{})
}

type ScedulerEngine struct {
	scheduler *cron.Cron
	logger    schedulerLogger
	taskList  map[int]int64
}

func NewCronEngine(l schedulerLogger) *ScedulerEngine {
	return &ScedulerEngine{
		scheduler: cron.New(cron.WithLogger(l)),
		logger:    l,
		taskList:  make(map[int]int64),
	}
}

func (se *ScedulerEngine) TaskList() map[int]int64 {
	return se.taskList
}

func (se *ScedulerEngine) Schedule(task map[string]interface{}, jobFunc func()) {
	var (
		wrapper  cron.JobWrapper
		schedule string
		id       int
	)

	// behaviour when next run is happens
	if v, ok := task["skip_next"]; ok {
		if boolVal, ok := v.(bool); ok && boolVal {
			wrapper = cron.SkipIfStillRunning(se.logger)
		} else {
			wrapper = cron.DelayIfStillRunning(se.logger)
		} // TODO!!!!! tests
	}

	funcJob := cron.FuncJob(jobFunc)

	if v, ok := task["schedule"]; ok {
		if stringVal, ok := v.(string); ok {
			schedule = stringVal
		}
	}

	if v, ok := task["id"]; ok {
		if intVal, ok := v.(int); ok {
			id = intVal
		}
	}

	// register job with its wrappers in cron
	if cronEntryID, err := se.scheduler.AddJob(
		schedule,
		cron.NewChain(
			wrapper,
			cron.Recover(se.logger),
		).Then(funcJob),
	); err != nil {
		se.logger.Critical(err, "cannot schedule task with id ", id)
	} else {
		se.logger.Info("Successfully scheduled task with id", id)
		se.taskList[int(cronEntryID)] = int64(id)
	}
}

func (se *ScedulerEngine) Remove(taskID int64) {
	var (
		jobID int
	)

	for jid, tid := range se.taskList {
		if tid == taskID {
			jobID = jid
			se.scheduler.Remove(cron.EntryID(jid))
		}
	}

	delete(se.taskList, jobID)
}

func (se *ScedulerEngine) Start() {
	se.scheduler.Start()
}
