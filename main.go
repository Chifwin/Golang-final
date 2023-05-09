package main

import (
	"golang-final/route"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := route.SetupAPI(":8080"); err != nil {
		log.Fatalln(err)
	}
}
