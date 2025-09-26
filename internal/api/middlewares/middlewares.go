package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

// ZerologMiddleware cria um middleware que usa zerolog para logging das requisições
func ZerologMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Processar requisição
		c.Next()

		// Log da requisição
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		// Determinar nível do log baseado no status code
		var logEvent *zerolog.Event
		if statusCode >= 500 {
			logEvent = zlog.Error()
		} else if statusCode >= 400 {
			logEvent = zlog.Warn()
		} else {
			logEvent = zlog.Info()
		}

		logEvent.
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("ip", clientIP).
			Str("user_agent", c.Request.UserAgent()).
			Msg("HTTP Request")
	}
}
