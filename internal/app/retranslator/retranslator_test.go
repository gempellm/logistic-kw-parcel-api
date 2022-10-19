package retranslator

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/gempellm/logistic-kw-parcel-api/internal/mocks"
	"github.com/gempellm/logistic-kw-parcel-api/internal/model"
	"github.com/golang/mock/gomock"
)

type TestCase uint8

const (
	Ok TestCase = iota
	NotOk
)

func generateEvents(amount uint64) []model.ParcelEvent {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	events := make([]model.ParcelEvent, amount)

	for i := range events {

		events[i] = model.ParcelEvent{
			ID:     uint64(i),
			Type:   model.Created,
			Status: model.Deferred,
			Entity: &model.Parcel{ID: r.Uint64()},
		}
	}

	return events
}

func getEvents(cfg Config, testCase TestCase) ([]model.ParcelEvent, error) {
	var events []model.ParcelEvent

	switch testCase {
	case Ok:
		events = generateEvents(cfg.ConsumeSize)
		return events, nil
	case NotOk:
		return nil, errors.New("mock error during Lock")
	default:
		return nil, errors.New("unknown test case")
	}

	//events := make([]model.ParcelEvent, cfg.Config.ConsumeSize)
}

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

	repo.EXPECT().Lock(cfg.ConsumeSize).Return(getEvents(cfg, Ok)).MinTimes(1)
	sender.EXPECT().Send(gomock.Any()).Return(nil).MinTimes(10)
	repo.EXPECT().Remove(gomock.Any()).Return(nil).MinTimes(10)
	repo.EXPECT().Add(gomock.Any()).Return(nil).MinTimes(10)
	repo.EXPECT().Unlock(gomock.Any()).Times(0)

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	time.Sleep(cfg.ConsumeTimeout * 2)
	retranslator.Close()
}

func Test_SendError(t *testing.T) {
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

	repo.EXPECT().Lock(cfg.ConsumeSize).Return(getEvents(cfg, Ok)).MinTimes(1)
	sender.EXPECT().Send(gomock.Any()).Return(errors.New("mock error")).MinTimes(1)
	repo.EXPECT().Remove(gomock.Any()).Return(nil).Times(0)
	repo.EXPECT().Add(gomock.Any()).Return(nil).Times(0)
	repo.EXPECT().Unlock(gomock.Any()).MinTimes(1)

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	time.Sleep(cfg.ConsumeTimeout * 2)
	retranslator.Close()
}

func Test_LockError(t *testing.T) {
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

	gomock.InOrder(
		repo.EXPECT().Lock(cfg.ConsumeSize).Return(getEvents(cfg, NotOk)).MinTimes(1),
		sender.EXPECT().Send(gomock.Any()).Return(errors.New("mock error")).Times(0),
		repo.EXPECT().Remove(gomock.Any()).Return(nil).Times(0),
		repo.EXPECT().Add(gomock.Any()).Return(nil).Times(0),
		repo.EXPECT().Unlock(gomock.Any()).Times(0),
	)

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	retranslator.Close()
}
