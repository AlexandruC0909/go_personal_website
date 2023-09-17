package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DbConnection struct {
	db *sql.DB
}

func NewDbConnection() (*DbConnection, error) {
	connString := "user=go_api_root dbname=postgres password=go_api_root sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DbConnection{
		db: db,
	}, nil
}
func (s *DbConnection) Init() error {
	return s.createUserTable()
}
