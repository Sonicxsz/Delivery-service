package server

import (
	"arabic/internal/store"
	"arabic/pkg/fs"
	"arabic/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

type Api struct {
	config *Config
	router *mux.Router
	store  *store.Store
	fs     *fs.FS
}

func New(config *Config) *Api {
	return &Api{
		config: config,
	}
}

func (api *Api) Start() error {
	if err := api.configureLogger(); err != nil {
		return err
	}
	defer logger.Log.Close()

	if err := api.configureStore(); err != nil {
		return err
	}

	api.configureFileSystem()

	api.configureRouter()

	return http.ListenAndServe(api.config.BindAddr, corsMiddleware(api.router))
}
