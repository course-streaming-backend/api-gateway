package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("warning: no .env file found: %v", err)
	}
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}

	app := &app{
		addr: ":" + addr,
	}

	app.run()
}
