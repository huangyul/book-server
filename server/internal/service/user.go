package service

import (
	"context"
	"errors"

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
	Login(ctx context.Context, username, password string) (int64, error)
	Edit(ctx context.Context, userID int64, username string) error
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

func (svc *userService) Login(ctx context.Context, username, password string) (int64, error) {
	user, err := svc.repo.FindByName(ctx, username)
	if errors.Is(err, errno.UserNotFound) {
		return 0, errno.UserPasswordIncorrect
	}
	if err != nil {
		return 0, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, errno.UserPasswordIncorrect
	}
	return user.ID, nil
}

func (svc *userService) Edit(ctx context.Context, userID int64, username string) error {
	_, err := svc.repo.FindById(ctx, userID)
	if errors.Is(err, errno.UserNotFound) {
		return errno.UserNotFound
	}
	if err != nil {
		return err
	}
	return svc.repo.UpdateById(ctx, userID, domain.User{
		Username: username,
	})
}

type JwtClaims struct {
	UserID    int64
	UserAgent string
	jwt.RegisteredClaims
}
