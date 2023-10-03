package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type PostsHandler interface {
	handleGetPosts(w http.ResponseWriter, r *http.Request) error
}

func (s *ApiRouter) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	templatesDir := os.Getenv("TEMPLATES_DIR")
	if templatesDir == "" {
		fmt.Println("TEMPLATES_DIR environment variable is not set.")
	}

	tmplPathBase := fmt.Sprintf("%s/ui/base.html", templatesDir)
	tmplPathNav := fmt.Sprintf("%s/ui/navbar.html", templatesDir)
	tmplPathContent := fmt.Sprintf("%s/posts/postsList.html", templatesDir)

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
	err = tmpl.Execute(w, nil)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
}
