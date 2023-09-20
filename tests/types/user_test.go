//go:generate go test -coverpkg=go_api/types -coverprofile=coverage.out

package tests

import (
	userType "go_api/types"
	"testing"
)

func TestNewUser(t *testing.T) {
	firstName := "John"
	lastName := "Doe"
	email := "Doe@john.com"
	password := "password1234"

	user,_ := userType.NewUser(firstName, lastName, email, password)

	if user == nil {
		t.Error("Expected a user object, but got nil")
	}

	if user.FirstName != firstName {
		t.Errorf("Expected FirstName to be %s, but got %s", firstName, user.FirstName)
	}

	if user.LastName != lastName {
		t.Errorf("Expected LastName to be %s, but got %s", lastName, user.LastName)
	}

	if user.Email != email {
		t.Errorf("Expected Email to be %s, but got %s", email, user.Email)
	}

	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be a valid time, but it's zero")
	}

	if user.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be a valid time, but it's zero")
	}
}

func TestCreateUserRequest(t *testing.T) {
	req := userType.CreateUserRequest{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:  "jane@smith.com",
		password : "password1234",
	}

	if req.FirstName != "Jane" {
		t.Errorf("Expected FirstName to be %s, but got %s", "Jane", req.FirstName)
	}

	if req.LastName != "Smith" {
		t.Errorf("Expected LastName to be %s, but got %s", "Smith", req.LastName)
	}
	if req.LastName != "jane@smith.com" {
		t.Errorf("Expected Email to be %s, but got %s", "jane@smith.com", req.LastName)
	}
}

func TestUpdateUserRequest(t *testing.T) {
	req := userType.UpdateUserRequest{
		ID: 3,
		FirstName: "Alice",
		LastName:  "Johnson",
	}

	if req.FirstName != "Alice" {
		t.Errorf("Expected FirstName to be %s, but got %s", "Alice", req.FirstName)
	}

	if req.LastName != "Johnson" {
		t.Errorf("Expected LastName to be %s, but got %s", "Johnson", req.LastName)
	}
}
