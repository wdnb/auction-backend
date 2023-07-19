package middleware

import (
	"auction-website/conf"
	db "auction-website/database/connectors/mysql"
	rdb "auction-website/database/connectors/redis"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	rdb *redis.Client
	db  *sqlx.DB
}

func Init(c *conf.Config) *Config {
	return &Config{
		rdb: rdb.GetClient(c.Redis),
		db:  db.GetClient(c.Mysql),
	}
}
