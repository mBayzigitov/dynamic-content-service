package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mBayzigitov/dynamic-content-service/internal/handler/banner"
	"github.com/mBayzigitov/dynamic-content-service/internal/repo"
	"github.com/mBayzigitov/dynamic-content-service/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type ApiServer struct {
	p *pgxpool.Pool
	config *AppConfig
	logger *zap.SugaredLogger
}

func New(config *AppConfig, logger *zap.SugaredLogger, p *pgxpool.Pool) *ApiServer {
	return &ApiServer{
		p: p,
		config: config,
		logger: logger,
	}
}

func (serv *ApiServer) Start() error {
	serv.logger.Info("Starting API server")

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	br := repo.NewBannerRepository(serv.p)
	bs := service.NewBannerService(br)

	bh := banner.NewHandler(bs)
	bh.RegisterRoutes(subrouter)

	return http.ListenAndServe(serv.config.ServerPort, router)
}