package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/order-service/domain"
	cp "marketplace/pkg/proto/common"
	eventsProto "marketplace/pkg/proto/events"
	pp "marketplace/pkg/proto/product"

	"github.com/google/uuid"
)

type CreateOrderUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID) (string, error)
}

type createOrderUseCase struct {
	orderRepository   domain.OrderRepository
	grpcProductClient domain.ProductClient
	grpcBasketClient  domain.BasketClient
	grpcPaymentClient domain.PaymentClient
	messaging         domain.Messaging
}

func NewCreateOrderUseCase(orderRepository domain.OrderRepository, grpcProductClient domain.ProductClient, grpcBasketClient domain.BasketClient, grpcPaymentClient domain.PaymentClient, messaging domain.Messaging) CreateOrderUseCase {
	return &createOrderUseCase{
		orderRepository:   orderRepository,
		grpcProductClient: grpcProductClient,
		grpcBasketClient:  grpcBasketClient,
		grpcPaymentClient: grpcPaymentClient,
		messaging:         messaging,
	}
}

func (u *createOrderUseCase) Execute(ctx context.Context, userID uuid.UUID) (string, error) {

	basket, err := u.grpcBasketClient.GetBasket(ctx, userID.String())
	if err != nil || basket == nil || len(basket.Items) == 0 {
		return "", fmt.Errorf("basket empty or error: %v", err)
	}

	orderID := uuid.New()
	var orderItems []domain.OrderItem
	var orderItemsData []*cp.OrderItemData

	for _, bItem := range basket.Items {
		pID := uuid.MustParse(bItem.ProductId)
		item := domain.OrderItem{
			ID:        uuid.New(),
			OrderID:   orderID,
			ProductID: pID,
			Quantity:  int(bItem.Quantity),
			Status:    domain.OrderPending,
		}
		orderItems = append(orderItems, item)

		orderItemsData = append(orderItemsData, &cp.OrderItemData{
			ProductId: bItem.ProductId,
			Quantity:  int32(bItem.Quantity),
		})
	}

	productResponse, err := u.grpcProductClient.ReserveStock(ctx, orderID.String(), orderItemsData)
	if err != nil {
		return "", fmt.Errorf("stock reservation failed: %w", err)
	}

	var totalPrice float64
	productInfoMap := make(map[string]*pp.ProductResponse)
	for _, p := range productResponse.Products {
		productInfoMap[p.Id] = p
	}
	for i := range orderItems {
		pIDstr := orderItems[i].ProductID.String()
		if info, ok := productInfoMap[pIDstr]; ok {
			orderItems[i].UnitPrice = info.Price
			orderItems[i].ProductName = info.Name
			orderItems[i].ProductImageUrl = info.ImageUrl
			orderItems[i].SellerID = uuid.MustParse(info.SellerId)
			totalPrice += info.Price * float64(orderItems[i].Quantity)
		}
	}

	newOrder := &domain.Order{
		ID:         orderID,
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     domain.OrderPending,
		Items:      orderItems,
	}

	if err := u.orderRepository.CreateOrder(ctx, newOrder); err != nil {

		return "", fmt.Errorf("failed to save order: %v", err)
	}

	payment, err := u.grpcPaymentClient.CreatePaymentSession(ctx, orderID.String(), userID.String(), "user@mail.com", totalPrice)
	if err != nil {
		return "", err
	}

	msg := &eventsProto.Message{
		Type:        eventsProto.MessageType_ORDER_CREATED,
		FromService: eventsProto.ServiceType_ORDER_SERVICE,
		Payload: &eventsProto.Message_OrderCreatedData{
			OrderCreatedData: &eventsProto.OrderCreatedData{
				OrderId:    orderID.String(),
				UserId:     userID.String(),
				TotalPrice: totalPrice,
				Items:      orderItemsData,
			},
		},
	}
	u.messaging.PublishMessage(ctx, msg)

	return payment.PaymentUrl, nil
}
