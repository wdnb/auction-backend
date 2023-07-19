package producer

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var producer *nsq.Producer

func Init() *nsq.Producer {
	if producer == nil {
		producer = initProducer()
	}
	return producer
}

func initProducer() *nsq.Producer {
	// 创建一个 nsqlookupd 配置
	lookupCfg := nsq.NewConfig()

	// 创建一个 nsqlookupd 生产者
	lookupProducer, err := nsq.NewProducer(GetNsqdAddress(), lookupCfg)
	if err != nil {
		zap.L().Error("NSQ:created nsq_producer failed, err", zap.Error(err))
		panic(err)
	}
	return lookupProducer
}

func GetNsqdAddress() string {
	return fmt.Sprintf(
		"%s:%d",
		viper.GetString("nsqd.host"),
		viper.GetInt("nsqd.tcp_port"))
}
