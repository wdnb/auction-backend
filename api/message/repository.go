package message

import (
	"auction-website/conf"
	mdb "auction-website/database/connectors/mongodb"
	rdb "auction-website/database/connectors/redis"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	rdb *redis.Client
	mdb *mongo.Client
}

func NewRepository(c *conf.Config) *Repository {
	return &Repository{
		rdb: rdb.GetClient(c.Redis),
		mdb: mdb.GetClient(c.MongoDB),
	}
}
