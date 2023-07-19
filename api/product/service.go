package product

import (
	"auction-website/conf"
	"auction-website/internal/global"
	"database/sql"
)

type Service struct {
	repo *Repository
}

func NewService(c *conf.Config) *Service {
	return &Service{
		repo: NewRepository(c),
	}
}

func (s *Service) CreateProduct(p *Product) (uint32, error) {
	pid, err := s.repo.CreateProduct(p)
	if err != nil {
		return 0, err
	}
	return pid, nil
}

// todo 增加where条件的参数
func (s *Service) GetProductList(page, pageSize uint32) ([]*Product, error) {
	products, err := s.repo.GetProductList(page, pageSize)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *Service) UpdateProduct(p *UpdateProduct) error {
	err := s.repo.UpdateProduct(p)
	if err != nil {
		return err
	}
	return nil
}

// GetProductDetail retrieves a single product by its ID
func (s *Service) GetProductDetail(productID uint32) (*Product, error) {
	prod, err := s.repo.GetProductByID(productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, global.ErrNotFound
		}
		return nil, err
	}
	return prod, nil
}

func (s *Service) DeleteProduct(productID uint32) error {
	err := s.repo.DeleteProduct(productID)
	if err != nil {
		return err
	}
	return nil
}
