package crawler

import (
	"direwolf/internal/domain/model/task"
	"direwolf/internal/domain/service/crawler"
)

type service struct {
	Engine crawler.Engine
}

func NewService(engine crawler.Engine) crawler.Crawler {
	return &service{
		Engine: engine,
	}
}

func (s *service) DoTasks() {}

func (s *service) GetTask() *task.CrawlerTask {
	return nil
}
