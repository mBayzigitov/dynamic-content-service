package main

import (
	"context"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mBayzigitov/dynamic-content-service/internal/apiserver"
	"go.uber.org/zap"
	"log"
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

	dbURL := config.Postgres.GetDbUrl()
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to a database: %v\n", err)
	}
	defer pool.Close()

	serv := apiserver.New(
		config,
		logger.Sugar(),
	)
	if err := serv.Start(); err != nil {
		log.Fatal(err)
	}
}
