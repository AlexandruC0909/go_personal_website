package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go_api/types"
)

type UserHandler interface {
	handleGetUsers(w http.ResponseWriter, r *http.Request) error
	handleUserById(w http.ResponseWriter, r *http.Request) error
	handleDeleteUser(w http.ResponseWriter, r *http.Request) error
	handleUpdateUser(w http.ResponseWriter, r *http.Request) error
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
