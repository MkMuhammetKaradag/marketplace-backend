package worker

import (
	"encoding/json"
	"marketplace/internal/product-service/domain"

	"github.com/hibiken/asynq"
)

const TaskUploadProductImage = "task:upload_product_image"
const TaskTrackProductView = "task:track_product_view"

type Worker struct {
	client *asynq.Client
}

func NewWorker(client *asynq.Client) *Worker {
	return &Worker{
		client: client,
	}
}

func (w *Worker) EnqueueImageUpload(payload domain.UploadImageTaskPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskUploadProductImage, data, asynq.MaxRetry(5))

	_, err = w.client.Enqueue(task)
	return err
}


