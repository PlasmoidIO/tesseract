package main

import "share/central/server"

func main() {
	s := server.CreateServer()
	s.Listen(8080)
}