package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
)

var (
	AuthorizedPrivateKey *rsa.PrivateKey
	expiredPrivateKey    *rsa.PrivateKey
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
const authorizedKID = "AuthorizedGoodKeyID"

func InitializeKeyStore() {
	generateKeys()
}

func generateKeys() {
	// Generate a good key pair
	var err error
	AuthorizedPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Error generating RSA keys: %v", err)
	}

	// Generate an expired key pair for demonstration purposes
	expiredPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Error generating expired RSA keys: %v", err)
	}
}

func JWKSHandler(w http.ResponseWriter, r *http.Request) {
	jwk := generateJWK(AuthorizedPrivateKey.Public().(*rsa.PublicKey), authorizedKID)
	resp := JWKS{
		Keys: []JWK{jwk},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	KID       string `json:"kid"`
	Algorithm string `json:"alg"`
	KeyType   string `json:"kty"`
	Use       string `json:"use"`
	N         string `json:"n"`
	E         string `json:"e"`
}

func generateJWK(pubKey *rsa.PublicKey, kid string) JWK {
	return JWK{
		KID:       kid,
		Algorithm: "RS256",
		KeyType:   "RSA",
		Use:       "sig",
		N:         base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes()),
		E:         base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pubKey.E)).Bytes()),
	}
}
