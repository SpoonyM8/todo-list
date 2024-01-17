package service

import (
	"fmt"
	"main/constants"
	"main/util"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

func HandleCreateSimpleTodo(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if err := verifyJWT(req.Header.Get("JWT")); err != nil {
		util.ThrowUnauthorisedRequest(writer)
		return
	}

}

func verifyJWT(tokenString string) error {
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
