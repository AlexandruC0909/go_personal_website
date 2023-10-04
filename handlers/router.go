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
	templates "go_api/templates"

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

type ApiFunc func(http.ResponseWriter, *http.Request) error

func NewAPIServer(listenAddress string, store database.Methods) *ApiRouter {
	return &ApiRouter{
		listenAddress: listenAddress,
		store:         store,
	}
}

func (s *ApiRouter) Run() {

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "https://localhost", "http://87.106.122.212", "https://87.106.122.212"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(middleware.Logger)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(filesDir)))

	router.NotFound(s.handleNotFound)

	flag.Parse()
	router.Get("/", s.handleHome)
	router.Get("/auth/login", s.handleLoginGET)
	router.Post("/auth/login", s.handleLoginPOST)
	router.Post("/auth/logout", s.handleLogout)
	router.Post("/auth/register", s.handleRegisterPOST)
	router.Get("/auth/register", s.handleRegisterGET)
	router.Route("/users", func(r chi.Router) {
		r.Use(JWTAuthMiddleware(s.store))
		r.Get("/", s.handleGetUsers)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", s.handleUserByIdGET)
			r.Delete("/", s.handleUserByIdDELETE)
			r.Put("/", s.handleUserByIdPUT)
			r.Post("/upload", s.UploadImages)
		})
	})

	router.Get("/posts", s.handleGetPosts)
	log.Println("JSON API server running on port:", s.listenAddress)

	http.ListenAndServe(s.listenAddress, router)
}

func (s *ApiRouter) handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "ui/home.html")
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

func (s *ApiRouter) handleNotFound(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/page404.html")
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

func (s *ApiRouter) handleError(w http.ResponseWriter, r *http.Request, err error) {
	WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
}
func (s *ApiRouter) handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	s.handleError(w, r, fmt.Errorf("method not allowed %s", r.Method))
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func getID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
