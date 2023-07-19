package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

var db *sqlx.DB

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func GetClient(c *Config) *sqlx.DB {
	if db == nil {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
			c.User, c.Password, c.Host, c.Port, c.DBName)

		db = initDB(dsn)
	}
	return db
}

func initDB(dsn string) *sqlx.DB {
	d := conn(dsn)
	d.SetMaxOpenConns(200)
	d.SetMaxIdleConns(100)
	return d
}

func conn(dsn string) *sqlx.DB {
	var maxRetries = 3                  // 最大重试次数
	var retryInterval = time.Second * 5 // 重试间隔
	for i := 0; i <= maxRetries; i++ {
		db, err := sqlx.Connect("mysql", dsn)
		if err == nil {
			return db
		}
		zap.L().Error("Failed to connect to Mysql", zap.Int("attempt", i+1), zap.Int("max_retries", maxRetries), zap.Error(err))
		if i == maxRetries {
			panic(fmt.Errorf("failed to connect to mysql after %d retries", maxRetries))
		}
		time.Sleep(retryInterval)
	}
	return db
}
