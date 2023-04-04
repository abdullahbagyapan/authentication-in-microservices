package app

import (
	"log"
	"microservices/token-service/internal/ports"
)

type Adapter struct {
	core ports.CorePort
}

func NewAppAdapter(core ports.CorePort) *Adapter {
	return &Adapter{core: core}
}

func (A Adapter) GenerateToken(id string) (string, error) {

	token, err := A.core.GenerateNewToken(id)

	if err != nil {
		log.Printf("Error generating token %v", err)
		return "", err
	}
	return token, nil
}

func (A Adapter) ValidateToken(token string) (string, bool, error) {

	jwt, claims, err := A.core.GetJWTAndClaims(token)

	if err != nil {
		log.Printf("Error parsing token %v", err)
		return "", false, err
	}

	isValid := A.core.ValidateToken(jwt)

	if !isValid {
		return "", false, nil
	}

	userId := A.core.GetUserIdByClaims(claims)

	return userId, isValid, nil

}
