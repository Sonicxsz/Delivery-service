package handlers

import (
	"arabic/internal/handlers/dto"
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

		userId, err := repo.Create(user)

		if err != nil {
			handleDatabaseError(w, err)
			return
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(NewSuccessMessage(
			"User created successfully",
			200,
			userId,
		))
	}
}

func Login(repo *store.UserRepository, jwtConfig *security.JWTConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var data dto.UserRequest

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			OnJsonDataParseError(w)
			return
		}
		u := verifyCredentials(w, repo, &data)
		generateToken(u.Id, jwtConfig, w)

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(NewSuccessMessage(
			"Your advertisement could be here.",
			200,
			u,
		))
	}
}
