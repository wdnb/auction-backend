package common

import (
	"auction-website/conf"
)

type Repository struct {
}

func NewRepository(c *conf.Config) *Repository {
	return &Repository{}
}
