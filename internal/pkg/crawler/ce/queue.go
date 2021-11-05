package ce

import (
	"context"
	"log"

	"github.com/gocolly/colly/v2"
	cq "github.com/gocolly/colly/v2/queue"
)

type Queue struct {
	collyQueue *cq.Queue
}

func NewQueue(workersNum int) *Queue {
	q, err := cq.New(
		workersNum,
		&cq.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)
	if err != nil {
		log.Fatalln(err)
	}
	return &Queue{
		collyQueue: q,
	}
}

func (q *Queue) AddRequest(
	ctx context.Context,
	arg interface{},
	r *colly.Request,
	f func(ctx context.Context, entity interface{}) error,
) error {
	if arg != nil {
		err := f(ctx, arg)
		if err != nil {
			return err
		}
	}

	err := q.collyQueue.AddRequest(r)
	if err != nil {
		return err
	}

	return nil
}

func (q *Queue) Run(c *colly.Collector) error {
	return q.collyQueue.Run(c)
}
