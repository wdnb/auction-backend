package req

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func GetRequestWithValidator(c *gin.Context, req any, should string) (any, error) {
	reqValue := reflect.ValueOf(req)

	if reqValue.Kind() != reflect.Ptr {
		return nil, errors.New("请求对象必须是指针类型")
	}

	reqType := reqValue.Elem().Type()
	reqValue = reflect.New(reqType)

	if should == "json" {
		if err := c.ShouldBindJSON(reqValue.Interface()); err != nil {
			return nil, err
		}
	} else if should == "uri" {
		if err := c.ShouldBindUri(reqValue.Interface()); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("绑定类型错误")
	}
	if err := validator.New().Struct(reqValue.Interface()); err != nil {
		return nil, err
	}

	return reqValue.Interface(), nil
}

func GetUid(c *gin.Context) uint32 {
	uid := c.MustGet("uid").(uint32)
	return uid
}
