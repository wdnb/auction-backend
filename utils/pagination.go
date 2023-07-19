package utils

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin"
)

type PaginationQ struct {
	//Ok   bool   `json:"ok"`
	Size uint32 `form:"size" json:"size"`
	Page uint32 `form:"page" json:"page"`
	List any    `json:"list" ` // save pagination list
	//Total uint        `json:"total"`
}

func GetPageParameter(c *gin.Context) (pageInt, pageSizeInt uint32, err error) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	//fmt.Println(page, pageSize)
	pageInt, err = StringToUint32(page)
	if err != nil {
		return 0, 0, err
	}
	pageSizeInt, err = StringToUint32(pageSize)
	if err != nil {
		return 0, 0, err
	}

	if pageSizeInt > 10 {
		return 0, 0, errors.New("pageSize must less than 10")
	}

	return pageInt, pageSizeInt, nil

}

func JsonPagination(lists any, page uint32, pageSize uint32) PaginationQ {

	//ok := true
	vi := reflect.ValueOf(lists)
	if IsBlank(vi) {
		//ok = false
		lists = make([]int, 0)
	}
	data := PaginationQ{
		//Ok:   ok,
		Page: page,
		Size: pageSize,
		List: lists,
	}
	return data
}

func Offset(pageNumber, itemsPerPage uint32) uint32 {
	offset := (pageNumber - 1) * itemsPerPage
	return offset
}
