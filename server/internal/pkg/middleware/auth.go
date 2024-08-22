package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/huangyul/book-server/internal/pkg/jwt"
)

type AuthMiddlewareBuilder struct {
	// 白名单
	whitePaths []string

	j jwt.JWT
}

func NewAuthMiddlewareBuilder(j jwt.JWT) *AuthMiddlewareBuilder {
	return &AuthMiddlewareBuilder{
		j: j,
	}
}

func (b *AuthMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, p := range b.whitePaths {
			if strings.HasPrefix(ctx.Request.URL.Path, p) {
				ctx.Next()
				return
			}
		}

		userID, err := b.j.AuthJWT(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		ctx.Set("userID", userID)

		ctx.Next()
	}
}

func (b *AuthMiddlewareBuilder) AddWhitePath(paths ...string) *AuthMiddlewareBuilder {
	b.whitePaths = append(b.whitePaths, paths...)
	return b
}
