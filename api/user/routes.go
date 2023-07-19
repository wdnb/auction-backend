package user

import (
	"auction-website/conf"
	"auction-website/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup, c *conf.Config, m *middleware.Config) {
	auth := m.Casbin()
	h := NewHandler(c)
	//userRoutes.POST("/register", h.registerHandler)
	r.POST("/user/login", h.loginHandler)
	r.POST("/user/login-phone", h.loginPhoneHandler)
	r.POST("/user/verification-code/:kind", h.getVerificationCode)
	perRoutes := r.Group("", auth.CheckPermissions())
	{
		perRoutes.GET("/user/profile/:uid", h.getUserProfile)
		perRoutes.PUT("/user/profile/:uid", h.updateUserProfile)
		perRoutes.POST("/user/shipping-address", h.createShippingAddress)
		perRoutes.PUT("/user/shipping-address/:id", h.updateUserAddressByID)
		perRoutes.GET("/user/shipping-address", h.getUserAddresses)
		perRoutes.DELETE("/user/shipping-address/:id", h.deleteUserAddressByID)
		perRoutes.GET("/user/center", h.getUserCenter)
		//perRoutes.PUT("", h.updateUserCenter)
	}
}
