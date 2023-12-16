package util

import (
	"encoding/json"
	"main/constants"
	"net/http"
)

func ThrowBadRequest(writer http.ResponseWriter, errorMessage string) {
	setContentTypeJson(writer)
	writer.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(writer).Encode(map[string]string{"error": errorMessage})
}

func ThrowUnauthorisedRequest(writer http.ResponseWriter) {
	setContentTypeJson(writer)
	writer.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(writer).Encode(map[string]string{"error": constants.UNAUTHORISED})
}

func setContentTypeJson(writer http.ResponseWriter) {
	writer.Header().Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
}