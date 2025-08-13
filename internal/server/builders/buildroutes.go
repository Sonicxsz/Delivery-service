package builders

import (
	"arabic/internal/handlers"
	"arabic/internal/middlewares"
	"arabic/store"

	"github.com/gorilla/mux"
)

func BuildRoutes(r *mux.Router, store *store.Store, jwtConfig *middlewares.JWTConfig) {
	r.HandleFunc("/register", handlers.CreateUser(store.User())).Methods("POST")
	r.HandleFunc("/login", handlers.Login(store.User(), jwtConfig)).Methods("POST")
}

func BuildProtectedRoutes(r *mux.Router, store *store.Store, jwtConfig *middlewares.JWTConfig) {
	JWTMiddleware := middlewares.NewJwtMiddleware(jwtConfig)

	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(JWTMiddleware.CheckJWT)

	// Будут защищенные api роуты
}
