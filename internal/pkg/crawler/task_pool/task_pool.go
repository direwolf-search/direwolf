package task_pool

import (
	"github.com/robfig/cron/v3"
)

type Cron cron.Cron

type Pool struct {
}

func (p *Pool) AddTask() (int, error) {
	return 0, nil
}
