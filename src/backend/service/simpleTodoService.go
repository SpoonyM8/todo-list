package service

import (
	"database/sql"
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
	var tokenString string = req.Header.Get("JWT")
	if err := util.VerifyJwt(tokenString); err != nil {
		util.ThrowUnauthorisedRequest(writer)
		return
	}

	var todo SimpleTodo
	json.NewDecoder(req.Body).Decode(&todo)

	if err := verifyCreateTodoRequestBody(todo); err != nil {
		util.ThrowBadRequest(writer, err.Error())
		return
	}

	username, err := util.GetUsernameFromJwt(tokenString)

	if err != nil {
		util.ThrowBadRequest(writer, err.Error())
		return
	}

	dbConn := req.Context().Value("db").(*sql.DB)

	var existingCount int
	err = dbConn.QueryRow(`
		SELECT count(*) FROM simple_todos st
		WHERE st.username=$1
		AND st.description=$2
		AND st.targetDate=$3
	`, username, todo.Description, todo.TargetDate).Scan(&existingCount)

	if err != nil || existingCount != 0 {
		fmt.Println(err)
		util.ThrowConflictRequest(writer, constants.GetTaskAlreadyDefinedMessage(todo.Description, todo.TargetDate))
		return
	}

	dbConn.Exec(`
		INSERT INTO simple_todos
		VALUES ($1, $2, $3)
	`, username, todo.Description, todo.TargetDate)
}

func verifyCreateTodoRequestBody(todo SimpleTodo) error {
	if len(todo.Description) == 0 {
		return fmt.Errorf(constants.DESCRIPTION_EMPTY)
	}

	if _, err := time.Parse("2006-01-02", todo.TargetDate); err != nil {
		return fmt.Errorf(constants.INVALID_TARGET_DATE)
	}

	if isDateBeforeNow(todo.TargetDate, time.Now().Format("2006-01-02")) {
		return fmt.Errorf(constants.INVALID_TARGET_DATE)
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
