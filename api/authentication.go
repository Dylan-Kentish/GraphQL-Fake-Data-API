package api

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type authenticationProvider struct {
	signingMethod jwt.SigningMethod
	secretKey     string
}

func NewAuthenticationProvider() IAuthenticationProvider {
	return &authenticationProvider{
		signingMethod: jwt.SigningMethodHS256,
		secretKey:     "This is very secret!",
	}
}

type IAuthenticationProvider interface {
	GetToken(id int) (string, error)
}

func (auth *authenticationProvider) GetToken(id int) (string, error) {
	claims := jwt.NewWithClaims(auth.signingMethod, jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(id),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
	})

	token, err := claims.SignedString([]byte(auth.secretKey))

	if err != nil {
		return "", err
	}

	return token, err
}
