package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"go_api/types"

	templates "go_api/templates"
)

type PostResponse struct {
	Posts []*types.Post
	Page  int
}

func (s *ApiRouter) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	pagination := r.Context().Value("pagination").(map[string]int)
	page := pagination["page"]

	posts, err := s.store.GetPosts(page)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
	data := PostResponse{
		Posts: posts,
		Page:  page + 1,
	}
	if page == 1 {

		tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "posts/postsList.html", "posts/postsPaginated.html")
		if err != nil {
			s.handleError(w, r, err)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			s.handleError(w, r, err)
			return
		}
	} else {
		tmpl, err := template.ParseFS(templates.Templates, "posts/postsPaginated.html")
		if err != nil {
			s.handleError(w, r, err)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			s.handleError(w, r, err)
			return
		}
	}

}

func (s *ApiRouter) handleGetPost(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	post, err := s.store.GetPost(id)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "post/postDetails.html")
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	err = tmpl.Execute(w, post)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
}

func (s *ApiRouter) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	if err := s.store.DeletePost(id); err != nil {
		s.handleError(w, r, err)
		return
	}

}

func (s *ApiRouter) handleEditPost(w http.ResponseWriter, r *http.Request) {

	updatePostReq := new(types.UpdatePostRequest)

	if err := json.NewDecoder(r.Body).Decode(updatePostReq); err != nil {
		s.handleError(w, r, err)
		return
	}

	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	post := types.UpdatePost(id, updatePostReq.Name, updatePostReq.Content)

	if err := s.store.UpdatePost(post); err != nil {
		s.handleError(w, r, err)
		return
	} else {
		post, err := s.store.GetPost(post.ID)
		tmpl, err := template.ParseFS(templates.Templates, "post/postRow.html")
		if err != nil {
			s.handleError(w, r, err)
			return
		}

		err = tmpl.Execute(w, post)
		if err != nil {
			s.handleError(w, r, err)
			return
		}
	}

}
