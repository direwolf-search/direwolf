package crawlall

import (
	"context"
	crawler2 "direwolf/internal/domain/repository/crawler"
	"log"

	"direwolf/internal/domain/service/crawler"
)

type CrawlAllUseCase struct {
	Context    context.Context
	Crawler    crawler.Crawler
	Repository crawler2.Repository
}

func NewCrawlAllUseCase(ctx context.Context, c crawler.Crawler, r crawler2.Repository) *CrawlAllUseCase {
	return &CrawlAllUseCase{
		Context:    ctx,
		Crawler:    c,
		Repository: r,
	}
}

func (cauc *CrawlAllUseCase) Run() {
	links, err := cauc.Repository.GetAll(cauc.Context)
	if err != nil {
		//return err
		log.Println(err)
	}
	cauc.Crawler.Crawl(links)
}
