package handlers

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	database "go_api/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "https://localhost", "http://http://87.106.122.212", "https://http://87.106.122.212"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(middleware.Logger)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "/static"))
	FileServer(router, "/static", filesDir)
	router.NotFound(makeHTTPHandleFunc(s.handleNotFound))
	flag.Parse()

	router.HandleFunc("/", makeHTTPHandleFunc(s.handleHome))
	router.HandleFunc("/auth/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/auth/logout", makeHTTPHandleFunc(s.handleLogout))

	router.HandleFunc("/auth/register", makeHTTPHandleFunc(s.handleRegister))
	router.HandleFunc("/users", withJWTAuth(makeHTTPHandleFunc(s.handleGetUsers), s.store))

	router.HandleFunc("/users/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleUserById), s.store))
	router.HandleFunc("/users/{id}/upload", withJWTAuth(makeHTTPHandleFunc(s.UploadImages), s.store))
	router.HandleFunc("/posts", withJWTAuth(makeHTTPHandleFunc(s.handleGetPosts), s.store))
	log.Println("JSON API server running on port:", s.listenAddress)

	http.ListenAndServe(s.listenAddress, router)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
func (s *ApiRouter) handleHome(w http.ResponseWriter, r *http.Request) error {

	templatesDir := os.Getenv("TEMPLATES_DIR")
	if templatesDir == "" {
		fmt.Println("TEMPLATES_DIR environment variable is not set.")
	}

	tmplPathBase := fmt.Sprintf("%s/ui/base.html", templatesDir)
	tmplPathContent := fmt.Sprintf("%s/ui/home.html", templatesDir)

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
func (s *ApiRouter) handleNotFound(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles("templates/ui/page404.html")
	if err != nil {
		return err
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
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
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
