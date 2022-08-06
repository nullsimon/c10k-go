package redis

import (
	"context"
	"github.com/go-redis/redis/v9"
)

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func InitQuantity(ctx context.Context, r *redis.Client, key string, quantity int) {
	err := r.Set(ctx, key, quantity, 0).Err()
	if err != nil {
		panic(err)
	}
}

func DecreaseQuantity(ctx context.Context, r *redis.Client, key string, quantity int64) bool {
	err := r.DecrBy(ctx, key, quantity).Err()
	if err != nil {
		panic(err)
		return false
	}
	return true
}
