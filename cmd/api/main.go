package main

import (
	"log"

	"go-library-rest-api/internal/router"
)

func main() {
	r := router.Setup()

	log.Println("server is running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
