package database

import (
	"database/sql"
	"fmt"
	"time"

	"go_api/types"

	_ "github.com/lib/pq"
)

type Methods interface {
	CreateUser(*types.User) error
	DeleteUser(int) error
	UpdateUser(*types.User) error
	GetUser(int) (*types.User, error)
	GetUserByEmail(string) (*types.User, error)
	GetUsers() ([]*types.User, error)
}

func (s *DbConnection) createUserTable() error {
	query := `create table if not exists users (
		id serial primary key,
		first_name varchar(100),
		last_name varchar(100),
		email varchar(100),
		password varchar(100),
		created_at timestamp,
		updated_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
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

func (s *DbConnection) CreateUser(user *types.User) error {
	query := `insert into users 
	(first_name, last_name, email, password, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *DbConnection) GetUser(id int) (*types.User, error) {
	rows, err := s.db.Query("select * from users where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("user %d not found", id)
}

func (s *DbConnection) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("select * from users where email = $1", email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("user %d not found", email)
}

func (s *DbConnection) UpdateUser(user *types.User) error {
	updateQuery := `update users set first_name = $1 , last_name = $2, updated_at= $3 where id = $4 RETURNING id`

	var userId int
	err := s.db.QueryRow(
		updateQuery,
		user.FirstName,
		user.LastName,
		time.Now(),
		user.ID,
	).Scan(&userId)

	if err == sql.ErrNoRows {
		return fmt.Errorf("user %d not found", user.ID)
	} else if err != nil {
		return err
	}

	return nil

}

func (s *DbConnection) DeleteUser(id int) error {
	deleteQuery := `DELETE FROM users WHERE id = $1 RETURNING id`

	var userID int
	err := s.db.QueryRow(deleteQuery, id).Scan(&userID)

	if err == sql.ErrNoRows {
		return fmt.Errorf("user with id %d not found", id)
	} else if err != nil {
		return err
	}

	return nil
}

func scanIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt)

	return user, err
}
