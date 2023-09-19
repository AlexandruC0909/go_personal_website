package database

import (
	"database/sql"
	"fmt"

	"go_api/types"

	_ "github.com/lib/pq"
)

type Methods interface {
	CreateUser(*types.User) error
	DeleteUser(int) error
	UpdateUser(*types.User) error
	GetUserById(int) (*types.User, error)
	GetUsers() ([]*types.User, error)
}

func (s *DbConnection) createUserTable() error {
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

func (s *DbConnection) CreateUser(user *types.User) error {
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

func (s *DbConnection) UpdateUser(user *types.User) error {

	query := `update users set first_name = $1 , last_name = $2 where id = $3`

	_, err := s.db.Exec(
		query,
		user.FirstName,
		user.LastName,
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *DbConnection) DeleteUser(id int) error {
	_, err := s.db.Query("delete from users where id = $1", id)
	return err
}

func (s *DbConnection) GetUserById(id int) (*types.User, error) {
	rows, err := s.db.Query("select * from users where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("user %d not found", id)
}

func (s *DbConnection) GetUsers() ([]*types.User, error) {
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
