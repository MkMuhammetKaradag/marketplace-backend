package postgres

import (
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

func (r *Repository) SaveImagesAndUpdateStatus(ctx context.Context, productID uuid.UUID, images []domain.ProductImage) error {
	// 1. Transaction başlat
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("transaction başlatılamadı: %w", err)
	}

	// Hata oluşursa işlemleri geri al (Rollback)
	// Eğer tx.Commit() çağrılırsa Rollback bir şey yapmaz
	defer tx.Rollback()

	// 2. Resimleri tabloya ekle
	const insertImageQuery = `
		INSERT INTO product_images (product_id, image_url, is_main, sort_order) 
		VALUES ($1, $2, $3, $4)`

	for i, img := range images {
		// İlk resmi otomatik olarak ana resim (is_main = true) yapıyoruz
		isMain := (i == 0)

		_, err := tx.ExecContext(ctx, insertImageQuery, productID, img.ImageURL, isMain, i)
		if err != nil {
			return fmt.Errorf("resim kaydedilirken hata (URL: %s): %w", img.ImageURL, err)
		}
	}

	// 3. Ürünün durumunu güncelle (Draft'tan Active'e)
	const updateStatusQuery = `
		UPDATE products 
		SET status = 'active', updated_at = NOW() 
		WHERE id = $1`

	result, err := tx.ExecContext(ctx, updateStatusQuery, productID)
	if err != nil {
		return fmt.Errorf("ürün durumu güncellenemedi: %w", err)
	}

	// Etkilenen satır sayısını kontrol et (Yanlış ID gönderilmiş olabilir)
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("güncellenecek ürün bulunamadı (ID: %s)", productID)
	}

	// 4. Her şey tamamsa onaylıyoruz
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit hatası: %w", err)
	}

	return nil
}
