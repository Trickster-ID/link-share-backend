package repositories

import (
	"context"
	"github.com/jackc/pgx/v5"
	"linkshare/app/global/db"
	"time"
)

type IGeneralRepository interface {
	BeginTransaction(ctx context.Context) pgx.Tx
	SetRedisCache(ctx context.Context, key string, value interface{}, expiration time.Duration)
	GetRedisCache(ctx context.Context, key string, destination interface{})
	DelRedisCache(ctx context.Context, key string)
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
	_ = r.redisDb.Set(ctx, key, value, expiration)
}

func (r *generalRepository) GetRedisCache(ctx context.Context, key string, destination interface{}) {
	_ = r.redisDb.Get(ctx, key, &destination)
}

func (r *generalRepository) DelRedisCache(ctx context.Context, key string) {
	err := r.redisDb.Del(ctx, key)
	if err != nil {
		// retry after second 2
		go func() {
			for i := 0; i < 5; i++ {
				time.Sleep(2 * time.Second)
				err := r.redisDb.Del(ctx, key)
				if err == nil {
					return
				}
			}
		}()
	}
}
