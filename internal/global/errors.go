package global

import "errors"

var (
	// Common errors
	ErrNotFound      = errors.New("目标不存在")
	ErrNotUpdate     = errors.New("目标不存在或已更新")
	ErrNotDelete     = errors.New("目标不存在或已删除")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInvalidInput  = errors.New("invalid input")

	// Permission errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	// Internal errors
	ErrInternal = errors.New("internal server error")
)
