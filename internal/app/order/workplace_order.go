package workplaceOrder

import (
	"github.com/ozonmp/omp-demo-api/internal/model"
	"sort"
	"sync"
)

type EventOrderer interface {
	AddEvent(event model.WorkplaceEvent)
	AddEvents(events []model.WorkplaceEvent)
	DeleteEvent(event model.WorkplaceEvent)
	IsEventAscOrdered(event model.WorkplaceEvent) bool
}

type workplaceOrder struct {
	orderMap map[uint64][]uint64
	mu       sync.Mutex
}

func NewOrderer() EventOrderer {
	return &workplaceOrder{
		orderMap: make(map[uint64][]uint64),
	}
}

func (o *workplaceOrder) AddEvents(events []model.WorkplaceEvent) {
	o.mu.Lock()
	defer o.mu.Unlock()

	for _, event := range events {
		if eventsId, ok := o.orderMap[event.Entity.ID]; !ok {
			o.orderMap[event.Entity.ID] = []uint64{event.ID}
		} else { //TODO Сортировку сделать после всей обработки, только для измененных слайсов.
			eventsId = append(eventsId, event.ID)
			sort.Slice(eventsId, func(i, j int) bool { return eventsId[i] < eventsId[j] })

			o.orderMap[event.Entity.ID] = eventsId
		}
	}
}

func (o *workplaceOrder) AddEvent(event model.WorkplaceEvent) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if eventsId, ok := o.orderMap[event.Entity.ID]; !ok {
		o.orderMap[event.Entity.ID] = []uint64{event.ID}
	} else {

		eventsId = append(eventsId, event.ID)
		sort.Slice(eventsId, func(i, j int) bool { return eventsId[i] < eventsId[j] })

		o.orderMap[event.Entity.ID] = eventsId
	}
}

func (o *workplaceOrder) DeleteEvent(event model.WorkplaceEvent) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if eventsId, ok := o.orderMap[event.Entity.ID]; ok {
		if ind := sort.Search(len(eventsId), func(i int) bool { return event.ID <= eventsId[i] }); ind != len(eventsId) {
			o.orderMap[event.Entity.ID] = append(eventsId[:ind], eventsId[ind+1:]...)
		}
	}
}

func (o *workplaceOrder) IsEventAscOrdered(event model.WorkplaceEvent) bool {
	o.mu.Lock()
	defer o.mu.Unlock()

	if eventsId, ok := o.orderMap[event.Entity.ID]; !ok {
		return true
	} else {
		if len(eventsId) < 1 {
			return true
		}

		if eventsId[0] >= event.ID {
			return true
		}

		return false
	}
}
