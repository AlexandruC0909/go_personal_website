package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

//go:embed migrations/*
var Migrations embed.FS

func NewPostgresDbConnection() (*DbConnection, error) {
	isProduction := os.Getenv("GO_ENV") == "production"

	if !isProduction {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	dbname := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")

	connString := "user=" + dbUser + " dbname=" + dbname + " password=" + dbPassword + " sslmode=disable"

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create database driver: %v", err)
	}

	sourceInstance, err := iofs.New(Migrations, "migrations")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create custom source driver: %v", err)
	}

	log.Printf("Using embedded migrations")
	m, err := migrate.NewWithInstance("iofs", sourceInstance, "postgres", driver)

	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		db.Close()
		return nil, fmt.Errorf("failed to apply migrations: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &DbConnection{
		db: db,
	}, nil
}
