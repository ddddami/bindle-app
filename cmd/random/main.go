package main

import (
	"fmt"
	"log"

	"github.com/ddddami/bindle/random"
)

func main() {
	randomStr, err := random.Generate(random.Options{Length: 10})
	if err != nil {
		log.Printf("Failed to generate random string")
	}
	fmt.Println("Random string (10)", randomStr)
}
