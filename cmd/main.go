package main

import (
	"database/sql"
	"fmt"

	"github.com/yordanos-habtamu/realstate/cmd/api"
	"github.com/yordanos-habtamu/realstate/config"
	"github.com/yordanos-habtamu/realstate/db"
	"github.com/yordanos-habtamu/realstate/internal/logger"
)

func main() {
	// Initialize logger (use "production" for prod, "development" for dev)
	logger.InitLogger("development")
	defer logger.Sync()

	logger.Log.Info("Starting SMARTHOMEZ Backend Server...")

	// Build DSN for PostgreSQL
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		config.Envs.DB_USER,
		config.Envs.DB_PWD,
		config.Envs.DB_ADDR,
		config.Envs.DB_NAME,
	)

	// Initialize PostgreSQL storage
	database, err := db.NewPostgresStorage(dsn)
	if err != nil {
		logger.Log.Fatalw("Failed to connect to database",
			"error", err,
		)
	}

	initStorage(database)

	server := api.NewApiServer(":"+config.Envs.PORT, database)
	if err := server.Run(); err != nil {
		logger.Log.Fatalw("Server failed to start",
			"error", err,
		)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		logger.Log.Fatalw("Database ping failed",
			"error", err,
		)
	}
	logger.Log.Info("Database connected successfully")
}
