package builders

import (
	"arabic/internal/handlers"
	security "arabic/internal/security/auth"
	"arabic/store"
	"net/http"

	"github.com/gorilla/mux"
)

func BuildRoutes(r *mux.Router, store *store.Store, jwtConfig *security.JWTConfig) {
	r.HandleFunc("/register", handlers.CreateUser(store.User())).Methods("POST")
	r.HandleFunc("/login", handlers.Login(store.User(), jwtConfig)).Methods("POST")
}

func BuildProtectedRoutes(r *mux.Router, store *store.Store, jwtConfig *security.JWTConfig) {
	JWTMiddleware := security.NewJwtMiddleware(jwtConfig)

	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(JWTMiddleware.CheckJWT)

	protected.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hgee"))
	}).Methods("POST")
	// Будут защищенные api роуты
}
