package database

import (
	"database/sql"
	"fmt"
	"time"

	"go_api/types"

	_ "github.com/lib/pq"
)

const getCardQuery = "SELECT c.id, c.name, c.content, c.position, c.type, c.parent_id FROM cards c ORDER BY c.position"

func (s *DbConnection) GetCards() ([]*types.Card, error) {
	rows, err := s.DB.Query(getCardQuery)
	if err != nil {
		return nil, err
	}

	cards := []*types.Card{}
	for rows.Next() {
		card, err := scanIntoCard(rows)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (s *DbConnection) CreateCard(card *types.Card) error {
	query := `insert into cards 
	(name, content, position, type, parent_id, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.DB.Query(
		query,
		card.Name,
		card.Content,
		card.Position,
		card.CardType,
		card.Parent,
		card.CreatedAt,
		card.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *DbConnection) GetCard(id int) (*types.Card, error) {
	rows, err := s.DB.Query(getCardQuery+"WHERE c.id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoCard(rows)
	}

	return nil, fmt.Errorf("card %d not found", id)
}

func (s *DbConnection) UpdateCard(card *types.Card) error {
	updateQuery := `update cards set name = $1 , content = $2, position = $3, parent_id = $4, updated_at= $5 where id = $6 RETURNING id`

	var cardId int
	err := s.DB.QueryRow(
		updateQuery,
		card.Name,
		card.Content,
		card.Position,
		card.Parent,
		time.Now(),
		card.ID,
	).Scan(&cardId)

	if err == sql.ErrNoRows {
		return fmt.Errorf("card %d not found", card.ID)
	} else if err != nil {
		return err
	}

	return nil

}

func (s *DbConnection) DeleteCard(id int) error {
	deleteQuery := `DELETE FROM cards WHERE id = $1 RETURNING id`

	var cardID int
	err := s.DB.QueryRow(deleteQuery, id).Scan(&cardID)

	if err == sql.ErrNoRows {
		return fmt.Errorf("card with id %d not found", id)
	} else if err != nil {
		return err
	}

	return nil
}

func (s *DbConnection) ReorderCards(id []string) error {
	updateQuery := `update cards set position = $1 where id = $2 RETURNING id`

	rows, err := s.DB.Query(getCardQuery)
	if err != nil {
		return err
	}

	cards := []*types.Card{}
	for rows.Next() {
		card, err := scanIntoCard(rows)
		if err != nil {
			return err
		}
		cards = append(cards, card)
	}

	for i, card := range cards {
		var cardId int
		updateErr := s.DB.QueryRow(
			updateQuery,
			id[i],
			card.ID,
		).Scan(&cardId)
		if updateErr != nil {
			return updateErr
		}
	}

	return nil

}

func scanIntoCard(rows *sql.Rows) (*types.Card, error) {
	card := new(types.Card)
	err := rows.Scan(
		&card.ID,
		&card.Name,
		&card.Content,
		&card.Position,
		&card.CardType,
		&card.Parent,
	)
	return card, err
}
