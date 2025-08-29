package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/Luawig/neoneuro/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const headerRequestID = "X-Request-ID"

// GinLogger logs each HTTP request with zap.
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// capture status/size after next handlers
		c.Next()

		lat := time.Since(start)
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.String("uri", c.Request.RequestURI),
			zap.String("ip", c.ClientIP()),
			zap.String("ua", c.Request.UserAgent()),
			zap.Int("size", c.Writer.Size()),
			zap.Duration("latency", lat),
		}
		if rid := c.Writer.Header().Get(headerRequestID); rid != "" {
			fields = append(fields, zap.String("rid", rid))
		}

		// log level by status
		switch {
		case c.Writer.Status() >= 500:
			logger.L().Error("http", fields...)
		case c.Writer.Status() >= 400:
			logger.L().Warn("http", fields...)
		default:
			logger.L().Info("http", fields...)
		}
	}
}

// GinRecovery recovers from panics and logs the stack with zap.
// If `APP_ENV=prod`, it hides the stack trace in response.
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.L().Error("panic",
					zap.Any("error", rec),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
		}()
		c.Next()
	}
}

// BodyDump copies request body up to max bytes for debugging.
// Use ONLY in dev; attach before handlers.
func BodyDump(max int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil {
			c.Next()
			return
		}
		var buf bytes.Buffer
		_, _ = io.CopyN(&buf, c.Request.Body, max)
		c.Request.Body = io.NopCloser(io.MultiReader(bytes.NewReader(buf.Bytes()), c.Request.Body))
		c.Set("req_body_snippet", buf.String())
		c.Next()
	}
}
