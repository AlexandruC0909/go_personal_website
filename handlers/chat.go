package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	templates "go_api/templates"
)

type ChantNickname struct {
	Nickname string `json:"nickname"`
}
type ChatHandler interface {
	handleChat(w http.ResponseWriter, r *http.Request) error
	handleChatLogin(w http.ResponseWriter, r *http.Request) error
}

func (s *ApiRouter) handleChat(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("nickname")
	if cookie != nil {
		tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "chat/chat.html")
		if err != nil {
			s.handleError(w, r, err)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			s.handleError(w, r, err)
			return
		}
	} else {
		tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "chat/login.html")
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

}

func (s *ApiRouter) handleChatLogin(w http.ResponseWriter, r *http.Request) {
	var req ChantNickname
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.handleError(w, r, err)
		return
	}
	domain := os.Getenv("DOMAIN")

	http.SetCookie(w, &http.Cookie{
		Name:     "nickname",
		Value:    req.Nickname,
		HttpOnly: true,
		Path:     "/",
		Domain:   domain,
	})
	w.Header().Set("HX-Redirect", "/chat")
}
