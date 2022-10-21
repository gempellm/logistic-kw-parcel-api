package repo

import (
	"context"
	"errors"

	"github.com/gempellm/logistic-parcel-api/internal/model"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	CreateParcel(ctx context.Context, name string) (*model.Parcel, error)
	DescribeParcel(ctx context.Context, parcelID uint64) (*model.Parcel, error)
	ListParcels(ctx context.Context, cursor, offset uint64) ([]*model.Parcel, error)
	RemoveParcel(ctx context.Context, parcelID uint64) (bool, error)
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

func (r *repo) CreateParcel(ctx context.Context, name string) (*model.Parcel, error) {
	var p *model.Parcel

	if len(Parcels) == 0 {
		p = &model.Parcel{ID: 0, Name: name, Created: timestamppb.Now()}
	} else {
		p = &model.Parcel{ID: (Parcels[len(Parcels)-1].ID) + 1, Name: name, Created: timestamppb.Now()}
	}

	Parcels = append(Parcels, p)

	return p, nil
}

func (r *repo) ListParcels(ctx context.Context, cursor, offset uint64) ([]*model.Parcel, error) {
	parcelsLength := uint64(len(Parcels))

	if parcelsLength == 0 {
		return nil, nil
	}

	lastID := uint64(parcelsLength - 1)
	if cursor > lastID {
		return nil, nil
	}

	var result []*model.Parcel = make([]*model.Parcel, 0)

	for i := cursor; ; i++ {
		if i > lastID {
			break
		}

		result = append(result, Parcels[i])
	}

	return result, nil
}

func (r *repo) RemoveParcel(ctx context.Context, parcelID uint64) (bool, error) {
	var ok bool

	for i, v := range Parcels {
		if v.ID == parcelID {
			newParcels := make([]*model.Parcel, 0)
			newParcels = append(newParcels, Parcels[:i]...)
			newParcels = append(newParcels, Parcels[i+1:]...)

			Parcels = newParcels
			ok = true
			break
		}
	}

	return ok, nil
}
