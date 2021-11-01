package retranslator

import (
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
}

func TestWithoutErrors(t *testing.T) {
	t.Parallel()

	fixture := mocks.Setup(t)
	defer fixture.TearDown()

	fixture.SetupSuccessfulLockMock(4, &eventsData, 1)

	fixture.SetupSuccessfulSendMock(gomock.Any(), 4)
	fixture.SetupSuccessfulRemoveMock(4)

	startRetranslator(&fixture)
}

func TestKafkaAndDBUpdErrors(t *testing.T) {
	t.Parallel()

	fixture := mocks.Setup(t)
	defer fixture.TearDown()

	fixture.SetupSuccessfulLockMock(4, &eventsData, 1)

	fixture.SetupSuccessfulSendMock(gomock.Any(), 4)
	fixture.SetupUnsuccessfulRemoveMock(4)

	startRetranslator(&fixture)
}

func TestLockErrors(t *testing.T) {
	t.Parallel()

	fixture := mocks.Setup(t)
	defer fixture.TearDown()

	fixture.SetupUnsuccessfulLockMock(4, 1)

	startRetranslator(&fixture)
}

func startRetranslator(fixture *mocks.RetranslatorMockFixture) {
	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  1,
		ConsumeSize:    4,
		ConsumeTimeout: 3 * time.Second,
		ProducerCount:  1,
		WorkerCount:    2,
		Repo:           fixture.Repo,
		Sender:         fixture.Sender,
	}

	retranslator := NewRetranslator(cfg)
	retranslator.Start()

	fixture.Wg.Wait()

	retranslator.Close()
}
