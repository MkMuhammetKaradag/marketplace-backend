// internal/basket-service/grpc_client/product_client.go

package grpc_client

import (
	"context"

	"marketplace/internal/basket-service/domain"
	pb "marketplace/pkg/proto/Product"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ProductServiceClient pb.ProductServiceClient
var conn *grpc.ClientConn

type productClient struct {
	client pb.ProductServiceClient
	conn   *grpc.ClientConn
}

func NewProductClient(grpcAddress string) (domain.ProductClient, error) {

	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &productClient{
		client: pb.NewProductServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *productClient) GetProductForBasket(ctx context.Context, id string) (*pb.ProductResponse, error) {

	req := &pb.GetProductRequest{Id: id}
	return c.client.GetProductForBasket(ctx, req)
}

func (c *productClient) Close() error {
	return c.conn.Close()
}
