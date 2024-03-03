// Kevin Le - KHL0041 - KevinLe2@my.unt.edu
// CSCE 3550 - Foundations of Cyber Security
// Project 1 - Basic JWKS Server

// This is a basic RESTful JWKS server developed in GO featuring an
// authentication endpoint, a handler for JWTs with expired keys,
// key expiry for enhanced security, and key deletion based on query
// parameters. This project will be improved upon in P2 and P3.


package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// start key store and generate initial key pair/store
	InitializeKeyStore()

	// Set up routing. Used Gorilla MUX for this
	r := mux.NewRouter()
	r.HandleFunc("/.well-known/jwks.json", JWKSHandler).Methods("GET")
	r.HandleFunc("/auth", AuthHandler).Methods("POST")

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}