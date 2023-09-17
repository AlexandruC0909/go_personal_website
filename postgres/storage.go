package postgres

import (
	"database/sql"
	"fmt"

	"go_api/types"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*types.User) error
	DeleteUser(int) error
	UpdateUser(*types.User) error
	GetUserById(int) (*types.User, error)
	GetUsers() ([]*types.User, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgressStore() (*PostgresStore, error) {
	connString := "user=go_api_root dbname=postgres password=go_api_root sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}
func (s *PostgresStore) Init() error {
	return s.createUserTable()
}

func (s *PostgresStore) createUserTable() error {
	query := `create table if not exists users (
		id serial primary key,
		first_name varchar(100),
		last_name varchar(100),
		number serial,
		encrypted_password varchar(100),
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateUser(user *types.User) error {
	query := `insert into users 
	(first_name, last_name, number, encrypted_password, created_at)
	values ($1, $2, $3, $4, $5)`

	_, err := s.db.Query(
		query,
		user.FirstName,
		user.LastName,
		user.Number,
		user.EncryptedPassword,
		user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateUser(user *types.User) error {
	query := `update users set 
	first_name = $1 , last_name = $2`

	_, err := s.db.Query(
		query,
		user.FirstName,
		user.LastName,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteUser(id int) error {
	_, err := s.db.Query("delete from users where id = $1", id)
	return err
}

func (s *PostgresStore) GetUserById(id int) (*types.User, error) {
	rows, err := s.db.Query("select * from users where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("user %d not found", id)
}

func (s *PostgresStore) GetUsers() ([]*types.User, error) {
	rows, err := s.db.Query("select * from users")
	if err != nil {
		return nil, err
	}

	users := []*types.User{}
	for rows.Next() {
		user, err := scanIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func scanIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Number,
		&user.EncryptedPassword,
		&user.CreatedAt)

	return user, err
}
