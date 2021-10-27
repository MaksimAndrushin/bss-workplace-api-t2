package retranslator

import (
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

	Repo      repo.EventRepo
	Sender    sender.EventSender
	BatchSize int
}

type retranslator struct {
	events     chan model.WorkplaceEvent
	consumer   consumer.Consumer
	producer   producer.Producer
	workerPool *workerpool.WorkerPool
}

func NewRetranslator(cfg Config) Retranslator {
	var events = make(chan model.WorkplaceEvent, cfg.ChannelSize)
	var workerPool = workerpool.New(cfg.WorkerCount)

	// Для реализации гарантии at-least-once при запуске ретранслятора снимаются все локи,
	// висящие дольше определенного времени. Реализация должна быть сделана в методе UnlockAll
	cfg.Repo.UnlockAll()

	var consumer = consumer.NewDbConsumer(
		cfg.ConsumerCount,
		cfg.ConsumeSize,
		cfg.ConsumeTimeout,
		cfg.Repo,
		events)
	var producer = producer.NewKafkaProducer(
		cfg.ProducerCount,
		cfg.Sender,
		events,
		workerPool,
		cfg.Repo,
		cfg.BatchSize)

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
