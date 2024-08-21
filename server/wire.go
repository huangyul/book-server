//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/huangyul/book-server/internal/repository"
	"github.com/huangyul/book-server/internal/repository/dao"
	"github.com/huangyul/book-server/internal/service"
	"github.com/huangyul/book-server/internal/web"
	"github.com/huangyul/book-server/ioc"
)

var userSet = wire.NewSet(
	dao.NewGORMUserDao,
	repository.NewUserRepository,
	service.NewUserService,
	web.NewUserHandler,
)

func InitWebServer() *gin.Engine {
	wire.Build(

		// 第三方依赖
		ioc.InitDB,

		// user
		userSet,

		ioc.InitWeb,
	)
	return gin.New()
}
