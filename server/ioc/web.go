package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/book-server/internal/pkg/jwt"
	"github.com/huangyul/book-server/internal/pkg/middleware"
	"github.com/huangyul/book-server/internal/web"
)

func InitWeb(u *web.UserHandler, mdls []gin.HandlerFunc) *gin.Engine {
	s := gin.New()

	s.Use(mdls...)

	u.RegisterRoutes(s)

	return s
}

func InitWebMiddlewares(j jwt.JWT) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.NewAuthMiddlewareBuilder(j).AddWhitePath("/user/login", "/user/signup").Build(),
	}
}
