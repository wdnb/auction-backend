package auction

import (
	"auction-website/conf"
	"auction-website/message/nsq/message_utils"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Service struct {
	c    *conf.Config
	repo *Repository
}

func NewService(c *conf.Config) *Service {
	return &Service{
		c:    c,
		repo: NewRepository(c),
	}
}

func (s *Service) CreateAuction(a *Auction) (uint32, error) {
	// Call the CreateAuction method of the repository to insert the new auction into the database
	auctionID, err := s.repo.CreateAuction(a)
	if err != nil {
		return 0, err
	}
	// Return the ID of the newly created auction
	return auctionID, nil
}

// CheckAuctions 每x分钟检测一次拍卖表状态 如果有将要拍卖结束的 推送延时队列 创建订单（nsq default 最久15m0s）
func (s *Service) CheckAuctions() error {
	//producer := s.Repository.producer
	nsq := s.repo.producer
	//db := s.Repository.db
	//get the current time
	now := time.Now().Unix()
	// get the auction time start and end offset
	auctionTimeStartOffset := viper.GetInt64("app.auction_time_start_offset")
	auctionTimeEndOffset := viper.GetInt64("app.auction_time_end_offset")
	//https://nsq.io/components/nsqd.html
	//range of 0-3600 s
	if auctionTimeEndOffset > 3600 {
		return errors.New("app.auction_time_end_offset is too large, max is 3600 s")
	}

	startTimeScope := now - auctionTimeStartOffset
	endTimeScope := now + auctionTimeEndOffset

	auctionRepo := NewRepository(s.c)
	topic := viper.GetString("nsq.auction.topic_name")

	auctions, err := auctionRepo.GetUnprocessedAuctions(startTimeScope, endTimeScope)
	if err != nil {
		return errors.New("Failed to query auctions:" + err.Error())
	}
	//fmt.Println(auctions)
	// loop through the auctions and send a message to NSQ for each one
	for _, a := range auctions {
		//订单类型为拍卖
		a.OrderType = "AUCTION"
		// get the highest bid for the auction
		bid, err := auctionRepo.GetHighestBid(a.ID)
		if err != nil {
			//所谓的流拍 不管是不是系统异常 先直接processed 127 避免死循环
			//todo 需要提供一个手动创建order的方法
			_ = s.UpdateProcessed(&Processed{ID: a.ID, Processed: 127})
			zap.L().Info("NSQ:Failed to get highest bid for auction", zap.Uint32("auction_id", a.ID), zap.Error(err))
			continue
		}
		//获得最高出价者
		a.BidID = bid.ID
		a.CustomerID = bid.CustomerID
		a.BidPrice = bid.BidPrice
		//a.AuctionID = bid.AuctionID

		delay := a.EndTime - now

		jsonMessage, err := CreateAuctionMessage(topic, a)
		if err != nil {
			return err
		}

		delayDurTime := time.Duration(delay) * time.Second

		// if the auction has already ended, create the order immediately
		if delay < 0 {
			err = nsq.Publish(topic, jsonMessage)
			if err != nil {
				zap.L().Error("NSQ:Failed to publish message to NSQ", zap.Error(err))
				continue
			}
			zap.L().Info("NSQ:Published to NSQ", zap.Any("auction", &a))
		} else {
			err = nsq.DeferredPublish(topic, delayDurTime, jsonMessage)
			if err != nil {
				zap.L().Error("NSQ:Failed to publish message to NSQ", zap.Error(err))
				continue
			}
			zap.L().Info("NSQ:Deferred published to NSQ", zap.Duration("delayDurTime", delayDurTime), zap.Any("auction", &a))
		}
		// mark the auction as processed
		//推送成功，修改processed字段状态为1 消费完毕后改成2
		err = s.UpdateProcessed(&Processed{ID: a.ID, Processed: 1})
		if err != nil {
			zap.L().Error("Failed to mark auction as processed", zap.Uint32("auction_id", a.ID), zap.Error(err))
			continue
		}

	}
	return nil
}

// CreateAuctionMessage creates a message to be published to NSQ for the given auction.
func CreateAuctionMessage(topic string, auction PayloadAuction) ([]byte, error) {
	//topic := viper.GetString("nsq.auction.topic_name")
	jsonMessage, err := message_utils.CreateMessage(topic, "create_auction_finished", auction)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %v", err)
	}
	return jsonMessage, nil
}

func (s *Service) getAllAuctions(page, pageSize uint32) ([]*List, error) {
	// Call the GetAllAuctions method of the repository to retrieve all auctions from the database
	auctions, err := s.repo.GetAllAuctions(page, pageSize)
	if err != nil {
		return nil, err
	}
	// Return the retrieved auctions
	return auctions, nil
}

// Define a method to retrieve a single auction by ID
func (s *Service) GetAuction(id uint32) (*Auction, error) {
	// Call the GetAuction method of the repository to retrieve the auction from the database
	a, err := s.repo.GetAuction(id)
	if err != nil {
		return nil, err
	}
	// Return the retrieved auction
	return a, nil
}

// Define a method to delete an auction by ID
func (s *Service) DeleteAuction(id uint32) error {
	// Call the DeleteAuction method of the repository to delete the auction from the database
	err := s.repo.DeleteAuction(id)
	if err != nil {
		return err
	}
	// Return nil if the auction was successfully deleted
	return nil
}

// Define a method to update an existing auction by ID
func (s *Service) UpdateAuction(a *Update) error {
	// Call the UpdateAuction method of the repository to update the auction in the database
	err := s.repo.UpdateAuction(a)
	if err != nil {
		return err
	}
	// Return nil if the auction was successfully updated
	return nil
}

func (s *Service) UpdateProcessed(a *Processed) error {
	// Call the UpdateAuction method of the repository to update the auction in the database
	err := s.repo.UpdateProcessed(a)
	if err != nil {
		return err
	}
	// Return nil if the auction was successfully updated
	return nil
}

// Define a method to create a new bid for an auction
func (s *Service) CreateBid(a *Bid) (uint32, error) {
	// Call the CreateBid method of the repository to insert the new bid into the database
	bidID, err := s.repo.CreateBid(a)
	if err != nil {
		return 0, err
	}
	// Return the ID of the newly created bid
	return bidID, nil
}

// Define a method to retrieve all bids for an auction by auction ID
func (s *Service) getAllBids(id, page, pageSize uint32) ([]*BidsList, error) {
	// Call the GetAllBids method of the repository to retrieve all bids for the auction from the database
	bids, err := s.repo.GetAllBids(id, page, pageSize)
	if err != nil {
		return nil, err
	}
	// Return the retrieved bids
	return bids, nil
}

// Define a method to retrieve the highest bid for an auction by auction ID
func (s *Service) GetHighestBid(auctionID uint32) (*Bid, error) {
	// Call the GetHighestBid method of the repository to retrieve the highest bid for the auction from the database
	highestBid, err := s.repo.GetHighestBid(auctionID)
	if err != nil {
		return nil, err
	}
	// Return the retrieved highest bid
	return highestBid, nil
}
