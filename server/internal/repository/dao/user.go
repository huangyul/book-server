package dao

import (
	"context"
	"errors"

	"github.com/huangyul/book-server/internal/pkg/errno"
	"gorm.io/gorm"
)

type UserDao interface {
	Create(ctx context.Context, user User) (int64, error)
}

type GORMUserDao struct {
	db *gorm.DB
}

// Create
func (g *GORMUserDao) Create(ctx context.Context, user User) (int64, error) {
	err := g.db.WithContext(ctx).Create(&user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return 0, errno.UserAlreadyExist
	}
	if err != nil {
		return 0, errno.InternalServerError

	}
	return user.ID, nil
}

func NewGORMUserDao(db *gorm.DB) UserDao {
	return &GORMUserDao{
		db: db,
	}
}

type User struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	Username  string `gorm:"uniqueIndex;type:varchar(255)"`
	Password  string
	CreatedAt int64
	UpdatedAt int64
}
