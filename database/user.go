package database

import (
	"database/sql"
	"fmt"
	"time"

	"go_api/types"

	_ "github.com/lib/pq"
)

const getUserQuery = "SELECT u.id, u.first_name, u.last_name, u.email, u.password, u.created_at, u.updated_at,image_url, r.id as role_id, r.name as role_name FROM users u JOIN roles r ON u.roles_id = r.id "

type Methods interface {
	CreateUser(*types.User) error
	DeleteUser(int) error
	UpdateUser(*types.User) error
	UpdateUserImage(*types.User) error
	GetUser(int) (*types.User, error)
	GetUserByEmail(string) (*types.User, error)
	GetUsers() ([]*types.User, error)
}

func (s *DbConnection) GetUsers() ([]*types.User, error) {
	rows, err := s.db.Query(getUserQuery)
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
	(first_name, last_name, email, password, created_at, updated_at, roles_id, image_url)
	values ($1, $2, $3, $4, $5, $6, 2, '')`

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
	rows, err := s.db.Query(getUserQuery+"WHERE u.id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("user %d not found", id)
}

func (s *DbConnection) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query(getUserQuery+"where email = $1", email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("user %s not found", email)
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

func (s *DbConnection) UpdateUserImage(user *types.User) error {
	updateQuery := `update users set image_url = $1 , updated_at= $2 where id = $3 RETURNING id`

	var userId int
	err := s.db.QueryRow(
		updateQuery,
		user.ImageURL,
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

func scanIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	role := types.Role{}
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.ImageURL,
		&role.ID,
		&role.Name,
	)
	user.Role = role
	return user, err
}
