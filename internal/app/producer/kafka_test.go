package producer

import (
	"context"
	"errors"
	"github.com/gammazero/workerpool"
	"github.com/ozonmp/omp-demo-api/internal/mocks"
	"github.com/ozonmp/omp-demo-api/internal/model"
	"testing"
	"time"
)

type KafkaProducer struct {
	events     chan model.WorkplaceEvent
	workerPool *workerpool.WorkerPool
	producer   Producer
}

func TestSendToKafkaSuccessful(t *testing.T) {
	t.Parallel()

	fixture := mocks.Setup(t)
	defer fixture.TearDown()

	kafkaTestProducer := NewTestKafkaProducer(2, fixture)
	defer kafkaTestProducer.Close()

	event := model.WorkplaceEvent{ID: 1, Type: 0, Status: 1, Entity: &model.Workplace{ID: 1}}

	fixture.Sender.EXPECT().Send(&event).Return(nil).Times(1)
	fixture.Repo.EXPECT().Remove([]uint64{event.ID}).Times(1)

	ctx, cancelFunc := context.WithCancel(context.Background())
	kafkaTestProducer.producer.Start(ctx)

	kafkaTestProducer.events <- event
	time.Sleep(1 * time.Second)

	cancelFunc()
}

func TestSendToKafkaUnsuccessful(t *testing.T) {
	t.Parallel()

	fixture := mocks.Setup(t)
	defer fixture.TearDown()

	kafkaTestProducer := NewTestKafkaProducer(2, fixture)
	defer kafkaTestProducer.Close()

	event := model.WorkplaceEvent{ID: 1, Type: 0, Status: 1, Entity: &model.Workplace{ID: 1}}

	fixture.Sender.EXPECT().Send(&event).Return(errors.New("Sending error")).Times(1)
	fixture.Repo.EXPECT().Unlock([]uint64{event.ID}).Times(1)

	ctx, cancelFunc := context.WithCancel(context.Background())
	kafkaTestProducer.producer.Start(ctx)

	kafkaTestProducer.events <- event
	time.Sleep(1 * time.Second)

	cancelFunc()
}

func TestSendToKafkaThreeSEventsAndOneUEvent(t *testing.T) {
	t.Parallel()

	fixture := mocks.Setup(t)
	defer fixture.TearDown()

	kafkaTestProducer := NewTestKafkaProducer(2, fixture)
	defer kafkaTestProducer.Close()

	events := []model.WorkplaceEvent{
		{ID: 1, Type: 0, Status: 1, Entity: &model.Workplace{ID: 1}},
		{ID: 2, Type: 0, Status: 1, Entity: &model.Workplace{ID: 1}},
		{ID: 3, Type: 0, Status: 1, Entity: &model.Workplace{ID: 1}},
		{ID: 4, Type: 0, Status: 1, Entity: &model.Workplace{ID: 1}}}

	fixture.Sender.EXPECT().Send(&events[3]).Return(nil).Times(1).After(
		fixture.Sender.EXPECT().Send(&events[2]).Return(errors.New("Sending error")).Times(1).After(
			fixture.Sender.EXPECT().Send(&events[1]).Return(nil).Times(1).After(
				fixture.Sender.EXPECT().Send(&events[0]).Return(nil).Times(1))))

	fixture.Repo.EXPECT().Remove([]uint64{events[0].ID}).Times(1)
	fixture.Repo.EXPECT().Remove([]uint64{events[1].ID}).Times(1)
	fixture.Repo.EXPECT().Unlock([]uint64{events[2].ID}).Times(1)
	fixture.Repo.EXPECT().Remove([]uint64{events[3].ID}).Times(1)

	ctx, cancelFunc := context.WithCancel(context.Background())
	kafkaTestProducer.producer.Start(ctx)

	for _, v := range events {
		kafkaTestProducer.events <- v
		time.Sleep(1 * time.Second)
	}

	cancelFunc()

}



func NewTestKafkaProducer(kafkaWorkers uint64, fixture mocks.RetranslatorMockFixture) KafkaProducer {
	events := make(chan model.WorkplaceEvent, 512)
	workerPool := workerpool.New(2)

	producer := NewKafkaProducer(
		kafkaWorkers,
		fixture.Sender,
		events,
		workerPool,
		fixture.Repo)

	return KafkaProducer{
		events:     events,
		workerPool: workerPool,
		producer:   producer,
	}
}

func (k *KafkaProducer) Close() {
	k.producer.Close()
}
