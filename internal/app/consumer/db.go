package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/ozonmp/omp-demo-api/internal/app/repo"
	"github.com/ozonmp/omp-demo-api/internal/model"
)

type Consumer interface {
	Start()
	Close()
}

type consumer struct {
	n      uint64
	events chan<- model.WorkplaceEvent

	repo repo.EventRepo

	batchSize uint64
	timeout   time.Duration

	ctxWithCancel context.Context
	wg   *sync.WaitGroup
}

type Config struct {
	n         uint64
	events    chan<- model.WorkplaceEvent
	repo      repo.EventRepo
	batchSize uint64
	timeout   time.Duration
}

func NewDbConsumer(
	ctx context.Context,
	n uint64,
	batchSize uint64,
	consumeTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.WorkplaceEvent) Consumer {

	var wg = &sync.WaitGroup{}

	return &consumer{
		n:         n,
		batchSize: batchSize,
		timeout:   consumeTimeout,
		repo:      repo,
		events:    events,
		wg:        wg,
		ctxWithCancel: ctx,
	}
}

func (c *consumer) Start() {
	for i := uint64(0); i < c.n; i++ {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout)
			for {
				select {
				case <-ticker.C:
					events, err := c.repo.Lock(c.batchSize)
					if err != nil {
						continue
					}
					for _, event := range events {
						c.events <- event
					}
				case <-c.ctxWithCancel.Done():
					return
				}
			}
		}()
	}
}

func (c *consumer) Close() {
	c.wg.Wait()
}
