package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"linkshare/app/global/helper"
	"os"
	"strconv"
	"time"
)

type IRedis interface {
	Client() *redis.Client
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	LPush(ctx context.Context, key string, value interface{}) error
	LPop(ctx context.Context, key string, dest interface{}) error
	RPush(ctx context.Context, key string, value interface{}) error
	RPop(ctx context.Context, key string, dest interface{}) error
	Llen(ctx context.Context, key string) (int64, error)
	LMove(ctx context.Context, source, dest, srcpos, destpos string) error
	LTrim(ctx context.Context, key string, start, stop int64) error
	Del(ctx context.Context, key string) error
	Ping(ctx context.Context) error
}

type Redis struct {
	redis *redis.Client
}

type RedisParam struct {
	Host     string
	Port     int
	Password string
	Database int
}

func RedisGetEnvVariable() *RedisParam {
	database, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		logrus.Fatal("fail to convert redis db to int, err:", err)
	}
	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		logrus.Fatal("fail to convert redis port to int, err:", err)
	}
	return &RedisParam{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     port,
		Password: os.Getenv("REDIS_PASSWORD"),
		Database: database,
	}
}

func (r *RedisParam) NewRedisParam() RedisParam {
	return *r
}
func NewRedisClient(param RedisParam) IRedis {
	ctx := context.Background()

	var rdb *redis.Client

	if len(param.Password) > 0 {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", param.Host, param.Port),
			DB:       param.Database,
			Password: param.Password,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", param.Host, param.Port),
			DB:   param.Database,
		})
	}

	status := rdb.Ping(ctx)
	if status.Err() != nil {
		logrus.Fatal("fail to connect to redis, err:", status.Err())
	}

	return &Redis{
		redis: rdb,
	}
}

func (rdb *Redis) Client() *redis.Client {
	return rdb.redis
}

func (rdb *Redis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"marshal_data": fmt.Sprintf("%+v", value),
			"call_from":    helper.GetCaller(2),
		}).Errorf("error marshaling while set redis, err: %v ", err)
		return err
	}
	compressedData, err := helper.GzipCompress(val)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"marshal_data": fmt.Sprintf("%+v", value),
			"call_from":    helper.GetCaller(2),
		}).Errorf("error compressing while set redis: %v", err)
		return err
	}
	return rdb.redis.Set(ctx, key, string(compressedData), ttl).Err()
}

func (rdb *Redis) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := rdb.redis.Get(ctx, key).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"redis_key": key,
			"call_from": helper.GetCaller(2),
		}).Tracef("error while get redis: %v", err)
		return err
	}
	decompressData, err := helper.GzipDecompress([]byte(val))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"redis_key": key,
			"call_from": helper.GetCaller(2),
		}).Tracef("error decompressing while get redis: %v", err)
		return err
	}

	err = json.Unmarshal(decompressData, &dest)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"redis_key": key,
			"call_from": helper.GetCaller(2),
		}).Tracef("error unmarshaling while get redis: %v", err)
		return err
	}
	return nil
}

func (rdb *Redis) LPush(ctx context.Context, key string, value interface{}) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = rdb.redis.LPush(ctx, key, string(val)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rdb *Redis) LPop(ctx context.Context, key string, dest interface{}) error {
	val, err := rdb.redis.LPop(ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), &dest)
	if err != nil {
		return err
	}
	return nil
}

func (rdb *Redis) RPush(ctx context.Context, key string, value interface{}) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = rdb.redis.RPush(ctx, key, string(val)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rdb *Redis) RPop(ctx context.Context, key string, dest interface{}) error {
	val, err := rdb.redis.RPop(ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), &dest)
	if err != nil {
		return err
	}
	return nil
}

func (rdb *Redis) Llen(ctx context.Context, key string) (int64, error) {
	llen, err := rdb.redis.LLen(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return llen, nil
}

func (rdb *Redis) LMove(ctx context.Context, source, dest, srcpos, destpos string) error {
	err := rdb.redis.LMove(ctx, source, dest, srcpos, destpos).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rdb *Redis) LTrim(ctx context.Context, key string, start, stop int64) error {
	err := rdb.redis.LTrim(ctx, key, start, stop).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rdb *Redis) Del(ctx context.Context, key string) error {
	err := rdb.redis.Del(ctx, key).Err()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"redis_key": key,
			"call_from": helper.GetCaller(2),
		}).Errorf("error while delete redis: %v", err)
		return err
	}
	return nil
}

func (rdb *Redis) Ping(ctx context.Context) error {
	status := rdb.redis.Ping(ctx)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}
