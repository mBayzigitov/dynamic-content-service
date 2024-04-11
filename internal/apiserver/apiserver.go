package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/mBayzigitov/dynamic-content-service/internal/handler/banner"
	"go.uber.org/zap"
	"net/http"
)

type ApiServer struct {
	config *AppConfig
	logger *zap.SugaredLogger
}

func New(config *AppConfig, logger *zap.SugaredLogger) *ApiServer {
	return &ApiServer{
		config: config,
		logger: logger,
	}
}

func (serv *ApiServer) Start() error {
	serv.logger.Info("Starting API server")

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	bh := banner.NewHandler()
	bh.RegisterRoutes(subrouter)

	return http.ListenAndServe(serv.config.ServerPort, router)
}