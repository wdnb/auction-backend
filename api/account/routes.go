package account

import (
	"auction-website/conf"
	"auction-website/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup, c *conf.Config, m *middleware.Config) {
	h := NewHandler(c)
	auth := m.Casbin()
	account := r.Group("", auth.CheckPermissions())
	{
		//商家入账和用户出账的记录，可以在其他表中记录，比如订单表。
		//在订单表中可以记录购买该订单的用户ID，订单金额，商家ID，订单状态等信息，订单状态可以包含待支付，已支付等状态，
		//商家可以根据订单状态和订单金额来计算自己的收益，用户也可以根据订单状态和订单金额来查看自己的消费情况。
		account.POST("/account/deposit", h.Deposit)
		account.POST("/account/withdraw", h.Withdraw)
		account.GET("/account/balance", h.GetBalance)
		account.GET("/account/withdrawal-record", h.WithdrawalRecord)
	}

}
