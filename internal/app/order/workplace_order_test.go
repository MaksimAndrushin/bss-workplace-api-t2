package workplaceOrder

import (
	"github.com/ozonmp/omp-demo-api/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWorkplaceOrderer(t *testing.T) {
	workplaceOrderer := NewOrderer()

	generateEvents(workplaceOrderer)

	ordererExpeced := false
	require.Equal(t, ordererExpeced, workplaceOrderer.IsEventAscOrdered(generateEvent(4, 0, 1, 1)),
		"Error")

	ordererExpeced = true
	require.Equal(t, ordererExpeced, workplaceOrderer.IsEventAscOrdered(generateEvent(1, 0, 1, 1)),
		"Error")
}

func generateEvents(workplaceOrderer EventOrderer) {
	workplaceOrderer.AddEvent(generateEvent(4, 0, 1, 1))
	workplaceOrderer.AddEvent(generateEvent(2, 0, 1, 2))
	workplaceOrderer.AddEvent(generateEvent(3, 0, 1, 3))
	workplaceOrderer.AddEvent(generateEvent(1, 0, 1, 1))
	workplaceOrderer.AddEvent(generateEvent(6, 0, 1, 2))
	workplaceOrderer.AddEvent(generateEvent(5, 0, 1, 2))
}

func generateEvent(eventId uint64, eventType model.EventType, eventStatus model.EventStatus, workplaceId uint64) model.WorkplaceEvent {

	return model.WorkplaceEvent{
		ID:     eventId,
		Type:   eventType,
		Status: eventStatus,
		Entity: &model.Workplace{ID: workplaceId},
	}
}
