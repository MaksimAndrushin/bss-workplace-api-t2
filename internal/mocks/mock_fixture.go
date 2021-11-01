package mocks

import (
	"github.com/golang/mock/gomock"
	"testing"
)

type RetranslatorMockFixture struct {
	Ctrl   *gomock.Controller
	Repo   *MockEventRepo
	Sender *MockEventSender
}

func Setup(t *testing.T) RetranslatorMockFixture {
	ctrl := gomock.NewController(t)

	return RetranslatorMockFixture{
		Ctrl:   ctrl,
		Repo:   NewMockEventRepo(ctrl),
		Sender: NewMockEventSender(ctrl),
	}
}

func (f RetranslatorMockFixture) TearDown() {
	f.Ctrl.Finish()
}
