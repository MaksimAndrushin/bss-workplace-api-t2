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

	dbBatchSize     int
	forceSyncPeriod time.Duration
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
		n:               n,
		sender:          sender,
		repo:            repo,
		events:          events,
		workerPool:      workerPool,
		wg:              wg,
		done:            done,
		dbBatchSize:     10,
		forceSyncPeriod: 5,
	}
}

func (p *producer) Start() {
	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)
		go func() {
			var unlockBatch = make([]uint64, 0, p.dbBatchSize)
			var removeBatch = make([]uint64, 0, p.dbBatchSize)

			defer p.wg.Done()
			defer p.flushUnlockBatch(&unlockBatch)
			defer p.flushRemoveBatch(&removeBatch)

			var syncTicker = time.NewTicker(p.forceSyncPeriod)
			defer syncTicker.Stop()

			for {
				select {
				case event := <-p.events:
					p.processEvent(event, &unlockBatch, &removeBatch)

				case <-syncTicker.C: // Сброс в БД накопившихся батчей, если они не сброшены дольше forceSyncPeriod секунд
					p.forceSyncEventWithDB(&unlockBatch, p.flushUnlockBatch)
					p.forceSyncEventWithDB(&removeBatch, p.flushRemoveBatch)

				case <-p.done:
					return
				}
			}
		}()
	}
}

func (p *producer) processEvent(event model.WorkplaceEvent, unlockBatch *[]uint64, removeBatch *[]uint64) {
	if err := p.sender.Send(&event); err != nil {
		log.Println(fmt.Sprintf("ERROR!!!! Event ID - %d not sended to kafka", event.ID))
		p.syncEventWithDB(event, unlockBatch, p.flushUnlockBatch)
	} else {
		p.syncEventWithDB(event, removeBatch, p.flushRemoveBatch)
	}
}

func (p *producer) forceSyncEventWithDB(dataBatch *[]uint64, syncFunc func(*[]uint64)) {
	if len(*dataBatch) <= 0 {
		return
	}

	var batchCopy = make([]uint64, len(*dataBatch))
	copy(batchCopy, *dataBatch)

	p.workerPool.Submit(func() {
		syncFunc(&batchCopy)
	})

	*dataBatch = (*dataBatch)[:0]
}

func (p *producer) syncEventWithDB(event model.WorkplaceEvent, dataBatch *[]uint64, syncFunc func(*[]uint64)) {
	*dataBatch = append(*dataBatch, event.ID)

	if len(*dataBatch) >= p.dbBatchSize {
		var batchCopy = make([]uint64, len(*dataBatch))
		copy(batchCopy, *dataBatch)

		p.workerPool.Submit(func() {
			syncFunc(&batchCopy)
		})

		*dataBatch = (*dataBatch)[:0]
	}
}

func (p *producer) flushUnlockBatch(unlockBatch *[]uint64) {
	if len(*unlockBatch) != 0 {
		if err := p.repo.Unlock(*unlockBatch); err != nil {
			log.Println(fmt.Sprintf("UNLOCK ERROR!!!! Event ID's - %v is not unlocked in DB", unlockBatch))
		}
	}
}

func (p *producer) flushRemoveBatch(removeBatch *[]uint64) {
	if len(*removeBatch) != 0 {
		if err := p.repo.Remove(*removeBatch); err != nil {
			log.Println(fmt.Sprintf("REMOVE ERROR!!!! Event ID's - %v is not deleted in DB", removeBatch))
		}
	}
}

func (p *producer) Close() {
	close(p.done)
	p.wg.Wait()
}
