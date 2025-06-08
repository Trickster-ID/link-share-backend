package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"linkshare/app/configuration"
	"linkshare/app/global/db"
	"linkshare/app/middleware"
	"linkshare/generated"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	configuration.InitialConfig()
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "main")
	}

	postgresParam := db.PostgresGetEnvVariable().NewPostgresParam()
	mongoParam := db.MongoGetEnvVariable().NewMongoParam()
	redisParam := db.RedisGetEnvVariable().NewRedisParam()

	switch os.Args[1] {
	case "main":
		app := fiber.New()
		configuration.FiberInitLogger(app)
		// Add RefreshTokenMiddleware to the /auth/refresh-token endpoint
		app.Use("/auth/refresh-token", middleware.RefreshTokenMiddleware())
		newServer := InitializeFiberServer(postgresParam, mongoParam, redisParam)
		generated.RegisterHandlers(app, newServer)

		startHttpServer(app)
	default:
		logrus.Fatal("invalid command")
	}

}

func startHttpServer(f *fiber.App) {

	// Channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	mainPort, err := strconv.Atoi(os.Getenv("MAIN_PORT"))
	if err != nil {
		logrus.Fatal(err)
	}
	go func() {
		if err := f.Listen(fmt.Sprintf(":%d", mainPort)); err != nil {
			logrus.Error(err)
			logrus.Fatal("shutting down http server")
		}
	}()

	// Block until a signal is received
	<-quit
	log.Print("Server is shutting down...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := f.ShutdownWithContext(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Print("Server stopped")
}
