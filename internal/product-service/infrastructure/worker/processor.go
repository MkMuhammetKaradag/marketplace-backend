// internal/product-service/infrastructure/worker/processor.go
package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"marketplace/internal/product-service/domain"
	"strings"

	"github.com/hibiken/asynq"
)

type TaskProcessor struct {
	server        *asynq.Server
	repo          domain.ProductRepository
	cloudinarySvc domain.ImageService
}

func NewTaskProcessor(redisOpt asynq.RedisClientOpt, repo domain.ProductRepository, cloudinarySvc domain.ImageService) *TaskProcessor {
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

	mux.HandleFunc(TaskUploadProductImage, p.ProcessUploadTask)
	mux.HandleFunc(TaskTrackProductView, p.ProcessTrackViewTask)

	log.Println("Worker Processor başlatılıyor...")
	return p.server.Run(mux)
}

func (p *TaskProcessor) ProcessUploadTask(ctx context.Context, t *asynq.Task) error {
	var payload domain.UploadImageTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	url, err := p.cloudinarySvc.UploadImageFromBytes(ctx, payload.ImageData, domain.UploadOptions{
		Folder:   "products",
		PublicID: fmt.Sprintf("%s_%d", payload.ProductID, payload.SortOrder),
	})
	if err != nil {
		if strings.Contains(err.Error(), "Invalid image file") {

			return fmt.Errorf("%w: %v", asynq.SkipRetry, err)
		}
		return err
	}

	img := []domain.ProductImage{{
		ImageURL:  url,
		IsMain:    payload.IsMain,
		SortOrder: payload.SortOrder,
	}}

	return p.repo.SaveImagesAndUpdateStatus(ctx, payload.ProductID, img)
}
func (p *TaskProcessor) ProcessTrackViewTask(ctx context.Context, t *asynq.Task) error {
	var payload domain.TrackProductViewPayload
	json.Unmarshal(t.Payload(), &payload)
	fmt.Println("Payload: ", payload)
	// 1. Ürün izlemeyi kaydet
	err := p.repo.TrackProductView(ctx, payload.UserID, payload.Embedding)
	if err != nil {
		return err
	}

	// 2. Etkileşimi ekle
	return p.repo.AddInteraction(ctx, payload.UserID, payload.ProductID, "view")
}
