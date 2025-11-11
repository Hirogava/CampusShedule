package main

import (
	"os"

	"github.com/Hirogava/CampusShedule/internal/config/environment"
	"github.com/Hirogava/CampusShedule/internal/config/logger"
	"github.com/Hirogava/CampusShedule/internal/repository/postgres"
	"github.com/Hirogava/CampusShedule/internal/transport/maxbot"
)

func main() {
	environment.LoadEnvFile(".env")

	logger.LogInit()
	logger.Logger.Info("Starting Swifty Gasprom backend server")

	dbConnStr := os.Getenv("DB_CONNECT_STRING")
	if dbConnStr == "" {
		logger.Logger.Fatal("DB_CONNECT_STRING environment variable is required")
	}
	logger.Logger.Info("Connecting to database", "connection_string", dbConnStr)

	manager := postgres.NewManager("postgres", dbConnStr)
	logger.Logger.Info("Database connection established successfully")

	logger.Logger.Info("Running database migrations")
	manager.Migrate()
	logger.Logger.Info("Database migrations completed successfully")

	maxbot.SetMaxConf(manager)
}