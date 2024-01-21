package util

import (
	"fmt"
	"main/constants"

	"github.com/dgrijalva/jwt-go"
)

var SECRET_KEY = []byte("SO DAMN SECRET!!!") // @TODO: store a key securely somewhere, obviously not as a const here. Hashicorp Vault? SSM?

func VerifyJWT(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf(constants.INVALID_JWT)
	}

	_, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return fmt.Errorf(constants.INVALID_CLAIMS)
	}

	return nil
}