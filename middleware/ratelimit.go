package middleware

import (
	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	lgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	"time"
)

func (c *Config) Limit() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Hour,
		Limit:  10000,
	}
	store, err := redis.NewStoreWithOptions(c.rdb, limiter.StoreOptions{
		Prefix: "api_rate_limiter",
	})
	if err != nil {
		panic(err)
	}
	middleware := lgin.NewMiddleware(limiter.New(store, rate))
	return middleware
}
