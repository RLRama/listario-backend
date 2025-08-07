package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RLRama/listario-backend/logger"
	"github.com/iris-contrib/middleware/throttler"
	"github.com/kataras/iris/v12"
	"github.com/rs/zerolog"
	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"
)

type varyByRemoteAddr struct{}

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

func (v *varyByRemoteAddr) Key(r *http.Request) string {
	return r.RemoteAddr
}

func NewRateLimiter(maxRatePerMinute int, maxBurst int) iris.Handler {
	store, err := memstore.NewCtx(65536)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create rate limiter memory store")
	}

	quota := throttled.RateQuota{
		MaxRate:  throttled.PerMin(maxRatePerMinute),
		MaxBurst: maxBurst,
	}

	rateLimiter, err := throttled.NewGCRARateLimiterCtx(store, quota)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create rate limiter")
	}

	limiter := &throttler.RateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &varyByRemoteAddr{},
		DeniedHandler: func(ctx iris.Context) {
			ctx.StatusCode(http.StatusTooManyRequests)
			ctx.JSON(iris.Map{
				"message": fmt.Sprintf("Too many requests, please try again in %s", ctx.GetHeader("Retry-After")),
				"status":  iris.StatusTooManyRequests,
			})
		},
	}

	return limiter.RateLimit
}
