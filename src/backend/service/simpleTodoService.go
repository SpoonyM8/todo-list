package service

import (
	"encoding/json"
	"fmt"
	"main/util"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

type SimpleTodo struct {
	Description string `json:"description"`
	TargetDate  string `json:"targetDate"`
}

// @TODO: Need to get the username from decoding JWT after verifying it

func HandleCreateSimpleTodo(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if err := util.VerifyJWT(req.Header.Get("JWT")); err != nil {
		util.ThrowUnauthorisedRequest(writer)
		return
	}

	var todo SimpleTodo
	json.NewDecoder(req.Body).Decode(&todo)

	if err := verifyCreateTodoRequestBody(todo); err != nil {
		util.ThrowBadRequest(writer, err.Error())
	}

}

func verifyCreateTodoRequestBody(todo SimpleTodo) error {
	if len(todo.Description) == 0 {
		return fmt.Errorf("description field must not be empty")
	}
	if match, _ := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2}$", todo.TargetDate); !match {
		return fmt.Errorf("targetDate field must be of the format YYYY-MM-DD")
	}

	return nil
}
