package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

type ClientAuth struct {
	Name     string `json:"name"`
	ClientID uint64 `json:"client_id"`
	jwt.StandardClaims
}

type AuthRequest struct {
	ClientID uint64 `json:"client_id" form:"client_id" xml:"client_id" binding:"required"`
	Secret   string `json:"secret" form:"secret" xml:"secret" binding:"required"`
}

const EXPIRES = 720 * time.Hour //30 days
const TOKEN_LOOKUP = "Authorization"
const TOKEN_HEAD_NAME = "Bearer"

var secret = []byte("MP^go7Kx&eHvsEMafpB66vFp")

func IssueToken(ID uint64, name string) (string, error) {
	now := time.Now().Unix()
	claims := &ClientAuth{
		Name:     name,
		ClientID: ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(EXPIRES).Unix(),
			NotBefore: now,
			IssuedAt:  now,
			Subject:   name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenString string) (*ClientAuth, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ClientAuth{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	client, ok := token.Claims.(*ClientAuth)
	if !ok {
		return nil, errors.New("invalid token.")
	}

	return client, nil
}

func ReFreshToken(client *ClientAuth) (string, error) {
	client.StandardClaims.ExpiresAt = time.Now().Add(EXPIRES).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, client)
	return token.SignedString(secret)
}
