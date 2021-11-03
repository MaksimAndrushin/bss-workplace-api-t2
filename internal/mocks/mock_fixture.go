package mocks

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/ozonmp/omp-demo-api/internal/model"
	"sync"
	"testing"
)

type RetranslatorMockFixture struct {
	Ctrl   *gomock.Controller
	Repo   *MockEventRepo
	Sender *MockEventSender
	Wg     sync.WaitGroup
	Mu     sync.Mutex
}

func Setup(t *testing.T) RetranslatorMockFixture {
	ctrl := gomock.NewController(t)

	return RetranslatorMockFixture{
		Ctrl:   ctrl,
		Repo:   NewMockEventRepo(ctrl),
		Sender: NewMockEventSender(ctrl),
		Wg:     sync.WaitGroup{},
	}
}

func (f *RetranslatorMockFixture) SetupUnsuccessfulSendMock(event *model.WorkplaceEvent, numTimes int) *gomock.Call {
	f.Wg.Add(numTimes)

	return f.Sender.EXPECT().
		Send(event).
		DoAndReturn(func(_ *model.WorkplaceEvent) error {
			defer f.Wg.Done()
			return errors.New("Sending error")
		}).Times(numTimes)
}

func (f *RetranslatorMockFixture) SetupSuccessfulSendMock(event interface{}, numTimes int) *gomock.Call {
	f.Wg.Add(numTimes)

	return f.Sender.EXPECT().
		Send(event).
		DoAndReturn(func(_ *model.WorkplaceEvent) error {
			defer f.Wg.Done()
			return nil
		}).Times(numTimes)
}

func (f *RetranslatorMockFixture) SetupUnsuccessfulRemoveMock(numTimes int) *gomock.Call {
	f.Wg.Add(numTimes)

	return f.Repo.EXPECT().
		Remove(gomock.Any()).
		DoAndReturn(func(_ []uint64) error {
			defer f.Wg.Done()
			return errors.New("Remove error")
		}).Times(numTimes)
}

func (f *RetranslatorMockFixture) SetupSuccessfulRemoveMock(numTimes int) *gomock.Call {
	f.Wg.Add(numTimes)

	return f.Repo.EXPECT().
		Remove(gomock.Any()).
		DoAndReturn(func(_ []uint64) error {
			defer f.Wg.Done()
			return nil
		}).Times(numTimes)
}

func (f *RetranslatorMockFixture) SetupSuccessfulUnlockMock(numTimes int) *gomock.Call {
	f.Wg.Add(numTimes)

	return f.Repo.EXPECT().Unlock(gomock.Any()).DoAndReturn(func(_ []uint64) error {
		defer f.Wg.Done()
		return nil
	}).Times(numTimes)
}

func (f *RetranslatorMockFixture) SetupUnsuccessfulLockMock(_ int, numTimes int) *gomock.Call {
	f.Wg.Add(numTimes)

	return f.Repo.EXPECT().
		Lock(uint64(4)).
		DoAndReturn(func(_ uint64) ([]model.WorkplaceEvent, error) {
			defer f.Wg.Done()
			return nil, errors.New("Lock error")
		}).Times(numTimes)
}

func (f *RetranslatorMockFixture) SetupSuccessfulLockMock(_ int, retEvents *[]model.WorkplaceEvent, numTimes int) *gomock.Call {
	f.Wg.Add(numTimes)

	return f.Repo.EXPECT().
		Lock(uint64(4)).
		DoAndReturn(func(_ uint64) ([]model.WorkplaceEvent, error) {
			defer f.Wg.Done()
			return *retEvents, nil
		}).Times(numTimes)
}


func (f RetranslatorMockFixture) TearDown() {
	f.Ctrl.Finish()
}
