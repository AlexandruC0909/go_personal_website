package database

import (
	"database/sql"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewPostgresDbConnection() (*DbConnection, error) {

	dbname := os.Getenv("DB_NAME")
	fmt.Printf("DB_NAME: %s\n", dbName)
	dbPassword := os.Getenv("DB_PASSWORD")
	fmt.Printf("Pass: %s\n", dbPassword)
	dbUser := os.Getenv("DB_USER")
	fmt.Printf("user: %s\n", dbUser)

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
