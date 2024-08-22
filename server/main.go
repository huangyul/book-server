package main

import (
	"github.com/huangyul/book-server/internal/pkg/bind"
	"github.com/spf13/viper"
)

func main() {

	initViper()

	server := InitWebServer()

	bind.InitTrans("zh")

	err := server.Run("127.0.0.1:8088")
	if err != nil {
		panic(err)
	}
}

func initViper() {
	viper.SetConfigName("book_server")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
