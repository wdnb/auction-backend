package api

import (
	"auction-website/api/account"
	"auction-website/api/auction"
	"auction-website/api/common"
	"auction-website/api/message"
	"auction-website/api/product"
	"auction-website/api/user"
	"auction-website/conf"
	"auction-website/docs"
	"auction-website/middleware"
	"auction-website/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func Setup(c *conf.Config) *gin.Engine {
	env := viper.GetString("http.mode")
	// 根据环境变量设置 Gin 的模式
	switch env {
	case "development":
		gin.SetMode(gin.DebugMode)
	case "testing":
		gin.SetMode(gin.TestMode)
	case "production":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	//初始化中间件
	m := middleware.Init(c)
	//全站接口访问频率控制 ws初始化成功后不增加计数器
	r.Use(m.Limit())
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	//swagger非测试环境不展示
	if env == "development" {
		//swagger http://127.0.0.1:8080/swagger/index.html#/
		docs.SwaggerInfo.BasePath = "/v1"
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	//初始路由
	version := r.Group("api/v1")

	user.Routes(version, c, m)
	product.Routes(version, c, m)
	auction.Routes(version, c, m)
	//order.Routes(r.Group("v1"), app)
	account.Routes(version, c, m)
	message.Routes(version, c, m)
	common.Routes(version, c, m)
	return r
}
