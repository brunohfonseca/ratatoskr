package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func ConnectRedis(RedisURL string) {
	// Parse the Redis URL to get connection options
	opts, err := redis.ParseURL(RedisURL)
	if err != nil {
		log.Fatal("Error parsing Redis URL:", err)
	}

	// Create client with parsed options
	client := redis.NewClient(opts)

	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("Connected to Redis")
}
