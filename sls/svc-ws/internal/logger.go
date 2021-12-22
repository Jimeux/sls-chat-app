package ws

import (
	"context"

	"go.uber.org/zap"
)

type contextKey int

const (
	KeyConnectionID contextKey = 0
	KeyRequestID    contextKey = 1
)

func withCtxValues(ctx context.Context, fields ...zap.Field) []zap.Field {
	// missing values will be empty strings
	connID, _ := ctx.Value(KeyConnectionID).(string)
	reqID, _ := ctx.Value(KeyRequestID).(string)
	return append(fields,
		zap.String("requestId", reqID),
		zap.String("connectionId", connID))
}

// Logger is a simple wrapper around a zap.Logger that provides
// convenience functions to reduce handler bloat.
type Logger struct {
	zap *zap.Logger
}

func NewLogger() *Logger {
	zapLogger, _ := zap.NewProduction()
	return &Logger{zap: zapLogger}
}

func NewNopLogger() *Logger {
	return &Logger{zap: zap.NewNop()}
}

func (l *Logger) Sync() {
	_ = l.zap.Sync()
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = withCtxValues(ctx, fields...)
	l.zap.Info(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, err error, fields ...zap.Field) {
	fields = append(withCtxValues(ctx, fields...), zap.Error(err))
	l.zap.Error(msg, fields...)
}
