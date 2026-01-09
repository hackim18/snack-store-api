package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestLogger(logger *logrus.Logger, slowThreshold time.Duration) gin.HandlerFunc {
	if logger == nil {
		logger = logrus.New()
	}

	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		ctx.Next()

		latency := time.Since(start)
		status := ctx.Writer.Status()
		clientIP := ctx.ClientIP()
		userAgent := ctx.Request.UserAgent()

		entry := logger.WithFields(logrus.Fields{
			"status":     status,
			"method":     method,
			"path":       path,
			"latency_ms": latency.Milliseconds(),
			"client_ip":  clientIP,
			"user_agent": userAgent,
		})

		if len(ctx.Errors) > 0 {
			entry = entry.WithField("errors", ctx.Errors.String())
		}

		if slowThreshold > 0 && latency >= slowThreshold {
			entry.Warn("slow request")
			return
		}

		entry.Info("request")
	}
}
