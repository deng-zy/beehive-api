package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gordon-zhiyong/beehive-api/pkg/bytesconv"
	"github.com/gordon-zhiyong/beehive-api/pkg/conf"
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

func IssueToken(ID uint64, name string) (string, error) {
	now := time.Now().Unix()
	expires := conf.Auth.GetDuration("expires")

	claims := &ClientAuth{
		Name:     name,
		ClientID: ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expires).Unix(),
			NotBefore: now,
			IssuedAt:  now,
			Subject:   name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(bytesconv.StringToBytes(conf.Auth.GetString("secret")))
}

func ParseToken(tokenString string) (*ClientAuth, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ClientAuth{}, func(token *jwt.Token) (interface{}, error) {
		return bytesconv.StringToBytes(conf.Auth.GetString("secret")), nil
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
	expires := conf.Auth.GetDuration("expires")
	client.StandardClaims.ExpiresAt = time.Now().Add(expires).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, client)
	return token.SignedString(bytesconv.StringToBytes(conf.Auth.GetString("secret")))
}
