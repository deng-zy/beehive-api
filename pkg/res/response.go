package res

import (
	"encoding/json"
	"fmt"
)

// JSONResult api response body
type JSONResult struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// JSONError api response error
type JSONError struct {
	code int
	err  error
}

func (j JSONError) Error() string {
	return j.err.Error()
}

func (j JSONError) String() string {
	return fmt.Sprintf("code:%d, message:%s", j.code, j.err.Error())
}

// MarshalJSON JSONError json encode
func (j JSONError) MarshalJSON() ([]byte, error) {
	return json.Marshal(&JSONResult{
		Success: false,
		Code:    j.code,
		Message: j.err.Error(),
		Data:    "",
	})
}

// JSONSuccess api reseponse success
func JSONSuccess() *JSONResult {
	return &JSONResult{
		Success: true,
		Code:    0,
		Message: "success",
		Data:    nil,
	}
}

// JSONData api response data
func JSONData(data interface{}) *JSONResult {
	return &JSONResult{
		Success: true,
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

// NewJSONError new api error
func NewJSONError(code int, err error) error {
	return &JSONError{
		code: code,
		err:  err,
	}
}
