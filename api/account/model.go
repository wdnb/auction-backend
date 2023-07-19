package account

import "time"

type Account struct {
	ID      uint32  `db:"id" json:"id"`
	UserID  uint32  `db:"user_id" json:"user_id"`
	Balance float64 `json:"balance" json:"balance" validate:"required"`
	Status  string  `db:"status" json:"status"`
	//CreatedAt time.Time       `db:"created_at" json:"created_at"`
	//UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

type Amount struct {
	Amount float64 `json:"amount" validate:"required"`
}

type WithdrawalRecord struct {
	ID        uint32    `db:"id" json:"-"`
	UserID    uint32    `db:"user_id" json:"user_id"`
	OrderNo   string    `db:"order_no" json:"order_no"`
	Amount    float64   `db:"amount" json:"amount"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
