package main

import "github.com/huangyul/book-server/internal/pkg/bind"

func main() {

	server := InitWebServer()

	bind.InitTrans("zh")

	err := server.Run("127.0.0.1:8088")
	if err != nil {
		panic(err)
	}
}
