package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis address
		Password: "",               // Set password if needed
		DB:       0,                // Default DB
	})

	// Ping Redis
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("❌ Could not connect to Redis:", err)
		return
	}

	fmt.Println("✅ Connected to Redis:", pong)
}
