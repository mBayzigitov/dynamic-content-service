package config

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"os"
)

type (
	AppConfig struct {
		ServerPort string `env:"SERVER_PORT"`
		Postgres   *Postgres
	}

	Postgres struct {
		Host     string `env:"PG_HOST"`
		Port     string `env:"PG_PORT"`
		DbName   string `env:"PG_DBNAME"`
		User     string `env:"PG_USER"`
		Password string `env:"PG_PASSWORD"`
	}
)

// Load each .env file from config/environ
func loadEnv() error {
	currentWorkingDirectory, _ := os.Getwd()
	envDir := currentWorkingDirectory + "/config/environ/"

	envFiles, _ := os.ReadDir(envDir)

	for _, env := range envFiles {
		fmt.Println(envDir + env.Name())
		err := godotenv.Load(envDir + env.Name())

		if err != nil {
			return err
		}
	}

	return nil
}

func GetAppConfig() (*AppConfig, error) {
	err := loadEnv()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	var config AppConfig
	err = envconfig.Process(ctx, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
