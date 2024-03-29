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

type GetTodosRequestBody struct {
	TargetDate string `jasn:"targetDate"`
}

func HandleCreateSimpleTodo(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	username, err := verifyJwtAndUsername(req.Header.Get("JWT"), writer)
	if err != nil {
		util.ThrowUnauthorisedRequest(writer)
		return
	}

	var todo SimpleTodo
	json.NewDecoder(req.Body).Decode(&todo)

	if err := verifyValidTodo(todo); err != nil {
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

	writer.WriteHeader(http.StatusNoContent)
}

func HandleGetTodos(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	username, err := verifyJwtAndUsername(req.Header.Get("JWT"), writer)
	if err != nil {
		util.ThrowUnauthorisedRequest(writer)
		return
	}

	var reqBody GetTodosRequestBody
	json.NewDecoder(req.Body).Decode(&reqBody)

	if err := verifyGetTodosRequestBody(reqBody.TargetDate); err != nil {
		util.ThrowBadRequest(writer, err.Error())
		return
	}

	dbConn := req.Context().Value("db").(*sql.DB)

	rows, err := dbConn.Query(`
		SELECT description, targetDate from simple_todos
		WHERE username=$1
		AND targetDate=$2
	`, username, reqBody.TargetDate)

	fmt.Println(err)

	var resp []SimpleTodo

	for rows.Next() {
		var row SimpleTodo
		rows.Scan(&row.Description, &row.TargetDate)
		resp = append(resp, row)
	}

	writer.Header().Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(resp)
}

func HandleDeleteSimpleTodo(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	username, err := verifyJwtAndUsername(req.Header.Get("JWT"), writer)
	if err != nil {
		util.ThrowUnauthorisedRequest(writer)
		return
	}

	var todo SimpleTodo
	json.NewDecoder(req.Body).Decode(&todo)

	if err := verifyValidTodo(todo); err != nil {
		util.ThrowBadRequest(writer, err.Error())
		return
	}

	dbConn := req.Context().Value("db").(*sql.DB)

	dbConn.Exec(`
		DELETE FROM simple_todos
		WHERE username=$1
		AND description=$2
		AND targetDate=$3
	`, username, todo.Description, todo.TargetDate)

	writer.WriteHeader(http.StatusNoContent)
}

func verifyGetTodosRequestBody(dateToCheck string) error {
	if _, err := time.Parse("2006-01-02", dateToCheck); err != nil {
		return fmt.Errorf(constants.INVALID_TARGET_DATE)
	}
	if isDateBeforeNow(dateToCheck, time.Now().Format("2006-01-02")) {
		return fmt.Errorf(constants.INVALID_TARGET_DATE)
	}
	return nil
}

func verifyValidTodo(todo SimpleTodo) error {
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
	dateSplit := strings.Split(dateToCheck, "-")
	if len(dateSplit) != 3 {
		return false
	}
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

func verifyJwtAndUsername(token string, writer http.ResponseWriter) (string, error) {
	if err := util.VerifyJwt(token); err != nil {
		return "", err
	}
	username, err := util.GetUsernameFromJwt(token)
	return username, err
}
