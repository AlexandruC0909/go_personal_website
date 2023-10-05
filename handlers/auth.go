package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	database "go_api/database"
	"go_api/types"
	user "go_api/types"

	templates "go_api/templates"

	"github.com/golang-jwt/jwt/v5"
)

func (s *ApiRouter) handleLoginGET(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "auth/login.html")
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

func (s *ApiRouter) handleLoginPOST(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.handleError(w, r, err)
		return
	}

	user, err := s.store.GetUserByEmail(req.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !user.ValidPassword(req.Password) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	token, err := createJWT(user)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
	domain := os.Getenv("DOMAIN")
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Domain:   domain,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "email",
		Value:    user.Email,
		HttpOnly: true,
		Path:     "/",
		Domain:   domain,
	})
	w.Header().Set("HX-Redirect", "/")
}

func (s *ApiRouter) handleRegisterGET(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "auth/register.html")
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

func (s *ApiRouter) handleRegisterPOST(w http.ResponseWriter, r *http.Request) {
	createAccReq := new(types.RegisterRequest)

	if err := json.NewDecoder(r.Body).Decode(createAccReq); err != nil {
		s.handleError(w, r, err)
		return
	}
	existingUser, _ := s.store.GetUserByEmail(createAccReq.Email)
	if existingUser != nil {
		WriteJSON(w, http.StatusForbidden, ApiError{Error: createAccReq.Email + " already used"})
		return
	}
	user, err := types.NewUser(createAccReq.FirstName, createAccReq.LastName, createAccReq.Email, createAccReq.Password)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
	if err := s.store.CreateUser(user); err != nil {
		s.handleError(w, r, err)
		return
	}
	token, err := createJWT(user)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
	domain := os.Getenv("DOMAIN")
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Domain:   domain,
	})
	w.Header().Set("HX-Redirect", "/")
}

func (s *ApiRouter) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		domain := os.Getenv("DOMAIN")
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			Domain:   domain,
			MaxAge:   -1,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "email",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			Domain:   domain,
			MaxAge:   -1,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		s.handleMethodNotAllowed(w, r)
	}
}

func JWTAuthMiddleware(s database.Methods) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Your middleware logic here
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
				// Handle other cases as needed
			}
			next.ServeHTTP(w, r)
		})
	}
}

func withRoleAuth(requiredRole string, handlerFunc http.HandlerFunc, s database.Methods) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	domain := os.Getenv("DOMAIN")
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Domain:   domain,
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

func permissionDenied(w http.ResponseWriter) error {
	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/page403.html")
	if err != nil {
		return nil
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return nil
	}

	return nil
}
