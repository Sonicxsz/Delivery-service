package handlers

import (
	"arabic/internal/model"
	"arabic/pkg/validator"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Message[T any] struct {
	Message string    `json:"message"`
	Error   bool      `json:"errors"`
	Code    int       `json:"code"`
	Data    T         `json:"data,omitempty"`
	Id      uuid.UUID `json:"id,omitempty"`
	Time    time.Time `json:"time,omitempty"`
	Path    string    `json:"path,omitempty"`
}

func NewSuccessMessage[T any](msg string, code int, data T) *Message[T] {
	return &Message[T]{
		Message: msg,
		Code:    code,
		Data:    data,
		Error:   false,
		Id:      uuid.New(),
		Time:    time.Now(),
		Path:    "",
	}
}

func NewErrorMessage(msg string, code int) *Message[interface{}] {
	return &Message[interface{}]{
		Message: msg,
		Code:    code,
		Error:   true, // лишнее поле
		Id:      uuid.New(),
		Time:    time.Now(),
		Path:    "",
	}
}

func UserValidator(user *model.User) (bool, error) {
	v := validator.New()
	v.CheckEmail(user.Email)
	v.CheckUsername(user.Username)
	v.CheckPassword(user.Password)

	hasErrors, err := v.HasErrors()

	if hasErrors {
		return hasErrors, errors.New(strings.Join(err, ", "))
	}

	return hasErrors, nil
}
