package crawler

import (
	"direwolf/internal/domain/repository"
	"direwolf/internal/domain/service/crawler"
)

type service struct {
	Engine     crawler.Engine
	Repository repository.CrawlerRepository
}

// NewService creates new crawler service.
// both arguments are assumed to be fully initialized before being passed
func NewService(e crawler.Engine, r repository.CrawlerRepository) crawler.Crawler {
	return &service{
		Engine:     e,
		Repository: r,
	}
}

// Crawl makes crawling of links
func (s *service) Crawl(links []string) {}
