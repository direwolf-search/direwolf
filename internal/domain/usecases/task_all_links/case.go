package task_all_links

import (
	"context"

	"direwolf/internal/domain/model/task"
	"direwolf/internal/domain/repository"
)

type UseCaseTaskAllLinks struct {
	Repository repository.Repository
	TaskChan   chan *task.CrawlerTask
}

func NewUseCaseTaskAllLinks(r repository.Repository) *UseCaseTaskAllLinks {
	return &UseCaseTaskAllLinks{
		Repository: r,
		TaskChan:   make(chan *task.CrawlerTask, 0),
	}
}

func (uct *UseCaseTaskAllLinks) Run(ctx context.Context) {

}
