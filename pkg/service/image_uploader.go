package service

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Image interface {
	Upload(ctx context.Context, file multipart.File) (ImageData, error)
	Delete(ctx context.Context, imageID string) error
}

type ImageServiceCloudinary struct {
	Cloudinary *cloudinary.Cloudinary
}

type ImageData struct {
	ImageURL string
	ImageID  string
}

func NewImageServiceCloudinary(cloudinary *cloudinary.Cloudinary) *ImageServiceCloudinary {
	return &ImageServiceCloudinary{Cloudinary: cloudinary}
}

func (s *ImageServiceCloudinary) Upload(ctx context.Context, file multipart.File) (ImageData, error) {
	resp, err := s.Cloudinary.Upload.Upload(ctx, file, uploader.UploadParams{})
	if err != nil {
		return ImageData{}, err
	}

	return ImageData{ImageURL: resp.URL, ImageID: resp.PublicID}, nil
}

func (s *ImageServiceCloudinary) Delete(ctx context.Context, imageID string) error {
	_, err := s.Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: imageID})
	if err != nil {
		return err
	}
	return nil
}

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = cld.Admin.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return cld, nil
}
