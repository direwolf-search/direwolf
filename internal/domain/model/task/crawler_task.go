package task

import (
	"time"
)

type CrawlerTask struct {
	name          string
	dateCreated   time.Time
	dateCompleted time.Time
	body          []string
	schedule      string
	once          bool
	complete      bool
}

func NewCrawlerTask(n string, b []string, o bool) *CrawlerTask {
	return &CrawlerTask{
		name:        n,
		dateCreated: time.Now(),
		body:        b,
		once:        o,
	}
}

func (ct *CrawlerTask) Name() string {
	return ct.name
}

func (ct *CrawlerTask) Schedule() string {
	return ct.schedule
}

func (ct *CrawlerTask) Body() []string {
	return ct.body
}

func (ct *CrawlerTask) DateCreated() time.Time {
	return ct.dateCreated
}

func (ct *CrawlerTask) DateCompleted() time.Time {
	return ct.dateCompleted
}

func (ct *CrawlerTask) Once() bool {
	return ct.once
}

func (ct *CrawlerTask) Complete() bool {
	return ct.complete
}
