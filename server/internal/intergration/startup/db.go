package startup

import (
	"github.com/huangyul/book-server/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"))

	if err != nil {
		panic(err)
	}

	// 构建表
	err = db.AutoMigrate(&dao.User{})
	if err != nil {
		panic(err)
	}

	return db
}
