package handlers

import (
	"arabic/internal/middlewares"
	"arabic/internal/model"
	"arabic/store"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

		user, err := repo.Create(user)

		if err != nil {
			OnDbError(w)
			return
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(NewSuccessMessage(
			"User created succesfully",
			200,
			user.Id,
		))
	}
}

func Login(repo *store.UserRepository, jwtConfig *middlewares.JWTConfig) http.HandlerFunc {
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

		w.WriteHeader(200)

		claims := middlewares.CustomClaims{
			UserID: u.Id,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
				Issuer:    jwtConfig.Issuer,
				Audience:  []string{jwtConfig.Audience},
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(jwtConfig.SecretJWTKey))

		if err != nil {
			println("Token error", err.Error())
			fmt.Printf("%+v\n", jwtConfig)
			OnDbError(w)
			return
		}

		json.NewEncoder(w).Encode(NewSuccessMessage(
			tokenString,
			200,
			u,
		))
	}
}
