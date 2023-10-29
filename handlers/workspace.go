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

	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "workspace/cardsList.html", "workspace/cardDraggable.html", "workspace/card.html")
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

func (s *ApiRouter) handleReorderCards(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		s.handleError(w, r, err)
	}

	cards, _ := r.PostForm["item"]
	var orderedCards []*types.Card
	if ordered, err := s.store.ReorderCards(cards); err != nil {
		s.handleError(w, r, err)
		return
	} else {
		orderedCards = ordered
	}
	tmpl, err := template.ParseFS(templates.Templates, "workspace/cardDraggable.html", "workspace/card.html")
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	err = tmpl.Execute(w, orderedCards)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
}
