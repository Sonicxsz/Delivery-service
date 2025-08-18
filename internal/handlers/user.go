package handlers

import (
	"arabic/internal/model"
	security "arabic/internal/security/auth"
	"arabic/store"
	"encoding/json"
	"net/http"
)

func CreateUser(repo *store.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initHeaders(w)

		user := &model.User{}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			OnJsonDataParseError(w)
			return
			// TODO LOGGER
		}

		hasErr, errors := UserValidator(user)
		if hasErr {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(NewErrorMessage(
				errors,
				400,
			))
			return
		}

		user, err := repo.Create(user)

		if err != nil {
			OnDbError(w)
			return
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(NewSuccessMessage(
			"User created successfully",
			200,
			user.Id,
		))
	}
}

func Login(repo *store.UserRepository, jwtConfig *security.JWTConfig) http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var data Request

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			OnJsonDataParseError(w)
			return
		}

		u, ok, err := repo.FindByEmail(data.Email)

		if err != nil {
			OnDbError(w)
			return
		}

		if !ok || data.Password != u.Password {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(NewErrorMessage(
				"Please check provided login and password",
				400,
			))
			return
		}

		token, err := security.GenerateJWT(u.Id, jwtConfig)
		if err != nil {
			w.WriteHeader(500)
			println("Token error", err.Error())
			OnDbError(w)
			return
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(NewSuccessMessage(
			token,
			200,
			u,
		))
	}
}
