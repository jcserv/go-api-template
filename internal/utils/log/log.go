package log

import (
	"context"
	"net/http"
	"strings"

	"github.com/jcserv/go-api-template/internal/utils/env"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func Init(isProd bool) *zap.Logger {
	var logger *zap.Logger
	if isProd {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	return logger
}

func GetLogger(ctx context.Context) *zap.Logger {
	isProd := env.GetString("ENVIRONMENT", "dev") == "production"

	if logger == nil {
		return Init(isProd)
	}
	return logger
}

func Error(ctx context.Context, msg string) {
	GetLogger(ctx).Error(msg)
}

func Fatal(ctx context.Context, msg string) {
	GetLogger(ctx).Fatal(msg)
}

func Info(ctx context.Context, msg string) {
	GetLogger(ctx).Info(msg)
}

func WithRequest(l *zap.Logger, r *http.Request) *zap.Logger {
	userAgent := r.UserAgent()
	os, device := parseUserAgent(userAgent)

	return l.With(
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("user_agent", userAgent),
		zap.String("request_id", r.Header.Get("X-Request-ID")),
		zap.String("user", r.Header.Get("X-User-ID")),
		zap.String("device", device),
		zap.String("os", os),
		zap.String("url", r.URL.String()),
	)
}

func parseUserAgent(userAgent string) (os, device string) {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "windows"):
		os = "windows"
	case strings.Contains(ua, "mac os"):
		os = "macos"
	case strings.Contains(ua, "linux"):
		os = "linux"
	case strings.Contains(ua, "android"):
		os = "android"
	case strings.Contains(ua, "ios"):
		os = "ios"
	default:
		os = "unknown"
	}

	switch {
	case strings.Contains(ua, "mobile"):
		device = "mobile"
	case strings.Contains(ua, "tablet"):
		device = "tablet"
	default:
		device = "desktop"
	}
	return
}
