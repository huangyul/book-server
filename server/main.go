package main

func main() {

	server := InitWebServer()

	err := server.Run("127.0.0.1:8088")
	if err != nil {
		panic(err)
	}
}
