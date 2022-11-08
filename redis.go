package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RDB struct {
	client *redis.Client
}

func newRDB() *RDB {
	return &RDB{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:16379",
			Password: "",
			DB:       0,
		}),
	}
}

func (r *RDB) call(ctx context.Context, args ...interface{}) (interface{}, error) {
	return r.client.Do(ctx, args...).Result()
}
