package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ApiError struct {
	Error string `json:"error"`
}

type APIServer struct {
	listenAddress string
	store Storage

}

func NewAPIServer(listenAddress string, store Storage) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store: store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/users", makeHTTPHandleFunc(s.handleUsers))
	router.HandleFunc("/users/{id}", makeHTTPHandleFunc(s.handleUserById))


	log.Println("JSON API server running on port:" , s.listenAddress)
	
	http.ListenAndServe(s.listenAddress, router)
}

func (s *APIServer) handleUsers(w http.ResponseWriter,r *http.Request) error {
	if r.Method == "GET"{
		return s.handleGetUsers(w,r)
	}
	if r.Method == "POST"{
		return s.handleCreateUser(w,r)
	}
	
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetUsers(w http.ResponseWriter,r *http.Request) error {
	users, err := s.store.GetUsers()

	if err!= nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, users)
}

func (s *APIServer) handleUserById(w http.ResponseWriter,r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		user, err := s.store.GetUserById(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, user)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteUser(w, r)
	}

	if r.Method == "PUT"{
		return s.handleUpdateUser(w,r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter,r *http.Request) error {
	createAccReq := new(CreateUserRequest)
	
	if err := json.NewDecoder(r.Body).Decode(createAccReq); err != nil {
		return err
	}
	
	user := NewUser(createAccReq.FirstName, createAccReq.LastName)
	
	if err := s.store.CreateUser(user); err!= nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleUpdateUser(w http.ResponseWriter,r *http.Request) error {
	updateUserReq := new(UpdateUserRequest)
	
	if err := json.NewDecoder(r.Body).Decode(updateUserReq); err != nil {
		return err
	}
	
	user := NewUser(updateUserReq.FirstName, updateUserReq.LastName)
	
	if err := s.store.UpdateUser(user); err!= nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleDeleteUser(w http.ResponseWriter,r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteUser(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

type apiFunc func (http.ResponseWriter, *http.Request) error 

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
