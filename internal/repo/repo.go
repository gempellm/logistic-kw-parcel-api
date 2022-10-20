package repo

import (
	"context"
	"errors"

	"github.com/gempellm/logistic-parcel-api/internal/model"
	"github.com/jmoiron/sqlx"
)

var ErrParcelNotFound error = errors.New("parcel not found")

var Parcels []*model.Parcel = []*model.Parcel{
	&model.Parcel{ID: 1, Name: "first parcel"},
	&model.Parcel{ID: 2, Name: "second parcel"},
	&model.Parcel{ID: 3, Name: "third parcel"},
	&model.Parcel{ID: 4, Name: "fourth parcel"},
	&model.Parcel{ID: 5, Name: "fifth parcel"},
}

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
	for _, v := range Parcels {
		if v.ID == parcelID {
			return v, nil
		}
	}
	return nil, ErrParcelNotFound
}
