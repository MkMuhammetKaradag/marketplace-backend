package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/order-service/domain"
	pp "marketplace/pkg/proto/Product"

	"github.com/google/uuid"
)

type CreateOrderUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID) error
}

type createOrderUseCase struct {
	basketRepository  domain.OrderRepository
	grpcProductClient domain.ProductClient
	grpcBasketClient  domain.BasketClient
	grpcPaymentClient domain.PaymentClient
}

func NewCreateOrderUseCase(basketRepository domain.OrderRepository, grpcProductClient domain.ProductClient, grpcBasketClient domain.BasketClient, grpcPaymentClient domain.PaymentClient) CreateOrderUseCase {
	return &createOrderUseCase{
		basketRepository:  basketRepository,
		grpcProductClient: grpcProductClient,
		grpcBasketClient:  grpcBasketClient,
		grpcPaymentClient: grpcPaymentClient,
	}
}

func (u *createOrderUseCase) Execute(ctx context.Context, userID uuid.UUID) error {

	basket, err := u.grpcBasketClient.GetBasket(ctx, userID.String())
	if err != nil || basket == nil {
		return fmt.Errorf("basket service error or basket empty: %v", err)
	}

	var ids []string
	for _, item := range basket.Items {
		ids = append(ids, item.ProductId)
	}

	productResponse, err := u.grpcProductClient.GetProductsByIds(ctx, ids)
	if err != nil || productResponse == nil {
		return fmt.Errorf("product service error: %v", err)
	}

	productMap := make(map[string]*pp.ProductResponse)
	for _, p := range productResponse.Products {
		productMap[p.Id] = p
	}

	orderID := uuid.New()
	var orderItems []domain.OrderItem
	var totalPrice float64

	for _, bItem := range basket.Items {
		productDetail, ok := productMap[bItem.ProductId]
		if !ok {
			return fmt.Errorf("product %s not found in product service", bItem.ProductId)
		}

		if productDetail.Stock < bItem.Quantity {
			return fmt.Errorf("not enough stock for product: %s", productDetail.Name)
		}

		itemPrice := productDetail.Price * float64(bItem.Quantity)
		totalPrice += itemPrice

		orderItems = append(orderItems, domain.OrderItem{
			ID:              uuid.New(),
			OrderID:         orderID,
			ProductID:       uuid.MustParse(bItem.ProductId),
			ProductName:     productDetail.Name,
			ProductImageUrl: productDetail.ImageUrl,
			UnitPrice:       productDetail.Price,
			Quantity:        int(bItem.Quantity),
			SellerID:        uuid.MustParse(productDetail.SellerId),
			Status:          domain.OrderPending,
		})
	}

	newOrder := &domain.Order{
		ID:         orderID,
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     domain.OrderPending,
		Items:      orderItems,
	}

	err = u.basketRepository.CreateOrder(ctx, newOrder)
	if err != nil {
		return fmt.Errorf("failed to save order: %v", err)
	}

	payment, err := u.grpcPaymentClient.CreatePaymentSession(ctx, orderID.String(), userID.String(), "email@test.com", totalPrice)
	if err != nil {
		return err
	}
	fmt.Println(payment)

	//  SON ADIM: Kafka'ya "OrderCreated" eventi at

	return nil

}
