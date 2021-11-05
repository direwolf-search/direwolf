package crawl

import (
	"context"

	"direwolf/internal/domain/repository"
	"direwolf/internal/domain/service/crawler"
)

type UseCaseCrawl struct {
	Crawler    crawler.Crawler
	Repository repository.Repository
}

func NewUseCase(c crawler.Crawler, r repository.Repository) *UseCaseCrawl {
	return &UseCaseCrawl{
		Crawler:    c,
		Repository: r,
	}
}

func (ucc *UseCaseCrawl) Run(ctx context.Context) {
	//ucc.Crawler.VisitAll(ctx, ucc.Repository.Insert, ucc.UrlsSource...)
	ucc.Crawler.DoTasks() // TODO:
}
