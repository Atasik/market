package logger

import (
	"context"

	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
	ctx    context.Context
}

func NewZapLogger(loggerType string, ctx context.Context) (*ZapLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{logger: logger, ctx: ctx}, nil
}

func (l *ZapLogger) Debug(msg string, fields map[string]interface{}) {
	l.addContextCommonFields(fields)

	l.logger.Debug(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Info(msg string, fields map[string]interface{}) {
	l.addContextCommonFields(fields)

	l.logger.Info(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Warn(msg string, fields map[string]interface{}) {
	l.addContextCommonFields(fields)

	l.logger.Warn(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Error(msg string, fields map[string]interface{}) {
	l.addContextCommonFields(fields)

	l.logger.Error(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Fatal(msg string, fields map[string]interface{}) {
	l.addContextCommonFields(fields)

	l.logger.Fatal(msg, zap.Any("args", fields))
}

func (l *ZapLogger) addContextCommonFields(fields map[string]interface{}) {
	if l.ctx != nil {
		for k, v := range l.ctx.Value("commonFields").(map[string]interface{}) {
			if _, ok := fields[k]; !ok {
				fields[k] = v
			}
		}
	}
}

func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}
