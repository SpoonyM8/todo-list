package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"encoding/hex"
	"github.com/julienschmidt/httprouter"
)

const CONTENT_TYPE = "Content-Type"
const APPLICATION_JSON = "application/json"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleRegister(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// up to here, implement hashing of password and store salt + hashed pword in db
	var reqBody User
	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		http.Error(writer, "Error decoding User JSON", http.StatusBadRequest)
		return
	}

	username := reqBody.Username
	password := reqBody.Password
	response := fmt.Sprintf("Successfully registered user with Name=%s", username)
	fmt.Println(password)

	writer.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]string{"response": response})
}
