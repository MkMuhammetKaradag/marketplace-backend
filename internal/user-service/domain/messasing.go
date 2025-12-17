package domain

import (
	"context"
	"marketplace/pkg/messaging"
)

type Messaging interface {
	Close() error
	PublishMessage(ctx context.Context, msg *messaging.Message) error
}
type MessageHandler interface {
	Handle(ctx context.Context, msg *messaging.Message) error
}
