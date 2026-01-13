// internal/seller-service/domain/messasing.go
package domain

import (
	"context"
	pb "marketplace/pkg/proto/events"
)

type Messaging interface {
	Close() error
	PublishMessage(ctx context.Context, msg *pb.Message) error
}
