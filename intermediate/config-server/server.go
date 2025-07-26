package main

import (
	"fmt"
	"os"
)

func main() {
	if secret, ok := os.LookupEnv("secret"); ok {
		fmt.Printf("secret: %v\n", secret)
	}
}
