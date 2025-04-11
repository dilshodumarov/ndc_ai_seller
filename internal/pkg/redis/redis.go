package redis

import (
	"sugurta/internal/pkg/config"
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	Client redis.Client
}

// NewInMemoryStorage is redis client
func New(cfg *config.Config) (*RedisDB, error) {
	db, err := strconv.Atoi(cfg.Redis.Name)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       db,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		fmt.Println("Error connecting to redis: ", err)
		log.Fatal(err)
		return nil, err
	}

	return &RedisDB{
		Client: *rdb,
	}, nil
}
