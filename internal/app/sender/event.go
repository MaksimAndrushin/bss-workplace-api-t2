package sender

import (
	"github.com/ozonmp/bss-workplace-api/internal/model"
)

type EventSender interface {
	Send(subdomain *model.WorkplaceEvent) error
}
