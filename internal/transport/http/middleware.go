package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jcserv/go-api-template/internal/utils/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func middlewareLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	l, _ := config.Build()
	return l
}

func LogIncomingRequests() mux.MiddlewareFunc {
	logger := middlewareLogger()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			wrapped := wrapResponseWriter(w)

			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
				r.Header.Set("X-Request-ID", requestID)
			}
			w.Header().Set("X-Request-ID", requestID)
			reqLogger := log.WithRequest(logger, r)

			defer func() {
				if err := recover(); err != nil {
					wrapped.WriteHeader(http.StatusInternalServerError)
					reqLogger.Error(fmt.Sprintf("panic: %v", err))
					return
				}
			}()

			next.ServeHTTP(wrapped, r)

			duration := time.Since(startTime)
			status := wrapped.Status()
			if status == 0 {
				status = 200
			}

			logMsg := fmt.Sprintf("%d | %s %s (%v)",
				status,
				r.Method,
				r.URL.Path,
				duration.Round(time.Millisecond),
			)

			if status >= 500 {
				reqLogger.Error(logMsg)
			} else if status >= 400 {
				reqLogger.Warn(logMsg)
			} else {
				reqLogger.Info(logMsg)
			}
		})
	}
}
