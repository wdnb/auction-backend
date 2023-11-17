package auction

import (
	"auction-website/conf"
	"auction-website/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup, c *conf.Config, m *middleware.Config) {
	h := NewHandler(c)
	auth := m.Casbin()
	// 创建拍卖品
	r.POST("/auction", auth.CheckPermissions(), h.createAuction)
	// 获取所有拍卖品
	r.GET("/auction", h.getAuctions)
	// 获取单个拍卖品
	r.GET("/auction/:id", h.getAuction)
	// 修改拍卖品
	r.PUT("/auction/:id", auth.CheckPermissions(), h.updateAuction)
	// 删除拍卖品
	r.DELETE("/auction/:id", auth.CheckPermissions(), h.deleteAuction)

	// 创建出价
	r.POST("/auction/:id/bid", auth.CheckPermissions(), h.createBid)
	// 获取所有出价
	r.GET("/auction/:id/bid", h.getBids)

	//拍卖结束 通知
	//r.POST("/auctions/:id/notify", utils.JWTMiddleware(), h.sendAuctionSuccessNotificationHandler)

	//r.GET("/auctions/consumer", h.createConsumer)
}
