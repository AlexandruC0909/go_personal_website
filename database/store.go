package database

import (
	"go_api/types"
)

type Methods interface {
	CreateUser(*types.User) error
	DeleteUser(int) error
	UpdateUser(*types.User) error
	UpdateUserImage(*types.User) error
	GetUser(int) (*types.User, error)
	GetUserByEmail(string) (*types.User, error)
	GetUsers() ([]*types.User, error)

	CreateCard(*types.Card) error
	DeleteCard(int) error
	UpdateCard(*types.Card) error
	ReorderCards([]string) error

	GetCard(int) (*types.Card, error)
	GetCards() ([]*types.Card, error)
}
