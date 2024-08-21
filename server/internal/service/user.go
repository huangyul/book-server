package service

import (
	"context"

	"github.com/huangyul/book-server/internal/domain"
	"github.com/huangyul/book-server/internal/pkg/errno"
	"github.com/huangyul/book-server/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=user.go -destination=./mock/user_mock.go -package=mocksvc
type UserService interface {
	SignUp(ctx context.Context, username, password string) (int64, error)
	Profile(ctx context.Context, userID int64) (domain.User, error)
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
