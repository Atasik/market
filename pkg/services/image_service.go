package services

import "mime/multipart"

type ImageService interface {
	Upload(file multipart.File) (string, error)
}
