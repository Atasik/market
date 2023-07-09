package service

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
)

type Config struct {
	Cloud  string
	Key    string
	Secret string
}

func NewCloudinary(cfg Config) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(cfg.Cloud, cfg.Key, cfg.Secret)
	if err != nil {
		return nil, err
	}

	_, err = cld.Admin.Ping(context.TODO())
	if err != nil {
		return nil, err
	}

	return cld, nil
}
