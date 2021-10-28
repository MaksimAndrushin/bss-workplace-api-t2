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
	orderMap sync.Map
}

func NewOrderer() EventOrderer {
	return &workplaceOrder{}
}

func (o *workplaceOrder) AddEvents(events []model.WorkplaceEvent) {
	for _, event := range events {
		if eventsIdIntrf, ok := o.orderMap.Load(event.Entity.ID); !ok {
			o.orderMap.Store(event.Entity.ID, []uint64{event.ID})
		} else { //TODO Сортировку сделать после всей обработки, только для измененных слайсов.
			eventsId := eventsIdIntrf.([]uint64)
			eventsId = append(eventsId, event.ID)
			sort.Slice(eventsId, func(i, j int) bool { return eventsId[i] < eventsId[j] })
			o.orderMap.Store(event.Entity.ID, eventsId)
		}
	}
}

func (o *workplaceOrder) AddEvent(event model.WorkplaceEvent) {
	if eventsIdIntrf, ok := o.orderMap.Load(event.Entity.ID); !ok {
		o.orderMap.Store(event.Entity.ID, []uint64{event.ID})
	} else {
		eventsId := eventsIdIntrf.([]uint64)
		eventsId = append(eventsId, event.ID)
		sort.Slice(eventsId, func(i, j int) bool { return eventsId[i] < eventsId[j] })

		o.orderMap.Store(event.Entity.ID, eventsId)
	}
}

func (o *workplaceOrder) DeleteEvent(event model.WorkplaceEvent) {
	if eventsIdIntrf, ok := o.orderMap.Load(event.Entity.ID); ok {
		eventsId := eventsIdIntrf.([]uint64)
		if ind := sort.Search(len(eventsId), func(i int) bool { return event.ID <= eventsId[i] }); ind != len(eventsId) {
			o.orderMap.Store(event.Entity.ID, append(eventsId[:ind], eventsId[ind+1:]...))
		}
	}
}

func (o *workplaceOrder) IsEventAscOrdered(event model.WorkplaceEvent) bool {
	if eventsIdIntrf, ok := o.orderMap.Load(event.Entity.ID); !ok {
		return true
	} else {
		eventsId := eventsIdIntrf.([]uint64)

		if len(eventsId) < 1 {
			return true
		}

		if eventsId[0] >= event.ID {
			return true
		}

		return false
	}
}
