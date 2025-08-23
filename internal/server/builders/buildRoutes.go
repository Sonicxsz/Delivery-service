package builders

import (
	"arabic/internal/handlers"
	security "arabic/internal/security/auth"
	"arabic/internal/service"
	"arabic/internal/store"
	"net/http"

	"github.com/gorilla/mux"
)

func BuildRoutes(r *mux.Router, store *store.Store, jwtConfig *security.JWTConfig) {
	authService := service.NewAuthService(store.UserRepository(), jwtConfig)
	tagService := service.NewTagService(store.TagRepository())

	r.HandleFunc("/register", handlers.CreateUser(authService)).Methods("POST")
	r.HandleFunc("/login", handlers.Login(authService)).Methods("POST")
	r.HandleFunc("/tag", handlers.CreateTag(tagService)).Methods("POST")
	r.HandleFunc("/all-tags", handlers.FindAllTags(tagService)).Methods("GET")
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
