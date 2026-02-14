package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"dvarapala/internal/db"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "dvarapala.db" // default local path
	}

	client, err := db.NewSQLiteClient(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer client.Close()

	fmt.Printf("Database initialized at %s\n", dbPath)

	port := ":8080"
	fmt.Printf("Starting server on port %s...\n", port)
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Dvarapala!")
	})

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
