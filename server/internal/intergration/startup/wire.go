//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/huangyul/book-server/internal/pkg/jwt"
	"github.com/huangyul/book-server/internal/repository"
	"github.com/huangyul/book-server/internal/repository/dao"
	"github.com/huangyul/book-server/internal/service"
	"github.com/huangyul/book-server/internal/web"
	"github.com/huangyul/book-server/ioc"
)

var thirdParty = wire.NewSet(
	InitDB,
)

var userSet = wire.NewSet(
	dao.NewGORMUserDao,
	repository.NewUserRepository,
	service.NewUserService,
	web.NewUserHandler,
)

func InitWebServer() *gin.Engine {

	wire.Build(
		thirdParty,

		userSet,

		jwt.NewJwtService,

		ioc.InitWebMiddlewares,
		ioc.InitWeb,
	)

	return &gin.Engine{}
}

func InitUserHandler() *web.UserHandler {
	wire.Build(
		thirdParty,
		jwt.NewJwtService,
		userSet,
	)

	return &web.UserHandler{}
}
