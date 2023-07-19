package order

type Order struct {
	ID                uint32  `db:"id"`
	CustomerID        uint32  `db:"customer_id"`
	ProductID         uint32  `db:"product_id"`
	AuctionID         uint32  `db:"auction_id"`
	BidID             uint32  `db:"bid_id"`
	ShippingAddressID uint32  `db:"shipping_address_id"`
	OrderTotalAmount  float64 `db:"order_total_amount"`
	PaymentStatus     string  `db:"payment_status"`
	OrderStatus       string  `db:"order_status"`
	PaymentMethod     string  `db:"payment_method"`
	OrderType         string  `db:"order_type"`
	OrderNumber       string  `db:"order_number"`
	InitTime          int64   `db:"init_time"`
	//CreatedAt       time.Time  `db:"created_at"`
	//UpdatedAt       time.Time  `db:"updated_at"`
	//DeletedAt       *time.Time `db:"deleted_at"`
}
