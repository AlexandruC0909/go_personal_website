package types

import (
	"time"
)

type UpdateCardRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	Position int    `json:"position"`
	Parent   *int   `json:"cards"`
}

type Card struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Position  int       `json:"position"`
	CardType  bool      `json:"cardType"`
	Parent    *int      `json:"cards"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type UpdateCardPositions struct {
	Items []int `json:"items"`
}

func NewCard(name string, content string, position int, cardType bool, parent *int) (*Card, error) {
	return &Card{
		Name:      name,
		Content:   content,
		Position:  position,
		CardType:  cardType,
		Parent:    parent,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func UpdateCard(id int, name, content string, position int, parent *int) *Card {

	return &Card{
		ID:        id,
		Name:      name,
		Content:   content,
		Position:  position,
		Parent:    parent,
		UpdatedAt: time.Now().UTC(),
	}
}
