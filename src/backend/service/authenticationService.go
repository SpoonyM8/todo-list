package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/constants"
	"main/util"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY = []byte("SO DAMN SECRET!!!") // @TODO: store a key securely somewhere, obviously not as a const here. Hashicorp Vault? SSM?

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleRegister(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	username, password, err := verifyAuthRequestBody(req)
	if err != nil {
		util.ThrowBadRequest(writer, err.Error())
		return
	}

	// @TODO: check username isn't taken in DB, else return 409. Username will be unique as user ID, so cannot be duplicates
	salt, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedAndSaltedPassword, _ := bcrypt.GenerateFromPassword([]byte(password+string(salt)), bcrypt.DefaultCost)
	fmt.Println(hashedAndSaltedPassword)
	// @TODO: store salt and hash in DB for user.

	var dbConn *sql.DB = req.Context().Value("db").(*sql.DB)
	dbConn.Exec(`
		INSERT INTO users (username, password, salt)
		VALUES ($1, $2, $3)
	`, username, hashedAndSaltedPassword, salt)

	fmt.Println(dbConn.Query("SELECT * FROM users"))
	returnJwt(writer, username)
}

func HandleLogin(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	/*
		JWT verification should not be done for HandleLogin. It should return a JWT instead

		if jwtError := verifyJWT(req.Header.Get(constants.ACCESS_TOKEN)); jwtError != nil {
			util.ThrowUnauthorisedRequest(writer)
			return
		}
	*/

	username, password, err := verifyAuthRequestBody(req)
	if err != nil {
		util.ThrowBadRequest(writer, err.Error())
		return
	}
	fmt.Println(username, password)
	/*
		@TODO: implement after DB is setup with username+password
		if er := verifyPassword(string(hashedAndSaltedPassword), password, string(salt)); er != nil {
			util.ThrowBadRequest(writer, constants.INVALID_USER_REQUEST)
			return
		}
	*/
	//
	returnJwt(writer, username)
}

func verifyAuthRequestBody(req *http.Request) (string, string, error) {
	var reqBody User
	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		return "", "", err
	}

	username := reqBody.Username
	password := reqBody.Password

	if len(username) < 6 || len(password) < 6 {
		err = fmt.Errorf(constants.USERNAME_LENGTH_VALIDATION)
		return username, password, err
	}

	return username, password, nil
}

func verifyPassword(storedPassword string, loginPassword string, salt string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginPassword+salt)); err != nil {
		return err
	}
	return nil
}

func generateJWT(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 3).Unix(),
	})
	fmt.Println(token)
	tokenString, _ := token.SignedString(SECRET_KEY)
	fmt.Println(tokenString)
	return tokenString
}

/*
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
*/

func returnJwt(writer http.ResponseWriter, username string) {
	writer.Header().Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]string{"JWT": generateJWT(username)})
}
