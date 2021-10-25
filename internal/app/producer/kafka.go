package producer

import (
	"fmt"
	"github.com/ozonmp/omp-demo-api/internal/app/repo"
	"log"
	"sync"
	"time"

	"github.com/ozonmp/omp-demo-api/internal/app/sender"
	"github.com/ozonmp/omp-demo-api/internal/model"

	"github.com/gammazero/workerpool"
)

type Producer interface {
	Start()
	Close()
}

type producer struct {
	n       uint64
	timeout time.Duration

	sender sender.EventSender
	events <-chan model.WorkplaceEvent

	workerPool *workerpool.WorkerPool

	wg   *sync.WaitGroup
	done chan bool

	repo repo.EventRepo
}

func NewKafkaProducer(
	n uint64,
	sender sender.EventSender,
	events <-chan model.WorkplaceEvent,
	workerPool *workerpool.WorkerPool,
	repo repo.EventRepo,
) Producer {

	var wg = &sync.WaitGroup{}
	var done = make(chan bool)

	return &producer{
		n:          n,
		sender:     sender,
		repo:       repo,
		events:     events,
		workerPool: workerPool,
		wg:         wg,
		done:       done,
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
					processEvent(p, event)
				case <-p.done:
					return
				}
			}
		}()
	}
}

func processEvent(p *producer, event model.WorkplaceEvent) {
	if err := p.sender.Send(&event); err != nil {
		log.Println(fmt.Sprintf("ERROR!!!! Event ID - %d not sended to kafka", event.ID))

		p.workerPool.Submit(func() {
			if err := p.repo.Unlock([]uint64{event.ID}); err != nil {
				log.Println(fmt.Sprintf("UNLOCK ERROR!!!! Event ID - %d is not unlocked in DB", event.ID))
			}
		})
	} else {
		p.workerPool.Submit(func() {
			if err := p.repo.Remove([]uint64{event.ID}); err != nil {
				log.Println(fmt.Sprintf("REMOVE ERROR!!!! Event ID - %d is not deleted in DB", event.ID))
			}
		})
	}
}

func (p *producer) Close() {
	close(p.done)
	p.wg.Wait()
}
