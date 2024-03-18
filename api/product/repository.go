package product

import (
	"auction-website/conf"
	db "auction-website/database/connectors/mysql"

	"auction-website/internal/global"
	"auction-website/utils"

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

func (r *Repository) CreateProduct(product *Product) (uint32, error) {
	query := `INSERT INTO product (user_id, name, description, category_id, start_price, quantity, image_url,bid_increment,fixed_price,start_time,end_time) 
VALUES (:user_id, :name, :description, :category_id, :start_price, :quantity, :image_url,:bid_increment,:fixed_price,:start_time,:end_time)`
	//_, err := r.db.Exec(query, product.Name, product.Description, product.BidPrice, product.ImageURL)
	result, err := r.db.NamedExec(query, product)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint32(id), err
}

func (r *Repository) GetProductList(page, pageSize uint32) ([]*Product, error) {
	offset := utils.Offset(page, pageSize)
	query := `SELECT id, user_id, name, description, category_id, start_price, quantity, image_url FROM product LIMIT ? OFFSET ?`
	var products []*Product
	err := r.db.Select(&products, query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	//for _, v := range products {
	//	fmt.Println(*v)
	//
	//}
	return products, nil
}

func (r *Repository) UpdateProduct(p *UpdateProduct) error {
	query := `UPDATE product SET
                   name=:name, description=:description, category_id=:category_id, start_price=:start_price, quantity=:quantity,image_url=:image_url, bid_increment=:bid_increment, fixed_price=:fixed_price, start_time=:start_time, end_time=:end_time
                   WHERE id=:id`
	result, err := r.db.NamedExec(query, p)
	num, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		return global.ErrNotUpdate
		//TODO 其他地方的更新目标不存在或已更新也修改
	}
	return nil
}

func (r *Repository) GetProductByID(id uint32) (*Product, error) {
	query := `SELECT id, user_id, name, description, category_id, start_price, quantity, image_url FROM product WHERE id=?`
	var p Product
	err := r.db.Get(&p, query, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) DeleteProduct(id uint32) error {
	query := `DELETE FROM product WHERE id=?`
	_, err := r.db.Exec(query, id)
	return err
}
