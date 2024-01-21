package service

import (
	"encoding/json"
	"fmt"
	"main/constants"
	"main/util"
	"net/http"
	"strconv"
	"strings"
	"time"

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
		return
	}
}

func verifyCreateTodoRequestBody(todo SimpleTodo) error {
	if len(todo.Description) == 0 {
		return fmt.Errorf(constants.DESCRIPTION_EMPTY)
	}

	if _, err := time.Parse("2006-01-02", todo.TargetDate); err != nil {
		fmt.Println(err)
		return fmt.Errorf(constants.INVALID_TARGET_DATE)
	}

	if isDateBeforeNow(todo.TargetDate, time.Now().Format("2006-01-02")) {
		return fmt.Errorf("")
	}

	return nil
}

func isDateBeforeNow(dateToCheck, now string) bool {
	// No need to check for invalid numbers after splitting as time.Parse checked already
	dateSplit := strings.Split(dateToCheck, "-")
	year, _ := strconv.Atoi(dateSplit[0])
	month, _ := strconv.Atoi(dateSplit[1])
	day, _ := strconv.Atoi(dateSplit[2])

	nowSplit := strings.Split(now, "-")
	currYear, _ := strconv.Atoi(nowSplit[0])
	currMonth, _ := strconv.Atoi(nowSplit[1])
	currDay, _ := strconv.Atoi(nowSplit[2])

	if year < currYear {
		return true
	}

	if year == currYear {
		if month < currMonth {
			return true
		} else if month == currMonth {
			return day < currDay
		} else {
			return false
		}
	}

	return false
}
