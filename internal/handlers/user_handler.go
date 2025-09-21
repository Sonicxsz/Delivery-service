package handlers

import (
	"arabic/internal/dto"
	"arabic/internal/model"
	"arabic/internal/service"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	service service.IUserService
}

func NewUserHandler(service service.IUserService) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) Create() http.HandlerFunc {
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

		user, err = u.service.CreateUser(r.Context(), user)
		if err != nil {
			handleServiceError(w, err, "CreateUser")
			return
		}
		respondSuccess(w, http.StatusCreated, user)
	}
}

func (u *UserHandler) Login() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.UserLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		user, token, err := u.service.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			handleServiceError(w, err, "Login")
			return
		}
		setAuthCookie(w, token)
		respondSuccess(w, http.StatusOK, user)
	}
}
