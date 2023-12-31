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

func (s *ApiRouter) handleLoginGet(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "auth/login.html")
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

func (s *ApiRouter) handleLoginPost(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.handleError(w, r, err)
		return
	}

	user, err := s.store.GetUserByEmail(req.Email)
	if err != nil {
		tmpl, err := template.ParseFS(templates.Templates, "ui/basicError.html")
		if err != nil {
			s.handleError(w, r, err)
			return
		}
		errorMessage := "Email not found."
		err = tmpl.Execute(w, errorMessage)
		if err != nil {
			s.handleError(w, r, err)
			return
		}
		return
	}

	if !user.ValidPassword(req.Password) {
		tmpl, err := template.ParseFS(templates.Templates, "ui/basicError.html")
		if err != nil {
			s.handleError(w, r, err)
			return
		}
		errorMessage := "Invalid password."
		err = tmpl.Execute(w, errorMessage)
		if err != nil {
			s.handleError(w, r, err)
			return
		}
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

func (s *ApiRouter) handleRegisterGet(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "auth/register.html")
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

func (s *ApiRouter) handleRegisterPost(w http.ResponseWriter, r *http.Request) {
	createAccReq := new(types.RegisterRequest)

	if err := json.NewDecoder(r.Body).Decode(createAccReq); err != nil {
		s.handleError(w, r, err)
		return
	}
	if createAccReq.Password != createAccReq.ConfirmPassword {
		tmpl, err := template.ParseFS(templates.Templates, "ui/basicError.html")
		if err != nil {
			s.handleError(w, r, err)
			return
		}
		errorMessage := "Passwords don't match."
		err = tmpl.Execute(w, errorMessage)
		if err != nil {
			s.handleError(w, r, err)
			return
		}
		return
	}
	existingUser, _ := s.store.GetUserByEmail(createAccReq.Email)
	if existingUser != nil {
		tmpl, err := template.ParseFS(templates.Templates, "ui/registerError.html")
		if err != nil {
			s.handleError(w, r, err)
			return
		}
		errorMessage := "Email Already Exists."
		err = tmpl.Execute(w, errorMessage)
		if err != nil {
			s.handleError(w, r, err)
			return
		}
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
			Name:     "nickname",
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
				cookie, err := r.Cookie("email")
				if err != nil {
					return
				}

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
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (s *ApiRouter) withRoleAuth(d database.Methods, requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secret := os.Getenv("JWT_SECRET")

			tokenString, err := extractTokenFromRequest(r)
			if err != nil {
				s.HandleGetUserRow(w, r)
				s.sendRoleAlert(w, r)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil {
				s.HandleGetUserRow(w, r)
				s.sendRoleAlert(w, r)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				email := claims["email"].(string)
				user, err := d.GetUserByEmail(email)
				if err != nil {
					s.HandleGetUserRow(w, r)
					s.sendRoleAlert(w, r)
					return
				}

				if user.Role.Name != requiredRole {
					s.HandleGetUserRow(w, r)
					s.sendRoleAlert(w, r)
					return
				}
			} else {
				s.HandleGetUserRow(w, r)
				s.sendRoleAlert(w, r)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func createJWT(user *user.User) (string, error) {

	expirationTime := time.Now().Add(3600 * time.Second)

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

	tokenString, err := extractTokenFromRequest(r)
	if err != nil {
		return err
	}

	claims := &types.LoginResponse{}

	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return err
	}

	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now())

	newToken, err := createJWT(user)
	if err != nil {
		return err
	}

	domain := os.Getenv("DOMAIN")
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newToken,
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
	tmpl, err := template.ParseFS(templates.Templates, "ui/base.html", "ui/navbar.html", "ui/page403.html")
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *ApiRouter) sendRoleAlert(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(templates.Templates, "ui/roleAlert.html")
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	err = tmpl.Execute(w, r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

}
