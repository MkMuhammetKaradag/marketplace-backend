package messaginghandler

import (
	"marketplace/internal/payment-service/domain"
	eventsProto "marketplace/pkg/proto/events"
)

type Handlers struct {
}

func NewMessageHandlers() *Handlers {
	return &Handlers{}
}

func SetupMessageHandlers() map[eventsProto.MessageType]domain.MessageHandler {

	return map[eventsProto.MessageType]domain.MessageHandler{}
}
