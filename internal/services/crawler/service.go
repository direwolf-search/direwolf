package crawler

import (
	crawler2 "direwolf/internal/domain/repository/crawler"
	"direwolf/internal/domain/service/crawler"
)

type service struct {
	Engine     crawler.Engine
	Repository crawler2.Repository
}

// NewService creates new crawler service.
// both arguments are assumed to be fully initialized before being passed
func NewService(e crawler.Engine, r crawler2.Repository) crawler.Crawler {
	return &service{
		Engine:     e,
		Repository: r,
	}
}

// Crawl makes crawling of links
func (s *service) Crawl(links []string) {}
