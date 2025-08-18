package main

import (
	"log"

	"github.com/Mihailo84/stakeholders-service/config"
	"github.com/Mihailo84/stakeholders-service/startup"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env variables")
	}

	cfg := config.NewConfig()
	server := startup.NewServer(cfg)
	server.Start()
}
