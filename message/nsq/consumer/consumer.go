package consumer

import (
	"auction-website/conf"
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 注册消费者
func InitConsumer(c *conf.Config) {
	go CreateOrderChannel(c)
}

// 生成消费者
func NewNSQConsumer(topic, channel, nsqdAddress string, handler nsq.HandlerFunc) {
	// 创建一个 nsqlookupd 配置
	lookupCfg := nsq.NewConfig()

	// 创建一个 nsqlookupd 消费者
	lookupConsumer, err := nsq.NewConsumer(topic, channel, lookupCfg)
	if err != nil {
		zap.L().Error("NSQ:创建消费者失败", zap.Error(err))
		panic(err)
	}
	// 设置消息处理函数
	lookupConsumer.AddHandler(handler)
	// 连接到 nsqlookupd
	err = lookupConsumer.ConnectToNSQLookupd(nsqdAddress)
	if err != nil {
		zap.L().Error("NSQ:连接nsqlookupd失败", zap.Error(err))
		panic(err)
	}
	// 阻塞等待消息
	select {}
}

func GetNsqlookupdAddress() string {
	return fmt.Sprintf(
		"%s:%d",
		viper.GetString("nsqlookupd.host"),
		viper.GetInt("nsqlookupd.http_port"))
}
