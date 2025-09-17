package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
    "github.com/joho/godotenv"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
     if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found, using system environment variables")
    }
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBHost := os.Getenv("DB_HOST")
	DBPort := os.Getenv("DB_PORT")
	DBName := os.Getenv("DB_NAME")
	sslMode := "disable"

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		DBUser, DBPassword, DBHost, DBPort, DBName, sslMode)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	// Automatically force dirty version if it exists
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatal(err)
	}

	if dirty {
		fmt.Printf("Database is dirty at version %d. Forcing...\n", version)
		if err := m.Force(int(version)); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Dirty migration forced successfully")
	}

	if len(os.Args) < 2 {
		log.Fatal("Specify migration command: up or down")
	}

	cmd := os.Args[1]
	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("Migrations applied successfully")
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Last migration rolled back")
	default:
		log.Fatal("Unknown command: use up or down")
	}
}
