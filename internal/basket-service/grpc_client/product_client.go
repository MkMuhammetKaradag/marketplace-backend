// internal/basket-service/grpc_client/product_client.go (Yeni dosya)

package grpc_client

import (
	"context"
	"log"
	"time"

	pb "marketplace/pkg/proto/Product"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ProductServiceClient pb.ProductServiceClient

var conn *grpc.ClientConn

func InitProductServiceClient(grpcAddress string) error {
	var err error

	conn, err = grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	ProductServiceClient = pb.NewProductServiceClient(conn)
	log.Printf("✅ Gateway, Product Servisine gRPC ile bağlandı: %s", grpcAddress)
	return nil
}

func CloseProductServiceClient() {
	if conn != nil {
		conn.Close()
	}
}

func GetProductForBasket(id string) (Product *pb.ProductResponse, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	req := &pb.GetProductRequest{Id: id}

	resp, err := ProductServiceClient.GetProductForBasket(ctx, req)

	if err != nil {
		log.Printf("🔒 gRPC doğrulama çağrısı başarısız: %v", err)
		return nil, err
	}

	return resp, nil
}
