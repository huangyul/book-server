package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/huangyul/book-server/internal/domain"
	"github.com/huangyul/book-server/internal/pkg/errno"
	"github.com/huangyul/book-server/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=user.go -destination=./mock/user_mock.go -package=mocksvc
type UserService interface {
	SignUp(ctx context.Context, username, password string) (int64, error)
	Profile(ctx context.Context, userID int64) (domain.User, error)
	Login(ctx context.Context, username, password string) error
}

type userService struct {
	repo repository.UserRepository
}

// Profile
func (svc *userService) Profile(ctx context.Context, userID int64) (domain.User, error) {
	return svc.repo.FindById(ctx, userID)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) SignUp(ctx context.Context, username, password string) (int64, error) {
	// 密码加密
	enPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errno.InternalServerError
	}
	return svc.repo.Create(ctx, domain.User{
		Username: username,
		Password: string(enPass),
	})
}

func (svc *userService) Login(ctx context.Context, username, password string) error {
	user, err := svc.repo.FindByName(ctx, username)
	if errors.Is(err, errno.UserNotFound) {
		return errno.BadRequest.SetMessage("账号或密码错误")
	}
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errno.BadRequest.SetMessage("账号或密码错误")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{})
	tokenStr, err := token.SignedString([]byte("secret"))
	if err != nil {
		return errno.InternalServerError
	}
	fmt.Print(tokenStr)
	return nil
}

type JwtClaims struct {
	UserID    int64
	UserAgent string
	jwt.RegisteredClaims
}
