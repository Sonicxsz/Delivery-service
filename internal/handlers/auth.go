package handlers

import (
	"arabic/internal/dto"
	"arabic/internal/model"
	"arabic/internal/service"
	"arabic/pkg/errors"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := &model.User{}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
			// TODO LOGGER
		}

		isNotValid, err := UserValidator(user)
		if isNotValid {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err = h.authService.CreateUser(r.Context(), user)
		if err != nil {
			var serviceErr *errors.ServiceError
			if errors2.As(err, &serviceErr) {
				respondError(w, serviceErr.Code, serviceErr.Message)
			} else {
				log.Printf("Unexpected error type in CreateUser: %v", err)
				respondError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}
		respondSuccess(w, http.StatusCreated, user)
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.UserLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		user, token, err := h.authService.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			var serviceErr *errors.ServiceError
			if errors2.As(err, &serviceErr) {
				respondError(w, serviceErr.Code, serviceErr.Message)
			} else {
				log.Printf("Unexpected error type in Login: %v", err)
				respondError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}
		setAuthCookie(w, token)
		respondSuccess(w, http.StatusOK, user)
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
