package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/mznrasil/my-blogs-be/internal/storage"
)

const PORT = ":8080"

var dataSourceName = os.Getenv("DATA_SOURCE_NAME")

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
