package order

import (
	"auction-website/conf"
	db "auction-website/database/connectors/mysql"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(c *conf.Config) *Repository {
	return &Repository{
		db: db.GetClient(c.Mysql),
	}
}

func (r *Repository) CreateOrder(o *Order) (uint32, error) {

	query := `INSERT INTO orders (customer_id, product_id, auction_id, bid_id, shipping_address_id, order_total_amount, payment_status,order_status,payment_method, order_type, order_number, init_time) 
            VALUES (:customer_id, :product_id, :auction_id, :bid_id, :shipping_address_id, :order_total_amount, :order_status, :payment_status, :payment_method, :order_type, :order_number, :init_time)`

	result, err := r.db.NamedExec(query, o)
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
