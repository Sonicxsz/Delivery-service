package handlers

import (
	"arabic/internal/model"
	"arabic/pkg/errors"
	"arabic/pkg/validator"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	errors2 "github.com/pkg/errors"
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
		return hasErrors, errors.NewServiceError(http.StatusBadRequest, strings.Join(err, ", "), nil)
	}

	return hasErrors, nil
}

func handleServiceError(w http.ResponseWriter, err error, operation string) {

	var serviceErr *errors.ServiceError
	if errors2.As(err, &serviceErr) {
		respondError(w, serviceErr.Code, serviceErr.Message)
	} else {
		log.Printf("Unexpected error type in %s: %v", operation, err)
		respondError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		// нужно изучить другие параметры
	})
}

func respondSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  status,
		"data":    data,
		"success": true, // нужен ли такой параметр если все и так понятно ?
	})
}

func respondError(w http.ResponseWriter, status int, message string) {
	pc, file, line, ok := runtime.Caller(1)
	var pathInfo string
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		pathInfo = fmt.Sprintf("%s:%d (%s)", file, line, funcName)
	} else {
		pathInfo = "unknown"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  status,
		"errors":  message,
		"success": false,
		"time":    time.Now(),
		"id":      uuid.New(),
		"path":    pathInfo,
	})
}
