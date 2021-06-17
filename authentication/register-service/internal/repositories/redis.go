package repositories

import "github.com/go-redis/redis/v8"

func NewRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:               "",
		Username:           "",
		Password:           "",
		DB:                 0,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		PoolSize:           0,
		MinIdleConns:       0,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		Limiter:            nil,
	})
}
