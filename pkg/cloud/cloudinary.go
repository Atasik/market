package cloud

import (
	"context"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
)

func NewCloudinary(cloud, key, secret string) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(cloud, key, secret)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = cld.Admin.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return cld, nil
}
