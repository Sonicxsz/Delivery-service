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

func (b *RouteBuilder) BuildRoutes(r *mux.Router) {
	authService := service.NewAuthService(b.store.UserRepo(), b.jwtConfig)

	r.HandleFunc("/register", handlers.CreateUser(authService, b.logger)).Methods("POST")
	r.HandleFunc("/login", handlers.Login(authService, b.logger)).Methods("POST")
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
