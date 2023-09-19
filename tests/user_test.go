package tests

import (
	"testing"
	userType "go_api/types"
)

func TestNewUser(t *testing.T) {
	firstName := "John"
	lastName := "Doe"

	user := userType.NewUser(firstName, lastName)

	if user == nil {
		t.Error("Expected a user object, but got nil")
	}

	if user.FirstName != firstName {
		t.Errorf("Expected FirstName to be %s, but got %s", firstName, user.FirstName)
	}

	if user.LastName != lastName {
		t.Errorf("Expected LastName to be %s, but got %s", lastName, user.LastName)
	}

	if user.Number <= 0 {
		t.Errorf("Expected Number to be a positive integer, but got %d", user.Number)
	}

	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be a valid time, but it's zero")
	}
}

func TestCreateUserRequest(t *testing.T) {
	req := userType.CreateUserRequest{
		FirstName: "Jane",
		LastName:  "Smith",
	}

	if req.FirstName != "Jane" {
		t.Errorf("Expected FirstName to be %s, but got %s", "Jane", req.FirstName)
	}

	if req.LastName != "Smith" {
		t.Errorf("Expected LastName to be %s, but got %s", "Smith", req.LastName)
	}
}

func TestUpdateUserRequest(t *testing.T) {
	req := userType.UpdateUserRequest{
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
