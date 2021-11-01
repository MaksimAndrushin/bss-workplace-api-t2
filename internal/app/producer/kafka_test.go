package producer

import (
	"errors"
	"github.com/gammazero/workerpool"
	"github.com/golang/mock/gomock"
	"github.com/ozonmp/omp-demo-api/internal/mocks"
	"github.com/ozonmp/omp-demo-api/internal/model"
	"testing"
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

	fixture.SetupSuccessfulSendMock(&event, 1)
	fixture.SetupSuccessfulRemoveMock(1)

	kafkaTestProducer.producer.Start()
	kafkaTestProducer.events <- event

	fixture.Wg.Wait()
}

func TestSendToKafkaUnsuccessful(t *testing.T) {
	t.Parallel()

	fixture := mocks.Setup(t)
	defer fixture.TearDown()

	kafkaTestProducer := NewTestKafkaProducer(2, fixture)
	defer kafkaTestProducer.Close()

	event := model.WorkplaceEvent{ID: 1, Type: 0, Status: 1, Entity: &model.Workplace{ID: 1}}

	fixture.SetupUnsuccessfulSendMock(&event, 1)
	fixture.SetupSuccessfulUnlockMock(1)

	kafkaTestProducer.producer.Start()
	kafkaTestProducer.events <- event

	fixture.Wg.Wait()
}

func TestSendToKafkaThreeSEventsAndOneUEvent(t *testing.T) {
	t.Parallel()

	fixture := mocks.Setup(t)
	defer fixture.TearDown()

	kafkaTestProducer := NewTestKafkaProducer(2, fixture)
	defer kafkaTestProducer.Close()

	sendingNum := 0

	fixture.Wg.Add(4)
	fixture.Sender.EXPECT().
		Send(gomock.Any()).
		DoAndReturn(func(_ *model.WorkplaceEvent) error {
			defer fixture.Wg.Done()

			fixture.Mu.Lock()
			defer fixture.Mu.Unlock()

			sendingNum++
			if sendingNum == 3 {
				fixture.SetupSuccessfulUnlockMock(1)
				return errors.New("Sending error")
			}

			return nil
		}).Times(4)

	fixture.SetupSuccessfulRemoveMock(3)

	kafkaTestProducer.producer.Start()

	for i := 0; i < 4; i++ {
		kafkaTestProducer.events <- model.WorkplaceEvent{ID: uint64(i), Type: 0, Status: 1, Entity: &model.Workplace{ID: 1}}
	}

	fixture.Wg.Wait()
}

func NewTestKafkaProducer(kafkaWorkers uint64, fixture mocks.RetranslatorMockFixture) KafkaProducer {
	events := make(chan model.WorkplaceEvent, 512)
	workerPool := workerpool.New(2)

	producer := NewKafkaProducer(
		kafkaWorkers,
		fixture.Sender,
		events,
		workerPool,
		fixture.Repo,
	)

	return KafkaProducer{
		events:     events,
		workerPool: workerPool,
		producer:   producer,
	}
}

func (k *KafkaProducer) Close() {
	k.producer.Close()
}
