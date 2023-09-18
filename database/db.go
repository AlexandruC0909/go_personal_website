package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DbConnection struct {
	db *sql.DB
}
