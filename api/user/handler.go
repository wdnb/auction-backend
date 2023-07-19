package user

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

//func (h *Handler) registerHandler(c *gin.Context) {
//	//userService := NewUserService()
//	var u user.User
//	if err := c.ShouldBindJSON(&u); err != nil {
//		resp.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid request body"), err)
//		return
//	}
//	nameIsExist, err := h.Service.UserNameIsExist(&u)
//	if err != nil {
//		resp.ErrorResponse(c, http.StatusBadRequest, ErrUserAlreadyExists, err)
//		return
//	}
//	if nameIsExist {
//		resp.ErrorResponse(c, http.StatusBadRequest, ErrUserAlreadyExists)
//		return
//	}
//	emailIsExist, err := h.Service.UserEmailIsExist(&u)
//	if err != nil {
//		resp.ErrorResponse(c, http.StatusBadRequest, ErrEmailAlreadyExists, err)
//		return
//	}
//	if emailIsExist {
//		resp.ErrorResponse(c, http.StatusBadRequest, ErrEmailAlreadyExists)
//		return
//	}
//	token, err := h.Service.CreateUser(&u)
//	if err != nil {
//		resp.ErrorResponse(c, http.StatusInternalServerError, global.ErrInternal, err)
//		return
//	}
//
//	resp.DataResponse(c, gin.H{"token": token})
//}

// getUserCenter godoc
// @Summary 获取用户中心信息
// @Schemes
// @Description Get user center information
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} resp.Response{} "请求成功，返回用户信息和用户详情"
// @Failure 404 {object} resp.ErrResponse{} "用户不存在"
// @Failure 500 {object} resp.ErrResponse{} "服务器内部错误"
// @Router /user/center [get]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) getUserCenter(c *gin.Context) {
	uid := req.GetUid(c)
	exist, err := h.Service.CheckUserExist(uid)
	if err != nil {
		resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrUserNotFound, err)
		return
	}
	if exist == false {
		resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrUserNotFound, err)
		return
	}
	u, err := h.Service.GetUserByID(uid)
	if err != nil {
		resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrUserNotFound, err)
		return
	}
	addresses, err := h.Service.GetUserAddresses(uid)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, global.ErrInternal, err)
		return
	}

	userInfo := map[string]any{
		"addresses": addresses,
	}

	resp.DataResponse(c, gin.H{"user": u, "user_info": userInfo})
}

// loginHandler godoc
// @Summary 用户登录
// @Schemes
// @Description User login with username and password
// @Tags user
// @Accept json
// @Produce json
// @Param user body user.Login true "登录信息"
// @Success 200 {object} resp.Response{} "登录成功返回token"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 401 {object} resp.ErrResponse{} "密码错误"
// @Failure 404 {object} resp.ErrResponse{} "用户不存在"
// @Failure 500 {object} resp.ErrResponse{} "服务器内部错误"
// @Router /user/login [post]
func (h *Handler) loginHandler(c *gin.Context) {
	u, err := req.GetRequestWithValidator(c, &Login{}, "json")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	cl := u.(*Login)
	token, err := h.Service.Login(cl.Username, cl.Password)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, err)
		} else if errors.Is(err, ErrIncorrectPassword) {
			resp.ErrorResponse(c, http.StatusUnauthorized, resp.ERROR, err)
		} else {
			resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, global.ErrInternal, err)
		}
		return
	}
	resp.DataResponse(c, token)
}

// @Summary 手机号登录
// @Schemes
// @Description 通过发送的验证码登陆
// @Tags user
// @Accept json
// @Produce json
// @Param 验证码 body user.VerificationCode true "通过验证码登录"
// @Success 200 {object} resp.Response{} "登录成功返回token"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 404 {object} resp.ErrResponse{} "用户不存在或者验证码错误"
// @Failure 500 {object} resp.ErrResponse{} "服务器内部错误"
// @Router /user/login-phone [post]
func (h *Handler) loginPhoneHandler(c *gin.Context) {

	request, err := req.GetRequestWithValidator(c, &VerificationCode{}, "json")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	clu := request.(*VerificationCode)
	token, err := h.Service.LoginByCode(clu, c.ClientIP())
	if err != nil {
		if errors.Is(err, ErrPhoneNotFound) {
			resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrPhoneNotFound)
		} else if errors.Is(err, ErrUserNotFound) {
			resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrUserNotFound)
		} else if errors.Is(err, ErrInvalidVerificationCode) {
			resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrInvalidVerificationCode)
		} else if errors.Is(err, ErrCodeNotFound) {
			resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrCodeNotFound)
		} else {
			resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, global.ErrInternal, err)
		}
		return
	}

	resp.DataResponse(c, token)

}

// getUserProfile godoc
// @Summary 获取用户资料
// @Schemes
// @Description Get user profile by ID
// @Tags user
// @Accept json
// @Produce json
// @Param uid path string true "User ID"
// @Success 200 {object} resp.Response{} "成功返回用户资料"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 404 {object} resp.ErrResponse{} "用户不存在"
// @Router /user/profile/{uid} [get]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) getUserProfile(c *gin.Context) {
	id := c.Param("uid")
	uid, err := utils.StringToUint32(id)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}

	u, err := h.Service.GetUserByID(uid)
	if err != nil {
		resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrUserNotFound, err)
		return
	}
	resp.DataResponse(c, u)
}

// updateUserProfile godoc
// @Summary 更新用户资料
// @Description Update user profile by user ID
// @Tags user
// @Accept json
// @Produce json
// @Param uid path string true "User ID"
// @Param body body user.UpdateUser true "User profile data"
// @Success 200 {object} user.UpdateUser "Updated user profile data"
// @Failure 400 {object} resp.ErrResponse{} "Invalid request or malformed JSON"
// @Failure 404 {object} resp.ErrResponse{} "User not found"
// @Failure 500 {object} resp.ErrResponse{} "Internal server error"
// @Router /user/profile/{uid} [put]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) updateUserProfile(c *gin.Context) {
	//userService := NewUserService()
	id := c.Param("uid")
	uid, err := utils.StringToUint32(id)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	exist, err := h.Service.CheckUserExist(uid)
	if err != nil {
		resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrUserNotFound, err)
		return
	}
	if exist == false {
		resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrUserNotFound, err)
		return
	}
	u, err := req.GetRequestWithValidator(c, &UpdateUser{}, "json")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	uu := u.(*UpdateUser)
	err = h.Service.UpdateUserByUserID(uid, uu)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, global.ErrInternal, err)
		return
	}

	resp.DataResponse(c, "修改成功")
}

// getUserAddresses godoc
// @Summary 获取用户地址
// @Description 获取用户地址列表
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {array} resp.Response{} "获取成功"
// @Failure 401 {object} resp.ErrResponse{} "未授权或Token过期"
// @Failure 500 {object} resp.ErrResponse{} "服务器内部错误"
// @Router /user/shipping-address [get]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) getUserAddresses(c *gin.Context) {
	uid := req.GetUid(c)
	addresses, err := h.Service.GetUserAddresses(uid)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, global.ErrInternal, err)
		return
	}
	resp.DataResponse(c, addresses)
}

// createShippingAddress godoc
// @Summary 创建收货地址
// @Schemes
// @Description Create a new shipping address for a user
// @Tags user
// @Accept json
// @Produce json
// @Param Address body user.CreateShippingAddress true "收货地址信息"
// @Success 200 {object} resp.Response{} "创建成功返回地址ID"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 404 {object} resp.ErrResponse{} "用户不存在"
// @Failure 500 {object} resp.ErrResponse{} "服务器内部错误"
// @Router /user/shipping-address [post]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) createShippingAddress(c *gin.Context) {
	uid := req.GetUid(c)
	request, err := req.GetRequestWithValidator(c, &CreateShippingAddress{}, "json")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	s := request.(*CreateShippingAddress)
	exist, err := h.Service.CheckUserExist(uid)
	if err != nil {
		resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrUserNotFound, err)
		return
	}
	if exist == false {
		resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrUserNotFound, err)
		return
	}
	s.UserID = uid

	sid, err := h.Service.CreateShippingAddress(s)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, errors.New("can not create address"), err)
		return
	}
	resp.DataResponse(c, sid)
}

// updateUserAddressByID godoc
// @Summary 更新用户地址信息
// @Description Update the shipping address of a user by ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "地址的主键id"
// @Param ShippingAddress body user.ShippingAddress true "Shipping address information"
// @Success 200 {object} user.ShippingAddress "Updated shipping address"
// @Failure 400 {object} resp.ErrResponse{} "Bad request"
// @Failure 401 {object} resp.ErrResponse{} "Unauthorized"
// @Failure 404 {object} resp.ErrResponse{} "User not found"
// @Failure 500 {object} resp.ErrResponse{} "Internal server error"
// @Router /user/shipping-address/{id} [put]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) updateUserAddressByID(c *gin.Context) {
	id := c.Param("id")
	sid, err := utils.StringToUint32(id)
	uid := req.GetUid(c)
	//s := c.MustGet("shipping_address").(user.ShippingAddress)
	request, err := req.GetRequestWithValidator(c, &ShippingAddress{}, "json")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	s := request.(*ShippingAddress)
	//fmt.Println(s)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	exist, err := h.Service.CheckUserExist(uid)
	if err != nil {
		resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrUserNotFound, err)
		return
	}
	if exist == false {
		resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrUserNotFound, err)
		return
	}
	s.UserID = uid
	s.ID = sid
	err = h.Service.UpdateUserAddressByID(s)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, errors.New("can not update address"), err)
		return
	}
	resp.DataResponse(c, s)
}

// deleteUserAddressByID godoc
// @Summary 删除用户地址
// @Schemes
// @Description Delete a user address by ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "Address ID"
// @Success 200 {object} resp.Response{} "地址删除成功"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 404 {object} resp.ErrResponse{} "用户不存在"
// @Failure 500 {object} resp.ErrResponse{} "服务器内部错误"
// @Router /user/shipping-address/{id} [delete]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) deleteUserAddressByID(c *gin.Context) {
	id := c.Param("id")
	sid, err := utils.StringToUint32(id)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	uid := req.GetUid(c)
	exist, err := h.Service.CheckUserExist(uid)
	if err != nil {
		resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrUserNotFound, err)
		return
	}
	if exist == false {
		resp.ErrorResponse(c, http.StatusNotFound, resp.ERROR, ErrUserNotFound, err)
		return
	}
	err = h.Service.DeleteUserAddressByID(uid, sid)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, errors.New("can not delete address"), err)
		return
	}
	resp.DataResponse(c, "Address deleted successfully")
}

//
//func (h *Handler) updateUserCenter(c *gin.Context) {
//	uid := req.GetUid(c)
//	exist, err := h.Service.CheckUserExist(uid)
//	if err != nil {
//		resp.ErrorResponse(c, http.StatusNotFound, ErrUserNotFound, err)
//		return
//	}
//	if exist == false {
//		resp.ErrorResponse(c, http.StatusNotFound, ErrUserNotFound, err)
//		return
//	}
//	var u user.UpdateUser
//	if err := c.ShouldBindJSON(&u); err != nil {
//		resp.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid request body"), err)
//		return
//	}
//	err = h.Service.UpdateUserByUserID(uid, &u)
//	if err != nil {
//		resp.ErrorResponse(c, http.StatusInternalServerError, global.ErrInternal, err)
//		return
//	}
//	addresses, err := h.Service.GetUserAddresses(uid)
//	if err != nil {
//		resp.ErrorResponse(c, http.StatusInternalServerError, global.ErrInternal, err)
//		return
//	}
//	resp.DataResponse(c, gin.H{"user": u, "addresses": addresses})
//}

// getVerificationCode godoc
// @Summary 获得验证码
// @Schemes
// @Description Get verification code for a specific kind
// @Tags user
// @Accept json
// @Produce json
// @Param kind path string true "Kind of verification code" Enums(login,register,reset_password)
// @Param Phone body user.CodeSendUser true "手机号"
// @Success 200 {object} resp.Response{} "登录成功返回token"
// @Failure 400 {object} resp.ErrResponse{} "非法请求"
// @Failure 404 {object} resp.ErrResponse{} "用户不存在或者验证码错误"
// @Failure 500 {object} resp.ErrResponse{} "服务器内部错误"
// @Router /user/verification-code/{kind} [post]
func (h *Handler) getVerificationCode(c *gin.Context) {
	request, err := req.GetRequestWithValidator(c, &KindValidator{}, "uri")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	kind := request.(*KindValidator)
	a, err := req.GetRequestWithValidator(c, &CodeSendUser{}, "json")
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	csu := a.(*CodeSendUser)

	_, err = h.Service.SendVerificationCode(csu.Phone, kind.Kind)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, resp.ERROR, global.ErrInternal, err)
		return
	}
	//todo 接入短信运营商
	resp.DataResponse(c, "验证码发送成功")
}
