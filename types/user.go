package types

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UpdateUserRequest struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
type UpdateUserImageRequest struct {
	ID       int    `json:"id"`
	ImageURL string `json:"image_url"`
}
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ImageURL  string    `json:"imageURL"`
}

func NewUser(firstName, lastName, email, password string) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(encpw),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func UpdateUser(id int, firstName, lastName string) *User {

	return &User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
	}
}

func (u *User) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw)) == nil
}
