package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	/*
		note:-
		In this code, &envelope{Error: message} creates a pointer to an envelope struct,
		and writeJSON accepts an any type(data any), which can be any type,
		including a pointer. The json.NewEncoder(w).Encode(data) function can handle both pointers and values,
		so it works without any issues
	*/

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)

}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	maxBytes := 1_048_573
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, &envelope{Error: message})
}
