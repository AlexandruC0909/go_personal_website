package database

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/lib/pq"
)

type DbConnection struct {
	db      *sql.DB
	migrate *migrate.Migrate
}
