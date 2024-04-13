package test

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mBayzigitov/dynamic-content-service/internal/apiserver"
	"github.com/mBayzigitov/dynamic-content-service/internal/handler/banner"
	"github.com/mBayzigitov/dynamic-content-service/internal/repo"
	"github.com/mBayzigitov/dynamic-content-service/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"log"
	"time"
)

type BannerHandlerSuite struct {
	suite.Suite
	router   *mux.Router
	pool     *pgxpool.Pool
	rediscli *redis.Client
}

func (suite *BannerHandlerSuite) SetupTest() {
	conf, err := apiserver.GetAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Second) // wait until environment is load completely

	dbURL := conf.Postgres.GetDbUrl()
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to a database: %v\n", err)
	}
	suite.pool = pool

	rediscli := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Url,
		Password: conf.Redis.Password,
		DB:       conf.Redis.Db,
	})
	suite.rediscli = rediscli

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	subrouter.Use(service.TokenValidationMiddleware)

	cr := repo.NewCacheRepo(rediscli)

	br := repo.NewBannerRepository(pool)
	bs := service.NewBannerService(br, cr)

	bh := banner.NewHandler(bs)
	bh.RegisterRoutes(subrouter)
	suite.router = router
}

func (suite *BannerHandlerSuite) TearDownSuite() {
	suite.pool.Close()
	suite.rediscli.Close()
}
