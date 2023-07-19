package product

import (
	"auction-website/conf"
	//errors2 "auction-website"
	"auction-website/internal/global"
	"auction-website/utils"
	"auction-website/utils/req"
	"auction-website/utils/resp"
	"errors"
	"fmt"
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

// createHandler godoc
// @Summary 创建产品
// @Schemes
// @Description Create a new product
// @Tags product
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param product body product.Product true "产品信息"
// @Success 200 {object} resp.Response{} "创建成功返回产品ID"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 500 {object} resp.ErrResponse{} "服务器内部错误"
// @Router /product [post]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) createHandler(c *gin.Context) {
	uid := req.GetUid(c)
	request, err := req.GetRequestWithValidator(c, &Product{}, "json")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	p := request.(*Product)
	//todo 这个如果用户修改了名字 jwt还没过期 name就是错的  解决方法是改名后刷新jwt
	//todo 也可在 JWT 中不包含用户名，而是在每次请求时通过其他方式来获取用户名，例如在数据库中进行查询。
	//写入uid
	p.UserID = uid
	pid, err := h.Service.CreateProduct(p)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, err)
		return
	}
	resp.DataResponse(c, pid)
}

// listHandler godoc
// @Summary 获取产品列表
// @Schemes
// @Description Get a list of products with pagination
// @Tags product
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} resp.Response{} "获取成功返回产品列表"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 500 {object} resp.ErrResponse{} "服务器内部错误"
// @Router /product/list [get]
func (h *Handler) listHandler(c *gin.Context) {
	page, pageSize, err := utils.GetPageParameter(c)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	lists, err := h.Service.GetProductList(page, pageSize)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, err)
		return
	}
	data := utils.JsonPagination(lists, page, pageSize)
	resp.DataResponse(c, data)
}

// detailHandler godoc
// @Summary 获取产品详情
// @Schemes
// @Description Get product detail by ID
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} resp.Response{} "成功返回产品详情"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 404 {object} resp.ErrResponse{} "产品不存在"
// @Failure 500 {object} resp.ErrResponse{} "服务器内部错误"
// @Router /product/{id} [get]
func (h *Handler) detailHandler(c *gin.Context) {
	request, err := req.GetRequestWithValidator(c, &global.IDValidator{}, "uri")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	//fmt.Println(request)
	g := request.(*global.IDValidator)
	fmt.Println(g)
	pid, err := utils.StringToUint32(g.ID)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, errors.New("invalid id number"), err)
		return
	}
	productDetail, err := h.Service.GetProductDetail(pid)
	if err != nil {
		if errors.Is(err, global.ErrNotFound) {
			resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, err)
		} else {
			resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, global.ErrInternal, err)
		}
		return
	}
	// Return the product detail
	resp.DataResponse(c, productDetail)
}

// updateHandler godoc
// @Summary 更新商品信息
// @Schemes
// @Description Update product information
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param body body product.UpdateProduct true "Updated product data"
// @Success 200 {object} resp.Response{} "更新成功返回成功消息"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 401 {object} resp.ErrResponse{} "未授权的访问"
// @Failure 404 {object} resp.ErrResponse{} "商品不存在"
// @Failure 500 {object} resp.ErrResponse{} "服务器内部错误"
// @Router /product/{id} [put]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) updateHandler(c *gin.Context) {
	uid := req.GetUid(c)
	id := c.Param("id")
	pid, err := utils.StringToUint32(id)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	request, err := req.GetRequestWithValidator(c, &UpdateProduct{}, "json")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	p := request.(*UpdateProduct)
	//fmt.Println(p)
	// Check if the user is authorized to update the product
	//获取商品信息
	productDetail, err := h.Service.GetProductDetail(pid)
	if err != nil {
		if errors.Is(err, global.ErrNotFound) {
			resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, err)
		} else {
			resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, global.ErrInternal, err)
		}
		return
	}
	//fmt.Println(productDetail)
	if productDetail.UserID != uid {
		resp.ErrorResponse(c, http.StatusUnauthorized, resp.ERROR, global.ErrUnauthorized, errors.New("恶意篡改行为"))
		return
	}
	// Call product service to update the product
	p.ID = productDetail.ID
	//fmt.Println(p)
	err = h.Service.UpdateProduct(p)

	if err != nil {
		if errors.Is(err, global.ErrNotUpdate) {
			resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, err)
		} else {
			resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, global.ErrInternal, err)
		}
		return
	}
	resp.DataResponse(c, "Succeed")
}

// deleteHandler godoc
// @Summary 删除产品
// @Schemes
// @Description Delete a product by ID
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} resp.Response{} "成功删除产品"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 401 {object} resp.ErrResponse{} "未授权的操作"
// @Failure 404 {object} resp.ErrResponse{} "产品不存在"
// @Failure 500 {object} resp.ErrResponse{} "服务器内部错误"
// @Router /product/{id} [delete]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) deleteHandler(c *gin.Context) {
	request, err := req.GetRequestWithValidator(c, &global.IDValidator{}, "uri")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	g := request.(*global.IDValidator)
	pid, err := utils.StringToUint32(g.ID)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}

	productToDelete, err := h.Service.GetProductDetail(pid)
	if err != nil {
		if errors.Is(err, global.ErrNotFound) {
			resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, err)
		} else {
			resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, global.ErrInternal, err)
		}
		return
	}

	// Check if the user is authorized to delete the product
	uid := req.GetUid(c)
	if productToDelete.UserID != uid {
		resp.ErrorResponse(c, http.StatusUnauthorized, resp.ERROR, global.ErrUnauthorized)
		return
	}

	// Call product service to delete the product
	err = h.Service.DeleteProduct(pid)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, errors.New("delete error"), err)
		return
	}
	// Return success message after deleting the product
	resp.DataResponse(c, "DeleteProduct success")
}
