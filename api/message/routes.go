package message

import (
	"auction-website/conf"
	"auction-website/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup, c *conf.Config, m *middleware.Config) {
	h := NewHandler(c)
	auth := m.Casbin()
	//im模块和消息通知模块
	imRoutes := r.Group("", auth.CheckPermissions())
	{
		imRoutes.GET("/im/ws", h.WsHandler)
	}
}
