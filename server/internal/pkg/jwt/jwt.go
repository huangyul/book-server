package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWT interface {
	// GenerateToken 通过 userid 生成 accessToken 和 refreshToken
	GenerateToken(userID int64) (string, string, error)
	// AuthJWT 验证 token
	AuthJWT(ctx *gin.Context) (int64, error)
	// RefreshToken 刷新accessToken
	RefreshToken(ctx *gin.Context) (string, error)
}

var _ JWT = (*jwtService)(nil)

type jwtService struct {
	key string
}

func NewJwtService() JWT {
	return &jwtService{
		key: "kBdeQtuBfgwSDD8VuMCjYsHofFJAUNFC",
	}
}

// generateToken 通过 userid 生成 accessToken 和 refreshToken
func (j *jwtService) GenerateToken(userID int64) (string, string, error) {
	accessClaims := LoginClaims{
		userID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
		},
	}
	accessToken, err := j.createToken(accessClaims)
	if err != nil {
		return "", "", fmt.Errorf("生成token失败，原因：%s", err.Error())
	}
	refreshClaims := LoginClaims{
		userID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	refreshToken, err := j.createToken(refreshClaims)
	if err != nil {
		return "", "", fmt.Errorf("生成token失败，原因：%s", err.Error())
	}
	return accessToken, refreshToken, nil
}

// AuthJWT 验证 token
func (j *jwtService) AuthJWT(ctx *gin.Context) (int64, error) {
	tokenStr := ctx.GetHeader("Authorization")
	if tokenStr == "" {
		return 0, fmt.Errorf("token为空")
	}
	strs := strings.Split(tokenStr, " ")
	if len(strs) != 2 {
		return 0, fmt.Errorf("token格式错误")
	}
	token, err := jwt.ParseWithClaims(strs[1], &LoginClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.key), nil
	})
	if err != nil {
		return 0, fmt.Errorf("token验证失败，原因：%s", err.Error())
	}
	if !token.Valid {
		return 0, fmt.Errorf("token验证失败")
	}
	return token.Claims.(*LoginClaims).userID, nil
}

// RefreshToken 刷新accessToken
func (j *jwtService) RefreshToken(ctx *gin.Context) (string, error) {
	userID, err := j.AuthJWT(ctx)
	if err != nil {
		return "", err
	}
	tokenClaims := LoginClaims{
		userID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
		},
	}
	return j.createToken(tokenClaims)
}

func (j *jwtService) createToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(j.key))
}

type LoginClaims struct {
	userID int64
	jwt.RegisteredClaims
}
