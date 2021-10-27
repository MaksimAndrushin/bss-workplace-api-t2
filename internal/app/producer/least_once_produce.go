package producer

import (
	"github.com/ozonmp/omp-demo-api/internal/app/repo"
	"log"
	"sync"
	"time"
)

type Repeater interface {
	Start()
	Close()
	ToUnlockRepeater(unlockBatch *[]uint64)
	ToRemoveRepeater(removeBatch *[]uint64)
}

type repeater struct {
	repeatPeriod time.Duration

	repo repo.EventRepo

	initCapacity uint64

	unlockRepeat []uint64
	urMutex      sync.Mutex

	removeRepeat []uint64
	rrMutex      sync.Mutex

	wg   *sync.WaitGroup
	done chan bool
}

func NewRepeater(
	repeatPeriod time.Duration,
	repo repo.EventRepo,
	repeaterInitSize uint64,
) Repeater {

	var wg = &sync.WaitGroup{}
	done := make(chan bool)

	return &repeater{
		repeatPeriod: repeatPeriod,
		repo:         repo,
		initCapacity: repeaterInitSize,
		unlockRepeat: make([]uint64, 0, repeaterInitSize),
		removeRepeat: make([]uint64, 0, repeaterInitSize),
		wg:           wg,
		done:         done,
	}
}

func (r *repeater) Start() {

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()

		var syncTicker = time.NewTicker(r.repeatPeriod)
		defer syncTicker.Stop()

		for {
			select {
			case <-syncTicker.C: // Сброс в БД накопившихся батчей, если они не сброшены дольше forceSyncPeriod секунд
				r.flushUnlockBatch()
				r.flushRemoveBatch()

			case <-r.done:
				r.flushUnlockBatch()
				r.flushRemoveBatch()
				return
			}
		}
	}()

}

func (r *repeater) ToUnlockRepeater(unlockBatch *[]uint64) {
	r.urMutex.Lock()
	defer r.urMutex.Unlock()

	r.unlockRepeat = append(r.unlockRepeat, *unlockBatch...)
}

func (r *repeater) ToRemoveRepeater(removeBatch *[]uint64) {
	r.rrMutex.Lock()
	defer r.rrMutex.Unlock()

	r.removeRepeat = append(r.removeRepeat, *removeBatch...)
}

func (r *repeater) flushUnlockBatch() {
	r.urMutex.Lock()
	defer r.urMutex.Unlock()

	if len(r.unlockRepeat) != 0 {
		if err := r.repo.Unlock(r.unlockRepeat); err != nil {
			log.Println("UNLOCK ERROR!!!!")
		} else {
			r.unlockRepeat = make([]uint64, 0, r.initCapacity)
		}
	}
}

func (r *repeater) flushRemoveBatch() {
	r.rrMutex.Lock()
	defer r.rrMutex.Unlock()

	if len(r.removeRepeat) != 0 {
		if err := r.repo.Remove(r.removeRepeat); err != nil {
			log.Println("UNLOCK ERROR!!!!")
		} else {
			r.removeRepeat = make([]uint64, 0, r.initCapacity)
		}
	}
}

func (r *repeater) Close() {
	close(r.done)
	r.wg.Wait()
}
