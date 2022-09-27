package consumer

import (
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

	done chan bool
	wg   *sync.WaitGroup
}

type Config struct {
	n         uint64
	events    chan<- model.ParcelEvent
	repo      repo.EventRepo
	batchSize uint64
	timeout   time.Duration
}

func NewDbConsumer(
	n uint64,
	batchSize uint64,
	consumerTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.ParcelEvent) Consumer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &consumer{
		n:         n,
		batchSize: batchSize,
		timeout:   consumerTimeout,
		repo:      repo,
		events:    events,
		wg:        wg,
		done:      done,
	}
}

func (c *consumer) Start() {
	for j := uint64(0); j < c.n; j++ {
		c.wg.Add(1)
		go func(i uint64) {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout)

			for {
				select {
				case <-ticker.C:
					events, err := c.repo.Lock(c.batchSize)
					// getEvents := func() ([]model.ParcelEvent, error) {
					// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
					// 	events := make([]model.ParcelEvent, 2)

					// 	for i := range events {
					// 		events[i] = model.ParcelEvent{
					// 			ID:     uint64(i),
					// 			Type:   model.Created,
					// 			Status: model.Deferred,
					// 			Entity: &model.Parcel{ID: r.Uint64()},
					// 		}
					// 	}

					// 	return events, nil
					// }

					// events, err := getEvents()
					if err != nil {
						continue
					}

					for _, event := range events {
						c.events <- event
					}

				case <-c.done:
					return
				}
			}
		}(j)
	}
}

func (c *consumer) Close() {
	close(c.done)
	c.wg.Wait()
}
