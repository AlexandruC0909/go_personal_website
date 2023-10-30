package types

import (
	"time"
)

type UpdatePostRequest struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Post struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	ImageURL  *string   `json:"image_url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewPost(name, content string, imageUrl *string) (*Post, error) {
	return &Post{
		Name:      name,
		Content:   content,
		ImageURL:  imageUrl,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func UpdatePost(id int, name, content string) *Post {

	return &Post{
		ID:        id,
		Name:      name,
		Content:   content,
		UpdatedAt: time.Now().UTC(),
	}
}
