package handlers

import (
	"arabic/internal/dto"
	"arabic/pkg/customError"
	"arabic/pkg/logger"
	"arabic/pkg/validator"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type SuccessMessage[T any] struct {
	Status  int  `json:"status"`
	Data    T    `json:"data,omitempty"`
	Success bool `json:"success"`
}

type ErrorMessage[T any] struct {
	Status  int       `json:"status"`
	Error   string    `json:"error"`
	Success bool      `json:"success"`
	Id      uuid.UUID `json:"id"`
	Time    time.Time `json:"time"`
	Path    string    `json:"path"`
}

func NewSuccessMessage[T any](status int, data T) *SuccessMessage[T] {
	return &SuccessMessage[T]{
		Status:  status,
		Data:    data,
		Success: true,
	}
}

func NewErrorMessage(path, error string, status int) *ErrorMessage[interface{}] {
	return &ErrorMessage[interface{}]{
		Status:  status,
		Error:   error,
		Success: false,
		Id:      uuid.New(),
		Time:    time.Now(),
		Path:    path,
	}
}

func UserValidator(user *dto.UserCreateRequest) (bool, error) {
	v := validator.New()
	v.CheckString(user.Email, "Email").IsEmail()
	v.CheckString(user.Username, "Username").IsValidUsername()
	v.CheckString(user.Password, "Password").IsPassword()

	hasErrors, err := v.HasErrors(), v.GetErrors()

	if hasErrors {
		return hasErrors, customError.NewServiceError(http.StatusBadRequest, strings.Join(err, ", "), nil)
	}

	return hasErrors, nil
}

func handleServiceError(w http.ResponseWriter, err error, operation string) {
	var serviceErr *customError.ServiceError
	if errors.As(err, &serviceErr) {
		respondError(w, serviceErr.Code, serviceErr.Message)
	} else {
		logger.Log.Error("Unexpected error type in %s: %v", operation, err)
		respondError(w, http.StatusInternalServerError, customError.Error500)
	}
}

func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})
}

func respondSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(NewSuccessMessage(status, data))
}

func respondError(w http.ResponseWriter, status int, error string) {
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
	json.NewEncoder(w).Encode(NewErrorMessage(pathInfo, error, status))
}
