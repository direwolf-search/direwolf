package crawler

import (
	"direwolf/internal/domain/repository"
	"direwolf/internal/domain/service/crawler"
)

type service struct {
	Engine     crawler.Engine
	Repository repository.CrawlerRepository
}

func NewService(e crawler.Engine, r repository.CrawlerRepository) crawler.Crawler {
	return &service{
		Engine:     e,
		Repository: r,
	}
}

func (s *service) Crawl(links []string) {}
