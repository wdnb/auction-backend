package account

import (
	"auction-website/conf"
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

// Deposit 处理用户充值请求,增加用户余额
func (h *Handler) Deposit(c *gin.Context) {
	uid := req.GetUid(c)
	request, err := req.GetRequestWithValidator(c, &Amount{}, "json")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	a := request.(*Amount)
	// 增加用户余额 TODO 这玩意其实要从payment进来
	if err := h.Service.Deposit(uid, a.Amount); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, err)
		return
	}
	resp.DataResponse(c, "deposit successful")
}

// GetBalance godoc
// @Summary 查询用户余额
// @Description 查询用户账户的余额
// @Tags account
// @Accept json
// @Produce json
// @Success 200 {object} resp.Response "成功"
// @Failure 500 {object} resp.ErrResponse "内部错误"
// @Router /account/balance [get]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) GetBalance(c *gin.Context) {
	// 查询用户当前的账户余额
	uid := req.GetUid(c)
	balance, err := h.Service.GetBalance(uid)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, errors.New("查询用户余额失败"), err)
		return
	}
	resp.DataResponse(c, balance)
}

// Withdraw godoc
// @Summary 用户提现
// @Description 处理用户提现请求,减少用户余额,发起转账
// @Tags account
// @Accept json
// @Produce json
// @Param request body account.Amount true "提现信息"
// @Success 200 {object} resp.Response "成功"
// @Failure 400 {object} resp.ErrResponse "请求错误"
// @Failure 500 {object} resp.ErrResponse "内部错误"
// @Router /account/withdraw [post]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) Withdraw(c *gin.Context) {
	uid := req.GetUid(c)
	request, err := req.GetRequestWithValidator(c, &Amount{}, "json")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	a := request.(*Amount)
	// 减少用户余额
	//TODO 使用error.new 完善错误信息
	if err := h.Service.Withdraw(uid, a.Amount); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, errors.New("申请提现失败"), err)
		return
	}
	// 发起转账 TODO 这里其实应该是调用 payment 的发起转账接口
	resp.DataResponse(c, "申请提现成功")
}

// WithdrawalRecord godoc
// @Summary 用户提现记录
// @Description 获取用户的提现记录
// @Tags account
// @Accept  json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} utils.PaginationQ "成功"
// @Failure 500 {object} resp.ErrResponse "内部错误"
// @Router /account/withdrawal-record [get]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
// TODO 提现失败回滚数据 发送通知（成功也发）
func (h *Handler) WithdrawalRecord(c *gin.Context) {
	uid := req.GetUid(c)

	page, pageSize, err := utils.GetPageParameter(c)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}

	lists, err := h.Service.WithdrawalRecord(uid, page, pageSize)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, errors.New("查询提现记录失败"), err)
		return
	}
	data := utils.JsonPagination(lists, page, pageSize)

	resp.DataResponse(c, data)
}
