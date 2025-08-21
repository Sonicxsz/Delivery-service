package server

import (
	"arabic/internal/server/builders"
	"arabic/store"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (a *Api) configureRouter() {
	router := mux.NewRouter()

	builder := builders.NewRouteBuilder(a.store, a.config.JWT, a.Logger)
	builder.BuildRoutes(router)
	builder.BuildProtectedRoutes(router)

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
