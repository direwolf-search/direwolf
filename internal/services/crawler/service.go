package crawler

import (
	"direwolf/internal/domain/service/crawler"
)

type service struct {
	Engine crawler.Engine // TODO: taskpool
}

func NewService(engine crawler.Engine) crawler.Crawler {
	return &service{
		Engine: engine,
	}
}

func (s *service) DoTasks() {
	// some engine actions here
}

func (s *service) GetTask() {
	//TODO implement me
	panic("implement me")
}
