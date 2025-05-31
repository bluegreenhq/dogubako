package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	pkgcontext "github.com/bluegreenhq/dogubako/context"
	"github.com/bluegreenhq/dogubako/model"
)

const MaxLogMessageLen = 4000

type loggerImpl struct {
	logger *slog.Logger
	level  slog.Level
}

func NewLogger(isProduction bool) model.Logger {
	var handler slog.Handler

	logLevel := getLogLevel(isProduction)

	if isProduction {
		handler = slog.NewJSONHandler(os.Stdout, getLogHandlerOptions(logLevel))
	} else {
		handler = slog.NewTextHandler(os.Stdout, getLogHandlerOptions(logLevel))
	}

	return &loggerImpl{
		logger: slog.New(handler),
		level:  logLevel,
	}
}

func getLogLevel(isProduction bool) slog.Level {
	if isProduction {
		return slog.LevelInfo
	}

	return slog.LevelDebug
}

func getLogHandlerOptions(logLevel slog.Level) *slog.HandlerOptions {
	level := new(slog.LevelVar)
	level.Set(logLevel)

	return &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}
}

func (l *loggerImpl) loggerWithContext(ctx context.Context) *slog.Logger {
	var logger = l.logger

	if ctx != nil {
		if requestID, ok := ctx.Value(pkgcontext.ContextKeyRequestID).(string); ok {
			logger = l.logger.With(slog.String(string(pkgcontext.ContextKeyRequestID), requestID))
		}
	}

	return logger
}

func getMessage(format string, v ...any) string {
	var message = fmt.Sprintf(format, v...)

	// Truncate the log message if it exceeds max length
	if len(message) > MaxLogMessageLen {
		// Adding "..." to indicate truncation
		message = message[:MaxLogMessageLen] + "..."
	}

	return message
}

func Infof(ctx context.Context, format string, v ...any) {
	logger := pkgcontext.FromContext(ctx)
	logger.Infof(ctx, format, v...)
}

func Debugf(ctx context.Context, format string, v ...any) {
	logger := pkgcontext.FromContext(ctx)
	logger.Debugf(ctx, format, v...)
}

func Warnf(ctx context.Context, format string, v ...any) {
	logger := pkgcontext.FromContext(ctx)
	logger.Warnf(ctx, format, v...)
}

func Errorf(ctx context.Context, format string, v ...any) {
	logger := pkgcontext.FromContext(ctx)
	logger.Errorf(ctx, format, v...)
}

func Fatalf(ctx context.Context, format string, v ...any) {
	logger := pkgcontext.FromContext(ctx)
	logger.Fatalf(ctx, format, v...)
}

func (l *loggerImpl) Infof(ctx context.Context, format string, v ...any) {
	if l.level > slog.LevelInfo {
		return
	}

	var pcs [1]uintptr

	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelInfo, getMessage(format, v...), pcs[0])
	_ = l.loggerWithContext(ctx).Handler().Handle(context.Background(), r)
}

func (l *loggerImpl) Debugf(ctx context.Context, format string, v ...any) {
	if l.level > slog.LevelDebug {
		return
	}

	var pcs [1]uintptr

	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelDebug, getMessage(format, v...), pcs[0])
	_ = l.loggerWithContext(ctx).Handler().Handle(context.Background(), r)
}

func (l *loggerImpl) Warnf(ctx context.Context, format string, v ...any) {
	if l.level > slog.LevelWarn {
		return
	}

	var pcs [1]uintptr

	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelWarn, getMessage(format, v...), pcs[0])
	_ = l.loggerWithContext(ctx).Handler().Handle(context.Background(), r)
}

func (l *loggerImpl) Errorf(ctx context.Context, format string, v ...any) {
	if l.level > slog.LevelError {
		return
	}

	var pcs [1]uintptr

	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelError, getMessage(format, v...), pcs[0])
	_ = l.loggerWithContext(ctx).Handler().Handle(context.Background(), r)
}

func (l *loggerImpl) Fatalf(ctx context.Context, format string, v ...any) {
	if l.level > slog.LevelError {
		return
	}

	var pcs [1]uintptr

	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelError, getMessage(format, v...), pcs[0])
	_ = l.loggerWithContext(ctx).Handler().Handle(context.Background(), r)

	os.Exit(1)
}
