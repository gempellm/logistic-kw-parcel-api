package repo

import "github.com/gempellm/logistic-kw-parcel-api/internal/model"

type EventRepo interface {
	Lock(n uint64) ([]model.ParcelEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.ParcelEvent) error
	Remove(eventIDs []uint64) error
}
