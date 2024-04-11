package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mBayzigitov/dynamic-content-service/internal/handler/banner"
	"github.com/mBayzigitov/dynamic-content-service/internal/repo"
	"github.com/mBayzigitov/dynamic-content-service/internal/service"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
)

type ApiServer struct {
	p      *pgxpool.Pool
	config *AppConfig
	logger *zap.SugaredLogger
	redis  *redis.Client
}

func New(config *AppConfig, logger *zap.SugaredLogger, p *pgxpool.Pool, rc *redis.Client) *ApiServer {
	return &ApiServer{
		p:      p,
		config: config,
		logger: logger,
		redis:  rc,
	}
}

func (serv *ApiServer) Start() error {
	serv.logger.Info("Starting API server")

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	subrouter.Use(service.TokenValidationMiddleware)

	cr := repo.NewCacheRepo(serv.redis)

	br := repo.NewBannerRepository(serv.p)
	bs := service.NewBannerService(br, cr)

	bh := banner.NewHandler(bs)
	bh.RegisterRoutes(subrouter)

	return http.ListenAndServe(serv.config.ServerPort, router)
}
