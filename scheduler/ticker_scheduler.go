package scheduler

import (
	"auction-website/api/auction"
	"auction-website/conf"
	rdb "auction-website/database/connectors/redis"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

type TickerScheduler struct {
	ticker          *time.Ticker
	c               *conf.Config
	rdb             *redis.Client
	redisAucLockKey string
}

func NewTickerScheduler(redisAucLockKey string, c *conf.Config) *TickerScheduler {
	return &TickerScheduler{
		redisAucLockKey: redisAucLockKey,
		c:               c,
		rdb:             rdb.GetClient(c.Redis),
	}
}

func (ts *TickerScheduler) AddTask() {
	ctx := context.Background()
	server := auction.NewService(ts.c)
	go func() {
		// acquire the global lock
		lock, err := ts.rdb.SetNX(ctx, ts.redisAucLockKey, "locked", time.Minute).Result()
		if err != nil {
			zap.L().Error("NSQ:Failed to acquire lock", zap.Error(err))
		} else if lock {
			// if lock is acquired, execute checkAuctions function
			err = server.CheckAuctions()
			if err != nil {
				zap.L().Error("NSQ:Failed to check auctions", zap.Error(err))
			} else {
				zap.L().Debug("NSQ:Auction check finished")
			}
			// release the global lock
			_, err = ts.rdb.Del(ctx, ts.redisAucLockKey).Result()
			if err != nil {
				zap.L().Error("NSQ:Failed to release lock", zap.Error(err))
			}
		} else {
			zap.L().Info("NSQ:Another instance is already checking auctions")
		}
	}()
}

func (ts *TickerScheduler) Start() {
	ts.ticker = time.NewTicker(viper.GetDuration("app.auction_run_interval") * time.Second)
	go func() {
		for range ts.ticker.C {
			ts.AddTask()
		}
	}()
}

func (ts *TickerScheduler) Stop() {
	ts.ticker.Stop()
}
