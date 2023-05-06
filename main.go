package main

import (
	"final/route"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := route.SetupAPI(); err != nil {
		log.Fatalln(err)
	}
}
