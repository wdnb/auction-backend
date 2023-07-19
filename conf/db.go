package conf

import (
	"auction-website/database/connectors/mongodb"
	"auction-website/database/connectors/mysql"
	rdb "auction-website/database/connectors/redis"
	"github.com/spf13/viper"
)

func Mysql() *db.Config {
	return &db.Config{
		User:     viper.GetString("mysql.user"),
		Password: viper.GetString("mysql.password"),
		Host:     viper.GetString("mysql.host"),
		Port:     "3306",
		DBName:   viper.GetString("mysql.dbname"),
	}
}

func MongoDB() *mdb.Config {
	return &mdb.Config{
		User:     viper.GetString("mongodb.username"),
		Password: viper.GetString("mongodb.password"),
		Host:     viper.GetString("mongodb.host"),
		Port:     "27017",
	}
}

func Redis() *rdb.Config {
	return &rdb.Config{
		Password: viper.GetString("redis.password"),
		Host:     viper.GetString("redis.host"),
		Port:     "6379",
		DBName:   viper.GetInt("redis.db"),
	}
}
