package main

import (
	"net/http"
	"os"

	"github.com/go-redis/redis"
)

var (
	API_PORT        string = ":" + getenv("API_PORT", "80")
	REVERSE_API_URI string = getenv("REVERSE_API_URI", "http://127.0.0.1:8001/reverse")
	REDIS_HOST      string = getenv("REDIS_HOST", "redishost:6379")
	redisClient     *redis.Client
)

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: REDIS_HOST,
		DB:   0,
	})
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	http.HandleFunc("/api", HandleRandom)
	http.ListenAndServe(API_PORT, nil)
}
