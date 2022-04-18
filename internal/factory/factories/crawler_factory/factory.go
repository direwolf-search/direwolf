package crawler_factory

import (
	"os"

	app_crawler "direwolf/internal/app/crawler"
	"direwolf/internal/domain"
	"direwolf/internal/domain/repository"
	service "direwolf/internal/domain/service/crawler"
	"direwolf/internal/domain/usecases/crawl_all"
	"direwolf/internal/factory"
	concrete "direwolf/internal/services/crawler"
)

type crawlerFactory struct{}

func NewCrawlerFactory() factory.AppFactory {
	return &crawlerFactory{}
}

func (cf *crawlerFactory) BuildApp(components []interface{}) domain.App {
	var (
		defaultEngineName     = os.Getenv("DW_DEFAULT_TOR_CRAWLER_ENGINE") // TODO:
		defaultRepositoryName = os.Getenv("DW_DEFAULT_TOR_CRAWLER_REPOSITORY")
		engine                service.Engine
		repo                  repository.Repository
	)

	for _, component := range components {
		if factory.GetName(component) == defaultEngineName {
			if m, ok := component.(service.Engine); ok {
				engine = m
			}
		}

		if factory.GetName(component) == defaultRepositoryName {
			if m, ok := component.(repository.Repository); ok {
				repo = m
			}
		}
	}

	crawler := concrete.NewService(engine)

	return app_crawler.NewAppCrawler(crawlall.NewCrawlAllUseCase(crawler, repo), defaultSchedule)
}
