package services

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ImageServiceCloudinary struct {
	Cloudinary *cloudinary.Cloudinary
}

func NewImageServiceCloudinary(cloudinary *cloudinary.Cloudinary) *ImageServiceCloudinary {
	return &ImageServiceCloudinary{Cloudinary: cloudinary}
}

func (serv *ImageServiceCloudinary) Upload(file multipart.File) (string, error) {
	resp, err := serv.Cloudinary.Upload.Upload(context.TODO(), file, uploader.UploadParams{})
	if err != nil {
		return "", err
	}

	return resp.URL, nil
}
