package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	database "go_api/database"
	"go_api/types"
	user "go_api/types"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler interface {
	handleLogin(w http.ResponseWriter, r *http.Request) error
	handleRefresh(w http.ResponseWriter, r *http.Request) error
	handleRegister(w http.ResponseWriter, r *http.Request) error
}

func (s *ApiRouter) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("../templates/auth/login.html")
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
			// Send a JSON response indicating user not found
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)

			return nil
		}

		if !user.ValidPassword(req.Password) {
			// Send a JSON response indicating user not found
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)

			return nil
		}

		token, err := createJWT(user)
		if err != nil {
			return err
		}
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("HX-Redirect", "/users/"+strconv.Itoa(user.ID))
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    token,
			Expires:  time.Now().Add(60 * time.Minute),
			HttpOnly: true,
			Path:     "/",
			Domain:   "localhost", // Set to the appropriate domain for your environment
		})

	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *ApiRouter) handleRefresh(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	user, err := s.store.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}

	token, err := refreshToken(w, r, user)
	if err != nil {
		return err
	}

	resp := types.LoginResponse{
		Token: token,
		Email: user.Email,
	}

	return WriteJSON(w, http.StatusOK, resp)
}

func (s *ApiRouter) handleRegister(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("../templates/auth/register.html")
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

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    token,
			Expires:  time.Now().Add(60 * time.Minute),
			HttpOnly: true,
			Path:     "/",
			Domain:   "localhost", // Set to the appropriate domain for your environment
		})
		http.Redirect(w, r, "/users/"+strconv.Itoa(user.ID), http.StatusMovedPermanently)
	}

	return fmt.Errorf("method not allowed %s", r.Method)

}

func createJWT(user *user.User) (string, error) {

	expirationTime := time.Now().Add(60 * time.Minute)

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

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
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

func refreshToken(w http.ResponseWriter, r *http.Request, user *user.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	tokenString := r.Header.Get("Authorization")
	splitToken := strings.Split(tokenString, "Bearer ")
	tokenString = splitToken[1]

	claims := &types.LoginResponse{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return tokenString, err
	}
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now())

	return createJWT(user)
}

func withJWTAuth(handlerFunc http.HandlerFunc, s database.Methods) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")
		tokenString, _ := extractTokenFromRequest(r)

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

			// Check if the user's role matches the required role
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

func extractTokenFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
