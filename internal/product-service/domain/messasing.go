package domain

import (
	"context"
	pb "marketplace/pkg/proto/events"
)

type Messaging interface {
	Close() error
	PublishMessage(ctx context.Context, msg *pb.Message) error
}
type MessageHandler interface {
	Handle(ctx context.Context, msg *pb.Message) error
}
