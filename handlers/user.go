package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"go_api/types"

	templates "go_api/templates"
)

func (s *ApiRouter) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	
	users, err := s.store.GetUsers()

	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "user/usersList.html", "user/userRow.html")
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

func (s *ApiRouter) handleGetUser(w http.ResponseWriter, r *http.Request) {
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

	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "user/userDetails.html")
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

func (s *ApiRouter) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	if err := s.store.DeleteUser(id); err != nil {
		s.handleError(w, r, err)
		return
	}

}

func (s *ApiRouter) handleEditUser(w http.ResponseWriter, r *http.Request) {

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
	} else {
		user, err := s.store.GetUser(user.ID)
		tmpl, err := template.ParseFS(templates.Templates, "user/userRow.html")
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

}

func (s *ApiRouter) handlgeGetUserEditRow(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
	user, err := s.store.GetUser(id)

	tmpl, err := template.ParseFS(templates.Templates, "user/userEditRow.html")
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

func (s *ApiRouter) HandleGetUserRow(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
	user, err := s.store.GetUser(id)

	tmpl, err := template.ParseFS(templates.Templates, "user/userRow.html")
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
