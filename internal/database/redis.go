package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis(RedisURL string) {
	ctx := context.Background()
	opts, err := redis.ParseURL(RedisURL)
	if err != nil {
		log.Fatal("Error parsing Redis URL:", err)
	}

	client := redis.NewClient(opts)

	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("Connected to Redis")
	RedisClient = client
}
