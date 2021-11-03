package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/bss-workplace-api/internal/model"
)

// Repo is DAO for workplace
type Repo interface {
	CreateWorkplace(ctx context.Context, foo string) (uint64, error)
	DescribeWorkplace(ctx context.Context, workplaceID uint64) (*model.Workplace, error)
	ListWorkplaces(ctx context.Context) ([]model.Workplace, error)
	RemoveWorkplace(ctx context.Context, workplaceID uint64) (bool, error)
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) CreateWorkplace(ctx context.Context, foo string) (uint64, error) {
	return 0, nil
}

func (r *repo) DescribeWorkplace(ctx context.Context, workplaceID uint64) (*model.Workplace, error) {
	return nil, nil
}

func (r *repo) ListWorkplaces(ctx context.Context) ([]model.Workplace, error) {
	return nil, nil
}

func (r *repo) RemoveWorkplace(ctx context.Context, workplaceID uint64) (bool, error) {
	return false, nil
}
