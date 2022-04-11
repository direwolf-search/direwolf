package crawler

import (
	"context"
	"direwolf/internal/domain/service/taskpool"
	"log"

	"github.com/robfig/cron/v3"

	"direwolf/internal/domain/model/task"
	"direwolf/internal/domain/repository"
	"direwolf/internal/domain/service/crawler"
	"direwolf/internal/factory/app"
)

type appCrawler struct {
	Crawler    crawler.Crawler
	Repository repository.Repository
	TaskPool   taskpool.TaskPool
	tasks      []*task.CrawlerTask // TODO: taskpool
}

func NewAppCrawler(crawler crawler.Crawler, taskPool taskpool.TaskPool, repo repository.Repository) app.App {
	return &appCrawler{
		Crawler:    crawler,
		Repository: repo,
		TaskPool:   taskPool,
	}
}

func (ac *appCrawler) Do(ctx context.Context) {
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
		log.Println("Crawler finished")

	}()
}
