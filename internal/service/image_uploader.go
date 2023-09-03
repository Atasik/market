package service

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

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
