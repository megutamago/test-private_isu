package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // Password (if any)
		DB:       0,                // Default database
	})

	// Ping the Redis server to check if it's reachable
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)

	// Set a value in the Redis cache
	err = client.Set(context.Background(), "mykey", "myvalue", 10*time.Second).Err()
	if err != nil {
		fmt.Println("Failed to set value in Redis:", err)
		return
	}

	// Get a value from the Redis cache
	value, err := client.Get(context.Background(), "mykey").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("Key not found in Redis cache")
		} else {
			fmt.Println("Failed to get value from Redis:", err)
		}
		return
	}
	fmt.Println("Value from Redis:", value)
}

