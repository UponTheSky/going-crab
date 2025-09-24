package main

import (
	"log"
)

func main() {
	server := NewServer("tcp", 8080)
	app := NewApp()

	if err := server.Run(app); err != nil {
		log.Fatal(err)
	}
}
