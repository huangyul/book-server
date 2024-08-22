package ioc

import (
	"fmt"

	"github.com/huangyul/book-server/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	type Config struct {
		Host     string
		Port     int
		User     string
		Password string
	}
	var cfg Config
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/book_server?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}

	dao.InitTables(db)

	return db
}
