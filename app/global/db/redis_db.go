package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type RedisRawCredential struct {
	Host     string
	Port     string
	Password string
	Database int
}

type RedisParam struct {
	RedisOptions *redis.Options
}

func (r *RedisRawCredential) NewRedisParam() RedisParam {
	return RedisParam{
		RedisOptions: &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
			Password: r.Password, // no password set
			DB:       r.Database, // use default DB
		},
	}
}

func RedisGetEnvVariable() *RedisRawCredential {
	database, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		logrus.Fatal("fail to convert redis db to int, err:", err)
	}
	return &RedisRawCredential{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Database: database,
	}
}
func NewRedisClient(redisParam RedisParam) *redis.Client {
	rdb := redis.NewClient(redisParam.RedisOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := rdb.Ping(ctx).Err()
	if err != nil {
		logrus.Fatal("fail to ping redis server")
	}
	return rdb
}
