package consumer

import (
	"auction-website/api/order"
	"auction-website/conf"
	"github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func CreateOrderChannel(c *conf.Config) {
	//创建订单成功后,更新process字段
	handler := func(message *nsq.Message) error {
		service := order.NewService(c)
		//service := order.NewOrderService(com)
		id, err := service.CreateOrder(message)
		if err != nil {
			zap.L().Error("NSQ:create order failed", zap.Error(err))
			return err
		}
		if id != 0 {
			zap.L().Info("NSQ:create order success", zap.Uint32("order_id", id))
		}
		return nil
	}

	// 创建一个 NSQ 消费者
	//使用配置文件 viper 读取 channel 名称的方法是不可靠的,因为:
	//1. 配置文件修改时,程序可能已经启动,会导致读取到错误的 channel 名称
	//2. 如果配置文件中指定的 channel 名称错误,也不会报错,会默默失败
	//3. 这种方式与 NSQ Golang 客户端的设计理念不符, channel 应该在代码中明确指定
	NewNSQConsumer(viper.GetString("nsq.auction.topic_name"), "create_auction_finished", GetNsqlookupdAddress(), handler)

}
