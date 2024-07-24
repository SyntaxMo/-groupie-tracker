package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /", rootHandler)
	http.HandleFunc("GET /artists/", artistHandler)

	fmt.Println("Listening on port 8080")
	log.Fatalf("Error starting server: %v", http.ListenAndServe(":8080", nil))
}
