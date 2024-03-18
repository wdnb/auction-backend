package auction

import (
	"auction-website/conf"
	db "auction-website/database/connectors/mysql"
	"auction-website/internal/global"
	"auction-website/message/nsq/message_utils"
	"auction-website/message/nsq/producer"
	"auction-website/utils"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Repository struct {
	producer *nsq.Producer
	db       *sqlx.DB
}

func NewRepository(c *conf.Config) *Repository {
	return &Repository{
		producer: producer.Init(),
		db:       db.GetClient(c.Mysql),
	}
}

// Create a new auction in the database
func (r *Repository) CreateAuction(a *Auction) (uint32, error) {
	// Use sqlx.NamedExec to insert a new auction into the database using data binding
	query := `INSERT INTO auction (name, product_id, product_uid, category_id, start_price,  fixed_price, bid_increment,processed, status, image_url, description, start_time, end_time) 
           VALUES (:name, :product_id, :product_uid, :category_id, :start_price, :fixed_price, :bid_increment,:processed, :status, :image_url, :description, :start_time, :end_time)`
	result, err := r.db.NamedExec(query, &a)
	if err != nil {
		return 0, err
	}
	// Retrieve the ID of the new auction and return it
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

// TODO 因当做筛选 where status=1...
func (r *Repository) GetAllAuctions(page uint32, pageSize uint32) ([]*List, error) {
	offset := utils.Offset(page, pageSize)
	query := `SELECT 
    id,name,product_id,product_uid,category_id,start_price,fixed_price,bid_increment,status,image_url,description,start_time,end_time
    FROM auction ORDER BY id DESC LIMIT :limit OFFSET :offset`

	// Bind the limit and offset values to the query using sqlx.Named
	query, args, err := sqlx.Named(query, map[string]interface{}{
		"limit":  pageSize,
		"offset": offset,
	})
	if err != nil {
		return nil, err
	}
	// Use sqlx.Select to execute the query and retrieve the results
	var auctions []*List
	//fmt.Println(auctions)
	err = r.db.Select(&auctions, query, args...)
	if err != nil {
		return nil, err
	}
	//fmt.Println(auctions)
	return auctions, nil
}

// Retrieve a single auction by ID from the database
func (r *Repository) GetAuction(id uint32) (*Auction, error) {
	// Define the query to retrieve the auction by ID
	query := `SELECT 
    id,name,product_id,product_uid,category_id,start_price,fixed_price,bid_increment,status,image_url,description,start_time,end_time
    FROM auction WHERE id = ?`

	// Use sqlx.Get to execute the query and retrieve the result
	var a Auction
	err := r.db.Get(&a, query, id)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Delete an auction by ID from the database
func (r *Repository) DeleteAuction(id uint32) error {
	// Define the query to delete the auction by ID
	query := `DELETE FROM auction WHERE id = ?`

	// Use sqlx.Exec to execute the query and delete the auction
	result, err := r.db.Exec(query, id)
	num, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		return global.ErrNotDelete
	}
	return nil

}

// Update an existing auction in the database
func (r *Repository) UpdateAuction(a *Update) error {
	// Define the query to update the auction by ID
	query := `UPDATE auction SET status=:status WHERE id=:id`

	// Use sqlx.NamedExec to execute the query and update the auction
	result, err := r.db.NamedExec(query, &a)
	num, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		return global.ErrNotUpdate
	}
	return nil
}

// Update an existing auction in the database
func (r *Repository) UpdateProcessed(a *Processed) error {
	//fmt.Println(a)
	// Define the query to update the auction by ID
	query := `UPDATE auction SET processed=:processed WHERE id=:id`

	// Use sqlx.NamedExec to execute the query and update the auction
	result, err := r.db.NamedExec(query, &a)
	num, err := result.RowsAffected()
	//fmt.Println(err)
	if err != nil {
		return err
	}
	if num == 0 {
		return errors.New("更新目标不存在或已更新！")
	}
	return nil
}

// GetActiveAuctions returns a list of active auctions within the given time scope.
func (s *Repository) GetUnprocessedAuctions(startTimeScope, endTimeScope int64) ([]PayloadAuction, error) {
	sql := `SELECT
		id, name, product_id, product_uid, category_id, start_price, fixed_price, bid_increment, status, image_url, description, start_time, end_time 
		FROM auction WHERE
		end_time >=? AND end_time <= ? AND processed=0 AND status='SOLD'`
	rows, err := s.db.Queryx(sql, startTimeScope, endTimeScope)
	if err != nil {
		return nil, fmt.Errorf("failed to query auctions: %v", err)
	}
	defer rows.Close()

	auctions := make([]PayloadAuction, 0)
	for rows.Next() {
		var a PayloadAuction
		if err := rows.StructScan(&a); err != nil {
			zap.L().Error("failed to scan auction", zap.Error(err))
			continue
		}
		auctions = append(auctions, a)
	}
	return auctions, nil
}

func checkAuctions_old(producer *nsq.Producer, db *sqlx.DB) error {
	// get the current time
	now := time.Now().Unix()
	// get the auction time start and end offset
	auctionTimeStartOffset := viper.GetInt64("app.auction_time_start_offset")
	auctionTimeEndOffset := viper.GetInt64("app.auction_time_end_offset")
	//nsq消息过期前的最大限制是15分钟 -max-msg-timeout
	//https://nsq.io/components/nsqd.html
	if auctionTimeEndOffset > 900 {
		return errors.New("app.auction_time_end_offset is too large, max is 900")
	}
	startTimeScope := now - auctionTimeStartOffset
	endTimeScope := now + auctionTimeEndOffset

	sql := `SELECT
		id, name, product_id, product_uid, category_id, start_price, fixed_price, bid_increment, status, image_url, description, start_time, end_time FROM auction WHERE
		end_time >=? AND end_time <= ? AND processed=0`
	rows, err := db.Queryx(sql, startTimeScope, endTimeScope)
	if err != nil {
		return errors.New("Failed to query auctions:" + err.Error())
	}
	//fmt.Println(rows.)
	defer rows.Close()
	// loop through the rows and send a message to NSQ for each auction
	for rows.Next() {
		var a PayloadAuction
		//!!和sqlselect字段查出来的不匹配会报错
		//err := rows.Scan(&a.ID, &a.Name, &a.ProductID, &a.ProductUID, &a.CategoryID, &a.StartPrice, &a.CurrentPrice,
		//&a.FixedPrice, &a.BidIncrement, &a.Status, &a.ImageURL, &a.Description, &a.StartTime, &a.EndTime)
		err := rows.StructScan(&a)
		//fmt.Println(a)
		if err != nil {
			zap.L().Error("NSQ:Failed to scan auction", zap.Error(err))
			continue
		}

		////获得最高出价者
		//service := auction2.NewAuctionService(s.Com)
		//bid, err := service.GetHighestBid(a.ID)
		//if err != nil {
		//	return 0, err
		//}
		//
		delay := a.EndTime - now

		topic := viper.GetString("nsq.auction.topic_name")
		jsonMessage, err := message_utils.CreateMessage(topic, "create_auction_finished", a)
		if err != nil {
			zap.L().Error("NSQ:Failed to create message", zap.Error(err))
		}
		//cc := 5 * time.Minute
		//message := fmt.Sprintf("%d|%d", auctionId, endTime)
		delayDurTime := time.Duration(delay) * time.Second
		//系统停止运行期间结束的拍卖 进入这里创建订单
		if delay < 0 {
			zap.L().Info("NSQ:end_time is less than current time")
			err = producer.Publish(topic, jsonMessage)
			continue
		} else {
			err = producer.DeferredPublish(topic, delayDurTime, jsonMessage)
			if err != nil {
				zap.L().Error("NSQ:Failed to publish message to NSQ", zap.Error(err))
			} else {
				zap.L().Info("NSQ:Deferred published to NSQ", zap.Duration("delayDurTime", delayDurTime), zap.Any("struct", &a))
			}
		}
	}
	return nil
}

func QueryWithLogging(db *sqlx.DB, query string, args ...interface{}) (*sqlx.Rows, error) {
	start := time.Now()
	rows, err := db.Queryx(query, args...)
	elapsed := time.Since(start)

	if err != nil {
		// 记录 SQL 查询错误、执行时间、查询语句和参数
		//log.Printf("SQL query error: %s\nElapsed time: %s\nSQL query: %s\nSQL parameters: %v", err.Error(), elapsed, query, args)
		zap.L().Error(fmt.Sprintf("SQL query error: %s\nElapsed time: %s\nSQL query: %s\nSQL parameters: %v", err.Error(), elapsed, query, args))
		//zap.L().Error("SQL query error", zap.Error(err))
	} else {
		// 记录 SQL 查询语句和参数
		log.Printf("SQL query: %s\nSQL parameters: %v", query, args)
	}

	return rows, err
}

// Create a new bid in the database for a specific auction
func (r *Repository) CreateBid(b *Bid) (uint32, error) {
	// Use sqlx.NamedExec to insert a new bid into the database using data binding
	query := `INSERT INTO bid (auction_id, customer_id, bid_price) VALUES (:auction_id, :customer_id, :bid_price)`
	result, err := r.db.NamedExec(query, &b)
	if err != nil {
		return 0, err
	}
	// Retrieve the ID of the new bid and return it
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

// Retrieve all bid for a specific auction from the database
func (r *Repository) GetAllBids(id, page, pageSize uint32) ([]*BidsList, error) {
	offset := utils.Offset(page, pageSize)
	// Define the query to retrieve all bid for a specific auction
	query := `SELECT id, auction_id, customer_id, bid_price FROM bid WHERE auction_id=:auction_id ORDER BY bid_price DESC LIMIT :limit OFFSET :offset`
	query, args, err := sqlx.Named(query, map[string]interface{}{
		"auction_id": id,
		"limit":      pageSize,
		"offset":     offset,
	})
	if err != nil {
		return nil, err
	}
	// Use sqlx.Select to execute the query and retrieve the results
	var bids []*BidsList
	err = r.db.Select(&bids, query, args...)
	if err != nil {
		return nil, err
	}
	return bids, nil
}

// Retrieve the highest bid for a specific auction from the database
func (r *Repository) GetHighestBid(id uint32) (*Bid, error) {
	//fmt.Println(id)
	// Define the query to retrieve the highest bid for a specific auction
	query := `SELECT id, customer_id, auction_id, bid_price FROM bid WHERE auction_id=? ORDER BY bid_price DESC LIMIT 1`
	// Use sqlx.Get to execute the query and retrieve the result
	var bid Bid
	err := r.db.Get(&bid, query, id)
	if err != nil {
		return nil, err
	}
	return &bid, nil
}
