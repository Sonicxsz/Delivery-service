package server

import (
	"arabic/store"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Api struct {
	config *Config
	router *mux.Router
	store  *store.Store
	Logger *logrus.Logger
}

func New(config *Config) *Api {
	return &Api{
		config: config,
		Logger: logrus.New(),
	}
}

func (api *Api) Start() error {

	if err := api.configureLogger(); err != nil {
		return err
	}

	if err := api.configureStore(); err != nil {
		return err
	}

	api.configureRouter()
	return http.ListenAndServe(api.config.BindAddr, api.router)
}
