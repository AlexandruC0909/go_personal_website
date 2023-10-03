package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"go_api/types"
)

func (s *ApiRouter) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.store.GetUsers()

	if err != nil {
		s.handleError(w, r, err)
		return
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
		s.handleError(w, r, err)
		return
	}
	err = tmpl.Execute(w, users)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
}

func (s *ApiRouter) handleUserByIdGET(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	user, err := s.store.GetUser(id)
	if err != nil {
		s.handleError(w, r, err)
		return
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
		s.handleError(w, r, err)
		return
	}
	err = tmpl.Execute(w, user)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
}

func (s *ApiRouter) handleUserByIdDELETE(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	if err := s.store.DeleteUser(id); err != nil {
		s.handleError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *ApiRouter) handleUserByIdPUT(w http.ResponseWriter, r *http.Request) {
	updateUserReq := new(types.UpdateUserRequest)

	if err := json.NewDecoder(r.Body).Decode(updateUserReq); err != nil {
		s.handleError(w, r, err)
		return
	}

	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	user := types.UpdateUser(id, updateUserReq.FirstName, updateUserReq.LastName)

	if err := s.store.UpdateUser(user); err != nil {
		s.handleError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, user)
}
