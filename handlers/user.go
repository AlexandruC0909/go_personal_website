package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

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
	templatesDir := os.Getenv("TEMPLATES_DIR")
	if templatesDir == "" {
		fmt.Println("TEMPLATES_DIR environment variable is not set.")
	}

	tmplPathBase := fmt.Sprintf("%s/ui/base.html", templatesDir)
	tmplPathNav := fmt.Sprintf("%s/ui/navbar.html", templatesDir)
	tmplPathContent := fmt.Sprintf("%s/user/usersList.html", templatesDir)

	files := []string{
		tmplPathBase,
		tmplPathNav,
		tmplPathContent,
	}
	tmpl, err := template.ParseFiles(files...)

	if err != nil {
		return err
	}
	err = tmpl.Execute(w, users)
	if err != nil {
		return err
	}
	return nil
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

		if err != nil {
			return err
		}
		templatesDir := os.Getenv("TEMPLATES_DIR")
		if templatesDir == "" {
			fmt.Println("TEMPLATES_DIR environment variable is not set.")
		}

		tmplPathBase := fmt.Sprintf("%s/ui/base.html", templatesDir)
		tmplPathNav := fmt.Sprintf("%s/ui/navbar.html", templatesDir)
		tmplPathContent := fmt.Sprintf("%s/user/userDetails.html", templatesDir)

		files := []string{
			tmplPathBase,
			tmplPathNav,
			tmplPathContent,
		}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			return err
		}
		err = tmpl.Execute(w, user)
		if err != nil {
			return err
		}
		return nil
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
