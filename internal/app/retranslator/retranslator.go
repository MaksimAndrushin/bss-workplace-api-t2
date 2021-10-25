package retranslator

import (
	"context"
	"time"

	"github.com/ozonmp/omp-demo-api/internal/app/consumer"
	"github.com/ozonmp/omp-demo-api/internal/app/producer"
	"github.com/ozonmp/omp-demo-api/internal/app/repo"
	"github.com/ozonmp/omp-demo-api/internal/app/sender"
	"github.com/ozonmp/omp-demo-api/internal/model"

	"github.com/gammazero/workerpool"
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

	Ctx context.Context
	CancelFunc context.CancelFunc
}

type retranslator struct {
	events     chan model.WorkplaceEvent
	consumer   consumer.Consumer
	producer   producer.Producer
	workerPool *workerpool.WorkerPool
	cancelFunc context.CancelFunc
}

func NewRetranslator(cfg Config) Retranslator {
	var events = make(chan model.WorkplaceEvent, cfg.ChannelSize)
	var workerPool = workerpool.New(cfg.WorkerCount)

	var consumer = consumer.NewDbConsumer(
		cfg.Ctx,
		cfg.ConsumerCount,
		cfg.ConsumeSize,
		cfg.ConsumeTimeout,
		cfg.Repo,
		events)
	var producer = producer.NewKafkaProducer(
		cfg.Ctx,
		cfg.ProducerCount,
		cfg.Sender,
		events,
		workerPool,
		cfg.Repo)

	return &retranslator{
		events:     events,
		consumer:   consumer,
		producer:   producer,
		workerPool: workerPool,
		cancelFunc: cfg.CancelFunc,
	}
}

func (r *retranslator) Start() {
	r.producer.Start()
	r.consumer.Start()
}

func (r *retranslator) Close() {
	r.cancelFunc()

	r.consumer.Close()
	r.producer.Close()

	r.workerPool.StopWait()
}
