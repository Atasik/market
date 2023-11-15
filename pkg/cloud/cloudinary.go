package cloud

import (
	"context"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
)

const timeout = 5 * time.Second

func NewCloudinary(cloud, key, secret string) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(cloud, key, secret)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err = cld.Admin.Ping(ctx); err != nil {
		return nil, err
	}

	return cld, nil
}
