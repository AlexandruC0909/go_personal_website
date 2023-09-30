package tests

import (
	authType "go_api/types"
	"testing"
)
func TestRegisterRequest(t *testing.T) {
	req := authType.RegisterRequest{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:  "jane@smith.com",
		Password : "password1234",
	}

	if req.FirstName != "Jane" {
		t.Errorf("Expected FirstName to be %s, but got %s", "Jane", req.FirstName)
	}

	if req.LastName != "Smith" {
		t.Errorf("Expected LastName to be %s, but got %s", "Smith", req.LastName)
	}
	if req.Email != "jane@smith.com" {
		t.Errorf("Expected Email to be %s, but got %s", "jane@smith.com", req.Email)
	}
}


func TestLoginRequest(t *testing.T) {
	req := authType.LoginRequest{
		Email:    "user@example.com",
		Password: "password1234",
	}

	if req.Email != "user@example.com" {
		t.Errorf("Expected Email to be %s, but got %s", "user@example.com", req.Email)
	}

	if req.Password != "password1234" {
		t.Errorf("Expected Password to be %s, but got %s", "password1234", req.Password)
	}
}

func TestLoginResponse(t *testing.T) {
	res := authType.LoginResponse{
		Email: "user@example.com",
		Token: "some-auth-token",
	}

	if res.Email != "user@example.com" {
		t.Errorf("Expected Email to be %s, but got %s", "user@example.com", res.Email)
	}

	if res.Token != "some-auth-token" {
		t.Errorf("Expected Token to be %s, but got %s", "some-auth-token", res.Token)
	}
}