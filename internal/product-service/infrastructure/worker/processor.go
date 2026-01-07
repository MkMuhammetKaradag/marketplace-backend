// internal/product-service/infrastructure/worker/processor.go
package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"marketplace/internal/product-service/domain"

	"github.com/hibiken/asynq"
)

type TaskProcessor struct {
	server        *asynq.Server
	repo          domain.ProductRepository
	cloudinarySvc domain.ImageService
}

func NewTaskProcessor(redisOpt asynq.RedisClientOpt, repo domain.ProductRepository, cloudinarySvc domain.ImageService) *TaskProcessor {
	// Concurrency: 5 -> Aynı anda en fazla 5 resim işlensin (Sunucuyu yormamak için)
	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 5,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
		},
	})

	return &TaskProcessor{
		server:        server,
		repo:          repo,
		cloudinarySvc: cloudinarySvc,
	}
}
func (p *TaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	// "task:upload_product_image" isminde bir görev gelirse ProcessUploadTask'ı çalıştır
	mux.HandleFunc(TaskUploadProductImage, p.ProcessUploadTask)

	log.Println("Worker Processor başlatılıyor...")
	return p.server.Run(mux)
}

func (p *TaskProcessor) ProcessUploadTask(ctx context.Context, t *asynq.Task) error {
	var payload domain.UploadImageTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	// 1. Cloudinary'ye yükle (Önceki adımda yazdığımız ImageFromBytes metodunu kullanıyor)
	url, err := p.cloudinarySvc.UploadImageFromBytes(ctx, payload.ImageData, domain.UploadOptions{
		Folder:   "products",
		PublicID: fmt.Sprintf("%s_%d", payload.ProductID, payload.SortOrder),
	})
	if err != nil {
		log.Printf("Cloudinary hatası: %v", err)
		return err // Hata dönerse Asynq otomatik olarak tekrar dener (Retry)
	}

	// 2. Veritabanına kaydet
	img := []domain.ProductImage{{
		ImageURL:  url,
		IsMain:    payload.IsMain,
		SortOrder: payload.SortOrder,
	}}

	return p.repo.SaveImagesAndUpdateStatus(ctx, payload.ProductID, img)
}
