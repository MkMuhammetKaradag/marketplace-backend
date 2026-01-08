package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

type ToggleFavoriteUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error
}

type toggleFavoriteUseCase struct {
	repository domain.ProductRepository
	worker     domain.Worker
}

func NewToggleFavoriteUseCase(repository domain.ProductRepository, worker domain.Worker) ToggleFavoriteUseCase {
	return &toggleFavoriteUseCase{
		repository: repository,
		worker:     worker,
	}
}

func (c *toggleFavoriteUseCase) Execute(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	userExists, err := c.repository.CheckLocalUserExists(ctx, userID)
	if err != nil || !userExists {
		return fmt.Errorf("unauthorized or user sync pending")
	}
	return c.worker.EnqueueToggleFavorite(domain.FavoritePayload{
		UserID:    userID,
		ProductID: productID,
	})

}
