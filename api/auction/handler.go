package auction

import (
	"auction-website/conf"
	"auction-website/internal/global"
	"auction-website/utils"
	"auction-website/utils/req"
	"auction-website/utils/resp"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(c *conf.Config) *Handler {
	return &Handler{
		Service: NewService(c),
	}
}

// @Summary 获取所有拍卖品
// @Description 获取分页后的所有拍卖品
// @Tags auction
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Success 200 {object} resp.Response{data=utils.PaginationQ} "拍卖品列表"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 500 {object} resp.ErrResponse{} "内部错误"
// @Router /auction [get]
func (h *Handler) getAllAuctions(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := "20"
	// Convert query parameters to integers
	pageInt, err := utils.StringToUint32(page)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	pageSizeInt, err := utils.StringToUint32(pageSize)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, errors.New("invalid page size"), err)
		return
	}

	lists, err := h.Service.getAllAuctions(pageInt, pageSizeInt)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, err)
		return
	}
	data := utils.JsonPagination(lists, pageInt, pageSizeInt)
	resp.DataResponse(c, data)
}

// @Summary 创建拍卖品
// @Description 创建一个新的拍卖品
// @Tags auction
// @Accept json
// @Produce json
// @Param auction body auction.Auction true "拍卖品信息"
// @Success 200 {object} resp.Response{} "拍卖品id"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 401 {object} resp.ErrResponse{} "没有权限"
// @Failure 500 {object} resp.ErrResponse{} "内部错误"
// @Router /auction [post]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) createAuction(c *gin.Context) {
	uid := req.GetUid(c)
	request, err := req.GetRequestWithValidator(c, &Auction{}, "json")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	a := request.(*Auction)
	//fmt.Println(p)
	//fmt.Println(uid)
	//todo 检测是否有权限创建拍品
	// 1:从请求中获取当前用户的信息，判断用户是否有权限创建拍卖
	// 2:检查用户是否为管理员或拥有创建拍卖的特殊权限，如果不是则返回错误信息
	//app := c.MustGet("app").(*utils.Components)
	a.ProductUID = uid
	id, err := h.Service.CreateAuction(a)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, errors.New("发布拍品失败"), err)
		return
	}
	resp.DataResponse(c, id)
}

// @Summary 获取拍卖品详情
// @Description 根据ID获取拍卖品详情
// @Tags auction
// @Accept json
// @Produce json
// @Param id path int true "拍卖品ID"
// @Success 200 {object} resp.Response{data=auction.Auction} "拍卖品详情"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 404 {object} resp.ErrResponse{} "拍卖品不存在"
// @Failure 500 {object} resp.ErrResponse{} "内部错误"
// @Router /auction/{id} [get]
func (h *Handler) getAuction(c *gin.Context) {
	id := c.Param("id")
	idInt, err := utils.StringToUint32(id)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	a, err := h.Service.GetAuction(idInt)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, errors.New("获取拍品信息失败"), err)
		return
	}
	resp.DataResponse(c, a)
}

// @Summary 更新拍卖品
// @Description 更新指定ID的拍卖品信息
// @Tags auction
// @Accept json
// @Produce json
// @Param id path int true "拍卖品ID"
// @Param auction body auction.Update true "拍卖品信息"
// @Success 200 {object} resp.Response{} "成功"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 401 {object} resp.ErrResponse{} "未认证"
// @Failure 403 {object} resp.ErrResponse{} "禁止访问"
// @Failure 404 {object} resp.ErrResponse{} "拍卖品不存在"
// @Failure 500 {object} resp.ErrResponse{} "内部错误"
// @Router /auction/{id} [put]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) updateAuction(c *gin.Context) {
	id := c.Param("id")
	idInt, err := utils.StringToUint32(id)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	request, err := req.GetRequestWithValidator(c, &Update{}, "json")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	a := request.(*Update)
	a.ID = idInt
	// TODO: 从请求中获取当前用户的信息，判断用户是否有权限更新拍卖
	err = h.Service.UpdateAuction(a)
	if err != nil {
		if errors.Is(err, global.ErrNotUpdate) {
			resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		} else {
			resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, global.ErrInternal, err)
		}
		return
	}
	resp.DataResponse(c, "更新拍品成功")
}

// @Summary 删除拍卖品
// @Description 删除指定ID的拍卖品
// @Tags auction
// @Accept json
// @Produce json
// @Param id path int true "拍卖品ID"
// @Success 200 {object} resp.Response{} "成功"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 401 {object} resp.ErrResponse{} "未认证"
// @Failure 403 {object} resp.ErrResponse{} "禁止访问"
// @Failure 404 {object} resp.ErrResponse{} "拍卖品不存在"
// @Failure 500 {object} resp.ErrResponse{} "内部错误"
// @Router /auction/{id} [delete]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) deleteAuction(c *gin.Context) {
	id := c.Param("id")
	idInt, err := utils.StringToUint32(id)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	// 从请求中获取当前用户的信息，判断用户是否有权限删除拍卖
	// 检查用户是否为管理员或拥有删除拍卖的特殊权限，如果不是则返回错误信息
	//if !isAdmin(c) {
	//	resp.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "没有权限删除拍品！")
	//	return
	//}

	// isAdmin函数的实现
	//func isAdmin(c *gin.Context) bool {
	//    role := c.GetString("role")
	//    return role == "admin"
	//}
	err = h.Service.DeleteAuction(idInt)
	if err != nil {
		if errors.Is(err, global.ErrNotDelete) {
			resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		} else {
			resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, global.ErrInternal, err)
		}
		return
	}
	resp.DataResponse(c, "删除拍品成功！")
}

// @Summary 获取拍卖品的所有竞拍
// @Description 获取指定拍卖品的竞拍列表,支持分页
// @Tags auction
// @Accept json
// @Produce json
// @Param id path int true "拍卖品ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} utils.PaginationQ "成功"
// @Failure 400 {object} resp.ErrResponse "请求错误"
// @Failure 404 {object} resp.ErrResponse "拍卖品不存在"
// @Failure 500 {object} resp.ErrResponse "内部错误"
// @Router /auction/{id}/bid [get]
func (h *Handler) getAllBids(c *gin.Context) {
	id := c.Param("id")
	idInt, err := utils.StringToUint32(id)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	page := c.DefaultQuery("page", "1")
	pageSize := "20"
	// Convert query parameters to integers
	pageInt, err := utils.StringToUint32(page)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	pageSizeInt, err := utils.StringToUint32(pageSize)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}

	lists, err := h.Service.getAllBids(idInt, pageInt, pageSizeInt)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, err)
		return
	}
	data := utils.JsonPagination(lists, pageInt, pageSizeInt)
	resp.DataResponse(c, data)
}

// @Summary 创建竞拍
// @Description 为指定拍卖品创建新的竞拍信息
// @Tags auction
// @Accept json
// @Produce json
// @Param id path int true "拍卖品ID"
// @Param bid body auction.Bid true "竞拍信息"
// @Success 200 {object} resp.Response "成功"
// @Failure 400 {object} resp.ErrResponse "请求错误"
// @Failure 401 {object} resp.ErrResponse "未认证"
// @Failure 403 {object} resp.ErrResponse "禁止访问"
// @Failure 404 {object} resp.ErrResponse "拍卖品不存在"
// @Failure 409 {object} resp.ErrResponse "竞拍冲突"
// @Failure 500 {object} resp.ErrResponse "内部错误"
// @Router /auction/{id}/bid [post]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
// todo 互斥锁 防止脏数据
func (h *Handler) createBid(c *gin.Context) {
	uid := req.GetUid(c)
	id := c.Param("id")
	request, err := req.GetRequestWithValidator(c, &Bid{}, "json")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	bid := request.(*Bid)
	idInt, err := utils.StringToUint32(id)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	// 检查用户是否有权限竞拍
	// 1:检查用户是否为拍品的创建者，如果是则返回错误信息
	// 2:检查用户是否已经竞拍过该拍品，如果是则返回错误信息
	// 3:检查用户的竞拍价格是否高于当前最高价，如果不是则返回错误信息
	a := Bid{
		CustomerID: uid,
		AuctionID:  idInt,
		BidPrice:   bid.BidPrice,
	}
	//fmt.Println(a)
	_, err = h.Service.CreateBid(&a)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, errors.New("出价失败"), err)
		return
	}
	resp.DataResponse(c, "出价成功")
}
