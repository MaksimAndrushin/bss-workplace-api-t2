package workplaceOrder

import (
	"github.com/ozonmp/omp-demo-api/internal/model"
	"testing"
)

func TestWorkplaceOrderer(t *testing.T) {

	var workplaceOrderer = NewOrderer()

	workplaceOrderer.AddEvent(model.WorkplaceEvent{ID: 4, Type: 0, Status: 1, Entity: &model.Workplace{ID: 1}})
	workplaceOrderer.AddEvent(model.WorkplaceEvent{ID: 2, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}})
	workplaceOrderer.AddEvent(model.WorkplaceEvent{ID: 3, Type: 0, Status: 1, Entity: &model.Workplace{ID: 3}})
	workplaceOrderer.AddEvent(model.WorkplaceEvent{ID: 1, Type: 0, Status: 1, Entity: &model.Workplace{ID: 1}})
	workplaceOrderer.AddEvent(model.WorkplaceEvent{ID: 6, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}})
	workplaceOrderer.AddEvent(model.WorkplaceEvent{ID: 5, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}})

	workplaceOrderer.DeleteEvent(model.WorkplaceEvent{ID: 15, Type: 0, Status: 1, Entity: &model.Workplace{ID: 2}})

	var isOrdererValue = workplaceOrderer.IsEventAscOrdered(model.WorkplaceEvent{ID: 4, Type: 0, Status: 1, Entity: &model.Workplace{ID: 1}})
	var ordererExpeced = false
	if isOrdererValue != ordererExpeced {
		t.Fatalf("Error")
	}

	isOrdererValue = workplaceOrderer.IsEventAscOrdered(model.WorkplaceEvent{ID: 1, Type: 0, Status: 1, Entity: &model.Workplace{ID: 1}})
	ordererExpeced = true
	if isOrdererValue != ordererExpeced {
		t.Fatalf("Error")
	}
}

