package db

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"linkshare/app/global/helper"
	"os"
	"time"
)

type MongoRawCredential struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type MongoParam struct {
	MongoUrl string
}

func (m *MongoRawCredential) NewMongoParam() MongoParam {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s", m.Username, m.Password, m.Host, m.Port)
	if !helper.IsValidIP(m.Host) {
		url = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", m.Username, m.Password, m.Host, m.Database)
	}
	return MongoParam{
		MongoUrl: url,
	}
}

func MongoGetEnvVariable() *MongoRawCredential {
	return &MongoRawCredential{
		Host:     os.Getenv("MONGO_HOST"),
		Port:     os.Getenv("MONGO_PORT"),
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Database: os.Getenv("MONGO_DATABASE"),
	}
}

func NewMongoClient(dbURL MongoParam) *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbURL.MongoUrl))
	if err != nil {
		logrus.Fatal("mongo connect err: ", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logrus.Fatal("mongo connect err: ", err)
	}
	return client
}
