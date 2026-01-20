// internal/order-service/grpc_client/product_client.go

package grpc_client

import (
	"context"

	"marketplace/internal/order-service/domain"
	pb "marketplace/pkg/proto/Product"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ProductServiceClient pb.ProductServiceClient
var productConn *grpc.ClientConn

type productClient struct {
	client pb.ProductServiceClient
	conn   *grpc.ClientConn
}

func NewProductClient(grpcAddress string) (domain.ProductClient, error) {

	productConn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &productClient{
		client: pb.NewProductServiceClient(productConn),
		conn:   productConn,
	}, nil
}

func (c *productClient) GetProductsByIds(ctx context.Context, ids []string) (*pb.GetProductsByIdsResponse, error) {

	req := &pb.GetProductsByIdsRequest{Ids: ids}
	return c.client.GetProductsByIds(ctx, req)
}

func (c *productClient) Close() error {
	return c.conn.Close()
}
