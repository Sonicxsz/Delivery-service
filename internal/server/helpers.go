package server

import (
	"arabic/internal/handlers"
	"arabic/internal/middlewares"
	"arabic/store"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Извлекаем claims (используем нашу функцию)
	_, err := middlewares.GetClaimsFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Hello"))
}

func (a *Api) configureRouter() {
	router := mux.NewRouter()
	JWTMiddleware := middlewares.NewJwtMiddleware(a.config.JWT)

	router.HandleFunc("/register", handlers.CreateUser(a.store.User())).Methods("POST")
	router.HandleFunc("/login", handlers.Login(a.store.User(), a.config.JWT)).Methods("POST")

	protected := router.PathPrefix("/api/v1").Subrouter()
	protected.Use(JWTMiddleware.CheckJWT)

	protected.HandleFunc("/getHello", testHandler).Methods("POST")

	a.router = router
}

func (a *Api) configureLogger() error {
	level, err := logrus.ParseLevel(a.config.LogLevel)

	if err != nil {
		return err
	}

	a.Logger.SetLevel(level)

	return nil
}

func (a *Api) configureStore() error {
	store := store.New(a.config.Storage)

	if err := store.Start(); err != nil {
		return err
	}
	a.store = store
	return nil
}
