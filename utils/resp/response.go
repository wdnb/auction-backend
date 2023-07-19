package resp

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Response struct {
	Code    int   `json:"code" default:"0"`
	TimeNow int64 `json:"time_now"`
	Data    any   `json:"data"`
	//Msg     string `json:"msg"`
}

type ErrResponse struct {
	Code    int    `json:"code" default:"1000"`
	TimeNow int64  `json:"time_now"`
	Msg     string `json:"msg"`
}

const (
	ERROR   = 1000
	SUCCESS = 0
)

func DataResponse(c *gin.Context, data interface{}) {
	code := SUCCESS
	//msg := ""
	c.AbortWithStatusJSON(http.StatusOK, Response{
		code,
		getTime(),
		data,
		//msg,
	})
}

func ErrorResponse(c *gin.Context, h int, code int, msg error, err ...error) {
	//组合error
	var t []string
	for _, e := range err {
		t = append(t, e.Error())
	}
	//记录日志
	if http.StatusInternalServerError == h {
		go zap.L().Error("ErrorResponse",
			zap.Strings("err", t),
			zap.String("msg", msg.Error()),
		)
	} else {
		go zap.L().Info("ErrorResponse",
			zap.Strings("err", t),
			zap.String("msg", msg.Error()),
		)
	}

	c.AbortWithStatusJSON(h, ErrResponse{
		code,
		getTime(),
		msg.Error(),
	})
}

func getTime() int64 {
	return time.Now().Unix()
}
