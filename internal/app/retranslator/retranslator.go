package retranslator

import (
	"time"

	"github.com/gammazero/workerpool"
	"github.com/gempellm/logistic-parcel-api/internal/app/consumer"
	"github.com/gempellm/logistic-parcel-api/internal/app/producer"
	"github.com/gempellm/logistic-parcel-api/internal/app/repo"
	"github.com/gempellm/logistic-parcel-api/internal/app/sender"
	"github.com/gempellm/logistic-parcel-api/internal/model"
)

type Retranslator interface {
	Start()
	Close()
}

type Config struct {
	ChannelSize uint64

	ConsumerCount  uint64
	ConsumeSize    uint64
	ConsumeTimeout time.Duration

	ProducerCount uint64
	WorkerCount   int

	Repo   repo.EventRepo
	Sender sender.EventSender
}

type retranslator struct {
	events     chan model.ParcelEvent
	consumer   consumer.Consumer
	producer   producer.Producer
	workerPool *workerpool.WorkerPool
}

func NewRetranslator(cfg Config) Retranslator {
	events := make(chan model.ParcelEvent, cfg.ChannelSize)
	workerPool := workerpool.New(cfg.WorkerCount)

	consumer := consumer.NewDbConsumer(
		cfg.ConsumerCount,
		cfg.ConsumeSize,
		cfg.ConsumeTimeout,
		cfg.Repo,
		events,
	)

	producer := producer.NewKafkaProducer(
		cfg.ProducerCount,
		cfg.Sender,
		events,
		workerPool,
	)

	return &retranslator{
		events:     events,
		consumer:   consumer,
		producer:   producer,
		workerPool: workerPool,
	}
}

func (r *retranslator) Start() {
	r.producer.Start()
	r.consumer.Start()
}

func (r *retranslator) Close() {
	r.consumer.Close()
	r.producer.Close()
	r.workerPool.StopWait()
}
