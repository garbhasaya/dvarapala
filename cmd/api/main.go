package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := ":8080"
	fmt.Printf("Starting server on port %s...\n", port)
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Dvarapala!")
	})

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
