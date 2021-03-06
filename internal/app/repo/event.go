package repo

import (
	"github.com/ozonmp/omp-demo-api/internal/model"
)

type EventRepo interface {
	Lock(n uint64) ([]model.WorkplaceEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.WorkplaceEvent) error
	Remove(eventIDs []uint64) error
}
