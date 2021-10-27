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
	{ID: 5, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}}}

func TestKafkaAndDBUpdErrors(t *testing.T) {
	t.Parallel()

	fixture := mocks.Setup(t)
	defer fixture.TearDown()

	fixture.Repo.EXPECT().Lock(gomock.Any()).Return(eventsData, nil).AnyTimes()

	fixture.Repo.EXPECT().Remove(gomock.Any()).Return(errors.New("Remove execution error")).AnyTimes().After(
		fixture.Repo.EXPECT().Remove(gomock.Any()).Return(nil).Times(2))

	fixture.Repo.EXPECT().Unlock(gomock.Any()).Return(errors.New("Unlock execution error")).AnyTimes().After(
		fixture.Repo.EXPECT().Unlock(gomock.Any()).Return(nil).Times(2))

	fixture.Sender.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes().After(
		fixture.Sender.EXPECT().Send(gomock.Any()).Return(errors.New("Sending to kafka error")).Times(6).After(
			fixture.Sender.EXPECT().Send(gomock.Any()).Return(nil).Times(6)))

	startRetranslator(fixture)
}

func TestKafkaErrors(t *testing.T) {
	t.Parallel()

	fixture := mocks.Setup(t)
	defer fixture.TearDown()

	fixture.Repo.EXPECT().Lock(gomock.Any()).Return(eventsData, nil).AnyTimes()
	fixture.Repo.EXPECT().Remove(gomock.Any()).Return(nil).AnyTimes()
	fixture.Repo.EXPECT().Unlock(gomock.Any()).Return(nil).AnyTimes()

	fixture.Sender.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes().After(
		fixture.Sender.EXPECT().Send(gomock.Any()).Return(errors.New("Sending to kafka error")).Times(5).After(
			fixture.Sender.EXPECT().Send(gomock.Any()).Return(nil).Times(3)))

	startRetranslator(fixture)
}

func TestLockErrors(t *testing.T) {
	t.Parallel()

	fixture := mocks.Setup(t)
	defer fixture.TearDown()

	fixture.Repo.EXPECT().Lock(gomock.Any()).Return(nil, errors.New("DB is down")).AnyTimes()
	fixture.Repo.EXPECT().Remove(gomock.Any()).Return(nil).AnyTimes()
	fixture.Repo.EXPECT().Unlock(gomock.Any()).Return(nil).AnyTimes()
	fixture.Sender.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()

	startRetranslator(fixture)
}

func TestWithoutErrors(t *testing.T) {
	t.Parallel()

	fixture := mocks.Setup(t)
	defer fixture.TearDown()

	fixture.Repo.EXPECT().Lock(gomock.Any()).Return(eventsData, nil).AnyTimes()
	fixture.Repo.EXPECT().Remove(gomock.Any()).Return(nil).AnyTimes()
	fixture.Repo.EXPECT().Unlock(gomock.Any()).Return(nil).AnyTimes()
	fixture.Sender.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()

	startRetranslator(fixture)
}

func startRetranslator(fixture mocks.RetranslatorMockFixture) {
	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 1 * time.Second,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           fixture.Repo,
		Sender:         fixture.Sender,
	}

	retranslator := NewRetranslator(cfg)
	retranslator.Start()

	time.Sleep(5 * time.Second)

	retranslator.Close()
}
