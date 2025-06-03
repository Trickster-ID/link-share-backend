package mongo_repo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"linkshare/app/constants"
	"linkshare/app/global/helper"
	"linkshare/app/global/model"
	"linkshare/app/models"
	"net/http"
)

type IAccessTokenSessionsRepository interface {
	GetByAccessToken(accessToken string, ctx context.Context) (*models.AccessTokenSession, *model.ErrorLog)
	Insert(request *models.AccessTokenSession, ctx context.Context)
}

type accessTokenSessionsRepository struct {
	mongo          *mongo.Client
	collectionName string
}

func NewAccessTokenSessionRepository(mongo *mongo.Client) IAccessTokenSessionsRepository {
	return &accessTokenSessionsRepository{
		mongo:          mongo,
		collectionName: constants.ACCESS_TOKEN_SESSIONS_COL,
	}
}

func (r *accessTokenSessionsRepository) GetByAccessToken(accessToken string, ctx context.Context) (*models.AccessTokenSession, *model.ErrorLog) {
	data := &models.AccessTokenSession{}
	col := r.mongo.Database(constants.MONGO_DATABASE_NAME).Collection(r.collectionName)
	err := col.FindOne(ctx, bson.M{"access_token": accessToken}).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			errLog := helper.WriteLog(err, http.StatusNotFound, "")
			return nil, errLog
		}
		errLog := helper.WriteLog(err, http.StatusInternalServerError, "error while get session")
		return nil, errLog
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
}
