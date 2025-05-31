package configuration

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func InitialConfig() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)

	err := godotenv.Load()
	if err != nil {
		logrus.Warnln("Error loading .env file, use os environment")
	}

	logLevelEnv := os.Getenv("LOG_LEVEL")

	// Default to Info level if not set or invalid
	logLevel, err := logrus.ParseLevel(strings.ToLower(logLevelEnv))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)
}

func InitialConfigForUnitTest() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
	// Default to Info level if not set or invalid
	logLevel, err := logrus.ParseLevel(strings.ToLower("panic"))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)
}

func FiberInitLogger(f *fiber.App) {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "info" || logLevel == "debug" || logLevel == "trace" {
		f.Use(logger.New())
	}
}
