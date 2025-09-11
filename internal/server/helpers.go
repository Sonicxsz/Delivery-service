package server

import (
	"arabic/internal/server/builders"
	"arabic/internal/store"
	"arabic/pkg/fs"
	"arabic/pkg/logger"
	"github.com/gorilla/mux"
)

func (a *Api) configureRouter() {
	router := mux.NewRouter()
	builder := &builders.Builder{
		Router:    router,
		Store:     a.store,
		JwtConfig: a.config.JWT,
		Fs:        a.fs,
	}

	builders.BuildRoutes(builder)
	builders.BuildProtectedRoutes(router, a.store, a.config.JWT)
	builders.BuildRoutesStatic(router, a.config.FS.Path)

	//a.config.FS.Image.Path
	a.router = router
}

func (a *Api) configureFileSystem() {
	a.fs = fs.New(a.config.FS)
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
