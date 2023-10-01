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

	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")

	connString := "user=" + dbUser + " dbname=" + dbName + " password=" + dbPassword + " sslmode=disable"
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

	return &DbConnection{
		db:      db,
		migrate: m,
	}, nil
}

func (c *DbConnection) CloseDB() {
	if c.db != nil {
		c.db.Close()
	}
}

func (c *DbConnection) RunMigrations() error {
	if err := c.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
