package main

import (
	"log"

	"github.com/mznrasil/my-blogs-be/internal/storage"
)

const (
	PORT           = ":8080"
	dataSourceName = "host=localhost port=5432 dbname=myblogs user=rasil password="
)

func main() {
	log.Println("Connecting to database...")
	db, err := storage.NewPostgresStore(dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to database!")
	defer db.DB.Close()

	server := NewAPIServer(PORT, db.DB)
	server.Run()
}
