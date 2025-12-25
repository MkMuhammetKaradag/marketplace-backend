// internal/user-service/infrastructure/cloudinary.go
package infrastructure

import (
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
)

func NewCloudinary(cloudName, apiKey, apiSecret string) (*cloudinary.Cloudinary, error) {
	fmt.Println("cloudinary config: ",
		cloudName,
		apiKey,
		apiSecret,
	)
	return cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
}
