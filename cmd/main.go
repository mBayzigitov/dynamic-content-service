package main

import (
	"context"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mBayzigitov/dynamic-content-service/internal/apiserver"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"time"
)

var (
	confPath string
)

func init() {
	flag.StringVar(&confPath, "config-path", "config/apiserver.toml", "path to a config file")
}

func main() {
	flag.Parse()

	config, err := apiserver.GetAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	_, err = toml.DecodeFile(confPath, config)
	if err != nil {
		log.Fatal(err)
	}

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	time.Sleep(10 * time.Second) // wait until environment is load completely

	dbURL := config.Postgres.GetDbUrl()
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to a database: %v\n", err)
	}
	defer pool.Close()

	rediscli := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Url,
		Password: config.Redis.Password,
		DB:       config.Redis.Db,
	})
	defer rediscli.Close()

	serv := apiserver.New(
		config,
		logger.Sugar(),
		pool,
		rediscli,
	)
	if err := serv.Start(); err != nil {
		log.Fatal(err)
	}
}
