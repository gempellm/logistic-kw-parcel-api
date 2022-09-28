package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gempellm/logistic-kw-parcel-api/internal/app/repo"
	"github.com/gempellm/logistic-kw-parcel-api/internal/model"
)

type Consumer interface {
	Start()
	Close()
}

type consumer struct {
	n      uint64
	events chan<- model.ParcelEvent

	repo repo.EventRepo

	batchSize uint64
	timeout   time.Duration

	//done chan bool
	wg *sync.WaitGroup
}

// type Config struct {
// 	n         uint64
// 	events    chan<- model.ParcelEvent
// 	repo      repo.EventRepo
// 	batchSize uint64
// 	timeout   time.Duration
// }

func NewDbConsumer(
	n uint64,
	batchSize uint64,
	consumerTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.ParcelEvent) Consumer {

	wg := &sync.WaitGroup{}
	//done := make(chan bool)

	return &consumer{
		n:         n,
		batchSize: batchSize,
		timeout:   consumerTimeout,
		repo:      repo,
		events:    events,
		wg:        wg,
		//done:      done,
	}
}

func (c *consumer) Start() {
	ctxBase := context.Background()
	ctx, _ := context.WithTimeout(ctxBase, c.timeout)
	for i := uint64(0); i < c.n; i++ {
		c.wg.Add(1)
		go func(j uint64) {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout)

			for {
				select {
				case <-ticker.C:
					events, err := c.repo.Lock(c.batchSize)

					if err != nil {
						fmt.Printf("<consumer-%d> error during c.repo.Lock(%v)\n", j, c.batchSize)
						continue
					}

					for _, event := range events {
						c.events <- event
					}

				case <-ctx.Done():
					return
				}
			}
		}(i)
	}
}

func (c *consumer) Close() {
	//close(c.done)
	c.wg.Wait()
}
