package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/bss-workplace-api/internal/model"
)

// Repo is DAO for workplace
type Repo interface {
	DescribeWorkplace(ctx context.Context, workplaceID uint64) (*model.Workplace, error)
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) DescribeWorkplace(ctx context.Context, workplaceID uint64) (*model.Workplace, error) {
	return nil, nil
}
