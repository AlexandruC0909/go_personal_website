package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	database "go_api/database"
	"go_api/types"
	user "go_api/types"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler interface {
	handleLogin(w http.ResponseWriter, r *http.Request) error
	handleRegister(w http.ResponseWriter, r *http.Request) error
}

func (s *ApiRouter) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		templatesDir := os.Getenv("TEMPLATES_DIR")
		if templatesDir == "" {
			fmt.Println("TEMPLATES_DIR environment variable is not set.")
			return
		}

		tmplPath := fmt.Sprintf("%s/auth/login.html", templatesDir)
		if err != nil {
			return err
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			return err
		}

		return nil
	}

	if r.Method == "POST" {

		var req types.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return err
		}

		user, err := s.store.GetUserByEmail(req.Email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)

			return nil
		}

		if !user.ValidPassword(req.Password) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)

			return nil
		}

		token, err := createJWT(user)
		if err != nil {
			return err
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			Domain:   "localhost",
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "email",
			Value:    user.Email,
			HttpOnly: true,
			Path:     "/",
			Domain:   "localhost",
		})
		w.Header().Set("HX-Redirect", "/home")

	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *ApiRouter) handleRegister(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		relativePath := "templates/auth/register.html"

		// Get the absolute path.
		absolutePath, err := filepath.Abs(relativePath)
		if err != nil {
			fmt.Println("Error:", err)
		}
		tmpl, err := template.ParseFiles(absolutePath)
		if err != nil {
			return err
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			return err
		}

		return nil
	}

	if r.Method == "POST" {
		createAccReq := new(types.RegisterRequest)

		if err := json.NewDecoder(r.Body).Decode(createAccReq); err != nil {
			return err
		}
		existingUser, _ := s.store.GetUserByEmail(createAccReq.Email)
		if existingUser != nil {
			return WriteJSON(w, http.StatusForbidden, ApiError{Error: createAccReq.Email + " already used"})
		}
		user, err := types.NewUser(createAccReq.FirstName, createAccReq.LastName, createAccReq.Email, createAccReq.Password)
		if err != nil {
			return err
		}
		if err := s.store.CreateUser(user); err != nil {
			return err
		}
		token, err := createJWT(user)
		if err != nil {
			return err
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			Domain:   "localhost",
		})
		w.Header().Set("HX-Redirect", "/home")
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *ApiRouter) handleLogout(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			Domain:   "localhost",
		})
		return nil
	}

	return fmt.Errorf("method not allowed %s", r.Method)

}

func withJWTAuth(handlerFunc http.HandlerFunc, s database.Methods) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("JWT_SECRET")

		fmt.Println("calling JWT auth middleware")
		tokenString, err := extractTokenFromRequest(r)
		if err != nil {
			permissionDenied(w)
			return
		}
		claims := &types.LoginResponse{}
		jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(secret), nil
		})
		if err != nil {
			permissionDenied(w)
			return
		}
		if time.Until(claims.ExpiresAt.Time) < 1*time.Second {
			cookie, _ := r.Cookie("email")

			user, _ := s.GetUserByEmail(cookie.Value)
			refreshToken(w, r, user)
		} else {

			token, err := validateJWT(tokenString)
			if err != nil {
				permissionDenied(w)
				return
			}
			if !token.Valid {
				permissionDenied(w)
				return
			}
			if err != nil {
				WriteJSON(w, http.StatusForbidden, ApiError{Error: "invalid token"})
				return
			}

		}

		handlerFunc(w, r)
	}
}

func withRoleAuth(requiredRole string, handlerFunc http.HandlerFunc, s database.Methods) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling role-based auth middleware")
		secret := os.Getenv("JWT_SECRET")

		tokenString, _ := extractTokenFromRequest(r)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			permissionDenied(w)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			email := claims["email"].(string)
			user, err := s.GetUserByEmail(email)
			if err != nil {
				permissionDenied(w)
				return
			}

			if user.Role.Name != requiredRole {
				permissionDenied(w)
				return
			}

			handlerFunc(w, r)
		} else {
			permissionDenied(w)
		}
	}
}

func createJWT(user *user.User) (string, error) {

	expirationTime := time.Now().Add(20 * time.Second)

	claims := &types.LoginResponse{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := &types.LoginResponse{}
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func refreshToken(w http.ResponseWriter, r *http.Request, user *user.User) error {
	secret := os.Getenv("JWT_SECRET")

	tokenString, _ := extractTokenFromRequest(r)
	claims := &types.LoginResponse{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now())

	token, err := createJWT(user)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Domain:   "localhost",
	})
	return nil
}

func extractTokenFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}
