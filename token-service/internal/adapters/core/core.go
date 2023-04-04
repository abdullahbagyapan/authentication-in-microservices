package core

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type Adapter struct {
}

const SECRET_KEY string = "KWQw4HyFUjtPu4KCXWxSTvzZCJ4ZJ8639aJybQBeaPRVjLwzH3uw7qWGMskKLxVg"

type JWTData struct {
	jwt.StandardClaims
}

func NewCoreAdapter() *Adapter {
	return new(Adapter)
}

func (A Adapter) GenerateNewToken(id string) (string, error) {

	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 6).Unix(), // 6 hours
			Id:        id,
		},
	}

	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenString.SignedString([]byte(SECRET_KEY))

	return token, err
}

func (A Adapter) ValidateToken(jwt *jwt.Token) bool {
	return jwt.Valid
}

func (A Adapter) GetUserIdByClaims(claims *JWTData) string {

	userId := claims.Id

	return userId
}

func (A Adapter) GetJWTAndClaims(token string) (*jwt.Token, *JWTData, error) {
	claims := &JWTData{}

	jwt, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return nil, nil, err
	}

	return jwt, claims, nil

}
