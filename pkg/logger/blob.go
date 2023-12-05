package logger

import (
	"log"
)

type BlobLogger struct {
}

func NewBlobLogger() *BlobLogger {
	return &BlobLogger{}
}

func (l *BlobLogger) Debug(msg string, fields map[string]interface{}) {
	log.Println(msg, fields)
}

func (l *BlobLogger) Info(msg string, fields map[string]interface{}) {
	log.Println(msg, fields)
}

func (l *BlobLogger) Warn(msg string, fields map[string]interface{}) {
	log.Println(msg, fields)
}

func (l *BlobLogger) Error(msg string, fields map[string]interface{}) {
	log.Println(msg, fields)
}

func (l *BlobLogger) Fatal(msg string, fields map[string]interface{}) {
	log.Println(msg, fields)
}
