package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

	if !user.ValidPassword(req.Password) {
		return fmt.Errorf("not authenticated")
	}

	token, err := createJWT(user)
	if err != nil {
		return err
	}

	resp := types.LoginResponse{
		Token: token,
		Email: user.Email,
	}

	return WriteJSON(w, http.StatusOK, resp)
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
	createAccReq := new(types.RegisterRequest)

	if err := json.NewDecoder(r.Body).Decode(createAccReq); err != nil {
		return err
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

	resp := types.LoginResponse{
		Token: token,
		Email: user.Email,
	}
	return WriteJSON(w, http.StatusOK, resp)
}

func createJWT(user *user.User) (string, error) {

	expirationTime := time.Now().Add(1 * time.Minute)

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

func withJWTAuth(handlerFunc http.HandlerFunc, s database.Methods) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]

		token, err := validateJWT(reqToken)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}
		userID, err := getID(r)
		if err != nil {
			permissionDenied(w)
			return
		}
		user, err := s.GetUser(userID)
		if err != nil {

			return
		}

		//claims := token.Claims.(jwt.MapClaims)
		if user.ID != 3 {
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

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := &types.LoginResponse{}
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
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
	return createJWT(user)
}
