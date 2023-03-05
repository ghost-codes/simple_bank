package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minKeySize = 32

type JWTMaker struct {
	secretKey string
}

// create and sign new token for user
func (jwtMaker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)

	return jwtToken.SignedString(jwtToken)
}

func (jwtMaker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, ErrorInvalidToken
		}
		return []byte(jwtMaker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		verr, ok := err.(jwt.ValidationError)

		if ok && errors.Is(verr.Inner, ErrorInvalidToken) {
			return nil, ErrorInvalidToken
		} else {
			return nil, ErrorExpiredToken

		}

	}
	payload, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return nil, ErrorInvalidToken
	}

	return payload, nil
}

// creates new JwtMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minKeySize {
		return nil, fmt.Errorf("Secret key must be at least %v characters", minKeySize)
	}
	return &JWTMaker{secretKey: secretKey}, nil
}
