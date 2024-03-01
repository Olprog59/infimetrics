package database

import (
	"fmt"
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/commons"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisDB struct {
	Client *redis.Client
}

func NewRedis() *RedisDB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", commons.REDIS_HOST, commons.REDIS_PORT),
		Password: commons.REDIS_PASSWORD,
		DB:       0,
	})
	return &RedisDB{
		Client: rdb,
	}
}

func (r *RedisDB) Ping() error {
	_, err := r.Client.Ping(r.Client.Context()).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisDB) Set(key string, value any) error {
	err := r.Client.Set(r.Client.Context(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisDB) SetWithTimeout(key string, value any, timeout time.Duration) error {
	err := r.Client.Set(r.Client.Context(), key, value, timeout).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisDB) Get(key string) (string, error) {
	val, err := r.Client.Get(r.Client.Context(), key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisDB) Del(key string) error {
	err := r.Client.Del(r.Client.Context(), key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisDB) Close() {
	err := r.Client.Close()
	if err != nil {
		golog.Err("Error closing redis connection: %s", err.Error())
		return
	}
}
