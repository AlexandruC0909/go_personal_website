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

func (s *ApiRouter) handleGetPosts(w http.ResponseWriter, r *http.Request) error {
	templatesDir := os.Getenv("TEMPLATES_DIR")
	if templatesDir == "" {
		fmt.Println("TEMPLATES_DIR environment variable is not set.")
	}

	tmplPathBase := fmt.Sprintf("%s/ui/base.html", templatesDir)
	tmplPathContent := fmt.Sprintf("%s/posts/postsList.html", templatesDir)

	files := []string{
		tmplPathBase,
		tmplPathContent,
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
