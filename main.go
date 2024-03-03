package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize key store and generate an initial key pair.
	InitializeKeyStore()

	// Set up routing.
	r := mux.NewRouter()
	r.HandleFunc("/.well-known/jwks.json", JWKSHandler).Methods("GET")
	r.HandleFunc("/auth", AuthHandler).Methods("POST")

	// Start the server.
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
