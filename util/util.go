package util

import (
	"encoding/json"
	"net/http"
)

/*
ContentType - "Content-Type" header: Define the message content-type
*/
const ContentType string = "Content-Type"

/*
ResponseOk - Facade http handler to return a successful result
*/
func ResponseOk(w http.ResponseWriter, body interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set(ContentType, "application/json")
	json.NewEncoder(w).Encode(body)
}

/*
ResponseError - Facade http handler to return an error code and message
*/
func ResponseError(w http.ResponseWriter, httpCode int, message string) {
	w.WriteHeader(httpCode)
	w.Header().Set(ContentType, "application/json")

	body := map[string]string{
		"error": message,
	}

	json.NewEncoder(w).Encode(body)
}
