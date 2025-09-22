package server

import (
	"arabic/internal/server/builders"
	"arabic/internal/store"
	"arabic/pkg/fs"
	"arabic/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
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
	builders.BuildRoutesStatic(router, a.config.FS.Path)

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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:4200" || origin == "http://localhost:5173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
