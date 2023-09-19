package types

import (
	"math/rand"
	"time"
)

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UpdateUserRequest struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type User struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Number            int64     `json:"number"`
	EncryptedPassword string    `json:"-"`
	CreatedAt         time.Time `json:"createdAt"`
}

func NewUser(firstName, lastName string) *User {

	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(1000000)),
		CreatedAt: time.Now().UTC(),
	}
}

func UpdateUser(id int, firstName, lastName string) *User {

	return &User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
	}
}
