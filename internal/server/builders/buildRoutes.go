package builders

import (
	"arabic/internal/handlers"
	security "arabic/internal/security/auth"
	"arabic/internal/service"
	"arabic/store"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type RouteBuilder struct {
	store     *store.Store
	jwtConfig *security.JWTConfig
	logger    *logrus.Logger
}

func NewRouteBuilder(store *store.Store, jwtConfig *security.JWTConfig, logger *logrus.Logger) *RouteBuilder {
	return &RouteBuilder{
		store:     store,
		jwtConfig: jwtConfig,
		logger:    logger,
	}
}

func BuildRoutes(r *mux.Router, store *store.Store, jwtConfig *security.JWTConfig) {
	authService := service.NewAuthService(store.UserRepo(), jwtConfig)

	r.HandleFunc("/register", handlers.CreateUser(authService)).Methods("POST")
	r.HandleFunc("/login", handlers.Login(authService)).Methods("POST")
}

func (b *RouteBuilder) BuildRoutes(r *mux.Router) {
	r.HandleFunc("/register", handlers.CreateUser(b.store.User(), b.logger)).Methods("POST")
	r.HandleFunc("/login", handlers.Login(b.store.User(), b.jwtConfig, b.logger)).Methods("POST")
}

func (b *RouteBuilder) BuildProtectedRoutes(r *mux.Router) {
	JWTMiddleware := security.NewJwtMiddleware(b.jwtConfig)

	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(JWTMiddleware.CheckJWT)

	protected.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Token is valid"))
	}).Methods("POST")
	// Будут защищенные api роуты
}
