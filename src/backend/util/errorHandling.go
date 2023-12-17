package util

import (
	"encoding/json"
	"main/constants"
	"net/http"
)

func ThrowBadRequest(writer http.ResponseWriter, errorMessage string) {
	setContentTypeJson(writer)
	setStatusCode(writer, http.StatusBadRequest)
	setErrorMessage(writer, errorMessage)
}

func ThrowUnauthorisedRequest(writer http.ResponseWriter) {
	setContentTypeJson(writer)
	setStatusCode(writer, http.StatusUnauthorized)
	setErrorMessage(writer, constants.UNAUTHORISED)
}

func ThrowConflictRequest(writer http.ResponseWriter, errorMessage string) {
	setContentTypeJson(writer)
	setStatusCode(writer, http.StatusConflict)
	setErrorMessage(writer, errorMessage)
}

func setStatusCode(writer http.ResponseWriter, status int) {
	writer.WriteHeader(status)
}

func setContentTypeJson(writer http.ResponseWriter) {
	writer.Header().Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
}

func setErrorMessage(writer http.ResponseWriter, errorMessage string) {
	json.NewEncoder(writer).Encode(map[string]string{"error": errorMessage})
}
