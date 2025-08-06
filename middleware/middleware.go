package middleware

import (
	"time"

	"github.com/RLRama/listario-backend/logger"
	"github.com/kataras/iris/v12"
	"github.com/rs/zerolog"
)

func RequestLogger() iris.Handler {
	return func(ctx iris.Context) {
		start := time.Now()
		ctx.Next()

		latency := time.Since(start)
		statusCode := ctx.ResponseWriter().StatusCode()
		path := ctx.Path()
		method := ctx.Method()
		clientIP := ctx.RemoteAddr()

		var event *zerolog.Event
		if statusCode >= 500 {
			event = logger.Error()
		} else if statusCode >= 400 {
			event = logger.Warn()
		} else {
			event = logger.Info()
		}

		event.
			Str("method", method).
			Str("path", path).
			Str("client_ip", clientIP).
			Int("status_code", statusCode).
			Dur("latency", latency).
			Msgf("Request %s %s completed with status %d in %s",
				method, path, statusCode, latency)
	}
}