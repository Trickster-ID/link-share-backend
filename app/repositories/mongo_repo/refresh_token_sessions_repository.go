package mongo_repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"linkshare/app/constants"
	"linkshare/app/global/helper"
	"linkshare/app/global/model"
	"linkshare/app/models"
	"net/http"
	"time"
)

type IRefreshTokenSessionsRepository interface {
	GetByRefreshToken(refreshToken string, ctx context.Context) (*models.RefreshTokenSession, *model.ErrorLog)
	Insert(request *models.RefreshTokenSession, ctx context.Context)
}

type refreshTokenSessionsRepository struct {
	mongo          *mongo.Client
	redis          *redis.Client
	collectionName string
}

func NewRefreshTokenSessionRepository(mongo *mongo.Client, redis *redis.Client) IRefreshTokenSessionsRepository {
	return &refreshTokenSessionsRepository{
		mongo:          mongo,
		redis:          redis,
		collectionName: constants.REFRESH_TOKEN_SESSIONS_COL,
	}
}

func (r *refreshTokenSessionsRepository) GetByRefreshToken(refreshToken string, ctx context.Context) (*models.RefreshTokenSession, *model.ErrorLog) {
	data := &models.RefreshTokenSession{}
	redisKey := fmt.Sprintf("%s:%s", r.collectionName, refreshToken)
	redisResult, err := r.redis.Get(ctx, redisKey).Result()
	if errors.Is(err, redis.Nil) {
		col := r.mongo.Database(constants.MONGO_DATABASE_NAME).Collection(r.collectionName)
		err = col.FindOne(ctx, bson.M{"refresh_token": refreshToken}).Decode(&data)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				errLog := helper.WriteLog(err, http.StatusNotFound, "")
				return nil, errLog
			}
			errLog := helper.WriteLog(err, http.StatusInternalServerError, "error while get session")
			return nil, errLog
		}
	} else if err != nil {
		errLog := helper.WriteLog(err, http.StatusInternalServerError, "error while get refresh session on redis")
		return nil, errLog
	} else {
		err = json.Unmarshal([]byte(redisResult), &data)
		if err != nil {
			errLog := helper.WriteLog(err, http.StatusInternalServerError, "error while get session")
			return nil, errLog
		}
	}
	return data, nil
}

func (r *refreshTokenSessionsRepository) Insert(request *models.RefreshTokenSession, ctx context.Context) {
	request.ObjectID = primitive.NewObjectID()
	col := r.mongo.Database(constants.MONGO_DATABASE_NAME).Collection(r.collectionName)
	_, err := col.InsertOne(ctx, request)
	if err != nil {
		helper.WriteLog(err, http.StatusInternalServerError, "error while insert session")
	}
	redisKey := fmt.Sprintf("%s:%s", r.collectionName, request.RefreshToken)
	marshalledData, err := json.Marshal(request)
	if err != nil {
		helper.WriteLog(err, http.StatusInternalServerError, "error while marshalling session")
	}
	err = r.redis.Set(ctx, redisKey, marshalledData, time.Until(request.Expired)).Err()
	if err != nil {
		helper.WriteLog(err, http.StatusInternalServerError, "error while insert session")
	}
}
