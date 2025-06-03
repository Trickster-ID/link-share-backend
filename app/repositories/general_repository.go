package repositories

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"linkshare/app/global/db"
	"time"
)

type IGeneralRepository interface {
	BeginTransaction(ctx context.Context) pgx.Tx
	SetRedisCache(ctx context.Context, key string, value interface{}, expiration time.Duration)
	GetRedisCache(ctx context.Context, key string, destination interface{}) error
}

type generalRepository struct {
	pgDb    db.PgxIface
	redisDb db.IRedis
}

func NewGeneralRepository(sql db.PgxIface, redisDb db.IRedis) IGeneralRepository {
	return &generalRepository{
		pgDb:    sql,
		redisDb: redisDb,
	}
}

func (r *generalRepository) BeginTransaction(ctx context.Context) pgx.Tx {
	tx, err := r.pgDb.Begin(ctx)
	if err != nil {
		return nil
	}
	return tx
}

func (r *generalRepository) SetRedisCache(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	err := r.redisDb.Set(ctx, key, value, expiration)
	if err != nil {
		logrus.Errorf("error while setting redis cache: %v", err)
	}
}

func (r *generalRepository) GetRedisCache(ctx context.Context, key string, destination interface{}) error {
	err := r.redisDb.Get(ctx, key, &destination)
	if err != nil {
		logrus.Errorf("error while setting redis cache: %v", err)
	}
	return err
}
