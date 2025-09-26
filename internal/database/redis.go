package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis(RedisURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

func DisconnectRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Printf("Erro ao desconectar do Redis: %v", err)
		} else {
			fmt.Println("✅ Disconnected from Redis")
		}
	}
}
