package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"snack-store-api/internal/constants"
	"snack-store-api/internal/messages"
	"snack-store-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/ulule/limiter/v3"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

func NewRateLimiter(viper *viper.Viper, redis *redis.Client) gin.HandlerFunc {
	rateStr := strings.TrimSpace(viper.GetString("RATE_LIMIT"))
	return NewRateLimiterWithRate(rateStr, redis)
}

func NewRateLimiterWithRate(rateStr string, redis *redis.Client) gin.HandlerFunc {
	rate := parseRate(rateStr)

	store, err := redisstore.NewStoreWithOptions(redis, limiter.StoreOptions{
		Prefix:   "rate_limiter",
		MaxRetry: 3,
	})
	if err != nil {
		panic(err)
	}

	limiterInstance := limiter.New(store, rate)

	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		if ip == "::1" {
			ip = "127.0.0.1"
		}

		limiterCtx, err := limiterInstance.Get(ctx.Request.Context(), ip)
		if err != nil {
			utils.HandleHTTPError(ctx, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err))
			return
		}

		ctx.Header(constants.RateLimitLimitHeader, fmt.Sprintf("%d", limiterCtx.Limit))
		ctx.Header(constants.RateLimitRemainingHeader, fmt.Sprintf("%d", limiterCtx.Remaining))
		ctx.Header(constants.RateLimitResetHeader, fmt.Sprintf("%d", limiterCtx.Reset))

		if limiterCtx.Reached {
			utils.HandleHTTPError(ctx, utils.Error(messages.TooManyRequests, http.StatusTooManyRequests, nil))
			return
		}

		ctx.Next()
	}
}

func parseRate(rateStr string) limiter.Rate {
	trimmed := strings.TrimSpace(rateStr)
	if trimmed == "" {
		trimmed = constants.DefaultRateLimit
	}

	rate, err := limiter.NewRateFromFormatted(trimmed)
	if err != nil {
		return limiter.Rate{
			Period: time.Minute,
			Limit:  60,
		}
	}

	return rate
}
