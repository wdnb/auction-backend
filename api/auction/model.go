package auction

// Define a struct Auction that matches the fields of the auctions table in the given SQL statement
type Auction struct {
	ID           uint32  `json:"id" db:"id"`
	Name         string  `json:"name" db:"name" validate:"required" example:"上拍品名"`
	ProductID    uint32  `json:"product_id" db:"product_id" validate:"number" example:"101"`
	ProductUID   uint32  `json:"product_uid" db:"product_uid" validate:"number" example:"101"`
	CategoryID   uint32  `json:"category_id" db:"category_id" validate:"number" example:"101"`
	StartPrice   float64 `json:"start_price" db:"start_price" validate:"number" example:"2.1"`
	FixedPrice   float64 `json:"fixed_price" db:"fixed_price" validate:"number" example:"100"`
	BidIncrement float64 `json:"bid_increment" db:"bid_increment" validate:"number" example:"10"`
	Processed    uint8   `json:"processed" db:"processed" validate:"number" example:"0"`
	Status       string  `db:"status" json:"status" validate:"required,oneof=upcoming active sold unsold expired cancelled" example:"upcoming"`
	ImageURL     string  `json:"image_url" db:"image_url" example:"https://example.com/image.jpg"`
	Description  string  `json:"description" db:"description" example:"这里是拍品描述"`
	StartTime    int64   `json:"start_time" db:"start_time" example:"1681398030"`
	EndTime      int64   `json:"end_time" db:"end_time" example:"1681819182"`
}

// PayloadAuction 用于消息队列生产者和消费者之间传输数据的结构体
type PayloadAuction struct {
	ID           uint32  `json:"id" db:"id"`
	Name         string  `json:"name" db:"name" validate:"required"`
	ProductID    uint32  `json:"product_id" db:"product_id" validate:"number"`
	ProductUID   uint32  `json:"product_uid" db:"product_uid" validate:"number"`
	CategoryID   uint32  `json:"category_id" db:"category_id" validate:"number"`
	StartPrice   float64 `json:"start_price" db:"start_price" validate:"number"`
	FixedPrice   float64 `json:"fixed_price" db:"fixed_price" validate:"number"`
	BidIncrement float64 `json:"bid_increment" db:"bid_increment" validate:"number"`
	Processed    uint8   `json:"processed" db:"processed" validate:"number"`
	Status       string  `db:"status" json:"status" validate:"required,oneof=upcoming active sold unsold expired cancelled"`
	ImageURL     string  `json:"image_url" db:"image_url"`
	Description  string  `json:"description" db:"description"`
	StartTime    int64   `json:"start_time" db:"start_time"`
	EndTime      int64   `json:"end_time" db:"end_time"`
	BidID        uint32  `json:"bid_id" db:"bid_id"`
	CustomerID   uint32  `json:"customer_id" db:"customer_id"`
	BidPrice     float64 `json:"bid_price" db:"bid_price"`
	OrderType    string  `json:"order_type" db:"order_type"`
}

// Define a struct AuctionList that will be used to generate a list of auctions
type List struct {
	Auction
}

type Update struct {
	ID     uint32 `db:"id" json:"-"`
	Status string `db:"status" json:"status" validate:"required,oneof=upcoming active sold unsold expired cancelled"`
}

type Processed struct {
	ID        uint32 `db:"id" json:"id"`
	Processed uint8  `db:"processed" json:"processed" validate:"required"`
}

type Bid struct {
	ID         uint32  `json:"-" db:"id"`
	CustomerID uint32  `json:"customer_id" db:"customer_id" validate:"number" example:"101"`
	AuctionID  uint32  `json:"-" db:"auction_id" validate:"number" example:"101"`
	BidPrice   float64 `json:"bid_price" db:"bid_price" validate:"required,number" example:"101.1"`
}

type BidsList struct {
	Bid
}

type BidsUpdate struct {
	ID     uint32  `json:"id" validate:"required,number"`
	Price  float64 `json:"bid_price" validate:"required,number"`
	Status string  `json:"status" validate:"required"`
}
