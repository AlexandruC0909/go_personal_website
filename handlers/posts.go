package handlers

import (
	"html/template"
	"net/http"

	templates "go_api/templates"
)

type PostsHandler interface {
	handleGetPosts(w http.ResponseWriter, r *http.Request) error
}

func (s *ApiRouter) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "posts/postsList.html")
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
}
