package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/joho/godotenv"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewPostgresDbConnection() (*DbConnection, error) {
	if os.Getenv("GO_ENV) != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	
	dbname := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")

	connString := "user=" + dbUser + " dbname=" + dbname + " password=" + dbPassword + " sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres",
		driver)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DbConnection{
		db: db,
	}, nil
}
