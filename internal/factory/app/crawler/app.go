package crawler

import (
	"context"
	"direwolf/internal/domain/model/task"
	"log"

	"github.com/robfig/cron"

	"direwolf/internal/domain/usecases/crawl"
	"direwolf/internal/factory/app"
)

type appCrawler struct {
	crawlerUseCase *crawl.UseCaseCrawl
	tasks          []*task.CrawlerTask
}

func NewAppCrawler(useCase *crawl.UseCaseCrawl, tasks ...*task.CrawlerTask) app.App {
	return &appCrawler{
		crawlerUseCase: useCase,
		tasks:          tasks,
	}
}

func (ac *appCrawler) Do(ctx context.Context) {
	go func() {
		ac.crawlerUseCase.Run(ctx)

		for _, task := range ac.tasks {
			cron := cron.New()
			if err := cron.AddFunc(task.Schedule(), func() {
				ac.crawlerUseCase.Run(ctx)
			}); err != nil {
				log.Println("Error adding to cron:", err)
			}
			cron.Start()
		}

		<-ctx.Done()
		log.Println("Shutdown service")

	}()
}
