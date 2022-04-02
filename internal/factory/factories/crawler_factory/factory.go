package crawler_factory

import (
	"os"

	"direwolf/internal/domain/repository"
	service "direwolf/internal/domain/service/crawler"
	"direwolf/internal/domain/usecases/crawl"
	"direwolf/internal/factory"
	"direwolf/internal/factory/app"
	app_crawler "direwolf/internal/factory/app/crawler"
	concrete "direwolf/internal/services/crawler"
	//"direwolf/internal/factory"
)

type crawlerFactory struct{}

func NewCrawlerFactory() factory.AppFactory {
	return &crawlerFactory{}
}

func (ef *crawlerFactory) BuildApp(components []interface{}) app.App {
	var (
		defaultEngineName     = os.Getenv("DW_DEFAULT_TOR_CRAWLER_ENGINE")     // TODO:
		defaultRepositoryName = os.Getenv("DW_DEFAULT_TOR_CRAWLER_REPOSITORY") // TODO:
		defaultSchedule       = os.Getenv("DW_DEFAULT_TOR_CRAWLER_SCHEDULE")
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

	return app_crawler.NewAppCrawler(crawl.NewUseCase(crawler, repo), defaultSchedule)
}
