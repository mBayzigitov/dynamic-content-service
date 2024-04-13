package repo

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheRepo struct {
	redcli *redis.Client
	c      context.Context
}

func NewCacheRepo(client *redis.Client) *CacheRepo {
	return &CacheRepo{
		redcli: client,
		c:      context.Background(),
	}
}

func (cr *CacheRepo) Set(key string, content string, ttl time.Duration) error {
	if err := cr.redcli.Set(cr.c, key, content, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (cr *CacheRepo) Get(key string) (string, error) {
	var content string
	if err := cr.redcli.Get(cr.c, key).Scan(&content); err != nil {
		if errors.Is(err, redis.Nil) {
			err = errors.New("redis: error while getting content")
		}

		return "", err
	}

	return content, nil
}
