package res

import (
	"encoding/json"
	"fmt"
)

type JsonResult struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type JsonError struct {
	code int
	err  error
}

func (j JsonError) Error() string {
	return j.err.Error()
}

func (j JsonError) String() string {
	return fmt.Sprintf("code:%d, message:%s", j.code, j.err.Error())
}

func (j JsonError) MarshalJSON() ([]byte, error) {
	return json.Marshal(&JsonResult{
		Success: false,
		Code:    j.code,
		Message: j.err.Error(),
		Data:    "",
	})
}

func JsonSuccess() *JsonResult {
	return &JsonResult{
		Success: true,
		Code:    0,
		Message: "success",
		Data:    nil,
	}
}

func JsonData(data interface{}) *JsonResult {
	return &JsonResult{
		Success: true,
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func NewJsonError(code int, err error) *JsonError {
	return &JsonError{
		code: code,
		err:  err,
	}
}
