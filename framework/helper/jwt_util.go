package helper

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserTypeEnum int

const (
	BUYER_TYPE UserTypeEnum = iota
	SELLER_TYPE
)

type UserJWTPayload struct {
	ID    int64
	Email string
	Name  string
	Type  UserTypeEnum
}

type JWTResponse struct {
	Data  interface{}  `json:"data"`
	Type  UserTypeEnum `json:"type"`
	Token string       `json:"token"`
}

func GenerateToken(payload *UserJWTPayload, secret []byte) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["id"] = payload.ID
	atClaims["email"] = payload.Email
	atClaims["name"] = payload.Name
	atClaims["type"] = payload.Type
	atClaims["exp"] = time.Now().Add(time.Minute * 180).Unix()

	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := sign.SignedString(secret)

	if err != nil {
		return "", err
	}
	return token, nil
}

func VerifyToken(tokenString string, secret []byte) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != t.Method {
			return nil, errors.New("unexpected signing method")
		}

		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	return parsedToken, nil
}

func ValidateToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	verifiedToken, err := VerifyToken(tokenString, secret)
	if err != nil {
		return nil, err
	}

	tokenClaims, ok := verifiedToken.Claims.(jwt.MapClaims)
	if !ok && !verifiedToken.Valid {
		return nil, errors.New("Unauthorized")
	}

	return tokenClaims, nil
}
