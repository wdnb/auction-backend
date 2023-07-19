package product

import (
	"auction-website/conf"
	"auction-website/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup, c *conf.Config, m *middleware.Config) {
	auth := m.Casbin()
	h := NewHandler(c)
	productRoutes := r.Group("")
	{
		productRoutes.GET("/product/list", h.listHandler)
		productRoutes.GET("/product/:id", h.detailHandler)
		productRoutes.POST("/product", auth.CheckPermissions(), h.createHandler)
		productRoutes.PUT("/product/:id", auth.CheckPermissions(), h.updateHandler)
		productRoutes.DELETE("/product/:id", auth.CheckPermissions(), h.deleteHandler)
	}
}
