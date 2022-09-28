package producer

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/gempellm/logistic-kw-parcel-api/internal/app/repo"
	"github.com/gempellm/logistic-kw-parcel-api/internal/app/sender"
	"github.com/gempellm/logistic-kw-parcel-api/internal/model"
)

type Producer interface {
	Start()
	Close()
}

type producer struct {
	n       uint64
	timeout time.Duration

	sender sender.EventSender
	events <-chan model.ParcelEvent

	workerPool *workerpool.WorkerPool

	wg   *sync.WaitGroup
	done chan bool

	repo repo.EventRepo
}

func NewKafkaProducer(
	n uint64,
	sender sender.EventSender,
	events <-chan model.ParcelEvent,
	workerPool *workerpool.WorkerPool,
	repo repo.EventRepo,
) Producer {
	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &producer{
		n:          n,
		sender:     sender,
		events:     events,
		workerPool: workerPool,
		wg:         wg,
		done:       done,
		repo:       repo,
	}
}

func (p *producer) Start() {
	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case event := <-p.events:
					err := p.sender.Send(&event)
					if err != nil {
						p.workerPool.Submit(func() {
							fmt.Printf("<producer> sender caught error while sending %+v: %v\n", event, err)
							err := p.repo.Unlock([]uint64{event.ID})
							if err != nil {
								log.Fatal("error occurred during p.repo.Unlock()")
							}
						})
					} else {
						p.workerPool.Submit(func() {
							fmt.Printf("<producer> sender successfully sent %+v\n", event)
							err := p.repo.Remove([]uint64{event.ID})
							if err != nil {
								log.Fatal("error occurred during p.repo.Remove()")
							}

							switch event.Type {
							case model.Created:
								event.Status = model.Processed
								err := p.repo.Add([]model.ParcelEvent{event})
								if err != nil {
									log.Fatal("error occurred during p.repo.Add()")
								}
								fmt.Printf("<producer> event.Status = model.Processed && p.repo.Add(): %+v\n", event)
							case model.Updated:
								fmt.Println("<producer> event.Type = model.Updated")
							case model.Removed:
								fmt.Println("<producer event.Type = model.Removed")
							default:
								log.Fatal("<producer> unknown event.Type")
							}
						})
					}
				case <-p.done:
					return
				}
			}
		}()
	}
}

func (p *producer) Close() {
	close(p.done)
	p.wg.Wait()
}
