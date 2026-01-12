package domain

import (
	"context"

	"github.com/google/uuid"
)

type ProductRepository interface {
	AddSeller(ctx context.Context, sellerID uuid.UUID, userID uuid.UUID, status string) error
	AddUser(ctx context.Context, userID uuid.UUID, username string, email string) error
	GetSellerIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	CreateProduct(ctx context.Context, p *Product) (uuid.UUID, error)
	UpdateProductEmbedding(ctx context.Context, id uuid.UUID, embedding []float32) error
	SaveImagesAndUpdateStatus(ctx context.Context, productID uuid.UUID, images []ProductImage) error
	CreateCategory(ctx context.Context, c *Category) error
	RemoveInteraction(ctx context.Context, userID uuid.UUID, productID uuid.UUID, interactionType string) error
	TrackProductView(ctx context.Context, userID uuid.UUID, productEmbedding []float32) error
	GetRecommendedProducts(ctx context.Context, userID uuid.UUID, limit int) ([]*Product, error)

	GetProduct(ctx context.Context, id uuid.UUID, userID *uuid.UUID) (*Product, error)
	AddInteraction(ctx context.Context, userID uuid.UUID, productID uuid.UUID, interactionType string) error
	SearchProducts(ctx context.Context, queryVector []float32, req SearchProductsParams) ([]*Product, error)
	ToggleFavorite(ctx context.Context, userID, productID uuid.UUID) (bool, error)
	CheckLocalUserExists(ctx context.Context, userID uuid.UUID) (bool, error)
	GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]*FavoriteItem, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error)
	UpdateProduct(ctx context.Context, p *Product) error
	SoftDeleteProduct(ctx context.Context, productID uuid.UUID) error
	SoftDeleteAllProductImages(ctx context.Context, productID uuid.UUID) error
	Close() error
}
