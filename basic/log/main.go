package main

import (
	"log"
	"os"
)

func main() {
	logger := log.New(
		os.Stderr,
		"[Going Crab] ",
		log.LstdFlags|log.Lshortfile|log.LUTC,
	)

	logger.Println("%v, test!", 42)
}
