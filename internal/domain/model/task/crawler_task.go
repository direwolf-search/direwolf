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
	skipNew       bool
}

func NewCrawlerTask(name, schedule string, skipNew bool) *CrawlerTask {
	return &CrawlerTask{
		name:        name,
		dateCreated: time.Now(),
		schedule:    schedule,
		skipNew:     skipNew,
	}
}

func (ct *CrawlerTask) FillBody(links []string) error {
	ct.body = links

	return nil
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
