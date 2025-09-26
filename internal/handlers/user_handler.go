package handlers

import (
	"arabic/internal/dto"
	"arabic/internal/service"
	"arabic/pkg/customError"
	"arabic/pkg/logger"
	security "arabic/pkg/security/auth"
	"encoding/json"
	"net/http"
	"strings"
)

type UserHandler struct {
	service service.IUserService
}

func NewUserHandler(service service.IUserService) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := &dto.UserCreateRequest{}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		isNotValid, err := UserValidator(user)
		if isNotValid {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = u.service.CreateUser(r.Context(), user)

		if err != nil {
			handleServiceError(w, err, "CreateUser")
			return
		}
		respondSuccess(w, http.StatusCreated, nil)
	}
}

func (u *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims, err := security.GetClaimsFromContext(r)

	if err != nil || claims.UserEmail == "" {
		logger.Log.Error("UserHandler -> Get -> err: " + err.Error())
		handleServiceError(w, customError.NewServiceError(http.StatusBadRequest, customError.ErrorAuthorize, nil), "User: Get")
		return
	}

	user, err := u.service.GetUser(r.Context(), claims.UserEmail)

	if err != nil {
		handleServiceError(w, err, "User: Get")
		return
	}

	respondSuccess(w, http.StatusOK, user)
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

func (u *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, err := security.GetClaimsFromContext(r)

	req := dto.UserUpdateRequest{}
	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		handleServiceError(w, customError.NewServiceError(http.StatusBadRequest, customError.ErrorParse, nil), "User: Decode error")
		return
	}

	if ok, errStrings := req.IsValid(); !ok {
		handleServiceError(w, customError.NewServiceError(http.StatusBadRequest, strings.Join(errStrings, "; "), nil), "User: Validation error")
		return
	}

	req.Id = claims.Id
	err = u.service.UpdateUserInfo(r.Context(), &req)

	if err != nil {
		handleServiceError(w, err, "User: Service error")
		return
	}

	respondSuccess(w, http.StatusOK, nil)
}

func (u *UserHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	claims, err := security.GetClaimsFromContext(r)

	req := dto.UserAddressUpdateRequest{}
	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		handleServiceError(w, customError.NewServiceError(http.StatusBadRequest, customError.ErrorParse, nil), "User: Decode error")
		return
	}

	if ok, errStrings := req.IsValid(); !ok {
		handleServiceError(w, customError.NewServiceError(http.StatusBadRequest, strings.Join(errStrings, "; "), nil), "User: Validation error")
		return
	}

	req.Id = claims.Id
	err = u.service.UpdateUserAddress(r.Context(), &req)

	if err != nil {
		handleServiceError(w, err, "User: Service error")
		return
	}

	respondSuccess(w, http.StatusOK, nil)
}
