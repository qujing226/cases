package main

func main() {
	server := InitWebServer()
	_ = server.Run(":8080")
}
