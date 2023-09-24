package handlers

import (
	"html/template"
	"net/http"
)

type PostsHandler interface {
	handleGetPosts(w http.ResponseWriter, r *http.Request) error
}

func (s *ApiRouter) handleGetPosts(w http.ResponseWriter, r *http.Request) error {
	files := []string{
		"../templates/ui/base.html",
		"../templates/posts/postsList.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return err
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}
	return nil
}
