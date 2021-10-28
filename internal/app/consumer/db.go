package consumer

import (
	"errors"
	workplaceOrder "github.com/ozonmp/omp-demo-api/internal/app/order"
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

	done chan interface{}
	wg   *sync.WaitGroup

	eventOrderer workplaceOrder.EventOrderer
}

type Config struct {
	n         uint64
	events    chan<- model.WorkplaceEvent
	repo      repo.EventRepo
	batchSize uint64
	timeout   time.Duration
}

func NewDbConsumer(
	n uint64,
	batchSize uint64,
	consumeTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.WorkplaceEvent,
	eventOrderer workplaceOrder.EventOrderer) Consumer {

	var wg = &sync.WaitGroup{}
	done := make(chan interface{})

	return &consumer{
		n:            n,
		batchSize:    batchSize,
		timeout:      consumeTimeout,
		repo:         repo,
		events:       events,
		wg:           wg,
		done:         done,
		eventOrderer: eventOrderer,
	}
}

/*
Чтобы был соблюден порядок выгрузки в кафку необходима гарантия в методе Lock:
- !!!! Каждая выборка возвращает евенты, у которых Workplace.ID будут присутствовать ТОЛЬКО на текущей инстанции
- !!!! Никакая инстанция НЕ МОЖЕТ залочить евенты с Workplace.ID, которые уже присутствуют в залоченых евентах
Эти гарантии необходимы для гарантии порядка выгрузки евентов на многоинстансной системе

- !!!! Метод лок возвращает евенты в порядке их возникновения

!!!!!!! Далее считаем что метод Lock выполняет все эти гарантии !!!!!!
*/

func (c *consumer) Start() {
	for i := uint64(0); i < c.n; i++ {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout)
			for {
				select {
				case <-ticker.C:
					if err := c.processTimerTick(); err != nil {
						continue
					}

				case <-c.done:
					return
				}
			}
		}()
	}
}

func (c *consumer) processTimerTick() error {
	events, err := c.repo.Lock(c.batchSize)
	if err != nil {
		return errors.New("Lock error")
	}

	c.eventOrderer.AddEvents(events)

	for _, event := range events {
		c.events <- event
	}

	return nil
}

func (c *consumer) Close() {
	close(c.done)
	c.wg.Wait()
}
