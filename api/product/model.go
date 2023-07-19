package product

// Product represents a product that can be sold on the auction website
type Product struct {
	ID           uint32  `db:"id" json:"id"`
	UserID       uint32  `db:"user_id" json:"user_id" validate:"required,numeric" description:"用户id"`
	Name         string  `db:"name" json:"name" validate:"min=6,max=30"`
	Description  string  `db:"description" json:"description" validate:"required"`
	ImageURL     string  `db:"image_url" json:"image_url" validate:"required"`
	CategoryID   uint32  `db:"category_id" json:"category_id" validate:"required,numeric"`
	StartPrice   float64 `db:"start_price" json:"start_price" validate:"required,numeric"`
	Quantity     uint32  `db:"quantity" json:"quantity" validate:"min=1,max=88888888"`
	BidIncrement float64 `db:"bid_increment" json:"bid_increment" validate:"required,numeric"`
	FixedPrice   float64 `db:"fixed_price" json:"fixed_price" validate:"required,numeric"`
	StartTime    int64   `db:"start_time" json:"start_time" validate:"required,numeric"`
	EndTime      int64   `db:"end_time" json:"end_time" validate:"required,numeric"`
}

type UpdateProduct struct {
	ID           uint32  `db:"id" json:"-"`
	UserID       uint32  `db:"user_id" json:"-"`
	Name         string  `db:"name" json:"name" validate:"min=6,max=30" example:"商品名"`
	Description  string  `db:"description" json:"description" validate:"json" example:"这里是商品描述"`
	ImageURL     string  `db:"image_url" json:"image_url" example:"https://example.com/image.jpg"`
	CategoryID   uint32  `db:"category_id" json:"category_id" validate:"numeric" example:"2"`
	StartPrice   float64 `db:"start_price" json:"start_price" validate:"numeric" example:"2.1"`
	Quantity     uint32  `db:"quantity" json:"quantity" validate:"min=1,max=88888888" example:"998"`
	BidIncrement float64 `db:"bid_increment" json:"bid_increment" validate:"numeric" example:"10"`
	FixedPrice   float64 `db:"fixed_price" json:"fixed_price" validate:"numeric" example:"100"`
	StartTime    int64   `db:"start_time" json:"start_time" validate:"numeric" example:"1681398030"`
	EndTime      int64   `db:"end_time" json:"end_time" validate:"numeric" example:"1681819182"`
}
