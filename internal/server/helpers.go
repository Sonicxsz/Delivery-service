package server

import (
	"arabic/internal/server/builders"
	"arabic/pkg/logger"
	"arabic/store"

	"github.com/gorilla/mux"
)

func (a *Api) configureRouter() {
	router := mux.NewRouter()

	builders.BuildRoutes(router, a.store, a.config.JWT)
	builders.BuildProtectedRoutes(router, a.store, a.config.JWT)

	a.router = router
}

func (a *Api) configureLogger() error {
	return logger.Init(a.config.LogLevel, a.config.LogDir)
}

func (a *Api) configureStore() error {
	store := store.New(a.config.Storage)

	if err := store.Start(); err != nil {
		return err
	}
	a.store = store
	return nil
}
