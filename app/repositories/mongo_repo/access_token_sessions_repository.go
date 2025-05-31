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

type IAccessTokenSessionsRepository interface {
	GetByAccessToken(accessToken string, ctx context.Context) (*models.AccessTokenSession, *model.ErrorLog)
	Insert(request *models.AccessTokenSession, ctx context.Context)
}

type accessTokenSessionsRepository struct {
	mongo          *mongo.Client
	redis          *redis.Client
	collectionName string
}

func NewAccessTokenSessionRepository(mongo *mongo.Client, redis *redis.Client) IAccessTokenSessionsRepository {
	return &accessTokenSessionsRepository{
		mongo:          mongo,
		redis:          redis,
		collectionName: constants.ACCESS_TOKEN_SESSIONS_COL,
	}
}

func (r *accessTokenSessionsRepository) GetByAccessToken(accessToken string, ctx context.Context) (*models.AccessTokenSession, *model.ErrorLog) {
	data := &models.AccessTokenSession{}
	redisKey := fmt.Sprintf("%s:%s", r.collectionName, accessToken)
	redisResult, err := r.redis.Get(ctx, redisKey).Result()
	if errors.Is(err, redis.Nil) {
		col := r.mongo.Database(constants.MONGO_DATABASE_NAME).Collection(r.collectionName)
		err = col.FindOne(ctx, bson.M{"access_token": accessToken}).Decode(&data)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				errLog := helper.WriteLog(err, http.StatusNotFound, "")
				return nil, errLog
			}
			errLog := helper.WriteLog(err, http.StatusInternalServerError, "error while get session")
			return nil, errLog
		}
	} else if err != nil {
		errLog := helper.WriteLog(err, http.StatusInternalServerError, "error while get access session on redis")
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

func (r *accessTokenSessionsRepository) Insert(request *models.AccessTokenSession, ctx context.Context) {
	request.ObjectID = primitive.NewObjectID()
	col := r.mongo.Database(constants.MONGO_DATABASE_NAME).Collection(r.collectionName)
	_, err := col.InsertOne(ctx, request)
	if err != nil {
		helper.WriteLog(err, http.StatusInternalServerError, "error while insert session")
	}
	redisKey := fmt.Sprintf("%s:%s", r.collectionName, request.AccessToken)
	marshalledData, err := json.Marshal(request)
	if err != nil {
		helper.WriteLog(err, http.StatusInternalServerError, "error while marshalling session")
	}
	err = r.redis.Set(ctx, redisKey, marshalledData, time.Until(request.Expired)).Err()
	if err != nil {
		helper.WriteLog(err, http.StatusInternalServerError, "error while insert session")
	}
}
