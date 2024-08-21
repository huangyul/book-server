package dao

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/huangyul/book-server/internal/pkg/errno"
	"gorm.io/gorm"
)

type UserDao interface {
	Create(ctx context.Context, user User) (int64, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByName(ctx context.Context, username string) (User, error)
}

type GORMUserDao struct {
	db *gorm.DB
}

// Create
func (g *GORMUserDao) Create(ctx context.Context, user User) (int64, error) {
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now
	err := g.db.WithContext(ctx).Create(&user).Error

	if strings.Contains(err.Error(), "1062") {
		return 0, errno.UserAlreadyExist
	}
	if err != nil {
		return 0, errno.InternalServerError

	}
	return user.ID, nil
}

func (g *GORMUserDao) FindById(ctx context.Context, id int64) (User, error) {
	var user User
	err := g.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errno.UserNotFound
	}
	if err != nil {
		return user, errno.InternalServerError
	}
	return user, nil
}

func (g *GORMUserDao) FindByName(ctx context.Context, username string) (User, error) {
	var user User
	err := g.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errno.UserNotFound
	}
	if err != nil {
		return user, errno.InternalServerError
	}
	return user, nil
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
