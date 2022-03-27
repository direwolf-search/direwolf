package task

type CrawlerTaskPool struct {
	tasks      []*CrawlerTask
	inProgress []string
}

func NewCrawlerTaskPool() *CrawlerTaskPool {
	return &CrawlerTaskPool{
		tasks:      make([]*CrawlerTask, 0),
		inProgress: make([]string, 0),
	}
}

func (p *CrawlerTaskPool) AddTask(t *CrawlerTask) {
	p.tasks = append(p.tasks, t)
}
