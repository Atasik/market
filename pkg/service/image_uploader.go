package service

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Image interface {
	Upload(file multipart.File) (ImageData, error)
	Delete(imageID string) error
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

func (s *ImageServiceCloudinary) Upload(file multipart.File) (ImageData, error) {
	// добавить таймауты
	resp, err := s.Cloudinary.Upload.Upload(context.TODO(), file, uploader.UploadParams{})
	if err != nil {
		return ImageData{}, err
	}

	return ImageData{ImageURL: resp.URL, ImageID: resp.PublicID}, nil
}

func (s *ImageServiceCloudinary) Delete(imageID string) error {
	// добавить таймауты
	_, err := s.Cloudinary.Upload.Destroy(context.TODO(), uploader.DestroyParams{PublicID: imageID})
	if err != nil {
		return err
	}
	return nil
}
