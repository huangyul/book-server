package repository

import (
	"context"
	"time"

	"github.com/huangyul/book-server/internal/domain"
	"github.com/huangyul/book-server/internal/repository/dao"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (int64, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
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

func (u *userRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	dUser, err := u.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return u.toDoman(dUser), nil
}

func (u *userRepository) toEntity(user domain.User) dao.User {
	return dao.User{
		Password: user.Password,
		Username: user.Username,
	}
}

func (u *userRepository) toDoman(user dao.User) domain.User {
	return domain.User{
		ID:        user.ID,
		Password:  user.Password,
		Username:  user.Username,
		CreatedAt: time.UnixMilli(user.CreatedAt),
		UpdatedAt: time.UnixMilli(user.UpdatedAt),
	}
}
