package db

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"

	"gharabakloo/search/internal/entity"
	"gharabakloo/search/pkg/myerr"
)

const coma = ","

var errNotFoundRedisHost = errors.New("not found redis host")

func Redis(ctx context.Context, cfg entity.RedisConfig) (entity.RedisClient, error) {
	db, err := strconv.Atoi(cfg.DB)
	if err != nil {
		return nil, myerr.Errorf(err)
	}

	hosts := strings.Split(cfg.Host, coma)
	if len(hosts) == 0 {
		return nil, myerr.Errorf(errNotFoundRedisHost)
	}

	var client entity.RedisClient
	if len(hosts) == 1 {
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Host,
			Password: cfg.Pass,
			DB:       db,
		})
	} else {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:          hosts,
			Password:       cfg.Pass,
			RouteByLatency: true,
		})
	}

	if err = client.Ping(ctx).Err(); err != nil {
		return nil, myerr.Errorf(err)
	}
	return client, nil
}
