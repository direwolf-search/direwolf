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
	skipNext      bool
}

func NewCrawlerTask(name, schedule string, skipNext bool) *CrawlerTask {
	return &CrawlerTask{
		name:        name,
		dateCreated: time.Now(),
		schedule:    schedule,
		skipNext:    skipNext,
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
