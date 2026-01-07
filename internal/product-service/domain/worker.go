package domain

import "github.com/google/uuid"

type Worker interface {
	EnqueueImageUpload(payload UploadImageTaskPayload) error
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
