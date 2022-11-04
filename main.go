package main

import (
	"log"

	"github.com/flash-cards-vocab/backend/internal/api"
)

func main() {
	server, err := api.NewServer()
	if err != nil {
		log.Panicln("Failed to Initialized postgres DB:", err)
	}
	server.Start()
}
