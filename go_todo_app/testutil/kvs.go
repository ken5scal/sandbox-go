package testutil

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"testing"
)

func OpenRedisForTest(t *testing.T) *redis.Client {
	t.Helper()
	port := 36379
	host := "127.0.0.1"
	if _, defined := os.LookupEnv("CI"); defined {
		port = 6379 // CI (Github Action)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("failed to connect redis: %s", err)
	}
	return client
}
