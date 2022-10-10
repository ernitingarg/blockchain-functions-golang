package utils

import (
	"encoding/json"
	"net/http"
)

const jsonContentType = "application/json"

// RespondJSON send Response with Json content
func RespondJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(payload)
}

// RespondJSONWithError send a response that container an error message
func RespondJSONWithError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}

// RequestData parse data from http request
func RequestData(r *http.Request) (map[string]string, error) {
	data := make(map[string]string)
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
