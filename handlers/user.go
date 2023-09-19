package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go_api/types"
)

type UserHandler interface {
	handleLogin(w http.ResponseWriter, r *http.Request) error
	handleUsers(w http.ResponseWriter, r *http.Request) error
	handleGetUsers(w http.ResponseWriter, r *http.Request) error
	handleUserById(w http.ResponseWriter, r *http.Request) error
	handleDeleteUser(w http.ResponseWriter, r *http.Request) error
	handleUpdateUser(w http.ResponseWriter, r *http.Request) error
}

func (s *ApiRouter) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	user, err := s.store.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}

	if !user.ValidPassword(req.Password) {
		return fmt.Errorf("not authenticated")
	}

	token, err := createJWT(user)
	if err != nil {
		return err
	}

	resp := types.LoginResponse{
		Token: token,
		Email: user.Email,
	}

	return WriteJSON(w, http.StatusOK, resp)
}

func (s *ApiRouter) handleUsers(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetUsers(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateUser(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *ApiRouter) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := s.store.GetUsers()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, users)
}

func (s *ApiRouter) handleUserById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		user, err := s.store.GetUser(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, user)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteUser(w, r)
	}

	if r.Method == "PUT" {
		return s.handleUpdateUser(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *ApiRouter) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	createAccReq := new(types.CreateUserRequest)

	if err := json.NewDecoder(r.Body).Decode(createAccReq); err != nil {
		return err
	}

	user, err := types.NewUser(createAccReq.FirstName, createAccReq.LastName, createAccReq.Email, createAccReq.Password)
	if err != nil {
		return err
	}
	if err := s.store.CreateUser(user); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *ApiRouter) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	updateUserReq := new(types.UpdateUserRequest)

	if err := json.NewDecoder(r.Body).Decode(updateUserReq); err != nil {
		return err
	}

	id, err := getID(r)
	if err != nil {
		return err
	}

	user := types.UpdateUser(id, updateUserReq.FirstName, updateUserReq.LastName)

	if err := s.store.UpdateUser(user); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *ApiRouter) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteUser(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}
