package main

import (
	"database/sql"
	"fmt"
	"log"
    "github.com/yordanos-habtamu/realstate/config"
    "github.com/yordanos-habtamu/realstate/cmd/api"
    "github.com/yordanos-habtamu/realstate/db"
	
)

func main() {
	// Build DSN for PostgreSQL
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		config.Envs.DB_USER,
		config.Envs.DB_PWD,
		config.Envs.DB_ADDR, // e.g. "localhost:5432"
		config.Envs.DB_NAME,
	)
	fmt.Println(dsn)

	// Initialize PostgreSQL storage
	db, err := db.NewPostgresStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewApiServer(":4000", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database connected successfully")
}
