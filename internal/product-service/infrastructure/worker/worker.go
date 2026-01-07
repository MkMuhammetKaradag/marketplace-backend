package worker

import (
	"encoding/json"
	"marketplace/internal/product-service/domain"

	"github.com/hibiken/asynq"
)

// Görev ismini dışarıdan erişilebilir yapmak için büyük harfle başlatabiliriz
const TaskUploadProductImage = "task:upload_product_image"

type Worker struct {
	client *asynq.Client
}

func NewWorker(client *asynq.Client) *Worker {
	return &Worker{
		client: client,
	}
}

// EnqueueImageUpload görevleri Redis kuyruğuna ekler
func (w *Worker) EnqueueImageUpload(payload domain.UploadImageTaskPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Görevi oluşturuyoruz
	task := asynq.NewTask(TaskUploadProductImage, data)

	// Kuyruğa gönderiyoruz
	_, err = w.client.Enqueue(task)
	return err
}
