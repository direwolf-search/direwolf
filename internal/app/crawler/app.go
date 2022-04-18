package crawler

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"

	"direwolf/internal/domain"
	"direwolf/internal/domain/repository"
	"direwolf/internal/domain/service/crawler"
	"direwolf/internal/domain/service/scheduler"
	crawlall "direwolf/internal/domain/usecases/crawl_all"
)

type appCrawler struct {
	Crawler    crawler.Crawler
	Repository repository.Repository
	Scheduler  scheduler.Scheduler
	Logger     domain.Logger
}

func NewAppCrawler(crawler crawler.Crawler, logger domain.Logger, taskPool scheduler.Scheduler, repo repository.Repository) domain.App {
	return &appCrawler{
		Crawler:    crawler,
		Repository: repo,
		Scheduler:  taskPool,
		Logger:     logger,
	}
}

func (ac *appCrawler) Do(ctx context.Context) {
	ac.Scheduler.Maintain("crawler_service")
	tasks, err := ac.Scheduler.GetTasks(ctx)
	if err != nil {
		ac.Logger.Error(err, "cannot get tasks for crawler_service: ")
	}
	useCase := crawlall.NewCrawlAllUseCase(ctx, ac.Crawler, ac.Repository)
	go func() {
		ac.crawlerUseCase.Crawler.GetTask()
		ac.crawlerUseCase.Run(ctx)

		for _, task := range ac.tasks {
			cron := cron.New()
			if _, err := cron.AddFunc(task.Schedule(), func() {
				ac.crawlerUseCase.Crawler.GetTask() // TODO: next steps
				ac.crawlerUseCase.Run(ctx)
			}); err != nil {
				log.Println("Error adding to cron:", err)
			}
			cron.Start()
		}

		<-ctx.Done()
		log.Println("Crawler stopped")

	}()
}
