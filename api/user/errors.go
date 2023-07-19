package user

import "errors"

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrPhoneNotFound           = errors.New("找不到手机号")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrEmailAlreadyExists      = errors.New("email already exists")
	ErrIncorrectPassword       = errors.New("incorrect password")
	ErrPhoneAlreadyExists      = errors.New("手机号已存在")
	ErrCodeNotFound            = errors.New("找不到验证码")
	ErrInvalidVerificationCode = errors.New("验证码错误")
)
