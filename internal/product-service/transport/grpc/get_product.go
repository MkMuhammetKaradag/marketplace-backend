package grpc_transport

import (
	"context"
	"marketplace/internal/product-service/domain"
	pb "marketplace/pkg/proto/Product"

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

func (h *AuthGrpcHandler) GetProductsByIds(ctx context.Context, req *pb.GetProductsByIdsRequest) (*pb.GetProductsByIdsResponse, error) {
	productIDs := req.GetIds()
	if len(productIDs) == 0 {
		return nil, nil
	}
	var productResponses []*pb.ProductResponse
	for _, productIDstr := range productIDs {
		productID, err := uuid.Parse(productIDstr)
		if err != nil {
			return nil, err
		}
		product, err := h.SessionRepo.GetProduct(ctx, productID, nil)
		if err != nil {
			return nil, err
		}
		IsActive := product.Status == "active"
		productResponses = append(productResponses, &pb.ProductResponse{
			Id:       product.ID.String(),
			Name:     product.Name,
			Price:    product.Price,
			Stock:    int32(product.StockCount),
			IsActive: IsActive,
		})
	}

	return &pb.GetProductsByIdsResponse{
		Products: productResponses,
	}, nil
}
