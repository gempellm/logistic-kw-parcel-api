package sender

import "github.com/gempellm/logistic-kw-parcel-api/internal/model"

type EventSender interface {
	Send(parcel *model.ParcelEvent) error
}
