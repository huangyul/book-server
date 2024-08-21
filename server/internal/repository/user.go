package repository

import (
	"context"

	"github.com/huangyul/book-server/internal/domain"
	"github.com/huangyul/book-server/internal/repository/dao"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (int64, error)
}

type userRepository struct {
	dao dao.UserDao
}

func NewUserRepository(dao dao.UserDao) UserRepository {
	return &userRepository{
		dao: dao,
	}
}

// Create 创建用户
func (u *userRepository) Create(ctx context.Context, user domain.User) (int64, error) {
	return u.dao.Create(ctx, u.toEntity(user))
}

func (u *userRepository) toEntity(user domain.User) dao.User {
	return dao.User{
		Password: user.Password,
		Username: user.Username,
	}
}
