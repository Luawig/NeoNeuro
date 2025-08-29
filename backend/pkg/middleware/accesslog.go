package middleware

import (
	"time"

	"github.com/Luawig/neoneuro/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
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
		if rid, ok := c.Get(HeaderRequestID); ok {
			fields = append(fields, zap.Any("rid", rid))
		}

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
