package database

import (
	"database/sql"
	"fmt"
	"time"

	"go_api/types"

	_ "github.com/lib/pq"
)

const getPostQuery = "SELECT p.id, p.name, p.content, image_url FROM posts p ORDER BY p.created_at"

func (s *DbConnection) GetPosts(page int) ([]*types.Post, error) {
	offset := (page - 1) * 10
	query := fmt.Sprintf("%s LIMIT %d OFFSET %d", getPostQuery, 10, offset)

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}

	posts := []*types.Post{}
	for rows.Next() {
		post, err := scanIntoPost(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *DbConnection) CreatePost(post *types.Post) error {
	query := `insert into posts 
	(name, content, image_url, created_at, updated_at)
	values ($1, $2, $3, $4, $5)`

	_, err := s.DB.Query(
		query,
		post.Name,
		post.Content,
		post.ImageURL,
		post.CreatedAt,
		post.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *DbConnection) GetPost(id int) (*types.Post, error) {
	rows, err := s.DB.Query(getPostQuery+"WHERE p.id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoPost(rows)
	}

	return nil, fmt.Errorf("post %d not found", id)
}

func (s *DbConnection) UpdatePost(post *types.Post) error {
	updateQuery := `update posts set name = $1 , content = $2, image_url = $3, updated_at= $5 where id = $6 RETURNING id`

	var postId int
	err := s.DB.QueryRow(
		updateQuery,
		post.Name,
		post.Content,
		post.ImageURL,
		time.Now(),
		post.ID,
	).Scan(&postId)

	if err == sql.ErrNoRows {
		return fmt.Errorf("post %d not found", post.ID)
	} else if err != nil {
		return err
	}

	return nil

}

func (s *DbConnection) DeletePost(id int) error {
	deleteQuery := `DELETE FROM posts WHERE id = $1 RETURNING id`

	var postID int
	err := s.DB.QueryRow(deleteQuery, id).Scan(&postID)

	if err == sql.ErrNoRows {
		return fmt.Errorf("post with id %d not found", id)
	} else if err != nil {
		return err
	}

	return nil
}

func scanIntoPost(rows *sql.Rows) (*types.Post, error) {
	post := new(types.Post)
	err := rows.Scan(
		&post.ID,
		&post.Name,
		&post.Content,
		&post.ImageURL,
	)
	return post, err
}
