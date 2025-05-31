//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"linkshare/app/controllers"
	"linkshare/app/global/db"
	"linkshare/app/repositories/mongo_repo"
	"linkshare/app/repositories/sql_repo"
	"linkshare/app/routes"
	"linkshare/app/security"
	"linkshare/app/usecases"
)

var connectionSet = wire.NewSet(
	db.NewPostgresClient,
	db.NewMongoClient,
	db.NewRedisClient,
)

var controllerSet = wire.NewSet(
	controllers.NewAuthController,
)

var useCaseSet = wire.NewSet(
	usecases.NewAuthUseCase,
)

var repositorySet = wire.NewSet(
	sql_repo.NewAuthRepository,
	mongo_repo.NewAccessTokenSessionRepository,
	mongo_repo.NewRefreshTokenSessionRepository,
)

var securitySet = wire.NewSet(
	security.NewJwtSecurity,
)

func InitializeFiberServer(postgresParam db.PostgresParam, mongoParam db.MongoParam, redisParam db.RedisParam) *fiber.App {
	wire.Build(
		connectionSet,
		controllerSet,
		useCaseSet,
		repositorySet,
		securitySet,
		routes.NewRouter,
	)
	return nil
}
