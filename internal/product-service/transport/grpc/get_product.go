package grpc_transport

import (
	"context"
	"marketplace/internal/product-service/domain"
	pb "marketplace/pkg/proto/product"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type AuthGrpcHandler struct {
	pb.UnimplementedProductServiceServer
	SessionRepo domain.ProductRepository
}

func NewAuthGrpcHandler(repo domain.ProductRepository) *AuthGrpcHandler {
	return &AuthGrpcHandler{
		SessionRepo: repo,
	}
}
func (h *AuthGrpcHandler) Register(gRPCServer *grpc.Server) {
	pb.RegisterProductServiceServer(gRPCServer, h)
}

func (h *AuthGrpcHandler) GetProductForBasket(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	productIDstr := req.GetId()
	if productIDstr == "" {
		return nil, nil
	}
	productID, err := uuid.Parse(productIDstr)
	if err != nil {
		return nil, nil
	}

	product, err := h.SessionRepo.GetProduct(ctx, productID, nil)
	if err != nil {
		return nil, nil
	}

	IsActive := product.Status == "active"
	return &pb.ProductResponse{
		Id:       product.ID.String(),
		Name:     product.Name,
		Price:    product.Price,
		Stock:    int32(product.StockCount),
		IsActive: IsActive,
	}, nil
}
