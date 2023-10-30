package database

import (
	"go_api/types"
)

type Methods interface {
	CreateUser(*types.User) error
	GetUsers() ([]*types.User, error)
	GetUser(int) (*types.User, error)
	UpdateUser(*types.User) error
	DeleteUser(int) error
	UpdateUserImage(*types.User) error
	GetUserByEmail(string) (*types.User, error)

	CreateCard(*types.Card) error
	GetCards() ([]*types.Card, error)
	GetCard(int) (*types.Card, error)
	UpdateCard(*types.Card) error
	DeleteCard(int) error
	ReorderCards([]string) ([]*types.Card, error)

	CreatePost(*types.Post) error
	GetPosts(page, limit int) ([]*types.Post, error)
	GetPost(int) (*types.Post, error)
	UpdatePost(*types.Post) error
	DeletePost(int) error
}
