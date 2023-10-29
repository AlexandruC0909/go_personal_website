package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"go_api/types"

	templates "go_api/templates"
)

func (s *ApiRouter) handleGetCards(w http.ResponseWriter, r *http.Request) {
	cards, err := s.store.GetCards()

	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "card/cardsList.html", "card/cardRow.html")
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	err = tmpl.Execute(w, cards)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
}

func (s *ApiRouter) handleGetCard(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	card, err := s.store.GetCard(id)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "card/cardDetails.html")
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	err = tmpl.Execute(w, card)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
}

func (s *ApiRouter) handleDeleteCard(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	if err := s.store.DeleteCard(id); err != nil {
		s.handleError(w, r, err)
		return
	}

}

func (s *ApiRouter) handleEditCard(w http.ResponseWriter, r *http.Request) {

	updateCardReq := new(types.UpdateCardRequest)

	if err := json.NewDecoder(r.Body).Decode(updateCardReq); err != nil {
		s.handleError(w, r, err)
		return
	}

	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	card := types.UpdateCard(id, updateCardReq.Name, updateCardReq.Content, updateCardReq.Position, updateCardReq.Parent)

	if err := s.store.UpdateCard(card); err != nil {
		s.handleError(w, r, err)
		return
	} else {
		card, err := s.store.GetCard(card.ID)
		tmpl, err := template.ParseFS(templates.Templates, "card/cardRow.html")
		if err != nil {
			s.handleError(w, r, err)
			return
		}

		err = tmpl.Execute(w, card)
		if err != nil {
			s.handleError(w, r, err)
			return
		}
	}

}

func (s *ApiRouter) handlgeGetCardEditRow(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
	card, err := s.store.GetCard(id)

	tmpl, err := template.ParseFS(templates.Templates, "card/cardEditRow.html")
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	err = tmpl.Execute(w, card)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

}

func (s *ApiRouter) HandleGetCardRow(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
	card, err := s.store.GetCard(id)

	tmpl, err := template.ParseFS(templates.Templates, "card/cardRow.html")
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	err = tmpl.Execute(w, card)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

}
