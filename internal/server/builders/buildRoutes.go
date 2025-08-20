package builders

import (
	"arabic/internal/handlers"
	security "arabic/internal/security/auth"
	"arabic/internal/service"
	"arabic/store"
	"net/http"

	"github.com/gorilla/mux"
)

func BuildRoutes(r *mux.Router, store *store.Store, jwtConfig *security.JWTConfig) {
	authService := service.NewAuthService(store.UserRepo(), jwtConfig)
	authHandler := handlers.NewAuthHandler(authService)

	r.HandleFunc("/register", authHandler.CreateUser()).Methods("POST")
	r.HandleFunc("/login", authHandler.Login()).Methods("POST")
}

func BuildProtectedRoutes(r *mux.Router, store *store.Store, jwtConfig *security.JWTConfig) {
	JWTMiddleware := security.NewJwtMiddleware(jwtConfig)

	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(JWTMiddleware.CheckJWT)

	protected.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Token is valid"))
	}).Methods("POST")
	// Будут защищенные api роуты
}
