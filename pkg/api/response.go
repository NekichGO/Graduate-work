package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseError struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "	")
	if err := encoder.Encode(data); err != nil {
		log.Printf("write json: %v", err)
	}
}

func writeError(w http.ResponseWriter, statusCode int, msg string) {
	writeJSON(w, statusCode, ResponseError{Error: msg})
}
