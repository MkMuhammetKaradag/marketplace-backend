package grpc_transport

import (
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"
	pb "marketplace/pkg/proto/product"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type ProductGrpcHandler struct {
	pb.UnimplementedProductServiceServer
	productRepo domain.ProductRepository
}

func NewProductGrpcHandler(repo domain.ProductRepository) *ProductGrpcHandler {
	return &ProductGrpcHandler{
		productRepo: repo,
	}
}
func (h *ProductGrpcHandler) Register(gRPCServer *grpc.Server) {
	pb.RegisterProductServiceServer(gRPCServer, h)
}

func (h *ProductGrpcHandler) GetProductForBasket(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	productIDstr := req.GetId()
	if productIDstr == "" {
		return nil, nil
	}
	productID, err := uuid.Parse(productIDstr)
	if err != nil {
		return nil, nil
	}

	product, err := h.productRepo.GetProduct(ctx, productID, nil)
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

func (h *ProductGrpcHandler) GetProductsByIds(ctx context.Context, req *pb.GetProductsByIdsRequest) (*pb.GetProductsByIdsResponse, error) {
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
		product, err := h.productRepo.GetProduct(ctx, productID, nil)
		if err != nil {
			return nil, err
		}
		IsActive := product.Status == "active"
		imageUrl := ""
		if len(product.Images) > 0 {
			imageUrl = product.Images[0].ImageURL
		} else {
			imageUrl = "https://placehold.co/150"
		}
		productResponses = append(productResponses, &pb.ProductResponse{
			Id:       product.ID.String(),
			Name:     product.Name,
			Price:    product.Price,
			ImageUrl: imageUrl,
			SellerId: product.SellerID.String(),
			Stock:    int32(product.StockCount),
			IsActive: IsActive,
		})
	}

	return &pb.GetProductsByIdsResponse{
		Products: productResponses,
	}, nil
}
func (h *ProductGrpcHandler) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
	items := req.GetItems()
	if len(items) == 0 {
		return nil, fmt.Errorf("items cannot be empty")
	}

	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, fmt.Errorf("invalid order id: %w", err)
	}

	var reserveItems []domain.OrderItemReserve
	for _, item := range items {
		pID, err := uuid.Parse(item.ProductId)
		if err != nil {
			return nil, fmt.Errorf("invalid product id: %w", err)
		}
		reserveItems = append(reserveItems, domain.OrderItemReserve{
			ProductID: pID,
			Quantity:  int(item.Quantity),
		})
	}

	reservedProducts, err := h.productRepo.ReserveStocks(ctx, orderID, reserveItems)
	if err != nil {
		fmt.Printf("Reservation failed for order %s: %v\n", orderID, err)
		return nil, err
	}

	var protoProducts []*pb.ProductResponse
	for _, p := range reservedProducts {
		protoProducts = append(protoProducts, &pb.ProductResponse{
			Id:       p.ID.String(),
			Name:     p.Name,
			Price:    p.Price,
			ImageUrl: p.ImageURL,
			SellerId: p.SellerID.String(),
		})
	}

	return &pb.ReserveStockResponse{
		Success:  true,
		Products: protoProducts,
	}, nil
}
