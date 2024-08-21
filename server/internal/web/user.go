package web

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huangyul/book-server/internal/pkg/bind"
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
		// 获取详情
		ug.GET("/profile/:id", u.Profile)
	}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	type Req struct {
		Username        string `json:"username" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}
	var req Req
	if err := bind.Bind(c, &req); err != nil {
		WriteResultErr(c, errno.BadRequest.SetMessage(err.Error()))
		return
	}
	if req.Password != req.ConfirmPassword {
		WriteResultErr(c, errno.BadRequest.SetMessage("两次密码不一致"))
		return
	}
	uId, err := h.svc.SignUp(c, req.Username, req.Password)
	if err != nil {
		WriteResultErr(c, errno.InternalServerError.SetMessage(err.Error()))
		return
	}
	WriteResult(c, uId)
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteResultErr(ctx, errno.BadRequest.SetMessage("非法id"))
		return
	}
	dUser, err := h.svc.Profile(ctx, id)
	if err != nil {
		WriteResultErr(ctx, errno.InternalServerError.SetMessage(err.Error()))
		return
	}
	user := UserResp{
		ID:        dUser.ID,
		Username:  dUser.Username,
		CreatedAt: dUser.CreatedAt.Format(time.DateOnly),
		UpdatedAt: dUser.UpdatedAt.Format(time.DateOnly),
	}
	WriteResult(ctx, user)
}

type UserResp struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
