package handler

import "errors"

var (
	ErrInvalidParam   = errors.New("invalid param")
	ErrObjectNotFound = errors.New("object not found")
	ErrAuthFail       = errors.New("auth failed")
)
