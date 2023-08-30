package src

import (
	"log"
	"os"

	"github.com/GokdenizCakir/stant_oyun/src/routes"
	"github.com/joho/godotenv"
)

func Init() {
	r := routes.Router
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("failed to run server: %v\n", err)
	}
}
