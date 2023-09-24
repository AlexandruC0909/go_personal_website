package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	database "go_api/database"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Replace with your front-end URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	})

	fs := http.FileServer(http.Dir("/uploads/"))

	router := mux.NewRouter()

	router.Use(c.Handler)
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", fs))
	router.HandleFunc("/auth/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/auth/refresh", makeHTTPHandleFunc(s.handleRefresh))
	router.HandleFunc("/auth/register", makeHTTPHandleFunc(s.handleRegister))
	router.HandleFunc("/users", withJWTAuth(makeHTTPHandleFunc(s.handleGetUsers), s.store))

	router.HandleFunc("/users/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleUserById), s.store))
	router.HandleFunc("/users/{id}/upload", withJWTAuth(makeHTTPHandleFunc(s.UploadImages), s.store))

	router.HandleFunc("/posts", withJWTAuth(makeHTTPHandleFunc(s.handleGetPosts), s.store))
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
