package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huangyul/book-server/internal/pkg/errno"
)

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func WriteResultErr(c *gin.Context, err *errno.Errno) {
	c.JSON(http.StatusOK, Result[any]{
		Code: err.Code,
		Msg:  err.Msg,
	})
}

func WriteResult[T any](c *gin.Context, code int, msg string, data T) {
	c.JSON(http.StatusOK, Result[T]{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func WrtieSuccess(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Result[any]{})
}
