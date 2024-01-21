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

	var dbConn *sql.DB = req.Context().Value("db").(*sql.DB)

	if isUsernameTaken(dbConn, username) {
		util.ThrowConflictRequest(writer, constants.USERNAME_TAKEN)
		return
	}

	salt, hashedAndSaltedPassword := generateHashAndSalt(password)

	dbConn.Exec(`
		INSERT INTO user_schema.users
		VALUES ($1, $2, $3)
	`, username, hashedAndSaltedPassword, salt)

	returnJwt(writer, username)
}

func HandleLogin(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	username, password, err := verifyAuthRequestBody(req)
	if err != nil {
		util.ThrowBadRequest(writer, err.Error())
		return
	}

	var dbConn *sql.DB = req.Context().Value("db").(*sql.DB)

	var storedPassword, storedSalt string

	dbConn.QueryRow(`
		SELECT u.password, u.salt
		FROM user_schema.users u
		WHERE u.username=$1
	`, username).Scan(&storedPassword, &storedSalt)

	if !doPasswordsMatch(storedPassword, password, storedSalt) {
		util.ThrowBadRequest(writer, constants.INVALID_USER_REQUEST)
		return
	}

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

func isUsernameTaken(dbConn *sql.DB, username string) bool {
	var usernameCount int

	dbConn.QueryRow(`
		SELECT count(*) FROM user_schema.users u
		WHERE u.username = $1	
	`, username).Scan(&usernameCount)

	return usernameCount != 0
}

func doPasswordsMatch(storedPassword string, loginPassword string, salt string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginPassword+salt)); err != nil {
		return false
	}
	return true
}

func generateJWT(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 3).Unix(),
	})
	tokenString, _ := token.SignedString(util.SECRET_KEY)
	return tokenString
}

func returnJwt(writer http.ResponseWriter, username string) {
	writer.Header().Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]string{"JWT": generateJWT(username)})
}

func generateHashAndSalt(password string) (string, string) {
	salt, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedAndSaltedPassword, _ := bcrypt.GenerateFromPassword([]byte(password+string(salt)), bcrypt.DefaultCost)
	return string(hashedAndSaltedPassword), string(salt)
}
