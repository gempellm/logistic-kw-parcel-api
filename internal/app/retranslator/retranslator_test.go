package retranslator

import (
	"math/rand"
	"testing"
	"time"

	"github.com/gempellm/logistic-kw-parcel-api/internal/mocks"
	"github.com/gempellm/logistic-kw-parcel-api/internal/model"
	"github.com/golang/mock/gomock"
)

func Test_EverythingIsOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 1 * time.Second, // default is: 10 * time.Second
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           repo,
		Sender:         sender,
	}

	getEvents := func() ([]model.ParcelEvent, error) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		events := make([]model.ParcelEvent, cfg.ConsumeSize)

		for i := range events {
			events[i] = model.ParcelEvent{
				ID:     uint64(i),
				Type:   model.Created,
				Status: model.Deferred,
				Entity: &model.Parcel{ID: r.Uint64()},
			}
		}

		return events, nil
	}

	repo.EXPECT().Lock(cfg.ConsumeSize).Return(getEvents()).AnyTimes()
	sender.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
	repo.EXPECT().Remove(gomock.Any()).Return(nil).AnyTimes()

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	time.Sleep(cfg.ConsumeTimeout * 2)
	retranslator.Close()
}
