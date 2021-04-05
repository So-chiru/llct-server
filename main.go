package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/so-chiru/llct-server/route"

	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	address := os.Getenv("ADDRESS")
	log.Println("Running on port " + address)

	http.ListenAndServe(address, route.Router())
}
