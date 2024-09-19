package helpers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, v any) {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		json.NewEncoder(w).Encode(nil)
		return
	}
}

func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteJSONError(w http.ResponseWriter, code int, message string) error {
	log.Println(message)
	return WriteJSON(w, code, APIError{
		Code:    code,
		Message: message,
	})
}

type APISuccess struct {
	Data    any    `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteJSONSuccess(w http.ResponseWriter, code int, message string, data any) error {
	return WriteJSON(w, code, APISuccess{
		Data:    data,
		Code:    code,
		Message: message,
	})
}
