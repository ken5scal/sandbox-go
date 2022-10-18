package store

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ken5scal/go_todo_app/config"
	"github.com/ken5scal/go_todo_app/entity"
	"time"
)

type KVS struct {
	Client *redis.Client
}

func NewKVS(ctx context.Context, cfg *config.Config) (*KVS, error) {
	c := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)})
	if err := c.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &KVS{c}, nil
}

func (k *KVS) Save(ctx context.Context, key string, userID entity.UserID) error {
	return k.Client.Set(ctx, key, int64(userID), 30*time.Minute).Err()
}

func (k *KVS) Load(ctx context.Context, key string) (entity.UserID, error) {
	id, err := k.Client.Get(ctx, key).Int64()
	if err != nil {
		return 0, fmt.Errorf("failed to get by %q: %w", key, ErrNotFound)
	}
	return entity.UserID(id), err
}
