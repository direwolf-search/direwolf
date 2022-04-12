package crawlall

import (
	"context"
	"log"

	"direwolf/internal/domain/repository"
	"direwolf/internal/domain/service/crawler"
)

type CrawlAllUseCase struct {
	Context    context.Context
	Crawler    crawler.Crawler
	Repository repository.Repository
}

func NewCrawlAllUseCase(ctx context.Context, c crawler.Crawler, r repository.Repository) *CrawlAllUseCase {
	return &CrawlAllUseCase{
		Context:    ctx,
		Crawler:    c,
		Repository: r,
	}
}

func (ca *CrawlAllUseCase) Run() {
	links, err := ca.Repository.GetAll(ca.Context)
	if err != nil {
		//return err
		log.Println(err)
	}
	ca.Crawler.Crawl(links)
}
