package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type BlobLogger struct {
	logger *zap.Logger
	ctx    context.Context
}

func NewBlobLogger() (*BlobLogger, error) {
	return &BlobLogger{}, nil
}

func (l *BlobLogger) Debug(msg string, fields map[string]interface{}) {
	fmt.Println(msg)
}

func (l *BlobLogger) Info(msg string, fields map[string]interface{}) {
	fmt.Println(msg)
}

func (l *BlobLogger) Warn(msg string, fields map[string]interface{}) {
	fmt.Println(msg)
}

func (l *BlobLogger) Error(msg string, fields map[string]interface{}) {
	fmt.Println(msg)
}

func (l *BlobLogger) Fatal(msg string, fields map[string]interface{}) {
	fmt.Println(msg)
}
