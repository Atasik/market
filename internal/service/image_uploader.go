package service

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ImageServiceCloudinary struct {
	cloudinary *cloudinary.Cloudinary
}

type ImageData struct {
	ImageURL string
	ImageID  string
}

func NewImageServiceCloudinary(cloudinary *cloudinary.Cloudinary) *ImageServiceCloudinary {
	return &ImageServiceCloudinary{cloudinary: cloudinary}
}

func (s *ImageServiceCloudinary) Upload(ctx context.Context, file multipart.File) (ImageData, error) {
	resp, err := s.cloudinary.Upload.Upload(ctx, file, uploader.UploadParams{})
	if err != nil {
		return ImageData{}, err
	}
	return ImageData{ImageURL: resp.URL, ImageID: resp.PublicID}, nil
}

func (s *ImageServiceCloudinary) Delete(ctx context.Context, imageID string) error {
	if _, err := s.cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: imageID}); err != nil {
		return err
	}
	return nil
}
