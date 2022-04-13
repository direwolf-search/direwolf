package crawler

import (
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

func (s *service) Crawl(links []string) {}
