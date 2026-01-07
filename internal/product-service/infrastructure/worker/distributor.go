package worker

import (
	"encoding/json"
	"marketplace/internal/product-service/domain"

	"github.com/hibiken/asynq"
)

func (w *Worker) EnqueueTrackView(payload domain.TrackProductViewPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = w.client.Enqueue(asynq.NewTask(TaskTrackProductView, data, asynq.MaxRetry(5)))
	return err
}
