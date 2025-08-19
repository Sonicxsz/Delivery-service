package handlers

import (
	"arabic/internal/handlers/dto"
	"arabic/internal/model"
	security "arabic/internal/security/auth"
	"arabic/pkg/validator"
	"arabic/store"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Message[T any] struct {
	Message string    `json:"message"`
	Error   bool      `json:"error"`
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

func initHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json")
}

func OnJsonDataParseError(w http.ResponseWriter) {
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(NewErrorMessage(
		"Cannot parse data, please check provided user data",
		400,
	))
}

func handleDatabaseError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	message := err.Error()
	if strings.Contains(err.Error(), "duplicate key") {
		message = getDuplicateField(err)
	}
	json.NewEncoder(w).Encode(NewErrorMessage(
		message,
		500,
	))
}

func getDuplicateField(err error) string {

	if strings.Contains(err.Error(), "email") {
		return "user with email %s already exists: %w"
	} else if strings.Contains(err.Error(), "username") {
		return "user with username %s already exists: %w"
	} else {
		return "user with id %s already exists: %w"
	}
}

func UserValidator(user *model.User) (bool, string) {
	v := validator.New()
	v.CheckEmail(user.Email)
	v.CheckUsername(user.Username)
	v.CheckPassword(user.Password)

	hasErrors, err := v.HasErrors()

	if hasErrors {
		return hasErrors, strings.Join(err, ", ")
	}

	return hasErrors, ""
}

func verifyCredentials(w http.ResponseWriter, repo *store.UserRepository, req *dto.UserRequest) *model.User {
	u, ok, err := repo.FindByEmail(req.Email)

	if err != nil {
		handleDatabaseError(w, err)
		return nil
	}

	// можно добавить защиту от timing attacks
	if !ok || bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(NewErrorMessage(
			"Please check provided login and password",
			400,
		))
		return nil
	}
	return u
}

func generateToken(userId string, jwtConfig *security.JWTConfig, w http.ResponseWriter) {
	token, err := security.GenerateJWT(userId, jwtConfig)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Token generation error %v", err.Error())
		handleDatabaseError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	})
	log.Println("Token generated successfully")
}
