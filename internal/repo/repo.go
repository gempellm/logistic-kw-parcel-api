package repo

import (
	"context"

	"github.com/gempellm/logistic-parcel-api/internal/model"
	"github.com/jmoiron/sqlx"
)

// Repo is DAO for parcel
type Repo interface {
	DescribeParcel(ctx context.Context, parcelID uint64) (*model.Parcel, error)
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) DescribeParcel(ctx context.Context, parcelID uint64) (*model.Parcel, error) {
	return nil, nil
}
