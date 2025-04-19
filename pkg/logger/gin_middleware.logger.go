package logger

import (
	"time"

	"github.com/gin-gonic/gin"
)

func GinZapMiddleware(l Logger, serviceName, env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		fields := []Field{
			{Key: "service_name", Value: serviceName},
			{Key: "env", Value: env},
			{Key: "method", Value: c.Request.Method},
			{Key: "path", Value: c.FullPath()},
			{Key: "status", Value: c.Writer.Status()},
			{Key: "duration", Value: duration.String()},
			{Key: "client_ip", Value: c.ClientIP()},
			{Key: "user_agent", Value: c.Request.UserAgent()},
		}

		fields = append(fields, TraceFieldsFromContext(c.Request.Context())...)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				fields = append(fields, Field{Key: "error", Value: e.Error()})
			}
			l.Error("request completed with errors", fields...)
		} else {
			l.Info("request completed", fields...)
		}
	}
}
