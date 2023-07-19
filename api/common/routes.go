package common

import (
	"auction-website/conf"
	"auction-website/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup, c *conf.Config, m *middleware.Config) {
	h := NewHandler(c)
	auth := m.Casbin()
	commonRoutes := r.Group("", auth.CheckPermissions())
	{
		commonRoutes.POST("/common/file/avatar", h.avatarUpload)
		commonRoutes.POST("/common/file/video", h.VideoUpload)
	}
}
