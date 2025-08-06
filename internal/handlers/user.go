package handlers

import (
	"arabic/internal/model"
	"arabic/store"
	"encoding/json"
	"net/http"
)

func CreateUser(repo *store.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initHeaders(w)

		user := &model.User{}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(NewErrorMessage(
				"Cannot parse data, please check provided user data",
				400,
			))

			// TODO LOGGER
		}

		user, err := repo.Create(user)

		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(NewErrorMessage(
				"Somethink went wrong while creating user, please try later...",
				500,
			))
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(NewSuccessMessage(
			"User created succesfully",
			200,
			user.Id,
		))
	}
}
