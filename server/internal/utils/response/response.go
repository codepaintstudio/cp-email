package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg,omitempty"`
}

const (
	SUCCESS                   = 0
	ERROR                     = 500
	INVALID_PARAMS            = 400
	UNAUTHORIZED              = 401
	FORBIDDEN                 = 403
	NOT_FOUND                 = 404
	FIELD_IS_EMPTY            = 90001
	ACCOUNT_OR_PASSWORD_ERROR = 90002
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: SUCCESS,
		Data: data,
	})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: ERROR,
		Msg:  msg,
	})
}