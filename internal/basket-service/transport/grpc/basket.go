package grpc_transport

import (
	"context"
	"marketplace/internal/basket-service/domain"
	pb "marketplace/pkg/proto/basket"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type BasketGrpcHandler struct {
	pb.UnimplementedBasketServiceServer
	BasketRepo domain.BasketRedisRepository
}

func NewBasketGrpcHandler(repo domain.BasketRedisRepository) *BasketGrpcHandler {
	return &BasketGrpcHandler{
		BasketRepo: repo,
	}
}
func (h *BasketGrpcHandler) Register(gRPCServer *grpc.Server) {
	pb.RegisterBasketServiceServer(gRPCServer, h)
}

func (h *BasketGrpcHandler) GetBasket(ctx context.Context, req *pb.GetBasketRequest) (*pb.BasketResponse, error) {
	userIDstr := req.GetUserId()
	if userIDstr == "" {
		return nil, nil
	}
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return nil, nil
	}

	basket, err := h.BasketRepo.GetBasket(ctx, userID.String())
	if err != nil {
		return nil, nil
	}
	products := make([]*pb.BasketItem, len(basket.Items))
	var totalPrice float64
	for i, product := range basket.Items {
		totalPrice += product.Price * float64(product.Quantity)
		products[i] = &pb.BasketItem{
			ProductId: product.ProductID.String(),
			Quantity:  int32(product.Quantity),
			Price:     product.Price,
		}
	}

	return &pb.BasketResponse{
		UserId:     basket.UserID.String(),
		Items:      products,
		TotalPrice: totalPrice,
	}, nil
}
