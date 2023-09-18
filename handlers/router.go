package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	database "go_api/database"

	"github.com/gorilla/mux"
)

type ApiRouter struct {
	listenAddress string
	store         database.Methods
}

type ApiError struct {
	Error string `json:"error"`
}

func NewAPIServer(listenAddress string, store database.Methods) *ApiRouter {
	return &ApiRouter{
		listenAddress: listenAddress,
		store:         store,
	}
}
func (s *ApiRouter) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/users", makeHTTPHandleFunc(s.handleUsers))
	router.HandleFunc("/users/{id}", makeHTTPHandleFunc(s.handleUserById))

	log.Println("JSON API server running on port:", s.listenAddress)

	http.ListenAndServe(s.listenAddress, router)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
