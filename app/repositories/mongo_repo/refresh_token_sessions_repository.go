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

type IRefreshTokenSessionsRepository interface {
	GetByRefreshToken(refreshToken string, ctx context.Context) (*models.RefreshTokenSession, *model.ErrorLog)
	Insert(request *models.RefreshTokenSession, ctx context.Context)
}

type refreshTokenSessionsRepository struct {
	mongo          *mongo.Client
	collectionName string
}

func NewRefreshTokenSessionRepository(mongo *mongo.Client) IRefreshTokenSessionsRepository {
	return &refreshTokenSessionsRepository{
		mongo:          mongo,
		collectionName: constants.REFRESH_TOKEN_SESSIONS_COL,
	}
}

func (r *refreshTokenSessionsRepository) GetByRefreshToken(refreshToken string, ctx context.Context) (*models.RefreshTokenSession, *model.ErrorLog) {
	data := &models.RefreshTokenSession{}
	col := r.mongo.Database(constants.MONGO_DATABASE_NAME).Collection(r.collectionName)
	err := col.FindOne(ctx, bson.M{"refresh_token": refreshToken}).Decode(&data)
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

func (r *refreshTokenSessionsRepository) Insert(request *models.RefreshTokenSession, ctx context.Context) {
	request.ObjectID = primitive.NewObjectID()
	col := r.mongo.Database(constants.MONGO_DATABASE_NAME).Collection(r.collectionName)
	_, err := col.InsertOne(ctx, request)
	if err != nil {
		helper.WriteLog(err, http.StatusInternalServerError, "error while insert session")
	}
}
