package web

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/book-server/internal/pkg/errno"
	"github.com/huangyul/book-server/internal/service"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (u *UserHandler) RegisterRoutes(g *gin.Engine) {
	ug := g.Group("/user")
	{
		// 注册接口
		ug.POST("/signup", u.SignUp)
	}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	type Req struct {
		Username        string `json:"username" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteResultErr(c, errno.BadRequest.SetMessage(err.Error()))
		return
	}
	if req.Password != req.ConfirmPassword {
		WriteResultErr(c, errno.BadRequest.SetMessage("两次密码不一致"))
		return
	}
	_, err := h.svc.SignUp(c, req.Username, req.Password)
	if err != nil {
		WriteResultErr(c, errno.InternalServerError.SetMessage(err.Error()))
		return
	}
	WrtieSuccess(c)
}
