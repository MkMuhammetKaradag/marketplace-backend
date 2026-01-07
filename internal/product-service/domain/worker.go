package domain

import "github.com/google/uuid"

type Worker interface {
	EnqueueImageUpload(payload UploadImageTaskPayload) error
	EnqueueTrackView(payload TrackProductViewPayload) error
}

type UploadImagePayload struct {
	ProductID uuid.UUID
	ImageData []byte
	FileName  string
	IsMain    bool
	SortOrder int
}
type UploadImageTaskPayload struct {
	ProductID uuid.UUID `json:"product_id"`
	ImageData []byte    `json:"image_data"`
	FileName  string    `json:"file_name"`
	IsMain    bool      `json:"is_main"`
	SortOrder int       `json:"sort_order"`
}

type TrackProductViewPayload struct {
	UserID    uuid.UUID `json:"user_id"`
	Embedding []float32 `json:"embedding"`
	ProductID uuid.UUID `json:"product_id"`
}
