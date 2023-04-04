package core

import (
	"github.com/google/uuid"
	"log"
	"testing"
)

var Token string

func TestGenerateNewToken(t *testing.T) {

	coreAdapter := NewCoreAdapter()

	token, err := coreAdapter.GenerateNewToken(uuid.New().String())

	if err != nil {
		log.Fatalf("Error generating new token %v", err)
	}

	Token = token

}

func TestValidateToken(t *testing.T) {

	coreAdapter := NewCoreAdapter()

	jwt, _, err := coreAdapter.GetTokenAndClaims(Token)

	if err != nil {
		log.Fatalf("Error parsing token %v", err)
	}

	isValid := coreAdapter.ValidateToken(jwt)

	log.Printf("Token is valid = %v", isValid)

}
