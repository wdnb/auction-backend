package message

import (
	"auction-website/conf"
)

type Service struct {
	repo *Repository
}

func NewService(c *conf.Config) *Service {
	return &Service{
		repo: NewRepository(c),
	}
}
