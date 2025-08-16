package main

import (
	"log"

	"github.com/VincentSamuelPaul/production-api/api"
	"github.com/VincentSamuelPaul/production-api/database"
)

func main() {
	store, err := database.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	server := api.NewAPIServer(":3000", store)
	server.Run()
}
