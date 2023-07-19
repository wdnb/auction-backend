package order

import (
	"auction-website/api/auction"
	"auction-website/api/user"
	"auction-website/conf"
	"auction-website/message/nsq/message_utils"
	"auction-website/utils"
	"fmt"
	"math/rand"
	"time"

	"github.com/nsqio/go-nsq"
)

type Service struct {
	userService    *user.Service
	auctionService *auction.Service
	repo           *Repository
}

func NewService(c *conf.Config) *Service {
	return &Service{
		userService:    user.NewService(c),
		auctionService: auction.NewService(c),
		repo:           NewRepository(c),
	}
}

// 检测到拍卖结束后 订阅发布模式 自动生成订单
func (s *Service) CreateOrder(message *nsq.Message) (uint32, error) {
	parseMessage, err := message_utils.ParseMessage(message.Body)
	if err != nil {
		return 0, err
	}
	var a auction.PayloadAuction
	err = utils.FromJSONBytes(parseMessage.Body, &a)
	if err != nil {
		return 0, err
	}

	//获得用户地址
	address, err := s.userService.GetActiveUserAddress(a.CustomerID)
	var aid uint32
	if err != nil {
		//查不到地址表id就置为0
		aid = 0
	} else {
		aid = address.ID
	}

	o := Order{
		CustomerID:        a.CustomerID,
		ProductID:         a.ProductID,
		AuctionID:         a.ID,
		BidID:             a.BidID,
		ShippingAddressID: aid,
		OrderTotalAmount:  a.BidPrice, //订单总金额
		PaymentStatus:     "PENDING",  //支付状态
		OrderStatus:       "PENDING",  //订单状态
		PaymentMethod:     "",         //支付方式
		//PaymentMethod:     sql.NullString{String: "Credit Card", Valid: true},
		OrderType:   a.OrderType,
		OrderNumber: GenerateOrderNumber(),
		InitTime:    time.Now().Unix(),
	}
	//
	id, err := s.repo.CreateOrder(&o)
	if err != nil {
		return id, err
	}
	//as := auction2.NewAuctionService(s.Com)

	var ap = auction.Processed{
		ID:        a.ID,
		Processed: 2, //推送成功，修改processed字段状态为1 消费完毕后改成2
	}
	err = s.auctionService.UpdateProcessed(&ap)
	if err != nil {
		return 0, err
	}

	//fmt.Println(id)
	return id, nil
}

// GenerateOrderNumber generates a unique order number using the current timestamp and a random number.
// It returns a string in the format "YYYYMMDDHHmmssRRRR", where YYYY is the year, MM is the month, DD is the day,
// HH is the hour, mm is the minute, ss is the second, and RRRR is a random 4-digit number.
func GenerateOrderNumber() string {
	t := time.Now()
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d%04d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), rand.Intn(10000))
}

// Place the above code inside the Service struct, just below the NewOrderService function.
