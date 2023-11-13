package main

import (
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	if err := store.Init(); err != nil {
		log.Fatal("Error initializing the database: ", err)
	}

	server := NewAPIServer(":8080", store)
	server.Run()
}
