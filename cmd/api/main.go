package main

import (
	"log"

	"go-library-rest-api/internal/config"
	"go-library-rest-api/internal/database"
	"go-library-rest-api/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresConnection(cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("database connected successfully")

	r := router.Setup(db)

	log.Printf("server is running on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
