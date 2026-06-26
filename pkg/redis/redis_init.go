package redis

import (
	"fish_shooting_admin_backend/configs"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	once   sync.Once
	client *redis.Client
)

func NewRedis() *redis.Client {
	redis_config := configs.Redis()

	once.Do(func() {

		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redis_config.RedisHost, redis_config.RedisPort),
			Password: redis_config.RedisPassword,
			DB:       redis_config.RedisDB,
		})
	
		pong, err := client.Ping(context.Background()).Result()
		if err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}
		fmt.Printf("Connection to redis successful: %s\n", pong)
	})

	return client
}
