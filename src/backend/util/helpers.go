package util

import (
	"fmt"
	"main/constants"

	"github.com/dgrijalva/jwt-go"
)

var SECRET_KEY = []byte("SO DAMN SECRET!!!") // @TODO: store a key securely somewhere, obviously not as a const here. Hashicorp Vault? SSM?

func VerifyJwt(tokenString string) error {
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

func GetUsernameFromJwt(tokenString string) (string, error) {
	// guaranteed no error as already checked in VerifyJwt
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})

	var claims jwt.MapClaims
	claims, _ = token.Claims.(jwt.MapClaims)
	for k, v := range claims {
		if k == "sub" {
			return v.(string), nil
		}
	}

	return "", fmt.Errorf(constants.INVALID_CLAIMS)
}
