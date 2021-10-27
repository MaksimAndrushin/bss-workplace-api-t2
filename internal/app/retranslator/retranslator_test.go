package retranslator

import (
	"errors"
	"github.com/ozonmp/omp-demo-api/internal/model"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ozonmp/omp-demo-api/internal/mocks"
)

var eventsData = []model.WorkplaceEvent{
	{ID: 1, Type: 0, Status: 1, Entity: &model.Workplace{ID: 1}},
	{ID: 2, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 3, Type: 0, Status: 1, Entity: &model.Workplace{ID: 3}},
	{ID: 4, Type: 0, Status: 1, Entity: &model.Workplace{ID: 1}},
	{ID: 5, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 6, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 7, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 8, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 9, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 10, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 11, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 12, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 13, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 14, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 15, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 16, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 17, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 18, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 19, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
	{ID: 20, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}},
}

func TestKafkaAndDBUpdErrors(t *testing.T) {
	t.Parallel()

	var ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	var repo = mocks.NewMockEventRepo(ctrl)
	var sender = mocks.NewMockEventSender(ctrl)

	repo.EXPECT().Lock(gomock.Any()).Return(eventsData, nil).AnyTimes()
	repo.EXPECT().UnlockAll().Return(nil).AnyTimes()

	repo.EXPECT().Remove(gomock.Any()).Return(errors.New("Remove execution error")).AnyTimes().After(
		repo.EXPECT().Remove(gomock.Any()).Return(nil).Times(2))

	repo.EXPECT().Unlock(gomock.Any()).Return(errors.New("Unlock execution error")).AnyTimes().After(
		repo.EXPECT().Unlock(gomock.Any()).Return(nil).Times(2))

	sender.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes().After(
		sender.EXPECT().Send(gomock.Any()).Return(errors.New("Sending to kafka error")).Times(6).After(
			sender.EXPECT().Send(gomock.Any()).Return(nil).Times(6)))

	startRetranslator(repo, sender)
}

func TestKafkaErrors(t *testing.T) {
	t.Parallel()

	var ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	var repo = mocks.NewMockEventRepo(ctrl)
	var sender = mocks.NewMockEventSender(ctrl)

	repo.EXPECT().Lock(gomock.Any()).Return(eventsData, nil).AnyTimes()
	repo.EXPECT().UnlockAll().Return(nil).AnyTimes()
	repo.EXPECT().Remove(gomock.Any()).Return(nil).AnyTimes()
	repo.EXPECT().Unlock(gomock.Any()).Return(nil).AnyTimes()

	sender.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes().After(
		sender.EXPECT().Send(gomock.Any()).Return(errors.New("Sending to kafka error")).Times(5))

	startRetranslator(repo, sender)
}

func TestLockErrors(t *testing.T) {
	t.Parallel()

	var ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	var repo = mocks.NewMockEventRepo(ctrl)
	var sender = mocks.NewMockEventSender(ctrl)

	repo.EXPECT().Lock(gomock.Any()).Return(nil, errors.New("DB is down")).AnyTimes()
	repo.EXPECT().UnlockAll().Return(nil).AnyTimes()
	repo.EXPECT().Remove(gomock.Any()).Return(nil).AnyTimes()
	repo.EXPECT().Unlock(gomock.Any()).Return(nil).AnyTimes()
	sender.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()

	startRetranslator(repo, sender)
}

func TestWithoutErrors(t *testing.T) {
	t.Parallel()

	var ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	var repo = mocks.NewMockEventRepo(ctrl)
	var sender = mocks.NewMockEventSender(ctrl)

	repo.EXPECT().Lock(gomock.Any()).Return(eventsData, nil).AnyTimes()
	repo.EXPECT().UnlockAll().Return(nil).AnyTimes()
	repo.EXPECT().Remove(gomock.Any()).Return(nil).AnyTimes()
	repo.EXPECT().Unlock(gomock.Any()).Return(nil).AnyTimes()
	sender.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()

	startRetranslator(repo, sender)
}

func startRetranslator(repo *mocks.MockEventRepo, sender *mocks.MockEventSender) {
	var cfg = Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 1 * time.Second,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           repo,
		Sender:         sender,
		BatchSize:      10,
	}

	var retranslator = NewRetranslator(cfg)
	retranslator.Start()

	time.Sleep(5 * time.Second)

	retranslator.Close()
}
