package handlers

import (
	"context"
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

	"go_api/chat"
	"go_api/database"
	"go_api/static"
	"go_api/templates"

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
	hub := chat.NewHub()
	go hub.Run()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, w, r)
	})
	router.Get("/", s.handleHome)
	router.Get("/auth/login", s.handleLoginGet)
	router.Post("/auth/login", s.handleLoginPost)
	router.Get("/auth/logout", s.handleLogout)
	router.Post("/auth/register", s.handleRegisterPost)
	router.Get("/auth/register", s.handleRegisterGet)
	router.Route("/users", func(r chi.Router) {
		r.Use(JWTAuthMiddleware(s.store))
		r.Get("/", s.handleGetUsers)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", s.handleGetUser)
			r.Get("/edit", s.handlgeGetUserEditRow)
			r.Get("/row", s.HandleGetUserRow)
			r.Post("/upload", s.handleUploadUserImages)
			r.With(s.withRoleAuth(s.store, "admin")).Put("/", s.handleEditUser)
			r.With(s.withRoleAuth(s.store, "admin")).Delete("/", s.handleDeleteUser)
		})
	})

	router.Route("/chat", func(r chi.Router) {
		r.Get("/", s.handleChat)
		r.Post("/login", s.handleChatLogin)
	})

	router.Route("/workspace", func(r chi.Router) {
		r.Get("/", s.handleGetCards)
		r.Post("/reorder", s.handleReorderCards)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", s.handleGetCard)
			r.Get("/edit", s.handleEditCard)
			r.Get("/delete", s.handleDeleteCard)
		})
	})

	router.Route("/posts", func(r chi.Router) {
		r.Use(JWTAuthMiddleware(s.store))
		r.Use(PaginationMiddleware)
		r.Get("/", s.handleGetPosts)
	})

	router.Handle("/static/css/", http.FileServer(http.FS(static.CssFiles)))
	router.Handle("/static/js/", http.FileServer(http.FS(static.JsFiles)))

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
		fileServer := http.FileServer(root)
		fs := http.StripPrefix(pathPrefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if strings.HasSuffix(path, ".css") {
				w.Header().Set("Content-Type", "text/css")
			} else if strings.HasSuffix(path, ".js") {
				w.Header().Set("Content-Type", "application/javascript")
			} else if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
				w.Header().Set("Content-Type", "image/jpeg")
			} else if strings.HasSuffix(path, ".png") {
				w.Header().Set("Content-Type", "image/png")
			} else if strings.HasSuffix(path, ".gif") {
				w.Header().Set("Content-Type", "image/gif")
			}
			fileServer.ServeHTTP(w, r)
		}))
		fs.ServeHTTP(w, r)
	})
}

func PaginationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling pagination middleware")
		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		offset := (page - 1) * 5
		limit := 5

		pagination2 := map[string]int{
			"page":   page,
			"limit":  limit,
			"offset": offset,
		}

		ctx := context.WithValue(r.Context(), "pagination", pagination2)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "ui/page404.html")
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
