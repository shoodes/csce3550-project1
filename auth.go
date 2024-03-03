package main

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt" // Ensure you're using the correct version
)

/*
	var KeyPair KeyPair
	var kid string
	mu.RLock()
	for k, v := range keys {

		kid = k
		KeyPair = v
		break
	}
	mu.RUnlock()
*/

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	expired, _ := strconv.ParseBool(r.URL.Query().Get("expired"))
	var signingKey *rsa.PrivateKey
	var kid string

	// Choose the key based on the 'expired' query parameter
	if expired {
		signingKey = expiredPrivateKey // Use the globally defined expired key
		kid = "expiredKeyID"           // A static ID for the expired key, adjust as needed
	} else {
		signingKey = AuthorizedPrivateKey // Use the globally defined good key
		kid = authorizedKID               // Use the static ID for the good key
	}

	// Create a new token with the specified claims, including the dynamic expiration
	claims := jwt.MapClaims{
		"iss": "jwks-server",
		"sub": "user123",
	}
	if expired {
		claims["exp"] = time.Now().Add(-1 * time.Hour).Unix() // Expired token
	} else {
		claims["exp"] = time.Now().Add(1 * time.Hour).Unix() // Valid token
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid

	// Sign the token
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		http.Error(w, "Failed to sign token", http.StatusInternalServerError)
		return
	}

	// Respond with the token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

/*
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	expired := r.URL.Query().Get("expired") == "true"

	var KeyPair KeyPair
	var kid string

	mu.RLock()
	for k, v := range keys {

		kid = k
		KeyPair = v
		break
	}
	mu.RUnlock()

	if expired {
		KeyPair.Expiry = time.Now().Add(-1 * time.Hour) //negative number to show expired
		//KeyPair.IsExpired = false
		//mu.Lock()
		//delete(keys, kid) // Assuming 'kid' identifies the key used for this operation
		//mu.Unlock()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "jwks-server",
		"sub": "user123",
		"exp": KeyPair.Expiry.Unix(),
		//"exp": KeyPair.Expiry == time.Now().Add(-1*time.Hour),
	})
	token.Header["kid"] = kid

	tokenString, err := token.SignedString(KeyPair.PrivateKey)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func encodeSegment(seg []byte) string {
	return jwt.EncodeSegment(seg)
}
*/
