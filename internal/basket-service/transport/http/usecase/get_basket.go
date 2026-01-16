package usecase

import (
	"context"
	"marketplace/internal/basket-service/domain"

	"github.com/google/uuid"
)

type GetBasketUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID) (*domain.BasketResponse, error)
}

type getBasketUseCase struct {
	basketRepository domain.BasketRedisRepository
}

func NewGetBasketUseCase(basketRepository domain.BasketRedisRepository) GetBasketUseCase {
	return &getBasketUseCase{
		basketRepository: basketRepository,
	}
}

func (u *getBasketUseCase) Execute(ctx context.Context, userID uuid.UUID) (*domain.BasketResponse, error) {

	basket, err := u.basketRepository.GetBasket(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	if basket == nil {
		return &domain.BasketResponse{UserID: userID.String(), Items: []domain.BasketItemResponse{}, TotalPrice: 0}, nil
	}

	var response domain.BasketResponse
	response.UserID = basket.UserID.String()

	var grandTotal float64 = 0

	for _, item := range basket.Items {
		subTotal := float64(item.Quantity) * item.Price
		grandTotal += subTotal

		response.Items = append(response.Items, domain.BasketItemResponse{
			ProductID: item.ProductID.String(),
			Name:      item.Name,
			Quantity:  item.Quantity,
			Price:     item.Price,
			ImageURL:  item.ImageURL,
			SubTotal:  subTotal,
		})
	}

	response.TotalPrice = grandTotal
	return &response, nil
}
