// internal/order-service/grpc_client/basket_client.go

package grpc_client

import (
	"context"

	"marketplace/internal/order-service/domain"
	pb "marketplace/pkg/proto/basket"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var BasketServiceClient pb.BasketServiceClient
var basketConn *grpc.ClientConn

type basketClient struct {
	client pb.BasketServiceClient
	conn   *grpc.ClientConn
}

func NewBasketClient(grpcAddress string) (domain.BasketClient, error) {

	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &basketClient{
		client: pb.NewBasketServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *basketClient) GetBasket(ctx context.Context, id string) (*pb.BasketResponse, error) {

	req := &pb.GetBasketRequest{UserId: id}
	return c.client.GetBasket(ctx, req)
}

func (c *basketClient) Close() error {
	return c.conn.Close()
}
