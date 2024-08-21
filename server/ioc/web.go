package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/book-server/internal/web"
)

func InitWeb(u *web.UserHandler) *gin.Engine {
	s := gin.New()

	u.RegisterRoutes(s)

	return s
}
