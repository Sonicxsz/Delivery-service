package handlers

import (
	"arabic/internal/model"
	"arabic/pkg/validator"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type Message[T any] struct {
	Message string    `json:"message"`
	Error   bool      `json:"error"`
	Code    int       `json:"code"`
	Data    T         `json:"data,omitempty"`
	Id      uuid.UUID `json:"id"`
	Time    time.Time `json:"time"`
}

func NewSuccessMessage[T any](msg string, code int, data T) *Message[T] {
	return &Message[T]{
		Message: msg,
		Code:    code,
		Data:    data,
		Error:   false,
	}
}

func NewErrorMessage(msg string, code int) *Message[interface{}] {
	return &Message[interface{}]{
		Message: msg,
		Code:    code,
		Error:   true,
		Id:      uuid.New(),
		Time:    time.Now(),
	}
}

func initHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json")
}

func OnJsonDataParseError(w http.ResponseWriter) {
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(NewErrorMessage(
		"Cannot parse data, please check provided user data",
		400,
	))
}

func OnDbError(w http.ResponseWriter) {
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(NewErrorMessage(
		"Something went wrong, please try later...",
		500,
	))
}

func UserValidator(user *model.User) (bool, string) {
	v := validator.New()
	v.CheckEmail(user.Email)
	v.CheckUsername(user.Username)
	v.CheckPassword(user.Password)

	hasErrors, err := v.HasErrors()

	if hasErrors {
		return hasErrors, strings.Join(err, ", ")
	}

	return hasErrors, ""
}
